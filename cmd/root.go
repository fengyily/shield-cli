package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	protocol  string
	server    string
	target    string
	apiServer  string
	verbose    bool
	tunnelPort int
)

var rootCmd = &cobra.Command{
	Use:   "shield",
	Short: "Shield CLI - Secure Tunnel Connector",
	Long:  "Shield CLI exposes internal network resources to the public server via secure tunnels.",
	RunE:  runShield,
}

func init() {
	rootCmd.Flags().StringVarP(&protocol, "type", "t", "", "Protocol type (e.g., ssh, http, https, tcp)")
	rootCmd.Flags().StringVarP(&target, "source", "s", "", "Target address in ip:port format (e.g., 10.0.0.2:22)")
	rootCmd.Flags().StringVarP(&apiServer, "server", "H", "https://console.yishield.com/raas", "API server URL")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose log output")
	rootCmd.Flags().IntVarP(&tunnelPort, "tunnel-port", "p", 62888, "Chisel tunnel server port")

	rootCmd.MarkFlagRequired("type")
	rootCmd.MarkFlagRequired("source")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
