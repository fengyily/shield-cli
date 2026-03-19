package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"shield-cli/web"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [port]",
	Short: "Start the Web management platform",
	Long:  "Start a local Web UI for managing Shield applications.\nDefault port is 8181.",
	Example: `  shield start          # Start on port 8181
  shield start 9090     # Start on port 9090`,
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
		go func() {
			<-sigCh
			fmt.Println("\n\033[33mShutting down...\033[0m")
			srv.Shutdown()
			os.Exit(0)
		}()

		return srv.Start()
	},
}
