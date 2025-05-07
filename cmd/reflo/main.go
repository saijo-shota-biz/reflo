package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"os"
	"os/signal"
	"time"
)

const (
	defaultFocus = 25 * time.Minute
	defaultBreak = 5 * time.Minute
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
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

	// タスク宣言
	fmt.Print("What will you do? > ")
	var goal string
	fmt.Scan(&goal)

	// タイマー開始
	start := time.Now().UTC()
	fmt.Printf("Focusing %v …\n", defaultFocus)
	if err := timer.New(defaultFocus).Wait(ctx); err != nil {
		printTimerError(err)
	}
	fmt.Print("\a")
	end := time.Now().UTC()

	// 振り返り
	fmt.Print("What have you done? > ")
	var retro string
	fmt.Scan(&retro)

	// ログ出力
	fmt.Printf("%v ~ %v\n", start.Format("2006-01-02 15:04"), end.Format("15:04"))
	log := logger.NewDefaultJsonLogger()
	err := log.Write(logger.Session{
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

	// もう一周？
	fmt.Print("One more session? (yes or no) > ")
	var oneMore string
	fmt.Scan(&oneMore)
	if oneMore == "yes" || oneMore == "y" {
		fmt.Println("========================")
		cmdStart()
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
			session.StartTime.Format("15:04"),
			session.EndTime.Format("15:04"),
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
