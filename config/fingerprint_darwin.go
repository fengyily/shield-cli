//go:build darwin

package config

import (
	"os/exec"
	"regexp"
	"strings"
)

func getPlatformMachineIDImpl() (string, error) {
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"IOPlatformUUID"\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1]), nil
	}

	return "", nil
}
