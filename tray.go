package main

import (
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type trayAction int

const (
	actionShow trayAction = 1
	actionHide trayAction = 2
	actionQuit trayAction = 3
)

// actionCh 接收来自平台托盘代码的菜单点击事件，在独立 goroutine 中分发。
var actionCh = make(chan trayAction, 8)

// initTray 初始化系统托盘。
// 在 App.startup 中调用；由平台文件（tray_darwin.go / tray_windows.go）提供具体实现。
func (a *App) initTray(ctx context.Context) {
	iconData, err := os.ReadFile("build/appicon.png")
	if err != nil {
		return
	}

	startTray(iconData, "DevBox · 开发工具箱")

	go func() {
		for act := range actionCh {
			switch act {
			case actionShow:
				runtime.WindowShow(ctx)
			case actionHide:
				runtime.WindowHide(ctx)
			case actionQuit:
				runtime.Quit(ctx)
				os.Exit(0)
			}
		}
	}()
}