//go:build windows

package updater

import "os"

// swapBinary replaces a running Windows executable by renaming it aside first.
// The locked .old file gets removed on next process start via cleanupStale.
func swapBinary(src, dst string) error {
	stale := dst + ".old"
	_ = os.Remove(stale)
	if err := os.Rename(dst, stale); err != nil {
		return err
	}
	if err := os.Rename(src, dst); err != nil {
		// Best-effort rollback so the service does not end up without a binary.
		_ = os.Rename(stale, dst)
		return err
	}
	return nil
}

// cleanupStale removes a leftover shield.exe.old from the previous upgrade.
// Safe to call on every startup; no-op if the file does not exist.
func cleanupStale(exePath string) {
	_ = os.Remove(exePath + ".old")
}
