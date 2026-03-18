package cmd

import (
	"fmt"
	"os"
	"strings"

	"shield-cli/config"

	"github.com/spf13/cobra"
)

var validProtocols = []string{"ssh", "http", "https", "rdp", "vnc", "telnet"}

// defaultPorts maps each protocol to its standard port
var defaultPorts = map[string]int{
	"ssh":    22,
	"http":   80,
	"https":  443,
	"rdp":    3389,
	"vnc":    5900,
	"telnet": 23,
}

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
	Use:          "shield <protocol> [ip:port] [flags]",
	Short:        "Shield CLI - Secure Tunnel Connector",
	Long:         "Shield CLI exposes internal network resources to the public server via secure tunnels.\n\nIf ip:port is omitted, defaults to 127.0.0.1 with the protocol's standard port.\nIf only ip is given, the protocol's standard port is used.",
	Example:      "  shield ssh                        # 127.0.0.1:22\n  shield ssh 2222                   # 127.0.0.1:2222\n  shield ssh 10.0.0.2               # 10.0.0.2:22\n  shield ssh 10.0.0.2:2222          # 10.0.0.2:2222\n  shield http 3000                  # 127.0.0.1:3000\n  shield rdp 192.168.1.100          # 192.168.1.100:3389\n  shield https 10.0.0.5 --invisible",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(0),
	RunE:         runShield,
}

func isValidProtocol(p string) bool {
	for _, v := range validProtocols {
		if strings.EqualFold(p, v) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.Flags().StringVarP(&protocol, "type", "t", "", "Protocol type (deprecated, use positional arg instead)")
	rootCmd.Flags().StringVarP(&target, "source", "s", "", "Target address (deprecated, use positional arg instead)")
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
