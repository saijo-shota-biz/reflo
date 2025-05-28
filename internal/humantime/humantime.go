package humantime

import (
	"fmt"
	"time"
)

func Span(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d / time.Hour)
	m := int(d % time.Hour / time.Minute)

	switch {
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh%02dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	case m > 0:
		return fmt.Sprintf("%02dm", m)
	default:
		return "0m"
	}
}
