package cli

import (
	"context"
	"errors"
	"github.com/saijo-shota-biz/reflo/internal/app"
	mock_logger "github.com/saijo-shota-biz/reflo/mock/logger"
	mock_prompt "github.com/saijo-shota-biz/reflo/mock/prompt"
	mock_runner "github.com/saijo-shota-biz/reflo/mock/runner"
	mock_timer "github.com/saijo-shota-biz/reflo/mock/timer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCLI_ParseCmd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    cmd
		wantErr bool
	}{
		{"startコマンドが入力されたとき、Startが返る", "start", Start, false},
		{"end-dayコマンドが入力されたとき、EndDayが返る", "end-day", EndDay, false},
		{"helpコマンドが入力されたとき、Helpが返る", "help", Help, false},
		{"start, end-day, help以外の未設定コマンドが入力されたとき、Unknownが返る", "unknown", Unknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCmd(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCLI_New(t *testing.T) {
	t.Run("refloのみで引数が少ないとき、エラーが返る", func(t *testing.T) {
		_, err := New([]string{"reflo"})
		require.Error(t, err)
	})

	t.Run("reflo startが入力されたとき、CLIがstartコマンドで初期化されたインスタンスが返却される", func(t *testing.T) {
		cli, err := New([]string{"reflo", "start"})
		require.NoError(t, err)
		require.Equal(t, Start, cli.command)
		require.NotNil(t, cli.app)
	})

	t.Run("初期化時にオプションが設定されていたら、CLIインスタンスの依存がデフォルトから上書きされて返却される", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeLogger := mock_logger.NewMockLogger(ctrl)
		fakeTimer := mock_timer.NewMockTimer(ctrl)
		fakeReader := mock_prompt.NewMockReader(ctrl)
		cli, err := New(
			[]string{"reflo", "start"},
			WithLogger(fakeLogger),
			WithTimer(fakeTimer),
			WithReader(fakeReader),
		)

		require.NoError(t, err)
		appImpl, ok := cli.app.(*app.App)
		require.True(t, ok)
		require.Same(t, appImpl.Logger, fakeLogger)
		require.Same(t, appImpl.Timer, fakeTimer)
		require.Same(t, appImpl.Reader, fakeReader)
	})
}

func TestCLI_Run(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		cmd       cmd
		setup     func(m *mock_runner.MockRunner)
		expectErr bool
	}{
		{
			name: "start ⇒ Runner.Start が 1 回呼ばれる",
			cmd:  Start,
			setup: func(m *mock_runner.MockRunner) {
				m.EXPECT().
					Start(ctx).
					Return(nil) // 正常終了
			},
		},
		{
			name: "end-day ⇒ Runner.EndDay が 1 回呼ばれる",
			cmd:  EndDay,
			setup: func(m *mock_runner.MockRunner) {
				m.EXPECT().
					EndDay().
					Return()
			},
		},
		{
			name: "help ⇒ Runner.Help が 1 回呼ばれる",
			cmd:  Help,
			setup: func(m *mock_runner.MockRunner) {
				m.EXPECT().
					Help().
					Return(nil)
			},
		},
		{
			name: "Runner.Start がエラーを返すと Run もエラー",
			cmd:  Start,
			setup: func(m *mock_runner.MockRunner) {
				m.EXPECT().
					Start(ctx).
					Return(errors.New("boom"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			runner := mock_runner.NewMockRunner(ctrl)
			tt.setup(runner) // 期待値の仕込み

			c := &CLI{
				app:     runner,
				command: tt.cmd,
			}

			err := c.Run(ctx)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
