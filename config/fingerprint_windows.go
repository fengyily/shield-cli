//go:build windows

package config

import (
	"os/exec"
	"regexp"
	"strings"
)

func getPlatformMachineIDImpl() (string, error) {
	cmd := exec.Command("reg", "query",
		`HKLM\SOFTWARE\Microsoft\Cryptography`,
		"/v", "MachineGuid")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`MachineGuid\s+REG_SZ\s+([^\s]+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1]), nil
	}

	return "", nil
}
