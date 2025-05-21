package fake

import "context"

type Runner struct {
	StartCalled  bool
	EndDayCalled bool
	HelpCalled   bool
	Err          error
}

func (f *Runner) Start(context.Context) error {
	f.StartCalled = true
	if f.Err != nil {
		return f.Err
	} else {
		return nil
	}
}

func (f *Runner) EndDay() {
	f.EndDayCalled = true
}

func (f *Runner) Help() error {
	f.HelpCalled = true
	if f.Err != nil {
		return f.Err
	} else {
		return nil
	}
}
