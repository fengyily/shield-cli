//go:build (darwin || windows) && withtray

package tray

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

var (
	dashboardURL string
	onQuitFunc   func()
)

// Run starts the system tray icon. This function blocks and must be called
// from the main goroutine on macOS. Pass the web server startup as onReady.
func Run(port int, onReady func(), onQuit func()) {
	dashboardURL = fmt.Sprintf("http://localhost:%d", port)
	onQuitFunc = onQuit

	systray.Run(func() {
		setupTray(port)
		if onReady != nil {
			onReady()
		}
	}, func() {
		if onQuitFunc != nil {
			onQuitFunc()
		}
	})
}

// Quit exits the system tray
func Quit() {
	systray.Quit()
}

// Available returns true if system tray is supported on this platform
func Available() bool {
	return true
}

func setupTray(port int) {
	systray.SetIcon(iconData)
	systray.SetTitle("")
	systray.SetTooltip(fmt.Sprintf("Shield CLI — localhost:%d", port))

	mOpen := systray.AddMenuItem("Open Dashboard", "Open Shield Web UI in browser")
	systray.AddSeparator()
	mInfo := systray.AddMenuItem(fmt.Sprintf("Port: %d", port), "")
	mInfo.Disable()
	mStatus := systray.AddMenuItem("Status: Running", "")
	mStatus.Disable()
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart Shield", "Restart Shield CLI service")
	mQuit := systray.AddMenuItem("Quit Shield", "Stop Shield CLI and exit")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser(dashboardURL)
			case <-mRestart.ClickedCh:
				restartService()
			case <-mQuit.ClickedCh:
				stopService()
				systray.Quit()
				return
			}
		}
	}()
}

// restartService restarts by exiting the process and letting the service
// manager (launchd KeepAlive / systemd Restart / Windows auto-restart) bring it back up.
func restartService() {
	// Don't unload the service — just quit the process.
	// KeepAlive (macOS) / Restart=on-failure (Linux) / auto start (Windows)
	// will restart us automatically.
	systray.Quit()
}

// stopService unloads the system service so KeepAlive won't restart us
func stopService() {
	switch runtime.GOOS {
	case "darwin":
		home, _ := exec.Command("sh", "-c", "echo $HOME").Output()
		plist := fmt.Sprintf("%s/Library/LaunchAgents/com.yishield.shield-cli.plist", trimNewline(string(home)))
		exec.Command("launchctl", "unload", plist).Run()
	case "windows":
		exec.Command("sc", "stop", "ShieldCLI").Run()
	}
}

func trimNewline(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
