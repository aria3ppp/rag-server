package stacktrace

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
)

type stackTraceHandler struct {
	slog.Handler
	opts Opts
}

var _ slog.Handler = (*stackTraceHandler)(nil)

func NewStackTraceHandler(handler slog.Handler, opts *Opts) *stackTraceHandler {
	if opts == nil {
		opts = &Opts{}
	}
	opts.apply()

	return &stackTraceHandler{Handler: handler, opts: *opts}
}

func (h *stackTraceHandler) Handle(ctx context.Context, record slog.Record) error {
	// Only add stack trace for error level logs
	if record.Level == slog.LevelError {
		// Capture stack trace
		const depth = 32
		var pcs [depth]uintptr
		n := runtime.Callers(h.opts.SkipFrames, pcs[:])
		frames := runtime.CallersFrames(pcs[:n])

		// Convert frames to string array
		var trace []string
		for {
			frame, more := frames.Next()
			trace = append(trace, frame.Function+"\n\t"+frame.File+":"+strconv.Itoa(frame.Line))
			if !more {
				break
			}
		}

		// Add stack trace as an attribute
		record.Add(slog.Any(h.opts.StackTraceKey, trace))
	}

	return h.Handler.Handle(ctx, record)
}
