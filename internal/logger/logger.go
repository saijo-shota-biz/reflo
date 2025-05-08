package logger

import (
	"os"
	"path/filepath"
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

func getDefaultRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".reflo/logs"
	}
	return filepath.Join(home, ".reflo/logs")
}
