---
title: Commands Reference — Full Shield CLI Parameter Guide
description: Complete reference for all Shield CLI commands, flags, address formats, and examples. Covers shield start, ssh, rdp, vnc, http, https, telnet, and clean commands.
head:
  - - meta
    - name: keywords
      content: Shield CLI commands, CLI reference, shield ssh, shield rdp, shield start, command flags, parameters
---

# Commands Reference

## Command Overview

| Command | Description |
|---|---|
| `shield start [port]` | Launch Web management dashboard (default port 8181) |
| `shield ssh [address]` | Create an SSH tunnel |
| `shield rdp [address]` | Create an RDP tunnel |
| `shield vnc [address]` | Create a VNC tunnel |
| `shield http [address]` | Create an HTTP tunnel |
| `shield https [address]` | Create an HTTPS tunnel |
| `shield telnet [address]` | Create a Telnet tunnel |
| `shield tcp <port\|address>` | Create a TCP port proxy (port required) |
| `shield udp <port\|address>` | Create a UDP port proxy (port required) |
| `shield install [--port]` | Install as system service (auto-start on boot) |
| `shield uninstall` | Uninstall system service |
| `shield clean` | Clear local credential cache |

## Address Format

`[address]` supports the following formats:

| Format | Example | Description |
|---|---|---|
| Omitted | `shield ssh` | Uses `127.0.0.1` + protocol default port |
| Port only | `shield ssh 2222` | Uses `127.0.0.1` + specified port |
| IP only | `shield ssh 10.0.0.5` | Uses specified IP + protocol default port |
| Full address | `shield ssh 10.0.0.5:2222` | Uses specified IP + specified port |

### Default Ports

| Protocol | Default Port |
|---|---|
| SSH | 22 |
| RDP | 3389 |
| VNC | 5900 |
| HTTP | 80 |
| HTTPS | 443 |
| Telnet | 23 |
| TCP | None (must specify) |
| UDP | None (must specify) |

## Global Flags

| Flag | Description | Example |
|---|---|---|
| `--username` | Target service username | `--username root` |
| `--auth-pass` | Target service password | `--auth-pass mypass` |
| `--server` | Custom server address | `--server https://my.server/raas` |

## SSH-Specific Flags

| Flag | Description | Example |
|---|---|---|
| `--private-key` | SSH private key file path | `--private-key ~/.ssh/id_rsa` |
| `--passphrase` | Private key passphrase | `--passphrase mypass` |
| `--enable-sftp` | Enable SFTP file transfer | `--enable-sftp` |

## Access Mode Flags

| Flag | Description | Example |
|---|---|---|
| `--visable` | Visible mode (default) | `--visable` |
| `--visable=<node>` | Visible mode with specific node | `--visable=HK` |
| `--invisible` | Invisible mode, requires auth code | `--invisible` |

## Service Management

| Command | Description |
|---|---|
| `shield install` | Install as system service with default port 8181 |
| `shield install --port 8182` | Install with custom port |
| `shield uninstall` | Remove system service |

### Install Flags

| Flag | Default | Description |
|---|---|---|
| `--port` | `8181` | Web UI port number |

The install command automatically detects port conflicts and suggests available alternatives. See [System Service Installation](/en/guide/system-service) for platform-specific details.

## Examples

```bash
# Simplest usage
shield ssh

# Full parameters
shield ssh 10.0.0.5:2222 --username root --auth-pass mypass --enable-sftp

# Web UI mode
shield start
shield start 9090

# TCP/UDP port proxy
shield tcp 3306                          # MySQL
shield tcp 192.168.1.10:6379             # Redis
shield udp 53                            # DNS

# Invisible mode RDP
shield rdp 10.0.0.5 --username Administrator --invisible

# Clear cache
shield clean

# Install as system service
shield install
shield install --port 8182

# Uninstall service
shield uninstall
```
