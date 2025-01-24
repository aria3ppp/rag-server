package app

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

func TestStartStateTransitions(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  int32
		expectedError error
	}{
		{
			name:          "Starting state returns ErrStartInProgress",
			initialState:  StateStarting,
			expectedError: ErrStartInProgress,
		},
		{
			name:          "Running state returns ErrAlreadyStarted",
			initialState:  StateRunning,
			expectedError: ErrAlreadyStarted,
		},
		{
			name:          "Shutdown state returns ErrAlreadyShutdown",
			initialState:  StateShutdown,
			expectedError: ErrAlreadyShutdown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.Default()
			app := New(func(ctx context.Context) error { return nil }, logger)

			// Manually set initial state
			app.state.Store(tc.initialState)

			// Attempt to start
			err := app.Start(context.Background())

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestShutdownStateTransitions(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  int32
		expectedError error
	}{
		{
			name:          "New state returns ErrNotRunning",
			initialState:  StateNew,
			expectedError: ErrNotRunning,
		},
		{
			name:          "Starting state returns ErrNotRunning",
			initialState:  StateStarting,
			expectedError: ErrNotRunning,
		},
		{
			name:          "Shutdown state returns nil (already shutdown)",
			initialState:  StateShutdown,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.Default()
			app := New(func(ctx context.Context) error { return nil }, logger)

			// Manually set initial state
			app.state.Store(tc.initialState)

			// Attempt to shutdown
			err := app.Shutdown(context.Background())

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
