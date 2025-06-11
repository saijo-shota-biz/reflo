package stopwatch

import (
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/humantime"
	"sync"
	"time"
)

type SimpleStopwatch struct {
	mu    sync.Mutex
	start time.Time
	end   time.Time
}

func NewSimpleStopwatch() Stopwatch {
	return &SimpleStopwatch{}
}

func (s *SimpleStopwatch) Start() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.start = time.Now().UTC()
	s.end = time.Time{}
	return s.start
}

func (s *SimpleStopwatch) Stop() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.end = time.Now().UTC()
	return s.end
}

func (s *SimpleStopwatch) Elapsed() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.elapsed()
}

func (s *SimpleStopwatch) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	span := humantime.Span(s.elapsed())
	return fmt.Sprintf(
		"\n🕑 %s (%s - %s)\n",
		span,
		s.start.In(time.Local).Format("15:04"),
		s.end.In(time.Local).Format("15:04"),
	)
}

func (s *SimpleStopwatch) elapsed() time.Duration {
	if s.start.IsZero() {
		return time.Duration(0)
	}
	if s.end.IsZero() {
		return time.Since(s.start)
	}
	return s.end.Sub(s.start)
}
