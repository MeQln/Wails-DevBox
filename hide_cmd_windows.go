//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// hideWindow 在 Windows 上隐藏子进程的命令行窗口。
// 解决 taskkill、ping 等系统命令弹出黑色控制台窗口的问题。
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}