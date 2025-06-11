package stopwatch

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const delta = 20 * time.Millisecond // 許容誤差

func TestSimpleStopwatch(t *testing.T) {
	t.Run("StartしてStopしたらTimeで開始と終了時刻が取得できる", func(t *testing.T) {
		// Arrange
		const sleeptime = 100 * time.Millisecond
		sw := NewSimpleStopwatch()

		// Act & Assert
		elapsed := sw.Elapsed()
		require.Equal(t, time.Duration(0), elapsed)

		// Act
		start := sw.Start()

		// Assert
		require.WithinDuration(t, time.Now().UTC(), start, delta)

		time.Sleep(sleeptime)

		// Act & Assert
		elapsed = sw.Elapsed()
		require.InDelta(t, float64(elapsed), float64(time.Since(start)), float64(delta))

		// Act
		end := sw.Stop()

		// Assert
		require.WithinDuration(t, time.Now().UTC(), end, delta)

		// Act & Assert
		elapsed = sw.Elapsed()
		require.InDelta(t, float64(sleeptime), float64(elapsed), float64(delta))
	})
}
