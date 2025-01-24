package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type otelHandler struct {
	slog.Handler
	opts Opts
}

var _ slog.Handler = (*otelHandler)(nil)

func NewHandlerWithTraceInfo(handler slog.Handler, opts *Opts) *otelHandler {
	if opts == nil {
		opts = &Opts{}
	}
	opts.apply()

	return &otelHandler{Handler: handler, opts: *opts}
}

func (h *otelHandler) Handle(ctx context.Context, record slog.Record) error {
	// add OpenTelemetry trace information to slog records
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		spanCtx := span.SpanContext()
		record.AddAttrs(
			slog.String(h.opts.TraceKey, spanCtx.TraceID().String()),
			slog.String(h.opts.SpanKey, spanCtx.SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, record)
}
