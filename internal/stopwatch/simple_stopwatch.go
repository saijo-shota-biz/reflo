package stopwatch

import (
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/humantime"
	"time"
)

type SimpleStopwatch struct {
	start time.Time
	end   time.Time
}

func NewSimpleStopwatch() Stopwatch {
	return &SimpleStopwatch{}
}

func (s *SimpleStopwatch) Start() {
	s.start = time.Now().UTC()
}

func (s *SimpleStopwatch) Stop() {
	s.end = time.Now().UTC()
	s.print()
}

func (s *SimpleStopwatch) Time() (time.Time, time.Time) {
	return s.start, s.end
}

func (s *SimpleStopwatch) Elapsed() time.Duration {
	if s.end.IsZero() {
		return time.Duration(0)
	}
	return s.end.Sub(s.start)
}

func (s *SimpleStopwatch) print() {
	span := humantime.Span(s.Elapsed())
	fmt.Println("")
	fmt.Printf(
		"🕑 %s (%s - %s)\n",
		span,
		s.start.In(time.Local).Format("15:04"),
		s.end.In(time.Local).Format("15:04"),
	)
	fmt.Println("")
}
