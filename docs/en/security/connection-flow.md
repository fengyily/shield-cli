---
title: Connection Flow — How Shield CLI Tunnels Work
description: Understand Shield CLI's complete connection flow from command execution to browser access. Covers WebSocket tunnel architecture, auto-reconnect, and port allocation.
head:
  - - meta
    - name: keywords
      content: Shield CLI architecture, connection flow, WebSocket tunnel, Chisel, auto-reconnect, tunnel architecture
---

# Connection Flow

Understand the complete process from running a command to browser access.

## Architecture Overview

```
Internal Service ←→ Shield CLI ←→ Public Gateway ←→ Browser
  (SSH/RDP/...)      (Encrypted)    (HTML5 Render)   (Any Device)
```

## Connection Steps

### CLI Mode

1. **Resolve address** — Determine target IP and port from protocol and input
2. **Load credentials** — Read or generate encrypted credentials (machine fingerprint-based)
3. **Authenticate** — Interactive prompt or command-line flags for target service credentials
4. **Call server API** — Send quick-setup request to the public gateway
5. **Establish main tunnel** — Create encrypted WebSocket tunnel
6. **Create resource tunnel** — Map ports to the target service on the main tunnel
7. **Assign URL** — Receive a unique Access URL
8. **Auto-open browser** — Open the Access URL in the default browser
9. **Start local API** — Launch health check and management API on port 4000-5000

### Web UI Mode

1. **Start** — `shield start` launches the web server and dashboard
2. **Add app** — User configures app in the UI (stored locally, encrypted)
3. **Click Connect** — Background async tunnel creation:
   - Call quick-setup API (with retry logic)
   - Establish main tunnel (shared across apps)
   - Create resource tunnel (per-app)
   - Poll target site until reachable
4. **Status updates** — Frontend polls status every 2 seconds (Connecting → Connected)
5. **Auto-open** — Browser opens the Access URL in a new tab on success
6. **Disconnect** — Closes resource tunnel, preserves main tunnel for reuse

## Auto-Reconnect

Shield CLI has built-in automatic reconnection:

- Immediately retries on disconnect detection
- Exponential backoff: 1s → 2s → 4s → 8s → 10s (maximum)
- Transparent to the user — no manual intervention needed

## Port Allocation

| Port | Purpose |
|---|---|
| `8181` (customizable) | Web UI dashboard |
| `4000-5000` | Local API server (health checks, tunnel management) |
| `62888` | Public gateway communication |
