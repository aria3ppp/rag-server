package clock

import (
	"time"

	"github.com/aria3ppp/rag-server/internal/rag/usecase"
)

type clock struct{}

var _ usecase.Clock = (*clock)(nil)

func NewClock() *clock {
	return &clock{}
}

func (c *clock) TimeNow() time.Time {
	return time.Now().UTC()
}
