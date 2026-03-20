//go:build (!darwin && !windows) || !withtray

package tray

// Run is a no-op on unsupported platforms.
// The onReady callback is called immediately.
func Run(port int, onReady func(), onQuit func()) {
	if onReady != nil {
		onReady()
	}
	// Block forever (will be interrupted by signal handler)
	select {}
}

// Quit is a no-op on unsupported platforms
func Quit() {}

// Available returns false on unsupported platforms
func Available() bool {
	return false
}
