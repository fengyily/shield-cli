---
title: System Service Installation — Run Shield CLI on Boot
description: Install Shield CLI as a system service that starts automatically with your operating system. Supports macOS (launchd), Linux (systemd), and Windows services.
head:
  - - meta
    - name: keywords
      content: Shield CLI service, install service, system service, launchd, systemd, Windows service, auto start, daemon
---

# System Service Installation

Shield CLI can be installed as a system service that starts automatically when your computer boots. This is ideal for always-on access to your internal services.

## Quick Start

```bash
# Install with default port (8181)
shield install

# Install with custom port
shield install --port 8182
```

After installation, the Web UI will be available at `http://localhost:8181` (or your chosen port) and will survive reboots.

On **macOS** and **Windows**, a system tray icon appears automatically — click it to open the Dashboard in your browser.

## Commands

### `shield install`

Installs Shield CLI as a system service.

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8181` | Web UI port number |

**What it does:**
1. Checks if the service is already installed
2. Verifies the specified port is available
3. Registers the service with the operating system
4. Starts the service immediately

### `shield uninstall`

Removes the Shield CLI system service.

```bash
shield uninstall
```

**What it does:**
1. Stops the running service
2. Removes it from automatic startup
3. Your configuration and credentials are **preserved**

## Port Configuration

By default, Shield uses port **8181** for the Web UI. If this port is occupied:

```bash
# Shield will detect the conflict and suggest an alternative
$ shield install
Error: port 8181 is already in use.
Try an available port: shield install --port 8182
```

You can specify any available port:

```bash
shield install --port 9090
```

## Platform Details

### macOS (launchd)

- **Type:** User-level Launch Agent (no sudo required)
- **Plist:** `~/Library/LaunchAgents/com.yishield.shield-cli.plist`
- **Logs:** `~/.shield-cli/logs/shield-cli.log`
- Starts automatically on user login
- Restarts on failure (KeepAlive)
- System tray icon for quick Dashboard access

### Linux (systemd)

- **Type:** System service (requires sudo)
- **Unit:** `/etc/systemd/system/shield-cli.service`
- **Logs:** `journalctl -u shield-cli`
- Starts after network is online
- Restarts on failure with 5s delay

```bash
# Manual service management
sudo systemctl status shield-cli
sudo systemctl restart shield-cli
sudo journalctl -u shield-cli -f
```

### Windows

- **Type:** Windows Service (requires Administrator)
- **Service Name:** `ShieldCLI`
- Starts automatically on boot
- System tray icon for quick Dashboard access
- Managed via Services console (`services.msc`) or `sc` command

```powershell
# Manual service management
sc query ShieldCLI
sc stop ShieldCLI
sc start ShieldCLI
```

## Reinstalling with a Different Port

To change the port of an installed service:

```bash
shield uninstall
shield install --port 9090
```

## Troubleshooting

### Port already in use

```bash
# Check what's using the port
# macOS / Linux
lsof -i :8181

# Windows
netstat -ano | findstr :8181
```

### Service won't start (Linux)

```bash
# Check service logs
sudo journalctl -u shield-cli --no-pager -n 50

# Verify the binary path
which shield
```

### Service won't start (macOS)

```bash
# Check system logs
log show --predicate 'eventMessage contains "shield"' --last 5m

# Verify plist
plutil ~/Library/LaunchAgents/com.yishield.shield-cli.plist
```
