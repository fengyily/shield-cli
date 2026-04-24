package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"shield-cli/service"
	"shield-cli/tray"
	"shield-cli/updater"
	"shield-cli/web"

	"github.com/spf13/cobra"
)

var noTray bool

var startCmd = &cobra.Command{
	Use:   "start [port]",
	Short: "Start the Web management platform",
	Long: `Start the Shield Web management platform.

If Shield is installed as a system service (via "shield install"),
this command starts the background service.

Otherwise, it runs the Web UI in the foreground.
Use Ctrl+C to stop the foreground process.`,
	Example: `  shield start              # Start service or foreground on port 8181
  shield start 9090         # Start on port 9090 (foreground only)
  shield start --no-tray    # Foreground without system tray icon`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := 8181
		if len(args) == 1 {
			p, err := strconv.Atoi(args[0])
			if err != nil || p < 1 || p > 65535 {
				return fmt.Errorf("invalid port: %s", args[0])
			}
			port = p
		}

		// If installed as system service and no custom port specified, start the service
		if service.IsInstalled() && len(args) == 0 {
			return startAsService()
		}

		// Otherwise, run in foreground
		return startForeground(port)
	},
}

// startAsService starts the installed system service
func startAsService() error {
	PrintBanner()

	port := service.GetInstalledPort()

	// Check if already running
	status, _ := service.Status()
	if status == "running" {
		fmt.Printf("  \033[1;32m✓ Service is already running\033[0m\n\n")
		fmt.Printf("  Web UI: \033[1;36mhttp://localhost:%d\033[0m\n\n", port)
		fmt.Printf("  \033[90mCommands:\033[0m\n")
		fmt.Printf("    shield stop         Stop the service\n")
		fmt.Printf("    shield uninstall    Remove the service\n\n")
		return nil
	}

	fmt.Printf("  \033[90mStarting Shield service...\033[0m\n\n")

	if err := service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w\n\n  Try: shield start %d  (run in foreground instead)", err, port)
	}

	fmt.Printf("  \033[1;32m✓ Service started successfully\033[0m\n\n")
	fmt.Printf("  Web UI: \033[1;36mhttp://localhost:%d\033[0m\n\n", port)
	fmt.Printf("  \033[90mCommands:\033[0m\n")
	fmt.Printf("    shield stop         Stop the service\n")
	fmt.Printf("    shield uninstall    Remove the service\n")
	fmt.Printf("    shield start %d    Run in foreground (debug)\n\n", port)
	return nil
}

// startForeground runs the Web UI in the foreground
func startForeground(port int) error {
	// Configure logging based on -v flag
	if verbose {
		setupLog(slog.LevelDebug)
	} else {
		setupLog(slog.LevelInfo)
	}

	PrintBanner()

	if exe, err := os.Executable(); err == nil {
		updater.CleanupStale(exe)
	}

	if service.IsInstalled() {
		fmt.Printf("  \033[1;33m⚠ Service is installed but starting in foreground mode\033[0m\n")
		fmt.Printf("  \033[90m  Use \"shield start\" (without port) to start the background service\033[0m\n\n")
	}

	srv, err := web.NewServer(port, Version, GitCommit, BuildTime)
	if err != nil {
		return err
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	shutdown := func() {
		fmt.Println("\n\033[33mShutting down...\033[0m")
		srv.Shutdown()
		os.Exit(0)
	}

	go func() {
		<-sigCh
		tray.Quit()
		shutdown()
	}()

	// If tray is available and not disabled, run with system tray
	if tray.Available() && !noTray {
		slog.Info("System tray enabled", "platform", "macOS/Windows")
		tray.Run(port, func() {
			// onReady: start the web server in a goroutine
			go func() {
				if err := srv.Start(); err != nil {
					slog.Error("Web server error", "error", err)
					tray.Quit()
				}
			}()
		}, shutdown)
		return nil
	}

	// No tray: run web server directly (blocks)
	return srv.Start()
}

func init() {
	startCmd.Flags().BoolVar(&noTray, "no-tray", false, "Disable system tray icon")
}
