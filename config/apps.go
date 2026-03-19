package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

const maxApps = 10

// AppConfig represents a saved application configuration
type AppConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Protocol    string    `json:"protocol"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Server      string    `json:"server"`
	TunnelPort  int       `json:"tunnel_port"`
	Visable     string    `json:"visable"`
	Invisible   bool      `json:"invisible"`
	Username    string    `json:"username,omitempty"`
	AuthPass    string    `json:"auth_pass,omitempty"`
	PrivateKey  string    `json:"private_key,omitempty"`
	Passphrase  string    `json:"passphrase,omitempty"`
	EnableSftp  bool      `json:"enable_sftp,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	SiteName    string    `json:"site_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastConnectedAt *time.Time `json:"last_connected_at,omitempty"`
}

// AppStore manages encrypted app configurations
type AppStore struct {
	path string
}

// NewAppStore creates an AppStore using the default storage path
func NewAppStore() *AppStore {
	return &AppStore{path: getAppStorePath()}
}

func getAppStorePath() string {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LOCALAPPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		dir = filepath.Join(dir, "ShieldCLI")
	default:
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".shield-cli")
	}
	return filepath.Join(dir, ".apps")
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// List returns all saved apps sorted by updated time (newest first)
func (s *AppStore) List() ([]AppConfig, error) {
	apps, err := s.load()
	if err != nil {
		return nil, err
	}
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].UpdatedAt.After(apps[j].UpdatedAt)
	})
	return apps, nil
}

// Get returns a single app by ID
func (s *AppStore) Get(id string) (*AppConfig, error) {
	apps, err := s.load()
	if err != nil {
		return nil, err
	}
	for i := range apps {
		if apps[i].ID == id {
			return &apps[i], nil
		}
	}
	return nil, fmt.Errorf("app not found: %s", id)
}

// Add creates a new app config, enforcing the max limit
func (s *AppStore) Add(app AppConfig) (*AppConfig, error) {
	apps, err := s.load()
	if err != nil {
		return nil, err
	}

	app.ID = generateID()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	// Apply defaults
	if app.Name == "" {
		app.Name = fmt.Sprintf("%s %s:%d", strings.ToUpper(app.Protocol), app.IP, app.Port)
	}
	if app.Server == "" {
		app.Server = "https://console.yishield.com/raas"
	}
	if app.TunnelPort == 0 {
		app.TunnelPort = 62888
	}
	if app.Visable == "" && !app.Invisible {
		app.Visable = "visable"
	}

	apps = append(apps, app)

	// Enforce max limit: remove oldest if over limit
	if len(apps) > maxApps {
		sort.Slice(apps, func(i, j int) bool {
			return apps[i].UpdatedAt.After(apps[j].UpdatedAt)
		})
		apps = apps[:maxApps]
	}

	if err := s.save(apps); err != nil {
		return nil, err
	}
	return &app, nil
}

// Update modifies an existing app config
func (s *AppStore) Update(id string, updated AppConfig) (*AppConfig, error) {
	apps, err := s.load()
	if err != nil {
		return nil, err
	}

	for i := range apps {
		if apps[i].ID == id {
			updated.ID = id
			updated.CreatedAt = apps[i].CreatedAt
			updated.UpdatedAt = time.Now()
			if updated.Server == "" {
				updated.Server = "https://console.yishield.com/raas"
			}
			if updated.TunnelPort == 0 {
				updated.TunnelPort = 62888
			}
			if updated.Visable == "" && !updated.Invisible {
				updated.Visable = "visable"
			}
			apps[i] = updated
			if err := s.save(apps); err != nil {
				return nil, err
			}
			return &updated, nil
		}
	}
	return nil, fmt.Errorf("app not found: %s", id)
}

// Delete removes an app config by ID
func (s *AppStore) Delete(id string) error {
	apps, err := s.load()
	if err != nil {
		return err
	}

	for i := range apps {
		if apps[i].ID == id {
			apps = append(apps[:i], apps[i+1:]...)
			return s.save(apps)
		}
	}
	return fmt.Errorf("app not found: %s", id)
}

// UpdateLastConnected sets LastConnectedAt to now for the given app ID and saves
func (s *AppStore) UpdateLastConnected(id string) error {
	apps, err := s.load()
	if err != nil {
		return err
	}
	for i := range apps {
		if apps[i].ID == id {
			now := time.Now()
			apps[i].LastConnectedAt = &now
			return s.save(apps)
		}
	}
	return fmt.Errorf("app not found: %s", id)
}

func (s *AppStore) load() ([]AppConfig, error) {
	encrypted, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []AppConfig{}, nil
		}
		return nil, fmt.Errorf("failed to read apps file: %w", err)
	}

	key, err := getDerivedKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get decryption key: %w", err)
	}

	plaintext, err := decryptAESGCM(encrypted, key)
	if err != nil {
		// Corrupted or fingerprint mismatch — reset to empty
		os.Remove(s.path)
		return []AppConfig{}, nil
	}

	var apps []AppConfig
	if err := json.Unmarshal(plaintext, &apps); err != nil {
		os.Remove(s.path)
		return []AppConfig{}, nil
	}
	return apps, nil
}

func (s *AppStore) save(apps []AppConfig) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	plaintext, err := json.Marshal(apps)
	if err != nil {
		return fmt.Errorf("failed to serialize apps: %w", err)
	}

	key, err := getDerivedKey()
	if err != nil {
		return fmt.Errorf("failed to get encryption key: %w", err)
	}

	encrypted, err := encryptAESGCM(plaintext, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt apps: %w", err)
	}

	return os.WriteFile(s.path, encrypted, 0600)
}
