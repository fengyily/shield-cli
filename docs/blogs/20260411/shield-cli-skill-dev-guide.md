> 当 AI Agent 需要调用外部能力时，Skill 就是它的"技能包"。本文以一个文旅素材搜索 Skill 为例，带你走完**本地开发 → 调试 → 发布 → 安装使用**的完整流程。核心工具只有一个 —— **Shield CLI**。

---

## 背景：什么是 Skill？

Skill 是符合 Agent Skills 规范的能力描述包。一个典型的 Skill 包含：

```
my-skill/
├── .well-known/agent-skills/index.json   # 发现文件（描述、权限、环境变量）
├── SKILL.md                               # 技能指令（Agent 读取并执行）
├── config.json                            # 运行时配置（API 地址等）
└── references/                            # 补充文档
```

Agent 通过读取 `SKILL.md` 了解如何调用你的 API，而 `index.json` 告诉平台这个 Skill 需要什么权限、支持哪些环境变量。

我们今天要开发的 Skill 很简单：**通过一个 POST API 搜索素材，将结果以表格形式展示给用户**。API 是现成的，我们要做的是把它"包装"成 Agent 能理解的 Skill。

---

## 一、安装 Shield CLI

Shield CLI 是 YiShield 提供的命令行工具，用于本地开发、调试和部署应用。

### macOS（推荐使用 Homebrew）

```bash
brew tap fengyily/tap
brew install shield-cli
```
### Windows（使用 PowerShell 一键安装）

```bash
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```
### 手动安装

前往 [Shield CLI 官方文档](https://docs.yishield.com/guide/install.html) 根据对应平台选择安装方式

### 验证安装

```bash
f1@F1s-MacBook-Pro ~ % shield --version
shield version 0.3.11
```
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/59099d37206d4574bfe82c1bb1ea9639.png)
```
f1@F1s-MacBook-Pro ~ % shield start 

   _____ __    _       __    __   ________    ____
  / ___// /_  (_)__   / /___/ /  / ____/ /   /  _/
  \__ \/ __ \/ // _ \/ // __  / / /   / /    / /
 ___/ / / / / //  __/ // /_/ / / /___/ /____/ /
/____/_/ /_/_/ \___/_/ \__,_/  \____/_____/___/
  Shield CLI - Secure Tunnel Connector

  ├─ Version:    0.3.11
  ├─ Go:         go1.25.8
  └─ Platform:   darwin/arm64

  ──────────────────────────────────────────────────

  ✓ Service is already running

  Web UI: http://localhost:8181

  Commands:
    shield stop         Stop the service
    shield uninstall    Remove the service

f1@F1s-MacBook-Pro ~ % 
```
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/f2a6fc7b96894d37ae431cc9cc9f38e1.png)


---

## 二、本地启动 API 服务

Shield CLI 最强大的能力之一是**将你的本地服务暴露为可调试的 API 端点**。

### 启动本地服务

假设你的后端 API 已经在本地开发完毕（比如一个 Python FastAPI 或 Go 服务），先把它跑起来：

```bash
# 示例：启动文旅素材搜索 API
curl 'http://localhost:8080/api/v1/resources/search' \
  --data-raw '{"query":"雁荡山","use_llm":true,"page":1,"page_size":1}'
```

服务启动后，通过 `http://127.0.0.1:8080` 即可在本地访问。

### 本地测试 API

用 `curl` 快速验证接口是否正常：

```bash
root@fenyi-MS-7E44:~# curl -X POST http://172.16.3.60:8080/api/v1/resources/search   -H "Content-Type: application/json"   -d '{"query": "雁荡山", "page_size": 1, "use_llm": false}'|jq
```

你应该能看到类似这样的响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 42,
    "items": [...],
    "page": 1,
    "page_size": 5
  }
}
```
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/8aff4d4fe0b148798862b98384b5df66.png)

---

## 三、WebUI 配置，获取外网访问地址

本地 API 跑通了，但 Agent 平台（如 OpenClaw）在云端，它无法直接访问你的 `127.0.0.1`。这时需要通过 Shield 的**隧道功能**将本地服务映射到外网。

### 启动隧道
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/0c642095a2d74df38d58fda9bc95c48d.png)

通过 `http://127.0.0.1:8181` 访问 Shield 的 Web 界面，添加一个 HTTP 协议，端口 8080 ，IP 为 127.0.0.1 的应用，点击连接，等待几秒，Shield 会为你分配一个外网可达的地址，格式类似：
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/d749553dd5534eefb25c59eac968a9b4.png)

```
https://d2128144f966ce17-yishield.cn01.apps.yishield.com
```

<!-- 截图占位：终端显示 shield tunnel 启动成功，输出外网地址 -->
> **[图片待补充]** Shield 隧道启动成功，显示分配的外网地址

### 更新 Skill 配置

拿到外网地址后，将其填入 Skill 的 `config.json`：

```json
{
  "api_origin": "https://d2128144f966ce17-yishield.cn01.apps.yishield.com",
  "api_version_path": "/api/v1",
  "preview_play_path": "/api/v1/play",
  "trusted_media_origins": [
    "https://d2128144f966ce17-yishield.cn01.apps.yishield.com"
  ]
}
```

这样 Agent 在远端调用时，请求会通过 Shield 隧道转发到你的本地服务。**改一行代码 → 刷新 → 立即生效**，无需重新部署。

---

## 四、将 Skill 目录上传到 ClawHub 并发布

本地调试完毕、功能确认无误后，就可以把 Skill 发布到 [ClawHub](https://clawhub.ai/) —— Agent Skill 的集中分发平台。


1. 打开 [https://clawhub.ai/](https://clawhub.ai/)，登录账号
2. 进入「我的 Skills」页面，点击「上传 Skill」
3. 选择 `文件夹`
4. 填写 Skill 名称、描述、分类等元信息
5. 点击「发布」

![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/56ac376986bd4f69ab98210c19e224a6.png) 
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/bed1cfbc0a8541e2906dacf96ec6fe44.png)


发布后，平台会扫描你的 Skill ，当状态变为 Benign 就可以被其他用户在 ClawHub 市场中搜索和安装了。

---

## 五、在 OpenClaw 中安装和调试 Skill

### 5.1 安装 Skill

1. 打开 OpenClaw 客户端
2. 进入「Skill 市场」或「技能管理」
3. 搜索你刚发布的 Skill 名称（如 `culturetour-skill`）
4. 点击「安装」

我这里用的是 ArkClaw 

![](https://i-blog.csdnimg.cn/direct/f1a03836414849baaff679afc993e8d9.png)
### 5.2 配置环境变量（可选）

如果你需要让 Skill 指向不同的 API 地址（比如在开发阶段指向 Shield 隧道地址），可以在 OpenClaw 中设置环境变量：

| 变量名 | 说明 |
|--------|------|
| `WENLV_API_ORIGIN` | 覆盖 config.json 中的 api_origin（仅站点根，不含 /api/v1） |
| `TRADE_API_BASE` | 交易 API 基址（预留） |

### 5.3 对话调试

安装完成后，直接在 OpenClaw 对话框中用自然语言触发 Skill：

```
用户：帮我搜索黄山相关的文旅素材
```

Agent 会自动调用 Skill 中定义的 API，返回结构化的搜索结果表格：


![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/034597a593014ac6b31d8bff5e504ae1.png)


你可以继续测试完整流程：

```
用户：选 1、3
用户：预览第 1 条的视频
用户：确认购买
```
### 5.4 实时调试技巧

这里是 Shield CLI 真正发光的地方 —— **本地服务还在跑着，隧道还连着**。你在 OpenClaw 里触发的每一次 API 调用，都会实时打到你的本地服务：

- **看日志**：本地终端实时输出请求日志，方便排查参数和响应
- **改代码**：修改 API 逻辑后重启服务，OpenClaw 里下一次对话就能看到效果
- **改 Skill**：调整 `SKILL.md` 中的指令描述，重新打包上传，测试 Agent 行为变化

```
[本地终端]
POST /api/v1/resources/search  200  {"query":"黄山","page_size":5}  32ms
POST /api/v1/resources/search  200  {"query":"西湖 日落","page_size":3}  28ms
```

<!-- 截图占位：本地终端显示的实时请求日志 -->
> **[图片待补充]** Shield 隧道实时转发请求到本地，终端可见完整日志

---

## 总结：为什么选择 Shield CLI？

| 传统方式 | Shield CLI 方式 |
|----------|----------------|
| 改代码 → 推送 → CI/CD → 部署 → 测试 | 改代码 → 本地重启 → 立即测试 |
| 搭建测试环境，配置域名、证书 | 一条 `shield http 8080` 命令搞定 |
| 调试靠日志平台，延迟高 | 本地终端实时看请求/响应 |
| API 和 Skill 分开调试 | 端到端打通，Agent 直接调用本地服务 |

**Shield CLI 的核心价值在于缩短反馈循环**：

1. **零部署调试** —— 本地服务通过隧道直接暴露给 Agent 平台，省去了"部署到测试环境"的等待时间
2. **所见即所得** —— 修改 API 逻辑后秒级生效，修改 Skill 描述后重新上传即可验证 Agent 行为
3. **全链路可观测** —— 请求从 Agent → Shield 隧道 → 本地服务，每一跳都有日志，问题定位快
4. **开发生产一致** —— 本地用的和线上跑的是同一套代码、同一个 API，不存在"本地能跑线上挂"的问题
5. **低门槛** —— 不需要公网服务器、不需要配置 Nginx、不需要申请域名证书，一个 CLI 工具全搞定

对于 Skill 开发者来说，Shield CLI 把原本需要多个工具、多次部署才能完成的工作流，压缩成了**本地开发 + 一条命令**。当你的 Skill 在 OpenClaw 里跑通了，把 `config.json` 里的地址换成生产环境，就是正式上线。

---

**开始动手吧！** 


*Shield CLI 开源地址：[https://github.com/fengyily/shield-cli](https://github.com/fengyily/shield-cli)* 如果你觉得有用，就点个星星吧

