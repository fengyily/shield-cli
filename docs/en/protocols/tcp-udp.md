---
title: TCP/UDP Port Proxy — Forward Any TCP/UDP Service
description: Create TCP/UDP port proxy tunnels with Shield CLI to forward MySQL, Redis, PostgreSQL, DNS, and any other service without opening a browser.
head:
  - - meta
    - name: keywords
      content: TCP proxy, UDP proxy, port forwarding, MySQL tunnel, Redis tunnel, DNS tunnel, Shield CLI TCP UDP
---

# TCP / UDP Port Proxy

Forward arbitrary TCP/UDP ports through Shield CLI, ideal for databases, caches, DNS, and other services that don't need browser rendering.

## How It Differs from Other Protocols

| | SSH/RDP/VNC/HTTP | TCP/UDP |
|---|---|---|
| Browser rendering | Yes (HTML5) | No |
| Auto-open browser | Yes | No |
| Connection method | Browser URL | Client tool via domain:port |
| Default port | Yes | None (must specify) |

## Quick Start

```bash
# TCP port proxy
shield tcp 3306                      # MySQL (127.0.0.1:3306)
shield tcp 6379                      # Redis (127.0.0.1:6379)
shield tcp 5432                      # PostgreSQL (127.0.0.1:5432)
shield tcp 192.168.1.10:3306         # Remote MySQL

# UDP port proxy
shield udp 53                        # DNS (127.0.0.1:53)
shield udp 514                       # Syslog (127.0.0.1:514)
shield udp 192.168.1.1:161           # SNMP
```

::: warning Port is required
TCP/UDP have no default port. Omitting the port will produce an error:
```
Error: port is required for tcp protocol

Usage: shield tcp <port> or shield tcp <ip:port>
```
:::

## Connection Guide

After the tunnel is established, Shield CLI prints connection info (no browser opens):

```
  📡 Connection Guide (TCP port proxy):
    your-app.cn01.apps.yishield.com:58379  →  127.0.0.1:3306

    Examples:
      telnet your-app.cn01.apps.yishield.com 58379
      mysql -h your-app.cn01.apps.yishield.com -P 58379 -u root
      redis-cli -h your-app.cn01.apps.yishield.com -p 58379
```

Use the printed **dedicated domain and port** with your preferred client tool.

## Use Cases

- **Database access** — remote MySQL, PostgreSQL, MongoDB, Redis
- **DNS forwarding** — proxy internal DNS servers
- **Message queues** — connect to RabbitMQ, Kafka, NATS
- **Custom services** — any service listening on a TCP/UDP port

## Address Format

| Input | Resolves To |
|---|---|
| `shield tcp 3306` | `127.0.0.1:3306` |
| `shield tcp 192.168.1.10:3306` | `192.168.1.10:3306` |
| `shield udp 53` | `127.0.0.1:53` |
| `shield udp 192.168.1.1:161` | `192.168.1.1:161` |
