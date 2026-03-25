package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"shield-cli/config"
	"shield-cli/plugin"
	"shield-cli/tunnel"
)

// ConnectParams holds the parameters needed to establish a connection
type ConnectParams struct {
	Protocol    string `json:"protocol"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Server      string `json:"server"`
	TunnelPort  int    `json:"tunnel_port"`
	Visable     string `json:"visable"`
	Invisible   bool   `json:"invisible"`
	Username    string `json:"username"`
	AuthPass    string `json:"auth_pass"`
	PrivateKey  string `json:"private_key"`
	Passphrase  string `json:"passphrase"`
	EnableSftp  bool   `json:"enable_sftp"`
	DisplayName string `json:"display_name"`
	SiteName    string `json:"site_name"`
	DBUser      string `json:"db_user"`
	DBPass      string `json:"db_pass"`
	DBName      string `json:"db_name"`
	DBReadOnly  bool   `json:"db_readonly"`
}

// ConnectResult holds the result of a connection attempt
type ConnectResult struct {
	Status       string `json:"status"` // connecting, connected, failed, disconnected
	SiteURL      string `json:"site_url,omitempty"`
	LocalURL     string `json:"local_url,omitempty"` // local plugin web UI URL (e.g. http://127.0.0.1:port)
	AuthURL      string `json:"auth_url,omitempty"`
	Error        string `json:"error,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	ResourcePort int    `json:"resource_port,omitempty"` // remote resource port for app tunnel
}

// ActiveConnection tracks an active tunnel session (app-level, resource tunnel only)
type ActiveConnection struct {
	AppConfigID  string
	ResourcePort string // remote resource port key for CloseTunnel
	Result       ConnectResult
	StopCh       chan struct{}
	PluginProc   *plugin.Process // non-nil if this connection uses a plugin
}

// ConnectionManager manages the shared main tunnel and per-app resource tunnels
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*ActiveConnection // keyed by app config ID

	credsMu sync.RWMutex
	creds   *config.Credentials

	// Shared main tunnel (established once at startup)
	mainMu       sync.Mutex
	mainMgr      *tunnel.Manager
	mainReady    bool
	mainError    string // last error message from tunnel setup
	mainStarting bool   // true while SetupMainTunnel is in progress
	localPort    int
}

// NewConnectionManager creates a new ConnectionManager with pre-loaded credentials
func NewConnectionManager(creds *config.Credentials) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ActiveConnection),
		creds:       creds,
	}
}

// GetCreds returns the shared credentials
func (cm *ConnectionManager) GetCreds() *config.Credentials {
	cm.credsMu.RLock()
	defer cm.credsMu.RUnlock()
	return cm.creds
}

// RefreshCreds reloads credentials from disk
func (cm *ConnectionManager) RefreshCreds() *config.Credentials {
	cm.credsMu.Lock()
	defer cm.credsMu.Unlock()
	newCreds, err := config.GetOrCreateCredentials()
	if err == nil {
		cm.creds = newCreds
	}
	return cm.creds
}

// UpdateCreds updates the shared credentials
func (cm *ConnectionManager) UpdateCreds(creds *config.Credentials) {
	cm.credsMu.Lock()
	defer cm.credsMu.Unlock()
	cm.creds = creds
}

// MainTunnelStatus returns the current main tunnel status and any error message.
// Status: "connected", "connecting", "disconnected"
func (cm *ConnectionManager) MainTunnelStatus() (status string, errMsg string) {
	cm.mainMu.Lock()
	defer cm.mainMu.Unlock()

	if cm.mainReady {
		return "connected", ""
	}
	if cm.mainStarting {
		return "connecting", ""
	}
	return "disconnected", cm.mainError
}

// SetupMainTunnel establishes the shared main tunnel (Server + API Tunnel).
// Called once at startup if connector info is available in credentials.
func (cm *ConnectionManager) SetupMainTunnel() error {
	cm.mainMu.Lock()

	if cm.mainReady {
		cm.mainMu.Unlock()
		return nil
	}

	creds := cm.GetCreds()
	if !creds.HasConnectorInfo() {
		slog.Info("No saved connector info, main tunnel will be established on first app connect")
		cm.mainMu.Unlock()
		return nil
	}

	cm.mainStarting = true
	cm.mainError = ""
	cm.mainMu.Unlock()

	err := cm.setupMainTunnelWithInfo(creds)

	cm.mainMu.Lock()
	cm.mainStarting = false
	if err != nil {
		cm.mainError = err.Error()
	}
	cm.mainMu.Unlock()

	return err
}

// setupMainTunnelWithInfo does the actual main tunnel setup.
// Callers must set mainStarting=true before calling and reset it after.
func (cm *ConnectionManager) setupMainTunnelWithInfo(creds *config.Credentials) error {
	// Determine local port
	localPort := creds.LocalPort
	if localPort > 0 {
		if ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort)); err != nil {
			localPort = 0
		} else {
			ln.Close()
		}
	}
	if localPort == 0 {
		var err error
		localPort, err = findAvailablePort(4000, 5000)
		if err != nil {
			return fmt.Errorf("no available local port: %w", err)
		}
	}

	connInfo := tunnel.ConnectionInfo{
		ExternalIP: creds.ExternalIP,
		ServerPort: creds.APIPort,
		TunnelPort: creds.TunnelPort,
		Username:   creds.ConnUsername,
		Password:   creds.ConnPassword,
	}

	mgr := tunnel.NewManager(connInfo)
	// Main tunnel: only API port mapping, no resource remotes
	if err := mgr.CreateMainTunnel(creds.APIPort, localPort); err != nil {
		return fmt.Errorf("failed to create main tunnel: %w", err)
	}

	// Start local API server first so it's ready when chisel connects
	go startConnLocalAPI(localPort)

	// Wait for chisel to connect — poll local health endpoint through the tunnel is not
	// possible from client side, so give chisel time to establish the connection.
	// Chisel typically connects within 3-5 seconds.
	time.Sleep(5 * time.Second)

	cm.mainMu.Lock()
	cm.mainMgr = mgr
	cm.mainReady = true
	cm.localPort = localPort
	cm.mainMu.Unlock()

	// Save localPort back to credentials
	creds.LocalPort = localPort
	creds.EncryptAndSave(config.GetCredentialFilePath())
	cm.UpdateCreds(creds)

	slog.Info("Main tunnel ready",
		"server", fmt.Sprintf("%s:%d", creds.ExternalIP, creds.TunnelPort),
		"api_tunnel", fmt.Sprintf("remote:%d ←→ localhost:%d", creds.APIPort, localPort),
	)
	return nil
}

// ensureMainTunnel ensures the main tunnel is ready, bootstrapping from quick-setup response if needed.
func (cm *ConnectionManager) ensureMainTunnel(resp *QuickSetupResponse, tunnelPort int) error {
	cm.mainMu.Lock()

	if cm.mainReady {
		cm.mainMu.Unlock()
		return nil
	}

	// Bootstrap: save connector info to credentials and establish main tunnel
	creds := cm.GetCreds()
	creds.ExternalIP = resp.Data.Connector.ExternalIP
	creds.APIPort = resp.Data.Connector.APIPort
	creds.TunnelPort = tunnelPort
	creds.ConnUsername = resp.Data.Connector.Username
	creds.ConnPassword = resp.Data.Connector.Password
	creds.ConnectorName = resp.Data.Connector.Username
	creds.Password = resp.Data.Connector.Password

	cm.mainStarting = true
	cm.mainError = ""
	cm.mainMu.Unlock()

	err := cm.setupMainTunnelWithInfo(creds)

	cm.mainMu.Lock()
	cm.mainStarting = false
	if err != nil {
		cm.mainError = err.Error()
	}
	cm.mainMu.Unlock()

	return err
}

// Connect establishes a tunnel connection for the given app config
func (cm *ConnectionManager) Connect(appID string, params ConnectParams) (*ConnectResult, error) {
	cm.mu.Lock()
	// If already connected, return existing result
	if existing, ok := cm.connections[appID]; ok {
		if existing.Result.Status == "connected" || existing.Result.Status == "connecting" {
			cm.mu.Unlock()
			return &existing.Result, nil
		}
		// Clean up failed/disconnected connection
		cm.disconnectLocked(appID)
	}

	conn := &ActiveConnection{
		AppConfigID: appID,
		Result: ConnectResult{
			Status: "connecting",
		},
		StopCh: make(chan struct{}),
	}
	cm.connections[appID] = conn
	cm.mu.Unlock()

	// Run connection in background
	go cm.doConnect(appID, params, conn)

	return &conn.Result, nil
}

// GetStatus returns the current status of a connection
func (cm *ConnectionManager) GetStatus(appID string) *ConnectResult {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conn, ok := cm.connections[appID]; ok {
		return &conn.Result
	}
	return nil
}

// Disconnect stops an active connection
func (cm *ConnectionManager) Disconnect(appID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.disconnectLocked(appID)
}

// DisconnectAll stops all active connections and closes the main tunnel
func (cm *ConnectionManager) DisconnectAll() {
	cm.mu.Lock()
	for id := range cm.connections {
		cm.disconnectLocked(id)
	}
	cm.mu.Unlock()

	cm.mainMu.Lock()
	if cm.mainMgr != nil {
		cm.mainMgr.CloseAll()
		cm.mainMgr = nil
		cm.mainReady = false
	}
	cm.mainMu.Unlock()
}

// ActiveCount returns the number of connections with status "connected" or "connecting"
func (cm *ConnectionManager) ActiveCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	count := 0
	for _, conn := range cm.connections {
		if conn.Result.Status == "connected" || conn.Result.Status == "connecting" {
			count++
		}
	}
	return count
}

func (cm *ConnectionManager) disconnectLocked(appID string) {
	if conn, ok := cm.connections[appID]; ok {
		close(conn.StopCh)
		// Stop plugin process if any
		if conn.PluginProc != nil {
			conn.PluginProc.Stop()
		}
		// Close only the resource tunnel for this app, not the main tunnel
		if conn.ResourcePort != "" && cm.mainMgr != nil {
			cm.mainMgr.CloseTunnel(conn.ResourcePort)
		}
		conn.Result.Status = "disconnected"
		delete(cm.connections, appID)
	}
}

// QuickSetupRequest mirrors the API request format
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

// QuickSetupResponse mirrors the API response format
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

func (cm *ConnectionManager) doConnect(appID string, params ConnectParams, conn *ActiveConnection) {
	slog.Info("Connecting app", "appID", appID, "protocol", params.Protocol, "target", fmt.Sprintf("%s:%d", params.IP, params.Port))

	// Check if this is a plugin protocol — if so, start the plugin first
	origProtocol := params.Protocol
	var pluginProc *plugin.Process
	var localURL string
	pluginIP := params.IP
	pluginPort := params.Port

	if info := findPluginForProtocol(params.Protocol); info != nil {
		slog.Info("Starting plugin for protocol", "protocol", params.Protocol, "plugin", info.Name)

		// Use DB credentials, falling back to auth credentials
		dbUser := params.DBUser
		dbPass := params.DBPass
		if dbUser == "" {
			dbUser = params.Username
		}
		if dbPass == "" {
			dbPass = params.AuthPass
		}

		cfg := plugin.PluginConfig{
			Host:     params.IP,
			Port:     params.Port,
			User:     dbUser,
			Pass:     dbPass,
			Database: params.DBName,
			ReadOnly: params.DBReadOnly,
		}
		proc, resp, err := plugin.StartPlugin(info, cfg)
		if err != nil {
			cm.setError(appID, fmt.Sprintf("Plugin failed to start: %v", err))
			return
		}
		pluginProc = proc
		conn.PluginProc = proc

		slog.Info("Plugin ready", "name", resp.Name, "version", resp.Version, "web_port", resp.WebPort)

		// Override: register the plugin's web port as HTTP
		params.Protocol = "http"
		pluginIP = "127.0.0.1"
		pluginPort = resp.WebPort
		params.IP = pluginIP
		params.Port = pluginPort
		localURL = fmt.Sprintf("http://127.0.0.1:%d", resp.WebPort)
		_ = origProtocol
		_ = pluginProc
	}

	// Use shared credentials
	creds := cm.GetCreds()

	// Call quick-setup API with retry (handles 401, 429, transient errors)
	slog.Debug("Calling quick-setup API", "server", params.Server)
	var err error
	var resp *QuickSetupResponse
	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = callQuickSetupAPI(params, creds)
		if err == nil {
			break
		}
		slog.Debug("Quick-setup API attempt failed", "attempt", attempt, "error", err)
		errMsg := err.Error()
		// Auth failure: refresh shared credentials and retry
		if strings.Contains(errMsg, "401") && creds.Password != "" {
			creds = cm.RefreshCreds()
			continue
		}
		// Rate limited: wait longer and retry
		if strings.Contains(errMsg, "429") {
			if attempt < maxAttempts {
				time.Sleep(time.Duration(attempt*3) * time.Second)
				continue
			}
		}
		// Transient errors: short retry
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
		slog.Error("Quick-setup API failed", "appID", appID, "error", err)
		cm.setError(appID, friendlyError(err))
		return
	}
	slog.Debug("Quick-setup API success",
		"connector", resp.Data.Connector.ConnectorName,
		"external_ip", resp.Data.Connector.ExternalIP,
		"api_port", resp.Data.Connector.APIPort,
		"site_url", resp.Data.App.SiteURL,
	)

	// Update credentials from response
	newCreds := cm.GetCreds()
	newCreds.ConnectorName = resp.Data.Connector.Username
	newCreds.Password = resp.Data.Connector.Password
	newCreds.ExternalIP = resp.Data.Connector.ExternalIP
	newCreds.APIPort = resp.Data.Connector.APIPort
	newCreds.TunnelPort = params.TunnelPort
	newCreds.ConnUsername = resp.Data.Connector.Username
	newCreds.ConnPassword = resp.Data.Connector.Password
	newCreds.EncryptAndSave(config.GetCredentialFilePath())
	cm.UpdateCreds(newCreds)

	// Ensure main tunnel is established (no-op if already ready)
	if err := cm.ensureMainTunnel(resp, params.TunnelPort); err != nil {
		cm.setError(appID, fmt.Sprintf("Failed to establish main tunnel: %v", err))
		return
	}

	// Add resource tunnel via dynamic tunnel on the shared main manager
	resource := resp.Data.App.Resource
	rport := fmt.Sprintf("%d", resource.Port)
	slog.Info("Creating app tunnel", "remote_port", rport, "target", fmt.Sprintf("%s:%d", params.IP, params.Port))
	if err := cm.mainMgr.CreateDynamicTunnel(rport, params.IP, fmt.Sprintf("%d", params.Port)); err != nil {
		slog.Error("Failed to create app tunnel", "appID", appID, "error", err)
		cm.setError(appID, fmt.Sprintf("Failed to create app tunnel: %v", err))
		return
	}

	// Poll site URL to activate the server-side route and wait for tunnel to be ready.
	// The server needs an incoming request on the site URL to set up routing through the tunnel.
	siteURL := resp.Data.App.SiteURL
	slog.Info("Activating app tunnel", "site_url", siteURL)
	if !waitForSiteReady(siteURL, 30*time.Second) {
		cm.setError(appID, "App tunnel activation timeout — site URL not reachable")
		return
	}
	cm.mainMgr.SetConnected(rport)
	slog.Info("App tunnel activated", "site_url", siteURL)

	// Build auth URL if invisible mode
	var authURL string
	if params.Invisible {
		apiKey := resp.Data.APIKey
		if apiKey.NHPServer != "" && apiKey.Code != "" {
			authURL = fmt.Sprintf("https://%s/plugins/auth?resid=%s&action=valid&format=redirect&passcode=%s",
				apiKey.NHPServer, apiKey.AppID, apiKey.Code,
			)
		}
	}

	// Update status — only now is the tunnel truly ready
	cm.mu.Lock()
	if c, ok := cm.connections[appID]; ok && c == conn {
		c.ResourcePort = rport
		c.Result = ConnectResult{
			Status:       "connected",
			SiteURL:      siteURL,
			LocalURL:     localURL,
			AuthURL:      authURL,
			AppID:        resp.Data.App.AppID,
			ResourcePort: resource.Port,
		}
	}
	cm.mu.Unlock()

	slog.Info("App connected",
		"appID", appID,
		"app_tunnel", fmt.Sprintf("remote:%s ←→ %s:%d", rport, params.IP, params.Port),
		"site_url", siteURL,
	)

	// Wait for stop signal
	<-conn.StopCh
	slog.Info("App disconnected", "appID", appID)
}

func (cm *ConnectionManager) setError(appID string, errMsg string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[appID]; ok {
		conn.Result.Status = "failed"
		conn.Result.Error = errMsg
	}
}

func maskPassword(pw string) string {
	if len(pw) <= 4 {
		return "****"
	}
	return pw[:2] + "****" + pw[len(pw)-2:]
}

func callQuickSetupAPI(params ConnectParams, creds *config.Credentials) (*QuickSetupResponse, error) {
	reqBody := QuickSetupRequest{
		Protocol:      params.Protocol,
		IP:            params.IP,
		Port:          params.Port,
		ConnectorName: creds.ConnectorName,
		Password:      creds.Password,
		Reset:         creds.Password == "",
		DisplayName:   params.DisplayName,
		SiteName:      params.SiteName,
		Visable:       params.Visable,
		Username:      params.Username,
		AuthPass:      params.AuthPass,
		PrivateKey:    params.PrivateKey,
		Passphrase:    params.Passphrase,
		EnableSftp:    params.EnableSftp,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("Failed to marshal quick-setup request", "error", err)
		return nil, err
	}

	url := strings.TrimRight(params.Server, "/") + "/api/public/quick-setup"

	slog.Debug("Quick-setup request",
		"url", url,
		"protocol", reqBody.Protocol,
		"ip", reqBody.IP,
		"port", reqBody.Port,
		"connector", reqBody.ConnectorName,
		"password", maskPassword(reqBody.Password),
		"reset", reqBody.Reset,
	)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		slog.Error("Quick-setup HTTP failed", "url", url, "error", err)
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	slog.Debug("Quick-setup response", "status", resp.StatusCode, "body_len", len(body))

	if resp.StatusCode != http.StatusOK {
		slog.Error("Quick-setup API error", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result QuickSetupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		slog.Error("Quick-setup parse failed", "error", err, "body", string(body))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 200 && result.Code != 0 {
		slog.Error("Quick-setup biz error", "code", result.Code, "message", result.Message)
		return nil, fmt.Errorf("API error (code=%d): %s", result.Code, result.Message)
	}

	slog.Info("Quick-setup success",
		"connector", result.Data.Connector.ConnectorName,
		"external_ip", result.Data.Connector.ExternalIP,
		"api_port", result.Data.Connector.APIPort,
		"app_id", result.Data.App.AppID,
		"site_url", result.Data.App.SiteURL,
		"resource_port", result.Data.App.Resource.Port,
	)

	return &result, nil
}

// tunnelCheckResponse is the response from the tunnel health-check API
type tunnelCheckResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// waitForSiteReady polls the site's tunnel API until it confirms the tunnel is active.
// POST {siteURL}/_webgate/api/tunnel with empty JSON body; code=0 means ready.
func waitForSiteReady(siteURL string, timeout time.Duration) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	apiURL := strings.TrimRight(siteURL, "/") + "/_webgate/api/tunnel"
	deadline := time.Now().Add(timeout)

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		resp, err := client.Post(apiURL, "application/json", bytes.NewReader([]byte("{}")))
		if err != nil {
			slog.Debug("Tunnel API not ready, retrying", "url", apiURL, "attempt", attempt, "error", err)
			time.Sleep(2 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result tunnelCheckResponse
		if err := json.Unmarshal(body, &result); err == nil && result.Code == 0 {
			slog.Debug("Tunnel API ready", "url", apiURL, "attempt", attempt)
			return true
		}
		slog.Debug("Tunnel API not ready, retrying", "url", apiURL, "attempt", attempt, "status", resp.StatusCode, "body", string(body))
		time.Sleep(2 * time.Second)
	}
	return false
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

func friendlyError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "EOF"):
		return "Server closed the connection unexpectedly"
	case strings.Contains(msg, "timeout"):
		return "Connection timed out"
	case strings.Contains(msg, "connection refused"):
		return "Connection refused — server may be down"
	case strings.Contains(msg, "no such host"):
		return "DNS resolution failed"
	case strings.Contains(msg, "429"):
		return "Rate limited — please wait"
	case strings.Contains(msg, "401"):
		return "Authentication failed — try clearing credentials"
	default:
		return msg
	}
}

// findPluginForProtocol returns the installed plugin info if the protocol is handled by a plugin.
func findPluginForProtocol(proto string) *plugin.PluginInfo {
	reg, err := plugin.LoadRegistry()
	if err != nil {
		return nil
	}
	info := reg.Find(proto)
	return info
}

func startConnLocalAPI(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	slog.Info("Local API starting", "addr", addr)
	http.ListenAndServe(addr, mux)
}
