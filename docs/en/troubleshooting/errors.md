---
title: Common Errors — Shield CLI Troubleshooting Guide
description: Solutions for common Shield CLI errors including connection timeout, authentication failures, port conflicts, and installation issues.
head:
  - - meta
    - name: keywords
      content: Shield CLI errors, troubleshooting, connection timeout, authentication failed, port in use, installation error
---

# Common Errors

## Connection Issues

### Connection Timeout

**Symptom:** Stuck on "Connecting" status for a long time

**Possible causes:**
- Target service is not running or port is incorrect
- Machine cannot reach the public gateway
- Firewall is blocking WebSocket connections

**Solutions:**
1. Verify target service is running: `telnet <ip> <port>`
2. Check connectivity to `console.yishield.com`
3. Ensure firewall allows outbound connections on port 62888

### Authentication Failed (401)

**Symptom:** Authentication failure or 401 error

**Possible cause:**
- Local credentials are expired or corrupted

**Solution:**
```bash
shield clean
```
Clear credentials and reconnect — new credentials will be generated automatically.

### Target Service Auth Failed

**Symptom:** Tunnel established successfully, but browser shows an authentication error

**Possible causes:**
- Wrong username or password
- SSH private key mismatch

**Solutions:**
- Verify username and password are correct
- For private keys, confirm the file path is correct and permissions are `600`

### Port Already in Use

**Symptom:** `shield start` fails to launch

**Possible cause:**
- Port 8181 is already in use by another program

**Solution:**
```bash
shield start 9090  # Use a different port
```

## Installation Issues

### Homebrew Install Failed

**Solution:**
```bash
brew update
brew tap fengyily/tap
brew install shield-cli
```

If it still fails, try the one-liner script:
```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### Permission Denied

**Symptom:** Permission denied when running on Linux

**Solution:**
```bash
chmod +x ./shield
```
