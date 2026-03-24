---
title: Plugin Development Guide — Build Custom Shield CLI Plugins
description: Complete guide to developing Shield CLI plugins, including protocol specification, project templates, Web UI development, and publishing workflow.
head:
  - - meta
    - name: keywords
      content: Shield CLI plugin development, custom plugin, Go plugin, plugin development, Shield extension
---

# Plugin Development Guide

This guide covers how to develop custom plugins for Shield CLI to support new protocols or service types.

## Concept

A Shield plugin is an **independent executable** with a simple responsibility:

1. Receive startup configuration from stdin (JSON)
2. Start a local Web service
3. Return ready status and Web port to stdout
4. Keep running until a stop signal is received

Shield handles everything else:
- Starting and stopping the plugin process
- Exposing the plugin's Web port via encrypted tunnel
- Providing a unified CLI experience

## Communication Protocol

Shield and plugins communicate via **stdin/stdout single-line JSON**. There are only three message types.

### 1. Start Request (stdin)

After launching the plugin, Shield sends a single JSON line via stdin:

```json
{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"xxx","database":"mydb"}}
```

Fields:

| Field | Type | Description |
|---|---|---|
| `action` | string | Always `"start"` |
| `config.host` | string | Target service IP |
| `config.port` | int | Target service port |
| `config.user` | string | Username (may be empty) |
| `config.pass` | string | Password (may be empty) |
| `config.database` | string | Database name (may be empty) |

### 2. Ready Response (stdout)

When ready, the plugin writes a single JSON line to stdout:

**Success:**
```json
{"status":"ready","web_port":19876,"name":"My Plugin","version":"0.1.0"}
```

**Failure:**
```json
{"status":"error","message":"cannot connect to service: connection refused"}
```

| Field | Type | Description |
|---|---|---|
| `status` | string | `"ready"` or `"error"` |
| `web_port` | int | Plugin Web service port (required when status=ready) |
| `name` | string | Plugin display name |
| `version` | string | Plugin version |
| `message` | string | Error message (when status=error) |

### 3. Stop Request (stdin)

When Shield exits, it sends via stdin:

```json
{"action":"stop"}
```

The plugin should exit gracefully. If it doesn't exit within 5 seconds, Shield will force-kill the process.

### Timeout

Shield waits up to **15 seconds** for the ready response. If the plugin doesn't respond, it will be terminated.

## Project Template

Here's a complete plugin project structure:

```
shield-plugin-example/
├── main.go              # Entry point, reads stdin, starts web server
├── handler.go           # HTTP handlers (business logic)
├── static/
│   └── index.html       # Web UI (embedded via embed.FS)
├── go.mod
├── go.sum
├── Makefile
└── .goreleaser.yml      # Multi-platform release config
```

### main.go Template

```go
package main

import (
    "embed"
    "encoding/json"
    "io/fs"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

//go:embed static/*
var staticFS embed.FS

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
            return // stdin closed
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

func handleStart(cfg PluginConfig) {
    // 1. Connect to target service and verify availability
    // conn, err := connectToService(cfg)
    // if err != nil {
    //     respond(StartResponse{Status: "error", Message: err.Error()})
    //     return
    // }

    // 2. Find an available port
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        respond(StartResponse{Status: "error", Message: err.Error()})
        return
    }
    webPort := listener.Addr().(*net.TCPAddr).Port

    // 3. Setup HTTP routes
    mux := http.NewServeMux()
    // mux.HandleFunc("/api/...", yourHandler)
    staticSub, _ := fs.Sub(staticFS, "static")
    mux.Handle("/", http.FileServer(http.FS(staticSub)))

    // 4. Return ready response
    respond(StartResponse{
        Status:  "ready",
        WebPort: webPort,
        Name:    "Example Plugin",
        Version: "0.1.0",
    })

    // 5. Start HTTP server
    go http.Serve(listener, mux)

    // 6. Wait for stop signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
}
```

### Key Points

1. **stdout is only for protocol messages** — do not print logs to stdout, use stderr
2. **Auto-assign ports** — use `net.Listen("tcp", "127.0.0.1:0")` for system-assigned ports
3. **embed.FS for static files** — Web UI is packaged into the binary, no external files needed
4. **Connect before responding** — verify the target service is reachable before returning `ready`

## Registering a Plugin

After development, register the plugin in Shield's `plugin/install.go`:

```go
var KnownPlugins = map[string]PluginInfo{
    "example": {
        Name:        "example",
        Source:      "your-org/shield-plugin-example",
        Protocols:   []string{"example"},
        DefaultPort: 9999,
    },
}
```

Or use local installation for testing:

```bash
# Build the plugin
go build -o shield-plugin-example .

# Install locally
shield plugin add example --from ./shield-plugin-example

# Test
shield example 127.0.0.1:9999
```

## Publishing a Plugin

### GoReleaser Configuration

```yaml
# .goreleaser.yml
project_name: shield-plugin-example

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    binary: shield-plugin-example

archives:
  - format: tar.gz
    name_template: "shield-plugin-example_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
```

### Release Workflow

```bash
# Tag
git tag v0.1.0
git push origin v0.1.0

# GoReleaser builds and publishes to GitHub Releases
goreleaser release
```

### Naming Convention

| Item | Convention |
|---|---|
| Repository | `shield-plugin-<name>` |
| Binary | `shield-plugin-<name>` |
| Release asset | `shield-plugin-<name>_<os>_<arch>.tar.gz` |

Shield downloads the platform-specific binary from GitHub Releases using this naming convention.

## Testing a Plugin

### Manual Protocol Testing

Test the stdin/stdout protocol directly without Shield:

```bash
echo '{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"root"}}' | ./shield-plugin-mysql
# Expected: {"status":"ready","web_port":xxxxx,"name":"MySQL Web Client","version":"0.1.0"}
```

### End-to-End Testing

```bash
# 1. Install locally
shield plugin add mysql --from ./shield-plugin-mysql

# 2. Start (connect to local MySQL)
shield mysql 127.0.0.1:3306 --db-user root --db-pass root

# 3. Verify the Web interface is accessible
# 4. Verify it's also accessible via the public URL
```

## Web UI Development Tips

### Recommended Stack

For Shield plugin Web interfaces, a pure frontend approach (no framework dependencies) works best:

- **HTML + Vanilla JS** — lightest weight, suitable for simple interfaces
- **Single HTML file** — all CSS/JS inline, embedded via `embed.FS`
- Avoid npm build toolchains — adds complexity with limited benefit for plugins

### Design Principles

1. **Lightweight** — single HTML file, under 50KB gzipped
2. **Responsive** — support narrow viewports (may open in embedded browsers)
3. **Safe defaults** — read-only mode by default, write operations require manual unlock
4. **Clear feedback** — loading states, success/failure notifications for all operations

## Reference Implementation

- **shield-plugin-mysql** — in-repo implementation at `plugins/mysql/`, serves as a reference for developing new plugins

## Next Steps

- [Plugin System Overview](/en/plugins/)
- [MySQL Plugin Documentation](/en/plugins/mysql)
- [Command Reference](/en/reference/commands)
