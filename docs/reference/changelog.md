---
title: 更新日志 — Shield CLI 版本历史
description: Shield CLI 完整版本发布记录。追踪每个版本的新功能、改进和修复。
head:
  - - meta
    - name: keywords
      content: Shield CLI 更新日志, 版本历史, 发布记录, 新功能, 变更记录
---

# 更新日志

Shield CLI 所有版本的变更记录。

## v0.3.x <Badge type="tip" text="最新" />

### v0.3.2 — PostgreSQL 插件 ER 图 {#v0.3.2}

**发布日期：2026-03-26**

#### 新功能

- **PostgreSQL 插件 v0.3.0** — 交互式 ER 图及表/字段/外键可视化管理
  - ER 图：SVG 渲染，支持缩放、平移、四种布局模式（Grid / Horizontal / Vertical / Center）
  - 拖拽创建外键：从字段拖动到目标表，自动匹配类型和字段名
  - 点选删除外键：点击关系线选中，按 Delete 键删除
  - 右键菜单：表头右键 → 重命名/删除表、新增字段；空白区右键 → 新建表
  - 表结构编辑器：点击表头齿轮图标，批量修改表名、字段名、字段类型
  - 字段操作：悬停字段行显示齿轮图标 → 编辑/删除字段
  - 动态表宽度：根据字段名和类型长度自动调整
  - 所有操作均有 SQL 预览弹窗，确认后执行
  - Read-only 模式下多层拦截（前端隐藏 + 后端 403）
  - SVG 导出、localStorage 持久化位置和布局
  - 侧边栏 Schema 节点新增 ER 图快捷入口

#### 文档

- 新增 [PostgreSQL 插件文档](/plugins/postgres)
- 插件系统概览页 PostgreSQL 从"即将发布"更新为正式可用

### v0.3.1 — MySQL 插件只读模式 {#v0.3.1}

**发布日期：2026-03-24**

#### 新功能

- **只读模式** — MySQL 插件支持服务端强制只读，远程用户无法绕过
  - CLI 新增 `--readonly` 参数强制只读模式
  - Web UI 添加/编辑应用时可勾选 **Read-Only Mode** 复选框
  - 只读状态全链路传递：CLI/Web UI → PluginConfig → 插件后端
  - 页面右上角标识仅展示当前模式，远程用户无法切换
  - 前后端双重拦截写操作

#### 文档

- 博文重写：[插件系统 & MySQL Web 管理](/blogs/blog-shield-cli-plugin-mysql)（以 Web UI 操作为主线）
- MySQL 插件文档更新只读模式说明（中英文）
- 命令参考新增 `--readonly` 参数

### v0.3.0 — 插件系统 & MySQL Web 客户端 {#v0.3.0}

**发布日期：2026-03-24**

#### 新功能

- **插件系统** — 通过独立二进制插件按需扩展协议支持，主程序零膨胀
  - `shield plugin add <name>` — 安装插件（支持从 GitHub Releases 下载或本地安装）
  - `shield plugin list` — 查看已安装插件
  - `shield plugin remove <name>` — 卸载插件
  - 插件通过 stdin/stdout JSON 协议与主程序通信，架构简洁可靠
- **MySQL 插件** — 浏览器内 Web 数据库管理客户端
  - `shield mysql 127.0.0.1:3306` — 一行命令暴露 MySQL Web 管理界面
  - 支持 `mysql` 和 `mariadb` 协议别名
  - Web 界面功能：数据库浏览、表过滤与翻页、表结构查看、SQL 执行
  - 结果排序（点击列头）、CSV 导出、一键复制
  - 默认只读模式，阻止写操作（可切换为读写模式）
  - 交互式凭证输入（缺少 `--db-user` / `--db-pass` 时自动提示）
  - 兼容 `--username` / `--auth-pass` 通用认证参数
- **数据库连接参数** — 新增 `--db-user`、`--db-pass`、`--db-name` 参数

#### 文档

- 新增[插件系统概览](/plugins/)文档（中英文）
- 新增 [MySQL 插件详细文档](/plugins/mysql)（中英文）
- 新增[插件开发指南](/plugins/development)（中英文）
- 命令参考新增插件管理和数据库参数说明

---

## v0.2.x

### v0.2.6 — TCP/UDP 端口代理 {#v0.2.6}

**发布日期：2026-03-24**

#### 新功能

- **TCP/UDP 端口代理** — 通过加密隧道转发任意 TCP/UDP 端口
  - `shield tcp 3306` — 代理 MySQL、Redis、PostgreSQL 等 TCP 服务
  - `shield udp 53` — 代理 DNS、Syslog 等 UDP 服务
  - 必须指定端口（tcp/udp 无默认端口）
  - 隧道建立后显示连接指南（专属域名 + 端口），不会自动打开浏览器
- **UDP over chisel** — UDP 隧道使用 chisel 原生 `/udp` 后缀进行 UDP 转发

#### 改进

- **隧道激活修复** — CLI 改为调用 `POST _webgate/api/tunnel`（与 Web UI 一致），而非简单 GET 请求
- **verbose 调试输出** — `-v` 参数现在在静默阶段也输出调试日志（quick-setup API、隧道创建、激活过程）

### v0.2.5 — Linux 包管理器支持 {#v0.2.5}

**发布日期：2026-03-23**

#### 新功能

- **APT / YUM 仓库** — 支持通过 `apt install shield-cli` 和 `yum install shield-cli` 安装
  - 基于 GitHub Pages 托管的 APT（Debian/Ubuntu）和 YUM（RHEL/CentOS/Fedora）包仓库
  - 添加仓库源后，`apt upgrade` / `yum update` 自动获取新版本
- **一键仓库配置脚本** — `setup-repo.sh` 自动检测包管理器并配置仓库源
- **install.sh 增强** — 支持 `--apt`、`--yum`、`--dnf` 参数快速配置包管理器安装

#### 文档

- README 和安装指南新增 APT / YUM 安装说明（中英文）

### v0.2.2 — Docker 支持 {#v0.2.2}

**发布日期：2026-03-22**

#### 新功能

- **Docker 支持** — 新增 Dockerfile，支持容器化部署
  - 多阶段构建，基于 Alpine 的轻量镜像
  - 多架构支持（linux/amd64、linux/arm64）
  - `--network host` 模式下可访问宿主机及内网资源
- **Docker 镜像自动发布** — CI 流水线新增 Docker 构建任务
  - 同时推送至 Docker Hub（`fengyily/shield-cli`）和 GHCR（`ghcr.io/fengyily/shield-cli`）
  - 自动生成语义化版本标签（`latest`、`0.2.2`、`0.2`）
- **监听地址可配置** — 新增 `SHIELD_LISTEN_HOST` 环境变量，支持自定义 Web UI 监听地址（默认 `127.0.0.1`，容器内自动设为 `0.0.0.0`）

#### 文档

- README 和安装指南新增 Docker 部署说明（中英文）
- 新增博客：[用 Docker 跑 Shield CLI](../blogs/blog-shield-cli-docker.md)

---

## v0.2.1

**发布日期：2026-03-20**

### 新功能

- **系统服务安装** — `shield install` 将 Shield 注册为系统服务，开机自动启动
  - macOS：launchd 用户代理（无需 sudo）
  - Linux：systemd 服务
  - Windows：Windows 服务
- **自定义端口** — `shield install --port 8182`，自动检测端口冲突并建议可用端口
- **系统托盘图标**（macOS 和 Windows）— 点击打开 Dashboard，支持重启和退出操作
- **异步隧道启动** — Web UI 即时可用，主隧道在后台连接
- **隧道状态 API** — `GET /api/tunnel` 接口供前端轮询隧道就绪状态

### 改进

- GoReleaser 拆分为桌面版（CGO + 托盘）和 Linux 版（纯 Go）构建
- 隧道连接中时，应用连接请求返回明确提示信息

---

## v0.2.0

**发布日期：2026-03-19**

### 新功能

- **Web UI 管理平台** — 浏览器端管理面板 `localhost:8181`
  - 添加、编辑、删除最多 10 个应用配置
  - 一键连接/断开，实时状态显示
  - AES-256-GCM 加密本地存储应用配置
- **持久化配置** — 应用配置文件加密保存到本地
- **多连接支持** — 最多 3 个并发活跃隧道连接
- **连接管理器** — 共享主隧道 + 每个应用独立的动态资源隧道

### 改进

- 重新设计 Logo 和品牌形象
- README 增加 Web UI 截图和示例

---

## v0.1.3

**发布日期：2026-03-18**

### 新功能

- **Windows 安装脚本** — PowerShell 一键安装
- **Linux 安装脚本** — curl 一键安装
- **双语 README** — 拆分为英文（`README.md`）和中文（`README_CN.md`）

### 改进

- 默认使用可见访问模式

---

## v0.1.2

**发布日期：2026-03-18**

### 新功能

- **Scoop 包管理** — Windows 上 `scoop install shield-cli`
- **deb / rpm 包** — Linux 原生包格式（通过 nfpm）
- **curl 安装器** — `curl -fsSL ... | sh` 一键安装
- **国内 CDN 镜像** — 基于 jsDelivr 的安装脚本，国内用户友好

---

## v0.1.1

**发布日期：2026-03-18**

### 改进

- **位置参数** — `shield ssh 10.0.0.5:2222` 代替 `--type ssh --source 10.0.0.5:2222`
- **智能默认值** — 省略 IP 默认 localhost，省略端口使用协议默认值
- 简化命令行用法，直觉化地址解析

---

## v0.1.0

**发布日期：2026-03-18**

### 新功能

- **GoReleaser 集成** — 自动化跨平台构建（macOS、Linux、Windows × amd64、arm64）
- **Homebrew tap** — `brew install shield-cli`
- **自动发布** — GitHub Actions CI/CD 流水线

---

## v0.0.1

**发布日期：2026-03-16**

### 首次发布

- 基于 Chisel 协议的核心隧道连接
- 支持协议：SSH、RDP、VNC、HTTP、HTTPS、Telnet
- AES-256-GCM 加密凭证存储，绑定机器指纹
- 可见模式与隐身模式
- 连接成功后自动打开浏览器
- 日志中密码自动脱敏
- GitHub Actions CI/CD 流水线
