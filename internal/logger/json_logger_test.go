package logger

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestJsonLogger_WriteReadRoundTrip(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	log := NewJsonLogger(fmt.Sprintf("%s/logs", dir))

	session1 := Session{
		StartTime: time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 4, 1, 12, 25, 0, 0, time.UTC),
		Goal:      "testGoal1",
		Retro:     "testRetro1",
	}

	session2 := Session{
		StartTime: time.Date(2025, 4, 1, 12, 30, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 4, 1, 12, 55, 0, 0, time.UTC),
		Goal:      "testGoal2",
		Retro:     "testRetro2",
	}

	// Act
	err1 := log.Write(session1)
	require.NoError(t, err1)
	err2 := log.Write(session2)
	require.NoError(t, err2)
	sessions, err3 := log.ReadDay()
	require.NoError(t, err3)

	// Assert
	require.Equal(t, 2, len(sessions))
	require.Equal(t, session1, sessions[0])
	require.Equal(t, session2, sessions[1])
	expected := filepath.Join(dir, "logs", fmt.Sprintf("%s.json", time.Now().Format("2006-01-02")))
	require.FileExists(t, expected)
}

func TestJsonLogger_CreatesDirIfMissing(t *testing.T) {
	// Arrange
	dir := fmt.Sprintf("%s/logs", t.TempDir())
	log := NewJsonLogger(dir)
	before := filepath.Join(dir, fmt.Sprintf("%s.json", time.Now().Format("2006-01-02")))
	require.NoFileExists(t, before)

	// Act
	err := log.Write(Session{})
	require.NoError(t, err)

	// Assert
	after := filepath.Join(dir, fmt.Sprintf("%s.json", time.Now().Format("2006-01-02")))
	require.FileExists(t, after)
}
func TestJsonLogger_ReadDayEmpty(t *testing.T) {
	// Arrange
	dir := fmt.Sprintf("%s/logs", t.TempDir())
	log := NewJsonLogger(dir)
	before := filepath.Join(dir, fmt.Sprintf("%s.json", time.Now().Format("2006-01-02")))
	require.NoFileExists(t, before)

	// Act
	sessions, err := log.ReadDay()

	// Assert
	require.Error(t, err)
	require.Equal(t, 0, len(sessions))
}

func TestJsonLogger_ConcurrentWrites(t *testing.T) {
	// Arrange
	dir := fmt.Sprintf("%s/logs", t.TempDir())
	log := NewJsonLogger(dir)

	const N = 50
	var wg sync.WaitGroup
	wg.Add(N)

	// Act
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			s := Session{
				StartTime: time.Now(),
				EndTime:   time.Now().Add(1 * time.Minute),
				Goal:      fmt.Sprintf("g%d", i),
				Retro:     "done",
			}
			if err := log.Write(s); err != nil {
				t.Errorf("write error: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Assert
	sessions, err := log.ReadDay()
	require.NoError(t, err)
	require.Equal(t, N, len(sessions))
}
