package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/humantime"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"time"
)

func (app *App) Start(ctx context.Context) error {
	for {
		// 目標入力
		goal, err := app.Reader.ReadLine("今回のフォーカスで“達成したいゴール”を入力してください")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("セッションを中断しました")
			return nil
		case err != nil:
			fmt.Println("input error:", err)
			return nil
		}

		// [時間計測開始]
		start := time.Now().UTC()

		// --- フォーカス ---
		fmt.Printf("Focusing %v …\n", app.Cfg.FocusDuration)
		if err := app.Timer.Focus(ctx); err != nil {
			printTimerError(err)
		}

		// 通知
		fmt.Print("\a")

		// --- 振り返り入力 ---
		retro, err := app.Reader.ReadLine("終わってみて、気づき・感想をどうぞ")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("セッションを中断しました")
			return nil
		case err != nil:
			fmt.Println("input error:", err)
			return nil
		}

		// [時間計測終了]
		end := time.Now().UTC()

		// 作業時間表示
		span := humantime.Span(end.Sub(start))
		fmt.Printf(
			"%s %s - %s\n",
			span,
			start.In(time.Local).Format("15:04"),
			end.In(time.Local).Format("15:04"),
		)

		// --- ログ書き込み ---
		err = app.Logger.Write(logger.Session{
			StartTime: start,
			EndTime:   end,
			Goal:      goal,
			Retro:     retro,
		})
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		// --- 休憩 ---
		fmt.Printf("Break %v …\n", app.Cfg.BreakDuration)
		if err := app.Timer.Break(ctx); err != nil {
			printTimerError(err)
		}

		// 通知
		fmt.Print("\a")

		// 次のセッションへ
		fmt.Println("")
		fmt.Println("— next session — ")
		fmt.Println("")
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
