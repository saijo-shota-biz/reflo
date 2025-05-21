package main

import (
	"context"
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/app"
	"github.com/saijo-shota-biz/reflo/internal/cli"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// デフォルト設定
	cfg := app.Config{
		DefaultFocus: 25 * time.Minute,
		DefaultBreak: 5 * time.Minute,
		PromptIn:     os.Stdin,
		PromptOut:    os.Stdout,
	}

	// 依存関係を整理
	l := logger.NewDefaultJsonLogger()

	// app, cliを構築
	a := app.New(cfg, l)
	c := cli.New(a)

	// 実行
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := c.Run(ctx, os.Args); err != nil {
		fmt.Println("CLI run failed", "err", err)
		os.Exit(1)
	}
}
