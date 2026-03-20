---
title: HTTP/HTTPS Tunnel — Expose Local Web Apps
description: Expose local or internal web applications to the public internet via Shield CLI HTTP/HTTPS tunnel. Preserves headers, cookies, and WebSocket connections.
head:
  - - meta
    - name: keywords
      content: HTTP tunnel, HTTPS tunnel, expose local web app, reverse proxy, localhost tunnel, Shield CLI HTTP
---

# HTTP / HTTPS

Expose local or internal web applications to the public internet via Shield CLI.

## Quick Connect

```bash
# HTTP - local port 80
shield http

# Expose local dev server
shield http 3000

# Expose internal web app
shield http 10.0.0.5:8080

# HTTPS
shield https
shield https 10.0.0.5:443
```

## Use Cases

### Local Development Preview

Share your in-progress web app with colleagues or clients:

```bash
# React / Vue / Next.js dev server
shield http 3000

# Python Flask / Django
shield http 5000

# Any local port
shield http 8080
```

### Internal App Access

Give external users temporary access to admin panels, dashboards, etc.:

```bash
shield http 10.0.0.5:8080
```

### HTTPS Services

If the target service uses HTTPS:

```bash
shield https 10.0.0.5:443
```

## Features

- Full HTTP request proxying with original headers and cookies preserved
- WebSocket support
- Automatic public HTTPS access URL

## Default Ports

| Input | Resolves To |
|---|---|
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield https` | `127.0.0.1:443` |
