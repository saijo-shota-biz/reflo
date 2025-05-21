package app

import (
	"fmt"
	"time"
)

func (app *App) EndDay() {
	sessions, err := app.Logger.ReadDay()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("=== Today's Summary ===")
	for i, session := range sessions {
		fmt.Printf("%d) %s ~ %s —  %s\n   %s\n",
			i+1,
			session.StartTime.In(time.Local).Format("15:04"),
			session.EndTime.In(time.Local).Format("15:04"),
			session.Goal,
			session.Retro)
	}
}
