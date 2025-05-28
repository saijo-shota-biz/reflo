package app

import (
	"context"
	"errors"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	mock_logger "github.com/saijo-shota-biz/reflo/mock/logger"
	mock_prompt "github.com/saijo-shota-biz/reflo/mock/prompt"
	mock_timer "github.com/saijo-shota-biz/reflo/mock/timer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"strings"
	"testing"
	"time"
)

func TestApp_Start(t *testing.T) {
	goalPrompt := "今回のフォーカスで“達成したいゴール”を入力してください"
	retroPrompt := "終わってみて、気づき・感想をどうぞ"

	tests := []struct {
		name      string
		setup     func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader)
		expectErr bool
	}{
		{
			name: "１セッションで依存先でエラーが起こらなかったら、正常終了する",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					// ゴール入力
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					// フォーカスタイマー
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					// 振り返り入力
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					// ログ書き込み
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(nil),
					// 休憩タイマー
					mt.EXPECT().Break(gomock.Any()).Return(nil),
					// ループ２周目ゴール入力
					mr.EXPECT().ReadLine(gomock.Any()).Return("", readline.ErrInterrupt),
				)
			},
			expectErr: false,
		},
		{
			name: "Goal読み取り時にCtrl+Cすると、タイマーが動かず、ログも書き込まれない",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("", readline.ErrInterrupt),
				)
				mt.EXPECT().Focus(gomock.Any()).Times(0)
				mr.EXPECT().ReadLine(retroPrompt).Times(0)
				ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0)
				mt.EXPECT().Break(gomock.Any()).Times(0)
			},
			expectErr: false,
		},
		{
			name: "Retro読み取り時にCtrl+Cすると、タイマーが動かず、ログも書き込まれない",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("", readline.ErrInterrupt),
				)
				ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0)
				mt.EXPECT().Break(gomock.Any()).Times(0)
			},
			expectErr: false,
		},
		{
			name: "Focusタイマーでのキャンセルが発生したときに、Retroへ進むこと",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(errors.New("boom")),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(nil),
					mt.EXPECT().Break(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(gomock.Any()).Return("", readline.ErrInterrupt),
				)
			},
			expectErr: false,
		},
		{
			name: "ログ書き込み時にエラーになった時、処理が終了すること",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(errors.New("boom")),
				)
				mt.EXPECT().Break(gomock.Any()).Times(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モック生成
			mReader := mock_prompt.NewMockReader(ctrl)
			mTimer := mock_timer.NewMockTimer(ctrl)
			mLogger := mock_logger.NewMockLogger(ctrl)

			tt.setup(mLogger, mTimer, mReader)

			cfg := Config{
				FocusDuration: time.Second,
				BreakDuration: time.Second,
				PromptIn:      io.NopCloser(strings.NewReader("")),
				PromptOut:     io.Discard,
			}

			a := New(cfg, mLogger, mTimer, mReader)

			ctx := context.Background()
			err := a.Start(ctx)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
