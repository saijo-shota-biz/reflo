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

func TestTerminalReader_ReadCommand(t *testing.T) {
	t.Run("Enterを押した時、nilが返る", func(t *testing.T) {
		r, w := io.Pipe()
		var out strings.Builder

		tr := &TerminalReader{In: r, Out: &out}

		go func() {
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("\n"))
			w.Close()
		}()

		err := tr.ReadCommand("次のセッションに進みますか？")
		require.NoError(t, err)
	})

	t.Run("Ctrl+Cを押した時、ErrInterruptが返る", func(t *testing.T) {
		r, w := io.Pipe()
		var out strings.Builder

		tr := &TerminalReader{In: r, Out: &out}

		go func() {
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte{readline.CharInterrupt})
			w.Close()
		}()

		err := tr.ReadCommand("次のセッションに進みますか？")
		require.ErrorIs(t, err, readline.ErrInterrupt)
	})

	t.Run("EOFの時、nilが返る", func(t *testing.T) {
		r, w := io.Pipe()
		var out strings.Builder

		tr := &TerminalReader{In: r, Out: &out}

		go func() {
			time.Sleep(10 * time.Millisecond)
			w.Close() // EOF
		}()

		err := tr.ReadCommand("次のセッションに進みますか？")
		require.NoError(t, err)
	})
}
