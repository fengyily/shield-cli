package updater

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Restart spawns a detached helper process that waits for the current
// shield-cli to exit, then brings it back up using the freshly installed
// binary. The caller should return from any in-flight HTTP handler and then
// os.Exit so the helper can take over.
func Restart(asService bool, port int) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{
		"__respawn",
		fmt.Sprintf("--pid=%d", os.Getpid()),
		fmt.Sprintf("--port=%d", port),
	}
	if asService {
		args = append(args, "--service")
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	detach(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	// Do not Wait; child is detached and outlives us.
	_ = cmd.Process.Release()
	return nil
}

// WaitForPidExit polls until the given pid no longer exists or timeout fires.
func WaitForPidExit(pid int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		p, err := os.FindProcess(pid)
		if err != nil {
			return
		}
		if err := p.Signal(nilSignal); err != nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}
