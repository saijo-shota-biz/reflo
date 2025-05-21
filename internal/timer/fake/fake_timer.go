package fake

import "context"

type Timer struct {
}

func (t *Timer) Focus(context.Context) error {
	return nil
}

func (t *Timer) Break(context.Context) error {
	return nil
}
