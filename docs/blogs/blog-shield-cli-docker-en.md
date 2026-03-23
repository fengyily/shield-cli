---
title: "Running Shield CLI with Docker: One Command to Containerize Your Tunnel"
description: Shield CLI now supports Docker deployment. A single docker run command gets you started, and with host networking mode you can access the host machine and LAN resources directly. This post covers the actual deployment process and gotchas we ran into.
date: 2025-01-25
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI Docker, tunnel Docker deployment, host network mode, containerized remote access, Docker networking
---

# Running Shield CLI with Docker: One Command to Containerize Your Tunnel

> We've previously introduced Shield CLI and compared it with other tools. This time, let's talk about a practical deployment topic: how to run it in Docker, and why containerizing a tunnel tool isn't as straightforward as it seems.

---

## Why Docker?

Shield CLI is a single binary — a quick `curl | sh` and you're done. In theory, Docker isn't necessary. But in practice, there are a few scenarios where containerization makes more sense:

**1. Keep production servers clean**

Production servers already have plenty of things running. Installing another binary isn't a big deal, but then comes auto-start configuration, systemd setup, upgrade maintenance — it adds up. With Docker, everything is cleanly isolated. A single `docker rm` wipes the slate clean.

**2. Unified deployment workflow**

If your team already uses Docker Compose or Kubernetes, deploying Shield CLI as a container keeps it consistent with the rest of your infrastructure. No special procedures to remember.

**3. Quick trial runs**

A new team member wants to try Shield CLI? No need to worry about Go versions or OS compatibility — just `docker run` and it's up.

---

## Getting It Running

The simplest approach:

```bash
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli
```

Open `http://localhost:8181` and the Web UI is right there — the same experience as a native install.

To use a custom port:

```bash
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli \
  shield start 9090
```

---

## About `--network host`

This is the most critical parameter when running Shield CLI in a container, and the biggest difference from containerizing a typical web application.

A regular web app in a container just needs `-p 8080:8080` for port mapping, since it only serves HTTP traffic. But Shield CLI's core function is **accessing the host network and LAN resources** — you need it to connect to RDP at `10.0.0.5`, SSH at `192.168.1.100`, and those addresses are unreachable under the default bridge network mode.

`--network host` lets the container use the host's network stack directly, with no network isolation. For Shield CLI's use case, this is essential.

Here's the difference at a glance:

```
Default bridge mode:
  Container → docker0 bridge → Host → LAN
  ❌ Container cannot see the 10.0.0.0/24 subnet

Host mode:
  Container ≡ Host (shared network stack)
  ✅ Container can reach everything the host can
```

---

## Note for macOS / Windows Users

`--network host` **only works on Linux**.

Docker Desktop on macOS and Windows runs inside a Linux VM. `--network host` binds to that VM's network, not your actual host network. So on a Mac, this flag has no effect.

On macOS / Windows, use port mapping instead:

```bash
docker run -d --name shield \
  -p 8181:8181 \
  --restart unless-stopped \
  fengyily/shield-cli
```

In this mode the Web UI works fine, but Shield CLI can only reach addresses accessible from the container's own network. If your target service is on the host machine (e.g., `127.0.0.1:22`), use `host.docker.internal:22` instead.

Honestly, on macOS / Windows a native binary install gives a better experience. Docker is better suited for Linux servers.

---

## A Gotcha: Listen Address Inside Containers

The first time we ran it, we hit an issue: the container was up, port mapping was in place, but `curl localhost:8181` just wouldn't connect.

After debugging, we found that Shield CLI defaults to listening on `127.0.0.1:8181`. That's fine on the host, but inside a container it's a problem — `127.0.0.1` is the container's own loopback, so external traffic (including Docker's port mapping) can't reach it.

The fix is the `SHIELD_LISTEN_HOST` environment variable. The official Docker image already defaults to `0.0.0.0`, so you won't hit this with the published image. But if you're building your own, remember to set it:

```bash
docker run -d --name shield \
  -e SHIELD_LISTEN_HOST=0.0.0.0 \
  --network host \
  fengyily/shield-cli
```

---

## Docker Compose Example

If you prefer Compose:

```yaml
services:
  shield:
    image: fengyily/shield-cli
    container_name: shield
    network_mode: host
    restart: unless-stopped
```

Just a few lines. `docker compose up -d` and you're done.

---

## Image Details

A quick overview of the image itself:

- **Base image**: `alpine:3.21` — keeps the final image small
- **Multi-architecture**: both `linux/amd64` and `linux/arm64` are available — works on x86 servers and ARM machines (Raspberry Pi, Oracle Cloud ARM instances, etc.)
- **Registries**: available on both `fengyily/shield-cli` (Docker Hub) and `ghcr.io/fengyily/shield-cli` (GitHub Container Registry)

---

## When to Use Docker, When Not To

| Scenario | Recommendation |
|---|---|
| Long-running Linux server | Docker + `--network host` — easy upgrades |
| Existing Docker Compose / K8s setup | Docker — keep deployment consistent |
| Quick trial | Docker — one command to get started |
| Daily use on macOS / Windows | Install the binary directly — better experience |
| Need system tray icon | Install directly — no desktop environment in containers |

---

## Final Thoughts

Containerization isn't the goal — reducing hassle is. Docker support for Shield CLI solves two core problems: standardized deployment and environment isolation. If you're already running on a Linux server, switching from `curl | sh` to `docker run` costs virtually nothing, and you save yourself the auto-start configuration.

Project: https://github.com/fengyily/shield-cli

Previous posts:
- [Shield CLI vs ngrok, frp, and Cloudflare Tunnel: A Technical Comparison](./blog-shield-cli-vs-tunnels-en.md)
