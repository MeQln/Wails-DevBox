//go:build windows

package main

import "os/exec"

// killProcess 在 Windows 上用 taskkill 强制结束进程树。
// Windows 无 SIGTERM 语义，与 Tauri 版的 Windows fallback 行为一致。
func killProcess(pid uint32) error {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", pidStr(pid))
	hideWindow(cmd)
	return cmd.Run()
}
