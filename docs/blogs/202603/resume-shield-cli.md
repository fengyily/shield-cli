# Shield CLI — 简历项目描述

> 根据实际需要选择合适的版本长度。

---

## 一句话版（适合项目列表）

开源远程访问工具，通过加密隧道将内网 RDP/SSH/VNC 服务直接渲染到浏览器，无需 VPN 和客户端安装。

---

## 简要版（适合简历项目栏，3-5 行）

**Shield CLI** — 开源安全隧道工具 | Go | Apache 2.0

- 基于 Go 开发的 CLI 工具，通过 WebSocket 加密隧道 + HTML5 协议渲染，实现浏览器直接访问内网 RDP 桌面、SSH 终端、VNC 会话等 6 种协议，无需 VPN 或客户端软件
- 设计并实现双层隧道架构（API 控制面 + 资源数据面），集成 Web UI 管理界面、系统托盘、系统服务安装等功能
- 搭建 GoReleaser + GitHub Actions CI/CD 流水线，支持跨平台编译（macOS/Linux/Windows）及 Homebrew/Scoop/apt/yum/Docker 多渠道分发
- 项目独立完成，3,700+ 行 Go 代码，已发布 v0.2.5 版本

---

## 详细版（适合项目经验详述）

### Shield CLI — 浏览器直连内网服务的开源安全隧道工具

**技术栈**: Go 1.25, Chisel (WebSocket Tunnel), Cobra CLI, HTML5, Docker, GitHub Actions, GoReleaser

**项目地址**: https://github.com/fengyily/shield-cli

**项目简介**:

独立设计并开发的开源命令行工具，解决传统隧道工具（ngrok/frp）只能做端口转发、用户仍需安装协议客户端的痛点。通过 HTML5 协议网关在浏览器中直接渲染 RDP 桌面、SSH 终端和 VNC 会话，实现"一条命令、一个 URL、浏览器即可访问"。

**核心工作**:

- **架构设计**: 设计双层隧道架构 — API Tunnel（持久化控制面）+ Resource Tunnel（按需数据面），实现动态隧道管理，支持 6 种协议（SSH/RDP/VNC/HTTP/HTTPS/Telnet）
- **安全体系**: 实现基于机器指纹（Hostname + MAC + Machine ID）的 AES-256-GCM 凭据加密方案，日志密码脱敏，文件权限 0600 控制
- **CLI 设计**: 基于 Cobra 构建智能参数推断系统（纯数字→端口、IP→自动补全默认端口），降低用户使用门槛
- **Web UI**: 开发内嵌式管理界面（纯 HTML5 + vanilla JS），支持应用配置管理、一键连接、状态监控、暗色主题
- **跨平台适配**: 实现 macOS（launchd）、Linux（systemd）、Windows 三平台系统服务安装与管理，macOS/Windows 系统托盘集成
- **CI/CD & 分发**: 搭建 GoReleaser + GitHub Actions 自动化流水线，支持跨平台编译、Docker 多阶段构建，接入 Homebrew/Scoop/apt/yum 四大包管理器及 jsDelivr CDN 国内镜像
- **文档 & 推广**: 基于 VitePress 构建中英双语文档站，编写技术对比文章（Shield CLI vs ngrok vs frp vs Cloudflare Tunnel）

**项目成果**:

- 3,700+ 行 Go 代码，31+ commits，已发布至 v0.2.5
- 支持 macOS / Linux / Windows 三平台 + Docker 容器化部署
- 6 种包管理器分发渠道（Homebrew / Scoop / apt / yum / Docker / curl）
- 中英双语文档站 + 技术博客

---

## 英文版（English Resume）

### Shield CLI — Open-Source Secure Tunnel Tool for Browser-Based Remote Access

**Tech Stack**: Go, WebSocket (Chisel), Cobra CLI, HTML5, Docker, GitHub Actions, GoReleaser

- Independently designed and built an open-source CLI tool that creates encrypted tunnels and renders internal services (RDP/SSH/VNC) directly in the browser via HTML5 — eliminating the need for VPN or client software
- Architected a dual-layer tunnel system (API control plane + on-demand resource data plane) supporting 6 protocols with smart parameter inference
- Implemented AES-256-GCM credential encryption using machine fingerprint-derived keys (Hostname + MAC + Machine ID)
- Built embedded Web UI dashboard for tunnel management with real-time status monitoring
- Delivered cross-platform support (macOS/Linux/Windows) with system service integration (launchd/systemd) and system tray
- Set up GoReleaser + GitHub Actions CI/CD pipeline with multi-platform builds and 6 distribution channels (Homebrew, Scoop, apt, yum, Docker, curl)
- 3,700+ lines of Go, 31+ commits, released v0.2.5

---

## 关键技术亮点（面试 talking points）

| 方面 | 亮点 |
|------|------|
| **架构** | 双层隧道设计（控制面/数据面分离），参考微服务 control plane 思想 |
| **安全** | 机器指纹 + AES-256-GCM 加密方案，凭据无法跨机器迁移 |
| **CLI 设计** | 智能参数推断，最少输入原则（`shield ssh` 等价于完整 IP:Port） |
| **跨平台** | 三平台系统服务 + 托盘，Go 条件编译处理平台差异 |
| **CI/CD** | GoReleaser 一键发布到 6 个分发渠道 |
| **差异化** | 业内唯一将隧道 + 协议渲染整合为单一命令行工具的开源方案 |
