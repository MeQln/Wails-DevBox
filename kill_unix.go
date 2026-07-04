//go:build !windows

package main

import (
	"os"
	"syscall"
)

// killProcess 在 Unix 上发 SIGTERM（优雅结束）；失败时 fallback 强制 kill，
// 对齐 Tauri 版 kill_with(Signal::Term).unwrap_or(kill()) 的行为。
func killProcess(pid uint32) error {
	p, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	if err := p.Signal(syscall.SIGTERM); err != nil {
		// 优雅结束失败（如权限不足）再尝试强制结束
		return p.Kill()
	}
	return nil
}
