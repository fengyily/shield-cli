package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const launchAgentLabel = "com.yishield.shield-cli"

var plistTemplate = template.Must(template.New("plist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{ .Label }}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{ .BinaryPath }}</string>
        <string>start</string>
        <string>{{ .Port }}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{ .LogDir }}/shield-cli.log</string>
    <key>StandardErrorPath</key>
    <string>{{ .LogDir }}/shield-cli.error.log</string>
    <key>WorkingDirectory</key>
    <string>{{ .HomeDir }}</string>
</dict>
</plist>
`))

func getPlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist")
}

func getLogDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".shield-cli", "logs")
}

// Install installs shield as a macOS launchd user agent
func Install(cfg Config) error {
	plistPath := getPlistPath()
	logDir := getLogDir()
	home, _ := os.UserHomeDir()

	// Create log directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create plist file
	f, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %w", err)
	}
	defer f.Close()

	data := struct {
		Label      string
		BinaryPath string
		Port       string
		LogDir     string
		HomeDir    string
	}{
		Label:      launchAgentLabel,
		BinaryPath: cfg.BinaryPath,
		Port:       fmt.Sprintf("%d", cfg.Port),
		LogDir:     logDir,
		HomeDir:    home,
	}

	if err := plistTemplate.Execute(f, data); err != nil {
		os.Remove(plistPath)
		return fmt.Errorf("failed to write plist: %w", err)
	}

	// Load the service
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load service: %s (%w)", string(output), err)
	}

	return nil
}

// Uninstall removes the shield launchd service
func Uninstall() error {
	plistPath := getPlistPath()

	// Check if plist exists
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return fmt.Errorf("service is not installed")
	}

	// Unload the service
	cmd := exec.Command("launchctl", "unload", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Warning: launchctl unload: %s\n", string(output))
	}

	// Remove the plist file
	if err := os.Remove(plistPath); err != nil {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	return nil
}

// IsInstalled checks if the service is already installed
func IsInstalled() bool {
	plistPath := getPlistPath()
	_, err := os.Stat(plistPath)
	return err == nil
}

// Status returns the current service status
func Status() (string, error) {
	plistPath := getPlistPath()
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return "not installed", nil
	}

	cmd := exec.Command("launchctl", "list", launchAgentLabel)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "installed (not running)", nil
	}
	_ = output
	return "running", nil
}

// GetInstalledPort reads the port from the installed plist (returns 0 if not found)
func GetInstalledPort() int {
	plistPath := getPlistPath()
	data, err := os.ReadFile(plistPath)
	if err != nil {
		return 0
	}
	// Simple parse: find the port argument after "start"
	content := string(data)
	// Look for <string>start</string> followed by <string>PORT</string>
	idx := 0
	for idx < len(content) {
		startIdx := indexOf(content, "<string>start</string>", idx)
		if startIdx < 0 {
			break
		}
		portStart := indexOf(content, "<string>", startIdx+len("<string>start</string>"))
		if portStart < 0 {
			break
		}
		portStart += len("<string>")
		portEnd := indexOf(content, "</string>", portStart)
		if portEnd < 0 {
			break
		}
		portStr := content[portStart:portEnd]
		var port int
		fmt.Sscanf(portStr, "%d", &port)
		if port > 0 {
			return port
		}
		idx = portEnd
	}
	return 8181
}

func indexOf(s, substr string, start int) int {
	if start >= len(s) {
		return -1
	}
	idx := 0
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
		idx++
	}
	return -1
}
