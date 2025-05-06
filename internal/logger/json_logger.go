package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type JsonLogger struct {
	mu   sync.Mutex
	root string
}

func NewJsonLogger(root string) Logger {
	return &JsonLogger{root: root}
}

func NewDefaultJsonLogger() Logger {
	return &JsonLogger{root: getDefaultRoot()}
}

func (l *JsonLogger) Write(session Session) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := os.MkdirAll(l.root, 0755); err != nil {
		return err
	}

	path := filepath.Join(l.root, time.Now().Format("2006-01-02")+".json")

	var list []Session
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &list)
	}

	list = append(list, session)

	b, _ := json.MarshalIndent(list, "", "  ")
	return os.WriteFile(path, b, 0644)
}

func (l *JsonLogger) ReadDay() ([]Session, error) {
	var list []Session

	path := filepath.Join(l.root, time.Now().Format("2006-01-02")+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(b, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
