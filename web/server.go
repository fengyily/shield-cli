package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"shield-cli/config"
)

//go:embed static
var staticFiles embed.FS

// Server represents the web management server
type Server struct {
	port    int
	store   *config.AppStore
	connMgr *ConnectionManager
}

// NewServer creates a new web server, pre-loading credentials at startup
func NewServer(port int) (*Server, error) {
	creds, err := config.GetOrCreateCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	return &Server{
		port:    port,
		store:   config.NewAppStore(),
		connMgr: NewConnectionManager(creds),
	}, nil
}

// Start starts the web server immediately, then establishes the main tunnel in the background
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/apps", s.handleApps)
	mux.HandleFunc("/api/apps/", s.handleAppByID)
	mux.HandleFunc("/api/rename/", s.handleRename)
	mux.HandleFunc("/api/connect/", s.handleConnect)
	mux.HandleFunc("/api/disconnect/", s.handleDisconnect)
	mux.HandleFunc("/api/status/", s.handleStatus)
	mux.HandleFunc("/api/tunnel", s.handleTunnelStatus)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	host := "127.0.0.1"
	if os.Getenv("SHIELD_LISTEN_HOST") != "" {
		host = os.Getenv("SHIELD_LISTEN_HOST")
	}
	addr := fmt.Sprintf("%s:%d", host, s.port)
	slog.Info("Web UI starting", "url", fmt.Sprintf("http://%s", addr))
	fmt.Printf("\n  Shield Web UI is running at:\n\n")
	fmt.Printf("    \033[1;36mhttp://%s\033[0m\n\n", addr)
	fmt.Printf("  \033[90mPress Ctrl+C to stop\033[0m\n\n")

	// Establish main tunnel in the background so Web UI is available immediately
	go func() {
		if err := s.connMgr.SetupMainTunnel(); err != nil {
			slog.Warn("Failed to establish main tunnel at startup, will retry on first connect", "error", err)
		}
	}()

	return http.ListenAndServe(addr, mux)
}

// Shutdown cleanly stops all connections
func (s *Server) Shutdown() {
	s.connMgr.DisconnectAll()
}

// handleApps handles GET (list) and POST (create) for /api/apps
func (s *Server) handleApps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		apps, err := s.store.List()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code": 500, "message": err.Error(),
			})
			return
		}
		// Attach connection status to each app
		type appWithStatus struct {
			config.AppConfig
			ConnStatus *ConnectResult `json:"conn_status,omitempty"`
		}
		result := make([]appWithStatus, len(apps))
		for i, app := range apps {
			result[i] = appWithStatus{AppConfig: app}
			if status := s.connMgr.GetStatus(app.ID); status != nil {
				result[i].ConnStatus = status
			}
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "data": result,
		})

	case http.MethodPost:
		// Check app count limit
		existingApps, err := s.store.List()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code": 500, "message": err.Error(),
			})
			return
		}
		if len(existingApps) >= 10 {
			writeJSON(w, http.StatusTooManyRequests, map[string]interface{}{
				"code": 429, "message": "Maximum 10 applications allowed. Please delete an app first.",
			})
			return
		}

		var app config.AppConfig
		if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code": 400, "message": "Invalid request body",
			})
			return
		}
		if app.Protocol == "" {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code": 400, "message": "protocol is required",
			})
			return
		}
		if app.IP == "" {
			app.IP = "127.0.0.1"
		}
		if app.Port == 0 {
			app.Port = defaultPort(app.Protocol)
		}
		created, err := s.store.Add(app)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code": 500, "message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "data": created,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleAppByID handles GET, PUT, DELETE for /api/apps/{id}
func (s *Server) handleAppByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Path[len("/api/apps/"):]
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "message": "app ID is required",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		app, err := s.store.Get(id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]interface{}{
				"code": 404, "message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "data": app,
		})

	case http.MethodPut:
		var app config.AppConfig
		if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code": 400, "message": "Invalid request body",
			})
			return
		}
		updated, err := s.store.Update(id, app)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]interface{}{
				"code": 404, "message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "data": updated,
		})

	case http.MethodDelete:
		// Disconnect if connected
		s.connMgr.Disconnect(id)
		if err := s.store.Delete(id); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]interface{}{
				"code": 404, "message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "message": "deleted",
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleRename handles PUT /api/rename/{id}
func (s *Server) handleRename(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/rename/"):]
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "message": "name is required",
		})
		return
	}

	app, err := s.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"code": 404, "message": err.Error(),
		})
		return
	}
	app.Name = body.Name
	updated, err := s.store.Update(id, *app)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500, "message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": 200, "data": updated,
	})
}

// handleConnect handles POST /api/connect/{id}
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/connect/"):]
	app, err := s.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"code": 404, "message": err.Error(),
		})
		return
	}

	// Check main tunnel readiness
	tunnelStatus, _ := s.connMgr.MainTunnelStatus()
	if tunnelStatus == "connecting" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"code": 503, "message": "Main tunnel is connecting, please wait...",
		})
		return
	}

	// Check concurrent connection limit
	if s.connMgr.ActiveCount() >= 3 {
		writeJSON(w, http.StatusTooManyRequests, map[string]interface{}{
			"code": 429, "message": "Maximum 3 concurrent connections allowed. Please disconnect an app first.",
		})
		return
	}

	params := ConnectParams{
		Protocol:    app.Protocol,
		IP:          app.IP,
		Port:        app.Port,
		Server:      app.Server,
		TunnelPort:  app.TunnelPort,
		Visable:     app.Visable,
		Invisible:   app.Invisible,
		Username:    app.Username,
		AuthPass:    app.AuthPass,
		PrivateKey:  app.PrivateKey,
		Passphrase:  app.Passphrase,
		EnableSftp:  app.EnableSftp,
		DisplayName: app.DisplayName,
		SiteName:    app.SiteName,
	}

	result, err := s.connMgr.Connect(id, params)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500, "message": err.Error(),
		})
		return
	}

	// Record last connected time
	_ = s.store.UpdateLastConnected(id)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": 200, "data": result,
	})
}

// handleDisconnect handles POST /api/disconnect/{id}
func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/disconnect/"):]
	s.connMgr.Disconnect(id)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": 200, "message": "disconnected",
	})
}

// handleStatus handles GET /api/status/{id}
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/status/"):]
	status := s.connMgr.GetStatus(id)
	if status == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": 200, "data": map[string]string{"status": "idle"},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": 200, "data": status,
	})
}

// handleTunnelStatus handles GET /api/tunnel — returns main tunnel readiness
func (s *Server) handleTunnelStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	status, errMsg := s.connMgr.MainTunnelStatus()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"status": status,
			"error":  errMsg,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func defaultPort(protocol string) int {
	ports := map[string]int{
		"ssh": 22, "http": 80, "https": 443,
		"rdp": 3389, "vnc": 5900, "telnet": 23,
	}
	if p, ok := ports[protocol]; ok {
		return p
	}
	return 80
}
