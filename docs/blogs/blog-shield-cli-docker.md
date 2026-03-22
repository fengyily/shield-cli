---
title: 用 Docker 跑 Shield CLI：一行命令把内网穿透丢进容器
description: Shield CLI 现在支持 Docker 部署，一行 docker run 即可启动，配合 host 网络模式直接访问宿主机和内网资源。本文记录实际部署过程和踩过的坑。
date: 2025-01-25
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI Docker, 内网穿透 Docker, Docker 部署隧道工具, host 网络模式, 容器化远程访问
---

# 用 Docker 跑 Shield CLI：一行命令把内网穿透丢进容器

> 之前写过 Shield CLI 的介绍和对比，这次聊一个实际部署的问题：怎么用 Docker 跑起来，以及为什么内网穿透工具跑在容器里没那么简单。

---

## 为什么要用 Docker

Shield CLI 本身是个单二进制文件，`curl | sh` 就能装，按理说不需要 Docker。但实际用下来，有几个场景确实容器化更合适：

**1. 服务器上不想污染环境**

生产服务器上已经跑了一堆东西，装一个新的二进制文件不算大事，但牵扯到开机自启、systemd 配置、升级维护，事情就多了。Docker 的好处是隔离干净，`docker rm` 一条命令删得一干二净。

**2. 统一部署流程**

如果团队已经在用 Docker Compose 或 K8s 管理服务，Shield CLI 作为一个基础设施组件，走容器化部署和其他服务保持一致，运维不用单独记一套流程。

**3. 快速试用**

新同事想试一下 Shield CLI，不用管 Go 版本、不用管操作系统，`docker run` 一把就起来了。

---

## 实际跑起来

最简单的方式：

```bash
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli
```

打开 `http://localhost:8181`，Web UI 就出来了，和直接装在宿主机上体验完全一样。

如果想指定端口：

```bash
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli \
  shield start 9090
```

---

## 关于 `--network host`

这是跑 Shield CLI 容器时最关键的一个参数，也是和普通 Web 应用容器化最大的区别。

一般的 Web 应用跑在容器里，用 `-p 8080:8080` 映射端口就行了，因为它只需要对外提供 HTTP 服务。但 Shield CLI 的核心功能是 **访问宿主机网络和内网资源** —— 你要用它连接 `10.0.0.5` 的 RDP、`192.168.1.100` 的 SSH，这些地址在默认的 bridge 网络模式下是不可达的。

`--network host` 让容器直接使用宿主机的网络栈，不做网络隔离。对于 Shield CLI 这种场景，这是必要的。

用一张图看区别：

```
默认 bridge 模式：
  容器 → docker0 网桥 → 宿主机 → 内网
  ❌ 容器看不到 10.0.0.0/24 网段

host 模式：
  容器 ≡ 宿主机（共享网络栈）
  ✅ 容器可以直接访问宿主机能访问的一切
```

---

## macOS / Windows 用户注意

`--network host` **只在 Linux 上生效**。

Docker Desktop 在 macOS 和 Windows 上是跑在一个 Linux 虚拟机里的，`--network host` 绑定的是那个虚拟机的网络，不是你的宿主机网络。所以在 Mac 上加了这个参数也没用。

macOS / Windows 下只能用端口映射：

```bash
docker run -d --name shield \
  -p 8181:8181 \
  --restart unless-stopped \
  fengyily/shield-cli
```

这种模式下 Web UI 可以正常使用，但 Shield CLI 只能访问容器自身网络能到达的地址。如果你的目标服务在宿主机本地（比如 `127.0.0.1:22`），需要改成 `host.docker.internal:22`。

老实说，macOS / Windows 下直接安装二进制文件体验更好，Docker 方式更适合 Linux 服务器。

---

## 踩过的一个坑：容器里监听地址的问题

第一次跑的时候发现一个问题：容器启动了，端口映射也做了，但 `curl localhost:8181` 就是连不上。

排查后发现 Shield CLI 默认监听 `127.0.0.1:8181`，这在宿主机上没问题，但在容器里就出事了 —— `127.0.0.1` 是容器自己的 loopback，外部流量（包括 Docker 的端口映射）进不来。

解决方式是通过环境变量 `SHIELD_LISTEN_HOST` 控制监听地址。Docker 镜像里已经默认设置了 `0.0.0.0`，所以直接拉官方镜像不会遇到这个问题。如果你是自己构建的，注意加上这个环境变量：

```bash
docker run -d --name shield \
  -e SHIELD_LISTEN_HOST=0.0.0.0 \
  --network host \
  fengyily/shield-cli
```

---

## Docker Compose 示例

如果你习惯用 Compose 管理服务：

```yaml
services:
  shield:
    image: fengyily/shield-cli
    container_name: shield
    network_mode: host
    restart: unless-stopped
```

就这么几行，`docker compose up -d` 搞定。

---

## 镜像细节

简单说一下镜像本身：

- **基础镜像**：`alpine:3.21`，最终镜像很小
- **多架构**：同时提供 `linux/amd64` 和 `linux/arm64`，x86 服务器和 ARM 机器（树莓派、Oracle Cloud ARM 实例等）都能跑
- **仓库地址**：`fengyily/shield-cli`（Docker Hub）和 `ghcr.io/fengyily/shield-cli`（GitHub Container Registry）都有，国内拉 Docker Hub 可能更快

---

## 什么时候该用 Docker，什么时候不该

| 场景 | 建议 |
|---|---|
| Linux 服务器长期运行 | Docker + `--network host`，升级方便 |
| 已有 Docker Compose / K8s 环境 | Docker，保持部署方式统一 |
| 快速试用 | Docker，一行命令起来 |
| macOS / Windows 日常使用 | 直接安装二进制文件，体验更好 |
| 需要系统托盘图标 | 直接安装，容器里没有桌面环境 |

---

## 最后

容器化不是目的，减少折腾才是。Shield CLI 的 Docker 支持解决的核心问题就两个：部署标准化和环境隔离。如果你本来就在 Linux 服务器上跑，从 `curl | sh` 换成 `docker run` 基本没有额外成本，还省了开机自启的配置。

项目地址：https://github.com/fengyily/shield-cli

之前的文章：
- [Shield CLI 与 ngrok、frp、Cloudflare Tunnel 的技术对比](./blog-shield-cli-vs-tunnels.md)
- [推荐一个思路不太一样的内网穿透工具](./zhihu-recommend-shield-cli.md)
