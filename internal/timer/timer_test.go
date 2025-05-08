package timer

import (
	"context"
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
	err := tm.Wait(context.Background())

	// Assert
	require.NoError(t, err)
	elapsed := time.Since(start)
	require.True(t, d-delta <= elapsed && elapsed <= d+delta)
}

func TestTimer_ZeroDurationReturnsImmediately(t *testing.T) {
	// Arrange
	d := 0 * time.Second
	tm := New(d)
	start := time.Now()

	// Act
	err := tm.Wait(context.Background())

	// Assert
	require.NoError(t, err)
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
		err := tm.Wait(context.Background())
		require.NoError(t, err)
		close(done)
	}()

	time.Sleep(d + delta)

	select {
	case <-done:
	default:
		t.Fatalf("first Wait did not finish in %v", d+delta)
	}

	start := time.Now()
	err := tm.Wait(context.Background())

	// Assert
	require.NoError(t, err)
	elapsed := time.Since(start)
	require.True(t, elapsed <= delta)
}

func TestTimer_WaitCancel(t *testing.T) {
	// Arrange
	d := 200 * time.Millisecond
	start := time.Now()

	// Act
	tm := New(d)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	err := tm.Wait(ctx)

	// Assert
	elapsed := time.Since(start)
	require.Less(t, elapsed, 100*time.Millisecond+20*time.Millisecond)
	require.ErrorIs(t, err, context.Canceled)
}
