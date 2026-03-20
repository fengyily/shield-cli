package cmd

import (
	"fmt"

	"shield-cli/service"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Shield system service",
	Long: `Remove the Shield CLI system service. This stops the running service and removes
it from automatic startup. Your configuration and credentials are preserved.`,
	Example: `  shield uninstall`,
	RunE:    runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) error {
	if !service.IsInstalled() {
		fmt.Println("Shield service is not installed.")
		return nil
	}

	port := service.GetInstalledPort()

	PrintBanner()

	fmt.Printf("  \033[1;33mUninstalling Shield CLI service...\033[0m\n\n")

	if err := service.Uninstall(); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}

	fmt.Printf("  \033[1;32m✓ Service uninstalled successfully!\033[0m\n\n")
	fmt.Printf("  \033[90mPort %d has been released.\033[0m\n", port)
	fmt.Printf("  \033[90mYour configuration and credentials are preserved.\033[0m\n")
	fmt.Printf("  \033[90mTo reinstall: shield install\033[0m\n\n")

	return nil
}
