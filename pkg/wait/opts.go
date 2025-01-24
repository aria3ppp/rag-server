package wait

import "time"

const (
	defaultInterval   = time.Second
	defaultMaxRetries = 10
)

type Opts struct {
	Interval   time.Duration
	MaxRetries int
}

func (opts *Opts) apply() {
	if opts.Interval == 0 {
		opts.Interval = defaultInterval
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = defaultMaxRetries
	}
}
