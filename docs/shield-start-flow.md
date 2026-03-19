# Shield Start 完整流程文档

## 总体架构

```
shield start [port]
    │
    ├── 1. 启动阶段（一次性）
    │   ├── 加载/创建 Credentials
    │   ├── 尝试建立主隧道（Server + API Tunnel）
    │   └── 启动 Web Server (默认 :8181)
    │
    ├── 2. 用户操作阶段（前端）
    │   ├── 添加应用 → 仅本地保存
    │   ├── 点击连接 → 触发后端连接流程
    │   └── 轮询状态 → 等待 connected 后跳转
    │
    └── 3. 连接阶段（后端，每个 App）
        ├── 调用 quick-setup API
        ├── 建立 App Tunnel（动态隧道）
        └── 等待隧道真正连通 → 返回 connected
```

---

## 阶段一：启动

### 1.1 入口 `cmd/start.go`

```
shield start [port]
    ↓
setupLog(verbose)              ← 配置日志级别（-v 启用 Debug）
PrintBanner()
web.NewServer(port)            ← 创建服务器
srv.Start()                    ← 启动（阻塞）
```

### 1.2 创建服务器 `web/server.go:NewServer`

```
NewServer(port)
    ↓
config.GetOrCreateCredentials()    ← 从磁盘加载加密凭证，或用机器指纹生成新的
    ↓
返回 Server{
    port:    port,
    store:   config.NewAppStore(),       ← 应用配置存储（本地加密）
    connMgr: NewConnectionManager(creds) ← 连接管理器（持有共享凭证）
}
```

### 1.3 启动服务 `web/server.go:Start`

```
srv.Start()
    ↓
s.connMgr.SetupMainTunnel()       ← 尝试建立主隧道
    ├── 有保存的 Connector 信息？
    │   ├── 是 → setupMainTunnelWithInfo() → 建主隧道 + 等待连通
    │   └── 否 → 跳过，等第一次连接 App 时再建
    ↓
注册 HTTP 路由：
    /api/apps        → handleApps (GET=列表, POST=创建)
    /api/apps/{id}   → handleAppByID (GET/PUT/DELETE)
    /api/connect/{id}   → handleConnect (POST)
    /api/disconnect/{id} → handleDisconnect (POST)
    /api/status/{id}    → handleStatus (GET)
    /                   → 静态文件 (index.html)
    ↓
http.ListenAndServe("127.0.0.1:{port}")
```

### 1.4 建立主隧道 `web/connect.go:setupMainTunnelWithInfo`

```
setupMainTunnelWithInfo(creds)
    ↓
① 确定本地端口
    ├── creds.LocalPort 有值且可用？→ 复用
    └── 否则 findAvailablePort(4000, 5000)
    ↓
② 创建 tunnel.Manager
    connInfo = {ExternalIP, APIPort, TunnelPort, ConnUsername, ConnPassword}
    ↓
③ mgr.CreateMainTunnel(APIPort, localPort)
    ├── 建 chisel 客户端，连接 ws://{ExternalIP}:{TunnelPort}
    ├── 映射: R:{APIPort}:localhost:{localPort}
    └── client.Start() 非阻塞返回，隧道在后台连接
    ↓
④ mgr.WaitReady("main", 30s)     ← 阻塞等待 chisel 真正 Connected
    ├── 监听 ConnectedCh（由 chiselLogWriter 触发）
    ├── chisel 日志出现 "Connected (Latency..." → close(ConnectedCh)
    └── 超时 30s → 返回失败
    ↓
⑤ 设置状态
    cm.mainMgr = mgr
    cm.mainReady = true
    ↓
⑥ 保存 localPort 到凭证
    ↓
⑦ 启动本地 API 服务
    go startConnLocalAPI(localPort, mgr, connInfo)
```

**主隧道映射关系：**
```
远程 Server:{TunnelPort}  ←chisel→  本机
远程 :{APIPort}  ←→  localhost:{localPort}   (API Tunnel)
```

---

## 阶段二：前端交互

### 2.1 页面加载 `web/static/index.html`

```
浏览器打开 http://127.0.0.1:8181
    ↓
loadApps()                     ← 立即加载应用列表
setInterval(loadApps, 5000)    ← 每 5 秒刷新
```

### 2.2 loadApps() — 获取应用列表

```
GET /api/apps
    ↓ server.go:handleApps
返回: [{AppConfig + conn_status}]
    ↓
前端渲染每个应用卡片：
    ├── idle        → 显示 [Connect] 按钮
    ├── connecting  → 显示 [Cancel] 按钮 + 脉冲动画
    ├── connected   → 显示 [Disconnect] + site_url 链接
    └── failed      → 显示错误信息
```

### 2.3 添加应用 — 仅本地保存

```
用户填写表单 → POST /api/apps
    ↓ server.go:handleApps (POST)
s.store.Add(app)               ← 仅保存到本地加密文件
                                  不调用任何远程 API
```

### 2.4 点击连接 `connectApp(id)`

```
用户点击 [Connect]
    ↓
POST /api/connect/{id}
    ↓ server.go:handleConnect
① 从 store 读取 AppConfig
② 构造 ConnectParams
③ cm.Connect(appID, params)
    ├── 创建 ActiveConnection{Status: "connecting"}
    ├── go cm.doConnect(appID, params, conn)   ← 后台执行
    └── 立即返回 {status: "connecting"}
    ↓
前端收到响应
    ├── toast("Connecting...")
    ├── loadApps()               ← 刷新 UI 显示 "Connecting..."
    └── startPolling(id)         ← 开始轮询
```

### 2.5 轮询状态 `startPolling(id)`

```
每 2 秒：
    GET /api/status/{id}
        ↓ server.go:handleStatus
        返回 ConnectResult{Status, SiteURL, AuthURL, Error, AppID}
    ↓
判断 status：
    ├── "connecting" → 继续轮询，刷新 UI
    ├── "connected"  → 停止轮询 → 跳转（见 2.6）
    └── "failed"     → 停止轮询 → 显示错误 toast
```

### 2.6 跳转逻辑（连接成功后）

```
status === "connected"
    ↓
clearInterval(pollingTimer)     ← 停止轮询
loadApps()                      ← 刷新 UI
    ↓
if (res.data.site_url) {
    toast("Connected! Opening access URL...")
    setTimeout(() => {
        window.open(site_url, '_blank')    ← 500ms 后新标签页打开
    }, 500)
}
```

---

## 阶段三：后端连接流程

### 3.1 doConnect 完整流程 `web/connect.go:doConnect`

```
doConnect(appID, params, conn)
    ↓
━━━ Step 1: 调用 Quick-Setup API ━━━
    creds = cm.GetCreds()
    callQuickSetupAPI(params, creds)
    ├── POST {server}/api/public/quick-setup
    │   Body: {protocol, ip, port, connector_name, password, ...}
    ├── 重试策略（最多 5 次）：
    │   ├── 401 → RefreshCreds() 重新加载凭证
    │   ├── 429 → 等待 attempt*3 秒
    │   └── EOF/timeout/refused → 等待 attempt 秒
    └── 返回: QuickSetupResponse
        ├── Connector: {ExternalIP, APIPort, Username, Password}
        ├── App: {AppID, SiteURL, Resource{IP, Port}}
        └── APIKey: {NHPServer, Code, AppID}
    ↓
━━━ Step 2: 更新凭证 ━━━
    newCreds.ConnectorName = resp.Connector.Username
    newCreds.Password = resp.Connector.Password
    newCreds.ExternalIP = resp.Connector.ExternalIP
    newCreds.APIPort = resp.Connector.APIPort
    ...
    newCreds.EncryptAndSave()
    cm.UpdateCreds(newCreds)
    ↓
━━━ Step 3: 确保主隧道就绪 ━━━
    cm.ensureMainTunnel(resp, tunnelPort)
    ├── cm.mainReady == true? → 跳过（已建立）
    └── 否 → setupMainTunnelWithInfo()（见 1.4）
    ↓
━━━ Step 4: 创建 App 动态隧道 ━━━
    rport = resp.App.Resource.Port
    cm.mainMgr.CreateDynamicTunnel(rport, params.IP, params.Port)
    ├── 建新的 chisel 客户端
    ├── 连接 ws://{ExternalIP}:{TunnelPort}
    ├── 映射: R:127.0.0.1:{rport}:{params.IP}:{params.Port}
    └── client.Start() 非阻塞返回
    ↓
━━━ Step 5: 等待隧道真正连通 ━━━
    cm.mainMgr.WaitReady(rport, 30s)     ← 关键！阻塞等待
    ├── 等待 ConnectedCh 被关闭
    ├── chisel 日志 "Connected (Latency..." → 触发
    ├── 成功 → 继续
    └── 超时 → setError("App tunnel connect timeout")
    ↓
━━━ Step 6: 构造结果 ━━━
    siteURL = resp.App.SiteURL
    authURL = (invisible mode ? 构造认证 URL : "")
    ↓
━━━ Step 7: 更新状态为 connected ━━━
    conn.Result = ConnectResult{
        Status:  "connected",      ← 前端轮询到此值后跳转
        SiteURL: siteURL,
        AuthURL: authURL,
        AppID:   resp.App.AppID,
    }
    ↓
━━━ Step 8: 等待断开信号 ━━━
    <-conn.StopCh                  ← 阻塞，直到用户点 Disconnect
```

### 3.2 隧道连通检测 `tunnel/manager.go`

```
chisel 内部日志流：
    "client: Connecting to ws://121.43.154.105:62888"
    "client: Connection error: server: access to 'R:...' denied"   ← 还没通
    "client: Retrying in 100ms..."
    "client: Connected (Latency 345.743458ms)"                     ← 真正连通！

检测机制：
    chiselLogWriter.Write(p)
        ↓
    if "Connected (Latency" in line:
        close(entry.ConnectedCh)     ← 通知 WaitReady
        ↓
    WaitReady(key, timeout)
        select {
        case <-entry.ConnectedCh:    ← 收到信号
            entry.Status = StatusConnected
            return true
        case <-time.After(timeout):
            return false
        }
```

---

## 断开连接流程

```
用户点击 [Disconnect]
    ↓
POST /api/disconnect/{id}
    ↓
cm.Disconnect(appID)
    ├── close(conn.StopCh)                    ← 通知 doConnect 退出
    ├── cm.mainMgr.CloseTunnel(resourcePort)  ← 关闭该 App 的动态隧道
    ├── conn.Result.Status = "disconnected"
    └── delete(cm.connections, appID)
    ↓
doConnect goroutine:
    <-conn.StopCh unblocks → 退出
    ↓
主隧道不受影响，继续运行
```

---

## 关键时序图

```
时间轴  │  前端                    │  后端                         │  chisel
────────┼──────────────────────────┼───────────────────────────────┼──────────────────
 T+0    │ Click [Connect]          │                               │
 T+0    │ POST /api/connect/{id}   │                               │
 T+0    │                          │ status="connecting"           │
 T+0    │ 收到 "connecting"         │ go doConnect()                │
 T+0    │ startPolling(2s)         │                               │
 T+1    │                          │ callQuickSetupAPI()           │
 T+2    │ poll → "connecting"      │                               │
 T+3    │                          │ API 返回                      │
 T+3    │                          │ ensureMainTunnel() (跳过)     │
 T+3    │                          │ CreateDynamicTunnel()         │ Connecting...
 T+4    │ poll → "connecting"      │ WaitReady() 阻塞中...         │ access denied
 T+5    │                          │                               │ Retrying...
 T+6    │ poll → "connecting"      │                               │
 T+8    │ poll → "connecting"      │                               │ Connected ✓
 T+8    │                          │ WaitReady() 返回 true         │
 T+8    │                          │ status="connected"            │
 T+10   │ poll → "connected" ✓     │                               │
 T+10   │ toast("Connected!")      │                               │
 T+10.5 │ window.open(site_url)    │                               │
```

---

## 关键参数表

| 参数 | 值 | 说明 |
|------|-----|------|
| Web 端口 | 8181 (默认) | `shield start [port]` |
| Tunnel 端口 | 62888 (默认) | chisel 服务端端口 |
| 本地 API 端口 | 4000-5000 | 自动选取，持久化到凭证 |
| Quick-Setup 超时 | 30s | HTTP 请求超时 |
| Quick-Setup 重试 | 5 次 | 401/429/瞬态错误 |
| WaitReady 超时 | 30s | 等隧道连通的超时 |
| 前端轮询间隔 | 2s | startPolling |
| 前端刷新间隔 | 5s | loadApps 定时刷新 |
| 跳转延迟 | 500ms | connected 后打开新标签 |
| chisel 最大重试间隔 | 10s | 断线重连间隔 |

---

## Credentials 数据结构

```json
{
  "connector_name": "shield_a1b2c3d4e5f6",
  "password": "xxx",
  "local_port": 4000,
  "external_ip": "121.43.154.105",
  "api_port": 63699,
  "tunnel_port": 62888,
  "conn_username": "connector_user",
  "conn_password": "connector_pass"
}
```

- 使用机器指纹派生 AES-256-GCM 密钥加密存储
- 路径: `~/.shield-cli/.credential` (macOS/Linux) 或 `%LOCALAPPDATA%\ShieldCLI\.credential` (Windows)
- `HasConnectorInfo()` 检查 ExternalIP/APIPort/TunnelPort/ConnUsername/ConnPassword 是否齐全
