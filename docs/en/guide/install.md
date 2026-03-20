---
title: Install Shield CLI — macOS, Linux, Windows
description: Install Shield CLI via Homebrew, Scoop, apt, rpm, one-liner scripts, or build from source. Supports macOS, Linux, Windows on amd64, arm64, and 386 architectures.
head:
  - - meta
    - name: keywords
      content: Shield CLI install, Homebrew, Scoop, Linux, macOS, Windows, download, one-liner
---

# Installation

Shield CLI supports multiple installation methods. Choose the one that fits your operating system.

## macOS

### Homebrew (Recommended)

```bash
brew tap fengyily/tap
brew install shield-cli
```

### One-Liner

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

## Windows

### Scoop (Recommended)

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### PowerShell One-Liner

```powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```

## Linux

### One-Liner

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### Debian / Ubuntu

```bash
sudo dpkg -i shield-cli_<version>_linux_amd64.deb
```

### RHEL / CentOS

```bash
sudo rpm -i shield-cli_<version>_linux_amd64.rpm
```

Download packages from [GitHub Releases](https://github.com/fengyily/shield-cli/releases).

## China Mirror

If GitHub is slow in your region, use the jsDelivr CDN mirror:

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

## Build from Source

```bash
git clone https://github.com/fengyily/shield-cli.git
cd shield-cli
go build -o shield .
```

Requires Go 1.25.0 or later.

## Verify Installation

```bash
shield --version
```

If you see a version number, the installation was successful.

## Supported Platforms

| OS | Architectures |
|---|---|
| macOS | amd64, arm64 (Apple Silicon) |
| Linux | amd64, arm64, 386, armv7 |
| Windows | amd64, arm64, 386 |

## Next Steps

- [Quick Start](./quickstart.md)
