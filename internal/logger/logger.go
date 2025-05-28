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

//go:generate mockgen -source=logger.go -destination=../../mock/logger/logger_mock.go -package=mock_logger

type Logger interface {
	Write(session Session) error
	ReadDay() ([]Session, error)
}
