---
title: Custom Server — Use a Private Shield Gateway
description: Configure Shield CLI to connect to a private server instead of the default public gateway. Ideal for organizations with compliance requirements or private deployments.
head:
  - - meta
    - name: keywords
      content: Shield CLI custom server, private server, self-hosted, enterprise deployment, --server flag
---

# Custom Server

By default, Shield CLI connects to the `https://console.yishield.com/raas` public service. If you've deployed a private server, use the `--server` flag to point to it.

## Usage

```bash
shield ssh 10.0.0.5 --server https://your-server.com/raas
```

## Use Cases

- Your organization has deployed a private Shield server
- Data must not traverse the public internet
- Compliance requirements prohibit external SaaS services

## Notes

- The private server must be compatible with your Shield CLI version
- Ensure the machine running Shield CLI can reach the custom server address
