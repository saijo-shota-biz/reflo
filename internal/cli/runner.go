package cli

import "context"

//go:generate mockgen -source=runner.go -destination=../../mock/runner/runner_mock.go -package=mock_runner
type Runner interface {
	Start(ctx context.Context) error
	EndDay()
	Help() error
}
