---
title: 常见问题 — Shield CLI FAQ
description: Shield CLI 常见问题解答，包括免费使用、平台支持、并发连接、安全机制、网络要求和配置存储等。
head:
  - - meta
    - name: keywords
      content: Shield CLI FAQ, 常见问题, 帮助, 故障排查, 技术支持
  - - script
    - type: application/ld+json
    - |
      {
        "@context": "https://schema.org",
        "@type": "FAQPage",
        "mainEntity": [
          {
            "@type": "Question",
            "name": "Shield CLI 是免费的吗？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Shield CLI 客户端是开源免费的。公共隧道服务有免费额度，具体请查看官网。"
            }
          },
          {
            "@type": "Question",
            "name": "支持哪些操作系统？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "macOS、Linux、Windows，支持 amd64、arm64、386、armv7 架构。"
            }
          },
          {
            "@type": "Question",
            "name": "需要安装什么依赖吗？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "不需要。Shield CLI 是单一可执行文件，下载即可运行。"
            }
          },
          {
            "@type": "Question",
            "name": "能同时连接几个应用？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Web UI 模式下最多同时连接 3 个应用。命令行模式下每次连接 1 个。"
            }
          },
          {
            "@type": "Question",
            "name": "连接断了会自动重连吗？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "会。Shield CLI 内置自动重连机制，使用指数退避策略，最大间隔 10 秒。"
            }
          },
          {
            "@type": "Question",
            "name": "数据传输安全吗？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "所有数据通过 WebSocket 加密隧道传输。本地凭证使用 AES-256-GCM 加密存储，密钥基于机器指纹派生。"
            }
          },
          {
            "@type": "Question",
            "name": "密码会被存储到服务端吗？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "不会。目标服务的密码仅在建立连接时使用，不会持久化存储在服务端。"
            }
          },
          {
            "@type": "Question",
            "name": "防火墙需要开放哪些端口？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Shield CLI 只需要出站连接到公网网关的 62888 端口（WebSocket）和 443 端口（HTTPS API），不需要开放任何入站端口。"
            }
          },
          {
            "@type": "Question",
            "name": "最多能保存几个应用？",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "最多 10 个应用配置，每个使用 AES-256-GCM 加密本地存储。"
            }
          }
        ]
      }
---

# 常见问题

## 基本问题

### Shield CLI 是免费的吗？

Shield CLI 客户端是开源免费的。公共隧道服务有免费额度，具体请查看 [官网](https://console.yishield.com)。

### 支持哪些操作系统？

macOS、Linux、Windows，支持 amd64、arm64、386、armv7 架构。

### 需要安装什么依赖吗？

不需要。Shield CLI 是单一可执行文件，下载即可运行。

## 连接问题

### 能同时连接几个应用？

Web UI 模式下最多同时连接 **3 个** 应用。命令行模式下每次连接 1 个。

### 连接断了会自动重连吗？

会。Shield CLI 内置自动重连机制，使用指数退避策略，最大间隔 10 秒。

### Access URL 有有效期吗？

Access URL 在隧道连接期间有效。断开连接后，URL 即失效。

### 多人可以同时通过同一个 URL 访问吗？

可以。同一个 Access URL 可以被多个浏览器同时访问。

## 安全问题

### 数据传输安全吗？

所有数据通过 WebSocket 加密隧道传输。本地凭证使用 AES-256-GCM 加密存储。

### 密码会被存储到服务端吗？

不会。目标服务的密码仅在建立连接时使用，不会持久化存储在服务端。

### 机器指纹是什么？

基于你机器硬件信息生成的唯一标识，用于派生加密密钥和标识连接器身份。不包含个人隐私信息。

## 网络问题

### 中国大陆访问 GitHub 慢怎么办？

使用 jsDelivr CDN 镜像安装：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

### 在公司网络（有代理）下能用吗？

Shield CLI 使用 WebSocket 建立隧道，需要确保网络允许 WebSocket 连接到公网。如果有 HTTP 代理，大多数情况下可以正常工作。

### 防火墙需要开放哪些端口？

Shield CLI 只需要**出站**连接到公网网关的 62888 端口（WebSocket），不需要开放任何入站端口。

## 配置问题

### 应用配置存在哪里？

加密存储在本地：
- macOS / Linux：`~/.shield-cli/`
- Windows：`%LOCALAPPDATA%\ShieldCLI\`

### 如何在多台机器之间同步配置？

目前不支持云同步。每台机器的配置和凭证是独立的（基于各自的机器指纹加密）。

### 最多能保存几个应用？

10 个。
