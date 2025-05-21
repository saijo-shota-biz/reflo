package app

import (
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"io"
	"time"
)

type Config struct {
	DefaultFocus time.Duration
	DefaultBreak time.Duration
	PromptIn     io.ReadCloser
	PromptOut    io.Writer
}

type App struct {
	cfg    Config
	logger logger.Logger
}

func New(cfg Config, l logger.Logger) *App {
	return &App{cfg: cfg, logger: l}
}
