package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const downloadTimeoutSecs = 120

// CleanupStale removes any leftover *.old file next to the given executable
// path. Only meaningful on Windows where an in-place upgrade renames the
// running binary aside; on Unix it is a no-op.
func CleanupStale(exePath string) { cleanupStale(exePath) }

// assetNames returns (archive, checksums) file names produced by goreleaser
// for the current platform. Matches .goreleaser.yaml archives section.
func assetNames() (archive, checksums string) {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	archive = fmt.Sprintf("shield-%s-%s.%s", runtime.GOOS, runtime.GOARCH, ext)
	checksums = "checksums.txt"
	return
}

// downloadURLs returns the archive URL and checksums URL for a release tag.
func downloadURLs(tag string) (archiveURL, checksumsURL string) {
	base := os.Getenv("SHIELD_UPDATE_ASSET_BASE")
	if base == "" {
		base = fmt.Sprintf("https://github.com/fengyily/shield-cli/releases/download/v%s", strings.TrimPrefix(tag, "v"))
	}
	archive, checksums := assetNames()
	return base + "/" + archive, base + "/" + checksums
}

// Apply downloads the release archive for the current platform, verifies its
// SHA256 against checksums.txt, extracts the `shield` binary, and atomically
// replaces the running executable. The caller is responsible for restarting
// the process afterward.
func Apply(ctx context.Context, latestTag string, progress func(stage string, pct int)) error {
	if latestTag == "" {
		return fmt.Errorf("latest version is empty")
	}
	if progress == nil {
		progress = func(string, int) {}
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}

	archiveURL, checksumsURL := downloadURLs(latestTag)
	archiveName, _ := assetNames()

	progress("download", 10)
	tmpDir, err := os.MkdirTemp("", "shield-upgrade-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, archiveName)
	if err := download(ctx, archiveURL, archivePath); err != nil {
		return fmt.Errorf("download archive: %w", err)
	}

	progress("verify", 50)
	checksumsPath := filepath.Join(tmpDir, "checksums.txt")
	if err := download(ctx, checksumsURL, checksumsPath); err != nil {
		return fmt.Errorf("download checksums: %w", err)
	}
	if err := verifySHA256(archivePath, checksumsPath, archiveName); err != nil {
		return fmt.Errorf("verify checksum: %w", err)
	}

	progress("extract", 70)
	binaryName := "shield"
	if runtime.GOOS == "windows" {
		binaryName = "shield.exe"
	}
	extractedBinary := filepath.Join(tmpDir, binaryName)
	if err := extractBinary(archivePath, binaryName, extractedBinary); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	progress("install", 90)
	if err := swapBinary(extractedBinary, exe); err != nil {
		return fmt.Errorf("install: %w", err)
	}

	progress("done", 100)
	return nil
}

func download(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "shield-cli-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func verifySHA256(archivePath, checksumsPath, archiveName string) error {
	data, err := os.ReadFile(checksumsPath)
	if err != nil {
		return err
	}
	var want string
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == archiveName {
			want = fields[0]
			break
		}
	}
	if want == "" {
		return fmt.Errorf("no checksum entry for %s", archiveName)
	}

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("sha256 mismatch: got %s want %s", got, want)
	}
	return nil
}

func extractBinary(archivePath, binaryName, dst string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractFromZip(archivePath, binaryName, dst)
	}
	return extractFromTarGz(archivePath, binaryName, dst)
}

func extractFromTarGz(archivePath, binaryName, dst string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("%s not found in archive", binaryName)
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) != binaryName {
			continue
		}
		return writeBinary(tr, dst)
	}
}

func extractFromZip(archivePath, binaryName, dst string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, zf := range zr.File {
		if filepath.Base(zf.Name) != binaryName {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return writeBinary(rc, dst)
	}
	return fmt.Errorf("%s not found in archive", binaryName)
}

func writeBinary(r io.Reader, dst string) error {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
