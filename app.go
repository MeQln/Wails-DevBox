package main

import (
	"context"
	"os"
	"path/filepath"
)

// App 是 Wails 绑定的根对象，所有暴露给前端的命令均挂在其上（同 package main，
// 按工具拆分到 tools_*.go）。ctx 在 startup 时保存，供 PingHost 发事件用。
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initTray(ctx)
}

// ReadFile 读取文件原始字节。对应 Tauri plugin-fs 的 readFile。
// Wails 会把 []byte 以 base64 字符串形式与前端交换。
func (a *App) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadTextFile 以 UTF-8 文本读取文件。对应 Tauri plugin-fs 的 readTextFile。
func (a *App) ReadTextFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile 写入字节到指定路径。对应 Tauri plugin-fs 的 writeFile。
// 写入前确保父目录存在，行为与前端「保存二维码图片」场景一致。
func (a *App) WriteFile(path string, data []byte) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}
