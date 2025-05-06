package timer

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const delta = 30 * time.Millisecond

func TestTimer_WaitExpiresNearDuration(t *testing.T) {
	// Arrange
	d := 3 * time.Second
	tm := New(d)
	start := time.Now()

	// Act
	tm.Wait()

	// Assert
	elapsed := time.Since(start)
	require.True(t, d-delta <= elapsed && elapsed <= d+delta)
}

func TestTimer_ZeroDurationReturnsImmediately(t *testing.T) {
	// Arrange
	d := 0 * time.Second
	tm := New(d)
	start := time.Now()

	// Act
	tm.Wait()

	// Assert
	elapsed := time.Since(start)
	require.True(t, elapsed <= delta)
}

func TestTimer_WaitIsIdempotent(t *testing.T) {
	// Arrange
	d := 1 * time.Second
	tm := New(d)

	// Act
	done := make(chan struct{})
	go func() {
		tm.Wait()
		close(done)
	}()

	time.Sleep(d + delta)

	select {
	case <-done:
	default:
		t.Fatalf("first Wait did not finish in %v", d+delta)
	}

	start := time.Now()
	tm.Wait()

	// Assert
	elapsed := time.Since(start)
	require.True(t, elapsed <= delta)
}
