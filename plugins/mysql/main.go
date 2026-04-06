package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed static/*
var staticFS embed.FS

// Protocol types (matches shield-cli/plugin package)

type StartRequest struct {
	Action string       `json:"action"`
	Config PluginConfig `json:"config,omitempty"`
}

type PluginConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user,omitempty"`
	Pass     string `json:"pass,omitempty"`
	Database string `json:"database,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`
}

type StartResponse struct {
	Status  string `json:"status"`
	WebPort int    `json:"web_port,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
	// If DB_HOST is set, run in standalone/Docker mode
	if os.Getenv("DB_HOST") != "" {
		standaloneMode()
		return
	}

	// Otherwise, use Shield CLI plugin protocol (stdin JSON)
	decoder := json.NewDecoder(os.Stdin)

	for {
		var req StartRequest
		if err := decoder.Decode(&req); err != nil {
			// stdin closed, exit
			return
		}

		switch req.Action {
		case "start":
			handleStart(req.Config)
		case "stop":
			os.Exit(0)
		}
	}
}

func respond(resp StartResponse) {
	json.NewEncoder(os.Stdout).Encode(resp)
}

func respondError(msg string) {
	respond(StartResponse{Status: "error", Message: msg})
}

// setupHTTP creates the HTTP mux with all API routes and static files.
func setupHTTP(db *sql.DB, cfg PluginConfig, hub *CollabHub) http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/databases", databasesHandler(db))
	mux.HandleFunc("/api/tables", tablesHandler(db))
	mux.HandleFunc("/api/schema", schemaHandler(db))
	mux.HandleFunc("/api/indexes", indexesHandler(db))
	mux.HandleFunc("/api/query", queryHandler(db, cfg.ReadOnly))
	mux.HandleFunc("/api/info", infoHandler(db, cfg))
	mux.HandleFunc("/api/er", erHandler(db))
	mux.HandleFunc("/api/export", exportSQLHandler(db))

	// WebSocket for ER collaboration
	mux.HandleFunc("/ws/er", collabHandler(hub))

	// Static files
	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/", http.FileServer(http.FS(staticSub)))

	return mux
}

// connectDB builds a DSN from the config, opens, and pings the database.
func connectDB(cfg PluginConfig) (*sql.DB, error) {
	user := cfg.User
	if user == "" {
		user = "root"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		user, cfg.Pass, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to MySQL at %s:%d: %w", cfg.Host, cfg.Port, err)
	}
	return db, nil
}

func handleStart(cfg PluginConfig) {
	db, err := connectDB(cfg)
	if err != nil {
		respondError(err.Error())
		return
	}

	// Find available port on loopback
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		respondError(fmt.Sprintf("failed to find available port: %v", err))
		return
	}
	webPort := listener.Addr().(*net.TCPAddr).Port

	// Respond with ready
	respond(StartResponse{
		Status:  "ready",
		WebPort: webPort,
		Name:    "MySQL Web Client",
		Version: "0.1.0",
	})

	hub := newCollabHub()
	go hub.run()

	// Start HTTP server in background
	go func() {
		if err := http.Serve(listener, setupHTTP(db, cfg, hub)); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for stop signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	db.Close()
}

// standaloneMode runs the plugin as a standalone server, reading config from
// environment variables. Intended for Docker / direct execution.
func standaloneMode() {
	port := 3306
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	readOnly := false
	if v := os.Getenv("DB_READONLY"); v == "true" || v == "1" {
		readOnly = true
	}
	webPort := 8080
	if v := os.Getenv("WEB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			webPort = p
		}
	}

	cfg := PluginConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Pass:     os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		ReadOnly: readOnly,
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", webPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	log.Printf("MySQL Web Client listening on http://%s", addr)

	hub := newCollabHub()
	go hub.run()

	go func() {
		if err := http.Serve(listener, setupHTTP(db, cfg, hub)); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	db.Close()
}
