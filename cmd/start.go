package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"shield-cli/tray"
	"shield-cli/web"

	"github.com/spf13/cobra"
)

var noTray bool

var startCmd = &cobra.Command{
	Use:   "start [port]",
	Short: "Start the Web management platform",
	Long:  "Start a local Web UI for managing Shield applications.\nDefault port is 8181.\nOn macOS and Windows, a system tray icon is shown for quick access.",
	Example: `  shield start              # Start on port 8181 with tray icon
  shield start 9090         # Start on port 9090
  shield start --no-tray    # Start without system tray icon`,
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

		// Configure logging based on -v flag
		if verbose {
			setupLog(slog.LevelDebug)
		} else {
			setupLog(slog.LevelInfo)
		}

		PrintBanner()

		srv, err := web.NewServer(port)
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
	},
}

func init() {
	startCmd.Flags().BoolVar(&noTray, "no-tray", false, "Disable system tray icon")
}
