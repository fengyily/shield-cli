# Shield CLI 视频教程脚本

> **视频标题**：Shield CLI — 一条命令，浏览器访问一切内部服务
> **时长**：约 5 分钟
> **风格**：纯屏幕录制 + AI 配音 + 字幕，无真人出镜
> **语言**：中文配音（可后续出英文版）
> **分辨率**：1920×1080

---

## 录制工具链

```
录制层: record-demo.sh --video          ← 屏幕内容 (ffmpeg)
叠加层: Cursor Pro (免费 macOS app)     ← 鼠标光圈 + 点击涟漪
自动层: capture-arch-gif.mjs --video    ← 架构图动画 (Puppeteer)
合成:   剪映 关键帧缩放                  ← 聚焦放大效果
```

### 录制前准备

1. 安装 **Cursor Pro** (Mac App Store 免费) → 开启 Highlight + Click Effect
2. 设置终端字体 **16pt**，窗口布局：左终端 + 右浏览器
3. 启动 VitePress: `npx vitepress dev docs --port 5173`
4. 确保内网演示环境就绪 (SSH server / Windows RDP / MySQL)

### 后期关键帧缩放技巧（剪映）

在以下时间点添加 scale 关键帧 (100%→140%→100%, 缓入缓出):
- 每次终端输入命令时 → 放大终端区域
- 点击 Web UI 按钮时 → 放大按钮区域
- 浏览器打开远程画面时 → 放大浏览器区域
- 每个聚焦持续 2-3s，过渡 0.5s

## 素材清单

| 编号 | 内容 | 录制命令 |
|------|------|---------|
| S1 | 架构图动画 | `node scripts/capture-arch-gif.mjs --video` |
| S2 | 安装过程 | `./docs/demo/record-demo.sh install --video` |
| S3 | Web UI 添加应用 + 连接 | `./docs/demo/record-demo.sh webui --video` |
| S4 | `shield ssh` 命令行 | `./docs/demo/record-demo.sh ssh-cli --video` |
| S5 | SSH 浏览器终端操作 | (S4 连续录制，不需要单独录) |
| S6 | `shield rdp` 命令 | `./docs/demo/record-demo.sh rdp --video` |
| S7 | RDP 浏览器桌面操作 | (S6 连续录制) |
| S8 | `shield mysql` 命令 | `./docs/demo/record-demo.sh mysql-cli --video` |
| S9 | MySQL Web 管理界面 | (S8 连续录制) |
| S10 | 手机浏览器打开 URL | `./docs/demo/record-demo.sh mobile --video` + Xcode Simulator |
| S11 | GitHub + 文档链接 | 截图即可 |

> **提示**: S4+S5、S6+S7、S8+S9 各自一次连续录制，终端执行命令后切到浏览器操作，后期剪辑分段。

---

## 分镜脚本

### 第一幕：痛点引入（0:00 - 0:30）

---

**画面 [0:00 - 0:10]**
黑屏渐入，逐行打字效果显示文字：

```
你是否遇到过这些场景？
  → 想远程访问公司内网的服务器，VPN 配置太复杂
  → 让同事临时看一下你的开发环境，得装一堆客户端
  → 出差时用手机操作一台 Windows 桌面
```

**配音**：
> 远程访问公司内网服务器，VPN 配置太复杂。让同事临时访问你的开发环境，对方还得装客户端。出差时想用手机操作一台 Windows，传统方案不是配置复杂，就是需要对方安装软件。

---

**画面 [0:10 - 0:20]**
终端画面，打字输入：

```bash
shield ssh 10.0.0.5
```

回车后显示公网 URL，浏览器自动打开 SSH 终端。

**配音**：
> 如果有一个工具，只需要一条命令，就能生成一个链接，对方在浏览器里直接操作——不需要 VPN，不需要安装任何客户端。

---

**画面 [0:20 - 0:30]**
切到架构图动画页面（S1），展示完整数据流动。

**配音**：
> 这就是 Shield CLI——一个浏览器优先的内网服务网关。SSH 终端、Windows 桌面、数据库管理、Web 应用，全部通过浏览器一条命令直达。

---

### 第二幕：安装（0:30 - 1:10）

---

**画面 [0:30 - 0:50]**
终端录制（S2），展示三种安装方式，每种 5 秒快速切换：

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Linux
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# Windows (PowerShell)
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

**配音**：
> 安装非常简单。macOS 用 Homebrew，Linux 一行 curl，Windows 用 Scoop。三个平台，都是一条命令。

---

**画面 [0:50 - 1:00]**
终端执行 `shield --version`，显示版本号。

**配音**：
> 安装完成后，验证一下版本。Shield CLI 是一个单二进制文件，没有任何依赖。

---

**画面 [1:00 - 1:10]**
终端执行 `shield start`，显示 Web UI 启动信息。

**配音**：
> 执行 shield start，启动本地管理面板。打开 localhost 8181，我们来看看 Web UI。

---

### 第三幕：Web UI 模式（1:10 - 2:10）

---

**画面 [1:10 - 1:30]**
浏览器打开 `http://localhost:8181`（S3），展示 Dashboard 空状态。

**配音**：
> 这是 Shield CLI 的 Web 管理面板。你可以在这里添加、管理和连接所有内网服务。我们来添加第一个应用。

---

**画面 [1:30 - 1:50]**
点击"添加应用"，填写：
- 协议：SSH
- 地址：10.0.0.5
- 端口：22
- 用户名 / 密码

保存。

**配音**：
> 选择协议，填入内网服务器的 IP 和端口，输入凭证，保存。配置会加密存储在本地。

---

**画面 [1:50 - 2:10]**
点击"连接"按钮，状态变为"已连接"，显示 Access URL。点击 URL，浏览器新标签页打开 SSH 终端。

**配音**：
> 点击连接，Shield CLI 会建立加密隧道，生成一个公网 URL。把这个链接发给任何人，对方在浏览器里打开就能直接操作 SSH 终端。不需要安装任何软件。

---

### 第四幕：场景一 — SSH 终端（2:10 - 2:50）

---

**画面 [2:10 - 2:30]**
命令行模式演示（S4）：

```bash
shield ssh 10.0.0.5
```

终端输出连接信息和 URL。

**配音**：
> 除了 Web UI，你也可以直接在终端创建隧道。一条命令，适合服务器环境和脚本自动化。

---

**画面 [2:30 - 2:50]**
浏览器中的 SSH 终端（S5），执行几个命令：

```bash
uname -a
ls -la
htop
```

**配音**：
> 浏览器中的终端体验和原生 SSH 客户端几乎一致。支持颜色、Tab 补全、快捷键。你正在通过加密隧道，在浏览器里操作一台远程服务器。

---

### 第五幕：场景二 — Windows 远程桌面（2:50 - 3:40）

---

**画面 [2:50 - 3:10]**
终端执行（S6）：

```bash
shield rdp 10.0.0.10
```

显示连接成功和 URL。

**配音**：
> 远程 Windows 桌面同样简单。shield rdp 加上 IP 地址，隧道建立后，浏览器自动打开。

---

**画面 [3:10 - 3:30]**
浏览器中的 RDP 桌面（S7），移动鼠标、打开文件管理器、右键菜单。

**配音**：
> 整个 Windows 桌面在浏览器中实时渲染。鼠标、键盘、剪贴板全部支持。对方不需要安装 RDP 客户端，一个链接就能远程操作。

---

**画面 [3:30 - 3:40]**
手机模拟器画面（S10），用 Safari 打开同一个 RDP URL，展示手机上的 Windows 桌面。

**配音**：
> 甚至可以在手机上打开同一个链接。出差时用手机操作 Windows 桌面，完全没问题。

---

### 第六幕：场景三 — MySQL 数据库管理（3:40 - 4:30）

---

**画面 [3:40 - 3:55]**
终端执行（S8）：

```bash
shield plugin add mysql
shield mysql 10.0.0.20:3306 --db-user root --db-pass ****
```

**配音**：
> Shield CLI 通过插件系统支持更多服务。安装 MySQL 插件后，一条命令就能在浏览器中管理数据库。

---

**画面 [3:55 - 4:15]**
MySQL Web 管理界面（S9），演示：
1. 左侧数据库列表
2. 点击一张表，查看数据
3. 执行一条 SQL 查询
4. 点击 CSV 导出

**配音**：
> 浏览器中的数据库管理界面。浏览表结构、翻页查看数据、执行 SQL 查询、一键导出 CSV。默认只读模式，防止误操作。

---

**画面 [4:15 - 4:30]**
Web UI 中添加 MySQL 应用，勾选"Read-Only Mode"复选框。

**配音**：
> 在 Web UI 中添加 MySQL 应用时，可以勾选只读模式。开启后，远程用户无法执行任何写操作，前后端双重拦截，安全可控。

---

### 第七幕：总结 + 号召行动（4:30 - 5:00）

---

**画面 [4:30 - 4:45]**
回到架构图动画页面（S1），粒子流动状态。覆盖文字：

```
SSH Terminal    →
DB Admin        → Shield CLI → 加密隧道 → 浏览器
Desktop / Web   →

一条命令，浏览器访问一切内部服务
```

**配音**：
> Shield CLI 不只是一个隧道工具，它是一个统一的浏览器入口。SSH 终端、Windows 桌面、数据库管理——通过一条命令，在浏览器中安全访问和操作任何内部服务。无需 VPN，无需安装客户端。

---

**画面 [4:45 - 5:00]**
展示 GitHub 页面和文档站（S11）：

```
GitHub:  github.com/fengyily/shield-cli
文档:    docs.yishield.com
安装:    brew install shield-cli
```

**配音**：
> Shield CLI 完全开源，Apache 2.0 协议。GitHub 搜索 shield-cli，Star 支持一下。完整文档在 docs.yishield.com。感谢观看，我们下期见。

---

## 配音生成指南

### ElevenLabs 设置
- **Voice**: 选择 "Josh"（沉稳专业）或 "Rachel"（清晰友好）
- **Stability**: 0.5（自然变化）
- **Clarity**: 0.75（清晰但不机械）
- **Style**: 0（不夸张）

### 逐段生成
将每个 `> 引用块` 的文字单独生成一段音频，方便后期对齐画面。

### 文件命名
```
audio-01-hook.mp3
audio-02-install.mp3
audio-03-webui.mp3
...
```

## VHS Tape 脚本参考

安装演示：

```tape
Output docs/demo/demo-install.gif
Set Shell "bash"
Set FontSize 16
Set Width 800
Set Height 400
Set Theme "Dracula"

Type "brew tap fengyily/tap && brew install shield-cli"
Enter
Sleep 3s
Type "shield --version"
Enter
Sleep 2s
```

SSH 演示：

```tape
Output docs/demo/demo-ssh-video.gif
Set Shell "bash"
Set FontSize 16
Set Width 800
Set Height 400
Set Theme "Dracula"

Type "shield ssh 10.0.0.5"
Enter
Sleep 5s
```

## 剪辑节奏参考

| 时间点 | 转场 | 备注 |
|--------|------|------|
| 0:00 | 淡入 | 黑屏→文字 |
| 0:10 | 硬切 | 文字→终端 |
| 0:20 | 硬切 | 终端→架构图 |
| 0:30 | 淡出淡入 | 架构图→安装 |
| 1:10 | 硬切 | 终端→浏览器 |
| 2:10 | 淡出淡入 | Web UI→CLI 模式 |
| 2:50 | 硬切 | SSH→RDP |
| 3:40 | 淡出淡入 | RDP→MySQL |
| 4:30 | 淡出淡入 | MySQL→总结 |

全片节奏：**快-慢-快**。开头抓注意力（快），中间演示讲清楚（慢），结尾干脆利落（快）。

## 背景音乐

推荐免费可商用的轻电子/Lo-fi 音乐：
- **YouTube Audio Library** — 搜索 "technology" 或 "corporate"
- **Pixabay Music** — 免费可商用
- 音量：配音的 10-15%，仅作氛围衬底
