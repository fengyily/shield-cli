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
	"golang.org/x/term"
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
	Reset         bool   `json:"reset,omitempty"`
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
			ID         int    `json:"id"`
			Code       string `json:"code"`
			NHPServer  string `json:"nhp_server"`
			KeyType    string `json:"key_type"`
			ExpireTime string `json:"expire_time"`
			AppID      string `json:"app_id"`
		} `json:"api_key"`
	} `json:"data"`
}

// muteStderr redirects os.Stderr to /dev/null to suppress chisel's direct writes.
// Returns a restore function.
func muteStderr() func() {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return func() {}
	}
	origStderr := os.Stderr
	os.Stderr = devNull
	return func() {
		os.Stderr = origStderr
		devNull.Close()
	}
}

func runShield(cmd *cobra.Command, args []string) error {
	// === Phase 1: Silent setup — suppress ALL output including chisel ===
	restoreStderr := muteStderr()
	log.SetOutput(io.Discard)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	// --invisible overrides visable to empty string
	if invisible {
		visable = ""
	}

	// Parse target address
	ip, port, err := parseTarget(target)
	if err != nil {
		restoreStderr()
		return fmt.Errorf("invalid target address: %w", err)
	}

	// Print banner while API call happens
	PrintBanner()
	fmt.Fprintf(os.Stdout, "  \033[90mConnecting...\033[0m")

	// Load or create credentials
	creds, err := config.GetOrCreateCredentials()
	if err != nil {
		restoreStderr()
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	// Call quick-setup API, auto-reset on auth failure, retry on transient errors
	var resp *QuickSetupResponse
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = callQuickSetup(ip, port, creds)
		if err == nil {
			break
		}
		// Auth failure: reset credentials and retry once
		if strings.Contains(err.Error(), "401") && creds.Password != "" {
			os.Remove(config.GetCredentialFilePath())
			creds, _ = config.GetOrCreateCredentials()
			continue
		}
		// Transient error (EOF, timeout, connection refused): retry
		errMsg := err.Error()
		if strings.Contains(errMsg, "EOF") ||
			strings.Contains(errMsg, "timeout") ||
			strings.Contains(errMsg, "connection refused") {
			if attempt < 3 {
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
		fmt.Fprintf(os.Stdout, "  \033[90mPlease check your network and try again.\033[0m\n\n")
		os.Exit(1)
	}

	// Save credentials from response
	newCreds := &config.Credentials{
		ConnectorName: resp.Data.Connector.Username,
		Password:      resp.Data.Connector.Password,
	}
	newCreds.EncryptAndSave(config.GetCredentialFilePath())

	// Find available local port
	localPort, err := findAvailablePort(4000, 5000)
	if err != nil {
		restoreStderr()
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

	// Create single chisel connection with both API tunnel and resource tunnel
	resource := resp.Data.App.Resource
	resourceRemote := fmt.Sprintf("R:%d:%s:%d", resource.Port, ip, port)

	err = mgr.CreateMainTunnel(resp.Data.Connector.APIPort, localPort, resourceRemote)
	if err != nil {
		restoreStderr()
		return fmt.Errorf("failed to create tunnel: %w", err)
	}

	// Activate tunnel (silent, no logs)
	siteURL := resp.Data.App.SiteURL
	activateTunnel(siteURL, 3, mgr)

	// === Phase 2: Clear screen and draw clean header ===
	fmt.Print("\033[2J\033[H") // Clear screen, cursor to top

	PrintBanner()
	headerLines := 11 // banner takes ~11 lines
	headerLines += printHeader(resp, resource.Port, ip, port)

	// === Phase 3: Set scroll region and enable logs ===
	termHeight := getTermHeight()
	fmt.Printf("\033[%d;%dr", headerLines+1, termHeight)
	fmt.Printf("\033[%d;1H", headerLines+1)

	// Restore stderr and enable log output
	restoreStderr()
	if verbose {
		setupLog(slog.LevelDebug)
	} else {
		setupLog(slog.LevelInfo)
	}

	fmt.Println("  \033[1;32m✓ Tunnel established successfully!\033[0m")
	fmt.Println()

	// Start local API server
	go startLocalAPI(localPort, mgr, connInfo)

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Reset scroll region and clean up
	fmt.Printf("\033[r")
	fmt.Printf("\033[%d;1H", termHeight)
	fmt.Println("\033[33mShutting down...\033[0m")
	mgr.CloseAll()
	return nil
}

// friendlyError converts technical errors into user-friendly messages.
func friendlyError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "EOF"):
		return "Server closed the connection unexpectedly (EOF)"
	case strings.Contains(msg, "timeout"):
		return "Connection timed out — server may be unreachable"
	case strings.Contains(msg, "connection refused"):
		return "Connection refused — server may be down"
	case strings.Contains(msg, "no such host"):
		return "DNS resolution failed — check the server URL"
	case strings.Contains(msg, "certificate"):
		return "TLS certificate error — check the server URL"
	case strings.Contains(msg, "429"):
		return "Rate limited — please wait a moment and try again"
	case strings.Contains(msg, "401"):
		return "Authentication failed — try: shield clean"
	default:
		return msg
	}
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
		Reset:         creds.Password == "",
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

	// Log full request body with password masked
	var maskedReq QuickSetupRequest
	maskedReq = reqBody
	maskedReq.Password = maskPassword(reqBody.Password)
	if maskedReq.AuthPass != "" {
		maskedReq.AuthPass = maskPassword(reqBody.AuthPass)
	}
	maskedJSON, _ := json.Marshal(maskedReq)
	slog.Debug("API request", "url", url, "body", string(maskedJSON))

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

// activateTunnel sends requests to warm up the tunnel, stops early if tunnel is connected.
func activateTunnel(siteURL string, times int, mgr *tunnel.Manager) {
	client := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < times; i++ {
		// Stop early if main tunnel is already connected
		if mgr.IsMainConnected() {
			return
		}
		resp, err := client.Get(siteURL)
		if err == nil {
			resp.Body.Close()
		}
		if i < times-1 {
			time.Sleep(1 * time.Second)
		}
	}
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

// getTermHeight returns the terminal height, defaults to 40 if detection fails.
func getTermHeight() int {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || h <= 0 {
		return 40
	}
	return h
}

// printHeader draws the fixed header area and returns the number of lines used.
func printHeader(resp *QuickSetupResponse, resourcePort int, targetIP string, targetPort int) int {
	lines := 0

	p := func(format string, a ...any) {
		fmt.Printf(format, a...)
		fmt.Println()
		lines++
	}

	p("  \033[1;32m✓ Tunnel established successfully!\033[0m")
	p("")
	p("  \033[1;33m⚡ Tunnel Mapping\033[0m")
	p("    \033[36mApp Tunnel:\033[0m   remote:%d  ←→  %s:%d", resourcePort, targetIP, targetPort)
	p("    \033[36mServer:\033[0m       %s:%d", resp.Data.Connector.ExternalIP, tunnelPort)
	p("")

	// Access URL
	p("  \033[1;36mAccess URL:\033[0m")
	p("    %s", resp.Data.App.SiteURL)
	p("")

	// Auth URL (only in invisible mode)
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

	p("  \033[90mProtocol: %s | Target: %s\033[0m", protocol, target)
	p("  \033[90mExpires: %s\033[0m", resp.Data.APIKey.ExpireTime)
	p("")
	p("  \033[90m──────────────────────────────────────────────────\033[0m")
	p("  \033[90mPress Ctrl+C to stop | Logs below ↓\033[0m")
	p("  \033[90m──────────────────────────────────────────────────\033[0m")

	return lines
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
			"code": 200,
			"data": tunnels,
		})
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	slog.Info("Local API server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Local API server error", "error", err)
	}
}
