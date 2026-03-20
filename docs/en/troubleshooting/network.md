---
title: Network Issues — Firewall, Proxy, and Connectivity Guide
description: Diagnose and resolve Shield CLI network issues. Covers connectivity checks, proxy environments, firewall configuration, auto-reconnect, and China optimization.
head:
  - - meta
    - name: keywords
      content: Shield CLI network, firewall, proxy, WebSocket, connectivity, China mirror, auto-reconnect
---

# Network Issues

## Connectivity Checklist

### Step 1: Check Target Service

Confirm the target service is accessible from the machine running Shield CLI:

```bash
# Check port connectivity
telnet <target-ip> <port>

# Or use nc
nc -zv <target-ip> <port>
```

### Step 2: Check Public Gateway

Confirm you can reach the Shield public gateway:

```bash
curl -I https://console.yishield.com
```

### Step 3: Check WebSocket

Shield CLI uses WebSocket for tunneling. If you're behind a proxy or firewall, ensure WebSocket connections are not being intercepted.

## Proxy Environments

If your network requires an HTTP proxy to access the internet, Shield CLI should work in most cases because WebSocket can traverse proxies via the HTTP CONNECT method.

If connections still fail, check with your network admin to confirm WebSocket connections are allowed through the proxy.

## Firewall Configuration

Shield CLI **only requires outbound connections** — no inbound ports need to be opened:

| Direction | Destination | Port | Protocol |
|---|---|---|---|
| Outbound | console.yishield.com | 62888 | WebSocket (WSS) |
| Outbound | console.yishield.com | 443 | HTTPS (API) |

## Auto-Reconnect

Shield CLI has built-in automatic reconnection:

- Retries immediately on disconnect detection
- Exponential backoff: 1s → 2s → 4s → 8s → 10s (maximum)
- Tunnel re-established automatically when network recovers
- No manual intervention needed

## China Optimization

If you're in mainland China, you may experience slow GitHub access:

**Use mirror for installation:**
```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

**Use nearby access node:**
```bash
shield ssh 10.0.0.5 --visable=HK
```
