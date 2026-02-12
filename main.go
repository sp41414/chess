package main

import (
	"context"
	"embed"

	"github.com/sp41414/chess/internal/engine"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type App struct {
	ctx   context.Context
	board *engine.Board
}

//go:embed all:internal/ui/dist
var assets embed.FS

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.board = engine.Init("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title: "Chess",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("error:", err.Error())
	}
}
