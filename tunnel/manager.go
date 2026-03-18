package tunnel

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	chclient "github.com/jpillora/chisel/client"
)

// TunnelStatus represents the current state of a tunnel
type TunnelStatus string

const (
	StatusConnecting   TunnelStatus = "connecting"
	StatusConnected    TunnelStatus = "connected"
	StatusDisconnected TunnelStatus = "disconnected"
)

// ConnectionInfo holds the chisel server connection details
type ConnectionInfo struct {
	ExternalIP string
	ServerPort int // API port (from response)
	TunnelPort int // Chisel tunnel port (default 62888)
	Username   string
	Password   string
}

// TunnelEntry represents an active tunnel with status tracking
type TunnelEntry struct {
	Rport  string
	Lip    string
	Lport  string
	Status TunnelStatus
	Client *chclient.Client
	Cancel context.CancelFunc
	closed sync.Once
}

// Manager manages chisel tunnel connections
type Manager struct {
	connInfo ConnectionInfo
	tunnels  map[string]*TunnelEntry
	mu       sync.RWMutex
}

// NewManager creates a new tunnel manager
func NewManager(info ConnectionInfo) *Manager {
	return &Manager{
		connInfo: info,
		tunnels:  make(map[string]*TunnelEntry),
	}
}

// CreateMainTunnel creates the main reverse tunnel with multiple remotes.
// The first remote maps the API port, additional remotes map resource ports.
func (m *Manager) CreateMainTunnel(remotePort, localPort int, extraRemotes ...string) error {
	remotes := []string{fmt.Sprintf("R:%d:localhost:%d", remotePort, localPort)}
	remotes = append(remotes, extraRemotes...)

	cfg := chclient.Config{
		Headers:          http.Header{},
		MaxRetryCount:    -1,
		MaxRetryInterval: 10 * time.Second,
		Server:           fmt.Sprintf("http://%s:%d", m.connInfo.ExternalIP, m.connInfo.TunnelPort),
		Remotes:          remotes,
		Auth:             fmt.Sprintf("%s:%s", m.connInfo.Username, m.connInfo.Password),
	}

	client, err := chclient.NewClient(&cfg)
	if err != nil {
		return fmt.Errorf("failed to create chisel client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	entry := &TunnelEntry{
		Rport:  fmt.Sprintf("%d", remotePort),
		Lip:    "localhost",
		Lport:  fmt.Sprintf("%d", localPort),
		Status: StatusConnecting,
		Client: client,
		Cancel: cancel,
	}

	m.mu.Lock()
	m.tunnels["main"] = entry
	m.mu.Unlock()

	go func() {
		defer func() {
			m.mu.Lock()
			if e, ok := m.tunnels["main"]; ok && e == entry {
				e.Status = StatusDisconnected
			}
			m.mu.Unlock()
			slog.Info("Main tunnel disconnected")
		}()

		if err := client.Start(ctx); err != nil {
			slog.Error("Main tunnel start error", "error", err)
			return
		}

		m.mu.Lock()
		if e, ok := m.tunnels["main"]; ok && e == entry {
			e.Status = StatusConnected
			slog.Info("Main tunnel connected", "remotes", cfg.Remotes)
		}
		m.mu.Unlock()

		if err := client.Wait(); err != nil {
			slog.Error("Main tunnel wait error", "error", err)
		}
	}()

	slog.Info("Main tunnel established",
		"remotes", remotes,
		"server", cfg.Server,
	)
	return nil
}

// CreateDynamicTunnel creates a dynamic tunnel for resource mapping
func (m *Manager) CreateDynamicTunnel(rport, lip, lport string) error {
	m.mu.Lock()
	if existing, ok := m.tunnels[rport]; ok {
		if existing.Status != StatusDisconnected && existing.Lip == lip && existing.Lport == lport {
			m.mu.Unlock()
			return nil
		}
		// Close existing (dead or config changed)
		m.closeEntryLocked(existing, rport)
	}

	entry := &TunnelEntry{
		Rport:  rport,
		Lip:    lip,
		Lport:  lport,
		Status: StatusConnecting,
	}
	m.tunnels[rport] = entry
	m.mu.Unlock()

	cfg := chclient.Config{
		Headers:          http.Header{},
		MaxRetryCount:    -1,
		MaxRetryInterval: 10 * time.Second,
		Server:           fmt.Sprintf("http://%s:%d", m.connInfo.ExternalIP, m.connInfo.TunnelPort),
		Remotes:          []string{fmt.Sprintf("R:127.0.0.1:%s:%s:%s", rport, lip, lport)},
		Auth:             fmt.Sprintf("%s:%s", m.connInfo.Username, m.connInfo.Password),
	}

	client, err := chclient.NewClient(&cfg)
	if err != nil {
		m.mu.Lock()
		delete(m.tunnels, rport)
		m.mu.Unlock()
		return fmt.Errorf("failed to create chisel client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	entry.Client = client
	entry.Cancel = cancel

	go func() {
		defer func() {
			m.mu.Lock()
			if e, ok := m.tunnels[rport]; ok && e == entry {
				e.Status = StatusDisconnected
			}
			m.mu.Unlock()
			slog.Info("Dynamic tunnel disconnected", "rport", rport)
		}()

		if err := client.Start(ctx); err != nil {
			slog.Error("Dynamic tunnel start error", "rport", rport, "error", err)
			return
		}

		m.mu.Lock()
		if e, ok := m.tunnels[rport]; ok && e == entry {
			e.Status = StatusConnected
			slog.Info("Dynamic tunnel connected", "rport", rport, "lip", lip, "lport", lport)
		}
		m.mu.Unlock()

		if err := client.Wait(); err != nil {
			slog.Error("Dynamic tunnel wait error", "rport", rport, "error", err)
		}
	}()

	slog.Info("Dynamic tunnel established",
		"remote", fmt.Sprintf("R:127.0.0.1:%s:%s:%s", rport, lip, lport),
	)
	return nil
}

// closeEntryLocked closes an entry (must be called with lock held)
func (m *Manager) closeEntryLocked(entry *TunnelEntry, key string) {
	entry.closed.Do(func() {
		if entry.Cancel != nil {
			entry.Cancel()
		}
		if entry.Client != nil {
			entry.Client.Close()
		}
	})
	delete(m.tunnels, key)
}

// CloseTunnel closes a specific tunnel by rport
func (m *Manager) CloseTunnel(rport string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, ok := m.tunnels[rport]; ok {
		m.closeEntryLocked(entry, rport)
		slog.Info("Tunnel closed", "rport", rport)
	}
}

// CloseAll closes all tunnels
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, entry := range m.tunnels {
		entry.closed.Do(func() {
			if entry.Cancel != nil {
				entry.Cancel()
			}
			if entry.Client != nil {
				entry.Client.Close()
			}
		})
		delete(m.tunnels, key)
		slog.Info("Tunnel closed", "rport", key)
	}
}

// IsMainConnected returns true if the main tunnel is in connected state.
func (m *Manager) IsMainConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.tunnels["main"]; ok {
		return e.Status == StatusConnected
	}
	return false
}

// List returns info about all active tunnels including status
func (m *Manager) List() map[string]map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]map[string]string)
	for key, entry := range m.tunnels {
		result[key] = map[string]string{
			"rport":  entry.Rport,
			"lip":    entry.Lip,
			"lport":  entry.Lport,
			"status": string(entry.Status),
		}
	}
	return result
}
