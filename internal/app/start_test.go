package app

import (
	"context"
	"errors"
	"github.com/chzyer/readline"
	"github.com/saijo-shota-biz/reflo/internal/logger"
	"github.com/saijo-shota-biz/reflo/internal/testutil"
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
	goalPrompt := "✏️ このセッションで“完了したいゴール”を入力してください"
	retroPrompt := "✏️ セッションを通しての気づき・感想をどうぞ"

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
			name: "Goal読み取り時にCtrl+Cすると、後続の処理を行わず、正常終了する",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("", readline.ErrInterrupt),
					mt.EXPECT().Focus(gomock.Any()).Times(0),
					mr.EXPECT().ReadLine(retroPrompt).Times(0),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0),
					mt.EXPECT().Break(gomock.Any()).Times(0),
				)

			},
			expectErr: false,
		},
		{
			name: "Goal読み取り時にエラーが発生すると、後続の処理を行わず、エラーを返す",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("", errors.New("boom")),
					mt.EXPECT().Focus(gomock.Any()).Times(0),
					mr.EXPECT().ReadLine(retroPrompt).Times(0),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0),
					mt.EXPECT().Break(gomock.Any()).Times(0),
				)

			},
			expectErr: true,
		},
		{
			name: "Focusタイマーでのキャンセルが発生したときに、Retroへ進む",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(context.Canceled),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(nil),
					mt.EXPECT().Break(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(gomock.Any()).Return("", readline.ErrInterrupt),
				)
			},
			expectErr: false,
		},
		{
			name: "Focusタイマーでのエラーが発生したときに、後続の処理を行わず、エラーを返す",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(errors.New("boom")),
					mr.EXPECT().ReadLine(retroPrompt).Times(0),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0),
					mt.EXPECT().Break(gomock.Any()).Times(0),
					mr.EXPECT().ReadLine(gomock.Any()).Times(0),
				)
			},
			expectErr: true,
		},
		{
			name: "Retro読み取り時にCtrl+Cしたら、後続の処理を行わず、正常終了する",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("", readline.ErrInterrupt),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0),
					mt.EXPECT().Break(gomock.Any()).Times(0),
					mr.EXPECT().ReadLine(gomock.Any()).Times(0),
				)
			},
			expectErr: false,
		},
		{
			name: "Retro読み取り時にエラーが発生したら、後続の処理を行わず、エラーを返す",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("", errors.New("boom")),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Times(0),
					mt.EXPECT().Break(gomock.Any()).Times(0),
					mr.EXPECT().ReadLine(gomock.Any()).Times(0),
				)
			},
			expectErr: true,
		},
		{
			name: "ログ書き込み時にエラーになった時、処理が終了すること",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(errors.New("boom")),
					mt.EXPECT().Break(gomock.Any()).Times(0),
				)
			},
			expectErr: true,
		},
		{
			name: "Breakタイマーでのキャンセルが発生したときに、次のセッションに進む",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(nil),
					mt.EXPECT().Break(gomock.Any()).Return(context.Canceled),
					mr.EXPECT().ReadLine(gomock.Any()).Return("", readline.ErrInterrupt),
				)
			},
			expectErr: false,
		},
		{
			name: "Breakタイマーでのエラーが発生したときに、後続の処理を行わず、エラーを返す",
			setup: func(ml *mock_logger.MockLogger, mt *mock_timer.MockTimer, mr *mock_prompt.MockReader) {
				gomock.InOrder(
					mr.EXPECT().ReadLine(goalPrompt).Return("Write docs", nil),
					mt.EXPECT().Focus(gomock.Any()).Return(nil),
					mr.EXPECT().ReadLine(retroPrompt).Return("Good job", nil),
					ml.EXPECT().Write(gomock.AssignableToTypeOf(logger.Session{})).Return(nil),
					mt.EXPECT().Break(gomock.Any()).Return(errors.New("boom")),
					mr.EXPECT().ReadLine(gomock.Any()).Times(0),
				)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testutil.SilenceStdout(func(t *testing.T) {
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

			err := a.Start()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}))
	}
}
