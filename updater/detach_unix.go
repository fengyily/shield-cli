//go:build !windows

package updater

import (
	"os/exec"
	"syscall"
)

var nilSignal = syscall.Signal(0)

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
