package timer

import (
	"context"
)

type Timer interface {
	Focus(ctx context.Context) error
	Break(ctx context.Context) error
}
