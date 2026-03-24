package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// PluginInfo describes an installed plugin.
type PluginInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Binary      string   `json:"binary"`
	Source      string   `json:"source"`
	Protocols   []string `json:"protocols"`
	DefaultPort int      `json:"default_port"`
	InstalledAt string   `json:"installed_at"`
}

// Registry holds the list of installed plugins.
type Registry struct {
	Plugins []PluginInfo `json:"plugins"`
}

// PluginsDir returns the path to the plugins directory.
func PluginsDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "ShieldCLI", "plugins")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".shield-cli", "plugins")
}

func registryPath() string {
	return filepath.Join(PluginsDir(), "registry.json")
}

// LoadRegistry reads the registry from disk.
func LoadRegistry() (*Registry, error) {
	data, err := os.ReadFile(registryPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{}, nil
		}
		return nil, err
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

// Save writes the registry to disk.
func (r *Registry) Save() error {
	dir := PluginsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(registryPath(), data, 0644)
}

// Find returns the plugin info for a given protocol name, or nil if not found.
func (r *Registry) Find(protocol string) *PluginInfo {
	for i := range r.Plugins {
		for _, p := range r.Plugins[i].Protocols {
			if p == protocol {
				return &r.Plugins[i]
			}
		}
	}
	return nil
}

// FindByName returns the plugin info by plugin name.
func (r *Registry) FindByName(name string) *PluginInfo {
	for i := range r.Plugins {
		if r.Plugins[i].Name == name {
			return &r.Plugins[i]
		}
	}
	return nil
}

// Register adds or updates a plugin in the registry.
func (r *Registry) Register(info PluginInfo) {
	info.InstalledAt = time.Now().Format(time.RFC3339)
	for i := range r.Plugins {
		if r.Plugins[i].Name == info.Name {
			r.Plugins[i] = info
			return
		}
	}
	r.Plugins = append(r.Plugins, info)
}

// Remove deletes a plugin from the registry and its binary.
func (r *Registry) Remove(name string) error {
	for i := range r.Plugins {
		if r.Plugins[i].Name == name {
			// Remove binary
			binPath := filepath.Join(PluginsDir(), r.Plugins[i].Binary)
			os.Remove(binPath)
			// Remove from list
			r.Plugins = append(r.Plugins[:i], r.Plugins[i+1:]...)
			return r.Save()
		}
	}
	return fmt.Errorf("plugin %q not found", name)
}
