package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
)

// Application state constants
const (
	StateNew int32 = iota
	StateStarting
	StateRunning
	StateShutdown
)

var (
	ErrAlreadyStarted  = errors.New("app is already started")
	ErrAlreadyShutdown = errors.New("app has been shut down, create a new instance to start again")
	ErrNotRunning      = errors.New("app is not running, call Start() first")
	ErrStartInProgress = errors.New("app start is already in progress")
)

// App provides core application lifecycle management
type App struct {
	logger *slog.Logger
	state  atomic.Int32
	done   chan struct{}

	// Hooks for custom implementation
	onStart func(context.Context) error
}

// New creates a new instance of BaseApp
func New(onStart func(context.Context) error, logger *slog.Logger) *App {
	app := &App{
		logger:  logger,
		onStart: onStart,
		done:    make(chan struct{}),
	}
	app.state.Store(StateNew)
	return app
}

// Start initiates the application
func (app *App) Start(ctx context.Context) error {
	// Try to transition from new to starting state
	if !app.state.CompareAndSwap(StateNew, StateStarting) {
		currentState := app.state.Load()
		switch currentState {
		case StateStarting:
			return ErrStartInProgress
		case StateRunning:
			return ErrAlreadyStarted
		case StateShutdown:
			return ErrAlreadyShutdown
		default:
			return fmt.Errorf("unexpected app state: %d", currentState)
		}
	}

	// Create a new context that will be canceled either by the parent context
	// or when done channel is closed
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Reset state on failure
	defer func() {
		if r := recover(); r != nil {
			app.state.Store(StateNew)
			panic(r)
		}
	}()

	// Transition to running before executing start function
	app.state.Store(StateRunning)

	// Handle shutdown in a separate goroutine
	go func() {
		select {
		case <-app.done:
		case <-ctx.Done():
		}
		app.logger.InfoContext(ctx, "initiating graceful shutdown")
		cancel()
	}()

	// Execute start function
	if err := app.onStart(ctx); err != nil {
		if err != context.Canceled {
			app.state.Store(StateNew)
			return fmt.Errorf("failed to start application: %w", err)
		}
	}

	app.logger.InfoContext(ctx, "shutdown completed")

	return nil
}

// Shutdown gracefully stops the application
func (app *App) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for {
		currentState := app.state.Load()
		switch currentState {
		case StateNew:
			return ErrNotRunning
		case StateStarting:
			return ErrNotRunning
		case StateRunning:
			if app.state.CompareAndSwap(StateRunning, StateShutdown) {
				close(app.done)
				return nil
			}
			continue
		case StateShutdown:
			return nil
		default:
			return fmt.Errorf("unexpected app state: %d", currentState)
		}
	}
}

// GetState returns the current state of the application
func (app *App) GetState() int32 {
	return app.state.Load()
}
