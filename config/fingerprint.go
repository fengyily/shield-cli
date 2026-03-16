package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

// GetMachineFingerprint gets the unique machine fingerprint
func GetMachineFingerprint() (string, error) {
	var parts []string

	// 1. Hostname
	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		parts = append(parts, hostname)
	}

	// 2. First valid MAC address
	mac, err := getFirstMACAddress()
	if err == nil && mac != "" {
		parts = append(parts, mac)
	}

	// 3. Platform-specific machine ID
	machineID, err := getPlatformMachineID()
	if err == nil && machineID != "" {
		parts = append(parts, machineID)
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("unable to get machine identification info")
	}

	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:]), nil
}

// getFirstMACAddress gets the first valid MAC address
func getFirstMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}

		nameLower := strings.ToLower(iface.Name)
		virtualPrefixes := []string{"docker", "br-", "veth", "virbr", "vnet", "tun", "tap", "vmnet", "vboxnet"}
		isVirtual := false
		for _, prefix := range virtualPrefixes {
			if strings.HasPrefix(nameLower, prefix) {
				isVirtual = true
				break
			}
		}
		if isVirtual {
			continue
		}

		return iface.HardwareAddr.String(), nil
	}

	return "", fmt.Errorf("no valid MAC address found")
}

func getPlatformMachineID() (string, error) {
	return getPlatformMachineIDImpl()
}
