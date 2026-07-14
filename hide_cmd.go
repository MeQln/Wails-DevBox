//go:build !windows

package main

import "os/exec"

// hideWindow 在非 Windows 平台为无操作。
func hideWindow(cmd *exec.Cmd) {}