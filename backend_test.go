package main

import (
	"strings"
	"testing"
)

// 以下单测对齐原 Tauri 项目 src-tauri/src/tools/*.rs 的 #[cfg(test)] 用例，
// 验证 Go 移植与 Rust 版行为等价。仅覆盖不依赖 ctx / 外部状态的纯逻辑。

// URL 编解码：与 JS encodeURIComponent 字面对齐（原 url 单测）
func TestUrlEncodeDecode(t *testing.T) {
	app := NewApp()
	if got := app.UrlEncode("hello world", false); got != "hello%20world" {
		t.Errorf("space → %%20: got %q", got)
	}
	if got := app.UrlEncode("a-_.!~*'()", false); got != "a-_.!~*'()" {
		t.Errorf("unreserved 不编码: got %q", got)
	}
	if got := app.UrlEncode("hello world\nfoo bar", true); got != "hello%20world\nfoo%20bar" {
		t.Errorf("multiline 每行独立: got %q", got)
	}
	if got := app.UrlDecode("%zz", false); got != "%zz" {
		t.Errorf("非法序列返回原文: got %q", got)
	}
	s := "中文 abc!"
	if app.UrlDecode(app.UrlEncode(s, false), false) != s {
		t.Error("roundtrip 失败")
	}
}

// Base64：roundtrip / url safe / 非法返回原文（原 base64 单测）
func TestBase64(t *testing.T) {
	app := NewApp()
	s := "中文 abc!"
	enc := app.Base64Encode(s, false)
	if app.Base64Decode(enc, false) != s {
		t.Error("roundtrip 失败")
	}
	safe := app.Base64Encode("???>>>@@@###", true)
	if strings.ContainsAny(safe, "+/=") {
		t.Errorf("url safe 不应含 +/=: %q", safe)
	}
	if app.Base64Decode(safe, true) != "???>>>@@@###" {
		t.Error("url safe roundtrip 失败")
	}
	if got := app.Base64Decode("!!!not-base64!!!", false); got != "!!!not-base64!!!" {
		t.Errorf("非法返回原文: got %q", got)
	}
	// 含折行空白仍可解码
	wrapped := enc[:4] + "\n" + enc[4:]
	if app.Base64Decode(wrapped, false) != s {
		t.Error("折行空白解码失败")
	}
}

// Hash：已知向量 + bytes 与 text 一致（原 hash 单测）
func TestHash(t *testing.T) {
	app := NewApp()
	r := app.HashText("abc")
	if r.Md5 != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("md5(abc): %s", r.Md5)
	}
	if r.Sha1 != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Errorf("sha1(abc): %s", r.Sha1)
	}
	if r.Sha256 != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Errorf("sha256(abc): %s", r.Sha256)
	}
	if app.HashBytes([]byte("hello")).Md5 != app.HashText("hello").Md5 {
		t.Error("bytes 与 text 哈希不一致")
	}
}

// UUID：v4/v7 版本位、长度、批量唯一、count 钳制（原 uuid 单测）
func TestUuid(t *testing.T) {
	app := NewApp()
	v := app.GenerateUuids(4, 1, false, true)
	if len(v) != 1 || len(v[0]) != 36 || v[0][14] != '4' {
		t.Errorf("v4 hyphenated 36 字符且版本位 4: %v", v)
	}
	if got := app.GenerateUuids(4, 1, false, false)[0]; len(got) != 32 || strings.Contains(got, "-") {
		t.Errorf("no hyphen 32 字符: %q", got)
	}
	if got := app.GenerateUuids(7, 1, false, true)[0]; got[14] != '7' {
		t.Errorf("v7 版本位 7: %c", got[14])
	}
	if len(app.GenerateUuids(4, 0, false, true)) != 1 {
		t.Error("count 0 → 1")
	}
	if len(app.GenerateUuids(4, 9999, false, true)) != 1000 {
		t.Error("count 超大 → 1000")
	}
	batch := app.GenerateUuids(4, 100, false, true)
	seen := map[string]bool{}
	for _, s := range batch {
		if seen[s] {
			t.Fatal("批量生成出现重复")
		}
		seen[s] = true
	}
}

// Password：长度 / 类别覆盖 / 错误分支 / 批量钳制（原 password 单测）
func TestPassword(t *testing.T) {
	app := NewApp()
	all := PasswordOptions{Length: 40, Upper: true, Lower: true, Digit: true, Symbol: true}
	p, err := app.GeneratePasswords(all, 1)
	if err != nil || len(p) != 1 || len([]rune(p[0])) != 40 {
		t.Errorf("长度匹配: %v %v", p, err)
	}
	if !strings.ContainsAny(p[0], "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
		!strings.ContainsAny(p[0], "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(p[0], "0123456789") ||
		!strings.ContainsAny(p[0], pwSymbol) {
		t.Errorf("每个启用类别至少一个字符: %q", p[0])
	}
	if _, err := app.GeneratePasswords(PasswordOptions{Length: 0, Upper: true}, 1); err == nil {
		t.Error("零长度应报错")
	}
	if _, err := app.GeneratePasswords(PasswordOptions{Length: 12}, 1); err == nil {
		t.Error("无字符类别应报错")
	}
	if v, _ := app.GeneratePasswords(PasswordOptions{Length: 16, Upper: true}, 99); len(v) != 10 {
		t.Errorf("批量钳制到 10: %d", len(v))
	}
	excl := PasswordOptions{Length: 200, Upper: true, Lower: true, Digit: true, ExcludeAmbiguous: true}
	pe, _ := app.GeneratePasswords(excl, 1)
	if strings.ContainsAny(pe[0], pwAmbiguous) {
		t.Errorf("剔除易混淆字符失败: %q", pe[0])
	}
}
