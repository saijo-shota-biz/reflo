package app

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
	"time"
)

func (app *App) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for {
		// --- 計画 ---
		goal, err := app.readLine("今回のフォーカスで“達成したいゴール”を入力してください")
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
		fmt.Printf("Focusing %v …\n", app.cfg.DefaultFocus)
		if err := timer.New(app.cfg.DefaultFocus).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")
		end := time.Now().UTC()

		// --- 振り返り ---
		retro, err := app.readLine("終わってみて、気づき・感想をどうぞ")
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
		err = app.logger.Write(logger.Session{
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
		fmt.Printf("Break %v …\n", app.cfg.DefaultBreak)
		if err := timer.New(app.cfg.DefaultBreak).Wait(ctx); err != nil {
			printTimerError(err)
		}
		fmt.Print("\a")

		fmt.Println("\n — next session — \n")
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

func (app *App) readLine(prompt string) (string, error) {
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
		Stdin:  app.cfg.PromptIn,
		Stdout: app.cfg.PromptOut,
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
