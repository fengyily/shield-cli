package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"shield-cli/config"
	"shield-cli/plugin"
	"shield-cli/tunnel"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// runWithPlugin handles protocol execution via a plugin.
// It starts the plugin process, gets the web port, then tunnels it as HTTP.
func runWithPlugin(cmd *cobra.Command, info *plugin.PluginInfo, rawTarget string) error {
	// Apply default port from plugin if needed
	if rawTarget == "" {
		rawTarget = fmt.Sprintf("127.0.0.1:%d", info.DefaultPort)
	} else if !strings.Contains(rawTarget, ":") && !strings.Contains(rawTarget, ".") {
		// Pure number: shield mysql 3307 => 127.0.0.1:3307
		rawTarget = fmt.Sprintf("127.0.0.1:%s", rawTarget)
	} else if !strings.Contains(rawTarget, ":") {
		// Only IP: shield mysql 10.0.0.2 => 10.0.0.2:<default>
		rawTarget = fmt.Sprintf("%s:%d", rawTarget, info.DefaultPort)
	}

	dbIP, dbPort, err := parseTarget(rawTarget)
	if err != nil {
		return fmt.Errorf("invalid target address: %w", err)
	}

	// Accept --username/--auth-pass as aliases for --db-user/--db-pass
	if dbUser == "" && authUser != "" {
		dbUser = authUser
	}
	if dbPass == "" && authPass != "" {
		dbPass = authPass
	}

	// Prompt for database credentials if not provided via flags
	if dbUser == "" || dbPass == "" {
		promptDBCredentials()
	}

	// === Phase 1: Silent setup ===
	restoreStderr := muteStderr()
	log.SetOutput(io.Discard)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	if invisible {
		visable = ""
	}

	PrintBanner()
	fmt.Fprintf(os.Stdout, "  \033[90mStarting %s plugin...\033[0m", info.Name)

	// Start the plugin process
	cfg := plugin.PluginConfig{
		Host:     dbIP,
		Port:     dbPort,
		User:     dbUser,
		Pass:     dbPass,
		Database: dbName,
		ReadOnly: dbReadOnly,
	}

	proc, pluginResp, err := plugin.StartPlugin(info, cfg)
	if err != nil {
		restoreStderr()
		fmt.Fprintf(os.Stdout, "\n\n  \033[1;31m✗ Plugin failed to start\033[0m\n")
		fmt.Fprintf(os.Stdout, "    %s\n\n", err)
		os.Exit(1)
	}
	defer proc.Stop()

	fmt.Fprintf(os.Stdout, "\r  \033[32m✓ Plugin ready: %s %s (port %d)\033[0m\n", pluginResp.Name, pluginResp.Version, pluginResp.WebPort)
	fmt.Fprintf(os.Stdout, "  \033[90mConnecting tunnel...\033[0m")

	// Now register the plugin's web port as an HTTP service through Shield tunnel
	creds, err := config.GetOrCreateCredentials()
	if err != nil {
		restoreStderr()
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	// Override protocol to http for the API call — the plugin serves HTTP
	origProtocol := protocol
	protocol = "http"
	defer func() { protocol = origProtocol }()

	// Call quick-setup with the plugin's web port as the target
	var resp *QuickSetupResponse
	maxAttempts := 5
	credReset := false
	pluginIP := "127.0.0.1"
	pluginPort := pluginResp.WebPort

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = callQuickSetup(pluginIP, pluginPort, creds)
		if err == nil {
			break
		}
		errMsg := err.Error()
		if strings.Contains(errMsg, "401") && !credReset {
			credReset = true
			os.Remove(config.GetCredentialFilePath())
			creds, _ = config.GetOrCreateCredentials()
			time.Sleep(2 * time.Second)
			continue
		}
		if strings.Contains(errMsg, "429") {
			if attempt < maxAttempts {
				time.Sleep(time.Duration(attempt*3) * time.Second)
				continue
			}
		}
		if strings.Contains(errMsg, "EOF") ||
			strings.Contains(errMsg, "timeout") ||
			strings.Contains(errMsg, "connection refused") {
			if attempt < maxAttempts {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
		}
		break
	}
	if err != nil {
		restoreStderr()
		fmt.Fprintf(os.Stdout, "\n\n  \033[1;31m✗ Connection failed\033[0m\n")
		fmt.Fprintf(os.Stdout, "    %s\n\n", friendlyError(err))
		os.Exit(1)
	}

	// Save credentials
	newCreds := &config.Credentials{
		ConnectorName: creds.ConnectorName,
		Password:      resp.Data.Connector.Password,
		ExternalIP:    resp.Data.Connector.ExternalIP,
		APIPort:       resp.Data.Connector.APIPort,
		TunnelPort:    tunnelPort,
		ConnUsername:   resp.Data.Connector.Username,
		ConnPassword:   resp.Data.Connector.Password,
	}

	localPort := creds.LocalPort
	if localPort > 0 {
		if ln, lnErr := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort)); lnErr != nil {
			localPort = 0
		} else {
			ln.Close()
		}
	}
	if localPort == 0 {
		localPort, err = findAvailablePort(4000, 5000)
		if err != nil {
			restoreStderr()
			return fmt.Errorf("failed to find available port: %w", err)
		}
	}

	newCreds.LocalPort = localPort
	newCreds.EncryptAndSave(config.GetCredentialFilePath())

	// Create tunnel
	connInfo := tunnel.ConnectionInfo{
		ExternalIP: resp.Data.Connector.ExternalIP,
		ServerPort: resp.Data.Connector.APIPort,
		TunnelPort: tunnelPort,
		Username:   resp.Data.Connector.Username,
		Password:   resp.Data.Connector.Password,
	}

	mgr := tunnel.NewManager(connInfo)

	resource := resp.Data.App.Resource
	resourceRemote := fmt.Sprintf("R:%d:%s:%d", resource.Port, pluginIP, pluginPort)

	err = mgr.CreateMainTunnel(resp.Data.Connector.APIPort, localPort, resourceRemote)
	if err != nil {
		restoreStderr()
		return fmt.Errorf("failed to create tunnel: %w", err)
	}

	siteURL := resp.Data.App.SiteURL
	activateTunnel(siteURL, 3, mgr)

	// === Phase 2: Clean header ===
	fmt.Print("\033[2J\033[H")
	PrintBanner()
	headerLines := 11

	p := func(format string, a ...any) {
		fmt.Printf(format, a...)
		fmt.Println()
		headerLines++
	}

	p("  \033[1;32m✓ Tunnel established successfully!\033[0m")
	p("")
	p("  \033[1;36mPlugin:\033[0m %s %s", pluginResp.Name, pluginResp.Version)
	p("  \033[1;36mTarget:\033[0m %s (%s:%d)", origProtocol, dbIP, dbPort)
	p("")
	p("  \033[1;33m⚡ Tunnel Mapping\033[0m")
	p("    \033[36mPlugin Web:\033[0m   localhost:%d  →  %s %s:%d", pluginPort, origProtocol, dbIP, dbPort)
	p("    \033[36mAPI Tunnel:\033[0m   remote:%d  ←→  localhost:%d", resp.Data.Connector.APIPort, localPort)
	p("    \033[36mApp Tunnel:\033[0m   remote:%d  ←→  localhost:%d", resource.Port, pluginPort)
	p("    \033[36mServer:\033[0m       %s:%d", resp.Data.Connector.ExternalIP, tunnelPort)
	p("")
	p("  \033[1;36mAccess URL:\033[0m")
	p("    %s", siteURL)
	p("")

	if invisible {
		apiKey := resp.Data.APIKey
		if apiKey.NHPServer != "" && apiKey.Code != "" {
			authURL := fmt.Sprintf("https://%s/plugins/auth?resid=%s&action=valid&format=redirect&passcode=%s",
				apiKey.NHPServer, apiKey.AppID, apiKey.Code,
			)
			p("  \033[1;36mAuth URL:\033[0m")
			p("    %s", authURL)
			p("")
		}
	}

	p("  \033[90mExpires: %s\033[0m", resp.Data.APIKey.ExpireTime)
	p("")
	p("  \033[90m──────────────────────────────────────────────────\033[0m")
	p("  \033[90mPress Ctrl+C to stop | Logs below ↓\033[0m")
	p("  \033[90m──────────────────────────────────────────────────\033[0m")

	// === Phase 3: Logs ===
	termHeight := getTermHeight()
	fmt.Printf("\033[%d;%dr", headerLines+1, termHeight)
	fmt.Printf("\033[%d;1H", headerLines+1)

	restoreStderr()
	if verbose {
		setupLog(slog.LevelDebug)
	} else {
		setupLog(slog.LevelInfo)
	}

	if !invisible && siteURL != "" {
		openBrowser(siteURL)
	}

	go startLocalAPI(localPort, mgr, connInfo)

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Printf("\033[r")
	fmt.Printf("\033[%d;1H", termHeight)
	fmt.Println("\033[33mShutting down plugin and tunnel...\033[0m")
	proc.Stop()
	mgr.CloseAll()
	return nil
}

// promptDBCredentials interactively asks for database username, password, and database name.
func promptDBCredentials() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("  \033[1;36m🔐 Database credentials (press Enter to skip)\033[0m")
	fmt.Println()

	if dbUser == "" {
		fmt.Print("  Username [root]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			dbUser = input
		} else {
			dbUser = "root"
		}
	}

	if dbPass == "" {
		fmt.Print("  Password: ")
		pwBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err == nil {
			dbPass = strings.TrimSpace(string(pwBytes))
		}
	}

	if dbName == "" {
		fmt.Print("  Database (optional): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			dbName = input
		}
	}

	fmt.Println()
	fmt.Printf("  \033[32m✓ Connecting as %s\033[0m\n", dbUser)
	if dbName != "" {
		fmt.Printf("  \033[32m  Database: %s\033[0m\n", dbName)
	}
	fmt.Println()
}
