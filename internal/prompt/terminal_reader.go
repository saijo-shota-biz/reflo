package prompt

import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"strings"
)

type TerminalReader struct {
	In  io.ReadCloser
	Out io.Writer
}

func NewTerminalReader(in io.ReadCloser, out io.Writer) *TerminalReader {
	return &TerminalReader{In: in, Out: out}
}

func (tr *TerminalReader) ReadLine(prompt string) (string, error) {
	done := false

	fmt.Printf("%v\n(Enterで改行 / Ctrl+Dで送信 / Ctrl+Cで終了) > \n", prompt)

	cfg := &readline.Config{
		Prompt:              "",
		UniqueEditLine:      false,
		ForceUseInteractive: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharDelete { // Ctrl-D で送信したい
				done = true
				return '\n', true // ← 改行を“入力した”ことにする
			}
			return r, true // 通常処理
		},
		Stdin:  tr.In,
		Stdout: tr.Out,
	}
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return "", fmt.Errorf("readline init: %w", err)
	}
	defer rl.Close()

	var sb strings.Builder
	for {
		line, err := rl.Readline()
		switch {
		case err == nil:
			sb.WriteString(line)

			if done {
				return sb.String(), nil
			}

			sb.WriteString("\n")
		case errors.Is(err, readline.ErrInterrupt):
			return "", err // 上位で Ctrl-C 判定
		case err == io.EOF:
			return sb.String(), nil
		default:
			return "", err // 予期せぬエラー
		}
	}
}
