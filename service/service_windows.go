package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const windowsServiceName = "ShieldCLI"
const windowsDisplayName = "Shield CLI - Secure Tunnel Connector"

// Install installs shield as a Windows service using sc.exe
func Install(cfg Config) error {
	binPath := fmt.Sprintf(`"%s" start %d`, cfg.BinaryPath, cfg.Port)

	cmd := exec.Command("sc", "create", windowsServiceName,
		"binPath=", binPath,
		"DisplayName=", windowsDisplayName,
		"start=", "auto",
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create service (try running as Administrator): %s (%w)", string(output), err)
	}

	// Set service description
	descCmd := exec.Command("sc", "description", windowsServiceName,
		"Shield CLI Web UI management platform - secure tunnel connector service")
	_ = descCmd.Run()

	// Start the service
	startCmd := exec.Command("sc", "start", windowsServiceName)
	if output, err := startCmd.CombinedOutput(); err != nil {
		fmt.Printf("Warning: service created but failed to start: %s\n", string(output))
	}

	return nil
}

// Uninstall removes the shield Windows service
func Uninstall() error {
	if !IsInstalled() {
		return fmt.Errorf("service is not installed")
	}

	// Stop the service first
	stopCmd := exec.Command("sc", "stop", windowsServiceName)
	_ = stopCmd.Run()

	// Delete the service
	cmd := exec.Command("sc", "delete", windowsServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete service (try running as Administrator): %s (%w)", string(output), err)
	}

	return nil
}

// IsInstalled checks if the Windows service exists
func IsInstalled() bool {
	cmd := exec.Command("sc", "query", windowsServiceName)
	err := cmd.Run()
	return err == nil
}

// Status returns the current service status
func Status() (string, error) {
	if !IsInstalled() {
		return "not installed", nil
	}

	cmd := exec.Command("sc", "query", windowsServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "installed (unknown)", nil
	}
	out := string(output)
	if strings.Contains(out, "RUNNING") {
		return "running", nil
	}
	if strings.Contains(out, "STOPPED") {
		return "installed (stopped)", nil
	}
	return "installed", nil
}

// GetInstalledPort reads the port from the Windows service registry
func GetInstalledPort() int {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\`+windowsServiceName,
		registry.QUERY_VALUE)
	if err != nil {
		return 0
	}
	defer key.Close()

	imagePath, _, err := key.GetStringValue("ImagePath")
	if err != nil {
		return 0
	}

	// Parse port from ImagePath: "path\to\shield.exe" start 8182
	parts := strings.Fields(imagePath)
	for i, p := range parts {
		if p == "start" && i+1 < len(parts) {
			var port int
			fmt.Sscanf(parts[i+1], "%d", &port)
			if port > 0 {
				return port
			}
		}
	}
	return 8181
}
