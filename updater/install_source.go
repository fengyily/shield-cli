package updater

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// InstallSource describes how the running binary was installed. The Web UI
// hides the self-upgrade button for sources we cannot safely overwrite.
type InstallSource struct {
	Kind     string `json:"kind"`     // "binary" | "brew" | "scoop" | "apt" | "rpm" | "docker" | "unknown"
	Writable bool   `json:"writable"` // whether the current process can replace the binary
	Path     string `json:"path"`
	Hint     string `json:"hint,omitempty"` // suggested upgrade command for managed installs
}

// DetectInstallSource inspects the current executable path to classify how it
// was installed.
func DetectInstallSource() InstallSource {
	src := InstallSource{Kind: "unknown"}

	exe, err := os.Executable()
	if err != nil {
		return src
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	src.Path = exe
	src.Kind = classify(exe)
	src.Hint = hintFor(src.Kind)
	src.Writable = canReplace(exe)
	return src
}

func classify(path string) string {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "docker"
	}
	p := strings.ToLower(filepath.ToSlash(path))
	switch {
	case strings.Contains(p, "/cellar/") || strings.Contains(p, "/homebrew/") || strings.Contains(p, "/linuxbrew/"):
		return "brew"
	case strings.Contains(p, "/scoop/"):
		return "scoop"
	case runtime.GOOS == "linux" && (strings.HasPrefix(p, "/usr/bin/") || strings.HasPrefix(p, "/usr/sbin/")):
		// deb/rpm install to /usr/bin; install.sh defaults to /usr/local/bin
		return packageKindLinux()
	}
	return "binary"
}

func packageKindLinux() string {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "apt"
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "rpm"
	}
	return "rpm"
}

func hintFor(kind string) string {
	switch kind {
	case "brew":
		return "brew upgrade shield-cli"
	case "scoop":
		return "scoop update shield-cli"
	case "apt":
		return "sudo apt update && sudo apt install --only-upgrade shield-cli"
	case "rpm":
		return "sudo yum update shield-cli"
	case "docker":
		return "docker pull fengyily/shield-cli && docker restart shield"
	}
	return ""
}

func canReplace(path string) bool {
	dir := filepath.Dir(path)
	// Directory must be writable so we can atomically rename a sibling file in.
	testFile, err := os.CreateTemp(dir, ".shield-upgrade-probe-*")
	if err != nil {
		return false
	}
	name := testFile.Name()
	testFile.Close()
	os.Remove(name)
	return true
}
