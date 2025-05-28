package cli

//go:generate mockgen -source=runner.go -destination=../../mock/runner/runner_mock.go -package=mock_runner
type Runner interface {
	Start() error
	EndDay()
	Help() error
}
