---
title: Changelog — Shield CLI Release History
description: Complete release history and changelog for Shield CLI. Track new features, improvements, and bug fixes across all versions.
head:
  - - meta
    - name: keywords
      content: Shield CLI changelog, release history, version history, updates, new features
---

# Changelog

All notable changes to Shield CLI are documented here.

## v0.2.1 <Badge type="tip" text="latest" />

**Released: 2026-03-20**

### New Features

- **System service installation** — `shield install` registers Shield as a system service that starts automatically on boot
  - macOS: launchd user agent (no sudo required)
  - Linux: systemd service
  - Windows: Windows Service
- **Custom port support** — `shield install --port 8182` with automatic port conflict detection and alternative suggestion
- **System tray icon** (macOS & Windows) — click to open Dashboard, with Restart and Quit options
- **Async tunnel startup** — Web UI starts immediately, main tunnel connects in the background
- **Tunnel status API** — `GET /api/tunnel` endpoint for frontend to poll tunnel readiness

### Improvements

- Split goreleaser into desktop (CGO + tray) and Linux (pure Go) builds
- App connections blocked with clear message while tunnel is still connecting

---

## v0.2.0

**Released: 2026-03-19**

### New Features

- **Web UI management platform** — browser-based dashboard at `localhost:8181`
  - Add, edit, delete up to 10 application profiles
  - One-click connect/disconnect with real-time status
  - Encrypted local storage for app configurations
- **Persistent configuration** — save application profiles with AES-256-GCM encrypted storage
- **Multi-connection support** — up to 3 concurrent active tunnel connections
- **Connection manager** — shared main tunnel with per-app dynamic resource tunnels

### Improvements

- Redesigned logo and branding
- Updated README with Web UI screenshots and examples

---

## v0.1.3

**Released: 2026-03-18**

### New Features

- **Windows installer** — PowerShell one-liner installation script
- **Linux installer** — curl-based install script
- **Bilingual README** — split into English (`README.md`) and Chinese (`README_CN.md`)

### Improvements

- Default to visible access mode

---

## v0.1.2

**Released: 2026-03-18**

### New Features

- **Scoop package** — `scoop install shield-cli` for Windows
- **deb / rpm packages** — native Linux package formats via nfpm
- **curl installer** — `curl -fsSL ... | sh` one-liner for Linux/macOS
- **China CDN mirror** — jsDelivr-based install script for users in China

---

## v0.1.1

**Released: 2026-03-18**

### Improvements

- **Positional arguments** — `shield ssh 10.0.0.5:2222` instead of `--type ssh --source 10.0.0.5:2222`
- **Smart defaults** — omit IP for localhost, omit port for protocol default
- Simplified CLI usage with intuitive address resolution

---

## v0.1.0

**Released: 2026-03-18**

### New Features

- **GoReleaser integration** — automated cross-platform builds (macOS, Linux, Windows × amd64, arm64)
- **Homebrew tap** — `brew install shield-cli`
- **Automated releases** — GitHub Actions CI/CD pipeline

---

## v0.0.1

**Released: 2026-03-16**

### Initial Release

- Core tunnel connectivity via Chisel protocol
- Supported protocols: SSH, RDP, VNC, HTTP, HTTPS, Telnet
- AES-256-GCM encrypted credential storage with machine fingerprint binding
- Visible and invisible access modes
- Auto-open browser on connection
- Password masking in all log output
- CI/CD pipeline with GitHub Actions
