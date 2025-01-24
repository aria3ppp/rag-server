package stacktrace

const (
	defaultStackTraceKey = "stack_trace"
)

type Opts struct {
	StackTraceKey string
	SkipFrames    int
}

func (opts *Opts) apply() {
	if len(opts.StackTraceKey) == 0 {
		opts.StackTraceKey = defaultStackTraceKey
	}
}
