package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"io"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "reflo - Pomodoro-style focus logger")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  reflo <command>")
	fmt.Fprintln(w, "COMMANDS:")
	fmt.Fprintln(w, "  start\tStart a focus session")
	fmt.Fprintln(w, "  end-day\tOutput daily summary")
	fmt.Fprintln(w, "  help\tShow this help message")
	fmt.Fprintln(w, "EXAMPLES:")
	fmt.Fprintln(w, "  reflo start")
	fmt.Fprintln(w, "  reflo end-day")
	w.Flush()
}

func cmdStart() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for {
		// --- 計画 ---
		goal, err := readLine("今回のフォーカスで“達成したいゴール”を入力してください")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("セッションを中断しました")
			return
		case err != nil:
			fmt.Println("input error:", err)
			return
		}

		// --- フォーカス ---
		start := time.Now().UTC()
		fmt.Printf("Focusing %v …\n", defaultFocus)
		if err := timer.New(defaultFocus).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")
		end := time.Now().UTC()

		// --- 振り返り ---
		retro, err := readLine("終わってみて、気づき・感想をどうぞ")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("セッションを中断しました")
			return
		case err != nil:
			fmt.Println("input error:", err)
			return
		}

		// --- ログ ---
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

		// --- 休憩 ---
		fmt.Printf("Break %v …\n", defaultBreak)
		if err := timer.New(defaultBreak).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")

		fmt.Println("\n — next session — \n")
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
	done := false

	fmt.Printf("%v\n(Enterで改行 / Ctrl+Dで送信 / Ctrl+Cで終了) > \n", prompt)

	cfg := &readline.Config{
		Prompt:              "",
		UniqueEditLine:      false,
		ForceUseInteractive: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharDelete { // Ctrl-D で送信したい
				done = true
				return '\n', true // ← 改行を“入力した”ことにする
			}
			return r, true // 通常処理
		},
	}
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return "", fmt.Errorf("readline init: %w", err)
	}
	defer rl.Close()

	var sb strings.Builder
	for {
		line, err := rl.Readline()
		switch {
		case err == nil:
			sb.WriteString(line)

			if done {
				return sb.String(), nil
			}

			sb.WriteString("\n")
		case errors.Is(err, readline.ErrInterrupt):
			return "", err // 上位で Ctrl-C 判定
		case err == io.EOF:
			return sb.String(), nil
		default:
			return "", err // 予期せぬエラー
		}
	}
}
