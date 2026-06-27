//go:build !windows

package bridge

import "os/exec"

func hideConsole(*exec.Cmd) {}
