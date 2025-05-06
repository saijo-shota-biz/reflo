package timer

import (
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

func (t *Timer) Wait() {
	<-t.ch
}
