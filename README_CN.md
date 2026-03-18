<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>基于浏览器的安全内网穿透工具</strong><br>
  通过浏览器直接访问内网的 RDP 远程桌面、VNC 屏幕、SSH 终端和 Web 服务 — 无需安装任何客户端软件。
</p>

<p align="center">
  <a href="#安装">安装</a> &bull;
  <a href="#快速开始">快速开始</a> &bull;
  <a href="#可见与隐身模式">可见与隐身模式</a> &bull;
  <a href="#命令参数">命令参数</a> &bull;
  <a href="../README.md">English</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

---

## 为什么选择 Shield CLI？

传统远程访问需要在每台设备上安装专用客户端（RDP 客户端、VNC Viewer、SSH 终端）。Shield CLI 通过安全网关将内网服务隧道化，所有操作都在**浏览器中完成**。

- **无需安装客户端** — 在任意浏览器中打开链接即可访问 RDP 桌面、VNC 屏幕或 SSH 终端
- **随处可用** — 手机、平板、受限的企业电脑，任何有浏览器的设备都能访问
- **一个命令** — 一条 `shield` 命令即可将任意内网服务暴露到公网

## 安装

### macOS (Homebrew)

```bash
brew tap fengyily/tap
brew install shield-cli
```

### Windows (Scoop)

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### Windows (PowerShell 一键安装)

```powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```

### Linux / macOS (curl 一键安装)

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### Debian / Ubuntu (.deb)

```bash
# 从 GitHub Releases 下载
sudo dpkg -i shield-cli_<version>_linux_amd64.deb
```

### RHEL / CentOS (.rpm)

```bash
sudo rpm -i shield-cli_<version>_linux_amd64.rpm
```

### 从源码编译

```bash
git clone https://github.com/fengyily/shield-cli.git
cd shield-cli
go build -o shield .
```

## 快速开始

```bash
# 暴露本地 SSH — 在浏览器中访问终端
shield ssh

# 暴露远程 RDP 桌面
shield rdp 10.0.0.5

# 暴露本地 3000 端口的开发服务
shield http 3000

# 暴露远程 VNC 服务（自定义端口）
shield vnc 10.0.0.10:5901
```

隧道建立后，在任意浏览器中打开 **Access URL** 即可访问。

### 智能默认值

| 命令 | 解析为 |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield ssh 10.0.0.2:2222` | `10.0.0.2:2222` |
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield vnc` | `127.0.0.1:5900` |
| `shield https` | `127.0.0.1:443` |
| `shield telnet` | `127.0.0.1:23` |

支持的协议：`ssh`、`http`、`https`、`rdp`、`vnc`、`telnet`

## 可见与隐身模式

隧道建立后会生成两个 URL：

- **Site URL** — 应用地址（如 `https://xxxx.hk01.apps.yishield.com`）。单独使用此 URL **无法访问**，需要授权。
- **Access URL** — 包含内嵌授权密钥的链接。拥有此 URL 的人可以直接访问服务。

### 隐身模式（默认）

服务**需要授权** — 仅凭 Site URL 无法访问，必须使用 Access URL。

```bash
shield rdp 10.0.0.5
shield ssh 10.0.0.2
```

**适用场景：** 生产服务器、故障处理期间的临时访问、敏感机器。

### 可见模式

服务**对未授权用户开放** — 任何知道 Site URL 的人都可以直接访问。

```bash
shield --visable ssh 10.0.0.2
shield --visable=HK rdp 10.0.0.5
```

**适用场景：** 公开演示环境、团队共享开发服务器、QA 预发布环境。

## 命令参数

```
shield <protocol> [ip:port] [flags]

参数:
  -H, --server string         API 服务器地址 (默认: https://console.yishield.com/raas)
  -p, --tunnel-port int       隧道服务器端口 (默认: 62888)
      --visable [过滤词]      启用可见模式 (可选: AC 节点名称过滤)
      --invisible             隐身模式，需要授权密钥
      --display-name string   连接器显示名称
      --site-name string      应用站点名称
      --username string       目标服务用户名 (SSH/RDP/VNC)
      --auth-pass string      目标服务密码 (SSH/RDP/VNC)
      --private-key string    SSH 私钥
      --passphrase string     SSH 私钥密码
      --enable-sftp           启用 SFTP (仅 SSH)
  -v, --verbose               启用详细日志输出
  -h, --help                  显示帮助信息

子命令:
  clean                       清除缓存的凭证
```

### 运行截图

![Shield CLI SSH](docs/images/shieldcli-ssh-001.jpg)

### 本地 API

Shield CLI 运行后会在 `127.0.0.1:<port>` 上提供本地管理接口：

| 接口 | 方法 | 说明 |
|---|---|---|
| `/health` | GET | 健康检查 |
| `/connectors` | GET | 列出所有活跃隧道 |
| `/connector?rport=&lip=&lport=` | GET | 创建动态隧道 |
| `/connector?rport=` | DELETE | 关闭隧道 |

## 工作原理

```
                      浏览器 (通过 HTML5 访问 RDP/VNC/SSH)
                                │
                                ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│  内网服务      │ ◄──► │  Shield CLI   │ ◄══► │  公网网关      │
│  10.0.0.5:    │ 本地  │  (隧道)       │chisel│  + Web UI    │
│  3389/5900/22 │      └──────────────┘wss://└──────────────┘
└──────────────┘
```

## 安全性

- 凭证使用 AES-256-GCM 加密，密钥基于机器指纹派生
- 所有日志输出中的密码均已脱敏
- 隧道连接使用带认证的 WebSocket 传输
- 凭证文件权限为 `0600`

## 许可证

MIT
