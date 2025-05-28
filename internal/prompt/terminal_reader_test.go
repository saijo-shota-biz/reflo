package prompt

import (
	"github.com/chzyer/readline"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTerminalReader_ReadLine(t *testing.T) {

	t.Run("Hello\nWorldと入力し、Ctrl+Dで送信した時、Hello\nWorldが返る", func(t *testing.T) {
		r, w := io.Pipe()
		var out strings.Builder

		tr := &TerminalReader{In: r, Out: &out}

		go func() {
			w.Write([]byte("Hello"))
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("\n"))
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("World"))
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte{readline.CharDelete})
			w.Close()
		}()

		got, err := tr.ReadLine("dummy prompt")
		require.NoError(t, err)
		require.Equal(t, "Hello\nWorld", got)
	})

	t.Run("Ctrl+Cを入力した時、空文字列とErrInterruptが返る", func(t *testing.T) {
		r, w := io.Pipe()
		var out strings.Builder

		tr := &TerminalReader{In: r, Out: &out}
		go func() {
			w.Write([]byte("Hello"))
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte{readline.CharInterrupt})
			w.Close()
		}()

		got, err := tr.ReadLine("dummy prompt")
		require.ErrorIs(t, err, readline.ErrInterrupt)
		require.Equal(t, "", got)
	})
}
