---
title: Install Shield CLI — macOS, Linux, Windows
description: Install Shield CLI via Homebrew, Scoop, apt, rpm, one-liner scripts, or build from source. Supports macOS, Linux, Windows on amd64, arm64, and 386 architectures.
head:
  - - meta
    - name: keywords
      content: Shield CLI install, Homebrew, Scoop, Docker, Linux, macOS, Windows, download, one-liner
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

### One-Click Installer (Recommended for beginners)

Download [`install.bat`](https://raw.githubusercontent.com/fengyily/shield-cli/main/install.bat) and double-click to run. It will automatically:

1. Detect your system architecture (AMD64 / ARM64)
2. Download the correct version
3. Install as a system service (auto-start on boot)
4. Open the Web UI in your browser

> A UAC permission prompt will appear — click "Yes" to continue.

### Scoop

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### PowerShell One-Liner

```powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```

## Linux

### One-Liner (Recommended)

Automatically detects apt / yum / dnf, adds the repository, and installs:

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash
```

Once installed, future updates are available via `apt upgrade` or `yum update`.

### APT (Debian / Ubuntu)

Manually add the repository:

```bash
echo "deb [trusted=yes] https://fengyily.github.io/linux-repo/apt stable main" \
  | sudo tee /etc/apt/sources.list.d/shield-cli.list
sudo apt update
sudo apt install shield-cli
```

### YUM / DNF (RHEL / CentOS / Fedora)

```bash
sudo tee /etc/yum.repos.d/shield-cli.repo <<EOF
[shield-cli]
name=Shield CLI Repository
baseurl=https://fengyily.github.io/linux-repo/yum
enabled=1
gpgcheck=0
EOF
sudo yum install shield-cli   # or: dnf install shield-cli
```

### Binary Install

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### Manual deb / rpm Install

Download packages from [GitHub Releases](https://github.com/fengyily/shield-cli/releases):

```bash
# Debian / Ubuntu
sudo dpkg -i shield-cli_<version>_amd64.deb

# RHEL / CentOS
sudo rpm -i shield-cli_<version>_amd64.rpm
```

## China Mirror

If GitHub is slow in your region, use the jsDelivr CDN mirror:

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

## Docker

```bash
# Use the prebuilt image (recommended)
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli

# Or build from source
docker build -t shield-cli https://github.com/fengyily/shield-cli.git
docker run -d --name shield --network host --restart unless-stopped shield-cli
```

`--network host` shares the host's network stack so Shield CLI can reach local and LAN services. Open `http://localhost:8181` after startup.

> **Note:** `--network host` only works on Linux. On macOS / Windows Docker Desktop, use port mapping instead:
>
> ```bash
> docker run -d --name shield -p 8181:8181 --restart unless-stopped fengyily/shield-cli
> ```

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
