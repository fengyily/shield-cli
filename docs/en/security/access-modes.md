---
title: Access Modes — Visible and Invisible Tunnel URLs
description: Shield CLI offers Visible mode (public URL) and Invisible mode (authorization code required) to control Access URL visibility and security level.
head:
  - - meta
    - name: keywords
      content: Shield CLI access modes, visible mode, invisible mode, access control, authorization code, tunnel security
---

# Access Modes

Shield CLI offers two access modes to control the visibility and security level of Access URLs.

## Visible Mode (Default)

The Access URL is publicly accessible — anyone with the link can connect. Suitable for internal teams and trusted scenarios.

```bash
# Default is visible mode
shield ssh 10.0.0.5

# Specify access node (e.g., Hong Kong)
shield ssh 10.0.0.5 --visable=HK
```

### Characteristics

- Access URL can be shared directly — open and use
- Great for temporary sharing and internal collaboration
- Use `--visable=<node>` to select a nearby access node

## Invisible Mode

The Access URL requires an additional authorization code. Suitable for sensitive services.

```bash
shield ssh 10.0.0.5 --invisible
```

### Characteristics

- Generates both an Access URL and an Auth URL
- Visitors must first enter an authorization code via the Auth URL
- Ideal for exposing sensitive internal services externally

## Choosing a Mode

| Scenario | Recommended Mode |
|---|---|
| Internal team collaboration | Visible |
| Temporary sharing with colleagues | Visible |
| Exposing production services | Invisible |
| External client demos | Invisible |
