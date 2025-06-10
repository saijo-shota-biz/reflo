package app

import (
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/notification"
	"github.com/saijo-shota-biz/reflo/internal/prompt"
	"github.com/saijo-shota-biz/reflo/internal/stopwatch"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"io"
	"time"
)

type Config struct {
	FocusDuration time.Duration
	BreakDuration time.Duration
	PromptIn      io.ReadCloser
	PromptOut     io.Writer
}

type App struct {
	Cfg       Config
	Logger    logger.Logger
	Timer     timer.Timer
	Reader    prompt.Reader
	Stopwatch stopwatch.Stopwatch
	Notifier  notification.Notifier
}

func New(cfg Config, l logger.Logger, t timer.Timer, r prompt.Reader, stopwatch stopwatch.Stopwatch, notifier notification.Notifier) *App {
	return &App{Cfg: cfg, Logger: l, Timer: t, Reader: r, Stopwatch: stopwatch, Notifier: notifier}
}
