package domain

import "time"

type QueryAsyncEvent struct {
	Content   string
	CreatedAt time.Time
	Done      bool
}
