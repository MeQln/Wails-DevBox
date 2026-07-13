//go:build !darwin

package main

import "context"

// startTray 是系统托盘初始化的平台特定实现。
// 非 macOS 平台暂为 stub，后续可添加 Windows 实现。
func startTray(iconData []byte, tooltip string) {
	// no-op
}

// startTrayActionHandler 非 darwin 平台为 stub，防止 goroutine 泄漏。
func startTrayActionHandler(ctx context.Context) {
	// no-op
}