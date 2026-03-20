---
title: 安装 Shield CLI — macOS、Linux、Windows
description: 通过 Homebrew、Scoop、apt、rpm、一键脚本或源码编译安装 Shield CLI。支持 macOS、Linux、Windows，amd64 和 arm64 架构。
head:
  - - meta
    - name: keywords
      content: Shield CLI 安装, Homebrew, Scoop, Linux, macOS, Windows, 下载, 一键安装, 中国镜像
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

### 一键安装

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

从 [GitHub Releases](https://github.com/fengyily/shield-cli/releases) 下载对应的安装包。

## 中国大陆镜像

如果 GitHub 访问较慢，可以使用 jsDelivr CDN 镜像：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

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
