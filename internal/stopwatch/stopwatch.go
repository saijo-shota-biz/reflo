package stopwatch

import "time"

//go:generate mockgen -source=stopwatch.go -destination=../../mock/stopwatch/stopwatch_mock.go -package=mock_stopwatch
type Stopwatch interface {
	Start()
	Stop()
	Time() (time.Time, time.Time)
	Elapsed() time.Duration
}
