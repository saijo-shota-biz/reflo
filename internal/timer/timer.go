package timer

import (
	"context"
)

//go:generate mockgen -source=timer.go -destination=../../mock/timer/timer_mock.go -package=mock_timer

type Timer interface {
	Focus(ctx context.Context) error
	Break(ctx context.Context) error
}
