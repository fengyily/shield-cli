package service

import (
	"fmt"
	"net"
	"time"
)

// Config holds the service installation configuration
type Config struct {
	Port       int    // Web UI port (default 8181)
	BinaryPath string // Absolute path to the shield binary
}

// CheckPortAvailable checks if a TCP port is available for binding
func CheckPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	// Small delay to ensure the port is fully released
	time.Sleep(50 * time.Millisecond)
	return true
}

// SuggestPort finds the next available port starting from the given port
func SuggestPort(port int) int {
	for p := port + 1; p <= port+100 && p <= 65535; p++ {
		if CheckPortAvailable(p) {
			return p
		}
	}
	return 0
}
