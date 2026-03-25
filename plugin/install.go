package plugin

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	mainRepo         = "fengyily/shield-cli"
	defaultSourceOrg = "fengyily"
)

// KnownPlugins maps plugin names to their GitHub repo and metadata.
var KnownPlugins = map[string]PluginInfo{
	"mysql": {
		Name:        "mysql",
		Source:      mainRepo, // built in-repo, assets attached to main release
		Protocols:   []string{"mysql", "mariadb"},
		DefaultPort: 3306,
	},
	"postgres": {
		Name:        "postgres",
		Source:      defaultSourceOrg + "/shield-plugin-postgres",
		Protocols:   []string{"postgres", "pg", "postgresql"},
		DefaultPort: 5432,
	},
	"sqlserver": {
		Name:        "sqlserver",
		Source:      defaultSourceOrg + "/shield-plugin-sqlserver",
		Protocols:   []string{"sqlserver", "mssql"},
		DefaultPort: 1433,
	},
}

// githubRelease represents a GitHub release API response.
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Install downloads and installs a plugin by name.
func Install(name string) (*PluginInfo, error) {
	known, ok := KnownPlugins[name]
	if !ok {
		return nil, fmt.Errorf("unknown plugin %q\n\nAvailable plugins: %s", name, AvailablePluginNames())
	}

	// Check if it's a built-in plugin (source starts with "builtin:")
	if strings.HasPrefix(known.Source, "builtin:") {
		return installBuiltin(name, &known)
	}

	// Fetch latest release from GitHub
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", known.Source)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d for %s", resp.StatusCode, known.Source)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Find matching asset for current platform
	osName := runtime.GOOS
	archName := runtime.GOARCH
	assetName := fmt.Sprintf("shield-plugin-%s_%s_%s.tar.gz", name, osName, archName)

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return nil, fmt.Errorf("no release asset found for %s/%s\n\nLooking for: %s", osName, archName, assetName)
	}

	// Download and extract
	binName := fmt.Sprintf("shield-plugin-%s", name)
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	if err := downloadAndExtract(downloadURL, binName); err != nil {
		return nil, fmt.Errorf("failed to install: %w", err)
	}

	// Register in registry
	info := known
	info.Version = release.TagName
	info.Binary = binName

	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	reg.Register(info)
	if err := reg.Save(); err != nil {
		return nil, err
	}

	return &info, nil
}

// installBuiltin handles plugins that are built into the shield-cli binary.
func installBuiltin(name string, known *PluginInfo) (*PluginInfo, error) {
	// For builtin plugins, we create a symlink or copy of the main binary
	// and register with a special binary name
	binName := fmt.Sprintf("shield-plugin-%s", name)
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	destPath := filepath.Join(PluginsDir(), binName)
	if err := os.MkdirAll(PluginsDir(), 0755); err != nil {
		return nil, err
	}

	// Copy the main binary as the plugin binary
	src, err := os.Open(execPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	info := *known
	info.Version = "builtin"
	info.Binary = binName

	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	reg.Register(info)
	if err := reg.Save(); err != nil {
		return nil, err
	}

	return &info, nil
}

func downloadAndExtract(url, binName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	dir := PluginsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Extract tar.gz
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip error: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar error: %w", err)
		}

		// Only extract the binary
		base := filepath.Base(header.Name)
		if base == binName || strings.TrimSuffix(base, ".exe") == strings.TrimSuffix(binName, ".exe") {
			destPath := filepath.Join(dir, binName)
			f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
			return nil
		}
	}

	return fmt.Errorf("binary %q not found in archive", binName)
}

// InstallFromLocal installs a plugin from a local binary path.
func InstallFromLocal(name, binaryPath string) (*PluginInfo, error) {
	known, ok := KnownPlugins[name]
	if !ok {
		return nil, fmt.Errorf("unknown plugin %q\n\nAvailable plugins: %s", name, AvailablePluginNames())
	}

	// Verify source binary exists
	if _, err := os.Stat(binaryPath); err != nil {
		return nil, fmt.Errorf("binary not found: %s", binaryPath)
	}

	binName := fmt.Sprintf("shield-plugin-%s", name)
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	destPath := filepath.Join(PluginsDir(), binName)
	if err := os.MkdirAll(PluginsDir(), 0755); err != nil {
		return nil, err
	}

	// Copy binary
	src, err := os.Open(binaryPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	info := known
	info.Version = "local"
	info.Binary = binName

	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	reg.Register(info)
	if err := reg.Save(); err != nil {
		return nil, err
	}

	return &info, nil
}

// UpgradeResult describes the result of an upgrade check or operation.
type UpgradeResult struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
	Upgraded       bool
}

// CheckUpdate checks if a newer version is available for a plugin without installing it.
func CheckUpdate(name string) (*UpgradeResult, error) {
	known, ok := KnownPlugins[name]
	if !ok {
		return nil, fmt.Errorf("unknown plugin %q", name)
	}

	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	installed := reg.FindByName(name)
	if installed == nil {
		return nil, fmt.Errorf("plugin %q is not installed", name)
	}

	current := installed.Version
	if current == "local" || current == "builtin" {
		return &UpgradeResult{Name: name, CurrentVersion: current, LatestVersion: current}, nil
	}

	latest, err := fetchLatestVersion(known.Source)
	if err != nil {
		return nil, err
	}

	return &UpgradeResult{
		Name:           name,
		CurrentVersion: current,
		LatestVersion:  latest,
	}, nil
}

// Upgrade upgrades a plugin to the latest version if available.
func Upgrade(name string) (*UpgradeResult, error) {
	known, ok := KnownPlugins[name]
	if !ok {
		return nil, fmt.Errorf("unknown plugin %q", name)
	}

	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	installed := reg.FindByName(name)
	if installed == nil {
		return nil, fmt.Errorf("plugin %q is not installed\n\nInstall first: shield plugin add %s", name, name)
	}

	current := installed.Version
	if current == "local" || current == "builtin" {
		return nil, fmt.Errorf("plugin %q was installed from %s, cannot upgrade from GitHub\n\nReinstall: shield plugin remove %s && shield plugin add %s", name, current, name, name)
	}

	latest, err := fetchLatestVersion(known.Source)
	if err != nil {
		return nil, err
	}

	if latest == current {
		return &UpgradeResult{Name: name, CurrentVersion: current, LatestVersion: latest, Upgraded: false}, nil
	}

	// Download and install the new version
	info, err := Install(name)
	if err != nil {
		return nil, fmt.Errorf("upgrade failed: %w", err)
	}

	return &UpgradeResult{
		Name:           name,
		CurrentVersion: current,
		LatestVersion:  info.Version,
		Upgraded:       true,
	}, nil
}

// UpgradeAll upgrades all installed plugins. Returns results for each.
func UpgradeAll() ([]UpgradeResult, error) {
	reg, err := LoadRegistry()
	if err != nil {
		return nil, err
	}

	var results []UpgradeResult
	for _, p := range reg.Plugins {
		res, err := Upgrade(p.Name)
		if err != nil {
			// Non-fatal: collect the error as a result
			results = append(results, UpgradeResult{
				Name:           p.Name,
				CurrentVersion: p.Version,
				LatestVersion:  "error: " + err.Error(),
			})
			continue
		}
		results = append(results, *res)
	}
	return results, nil
}

// fetchLatestVersion fetches the latest release tag from GitHub.
func fetchLatestVersion(source string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", source)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}
	return release.TagName, nil
}

// AvailablePluginNames returns a comma-separated list of known plugin names.
func AvailablePluginNames() string {
	names := make([]string, 0, len(KnownPlugins))
	for name := range KnownPlugins {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}
