---
title: SSH Tunnel — Browser-Based SSH Terminal
description: Access a full SSH terminal in your browser via Shield CLI. Supports password and private key authentication, SFTP file transfer, and xterm.js-based terminal with copy/paste.
head:
  - - meta
    - name: keywords
      content: SSH tunnel, SSH browser, web terminal, xterm.js, SFTP, Shield CLI SSH, remote SSH
---

# SSH

Access a full SSH terminal in your browser via Shield CLI, with password and private key authentication support.

## Quick Connect

```bash
# Connect to localhost
shield ssh

# Connect to a specific server
shield ssh 10.0.0.5

# Specify port
shield ssh 10.0.0.5:2222
```

## Authentication

### Password

```bash
shield ssh 10.0.0.5 --username root --auth-pass mypassword
```

Omit flags for interactive prompts.

### Private Key

```bash
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa
```

### Encrypted Private Key

```bash
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa --passphrase mypass
```

## SFTP File Transfer

Enable SFTP to upload and download files through the browser:

```bash
shield ssh 10.0.0.5 --enable-sftp
```

## Browser Terminal

Once connected, you get a full xterm.js-based terminal in the browser:

- Full terminal interaction (vim, top, etc.)
- Copy and paste support
- Auto-resizing to browser window

## Default Ports

| Input | Resolves To |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.5` | `10.0.0.5:22` |
