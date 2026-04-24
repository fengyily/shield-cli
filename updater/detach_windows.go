//go:build windows

package updater

import (
	"os"
	"os/exec"
	"syscall"
)

// nilSignal on Windows: os.FindProcess.Signal(os.Interrupt) fails on a dead
// pid; we use a zero-value equivalent via the os package for WaitForPidExit.
var nilSignal os.Signal = syscall.Signal(0)

const createNewProcessGroup = 0x00000200
const detachedProcess = 0x00000008

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: createNewProcessGroup | detachedProcess,
	}
}
