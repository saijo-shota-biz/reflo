package main

import (
	"context"
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/cli"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	c, err := cli.New(os.Args)
	if err != nil {
		fmt.Println("CLI init failed", "err", err)
		os.Exit(1)
	}

	// 実行
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := c.Run(ctx); err != nil {
		fmt.Println("CLI run failed", "err", err)
		os.Exit(1)
	}
}
