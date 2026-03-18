package cmd

import (
	"fmt"
	"os"

	"shield-cli/config"

	"github.com/spf13/cobra"
)

var (
	protocol    string
	server      string
	target      string
	apiServer   string
	verbose     bool
	tunnelPort  int
	visable     string
	invisible   bool
	displayName string
	siteName    string
	authUser    string
	authPass    string
	privateKey  string
	passphrase  string
	enableSftp  bool
)

var rootCmd = &cobra.Command{
	Use:          "shield",
	Short:        "Shield CLI - Secure Tunnel Connector",
	Long:         "Shield CLI exposes internal network resources to the public server via secure tunnels.",
	SilenceUsage: true,
	RunE:         runShield,
}

func init() {
	rootCmd.Flags().StringVarP(&protocol, "type", "t", "", "Protocol type (e.g., ssh, http, https, tcp)")
	rootCmd.Flags().StringVarP(&target, "source", "s", "", "Target address in ip:port format (e.g., 10.0.0.2:22)")
	rootCmd.Flags().StringVarP(&apiServer, "server", "H", "https://console.yishield.com/raas", "API server URL")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose log output")
	rootCmd.Flags().IntVarP(&tunnelPort, "tunnel-port", "p", 62888, "Chisel tunnel server port")
	rootCmd.Flags().StringVar(&visable, "visable", "visable", "AC node name filter (default: visable)")
	rootCmd.Flags().BoolVar(&invisible, "invisible", false, "Invisible mode: require Access URL with authorization key")
	rootCmd.Flags().StringVar(&displayName, "display-name", "", "Connector display name")
	rootCmd.Flags().StringVar(&siteName, "site-name", "", "Application site name")
	rootCmd.Flags().StringVar(&authUser, "username", "", "Target service username (SSH/RDP/VNC)")
	rootCmd.Flags().StringVar(&authPass, "auth-pass", "", "Target service password (SSH/RDP/VNC)")
	rootCmd.Flags().StringVar(&privateKey, "private-key", "", "SSH private key")
	rootCmd.Flags().StringVar(&passphrase, "passphrase", "", "SSH private key passphrase")
	rootCmd.Flags().BoolVar(&enableSftp, "enable-sftp", false, "Enable SFTP (SSH only)")

	rootCmd.MarkFlagRequired("type")
	rootCmd.MarkFlagRequired("source")

	// Subcommand: clear cached credentials
	rootCmd.AddCommand(&cobra.Command{
		Use:   "clean",
		Short: "Clear cached credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.GetCredentialFilePath()
			if err := os.Remove(path); err != nil {
				if os.IsNotExist(err) {
					fmt.Println("No cached credentials found.")
					return nil
				}
				return fmt.Errorf("failed to remove credentials: %w", err)
			}
			fmt.Printf("Credentials cleared: %s\n", path)
			return nil
		},
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
