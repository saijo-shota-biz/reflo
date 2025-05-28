package timer

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	focusDur = 300 * time.Millisecond
	breakDur = 100 * time.Millisecond
	delta    = 20 * time.Millisecond // 許容誤差
)

func TestRealTimer_FocusAndBreak(t *testing.T) {
	// Arrange
	config := Config{
		Focus: focusDur,
		Break: breakDur,
	}
	tm := NewRealTimer(config)

	t.Run("Focusしたら設定したFocus時間待つ", func(t *testing.T) {
		// Arrange
		start := time.Now()

		// Act
		err := tm.Focus(context.Background())

		// Assert
		require.NoError(t, err)
		require.InDelta(t, focusDur, time.Since(start), float64(delta))
	})

	t.Run("BreakしたらBreak時間待つ", func(t *testing.T) {
		// Arrange
		start := time.Now()

		// Act
		err := tm.Break(context.Background())

		// Assert
		require.NoError(t, err)
		require.InDelta(t, breakDur, time.Since(start), float64(delta))
	})
}

func TestRealTimer_ZeroDurationReturnsImmediately(t *testing.T) {
	// Arrange
	config := Config{
		Focus: 0 * time.Second,
		Break: 0 * time.Second,
	}
	tm := NewRealTimer(config)

	t.Run("Focus時間が0秒の設定のとき、即時リターンされる", func(t *testing.T) {
		start := time.Now()

		// Act
		err := tm.Focus(context.Background())

		// Assert
		require.NoError(t, err)
		elapsed := time.Since(start)
		require.LessOrEqual(t, elapsed, delta)
	})

	t.Run("Break時間が0秒の設定のとき、即時リターンされる", func(t *testing.T) {
		start := time.Now()

		// Act
		err := tm.Break(context.Background())

		// Assert
		require.NoError(t, err)
		elapsed := time.Since(start)
		require.LessOrEqual(t, elapsed, delta)
	})
}

func TestRealTimer_RepeatFocusWaitsAgain(t *testing.T) {
	t.Run("タイマーを２回スタートしたら、２回とも指定時間待つ", func(t *testing.T) {
		// Arrange
		config := Config{Focus: focusDur}
		tm := NewRealTimer(config)

		// Act And Assert
		start1 := time.Now()
		require.NoError(t, tm.Focus(context.Background()))
		require.InDelta(t, focusDur, time.Since(start1), float64(delta))

		start2 := time.Now()
		require.NoError(t, tm.Focus(context.Background()))
		require.InDelta(t, focusDur, time.Since(start2), float64(delta))
	})
}
func TestRealTimer_WaitCancel(t *testing.T) {
	t.Run("タイマーを途中でキャンセルしたら、ブロックが解除される", func(t *testing.T) {
		// Arrange
		config := Config{Focus: 200 * time.Millisecond}
		tm := NewRealTimer(config)
		start := time.Now()

		// Act
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()
		err := tm.Focus(ctx)

		// Assert
		elapsed := time.Since(start)
		require.Less(t, elapsed, 100*time.Millisecond+20*time.Millisecond)
		require.ErrorIs(t, err, context.Canceled)
	})
}
