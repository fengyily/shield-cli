---
title: 一个 CLI 工具的开源迭代记录：从单二进制到全平台分发
description: 记录 Shield CLI 在 v0.2.x 阶段的工程实践——Docker 容器化、Linux 包管理器分发、TCP/UDP 代理、系统服务集成。聚焦开源工具链选型和踩坑，而非产品本身。
date: 2026-03-23
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: 开源 CLI 工具, GoReleaser, Docker 多架构构建, APT YUM 仓库, GitHub Actions, 系统服务, Go 开源实践
---

# 一个 CLI 工具的开源迭代记录：从单二进制到全平台分发

> 这不是一篇产品介绍。这是一个用 Go 写的 CLI 工具在 v0.2.x 阶段的工程迭代记录——怎么把一个"能跑"的二进制，变成一个用户在任何平台都能用一行命令装上的工具。过程中用到的工具链、踩过的坑、做过的取舍，可能对同样在做开源 CLI 的人有参考价值。

---

## 背景

Shield CLI 是一个用 Go 写的内网穿透工具，核心功能是通过 Chisel 协议建立加密隧道，支持 SSH/RDP/VNC/HTTP 等协议的浏览器内访问。

v0.1.x 阶段做完了核心功能，通过 GoReleaser 交叉编译出 macOS / Linux / Windows 的二进制，放到 GitHub Release 上，用户 `curl | sh` 或 `brew install` 能装上。

但实际推出去之后发现，"能装"和"好装"之间还差着不少工程量。v0.2.x 主要在填这个坑。

---

## 一、Web UI 和系统服务：从命令行工具到常驻服务

### 问题

CLI 工具默认是前台进程，终端一关就断了。对于隧道这种需要长时间运行的服务，这不够用。

### 做法

v0.2.0 加了一个内嵌的 Web UI（`shield start` 启动，默认 `localhost:8181`），用浏览器管理多个应用连接。v0.2.1 接着做了系统服务注册：

```bash
# 注册为系统服务，开机自启
shield install

# 指定端口
shield install --port 8182
```

三个平台走的是不同的底层机制：

| 平台 | 机制 | 备注 |
|------|------|------|
| macOS | launchd 用户代理 | 不需要 sudo |
| Linux | systemd 服务 | 标准做法 |
| Windows | Windows Service | 需要管理员权限 |

同时 macOS 和 Windows 加了系统托盘图标，点击可以快速打开 Dashboard、重启、退出。

**踩坑记录**：系统托盘依赖 CGO（底层用到了各平台的原生 GUI 库），但 Linux 服务器通常没有桌面环境，也不需要托盘。所以 GoReleaser 配置拆成了两套：桌面版（macOS/Windows，CGO_ENABLED=1）和服务器版（Linux，纯 Go 编译）。这是一个在 CI 里调了很久才跑通的东西——交叉编译 + CGO 基本上是噩梦级别的组合，最后 macOS 和 Windows 各自在对应平台的 runner 上原生编译才解决。

---

## 二、Docker 容器化：看似简单，实际有坑

### 为什么要做

有用户反馈想在服务器上容器化部署，和已有的 Docker Compose 栈统一管理。

### Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o shield-cli .

FROM alpine:3.21
COPY --from=builder /app/shield-cli /usr/local/bin/shield
ENV SHIELD_LISTEN_HOST=0.0.0.0
ENTRYPOINT ["shield", "start"]
```

多阶段构建，最终镜像基于 Alpine，很常规。但有两个细节不常规：

**1. 必须用 `--network host`**

一般 Web 应用容器 `-p 8080:8080` 就行了。但内网穿透工具的核心功能是访问宿主机网络和内网资源——`10.0.0.0/24` 网段在 bridge 模式下不可达。`--network host` 让容器共享宿主机网络栈，这是这类工具容器化的必要条件。

**2. 监听地址的问题**

容器内 `127.0.0.1` 是容器自己的 loopback，外部流量进不来。所以 Docker 镜像默认把 `SHIELD_LISTEN_HOST` 设为 `0.0.0.0`。第一次上线时漏掉了这个，导致好几个用户反馈"容器启动了但访问不了"。

### CI 自动构建

GitHub Actions 里用 `docker/build-push-action` 做多架构构建（amd64 + arm64），同时推到 Docker Hub 和 GHCR。语义化标签（`latest`、`0.2.2`、`0.2`）通过 `docker/metadata-action` 自动生成。

这套流程现在是标准模板了，任何 Go 项目都可以直接抄：

```yaml
- uses: docker/metadata-action@v5
  with:
    images: |
      fengyily/shield-cli
      ghcr.io/fengyily/shield-cli
    tags: |
      type=semver,pattern={{version}}
      type=semver,pattern={{major}}.{{minor}}
      type=raw,value=latest
```

---

## 三、Linux 包管理器：APT 和 YUM 仓库搭建

### 为什么要做

用户已经可以通过 `curl | sh` 安装了，但 Linux 运维习惯的是 `apt install` / `yum install`。更重要的是，包管理器支持 `apt upgrade` 自动更新，不用手动重跑安装脚本。

### 技术方案

没有用第三方包管理托管服务（Packagecloud 等要收费），而是基于 **GitHub Pages** 自建仓库：

- APT 仓库：用 `dpkg-scanpackages` 生成 `Packages.gz`，用 GPG 签名 `Release` 文件
- YUM 仓库：用 `createrepo` 生成 `repodata/`，RPM 包用 GPG 签名

整个仓库托管在一个独立的 GitHub repo 的 `gh-pages` 分支上，GitHub Actions 在每次 Release 时自动把新的 deb/rpm 包推进去并重新生成索引。

用户配置仓库源：

```bash
# Debian / Ubuntu
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
sudo apt update && sudo apt install shield-cli

# RHEL / CentOS / Fedora
sudo tee /etc/yum.repos.d/shield-cli.repo <<EOF
[shield-cli]
name=Shield CLI Repository
baseurl=https://fengyily.github.io/linux-repo/yum
enabled=1
gpgcheck=0
EOF
sudo yum install shield-cli   # 或 dnf install shield-cli
```

### 还做了一个安装检测脚本

`install.sh` 加了 `--apt` / `--yum` 参数，自动检测系统类型并配置对应的包管理器源：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh -s -- --apt
```

---

## 四、TCP/UDP 端口代理：协议层的扩展

### 之前的状态

Shield CLI 最初只支持 SSH、RDP、VNC、HTTP、HTTPS、Telnet 这些有明确语义的协议——它在网关侧做协议渲染（比如 SSH 转成 Web Terminal，RDP 转成 HTML5 Canvas），所以每个协议需要对应的网关支持。

### 新需求

有用户需要代理 MySQL（3306）、Redis（6379）、PostgreSQL（5432）等 TCP 服务。这些不需要浏览器渲染，纯端口转发就够了。还有少量 DNS（53）、Syslog 等 UDP 场景。

### 实现

```bash
# TCP 代理
shield tcp 3306              # 本地 MySQL
shield tcp 10.0.0.5:6379     # 远程 Redis

# UDP 代理
shield udp 53                # DNS
shield udp 10.0.0.5:514      # Syslog
```

技术上，TCP 走 Chisel 的标准反向隧道；UDP 用 Chisel 原生的 `/udp` 后缀做 UDP over WebSocket 转发。

和 SSH/RDP 等协议的关键区别：TCP/UDP **没有默认端口**，所以 CLI 强制要求用户指定端口号。隧道建立后不会自动打开浏览器（因为没有 Web UI 可看），而是在终端打印连接指南：

```
  📡 Connection Guide (TCP port proxy):
    gateway.example.com:48721  →  10.0.0.5:3306

    Examples:
      mysql -h gateway.example.com -P 48721 -u root
      redis-cli -h gateway.example.com -p 48721
```

---

## 五、分发矩阵总结

做完 v0.2.x 这一轮之后，Shield CLI 的安装方式变成了这样：

| 方式 | 命令 | 平台 |
|------|------|------|
| Homebrew | `brew install shield-cli` | macOS |
| Scoop | `scoop install shield-cli` | Windows |
| APT | `apt install shield-cli` | Debian/Ubuntu |
| YUM | `yum install shield-cli` | RHEL/CentOS/Fedora |
| Docker | `docker run fengyily/shield-cli` | 任意 Linux |
| curl | `curl ... \| sh` | macOS/Linux |
| PowerShell | 一键脚本 | Windows |
| 二进制 | GitHub Release 下载 | 全平台 |

这些不是一次性做完的，是在半个月内迭代加上去的。每一种安装方式背后都有对应的 CI 流水线在维护——GoReleaser 管二进制 + Homebrew + Scoop + deb/rpm 包，Docker 走单独的 workflow，APT/YUM 仓库也是独立的 workflow。

---

## 开源工具链清单

把这次用到的工具列一下，都是开源的，别的 Go 项目可以直接参考：

| 工具 | 用途 |
|------|------|
| [GoReleaser](https://goreleaser.com/) | 交叉编译 + 打包 + Homebrew/Scoop/deb/rpm |
| [docker/build-push-action](https://github.com/docker/build-push-action) | 多架构 Docker 构建和推送 |
| [docker/metadata-action](https://github.com/docker/metadata-action) | Docker 标签自动生成 |
| [nfpm](https://nfpm.goreleaser.com/) | 不需要 dpkg-deb 也能打 deb/rpm 包（GoReleaser 内置） |
| GitHub Pages | APT/YUM 仓库静态托管 |
| GitHub Actions Secrets | GPG 密钥、Docker 凭证管理 |

---

## 几个教训

**1. 不要低估分发的工作量。** 核心功能可能只占 40% 的工作量，剩下的全在"让用户装得上"这件事上。Homebrew tap 配置、Scoop bucket manifest、deb/rpm 打包参数、Docker 多架构、APT/YUM 仓库签名、install.sh 的各种 edge case……每一个都不难，但加起来很花时间。

**2. CGO 和交叉编译是两个互斥的目标。** 如果你的 Go 项目依赖 CGO（GUI 库、SQLite 等），老老实实在目标平台上原生编译。不要试图在 Linux CI 上交叉编译 macOS 的 CGO 项目，那条路走不通。

**3. Docker + 内网穿透 = `--network host`。** 这个组合很违反直觉，因为 Docker 的核心价值之一是网络隔离。但对于需要访问宿主机网络的工具，host 模式是唯一选择。在文档里一定要把这个说清楚，否则用户的第一反应永远是"为什么容器里连不上"。

**4. GitHub Pages 做包仓库足够用了。** 不需要 Packagecloud，不需要 Artifactory。一个 gh-pages 分支 + GitHub Actions 自动更新索引，对于中小型开源项目完全够。省钱，可控。

---

## 最后

这篇文章记录的不是 Shield CLI 本身的功能，而是一个开源 CLI 工具在 **分发和部署** 层面的工程实践。如果你也在做一个 Go CLI 项目，正在纠结怎么让用户更方便地安装和运行，这里面的工具链和踩坑经验应该能帮上忙。

所有构建配置和 CI 脚本都在仓库里，可以直接参考：

- GoReleaser 配置：`.goreleaser.yaml`
- Docker 构建：`Dockerfile` + `.github/workflows/docker.yml`
- APT/YUM 仓库：`.github/workflows/update-repo.yml`

项目地址：https://github.com/fengyily/shield-cli

如果这篇文章对你有帮助，去仓库点个 Star 就是最好的支持。用的过程中遇到问题或者有想法，欢迎直接提 [Issue](https://github.com/fengyily/shield-cli/issues)。
