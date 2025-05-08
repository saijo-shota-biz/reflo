package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	defaultFocus = 25 * time.Minute
	defaultBreak = 5 * time.Minute
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "start":
		cmdStart()
	case "end-day":
		cmdEndDay()
	case "help":
		showHelp()
	case "version":

	}
}

func showHelp() {
	fmt.Println(`
reflo - Reflect, Flow, Log

Usage:
  reflo start    # start a focus session
  reflo end-day  # output daily summary
  reflo help     # show this message`,
	)
}

func cmdStart() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for {
		// タスク宣言
		goal, err := readLine("What will you do? > ")
		if err != nil {
			fmt.Println("input error:", err)
			return
		}

		// タイマー開始
		start := time.Now().UTC()
		fmt.Printf("Focusing %v …\n", defaultFocus)
		if err := timer.New(defaultFocus).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")
		end := time.Now().UTC()

		// 振り返り
		retro, err := readLine("What have you done? > ")
		if err != nil {
			fmt.Println("input error:", err)
			return
		}

		// ログ出力
		fmt.Printf("%v ~ %v\n", start.Format("2006-01-02 15:04"), end.Format("15:04"))
		log := logger.NewDefaultJsonLogger()
		err = log.Write(logger.Session{
			StartTime: start,
			EndTime:   end,
			Goal:      goal,
			Retro:     retro,
		})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 休憩
		fmt.Printf("Break %v …\n", defaultBreak)
		if err := timer.New(defaultBreak).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")

		// もう一周？
		ans, err := readLine("One more session? (yes or no) > ")
		if err != nil {
			fmt.Println("input error:", err)
			return
		}
		if !strings.HasPrefix(strings.ToLower(ans), "y") {
			break
		}
		fmt.Println("========================")
	}
}

func cmdEndDay() {
	log := logger.NewDefaultJsonLogger()
	sessions, err := log.ReadDay()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("=== Today's Summary ===")
	for i, session := range sessions {
		fmt.Printf("%d) %s ~ %s —  %s\n   %s\n",
			i+1,
			session.StartTime.In(time.Local).Format("15:04"),
			session.EndTime.In(time.Local).Format("15:04"),
			session.Goal,
			session.Retro)
	}
}

func printTimerError(err error) {
	switch {
	case errors.Is(err, context.Canceled):
		fmt.Println("Canceled by user")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("Timed out")
	default:
		fmt.Println("Error:", err)
	}
}

func readLine(prompt string) (string, error) {
	fmt.Print(prompt)

	// fmt.Scanを使わないのはfmt.Scanではスペースが入れられないため。
	// "Set a goal" -> "Set"しか取得できない
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	// 1. タイプミスで前後にスペースを入れた → Goal / Goal のままログ保存すると、テストや JSON で比較しづらい
	// 2. Scanner は \n 区切りで読み込むが、Windows 環境の \r\n だと \r が残ることがある（"\r"）
	//     → yes\r と比較すると strings.ToLower(oneMore) が "yes\r" になり、HasPrefix("y") が通らない
	// 上記の理由によりstrings.TrimSpaceを行う
	return strings.TrimSpace(scanner.Text()), nil
}
