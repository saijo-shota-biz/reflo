package app

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func (app *App) Help() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "reflo - Pomodoro-style focus logger")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  reflo <command>")
	fmt.Fprintln(w, "COMMANDS:")
	fmt.Fprintln(w, "  start\tStart a focus session")
	fmt.Fprintln(w, "  end-day\tOutput daily summary")
	fmt.Fprintln(w, "  help\tShow this help message")
	fmt.Fprintln(w, "EXAMPLES:")
	fmt.Fprintln(w, "  reflo start")
	fmt.Fprintln(w, "  reflo end-day")
	w.Flush()
}
