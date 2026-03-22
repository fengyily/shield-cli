---
title: Shield CLI 与 ngrok、frp、Cloudflare Tunnel 的技术对比
description: 深入对比 Shield CLI 与 ngrok、frp、Cloudflare Tunnel 在协议支持、配置复杂度、安全模型、架构和价格方面的差异。Shield CLI 是唯一支持浏览器内 RDP/VNC/SSH 渲染的隧道工具。
date: 2025-01-15
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI vs ngrok, Shield CLI vs frp, Shield CLI vs Cloudflare Tunnel, 隧道工具对比, 内网穿透对比, RDP浏览器
  - - script
    - type: application/ld+json
    - |
      {
        "@context": "https://schema.org",
        "@type": "Article",
        "headline": "Shield CLI 与主流隧道工具的技术对比",
        "description": "深入对比 Shield CLI 与 ngrok、frp、Cloudflare Tunnel 在协议支持、配置复杂度、安全模型、架构和价格方面的差异",
        "author": {"@type": "Organization", "name": "Shield CLI"},
        "publisher": {"@type": "Organization", "name": "Shield CLI", "url": "https://docs.yishield.com"},
        "datePublished": "2025-01-15"
      }
---

# 我用一条命令把内网的 RDP 桌面开到了浏览器里 —— Shield CLI 与主流隧道工具的技术对比

> 最近在折腾远程运维的方案，需要把内网的 Windows 远程桌面、Linux SSH 暴露给外部协作者使用。试了 ngrok、frp、Cloudflare Tunnel 之后，发现了一个思路不太一样的工具 Shield CLI，花了些时间深入对比，记录一下技术细节。

## 先说场景：为什么"隧道"不等于"远程访问"

大多数人用隧道工具的典型场景是：本地跑了个 Web 服务，想给外部临时访问一下。ngrok 在这个场景下几乎是标配。

但如果你的场景是这样的：

- 需要让客户通过浏览器直接操作一台内网 Windows 的远程桌面（RDP）
- 给外包团队临时开一个 SSH 终端，但不想让他们装任何客户端
- 演示环境中的 VNC 桌面需要一个可分享的链接

这时候你会发现，ngrok 能打通 TCP 隧道没错，但对方还是得装 RDP 客户端、配置 SSH 工具。**隧道解决的是"网络可达"，但没解决"终端可用"**。

Shield CLI 的做法是：隧道打通之后，在网关侧直接提供 HTML5 的 Web 终端（基于 Apache Guacamole 等协议网关），用户拿到一个 HTTPS 链接，浏览器打开就是 RDP 桌面或 SSH 终端。

这是两个产品本质思路的区别。下面进入逐项对比。

---

## 一、协议支持：谁真正做了"远程桌面"

| 协议 | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|------|-----------|-------|-----|-------------------|
| HTTP/HTTPS | ✅ | ✅ | ✅ | ✅ |
| TCP 通用 | ✅ (通过具体协议) | ✅ | ✅ | ✅ (Spectrum) |
| UDP | ❌ | ❌ (付费) | ✅ | ❌ |
| SSH（浏览器终端） | ✅ 内置 Web Terminal | ❌ 仅 TCP 转发 | ❌ 仅 TCP 转发 | ✅ (需配置 Access) |
| RDP（浏览器桌面） | ✅ 内置 Web Desktop | ❌ | ❌ | ❌ |
| VNC（浏览器桌面） | ✅ 内置 Web Desktop | ❌ | ❌ | ❌ |
| Telnet | ✅ | ❌ | ❌ | ❌ |
| SFTP 文件传输 | ✅ (SSH 模式下 `--enable-sftp`) | ❌ | ❌ | ❌ |

**关键区别**：ngrok 和 frp 做的是 **L4 层的端口转发**——它把远端的 3389 端口映射到公网，但用户依然需要自己启动 mstsc.exe（Windows 远程桌面客户端）去连接。Shield CLI 做的是 **L7 层的协议渲染**——远端服务通过 HTML5 在浏览器中直接呈现，零客户端安装。

来看实际命令对比。同样是暴露内网一台 Windows 的 RDP：

```bash
# ngrok：打通隧道，但对方需要装 RDP 客户端
ngrok tcp 3389
# 输出：tcp://0.tcp.ngrok.io:12345
# 对方需要：打开远程桌面连接 → 输入 0.tcp.ngrok.io:12345 → 登录

# frp：需要先部署 frps 服务端 + 写配置文件
# frpc.toml:
# [[proxies]]
# name = "rdp"
# type = "tcp"
# localIP = "127.0.0.1"
# localPort = 3389
# remotePort = 7001
frpc -c frpc.toml
# 对方需要：同 ngrok，必须有 RDP 客户端

# Shield CLI：一条命令，浏览器直接打开
shield rdp --username admin --auth-pass mypass
# 输出：https://xxxx-yishield.ac.example.com
# 对方需要：点开链接，完事
```

对于 SSH 也一样：

```bash
# Shield CLI
shield ssh 10.0.0.5 --username root
# 浏览器中直接出现 Web 终端，支持 SFTP 文件上传下载

# ngrok
ngrok tcp 22
# 对方需要：ssh -p 12345 root@0.tcp.ngrok.io
# 还需要处理 known_hosts、密钥等问题
```

---

## 二、配置复杂度：从"一条命令"到"一堆配置文件"

### Shield CLI 的智能默认值

Shield CLI 在 CLI 参数设计上做了大量的默认推断，能最大限度减少输入：

```bash
shield ssh                  # 等价于 127.0.0.1:22
shield ssh 2222             # 等价于 127.0.0.1:2222（纯数字 → 端口）
shield ssh 10.0.0.5         # 等价于 10.0.0.5:22（IP → 用默认端口）
shield ssh 10.0.0.5:2222    # 完整指定
shield rdp                  # 等价于 127.0.0.1:3389
shield vnc 10.0.0.10:5901   # 完整指定
shield http 3000            # 等价于 127.0.0.1:3000
```

这套规则写在 `cmd/helpers.go` 里，逻辑是：纯数字 → 端口（用 127.0.0.1），包含 `.` 或 `:` → IP 或 IP:Port，空 → 默认 IP + 默认端口。每种协议有自己的默认端口（SSH=22, RDP=3389, VNC=5900, HTTP=80, HTTPS=443, Telnet=23）。

### 与 frp 的配置对比

frp 是自建隧道的经典方案，但配置门槛较高。来看一个完整的 SSH 转发配置：

```toml
# frps.toml（服务端 —— 你需要一台公网机器）
bindPort = 7000

# frpc.toml（客户端）
serverAddr = "your-server.com"
serverPort = 7000

[[proxies]]
name = "ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 6000
```

这里涉及：部署服务端、配置端口映射、管理配置文件、维护公网服务器。Shield CLI 不需要你管服务端——公网网关由 Shield 的基础设施提供（类似 ngrok 的模式），但同时 CLI 侧开源（Apache 2.0）。

### 与 ngrok 的配置对比

ngrok 的单条命令确实简洁：

```bash
ngrok http 8080
```

但在多协议场景下（比如同时需要 SSH + RDP + 一个 HTTP 服务），ngrok 需要配置文件：

```yaml
# ngrok.yml
tunnels:
  ssh:
    proto: tcp
    addr: 22
  rdp:
    proto: tcp
    addr: 3389
  web:
    proto: http
    addr: 8080
```

Shield CLI 通过 Web UI（`shield start`）管理多个服务，最多保存 10 个应用配置，界面上点击 Connect/Disconnect 即可动态管理，不需要配置文件。

### 与 Cloudflare Tunnel 的配置对比

Cloudflare Tunnel 的配置最为复杂（但也最强大）：

```yaml
# config.yml
tunnel: your-tunnel-id
credentials-file: /root/.cloudflared/your-tunnel-id.json

ingress:
  - hostname: ssh.example.com
    service: ssh://localhost:22
  - hostname: rdp.example.com
    service: rdp://localhost:3389
  - service: http_status:404
```

你还需要：Cloudflare 账户 → 添加域名 → 创建 Tunnel → 配置 DNS → 配置 Access Policy。对于企业级持久化部署这无可厚非，但对于"临时给人开个远程桌面"来说，杀鸡用牛刀了。

---

## 三、安全模型对比

### 凭证存储

| 维度 | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|------|-----------|-------|-----|-------------------|
| 凭证存储方式 | AES-256-GCM 加密本地文件 | Token 明文存于 `~/.ngrok2/ngrok.yml` | Token 在配置文件中 | JSON 文件存于 `~/.cloudflared/` |
| 密钥来源 | 机器指纹 SHA256（主机名 + MAC + Machine ID） | 用户账户 Token | 用户自定义 | Cloudflare 颁发 |
| 跨机器迁移 | ❌ 加密文件绑定机器，迁移后自动失效 | ✅ Token 可迁移 | ✅ 配置文件可迁移 | ✅ 凭证文件可迁移 |

Shield CLI 的做法比较有意思：它用 **机器指纹** 作为 AES-256-GCM 的加密密钥。指纹由三部分组成：

1. 主机名（`os.Hostname()`）
2. 第一块物理网卡的 MAC 地址（跳过 docker/br-/veth/virbr 等虚拟接口）
3. 平台级 Machine ID（Linux: `/etc/machine-id`，macOS: `IOPlatformUUID`，Windows: 注册表 `MachineGuid`）

三者拼接后 SHA256 取摘要，作为 AES 密钥。这意味着：

- **凭证泄露风险降低**：即使有人拷走了 `~/.shield-cli/.credential` 文件，在另一台机器上也解不开（机器指纹不同）
- **不需要"登录"操作**：首次使用时自动生成凭证并注册到服务端，身份与机器绑定
- **trade-off**：换机器或重装系统后需要 `shield clean` 重置凭证

对比 ngrok 的 `ngrok config add-authtoken <token>` 方式——Token 是明文的，拷到另一台机器就能用，方便但风险更高。

### 访问控制

| 方式 | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|------|-----------|-------|-----|-------------------|
| 公开链接 | ✅ Visible 模式（默认） | ✅ 默认公开 | ✅ 默认公开 | ❌ 需要 Access 策略 |
| 授权访问 | 🔜 Invisible 模式（计划中） | ✅ IP 白名单/OAuth（付费） | ✅ 需自行实现 | ✅ Access（零信任） |
| 链接有效期 | 24 小时 API Key 自动刷新 | 免费版 2 小时/8 小时 | 无限（服务端运行期间） | 无限（Tunnel 运行期间） |

Shield CLI 目前默认是 Visible 模式——生成的 HTTPS 链接任何人都能访问。但服务端对每个 API Key 设置了 **24 小时有效期**，过期后会自动刷新。相比 ngrok 免费版的 2 小时限制（2024 年后调整过），Shield 的免费额度更宽松。

Cloudflare Tunnel 在安全性上是最强的——你可以配置完整的零信任策略（邮箱验证、SAML SSO、IP 限制等），但这也意味着更重的配置成本。

### 密码处理

Shield CLI 在日志中做了密码脱敏处理——只显示首尾各 2 个字符：

```
Connecting to 10.0.0.5:22 with password: my****ss
```

SSH 私钥通过 `--private-key` 参数传入文件路径，不会在命令行中暴露密钥内容。凭证在传输到服务端后，存储在 `main_app_config` 中用于协议网关的认证。

---

## 四、架构对比：Chisel vs ngrok 的私有协议

### Shield CLI 的双层隧道架构

Shield CLI 底层使用 [Chisel](https://github.com/jpillora/chisel)（一个基于 WebSocket 的 TCP 隧道库）。它建立了两条隧道：

```
┌──────────────┐                      ┌──────────────┐                    ┌──────────────┐
│ 内网服务      │ ←── 本地网络 ──→     │ Shield CLI   │ ←── WebSocket ──→  │ 公网网关      │
│ RDP/SSH/VNC  │                      │  (chisel     │     (wss://)      │ + 协议渲染    │
│ 10.0.0.5     │                      │   client)    │                    │ (Guacamole)  │
└──────────────┘                      └──────────────┘                    └──────────────┘
                                          │                                      │
                                     隧道 1: API 隧道                        隧道 2: 资源隧道
                                     (管理通道，持久化)                    (数据通道，按需创建)
```

**API 隧道**（Main Tunnel）：首次连接时建立，将本地 REST API 端口映射到公网，用于动态管理后续的资源隧道。这条隧道一直保持。

**资源隧道**（Resource Tunnel）：每个应用按需创建一条独立的 chisel 连接，映射目标服务端口到公网网关。最多 3 条并发。

这种设计的好处是：API 隧道提供了一个"控制面"，网关侧可以通过它动态地增删资源隧道，而不需要用户手动操作。

### ngrok 的架构

ngrok 使用自研的私有协议，客户端通过 TLS 连接到 ngrok 的边缘服务器，协议细节不公开。好处是性能优化可以做到极致，坏处是完全依赖 ngrok 的基础设施，无法审计。

### frp 的架构

frp 使用自定义的二进制协议（或可选的 KCP/QUIC），客户端和服务端都开源，你可以完全自建。但没有协议渲染层——它只做端口转发。

### Cloudflare Tunnel 的架构

cloudflared 通过 QUIC 协议连接到 Cloudflare 的全球边缘网络，自动利用 Anycast 选择最近的节点。在基础设施层面这是最强的方案（200+ 数据中心），但你的所有流量都经过 Cloudflare。

---

## 五、本地管理体验

### Shield CLI 的 Web UI

```bash
shield start
# 浏览器自动打开 http://localhost:8181
```

这会启动一个本地 Web 管理界面，功能包括：

- **应用管理**：添加/编辑/删除应用配置（协议、目标 IP:Port、凭证等）
- **一键连接**：点击 Connect 按钮，后台自动建立隧道，成功后弹出访问链接
- **状态监控**：实时显示每个应用的连接状态（idle / connecting / connected / failed）
- **配置持久化**：最多 10 个应用配置，AES-256-GCM 加密存储
- **深色/浅色主题**：支持切换

前端是纯 HTML5 + 原生 JS（约 1500 行），嵌入在二进制中，无外部依赖。后端是标准的 REST API。

### ngrok 的管理

ngrok 免费版没有本地 UI。你可以通过 `http://localhost:4040` 查看请求日志（仅 HTTP 隧道），但无法管理多个隧道。全功能管理在 ngrok Dashboard（SaaS），需要注册账户。

### frp 的管理

frp 有一个可选的 Dashboard（`frps` 启动时开启），可以查看代理列表和流量统计。但它在服务端运行，不是客户端本地的管理界面。且 UI 比较基础。

### Cloudflare Tunnel 的管理

通过 Cloudflare Zero Trust Dashboard（SaaS）管理。功能最全面（流量分析、访问策略、审计日志），但依赖云端。

---

## 六、部署与分发

| 维度 | Shield CLI | ngrok | frp | Cloudflare Tunnel |
|------|-----------|-------|-----|-------------------|
| 安装方式 | Homebrew / Scoop / curl / dpkg / rpm / 源码编译 | Homebrew / apt / choco / snap / 官网下载 | GitHub Release 下载 / 源码编译 | Homebrew / apt / 官网下载 |
| 二进制大小 | ~15 MB | ~25 MB | ~12 MB (frpc) | ~35 MB |
| 平台支持 | Linux/macOS/Windows (amd64/arm64/386) | Linux/macOS/Windows/FreeBSD | Linux/macOS/Windows/FreeBSD + 更多 | Linux/macOS/Windows |
| 中国镜像 | ✅ jsDelivr CDN 镜像 | ❌ 需要翻墙下载 | ✅ 国内可直连 GitHub | ❌ 需要翻墙 |
| 开源协议 | Apache 2.0 (CLI 侧) | 私有 | Apache 2.0 (全部) | 私有 |
| 自建服务端 | 🔜 计划中 | ❌ | ✅ 完全支持 | ❌ |

在中国大陆的网络环境下，Shield CLI 提供了 jsDelivr CDN 镜像安装：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

ngrok 和 Cloudflare Tunnel 的下载在国内经常受阻，这是一个实际的痛点。

---

## 七、价格与限制

| 维度 | Shield CLI | ngrok (Free) | ngrok (Personal $8/mo) | frp | Cloudflare Tunnel |
|------|-----------|-------------|----------------------|-----|-------------------|
| 费用 | 免费 | 免费 | $8/月 | 免费（自建需服务器） | 免费（需域名在 CF） |
| 隧道数量 | 3 并发 | 1 agent / 1 domain | 2 agents / 1 domain | 无限制 | 无限制 |
| 保存配置数 | 10 个 | N/A | N/A | 配置文件不限 | 不限 |
| 带宽限制 | 未声明 | 1 GB/月 | 1 GB/月 | 取决于服务器 | 未声明 |
| 连接时长 | 24 小时（自动续期） | 2 小时（需重连） | 无限 | 无限 | 无限 |
| TCP 隧道 | ✅ 免费 | ❌ 需付费 | ✅ | ✅ | ✅ (Spectrum 付费) |
| 自定义域名 | 🔜 计划中 | ❌ 需付费 | ✅ | ✅ | ✅ |
| 访问日志 | 本地日志 | Dashboard | Dashboard | Dashboard | Dashboard + 分析 |

一个经常被忽略的细节：**ngrok 免费版不支持 TCP 隧道**。也就是说，你无法免费用 ngrok 转发 SSH（22端口）或 RDP（3389端口）。Shield CLI 的 TCP 类协议（SSH/RDP/VNC/Telnet）全部免费可用。

---

## 八、实际使用体验对比

### 场景 1：给客户演示内网系统

```bash
# Shield CLI：一条命令，发链接
shield http 3000
# → https://abc123-yishield.ac.example.com
# 客户点击链接直接看到你的应用

# ngrok：类似，但免费版链接 2 小时过期
ngrok http 3000
# → https://abc123.ngrok-free.app
# 客户访问时会看到 ngrok 的警告页（免费版）
```

ngrok 免费版有一个 interstitial 警告页面（"You are about to visit..."），在给客户演示时观感不佳。Shield CLI 没有这个限制。

### 场景 2：远程协助 Windows 桌面

```bash
# Shield CLI：浏览器中直接出现 Windows 桌面
shield rdp 10.0.0.100 --username admin --auth-pass P@ssw0rd
# → https://xxx-yishield.ac.example.com
# 对方在浏览器里操作完整的 Windows 桌面

# 其他工具：都无法做到浏览器内渲染 RDP
# ngrok: ngrok tcp 3389 → 对方需要装 RDP 客户端
# 如果对方是 Mac/Linux 用户，还需要额外安装 Microsoft Remote Desktop 或 Remmina
```

这个场景是 Shield CLI 的核心优势所在。其他隧道工具在这里只能提供"网络可达"，用户仍需自行解决"客户端兼容性"的问题。

### 场景 3：管理多个内网服务

```bash
# Shield CLI：启动 Web UI 统一管理
shield start
# 在浏览器中添加多个应用，点击连接/断开

# frp：需要编辑配置文件，重启客户端
# ngrok：需要编写 ngrok.yml，或开多个终端窗口
```

---

## 九、局限性与 trade-off

公平起见，Shield CLI 也有其当前阶段的不足：

1. **服务端不开源**：网关服务由 Shield 官方运营（console.yishield.com），目前不能自建。这意味着数据会经过第三方服务器。路线图中有自建部署计划，但尚未发布。frp 在这一点上完胜。

2. **并发限制**：最多 3 条并发连接、10 个保存配置。对于个人和小团队够用，但企业级场景不足。

3. **无 UDP 支持**：底层 Chisel 基于 WebSocket（TCP），不支持 UDP 协议。frp 在这方面更全面。

4. **访问控制较弱**：当前只有"有链接就能访问"的 Visible 模式。Invisible 模式（需要额外授权密钥）在规划中但未上线。Cloudflare Access 的零信任方案在安全性上领先一个量级。

5. **社区生态**：作为新项目，社区规模远不及 ngrok（GitHub 25k+ stars）和 frp（80k+ stars）。遇到问题可能需要直接看源码。

6. **依赖外部服务**：虽然 CLI 开源，但核心功能依赖 Shield 的公网网关。如果服务不可用，工具就无法使用。这与 frp 的完全自主可控形成对比。

---

## 十、选型建议

| 你的需求 | 推荐方案 | 理由 |
|---------|---------|------|
| 临时暴露本地 Web 服务给同事 | ngrok 或 Shield CLI | 都是一条命令搞定 |
| 远程桌面（RDP/VNC）通过浏览器访问 | **Shield CLI** | 唯一支持浏览器内渲染桌面协议的方案 |
| 完全自建、不经过任何第三方 | **frp** | 全开源，服务端自部署 |
| 企业级零信任远程访问 | **Cloudflare Tunnel + Access** | 最完善的安全策略引擎 |
| 大陆网络环境下使用 | **Shield CLI 或 frp** | 有国内可达的安装和服务节点 |
| SSH + SFTP 文件传输一站式 | **Shield CLI** | 浏览器内 SSH + SFTP 开箱即用 |
| 需要 UDP 转发（游戏、DNS） | **frp** | 唯一支持 UDP 的方案 |
| 预算敏感、需要 TCP 隧道 | **Shield CLI** | ngrok TCP 隧道需要付费 |

---

## 总结

Shield CLI 不是要替代 ngrok 或 frp——它们解决的问题重叠但不相同。**如果你的核心诉求是"让别人通过浏览器直接操作内网的桌面或终端"，Shield CLI 是目前唯一用一条命令就能做到的工具**。它把隧道工具和协议网关整合成了一个工作流，省去了中间"装客户端"的步骤。

但如果你需要完全自主可控的基础设施（frp）、企业级的零信任安全策略（Cloudflare），或者只是简单地转发一个 HTTP 服务（ngrok），那些工具各自有不可替代的优势。

技术选型从来都是 trade-off，希望这篇对比能帮你在具体场景下做出更合适的选择。

---

*Shield CLI 开源地址：https://github.com/fengyily/shield-cli*
*协议：Apache 2.0*
