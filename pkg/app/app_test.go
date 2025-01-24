package app_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/aria3ppp/rag-server/pkg/app"
)

func TestNewApp(t *testing.T) {
	logger := slog.Default()
	onStart := func(ctx context.Context) error { return nil }

	application := app.New(onStart, logger)

	if application == nil {
		t.Fatal("NewApp returned nil")
	}

	if application.GetState() != app.StateNew {
		t.Errorf("expected initial state to be StateNew, got %d", application.GetState())
	}
}

func TestStartSuccess(t *testing.T) {
	started := make(chan struct{})
	logger := slog.Default()
	onStart := func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		return nil
	}

	application := app.New(onStart, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start app in goroutine since it blocks until shutdown
	go func() {
		if err := application.Start(ctx); err != nil {
			t.Errorf("unexpected error on start: %v", err)
		}
	}()

	// Wait for start
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("app failed to start within timeout")
	}

	if state := application.GetState(); state != app.StateRunning {
		t.Errorf("expected state Running, got %d", state)
	}
}

func TestStartFailure(t *testing.T) {
	logger := slog.Default()
	expectedErr := errors.New("start failure")
	onStart := func(ctx context.Context) error {
		return expectedErr
	}

	application := app.New(onStart, logger)

	err := application.Start(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if state := application.GetState(); state != app.StateNew {
		t.Errorf("expected state to reset to New after failure, got %d", state)
	}
}

func TestMultipleStarts(t *testing.T) {
	logger := slog.Default()
	onStart := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}

	application := app.New(onStart, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// First start
	go application.Start(ctx)

	// Wait for app to be in running state
	time.Sleep(100 * time.Millisecond)

	// Second start should fail
	err := application.Start(context.Background())
	if !errors.Is(err, app.ErrAlreadyStarted) {
		t.Errorf("expected ErrAlreadyStarted, got %v", err)
	}
}

func TestShutdown(t *testing.T) {
	logger := slog.Default()
	onStart := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}

	application := app.New(onStart, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start app
	go application.Start(ctx)

	// Wait for app to be running
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	err := application.Shutdown(context.Background())
	if err != nil {
		t.Errorf("unexpected error on shutdown: %v", err)
	}

	if state := application.GetState(); state != app.StateShutdown {
		t.Errorf("expected state Shutdown, got %d", state)
	}
}

func TestShutdownNotRunning(t *testing.T) {
	logger := slog.Default()
	application := app.New(func(ctx context.Context) error { return nil }, logger)

	err := application.Shutdown(context.Background())
	if !errors.Is(err, app.ErrNotRunning) {
		t.Errorf("expected ErrNotRunning, got %v", err)
	}
}

func TestContextCancellation(t *testing.T) {
	logger := slog.Default()
	started := make(chan struct{})
	onStart := func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		return ctx.Err()
	}

	application := app.New(onStart, logger)
	ctx, cancel := context.WithCancel(context.Background())

	// Start app
	go func() {
		if err := application.Start(ctx); err != nil && err != context.Canceled {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// Wait for start
	<-started

	// Cancel context
	cancel()

	// Wait for shutdown to complete
	time.Sleep(100 * time.Millisecond)

	if state := application.GetState(); state != app.StateRunning {
		t.Errorf("expected state Running after context cancellation, got %d", state)
	}
}

func TestStartPanic(t *testing.T) {
	logger := slog.Default()
	onStart := func(ctx context.Context) error {
		panic("test panic")
	}

	application := app.New(onStart, logger)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic was not propagated")
		}

		if state := application.GetState(); state != app.StateNew {
			t.Errorf("expected state to reset to New after panic, got %d", state)
		}
	}()

	_ = application.Start(context.Background())
}

func TestShutdownConcurrent(t *testing.T) {
	logger := slog.Default()
	application := app.New(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}, logger)

	// Start the app
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go application.Start(ctx)

	// Wait for app to be in running state
	time.Sleep(100 * time.Millisecond)

	// Try multiple concurrent shutdowns
	const numGoroutines = 5
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			errChan <- application.Shutdown(context.Background())
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-errChan
		if err != nil && err != app.ErrNotRunning {
			t.Errorf("unexpected error from concurrent shutdown: %v", err)
		}
	}

	// Verify final state
	if state := application.GetState(); state != app.StateShutdown {
		t.Errorf("expected final state Shutdown, got %d", state)
	}
}

func TestStartShutdownRace(t *testing.T) {
	logger := slog.Default()
	startCalled := make(chan struct{})
	application := app.New(func(ctx context.Context) error {
		close(startCalled)
		<-ctx.Done()
		return nil
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the app in a goroutine
	go application.Start(ctx)

	// Wait for start to be called
	<-startCalled

	// Try to start again while shutting down
	go func() {
		err := application.Shutdown(context.Background())
		if err != nil {
			t.Errorf("unexpected shutdown error: %v", err)
		}
	}()

	// Try to start while shutdown is in progress
	err := application.Start(context.Background())
	if err == nil {
		t.Error("expected error when starting during shutdown")
	}
}
