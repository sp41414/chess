package main

import (
	"github.com/sp41414/chess/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := ui.CreateGameWindow(a)
	w.ShowAndRun()
}
