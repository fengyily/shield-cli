package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"shield-cli/config"
	"shield-cli/tunnel"

	"github.com/spf13/cobra"
)

// filterWriter filters out noisy chisel retry/denied log lines
type filterWriter struct {
	out io.Writer
}

func (f *filterWriter) Write(p []byte) (n int, err error) {
	line := string(p)
	// Filter out chisel retry spam and access denied noise
	if strings.Contains(line, "Retrying in") ||
		strings.Contains(line, "access to") ||
		strings.Contains(line, "Connection error: server:") {
		return len(p), nil // discard silently
	}
	return f.out.Write(p)
}

// setupLog configures slog and standard log based on verbose flag.
// In normal mode: only warnings and errors are shown during setup, info enabled after connection.
// In verbose mode: all levels shown from the start.
func setupLog(level slog.Level) {
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(handler))

	if level <= slog.LevelDebug {
		log.SetOutput(&filterWriter{out: os.Stderr})
	} else {
		// Suppress chisel's internal log.Printf in non-verbose mode
		log.SetOutput(io.Discard)
	}
}

// maskPassword masks a password string for safe logging
func maskPassword(pw string) string {
	if len(pw) <= 4 {
		return "****"
	}
	return pw[:2] + "****" + pw[len(pw)-2:]
}

// API request/response types
type QuickSetupRequest struct {
	Protocol      string `json:"protocol"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	ConnectorName string `json:"connector_name"`
	Password      string `json:"password"`
	DisplayName   string `json:"display_name,omitempty"`
	SiteName      string `json:"site_name,omitempty"`
	Visable       string `json:"visable,omitempty"`
	Username      string `json:"username,omitempty"`
	AuthPass      string `json:"auth_pass,omitempty"`
	PrivateKey    string `json:"private_key,omitempty"`
	Passphrase    string `json:"passphrase,omitempty"`
	EnableSftp    bool   `json:"enable_sftp,omitempty"`
}

type QuickSetupResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Connector struct {
			ID            int    `json:"id"`
			ConnectorName string `json:"connector_name"`
			DisplayName   string `json:"display_name"`
			APIPort       int    `json:"api_port"`
			ExternalIP    string `json:"external_ip"`
			Username      string `json:"username"`
			Password      string `json:"password"`
		} `json:"connector"`
		App struct {
			ID       int    `json:"id"`
			AppID    string `json:"app_id"`
			SiteName string `json:"site_name"`
			SiteURL  string `json:"site_url"`
			Protocol string `json:"protocol"`
			Resource struct {
				IP       string `json:"ip"`
				Port     int    `json:"port"`
				AcID     string `json:"ac_id"`
				Hostname string `json:"hostname"`
				Maskhost bool   `json:"maskhost"`
				Protocol string `json:"protocol"`
			} `json:"resource"`
		} `json:"app"`
		APIKey struct {
			ID        int    `json:"id"`
			Code      string `json:"code"`
			NHPServer string `json:"nhp_server"`
			KeyType   string `json:"key_type"`
			ExpireTime string `json:"expire_time"`
			AppID     string `json:"app_id"`
		} `json:"api_key"`
	} `json:"data"`
}

func runShield(cmd *cobra.Command, args []string) error {
	// === Phase 1: Setup ===
	if verbose {
		setupLog(slog.LevelDebug)
	} else {
		setupLog(slog.LevelWarn)
	}

	PrintBanner()
	fmt.Println("  \033[90mConnecting...\033[0m")

	// Parse target address
	ip, port, err := parseTarget(target)
	if err != nil {
		return fmt.Errorf("invalid target address: %w", err)
	}

	// Load or create credentials
	creds, err := config.GetOrCreateCredentials()
	if err != nil {
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	// Call quick-setup API
	resp, err := callQuickSetup(ip, port, creds)
	if err != nil {
		return fmt.Errorf("failed to call quick-setup API: %w", err)
	}

	// Save credentials from response
	newCreds := &config.Credentials{
		ConnectorName: resp.Data.Connector.Username,
		Password:      resp.Data.Connector.Password,
	}
	credPath := config.GetCredentialFilePath()
	newCreds.EncryptAndSave(credPath) // best-effort save

	// Find available local port
	localPort, err := findAvailablePort(4000, 5000)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	// Create tunnel manager
	connInfo := tunnel.ConnectionInfo{
		ExternalIP: resp.Data.Connector.ExternalIP,
		ServerPort: resp.Data.Connector.APIPort,
		TunnelPort: tunnelPort,
		Username:   resp.Data.Connector.Username,
		Password:   resp.Data.Connector.Password,
	}

	mgr := tunnel.NewManager(connInfo)

	// Create main tunnel: maps remote api_port to local API server
	err = mgr.CreateMainTunnel(resp.Data.Connector.APIPort, localPort)
	if err != nil {
		return fmt.Errorf("failed to create main tunnel: %w", err)
	}

	// Create resource tunnel: maps remote resource port to local target
	resource := resp.Data.App.Resource
	mgr.CreateDynamicTunnel(
		strconv.Itoa(resource.Port),
		ip,
		strconv.Itoa(port),
	)

	// === Phase 2: Print tunnel mapping & connection info ===
	fmt.Println()
	fmt.Printf("  \033[1;33m⚡ Tunnel Mapping\033[0m\n")
	fmt.Printf("    \033[36mApp Tunnel:\033[0m   remote:%d  ←→  %s:%d\n", resource.Port, ip, port)
	fmt.Printf("    \033[36mServer:\033[0m       %s:%d\n", resp.Data.Connector.ExternalIP, tunnelPort)
	fmt.Println()

	printAccessInfo(resp)

	// === Phase 3: Enable info-level logs — scroll below from here ===
	if verbose {
		setupLog(slog.LevelDebug)
	} else {
		setupLog(slog.LevelInfo)
	}

	// Start local API server
	go startLocalAPI(localPort, mgr, connInfo)

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n\033[33mShutting down...\033[0m")
	mgr.CloseAll()
	return nil
}

func parseTarget(target string) (string, int, error) {
	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("expected format ip:port, got %q", target)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port %q: %w", parts[1], err)
	}

	return parts[0], port, nil
}

func callQuickSetup(ip string, port int, creds *config.Credentials) (*QuickSetupResponse, error) {
	reqBody := QuickSetupRequest{
		Protocol:      protocol,
		IP:            ip,
		Port:          port,
		ConnectorName: creds.ConnectorName,
		Password:      creds.Password,
		DisplayName:   displayName,
		SiteName:      siteName,
		Visable:       visable,
		Username:      authUser,
		AuthPass:      authPass,
		PrivateKey:    privateKey,
		Passphrase:    passphrase,
		EnableSftp:    enableSftp,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(apiServer, "/") + "/api/public/quick-setup"
	slog.Debug("API request", "url", url, "connector", creds.ConnectorName, "password", maskPassword(creds.Password))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Mask password in response body for logging
	var maskedBody json.RawMessage
	{
		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err == nil {
			if data, ok := raw["data"].(map[string]interface{}); ok {
				if conn, ok := data["connector"].(map[string]interface{}); ok {
					if pw, ok := conn["password"].(string); ok {
						conn["password"] = maskPassword(pw)
					}
				}
			}
			maskedBody, _ = json.Marshal(raw)
		}
	}
	if maskedBody != nil {
		slog.Debug("API response", "status", resp.StatusCode, "body", string(maskedBody))
	} else {
		slog.Debug("API response", "status", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result QuickSetupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Accept code=200 or code=0 (some APIs return 0 for success)
	if result.Code != 200 && result.Code != 0 {
		return nil, fmt.Errorf("API error (code=%d): %s", result.Code, result.Message)
	}

	slog.Debug("API parsed",
		"connector", result.Data.Connector.ConnectorName,
		"external_ip", result.Data.Connector.ExternalIP,
		"api_port", result.Data.Connector.APIPort,
		"site_url", result.Data.App.SiteURL,
	)

	return &result, nil
}

func findAvailablePort(min, max int) (int, error) {
	for port := min; port <= max; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port in range %d-%d", min, max)
}

func printAccessInfo(resp *QuickSetupResponse) {
	fmt.Println()
	fmt.Println("  \033[1;32m✓ Tunnel established successfully!\033[0m")
	fmt.Println()

	// Print site URL
	fmt.Printf("  \033[1;36mSite URL:\033[0m\n")
	fmt.Printf("    %s\n", resp.Data.App.SiteURL)
	fmt.Println()

	// Build and print access URL from api_key
	apiKey := resp.Data.APIKey
	if apiKey.NHPServer != "" && apiKey.Code != "" {
		accessURL := fmt.Sprintf("https://%s/plugins/auth?resid=%s&action=valid&format=redirect&passcode=%s",
			apiKey.NHPServer,
			apiKey.AppID,
			apiKey.Code,
		)
		fmt.Printf("  \033[1;36mAccess URL:\033[0m\n")
		fmt.Printf("    %s\n", accessURL)
		fmt.Println()
	}

	fmt.Printf("  \033[90mProtocol: %s | Target: %s\033[0m\n", protocol, target)
	fmt.Printf("  \033[90mExpires: %s\033[0m\n", resp.Data.APIKey.ExpireTime)
	fmt.Println()
	fmt.Println("  \033[90m──────────────────────────────────────────────────\033[0m")
	fmt.Println()
	fmt.Println("  \033[90mPress Ctrl+C to stop\033[0m")
	fmt.Println()
}

func startLocalAPI(port int, mgr *tunnel.Manager, connInfo tunnel.ConnectionInfo) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})

	mux.HandleFunc("/connector", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			rport := r.URL.Query().Get("rport")
			lip := r.URL.Query().Get("lip")
			lport := r.URL.Query().Get("lport")

			if rport == "" || lip == "" || lport == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    400,
					"message": "rport, lip, lport are required",
				})
				return
			}

			slog.Info("Dynamic tunnel request", "rport", rport, "lip", lip, "lport", lport)

			err := mgr.CreateDynamicTunnel(rport, lip, lport)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    500,
					"message": fmt.Sprintf("Failed to create tunnel: %v", err),
				})
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    200,
				"message": fmt.Sprintf("Tunnel created: R:%s:%s:%s", rport, lip, lport),
			})

		case http.MethodDelete:
			rport := r.URL.Query().Get("rport")
			if rport == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    400,
					"message": "rport is required",
				})
				return
			}

			mgr.CloseTunnel(rport)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    200,
				"message": "Tunnel closed",
			})

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/connectors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tunnels := mgr.List()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"data":    tunnels,
		})
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	slog.Info("Local API server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Local API server error", "error", err)
	}
}
