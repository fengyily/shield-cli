---
title: 推荐 Shield CLI — 一个思路不同的内网穿透工具
description: Shield CLI 不只是打隧道，而是在网关侧做协议渲染。SSH、RDP、VNC 通过 HTML5 在浏览器中直接呈现，对方不需要装任何客户端。
date: 2025-01-20
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI 推荐, 内网穿透工具, SSH Web终端, RDP浏览器, 远程桌面, 零客户端
---

# 推荐一个思路不太一样的内网穿透工具：Shield CLI

> 大多数隧道工具解决的是"网络可达"，但如果对方没有客户端，隧道打通了也白搭。Shield CLI 的做法是：隧道 + 浏览器内协议渲染，一条命令生成链接，对方打开浏览器就能操作 RDP 桌面、SSH 终端、VNC 会话。

---

## 先说场景

大多数内网穿透工具（ngrok、frp、Cloudflare Tunnel）解决的是 **"网络可达"** 的问题——把内网端口映射到公网。但实际工作中经常遇到这种情况：

- 要给客户远程看一台内网 Windows 的桌面，对方是 Mac 用户，还得让人家装 Microsoft Remote Desktop
- 给外包同事临时开个 SSH 权限，对方说"我电脑上没有终端工具"
- VNC 桌面想分享给同事看，结果光是教对方配 VNC Viewer 就花了 20 分钟

**隧道打通了，但对方没有客户端，等于白搭。**

---

## Shield CLI 的思路不一样

它不只是打隧道，而是在网关侧直接做了 **协议渲染**——SSH、RDP、VNC 都通过 HTML5 在浏览器里直接呈现。

```bash
# 一条命令暴露内网 RDP 桌面
shield rdp 10.0.0.100 --username admin --auth-pass P@ssw0rd
# 输出一个 HTTPS 链接，对方浏览器打开就是完整的 Windows 桌面

# SSH 终端也一样
shield ssh 10.0.0.5
# 浏览器里直接出现终端，还支持 SFTP 文件传输

# 本地 Web 服务
shield http 3000
```

对方 **不需要装任何客户端**，手机、平板、锁死的公司电脑，有浏览器就行。

---

## 和其他工具的关键区别

| | Shield CLI | ngrok | frp |
|---|---|---|---|
| RDP/VNC 浏览器内渲染 | ✅ | ❌ 需客户端 | ❌ 需客户端 |
| SSH Web 终端 | ✅ 内置 | ❌ 仅 TCP 转发 | ❌ 仅 TCP 转发 |
| TCP 隧道 | 免费 | 需付费 | 免费（需自建服务端） |
| 配置复杂度 | 一条命令 | 一条命令（HTTP）/配置文件（多隧道） | 需部署服务端 + 写配置 |
| 国内可用性 | ✅ jsDelivr 镜像 | 下载经常受阻 | ✅ |

---

## 几个实用的点

### 智能默认值

`shield ssh` 自动解析为 `127.0.0.1:22`，`shield ssh 2222` 解析为 `127.0.0.1:2222`，`shield ssh 10.0.0.5` 解析为 `10.0.0.5:22`，基本不用记参数。

### Web UI 管理

`shield start` 启动本地管理面板（localhost:8181），可以保存最多 10 个应用配置，点击连接/断开，适合日常管理多个内网服务。

### 安全方面

凭证用 AES-256-GCM 加密存储，密钥绑定机器指纹（主机名 + MAC + Machine ID），文件拷到别的机器解不开。日志里密码自动脱敏。

---

## 安装

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli

# Linux / macOS 一键安装
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

---

## 客观说下不足

1. **服务端目前不开源**，数据经过官方网关，介意的话 frp 全栈自建更合适
2. **并发上限 3 条隧道**，个人和小团队够用，大规模不行
3. **不支持 UDP**，有 UDP 需求选 frp
4. **访问控制还比较基础**，目前是"有链接就能访问"，零信任方案还在规划中

---

## 总结

如果你的需求是 **"让别人通过浏览器直接操作内网桌面或终端"**，Shield CLI 是目前唯一一条命令就能搞定的方案。它把隧道 + 协议网关合成了一个工作流，省掉了"让对方装客户端"这个最大的摩擦点。

如果只是转发 HTTP 服务，ngrok 依然好用；如果要完全自主可控，frp 无可替代。工具不冲突，看场景选。

开源地址：https://github.com/fengyily/shield-cli （Apache 2.0）
