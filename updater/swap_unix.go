//go:build !windows

package updater

import "os"

// swapBinary atomically replaces dst with src. On Unix the running process
// continues executing the original inode, so the rename is safe.
func swapBinary(src, dst string) error {
	if err := os.Chmod(src, 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// cleanupStale is a no-op on Unix; the old inode is freed automatically.
func cleanupStale(string) {}
