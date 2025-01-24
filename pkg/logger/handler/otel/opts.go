package otel

const (
	defaultTraceKey = "trace_id"
	defaultSpanKey  = "span_id"
)

type Opts struct {
	SpanKey  string
	TraceKey string
}

func (opts *Opts) apply() {
	if len(opts.SpanKey) == 0 {
		opts.SpanKey = defaultSpanKey
	}
	if len(opts.TraceKey) == 0 {
		opts.TraceKey = defaultTraceKey
	}
}
