package logger

import (
	"time"
)

type Session struct {
	StartTime time.Time
	EndTime   time.Time
	Goal      string
	Retro     string
}

type Logger interface {
	Write(session Session) error
	ReadDay() ([]Session, error)
}
