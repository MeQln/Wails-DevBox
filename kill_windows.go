//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// killProcess 在 Windows 上用 taskkill 强制结束进程树。
// Windows 无 SIGTERM 语义，与 Tauri 版的 Windows fallback 行为一致。
func killProcess(pid uint32) error {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", pidStr(pid))
	hideWindow(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// taskkill 失败时输出具体原因（如"拒绝访问"），需解码 OEM 编码为 UTF-8
		msg := strings.TrimSpace(toUTF8(out))
		// 去掉 "错误:" / "ERROR:" 等语言前缀，只保留原因
		if _, after, found := strings.Cut(msg, ": "); found {
			msg = after
		}
		if msg != "" {
			return fmt.Errorf("%s", msg)
		}
		return err
	}
	return nil
}
