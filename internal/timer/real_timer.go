package timer

import (
	"context"
	"time"
)

type Config struct {
	Focus time.Duration
	Break time.Duration
}

type realTimer struct {
	cfg Config
}

func New(cfg Config) Timer {
	return &realTimer{cfg: cfg}
}

func (t *realTimer) Focus(ctx context.Context) error {
	return run(t.cfg.Focus, ctx)
}

func (t *realTimer) Break(ctx context.Context) error {
	return run(t.cfg.Break, ctx)
}

func run(d time.Duration, ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if d <= 0 {
		return nil
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
