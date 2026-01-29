package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	colorLight     = hexToRGBA(0xeeeed2)
	colorDark      = hexToRGBA(0x769656)
	colorHighlight = hexToRGBA(0xbaca44)
	colorWhite     = hexToRGBA(0xffffff)
	colorBlack     = hexToRGBA(0x000000)
)

func hexToRGBA(hex uint32) color.RGBA {
	return color.RGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 255,
	}
}

func createSquare(id int) fyne.CanvasObject {
	row := id / 8
	col := id % 8
	bg := colorLight

	if (row+col)%2 != 0 {
		bg = colorDark
	}

	rect := canvas.NewRectangle(bg)
	rect.SetMinSize(fyne.NewSize(60, 60))

	button := widget.NewButton("", func() {
		fmt.Printf("TODO PRINT: clicked square %d\n", id)
	})
	button.Importance = widget.LowImportance

	return container.NewStack(button, rect)
}

func createBoardGrid(w fyne.Window) {
	grid := container.NewGridWithColumns(8)
	for i := range 64 {
		square := createSquare(i)
		grid.Add(square)
	}
	board := container.New(layout.NewCenterLayout(), grid)
	w.SetContent(board)
	w.Resize(fyne.NewSize(600, 600))
}

func CreateGameWindow(app fyne.App) fyne.Window {
	w := app.NewWindow("Chess")
	createBoardGrid(w)
	return w
}
