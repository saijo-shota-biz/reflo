package cli

import "context"

type Runner interface {
	Start(ctx context.Context) error
	EndDay()
	Help() error
}
