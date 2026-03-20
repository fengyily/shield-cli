---
title: FAQ — Shield CLI Frequently Asked Questions
description: Answers to common questions about Shield CLI including pricing, platform support, concurrent connections, security, network requirements, and configuration storage.
head:
  - - meta
    - name: keywords
      content: Shield CLI FAQ, frequently asked questions, help, troubleshooting, support
  - - script
    - type: application/ld+json
    - |
      {
        "@context": "https://schema.org",
        "@type": "FAQPage",
        "mainEntity": [
          {
            "@type": "Question",
            "name": "Is Shield CLI free?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "The Shield CLI client is open source and free. The public tunnel service has a free tier — check the website for details."
            }
          },
          {
            "@type": "Question",
            "name": "Which operating systems does Shield CLI support?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "macOS, Linux, and Windows, with amd64, arm64, 386, and armv7 architectures."
            }
          },
          {
            "@type": "Question",
            "name": "Are there any dependencies required to run Shield CLI?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "No. Shield CLI is a single binary — download and run. No dependencies needed."
            }
          },
          {
            "@type": "Question",
            "name": "How many apps can I connect simultaneously with Shield CLI?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "In Web UI mode, up to 3 concurrent connections. CLI mode connects to 1 service at a time."
            }
          },
          {
            "@type": "Question",
            "name": "Does Shield CLI auto-reconnect on disconnect?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Yes. Shield CLI has built-in auto-reconnect with exponential backoff, up to a maximum interval of 10 seconds."
            }
          },
          {
            "@type": "Question",
            "name": "Is the data transfer secure in Shield CLI?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "All data is transmitted through WebSocket encrypted tunnels. Local credentials use AES-256-GCM encryption with machine fingerprint-derived keys."
            }
          },
          {
            "@type": "Question",
            "name": "Are passwords stored on the server?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "No. Target service passwords are used only during connection establishment and are not persisted on the server."
            }
          },
          {
            "@type": "Question",
            "name": "What firewall ports does Shield CLI need?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Shield CLI only requires outbound connections to the public gateway on port 62888 (WebSocket) and port 443 (HTTPS API). No inbound ports need to be opened."
            }
          },
          {
            "@type": "Question",
            "name": "How many app profiles can I save in Shield CLI?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Up to 10 app profiles can be saved, each encrypted locally with AES-256-GCM."
            }
          }
        ]
      }
---

# FAQ

## General

### Is Shield CLI free?

The Shield CLI client is open source and free. The public tunnel service has a free tier — check the [website](https://console.yishield.com) for details.

### Which operating systems are supported?

macOS, Linux, and Windows, with amd64, arm64, 386, and armv7 architectures.

### Are there any dependencies?

No. Shield CLI is a single binary — download and run.

## Connections

### How many apps can I connect simultaneously?

In Web UI mode, up to **3** concurrent connections. CLI mode connects to 1 service at a time.

### Does it auto-reconnect on disconnect?

Yes. Shield CLI has built-in auto-reconnect with exponential backoff, up to a maximum interval of 10 seconds.

### Do Access URLs expire?

Access URLs are valid while the tunnel is connected. They become invalid once disconnected.

### Can multiple people use the same URL?

Yes. The same Access URL can be accessed from multiple browsers simultaneously.

## Security

### Is the data transfer secure?

All data is transmitted through WebSocket encrypted tunnels. Local credentials use AES-256-GCM encryption.

### Are passwords stored on the server?

No. Target service passwords are used only during connection establishment and are not persisted on the server.

### What is a machine fingerprint?

A unique identifier generated from your machine's hardware information. It's used to derive encryption keys and identify the connector. It contains no personal information.

## Network

### GitHub is slow in China — what can I do?

Use the jsDelivr CDN mirror:

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

### Does it work behind a corporate proxy?

Shield CLI uses WebSocket for tunneling. WebSocket can typically traverse HTTP proxies via the CONNECT method. If it doesn't work, check with your network admin whether WebSocket connections are allowed.

### What firewall ports are needed?

Shield CLI only requires **outbound** connections to the public gateway on port 62888 (WebSocket). No inbound ports need to be opened.

## Configuration

### Where are app configurations stored?

Encrypted locally:
- macOS / Linux: `~/.shield-cli/`
- Windows: `%LOCALAPPDATA%\ShieldCLI\`

### Can I sync configurations across machines?

Not currently. Each machine's configuration and credentials are independent (encrypted with each machine's unique fingerprint).

### How many apps can I save?

10.
