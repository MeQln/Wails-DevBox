package main

import (
	"context"
	"os"
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
	if err != nil || len(iconData) == 0 {
		return
	}

	startTray(iconData, "DevBox · 开发工具箱")
	startTrayActionHandler(ctx)
}