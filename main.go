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
	return &App{
		board: engine.Init("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// wrappers for wails bindings from go to the frontend

func (a *App) NewGame() {
	a.board = engine.Init("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func (a *App) GetFEN() string {
	return a.board.GetFEN()
}

func (a *App) GetMoves() []engine.Move {
	return a.board.GetMoves()
}

func (a *App) IsInCheck() bool {
	return a.board.IsInCheck()
}

func (a *App) IsCheckmate() bool {
	return a.board.IsCheckmate()
}

func (a *App) IsStalemate() bool {
	return a.board.IsStalemate()
}

func (a *App) IsFiftyMoveRule() bool {
	return a.board.IsFiftyMoveRule()
}

func (a *App) IsInsufficientMaterial() bool {
	return a.board.IsInsufficientMaterial()
}

func (a *App) IsThreefoldRepetition() bool {
	return a.board.IsThreefoldRepetition()
}

// func (a *App) IsDraw() bool {
// 	return a.board.IsDraw()
// }

func (a *App) PlayMove(m engine.Move) engine.Undo {
	return a.board.PlayMove(m)
}

func (a *App) UndoMove(m engine.Move, u engine.Undo) {
	a.board.UndoMove(m, u)
}

func (a *App) GetPieces() map[int]string {
	return a.board.GetPieces()
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title: "Chess",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("error:", err.Error())
	}
}
