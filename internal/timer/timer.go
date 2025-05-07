package timer

import (
	"context"
	"time"
)

type Timer struct {
	ch chan struct{}
}

func New(d time.Duration) *Timer {
	t := &Timer{ch: make(chan struct{})}

	if d <= 0 {
		close(t.ch)
		return t
	}

	time.AfterFunc(d, func() { close(t.ch) })

	return t
}

func (t *Timer) Wait(ctx context.Context) error {
	select {
	case <-t.ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
