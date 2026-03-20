package service

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"text/template"
)

const systemdServiceName = "shield-cli"
const systemdServicePath = "/etc/systemd/system/shield-cli.service"

var systemdTemplate = template.Must(template.New("systemd").Parse(`[Unit]
Description=Shield CLI - Secure Tunnel Connector Web UI
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User={{ .User }}
ExecStart={{ .BinaryPath }} start {{ .Port }}
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=shield-cli
WorkingDirectory={{ .HomeDir }}

[Install]
WantedBy=multi-user.target
`))

// Install installs shield as a systemd service on Linux
func Install(cfg Config) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	home, _ := os.UserHomeDir()

	// Create systemd unit file (requires root)
	f, err := os.Create(systemdServicePath)
	if err != nil {
		return fmt.Errorf("failed to create systemd service file (try with sudo): %w", err)
	}
	defer f.Close()

	data := struct {
		User       string
		BinaryPath string
		Port       string
		HomeDir    string
	}{
		User:       currentUser.Username,
		BinaryPath: cfg.BinaryPath,
		Port:       fmt.Sprintf("%d", cfg.Port),
		HomeDir:    home,
	}

	if err := systemdTemplate.Execute(f, data); err != nil {
		os.Remove(systemdServicePath)
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd and enable the service
	if err := runCmd("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}
	if err := runCmd("systemctl", "enable", systemdServiceName); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}
	if err := runCmd("systemctl", "start", systemdServiceName); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Uninstall removes the shield systemd service
func Uninstall() error {
	if _, err := os.Stat(systemdServicePath); os.IsNotExist(err) {
		return fmt.Errorf("service is not installed")
	}

	// Stop and disable the service
	_ = runCmd("systemctl", "stop", systemdServiceName)
	_ = runCmd("systemctl", "disable", systemdServiceName)

	// Remove the service file
	if err := os.Remove(systemdServicePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	_ = runCmd("systemctl", "daemon-reload")

	return nil
}

// IsInstalled checks if the systemd service exists
func IsInstalled() bool {
	_, err := os.Stat(systemdServicePath)
	return err == nil
}

// Status returns the current service status
func Status() (string, error) {
	if !IsInstalled() {
		return "not installed", nil
	}

	cmd := exec.Command("systemctl", "is-active", systemdServiceName)
	output, err := cmd.Output()
	if err != nil {
		return "installed (not running)", nil
	}
	status := string(output)
	if len(status) > 0 && status[len(status)-1] == '\n' {
		status = status[:len(status)-1]
	}
	if status == "active" {
		return "running", nil
	}
	return "installed (" + status + ")", nil
}

// GetInstalledPort reads the port from the installed service file (returns 0 if not found)
func GetInstalledPort() int {
	data, err := os.ReadFile(systemdServicePath)
	if err != nil {
		return 0
	}
	content := string(data)
	var port int
	// Parse ExecStart line to extract port
	idx := 0
	for idx < len(content) {
		lineStart := idx
		lineEnd := lineStart
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		line := content[lineStart:lineEnd]
		if len(line) > 10 && line[:10] == "ExecStart=" {
			// Find "start" followed by port number
			fmt.Sscanf(extractPortFromExecStart(line), "%d", &port)
			if port > 0 {
				return port
			}
		}
		idx = lineEnd + 1
	}
	return 8181
}

func extractPortFromExecStart(line string) string {
	// ExecStart=/path/to/shield start 8182
	parts := splitSpaces(line)
	for i, p := range parts {
		if p == "start" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func splitSpaces(s string) []string {
	var parts []string
	start := -1
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				parts = append(parts, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	return parts
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", string(output), err)
	}
	return nil
}
