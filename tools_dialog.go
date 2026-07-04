package main

import (
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// DialogFilter 对齐前端 { name, extensions } 结构，供 OpenDialog/SaveDialog 使用。
type DialogFilter struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}

// toWailsFilters 把前端 filter 转为 Wails runtime.FileFilter。
// extensions: ["png","jpg"] → Pattern: "*.png;*.jpg"。
func toWailsFilters(filters []DialogFilter) []wruntime.FileFilter {
	out := make([]wruntime.FileFilter, 0, len(filters))
	for _, f := range filters {
		var patterns []string
		for _, ext := range f.Extensions {
			ext = strings.TrimSpace(strings.TrimPrefix(ext, "."))
			if ext == "" {
				continue
			}
			patterns = append(patterns, "*."+ext)
		}
		if f.Name == "" && len(patterns) == 0 {
			continue
		}
		out = append(out, wruntime.FileFilter{
			DisplayName: f.Name,
			Pattern:     strings.Join(patterns, ";"),
		})
	}
	return out
}

// OpenDialog 打开文件选择对话框（单选）。取消返回空字符串（不报错），对齐
// Tauri plugin-dialog open 返回 null 的语义（前端封装层把空串转 null）。
func (a *App) OpenDialog(title string, filters []DialogFilter) (string, error) {
	return wruntime.OpenFileDialog(a.ctx, wruntime.OpenDialogOptions{
		Title:   title,
		Filters: toWailsFilters(filters),
	})
}

// SaveDialog 打开保存对话框。defaultPath 作为默认文件名。取消返回空字符串。
func (a *App) SaveDialog(title string, defaultPath string, filters []DialogFilter) (string, error) {
	return wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultPath,
		Filters:         toWailsFilters(filters),
	})
}
