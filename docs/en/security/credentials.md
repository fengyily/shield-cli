---
title: Credential Management — AES-256-GCM Encrypted Storage
description: Shield CLI uses AES-256-GCM encryption with machine fingerprint-derived keys to protect locally stored credentials. Passwords are never stored and always masked in logs.
head:
  - - meta
    - name: keywords
      content: Shield CLI credentials, AES-256-GCM, machine fingerprint, encryption, security, credential storage
---

# Credentials

Shield CLI uses strong encryption to protect all locally stored credentials.

## Machine Fingerprint

Each machine generates a unique fingerprint used for:

- Deriving encryption keys
- Identifying the connector (format: `shield_<12-char-fingerprint>`)

The fingerprint is derived from machine hardware information using platform-specific methods, remaining consistent until the OS is reinstalled.

## Encryption

| Item | Detail |
|---|---|
| Algorithm | AES-256-GCM |
| Key | SHA256(machine fingerprint) |
| Purpose | Encrypt locally stored credentials and app configurations |

## Storage Location

| Platform | Path |
|---|---|
| macOS / Linux | `~/.shield-cli/.credential` |
| Windows | `%LOCALAPPDATA%\ShieldCLI\.credential` |

File permissions are set to `0600` — readable and writable only by the current user.

## What's Stored

The local credential file contains (encrypted):

- Connector name
- Connector token
- Assigned server port
- Server address

These are automatically generated on first connection and reused for subsequent connections.

## Password Security

- Passwords are hidden during interactive input
- All password content is masked in log output (shown as `***`)
- Passwords are not stored locally — used only during connection establishment

## Clearing Credentials

To reset credentials (e.g., switching accounts):

```bash
shield clean
```

This clears the cached credential file. New credentials will be generated on the next connection.

See [Clear Cache](../config/clean.md).
