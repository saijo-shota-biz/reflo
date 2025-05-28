package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/humantime"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"os"
	"os/signal"
	"time"
)

func (app *App) Start() error {
	for {
		// 目標入力
		goal, err := app.Reader.ReadLine("✏️ このセッションで“完了したいゴール”を入力してください")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("⚠️ セッションをキャンセルしました。")
			fmt.Println("次のセッションを始めるには `reflo start` を再実行してください。")
			return nil
		case err != nil:
			fmt.Printf("💥 目標の読み取りに失敗しました: %v\n", err)
			return nil
		}

		// 時間計測開始
		start := time.Now().UTC()

		// フォーカス
		fmt.Println("")
		fmt.Printf("⏳ 作業開始 %v …\n", app.Cfg.FocusDuration)
		fmt.Println("")
		focusCtx, stopFocus := signal.NotifyContext(context.Background(), os.Interrupt)
		if err := app.Timer.Focus(focusCtx); err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				fmt.Println("⚠️ フォーカスタイマーをスキップしました。")
			default:
				fmt.Printf("💥 タイマー処理で予期せぬエラーが発生しました: %v\n", err)
			}
		}
		stopFocus() //Signal channelを開放

		// 通知
		fmt.Print("\a")

		// 振り返り入力
		retro, err := app.Reader.ReadLine("✏️ セッションを通しての気づき・感想をどうぞ")
		switch {
		case errors.Is(err, readline.ErrInterrupt):
			fmt.Println("⚠️ 振り返りをキャンセルしました。")
		case err != nil:
			fmt.Printf("💥 振り返りの取得に失敗しました: %v\n", err)
		}

		// [時間計測終了]
		end := time.Now().UTC()

		// 作業時間表示
		span := humantime.Span(end.Sub(start))
		fmt.Printf(
			"🕑 作業時間: %s (%s - %s)\n",
			span,
			start.In(time.Local).Format("15:04"),
			end.In(time.Local).Format("15:04"),
		)

		// ログ書き込み
		err = app.Logger.Write(logger.Session{
			StartTime: start,
			EndTime:   end,
			Goal:      goal,
			Retro:     retro,
		})
		if err != nil {
			fmt.Printf("💥 ログ保存に失敗しました: %v\n", err)
			return nil
		}

		// 休憩
		fmt.Println("")
		fmt.Printf("⏳ 休憩開始 %v …\n", app.Cfg.BreakDuration)
		fmt.Println("")
		breakCtx, stopBreak := signal.NotifyContext(context.Background(), os.Interrupt)
		if err := app.Timer.Break(breakCtx); err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				fmt.Println("⚠️ 休憩をスキップしました。すぐ次のセッションを開始します。")
			default:
				fmt.Printf("💥 タイマー処理で予期せぬエラーが発生しました: %v\n", err)
			}
		}
		stopBreak()

		// 通知
		fmt.Print("\a")

		// 次のセッションへ
		fmt.Println("")
		fmt.Println("▶️ 次のセッションへ")
		fmt.Println("")
	}
}
