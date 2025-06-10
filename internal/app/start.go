package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/prompt"
	"github.com/saijo-shota-biz/reflo/internal/timer"
	"os"
	"os/signal"
	"time"
)

func (app *App) Start() error {
	for {
		goal, canceled, err := readGoal(app.Reader)
		if canceled {
			return nil
		}
		if err != nil {
			return err
		}

		app.Stopwatch.Start()

		err = doFocus(app.Timer, app.Cfg.FocusDuration)
		if err != nil {
			return err
		}

		fmt.Print("\a")

		retro, canceled, err := readRetro(app.Reader)
		if canceled {
			return nil
		}
		if err != nil {
			return err
		}

		app.Stopwatch.Stop()

		start, end := app.Stopwatch.Time()
		err = saveSession(app.Logger, start, end, goal, retro)
		if err != nil {
			return err
		}

		err = doBreak(app.Timer, app.Cfg.BreakDuration)
		if err != nil {
			return err
		}

		fmt.Print("\a")

		printNextSession()
	}
}

func readGoal(reader prompt.Reader) (string, bool, error) {
	goal, err := reader.ReadLine("✏️ このセッションで“完了したいゴール”を入力してください")
	switch {
	case errors.Is(err, readline.ErrInterrupt):
		fmt.Println("⚠️ セッションをキャンセルしました。")
		fmt.Println("次のセッションを始めるには `reflo start` を再実行してください。")
		return "", true, nil
	case err != nil:
		fmt.Printf("💥 目標の読み取りに失敗しました: %v\n", err)
		return "", false, err
	}
	return goal, false, nil
}

func readRetro(reader prompt.Reader) (string, bool, error) {
	retro, err := reader.ReadLine("✏️ セッションを通しての気づき・感想をどうぞ")
	switch {
	case errors.Is(err, readline.ErrInterrupt):
		fmt.Println("⚠️ 振り返りをキャンセルしました。")
		return "", true, nil
	case err != nil:
		fmt.Printf("💥 振り返りの取得に失敗しました: %v\n", err)
		return "", false, err
	}
	return retro, false, nil
}

func doFocus(timer timer.Timer, duration time.Duration) error {
	fmt.Println("")
	fmt.Printf("⏳ 作業開始 %v …\n", duration)
	fmt.Println("")
	focusCtx, stopFocus := signal.NotifyContext(context.Background(), os.Interrupt)
	if err := timer.Focus(focusCtx); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fmt.Println("⚠️ フォーカスタイマーをスキップしました。")
			return nil
		default:
			fmt.Printf("💥 タイマー処理で予期せぬエラーが発生しました: %v\n", err)
			return err
		}
	}
	stopFocus() //Signal channelを開放
	return nil
}

func doBreak(timer timer.Timer, duration time.Duration) error {
	fmt.Println("")
	fmt.Printf("⏳ 休憩開始 %v …\n", duration)
	fmt.Println("")
	breakCtx, stopBreak := signal.NotifyContext(context.Background(), os.Interrupt)
	if err := timer.Break(breakCtx); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fmt.Println("⚠️ 休憩をスキップしました。すぐ次のセッションを開始します。")
			return nil
		default:
			fmt.Printf("💥 タイマー処理で予期せぬエラーが発生しました: %v\n", err)
			return err
		}
	}
	stopBreak()
	return nil
}

func saveSession(log logger.Logger, start time.Time, end time.Time, goal string, retro string) error {
	err := log.Write(logger.Session{
		StartTime: start,
		EndTime:   end,
		Goal:      goal,
		Retro:     retro,
	})
	if err != nil {
		fmt.Printf("💥 ログ保存に失敗しました: %v\n", err)
		return err
	}
	return nil
}

func printNextSession() {
	fmt.Println("")
	fmt.Println("▶️ 次のセッションへ")
	fmt.Println("")
}
