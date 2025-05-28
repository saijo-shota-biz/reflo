package cli

import (
	"context"
	"errors"
	"github.com/saijo-shota-biz/reflo/internal/app"
	faker "github.com/saijo-shota-biz/reflo/internal/cli/fake"
	fakel "github.com/saijo-shota-biz/reflo/internal/logger/fake"
	fakep "github.com/saijo-shota-biz/reflo/internal/prompt/fake"
	faket "github.com/saijo-shota-biz/reflo/internal/timer/fake"
	"github.com/stretchr/testify/require"
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
		fakeLogger := &fakel.Logger{}
		fakeTimer := &faket.Timer{}
		fakeReader := &fakep.Reader{}
		cli, err := New(
			[]string{"reflo", "start"},
			WithLogger(fakeLogger),
			WithTimer(fakeTimer),
			WithReader(fakeReader),
		)

		require.NoError(t, err)
		appImpl, ok := cli.app.(*app.App)
		require.True(t, ok, "app should be an *app.App")
		require.Same(t, appImpl.Logger, fakeLogger, "logger not overridden")
		require.Same(t, appImpl.Timer, fakeTimer, "timer not overridden")
	})
}

func TestCLI_Run(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		cmd          cmd
		wantStart    bool
		wantEndDay   bool
		wantHelp     bool
		mockErr      error
		expectRunErr bool
	}{
		{"startコマンドの時、runner.Start関数が呼び出される", Start, true, false, false, nil, false},
		{"end-dayコマンドの時、runner.EndDay関数が呼び出される", EndDay, false, true, false, nil, false},
		{"helpコマンドの時、runner.Help関数が呼び出される", Help, false, false, true, nil, false},
		{"runner.Start関数でエラーが発生した時、エラーが返る", Start, true, false, false, errors.New("boom"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &faker.Runner{Err: tt.mockErr}
			cli := &CLI{app: f, command: tt.cmd}

			err := cli.Run(ctx)

			if tt.expectRunErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantStart, f.StartCalled)
			require.Equal(t, tt.wantEndDay, f.EndDayCalled)
			require.Equal(t, tt.wantHelp, f.HelpCalled)
		})
	}
}
