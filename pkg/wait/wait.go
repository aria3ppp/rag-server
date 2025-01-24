package wait

import (
	"testing"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-cmp/cmp"
)

func Until(t *testing.T, opts *Opts, retryFn func() error) {
	t.Helper()

	if opts == nil {
		opts = &Opts{}
	}
	opts.apply()

	if err := backoff.Retry(
		retryFn,
		backoff.WithMaxRetries(
			backoff.NewConstantBackOff(opts.Interval),
			uint64(opts.MaxRetries),
		),
	); err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}
}
