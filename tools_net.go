package main

import (
	"bufio"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// PingLine 对齐 Tauri net::PingLine 事件载荷。
type PingLine struct {
	Host string `json:"host"`
	Line string `json:"line"`
}

// PortCheckResult 对齐 Tauri net::PortCheckResult。
type PortCheckResult struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Open      bool   `json:"open"`
	LatencyMs uint64 `json:"latency_ms"`
	Message   string `json:"message"`
}

// pingArgs 按平台拼装 ping 参数。-W 在 macOS 是毫秒、Linux 是秒、Windows 用 -w 毫秒，单位不同须分支。
func pingArgs(host string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"-n", "4", "-w", "2000", host}
	case "darwin":
		return []string{"-c", "4", "-W", "2000", host}
	default:
		return []string{"-c", "4", "-W", "2", host}
	}
}

// PingHost 对齐 Tauri net::ping_host：逐行读系统 ping 的 stdout，通过 `ping:line`
// 事件实时推送，返回 ping 退出码是否成功（主机不可达时 false，但不算 Go 侧错误）。
func (a *App) PingHost(host string) (bool, error) {
	cmd := exec.Command("ping", pingArgs(host)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	if err := cmd.Start(); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		wruntime.EventsEmit(a.ctx, "ping:line", PingLine{Host: host, Line: scanner.Text()})
	}
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckPort 对齐 Tauri net::check_port：3 秒超时内能否完成 TCP 三次握手。
// 地址无法解析时返回 open=false + "无法解析地址"，不报错（与原版一致）。
func (a *App) CheckPort(host string, port int) PortCheckResult {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	latencyMs := uint64(time.Since(start).Milliseconds())
	if err != nil {
		if _, ok := err.(*net.DNSError); ok {
			return PortCheckResult{Host: host, Port: port, Open: false, LatencyMs: 0, Message: "无法解析地址"}
		}
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "lookup") || strings.Contains(msg, "no such host") {
			return PortCheckResult{Host: host, Port: port, Open: false, LatencyMs: 0, Message: "无法解析地址"}
		}
		return PortCheckResult{Host: host, Port: port, Open: false, LatencyMs: latencyMs, Message: err.Error()}
	}
	conn.Close()
	return PortCheckResult{Host: host, Port: port, Open: true, LatencyMs: latencyMs, Message: "连接成功"}
}
