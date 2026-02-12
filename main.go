package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:internal/ui/dist
var assets embed.FS

func main() {
	err := wails.Run(&options.App{
		Title: "Chess",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
	})

	if err != nil {
		println("error:", err.Error())
	}
}
