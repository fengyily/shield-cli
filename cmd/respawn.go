package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"shield-cli/service"
	"shield-cli/updater"

	"github.com/spf13/cobra"
)

// respawnCmd is invoked by the updater after swapping the binary. It waits
// for the old process to exit, then brings the service (or a foreground
// instance on the given port) back up. Hidden from help because end users
// should never call it directly.
var respawnCmd = &cobra.Command{
	Use:    "__respawn",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		parentPid, _ := cmd.Flags().GetInt("pid")
		port, _ := cmd.Flags().GetInt("port")
		asService, _ := cmd.Flags().GetBool("service")

		if parentPid > 0 {
			updater.WaitForPidExit(parentPid, 30*time.Second)
		}
		// Give the OS a beat to release the old binary handle (Windows).
		time.Sleep(500 * time.Millisecond)

		if asService {
			if err := service.Start(); err != nil {
				return fmt.Errorf("respawn: service start: %w", err)
			}
			return nil
		}

		exe, err := os.Executable()
		if err != nil {
			return err
		}
		args2 := []string{"start"}
		if port > 0 {
			args2 = append(args2, fmt.Sprintf("%d", port))
		}
		c := exec.Command(exe, args2...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Start()
	},
}

func init() {
	respawnCmd.Flags().Int("pid", 0, "parent pid to wait for")
	respawnCmd.Flags().Int("port", 0, "port for foreground mode")
	respawnCmd.Flags().Bool("service", false, "restart as system service")
	rootCmd.AddCommand(respawnCmd)
}
