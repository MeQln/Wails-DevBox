package main

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"
)

// base64Engine 返回与 Rust 对齐的 engine：
// url_safe → URL_SAFE_NO_PAD（RawURLEncoding），否则 STANDARD（StdEncoding）。
func base64Engine(urlSafe bool) *base64.Encoding {
	if urlSafe {
		return base64.RawURLEncoding
	}
	return base64.StdEncoding
}

// Base64Encode 对齐 Tauri base64_encode：UTF-8 字节 → base64。
func (a *App) Base64Encode(text string, urlSafe bool) string {
	return base64Engine(urlSafe).EncodeToString([]byte(text))
}

// Base64Decode 对齐 Tauri base64_decode：先剥离所有空白（RFC 4648 §3.1 允许折行），
// 再解码；解码失败或非合法 UTF-8 返回原文（不报错，遵循项目「错误处理反原则」）。
func (a *App) Base64Decode(text string, urlSafe bool) string {
	stripped := strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, text)
	decoded, err := base64Engine(urlSafe).DecodeString(stripped)
	if err != nil {
		return text
	}
	if !utf8.Valid(decoded) {
		return text
	}
	return string(decoded)
}
