//go:build linux

package config

import (
	"os"
	"strings"
)

func getPlatformMachineIDImpl() (string, error) {
	paths := []string{
		"/etc/machine-id",
		"/var/lib/dbus/machine-id",
		"/sys/class/dmi/id/product_uuid",
	}

	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil {
			id := strings.TrimSpace(string(data))
			if id != "" {
				return id, nil
			}
		}
	}

	return "", nil
}
