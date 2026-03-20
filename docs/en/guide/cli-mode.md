---
title: CLI Mode — Terminal-Based Tunnel Creation
description: Use Shield CLI in command-line mode for servers, scripting, and automation. Create SSH, RDP, VNC, HTTP tunnels with smart defaults and flexible authentication.
head:
  - - meta
    - name: keywords
      content: Shield CLI command line, CLI mode, terminal, scripting, automation, SSH tunnel command
---

# CLI Mode

CLI mode is ideal for server environments, scripting, or when you just need to quickly create a single tunnel.

## Basic Usage

```bash
shield <protocol> [address]
```

### Examples

```bash
# Connect to local SSH (127.0.0.1:22)
shield ssh

# Connect to a remote RDP
shield rdp 10.0.0.5

# Expose a local web app on port 3000
shield http 3000

# Connect to VNC with specific IP and port
shield vnc 10.0.0.10:5901
```

## Authentication

### Interactive Input (Default)

Without auth flags, Shield CLI prompts interactively:

```bash
shield ssh 10.0.0.5
# → Prompts for username
# → Prompts for password (hidden input)
```

### Command-Line Flags

```bash
# Username + password
shield ssh 10.0.0.5 --username root --auth-pass mypassword

# SSH private key
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa

# Encrypted private key
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa --passphrase mypass
```

## Access Modes

```bash
# Visible mode (default) — public Access URL
shield ssh 10.0.0.5

# Specify access node
shield ssh 10.0.0.5 --visable=HK

# Invisible mode — requires authorization code
shield ssh 10.0.0.5 --invisible
```

See [Access Modes](../security/access-modes.md) for details.

## SFTP Support

Enable SFTP file transfer with SSH:

```bash
shield ssh 10.0.0.5 --enable-sftp
```

## Connection Output

After establishing a tunnel, the terminal displays:

```
Shield CLI v1.x.x

Protocol: SSH
Target:   10.0.0.5:22
Status:   Connected

Access URL: https://xxxxx.yishield.com
```

Share the Access URL with anyone who needs access.

## Exit

Press `Ctrl+C` to disconnect and exit.

## Best For

- Headless Linux servers
- Shell scripts and automation
- Quick single-tunnel connections
- CI/CD environments

## Next Steps

- [Web UI Mode](./web-ui.md) — graphical multi-app management
- [Commands Reference](../reference/commands.md) — full parameter list
