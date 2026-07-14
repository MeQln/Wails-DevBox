//go:build !windows

package main

// decodeOEM 在非 Windows 平台为无操作。
// Unix/macOS 系统命令输出已为 UTF-8，无需编码转换。
func decodeOEM(raw []byte) []byte {
	return raw
}