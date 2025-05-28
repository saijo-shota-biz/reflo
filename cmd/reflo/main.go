package main

import (
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/cli"
	"os"
)

func main() {

	c, err := cli.New(os.Args)
	if err != nil {
		fmt.Println("CLI init failed", "err", err)
		os.Exit(1)
	}

	// 実行
	if err := c.Run(); err != nil {
		fmt.Println("CLI run failed", "err", err)
		os.Exit(1)
	}
}
