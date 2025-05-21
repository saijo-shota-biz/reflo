package cli

import (
	"context"
	"github.com/saijo-shota-biz/reflo/internal/app"
	"os"
)

type CLI struct {
	app *app.App
}

func New(app *app.App) *CLI {
	return &CLI{app: app}
}

func (c *CLI) Run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		c.app.Help()
		return nil
	}

	switch os.Args[1] {
	case "start":
		c.app.Start()
	case "end-day":
		c.app.EndDay()
	case "help":
		c.app.Help()
	}

	return nil
}
