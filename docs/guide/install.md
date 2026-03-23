---
title: 安装 Shield CLI — macOS、Linux、Windows
description: 通过 Homebrew、Scoop、apt、rpm、一键脚本或源码编译安装 Shield CLI。支持 macOS、Linux、Windows，amd64 和 arm64 架构。
head:
  - - meta
    - name: keywords
      content: Shield CLI 安装, Homebrew, Scoop, Docker, Linux, macOS, Windows, 下载, 一键安装, 中国镜像
---

# 安装

Shield CLI 支持多种安装方式，选择适合你操作系统的方式即可。

## macOS

### Homebrew（推荐）

```bash
brew tap fengyily/tap
brew install shield-cli
```

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

## Windows

### Scoop（推荐）

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### PowerShell 一键安装

```powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```

## Linux

### 一键安装（推荐）

自动检测 apt / yum / dnf，添加仓库源并安装：

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash
```

安装后可通过 `apt upgrade` 或 `yum update` 自动获取新版本。

### APT（Debian / Ubuntu）

手动添加仓库源：

```bash
echo "deb [trusted=yes] https://fengyily.github.io/linux-repo/apt stable main" \
  | sudo tee /etc/apt/sources.list.d/shield-cli.list
sudo apt update
sudo apt install shield-cli
```

### YUM / DNF（RHEL / CentOS / Fedora）

```bash
sudo tee /etc/yum.repos.d/shield-cli.repo <<EOF
[shield-cli]
name=Shield CLI Repository
baseurl=https://fengyily.github.io/linux-repo/yum
enabled=1
gpgcheck=0
EOF
sudo yum install shield-cli   # 或 dnf install shield-cli
```

### 二进制直装

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### 手动安装 deb / rpm

从 [GitHub Releases](https://github.com/fengyily/shield-cli/releases) 下载对应的安装包：

```bash
# Debian / Ubuntu
sudo dpkg -i shield-cli_<version>_amd64.deb

# RHEL / CentOS
sudo rpm -i shield-cli_<version>_amd64.rpm
```

## 中国大陆镜像

如果 GitHub 访问较慢，可以使用 jsDelivr CDN 镜像：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

## Docker

```bash
# 使用预构建镜像（推荐）
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli

# 或从源码构建
docker build -t shield-cli https://github.com/fengyily/shield-cli.git
docker run -d --name shield --network host --restart unless-stopped shield-cli
```

`--network host` 让容器直接使用宿主机网络栈，可访问宿主机及内网资源。启动后访问 `http://localhost:8181`。

> **注意：** `--network host` 仅在 Linux 上生效。macOS / Windows Docker Desktop 请改用端口映射：
>
> ```bash
> docker run -d --name shield -p 8181:8181 --restart unless-stopped fengyily/shield-cli
> ```

## 从源码编译

```bash
git clone https://github.com/fengyily/shield-cli.git
cd shield-cli
go build -o shield .
```

需要 Go 1.25.0 或更高版本。

## 验证安装

```bash
shield --version
```

如果看到版本号输出，说明安装成功。

## 支持的平台和架构

| 操作系统 | 架构 |
|---|---|
| macOS | amd64, arm64 (Apple Silicon) |
| Linux | amd64, arm64, 386, armv7 |
| Windows | amd64, arm64, 386 |

## 下一步

- [5 分钟上手教程](./quickstart.md)
