package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"shield-cli/service"

	"github.com/spf13/cobra"
)

var installPort int

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Shield as a system service",
	Long: `Install Shield CLI as a system service that starts automatically with the system.
The Web UI will be available at http://localhost:<port> after installation.

Platform support:
  macOS:   Installs as a launchd user agent (no sudo required)
  Linux:   Installs as a systemd service (requires sudo)
  Windows: Installs as a Windows service (requires Administrator)`,
	Example: `  shield install              # Install with default port 8181
  shield install --port 8182  # Install with custom port
  shield install --port 9090  # Use port 9090`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().IntVar(&installPort, "port", 8181, "Web UI port")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Validate port
	if installPort < 1 || installPort > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", installPort)
	}

	// Check if already installed
	if service.IsInstalled() {
		existingPort := service.GetInstalledPort()
		return fmt.Errorf("Shield service is already installed (port %d).\nTo reinstall, run: shield uninstall && shield install --port %d", existingPort, installPort)
	}

	// Check port availability
	if !service.CheckPortAvailable(installPort) {
		suggested := service.SuggestPort(installPort)
		if suggested > 0 {
			return fmt.Errorf("port %d is already in use.\nTry an available port: shield install --port %d", installPort, suggested)
		}
		return fmt.Errorf("port %d is already in use. Please specify a different port with --port", installPort)
	}

	// Resolve binary path
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine binary path: %w", err)
	}
	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to resolve binary path: %w", err)
	}
	// Resolve symlinks to get the actual binary
	binaryPath, err = filepath.EvalSymlinks(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	PrintBanner()

	fmt.Printf("  \033[1;33mInstalling Shield CLI as a system service...\033[0m\n\n")
	fmt.Printf("  \033[90m├─\033[0m Binary:     \033[32m%s\033[0m\n", binaryPath)
	fmt.Printf("  \033[90m├─\033[0m Port:       \033[32m%d\033[0m\n", installPort)
	fmt.Printf("  \033[90m└─\033[0m Platform:   \033[32m%s\033[0m\n", platformInfo())
	fmt.Println()

	cfg := service.Config{
		Port:       installPort,
		BinaryPath: binaryPath,
	}

	if err := service.Install(cfg); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("  \033[1;32m✓ Service installed successfully!\033[0m\n\n")
	fmt.Printf("  Shield Web UI is running at:\n\n")
	fmt.Printf("    \033[1;36mhttp://localhost:%d\033[0m\n\n", installPort)
	fmt.Printf("  \033[90mThe service will start automatically on system boot.\033[0m\n")
	fmt.Printf("  \033[90mTo uninstall: shield uninstall\033[0m\n\n")

	return nil
}

func platformInfo() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS (launchd user agent)"
	case "linux":
		return "Linux (systemd service)"
	case "windows":
		return "Windows (Windows Service)"
	default:
		return runtime.GOOS
	}
}
