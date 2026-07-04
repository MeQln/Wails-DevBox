package main

import (
	"net/url"
	"strings"
)

// urlEncodeOne 对齐 JS encodeURIComponent：不编码 A-Z a-z 0-9 - _ . ! ~ * ' ( )，
// 其余字符按 UTF-8 字节 percent-encode。net/url 无直接等价函数，故手写。
func urlEncodeOne(s string) string {
	const safe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.!~*'()"
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(safe, r) {
			b.WriteRune(r)
			continue
		}
		// 非 ASCII 或保留字符：按 UTF-8 字节逐字节 %xx
		for _, byt := range []byte(string(r)) {
			b.WriteByte('%')
			b.WriteString(hexUpper2(byt))
		}
	}
	return b.String()
}

// hexUpper2 返回单字节的两位大写十六进制（percent-encoding 规范为大写）。
func hexUpper2(b byte) string {
	const hexd = "0123456789ABCDEF"
	return string([]byte{hexd[b>>4], hexd[b&0x0F]})
}

// urlDecodeOne 对齐 JS decodeURIComponent：失败时返回原文（不报错），与 Rust 版一致。
// 用 PathUnescape（不把 + 转空格，符合 decodeURIComponent 语义）。
func urlDecodeOne(s string) string {
	out, err := url.PathUnescape(s)
	if err != nil {
		return s
	}
	return out
}

// UrlEncode 对齐 Tauri url_encode：multiline 时按 \n 分割，每行独立编码。
func (a *App) UrlEncode(text string, multiline bool) string {
	if multiline {
		lines := strings.Split(text, "\n")
		for i, l := range lines {
			lines[i] = urlEncodeOne(l)
		}
		return strings.Join(lines, "\n")
	}
	return urlEncodeOne(text)
}

// UrlDecode 对齐 Tauri url_decode：multiline 时每行独立解码，失败返回该行原文。
func (a *App) UrlDecode(text string, multiline bool) string {
	if multiline {
		lines := strings.Split(text, "\n")
		for i, l := range lines {
			lines[i] = urlDecodeOne(l)
		}
		return strings.Join(lines, "\n")
	}
	return urlDecodeOne(text)
}
