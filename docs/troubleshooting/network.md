---
title: 网络问题 — 防火墙、代理和连通性排查
description: 排查 Shield CLI 网络问题，涵盖连通性检查、代理环境、防火墙配置、断线重连和中国大陆优化。
head:
  - - meta
    - name: keywords
      content: Shield CLI 网络, 防火墙, 代理, WebSocket, 连通性, 中国镜像, 断线重连
---

# 网络问题

## 连通性检查

### 步骤一：检查目标服务

确认目标服务可以从 Shield CLI 所在机器访问：

```bash
# 检查端口连通性
telnet <目标IP> <端口>

# 或使用 nc
nc -zv <目标IP> <端口>
```

### 步骤二：检查公网连通性

确认可以连接到 Shield 公网网关：

```bash
curl -I https://console.yishield.com
```

### 步骤三：检查 WebSocket

Shield CLI 使用 WebSocket 建立隧道。如果你在代理或防火墙后面，确保 WebSocket 连接未被拦截。

## 代理环境

如果你的网络需要通过 HTTP 代理访问外网，大多数情况下 Shield CLI 可以正常工作，因为 WebSocket 可以通过 HTTP CONNECT 方法穿过代理。

如果仍然无法连接，联系网络管理员确认代理是否允许 WebSocket 连接。

## 防火墙配置

Shield CLI **仅需出站连接**，不需要开放任何入站端口：

| 方向 | 目标 | 端口 | 协议 |
|---|---|---|---|
| 出站 | console.yishield.com | 62888 | WebSocket (WSS) |
| 出站 | console.yishield.com | 443 | HTTPS (API) |

## 断线重连

Shield CLI 内置自动重连机制：

- 检测到断线后立即重试
- 使用指数退避：1s → 2s → 4s → 8s → 10s（最大）
- 网络恢复后自动重新建立隧道
- 无需手动干预

## 中国大陆优化

如果在中国大陆使用，可能会遇到 GitHub 访问较慢的问题：

**安装时使用镜像：**
```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

**使用就近节点：**
```bash
shield ssh 10.0.0.5 --visable=HK
```
