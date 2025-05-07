package timer

import (
	"context"
	"sync"
	"time"
)

type Timer struct {
	once sync.Once
	ch   chan struct{}
}

func New(duration time.Duration) *Timer {
	t := &Timer{ch: make(chan struct{})}

	go func() {
		time.Sleep(duration)
		t.once.Do(func() {
			close(t.ch)
		})
	}()

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
