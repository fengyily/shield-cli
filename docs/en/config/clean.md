---
title: Clear Cache — Reset Shield CLI Credentials
description: Use shield clean to clear locally cached credentials. New credentials are auto-generated on next connection. App configurations are not affected.
head:
  - - meta
    - name: keywords
      content: Shield CLI clean, clear cache, reset credentials, troubleshooting
---

# Clear Cache

The `shield clean` command clears locally cached credential information.

## Usage

```bash
shield clean
```

## What Gets Cleared

- Local credential file (connector name, token, etc.)
- New credentials will be automatically generated on the next connection

## When to Use

- Switched server accounts
- Credential file is corrupted and causing connection failures
- Need to reset connector identity
- Troubleshooting connection issues

## Notes

- Saved app configurations are **not** affected
- New credentials are generated automatically on the next connection — no extra steps needed
