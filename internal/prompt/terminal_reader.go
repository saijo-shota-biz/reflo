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

	if _, err := fmt.Fprintf(tr.Out, "%v\n(Enterで改行 / Ctrl+Dで送信 / Ctrl+Cで終了) > \n", prompt); err != nil {
		return "", err
	}

	cfg := &readline.Config{
		Prompt:              "",
		UniqueEditLine:      false,
		ForceUseInteractive: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			// Ctrl + Dで送信する。
			// この処理がないと > hogehoge^Dで送信できない。
			if r == readline.CharDelete {
				done = true
				return '\n', true
			}
			return r, true
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
			return "", err
		case err == io.EOF:
			return sb.String(), nil
		default:
			return "", err
		}
	}
}

func (tr *TerminalReader) ReadCommand(prompt string) error {
	if _, err := fmt.Fprintf(tr.Out, "%v\n(Enterで続行 / Ctrl+Cで終了) > \n", prompt); err != nil {
		return err
	}

	cfg := &readline.Config{
		Prompt:              "",
		UniqueEditLine:      false,
		ForceUseInteractive: true,
		Stdin:               tr.In,
		Stdout:              tr.Out,
	}

	rl, err := readline.NewEx(cfg)
	if err != nil {
		return fmt.Errorf("readline init: %w", err)
	}
	defer rl.Close()

	_, err = rl.Readline()
	if err != nil {
		if errors.Is(err, readline.ErrInterrupt) {
			return err
		}
		if err == io.EOF {
			return nil
		}
		return err
	}

	return nil
}
