package main

import (
	"crypto/rand"
	"errors"
	"strings"
)

const (
	pwUpper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	pwLower     = "abcdefghijklmnopqrstuvwxyz"
	pwDigit     = "0123456789"
	pwSymbol    = "!@#$%^&*()-_=+[]{};:,.?/"
	pwAmbiguous = "Il1O0o"
	pwMaxCount  = 10
)

// PasswordOptions 对齐 Tauri password::PasswordOptions（serde camelCase）。
type PasswordOptions struct {
	Length           int  `json:"length"`
	Upper            bool `json:"upper"`
	Lower            bool `json:"lower"`
	Digit            bool `json:"digit"`
	Symbol           bool `json:"symbol"`
	ExcludeAmbiguous bool `json:"excludeAmbiguous"`
}

// randInt 返回 [0, n) 的密码学安全随机整数。
func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	v := uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	return int(v % uint64(n))
}

// stripAmbiguous 剔除 Il1O0o 等肉眼易混字符。
func stripAmbiguous(s string) string {
	var b strings.Builder
	for _, r := range s {
		if !strings.ContainsRune(pwAmbiguous, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// buildPools 收集启用类别池（已按需剔除易混淆字符），返回 (各类别池, 合并池)。
func buildPools(opts PasswordOptions) ([]string, string) {
	pick := func(on bool, base string) (string, bool) {
		if !on {
			return "", false
		}
		if opts.ExcludeAmbiguous {
			return stripAmbiguous(base), true
		}
		return base, true
	}
	var cats []string
	var full strings.Builder
	for _, base := range []string{pwUpper, pwLower, pwDigit, pwSymbol} {
		on := false
		switch base {
		case pwUpper:
			_, on = pick(opts.Upper, base)
		case pwLower:
			_, on = pick(opts.Lower, base)
		case pwDigit:
			_, on = pick(opts.Digit, base)
		case pwSymbol:
			_, on = pick(opts.Symbol, base)
		}
		if !on {
			continue
		}
		pool := base
		if opts.ExcludeAmbiguous {
			pool = stripAmbiguous(base)
		}
		cats = append(cats, pool)
		full.WriteString(pool)
	}
	return cats, full.String()
}

// genOne 生成单条密码：每个启用类别先各取一个保证覆盖，剩余从合并池随机填充，再 shuffle。
func genOne(opts PasswordOptions, cats []string, full string) string {
	fullChars := []rune(full)
	out := make([]rune, 0, opts.Length)
	for _, pool := range cats {
		if len(out) >= opts.Length {
			break
		}
		pc := []rune(pool)
		if len(pc) > 0 {
			out = append(out, pc[randInt(len(pc))])
		}
	}
	for len(out) < opts.Length {
		if len(fullChars) > 0 {
			out = append(out, fullChars[randInt(len(fullChars))])
		} else {
			break
		}
	}
	// Fisher-Yates shuffle（crypto/rand）
	for i := len(out) - 1; i > 0; i-- {
		j := randInt(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

// GeneratePasswords 对齐 Tauri password::generate_passwords。
func (a *App) GeneratePasswords(opts PasswordOptions, count int) ([]string, error) {
	if opts.Length == 0 {
		return nil, errors.New("长度必须大于 0")
	}
	cats, full := buildPools(opts)
	if full == "" {
		return nil, errors.New("至少选择一个字符类别")
	}
	n := count
	if n < 1 {
		n = 1
	}
	if n > pwMaxCount {
		n = pwMaxCount
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, genOne(opts, cats, full))
	}
	return out, nil
}
