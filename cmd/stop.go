package cmd

import (
	"fmt"

	"shield-cli/service"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Shield service",
	Long:  "Stop the Shield background service.\nOnly works when Shield is installed as a system service (via \"shield install\").",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !service.IsInstalled() {
			return fmt.Errorf("Shield is not installed as a service\n\n  Use \"shield install\" to install it first")
		}

		PrintBanner()

		status, _ := service.Status()
		if status != "running" {
			fmt.Printf("  \033[90mService is not running\033[0m\n\n")
			return nil
		}

		fmt.Printf("  \033[90mStopping Shield service...\033[0m\n\n")

		if err := service.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		fmt.Printf("  \033[1;32m✓ Service stopped\033[0m\n\n")
		fmt.Printf("  \033[90mCommands:\033[0m\n")
		fmt.Printf("    shield start        Start the service again\n")
		fmt.Printf("    shield uninstall    Remove the service\n\n")
		return nil
	},
}
