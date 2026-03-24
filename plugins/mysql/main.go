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
}

type StartResponse struct {
	Status  string `json:"status"`
	WebPort int    `json:"web_port,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
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

func handleStart(cfg PluginConfig) {
	// Build DSN
	user := cfg.User
	if user == "" {
		user = "root"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		user, cfg.Pass, cfg.Host, cfg.Port, cfg.Database)

	// Test connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		respondError(fmt.Sprintf("failed to open connection: %v", err))
		return
	}
	if err := db.Ping(); err != nil {
		respondError(fmt.Sprintf("cannot connect to MySQL at %s:%d: %v", cfg.Host, cfg.Port, err))
		return
	}

	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		respondError(fmt.Sprintf("failed to find available port: %v", err))
		return
	}
	webPort := listener.Addr().(*net.TCPAddr).Port

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/databases", databasesHandler(db))
	mux.HandleFunc("/api/tables", tablesHandler(db))
	mux.HandleFunc("/api/schema", schemaHandler(db))
	mux.HandleFunc("/api/query", queryHandler(db))
	mux.HandleFunc("/api/info", infoHandler(db, cfg))

	// Static files
	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/", http.FileServer(http.FS(staticSub)))

	// Respond with ready
	respond(StartResponse{
		Status:  "ready",
		WebPort: webPort,
		Name:    "MySQL Web Client",
		Version: "0.1.0",
	})

	// Start HTTP server in background
	go func() {
		if err := http.Serve(listener, mux); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for stop signal or stdin close
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	db.Close()
}
