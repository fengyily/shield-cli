---
title: Shield CLI vs ngrok vs frp vs Cloudflare Tunnel — Technical Comparison
description: In-depth comparison of Shield CLI with ngrok, frp, and Cloudflare Tunnel covering protocol support, configuration complexity, security models, architecture, and pricing. Shield CLI is the only tunnel tool that renders RDP/VNC/SSH directly in the browser.
date: 2025-01-15
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI vs ngrok, Shield CLI vs frp, Shield CLI vs Cloudflare Tunnel, tunnel tool comparison, remote access comparison, RDP browser
  - - script
    - type: application/ld+json
    - |
      {
        "@context": "https://schema.org",
        "@type": "Article",
        "headline": "Shield CLI vs ngrok vs frp vs Cloudflare Tunnel — Technical Comparison",
        "description": "In-depth comparison of Shield CLI with ngrok, frp, and Cloudflare Tunnel covering protocol support, configuration, security, architecture, and pricing",
        "author": {"@type": "Organization", "name": "Shield CLI"},
        "publisher": {"@type": "Organization", "name": "Shield CLI", "url": "https://docs.yishield.com"},
        "datePublished": "2025-01-15"
      }
---

# I Opened an Intranet RDP Desktop in a Browser with One Command — A Technical Comparison of Shield CLI vs Popular Tunnel Tools

> I've been exploring remote operations solutions recently, needing to expose intranet Windows Remote Desktop and Linux SSH to external collaborators. After trying ngrok, frp, and Cloudflare Tunnel, I discovered Shield CLI — a tool with a fundamentally different approach. I spent some time doing an in-depth comparison and documented the technical details here.

## First, the Scenario: Why "Tunneling" Doesn't Equal "Remote Access"

The typical use case for tunnel tools is: you have a local web service running and want to give someone external temporary access. ngrok is practically the standard for this scenario.

But what if your scenario looks like this:

- You need a client to operate an intranet Windows Remote Desktop (RDP) directly through a browser
- You want to give an outsourced team temporary SSH terminal access without requiring them to install any client software
- A VNC desktop in a demo environment needs a shareable link

This is where you'll find that, yes, ngrok can establish a TCP tunnel, but the other party still needs to install an RDP client or configure an SSH tool. **Tunneling solves "network reachability," but not "terminal usability."**

Shield CLI's approach: after the tunnel is established, it provides an HTML5 Web terminal directly at the gateway (based on protocol gateways like Apache Guacamole). Users get an HTTPS link — open it in a browser and you have an RDP desktop or SSH terminal.

This is the fundamental difference in approach between these products. Let's dive into the detailed comparison.

---

## 1. Protocol Support: Who Actually Delivers "Remote Desktop"

| Protocol | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|----------|-----------|-------|-----|-------------------|
| HTTP/HTTPS | ✅ | ✅ | ✅ | ✅ |
| Generic TCP | ✅ (via specific protocols) | ✅ | ✅ | ✅ (Spectrum) |
| UDP | ❌ | ❌ (paid) | ✅ | ❌ |
| SSH (browser terminal) | ✅ Built-in Web Terminal | ❌ TCP forwarding only | ❌ TCP forwarding only | ✅ (requires Access config) |
| RDP (browser desktop) | ✅ Built-in Web Desktop | ❌ | ❌ | ❌ |
| VNC (browser desktop) | ✅ Built-in Web Desktop | ❌ | ❌ | ❌ |
| Telnet | ✅ | ❌ | ❌ | ❌ |
| SFTP file transfer | ✅ (in SSH mode with `--enable-sftp`) | ❌ | ❌ | ❌ |

**Key difference**: ngrok and frp perform **L4 port forwarding** — they map the remote port 3389 to the public internet, but users still need to launch mstsc.exe (Windows Remote Desktop Client) to connect. Shield CLI performs **L7 protocol rendering** — remote services are rendered directly in the browser via HTML5, with zero client installation.

Let's look at actual command comparisons. Exposing an intranet Windows machine's RDP:

```bash
# ngrok: establishes the tunnel, but the other party needs an RDP client
ngrok tcp 3389
# Output: tcp://0.tcp.ngrok.io:12345
# Other party needs to: open Remote Desktop Connection → enter 0.tcp.ngrok.io:12345 → log in

# frp: requires deploying an frps server + writing config files
# frpc.toml:
# [[proxies]]
# name = "rdp"
# type = "tcp"
# localIP = "127.0.0.1"
# localPort = 3389
# remotePort = 7001
frpc -c frpc.toml
# Other party needs to: same as ngrok, must have an RDP client

# Shield CLI: one command, open directly in browser
shield rdp --username admin --auth-pass mypass
# Output: https://xxxx-yishield.ac.example.com
# Other party needs to: click the link, done
```

Same for SSH:

```bash
# Shield CLI
shield ssh 10.0.0.5 --username root
# Web terminal appears directly in browser, supports SFTP file upload/download

# ngrok
ngrok tcp 22
# Other party needs to: ssh -p 12345 root@0.tcp.ngrok.io
# Plus dealing with known_hosts, keys, etc.
```

---

## 2. Configuration Complexity: From "One Command" to "A Pile of Config Files"

### Shield CLI's Smart Defaults

Shield CLI's CLI parameter design includes extensive default inference to minimize input:

```bash
shield ssh                  # equivalent to 127.0.0.1:22
shield ssh 2222             # equivalent to 127.0.0.1:2222 (pure number → port)
shield ssh 10.0.0.5         # equivalent to 10.0.0.5:22 (IP → use default port)
shield ssh 10.0.0.5:2222    # fully specified
shield rdp                  # equivalent to 127.0.0.1:3389
shield vnc 10.0.0.10:5901   # fully specified
shield http 3000            # equivalent to 127.0.0.1:3000
```

This logic is implemented in `cmd/helpers.go`: pure number → port (uses 127.0.0.1), contains `.` or `:` → IP or IP:Port, empty → default IP + default port. Each protocol has its own default port (SSH=22, RDP=3389, VNC=5900, HTTP=80, HTTPS=443, Telnet=23).

### Configuration Comparison with frp

frp is the classic self-hosted tunnel solution, but has a higher configuration barrier. Here's a complete SSH forwarding config:

```toml
# frps.toml (server side — you need a public-facing machine)
bindPort = 7000

# frpc.toml (client side)
serverAddr = "your-server.com"
serverPort = 7000

[[proxies]]
name = "ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 6000
```

This involves: deploying a server, configuring port mappings, managing config files, and maintaining a public server. Shield CLI doesn't require you to manage a server — the public gateway is provided by Shield's infrastructure (similar to ngrok's model), while the CLI side is open source (Apache 2.0).

### Configuration Comparison with ngrok

ngrok's single command is indeed concise:

```bash
ngrok http 8080
```

But in multi-protocol scenarios (e.g., needing SSH + RDP + an HTTP service simultaneously), ngrok requires a config file:

```yaml
# ngrok.yml
tunnels:
  ssh:
    proto: tcp
    addr: 22
  rdp:
    proto: tcp
    addr: 3389
  web:
    proto: http
    addr: 8080
```

Shield CLI manages multiple services through a Web UI (`shield start`), supporting up to 10 saved application configurations. You can dynamically manage connections by clicking Connect/Disconnect in the interface — no config files needed.

### Configuration Comparison with Cloudflare Tunnel

Cloudflare Tunnel has the most complex configuration (but is also the most powerful):

```yaml
# config.yml
tunnel: your-tunnel-id
credentials-file: /root/.cloudflared/your-tunnel-id.json

ingress:
  - hostname: ssh.example.com
    service: ssh://localhost:22
  - hostname: rdp.example.com
    service: rdp://localhost:3389
  - service: http_status:404
```

You also need: Cloudflare account → add domain → create Tunnel → configure DNS → configure Access Policy. For enterprise-grade persistent deployments this is justified, but for "temporarily giving someone remote desktop access," it's using a sledgehammer to crack a nut.

---

## 3. Security Model Comparison

### Credential Storage

| Dimension | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|-----------|-----------|-------|-----|-------------------|
| Credential storage | AES-256-GCM encrypted local file | Token stored in plaintext in `~/.ngrok2/ngrok.yml` | Token in config file | JSON file stored in `~/.cloudflared/` |
| Key source | Machine fingerprint SHA256 (hostname + MAC + Machine ID) | User account token | User-defined | Issued by Cloudflare |
| Cross-machine migration | ❌ Encrypted file bound to machine, invalidated after migration | ✅ Token is portable | ✅ Config file is portable | ✅ Credential file is portable |

Shield CLI's approach is interesting: it uses the **machine fingerprint** as the AES-256-GCM encryption key. The fingerprint consists of three parts:

1. Hostname (`os.Hostname()`)
2. MAC address of the first physical network interface (skipping docker/br-/veth/virbr and other virtual interfaces)
3. Platform-level Machine ID (Linux: `/etc/machine-id`, macOS: `IOPlatformUUID`, Windows: Registry `MachineGuid`)

The three are concatenated and SHA256 hashed to derive the AES key. This means:

- **Reduced credential leak risk**: Even if someone copies the `~/.shield-cli/.credential` file, it can't be decrypted on another machine (different machine fingerprint)
- **No "login" required**: Credentials are automatically generated and registered with the server on first use, with identity bound to the machine
- **Trade-off**: Switching machines or reinstalling the OS requires `shield clean` to reset credentials

Compare this with ngrok's `ngrok config add-authtoken <token>` approach — the token is in plaintext and can be copied to another machine and used immediately. Convenient but higher risk.

### Access Control

| Method | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|--------|-----------|-------|-----|-------------------|
| Public links | ✅ Visible mode (default) | ✅ Public by default | ✅ Public by default | ❌ Requires Access policy |
| Authorized access | 🔜 Invisible mode (planned) | ✅ IP whitelist/OAuth (paid) | ✅ Self-implemented | ✅ Access (zero trust) |
| Link validity | 24-hour API Key auto-refresh | Free tier: 2 hours/8 hours | Unlimited (while server runs) | Unlimited (while Tunnel runs) |

Shield CLI currently defaults to Visible mode — the generated HTTPS link is accessible by anyone. However, the server sets a **24-hour validity period** for each API Key, which auto-refreshes upon expiration. Compared to ngrok's free tier 2-hour limit (adjusted after 2024), Shield's free quota is more generous.

Cloudflare Tunnel is the strongest on security — you can configure comprehensive zero-trust policies (email verification, SAML SSO, IP restrictions, etc.), but this also means heavier configuration overhead.

### Password Handling

Shield CLI applies password masking in logs — showing only the first and last 2 characters:

```
Connecting to 10.0.0.5:22 with password: my****ss
```

SSH private keys are passed via the `--private-key` parameter as a file path, so key contents are never exposed on the command line. After credentials are transmitted to the server, they are stored in `main_app_config` for protocol gateway authentication.

---

## 4. Architecture Comparison: Chisel vs ngrok's Proprietary Protocol

### Shield CLI's Dual-Layer Tunnel Architecture

Shield CLI uses [Chisel](https://github.com/jpillora/chisel) (a WebSocket-based TCP tunnel library) under the hood. It establishes two tunnels:

```
┌──────────────────┐                    ┌──────────────┐                    ┌──────────────────┐
│ Intranet Service │ ←── Local Net ──→  │ Shield CLI   │ ←── WebSocket ──→  │ Public Gateway   │
│ RDP/SSH/VNC      │                    │  (chisel     │     (wss://)       │ + Protocol       │
│ 10.0.0.5         │                    │   client)    │                    │   Rendering      │
└──────────────────┘                    └──────────────┘                    │   (Guacamole)    │
                                            │                               └──────────────────┘
                                       Tunnel 1: API Tunnel                        │
                                       (Control channel,               Tunnel 2: Resource Tunnel
                                        persistent)                   (Data channel, on-demand)
```

**API Tunnel** (Main Tunnel): Established on first connection, maps the local REST API port to the public network for dynamic management of subsequent resource tunnels. This tunnel is maintained persistently.

**Resource Tunnel**: Each application creates an independent chisel connection on demand, mapping the target service port to the public gateway. Up to 3 concurrent connections.

The benefit of this design: the API tunnel provides a "control plane," allowing the gateway to dynamically add or remove resource tunnels without requiring manual user intervention.

### ngrok's Architecture

ngrok uses a proprietary self-developed protocol. The client connects to ngrok's edge servers via TLS, with protocol details undisclosed. The upside is that performance can be optimized to the extreme; the downside is complete dependence on ngrok's infrastructure with no way to audit.

### frp's Architecture

frp uses a custom binary protocol (or optional KCP/QUIC). Both client and server are open source, so you can fully self-host. However, there's no protocol rendering layer — it only does port forwarding.

### Cloudflare Tunnel's Architecture

cloudflared connects to Cloudflare's global edge network via QUIC protocol, automatically using Anycast to select the nearest node. At the infrastructure level, this is the most robust solution (200+ data centers), but all your traffic passes through Cloudflare.

---

## 5. Local Management Experience

### Shield CLI's Web UI

```bash
shield start
# Browser automatically opens http://localhost:8181
```

This launches a local web management interface with features including:

- **Application management**: Add/edit/delete app configurations (protocol, target IP:Port, credentials, etc.)
- **One-click connect**: Click the Connect button, tunnel is established in the background, access link pops up on success
- **Status monitoring**: Real-time display of each application's connection status (idle / connecting / connected / failed)
- **Persistent configuration**: Up to 10 application configs, AES-256-GCM encrypted storage
- **Dark/Light theme**: Toggle support

The frontend is pure HTML5 + vanilla JS (~1500 lines), embedded in the binary with no external dependencies. The backend is a standard REST API.

### ngrok's Management

ngrok's free tier has no local UI. You can view request logs (HTTP tunnels only) via `http://localhost:4040`, but you can't manage multiple tunnels. Full management is on the ngrok Dashboard (SaaS), requiring account registration.

### frp's Management

frp has an optional Dashboard (enabled when launching `frps`) for viewing proxy lists and traffic statistics. However, it runs on the server side, not as a local client management interface. The UI is also quite basic.

### Cloudflare Tunnel's Management

Managed through the Cloudflare Zero Trust Dashboard (SaaS). The most comprehensive feature set (traffic analytics, access policies, audit logs), but cloud-dependent.

---

## 6. Deployment & Distribution

| Dimension | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|-----------|-----------|-------|-----|-------------------|
| Installation | Homebrew / Scoop / curl / dpkg / rpm / source build | Homebrew / apt / choco / snap / official download | GitHub Release download / source build | Homebrew / apt / official download |
| Binary size | ~15 MB | ~25 MB | ~12 MB (frpc) | ~35 MB |
| Platform support | Linux/macOS/Windows (amd64/arm64/386) | Linux/macOS/Windows/FreeBSD | Linux/macOS/Windows/FreeBSD + more | Linux/macOS/Windows |
| China mirror | ✅ jsDelivr CDN mirror | ❌ Requires VPN to download | ✅ GitHub directly accessible in China | ❌ Requires VPN |
| License | Apache 2.0 (CLI side) | Proprietary | Apache 2.0 (full) | Proprietary |
| Self-hosted server | 🔜 Planned | ❌ | ✅ Fully supported | ❌ |

For users in mainland China, Shield CLI provides a jsDelivr CDN mirror for installation:

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

Downloading ngrok and Cloudflare Tunnel in China is often blocked — a real pain point.

---

## 7. Pricing & Limitations

| Dimension | Shield CLI | ngrok (Free) | ngrok (Personal $8/mo) | frp | Cloudflare Tunnel |
|-----------|-----------|-------------|----------------------|-----|-------------------|
| Cost | Free | Free | $8/month | Free (server costs for self-hosting) | Free (domain must be on CF) |
| Tunnel count | 3 concurrent | 1 agent / 1 domain | 2 agents / 1 domain | Unlimited | Unlimited |
| Saved configs | 10 | N/A | N/A | Unlimited config files | Unlimited |
| Bandwidth limit | Not stated | 1 GB/month | 1 GB/month | Depends on server | Not stated |
| Connection duration | 24 hours (auto-renewal) | 2 hours (requires reconnect) | Unlimited | Unlimited | Unlimited |
| TCP tunnels | ✅ Free | ❌ Paid only | ✅ | ✅ | ✅ (Spectrum, paid) |
| Custom domains | 🔜 Planned | ❌ Paid only | ✅ | ✅ | ✅ |
| Access logs | Local logs | Dashboard | Dashboard | Dashboard | Dashboard + analytics |

An often-overlooked detail: **ngrok's free tier does not support TCP tunnels**. This means you cannot use ngrok for free to forward SSH (port 22) or RDP (port 3389). Shield CLI's TCP-based protocols (SSH/RDP/VNC/Telnet) are all free to use.

---

## 8. Real-World Usage Comparison

### Scenario 1: Demoing an Intranet System to a Client

```bash
# Shield CLI: one command, share the link
shield http 3000
# → https://abc123-yishield.ac.example.com
# Client clicks the link and sees your application directly

# ngrok: similar, but free tier links expire in 2 hours
ngrok http 3000
# → https://abc123.ngrok-free.app
# Client sees an ngrok warning page when visiting (free tier)
```

ngrok's free tier has an interstitial warning page ("You are about to visit..."), which looks unprofessional during client demos. Shield CLI doesn't have this limitation.

### Scenario 2: Remote Assistance for a Windows Desktop

```bash
# Shield CLI: Windows desktop appears directly in browser
shield rdp 10.0.0.100 --username admin --auth-pass P@ssw0rd
# → https://xxx-yishield.ac.example.com
# The other party operates a full Windows desktop in the browser

# Other tools: none can render RDP in a browser
# ngrok: ngrok tcp 3389 → other party needs an RDP client
# If the other party is on Mac/Linux, they also need to install Microsoft Remote Desktop or Remmina
```

This scenario is Shield CLI's core advantage. Other tunnel tools can only provide "network reachability" here, while users still need to solve "client compatibility" on their own.

### Scenario 3: Managing Multiple Intranet Services

```bash
# Shield CLI: launch Web UI for unified management
shield start
# Add multiple applications in the browser, click connect/disconnect

# frp: need to edit config files, restart client
# ngrok: need to write ngrok.yml, or open multiple terminal windows
```

---

## 9. Limitations & Trade-offs

In fairness, Shield CLI has its current-stage shortcomings:

1. **Server is not open source**: The gateway service is operated by Shield officially (console.yishield.com) and currently cannot be self-hosted. This means data passes through a third-party server. Self-hosted deployment is on the roadmap but hasn't been released yet. frp wins decisively on this point.

2. **Concurrency limits**: Maximum 3 concurrent connections and 10 saved configurations. Sufficient for individuals and small teams, but inadequate for enterprise scenarios.

3. **No UDP support**: The underlying Chisel is based on WebSocket (TCP) and doesn't support UDP protocols. frp is more comprehensive in this regard.

4. **Weak access control**: Currently only Visible mode — "anyone with the link can access." Invisible mode (requiring additional authorization keys) is planned but not yet available. Cloudflare Access's zero-trust approach is an order of magnitude ahead in security.

5. **Community ecosystem**: As a new project, the community is far smaller than ngrok (GitHub 25k+ stars) and frp (80k+ stars). You may need to read the source code directly when encountering issues.

6. **External service dependency**: While the CLI is open source, core functionality depends on Shield's public gateway. If the service is unavailable, the tool becomes unusable. This contrasts with frp's fully self-contained approach.

---

## 10. Selection Guide

| Your Need | Recommended Solution | Reason |
|-----------|---------------------|--------|
| Temporarily expose a local web service to colleagues | ngrok or Shield CLI | Both are one-command solutions |
| Remote desktop (RDP/VNC) via browser | **Shield CLI** | The only solution that renders desktop protocols in the browser |
| Fully self-hosted, no third parties | **frp** | Fully open source, self-deployed server |
| Enterprise-grade zero-trust remote access | **Cloudflare Tunnel + Access** | Most comprehensive security policy engine |
| Usage in mainland China network environment | **Shield CLI or frp** | Domestically reachable installation and service nodes |
| SSH + SFTP file transfer all-in-one | **Shield CLI** | Browser-based SSH + SFTP out of the box |
| UDP forwarding (gaming, DNS) | **frp** | The only solution supporting UDP |
| Budget-sensitive, need TCP tunnels | **Shield CLI** | ngrok TCP tunnels require payment |

---

## Conclusion

Shield CLI isn't trying to replace ngrok or frp — the problems they solve overlap but aren't identical. **If your core need is "letting others directly operate an intranet desktop or terminal through a browser," Shield CLI is currently the only tool that can do it with a single command.** It integrates tunnel tools and protocol gateways into a single workflow, eliminating the intermediate step of "installing a client."

But if you need fully self-controlled infrastructure (frp), enterprise-grade zero-trust security policies (Cloudflare), or simply want to forward an HTTP service (ngrok), those tools each have irreplaceable advantages.

Technology selection is always about trade-offs. I hope this comparison helps you make a more informed choice for your specific scenario.

---

*Shield CLI open source: https://github.com/fengyily/shield-cli*
*License: Apache 2.0*
