package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortEntry 对齐 Tauri port::PortEntry。
type PortEntry struct {
	Port        uint16 `json:"port"`
	Pid         uint32 `json:"pid"`
	ProcessName string `json:"process_name"`
	Address     string `json:"address"`
}

// ListPorts 对齐 Tauri port::list_ports：列出 TCP LISTEN 监听端口及其占用进程。
// gopsutil 在 macOS/Linux 上获取其他用户进程的 PID 可能受限（与 netstat2 同样的系统限制）。
func (a *App) ListPorts() ([]PortEntry, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return nil, err
	}
	out := make([]PortEntry, 0, len(conns))
	for _, c := range conns {
		if c.Status != "LISTEN" {
			continue
		}
		if c.Laddr.Port == 0 {
			continue
		}
		if c.Pid == 0 {
			continue
		}
		name := "unknown"
		if p, err := process.NewProcess(int32(c.Pid)); err == nil {
			if n, err := p.Name(); err == nil && n != "" {
				name = n
			}
		}
		out = append(out, PortEntry{
			Port:        uint16(c.Laddr.Port),
			Pid:         uint32(c.Pid),
			ProcessName: name,
			Address:     c.Laddr.IP,
		})
	}
	return out, nil
}

// KillPort 对齐 Tauri port::kill_port：结束占用端口的进程。
// 先确认进程存在，再调用平台特定的 killProcess。
func (a *App) KillPort(pid uint32) error {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("进程 %d 不存在", pid)
	}
	if _, err := p.Name(); err != nil {
		return fmt.Errorf("进程 %d 不存在", pid)
	}
	if err := killProcess(pid); err != nil {
		return errors.New("结束失败")
	}
	return nil
}

// pidStr 仅供 kill_windows.go 复用，避免重复 strconv。
func pidStr(pid uint32) string { return strconv.FormatUint(uint64(pid), 10) }
