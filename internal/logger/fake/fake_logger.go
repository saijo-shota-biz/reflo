package fake

import "github.com/saijo-shota-biz/reflo/internal/logger"

type Logger struct {
}

func (f *Logger) Write(logger.Session) error {
	return nil
}

func (f *Logger) ReadDay() ([]logger.Session, error) {
	return nil, nil
}
