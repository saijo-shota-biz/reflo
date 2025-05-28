package testutil

import (
	"os"
	"testing"
)

// SilenceStdout wraps a test body and discards stdout/stderr while it runs.
//
//	t.Run(name, testutil.SilenceStdout(func(t *testing.T) {
//	    ... noisy code ...
//	}))
func SilenceStdout(body func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		null, _ := os.Open(os.DevNull)
		defer null.Close()

		origOut := os.Stdout
		origErr := os.Stderr
		os.Stdout = null
		os.Stderr = null
		defer func() {
			os.Stdout = origOut
			os.Stderr = origErr
		}()

		body(t) // run the real test
	}
}
