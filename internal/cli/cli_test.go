package cli

import (
	"context"
	"errors"
	"github.com/saijo-shota-biz/reflo/internal/app"
	faker "github.com/saijo-shota-biz/reflo/internal/cli/fake"
	fakel "github.com/saijo-shota-biz/reflo/internal/logger/fake"
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
		{"start OK", "start", Start, false},
		{"end-day OK", "end-day", EndDay, false},
		{"help OK", "help", Help, false},
		{"unknown command", "unknown", Unknown, true},
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
	t.Run("too few args", func(t *testing.T) {
		_, err := New([]string{"reflo"})
		require.Error(t, err)
	})

	t.Run("happy path", func(t *testing.T) {
		cli, err := New([]string{"reflo", "start"})
		require.NoError(t, err)
		require.Equal(t, Start, cli.command)
		require.NotNil(t, cli.app)
	})

	t.Run("has options", func(t *testing.T) {
		fakeLogger := &fakel.Logger{}
		fakeTimer := &faket.Timer{}
		cli, err := New(
			[]string{"reflo", "start"},
			WithLogger(fakeLogger),
			WithTimer(fakeTimer),
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
		{"calls Start", Start, true, false, false, nil, false},
		{"calls EndDay", EndDay, false, true, false, nil, false},
		{"calls Help", Help, false, false, true, nil, false},
		{"propagates error", Start, true, false, false, errors.New("boom"), true},
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
