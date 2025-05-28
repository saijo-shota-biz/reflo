package cli

import (
	"context"
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/app"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/prompt"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"os"
	"time"
)

type cmd int

const (
	Start cmd = iota
	EndDay
	Help
	Unknown = -1
)

var cmdStrings = map[string]cmd{
	"start":   Start,
	"end-day": EndDay,
	"help":    Help,
}

func (c cmd) String() string {
	for k, v := range cmdStrings {
		if v == c {
			return k
		}
	}
	return "unknown command"
}

func parseCmd(s string) (cmd, error) {
	if c, ok := cmdStrings[s]; ok {
		return c, nil
	}
	return Unknown, fmt.Errorf("unsupported command: %s", s)
}

type Option func(*deps)
type deps struct {
	runner Runner
	logger logger.Logger
	timer  timer.Timer
	reader prompt.Reader
}

func defaultDeps(cfg app.Config) *deps {
	l := logger.NewDefaultJsonLogger()
	t := timer.NewRealTimer(timer.Config{
		Focus: cfg.FocusDuration,
		Break: cfg.BreakDuration,
	})
	r := prompt.NewTerminalReader(cfg.PromptIn, cfg.PromptOut)
	return &deps{
		logger: l,
		timer:  t,
		reader: r,
	}
}

func WithRunner(r Runner) Option {
	return func(d *deps) {
		d.runner = r
	}
}

func WithLogger(l logger.Logger) Option {
	return func(d *deps) {
		d.logger = l
	}
}

func WithTimer(t timer.Timer) Option {
	return func(d *deps) {
		d.timer = t
	}
}

func WithReader(r prompt.Reader) Option {
	return func(d *deps) {
		d.reader = r
	}
}

type CLI struct {
	app     Runner
	command cmd
}

func New(args []string, opts ...Option) (*CLI, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}

	command, err := parseCmd(args[1])
	if err != nil {
		return nil, err
	}

	// TODO 引数のオプションからConfigを作成
	cfg := app.Config{
		FocusDuration: 25 * time.Minute,
		BreakDuration: 5 * time.Minute,
		PromptIn:      os.Stdin,
		PromptOut:     os.Stdout,
	}

	d := defaultDeps(cfg)
	for _, o := range opts {
		o(d)
	}

	if d.runner == nil {
		d.runner = app.New(cfg, d.logger, d.timer, d.reader)
	}

	return &CLI{app: d.runner, command: command}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	switch c.command {
	case Start:
		if err := c.app.Start(ctx); err != nil {
			return err
		}
	case EndDay:
		c.app.EndDay()
	case Help:
		if err := c.app.Help(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown command: %s", c.command)
	}

	return nil
}
