package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "DevBox · 开发工具箱",
		Width:             1150,
		Height:            850,
		MinWidth:          1150,
		MinHeight:         850,
		HideWindowOnClose: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// 浅色底，对齐 tokens.css 的 --bg #f6f5f3，避免首屏深色闪烁
		BackgroundColour: &options.RGBA{R: 246, G: 245, B: 243, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}