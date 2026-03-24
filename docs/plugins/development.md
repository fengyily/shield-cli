---
title: 插件开发指南 — 为 Shield CLI 开发自定义插件
description: 完整的 Shield CLI 插件开发教程，包括协议规范、项目模板、Web UI 开发和发布流程。
head:
  - - meta
    - name: keywords
      content: Shield CLI 插件开发, 自定义插件, Go插件, plugin development, Shield扩展
---

# 插件开发指南

本指南介绍如何为 Shield CLI 开发自定义插件，让 Shield 支持新的协议或服务类型。

## 概念

一个 Shield 插件本质上是一个**独立的可执行文件**，职责很简单：

1. 从 stdin 接收启动配置（JSON）
2. 在本地启动一个 Web 服务
3. 向 stdout 返回就绪状态和 Web 端口号
4. 保持运行，直到收到停止信号

Shield 主程序负责：
- 启动和停止插件进程
- 将插件的 Web 端口通过加密隧道暴露到公网
- 提供统一的命令行体验

## 通信协议

Shield 与插件之间通过 **stdin/stdout 单行 JSON** 通信。协议只有三种消息。

### 1. 启动请求（stdin）

Shield 启动插件后，立即通过 stdin 发送一行 JSON：

```json
{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"xxx","database":"mydb"}}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `action` | string | 固定为 `"start"` |
| `config.host` | string | 目标服务 IP |
| `config.port` | int | 目标服务端口 |
| `config.user` | string | 用户名（可能为空） |
| `config.pass` | string | 密码（可能为空） |
| `config.database` | string | 数据库名（可能为空） |

### 2. 就绪响应（stdout）

插件准备就绪后，向 stdout 写入一行 JSON：

**成功：**
```json
{"status":"ready","web_port":19876,"name":"My Plugin","version":"0.1.0"}
```

**失败：**
```json
{"status":"error","message":"cannot connect to service: connection refused"}
```

| 字段 | 类型 | 说明 |
|---|---|---|
| `status` | string | `"ready"` 或 `"error"` |
| `web_port` | int | 插件 Web 服务端口（status=ready 时必须） |
| `name` | string | 插件显示名称 |
| `version` | string | 插件版本号 |
| `message` | string | 错误信息（status=error 时） |

### 3. 停止请求（stdin）

Shield 退出时通过 stdin 发送：

```json
{"action":"stop"}
```

插件收到后应优雅退出。如果 5 秒内未退出，Shield 会强制终止进程。

### 超时

Shield 等待就绪响应的超时时间为 **15 秒**。如果插件在 15 秒内未响应，会被终止。

## 项目模板

以下是一个完整的插件项目结构：

```
shield-plugin-example/
├── main.go              # 入口，读取 stdin，启动 web server
├── handler.go           # HTTP 处理器（业务逻辑）
├── static/
│   └── index.html       # Web 界面（embed.FS 嵌入）
├── go.mod
├── go.sum
├── Makefile
└── .goreleaser.yml      # 多平台发布配置
```

### main.go 模板

```go
package main

import (
    "embed"
    "encoding/json"
    "io/fs"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

//go:embed static/*
var staticFS embed.FS

type StartRequest struct {
    Action string       `json:"action"`
    Config PluginConfig `json:"config,omitempty"`
}

type PluginConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user,omitempty"`
    Pass     string `json:"pass,omitempty"`
    Database string `json:"database,omitempty"`
}

type StartResponse struct {
    Status  string `json:"status"`
    WebPort int    `json:"web_port,omitempty"`
    Name    string `json:"name,omitempty"`
    Version string `json:"version,omitempty"`
    Message string `json:"message,omitempty"`
}

func main() {
    decoder := json.NewDecoder(os.Stdin)

    for {
        var req StartRequest
        if err := decoder.Decode(&req); err != nil {
            return // stdin closed
        }

        switch req.Action {
        case "start":
            handleStart(req.Config)
        case "stop":
            os.Exit(0)
        }
    }
}

func respond(resp StartResponse) {
    json.NewEncoder(os.Stdout).Encode(resp)
}

func handleStart(cfg PluginConfig) {
    // 1. 连接目标服务，验证可用性
    // conn, err := connectToService(cfg)
    // if err != nil {
    //     respond(StartResponse{Status: "error", Message: err.Error()})
    //     return
    // }

    // 2. 找一个可用端口
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        respond(StartResponse{Status: "error", Message: err.Error()})
        return
    }
    webPort := listener.Addr().(*net.TCPAddr).Port

    // 3. 设置 HTTP 路由
    mux := http.NewServeMux()
    // mux.HandleFunc("/api/...", yourHandler)
    staticSub, _ := fs.Sub(staticFS, "static")
    mux.Handle("/", http.FileServer(http.FS(staticSub)))

    // 4. 返回就绪响应
    respond(StartResponse{
        Status:  "ready",
        WebPort: webPort,
        Name:    "Example Plugin",
        Version: "0.1.0",
    })

    // 5. 启动 HTTP 服务
    go http.Serve(listener, mux)

    // 6. 等待停止信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
}
```

### 关键要点

1. **stdout 只用于协议通信** — 不要在 stdout 输出日志，使用 stderr
2. **端口自动分配** — 使用 `net.Listen("tcp", "127.0.0.1:0")` 让系统分配端口
3. **embed.FS 嵌入静态文件** — Web 界面打包进二进制，无需外部文件
4. **先连接后响应** — 在返回 `ready` 之前确认目标服务可连接，避免用户看到空白页面

## 注册插件

开发完成后，需要在 Shield 主程序的 `plugin/install.go` 中注册：

```go
var KnownPlugins = map[string]PluginInfo{
    "example": {
        Name:        "example",
        Source:      "your-org/shield-plugin-example",
        Protocols:   []string{"example"},
        DefaultPort: 9999,
    },
}
```

或者使用本地安装进行测试：

```bash
# 编译插件
go build -o shield-plugin-example .

# 本地安装
shield plugin add example --from ./shield-plugin-example

# 测试
shield example 127.0.0.1:9999
```

## 发布插件

### GoReleaser 配置

```yaml
# .goreleaser.yml
project_name: shield-plugin-example

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    binary: shield-plugin-example

archives:
  - format: tar.gz
    name_template: "shield-plugin-example_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
```

### 发布流程

```bash
# 打标签
git tag v0.1.0
git push origin v0.1.0

# GoReleaser 自动构建并发布到 GitHub Releases
goreleaser release
```

### 命名规范

| 项目 | 规范 |
|---|---|
| 仓库名 | `shield-plugin-<name>` |
| 二进制名 | `shield-plugin-<name>` |
| Release asset | `shield-plugin-<name>_<os>_<arch>.tar.gz` |

Shield 安装时会按此命名规范从 GitHub Releases 下载对应平台的二进制。

## 测试插件

### 手动测试协议

不通过 Shield，直接测试插件的 stdin/stdout 协议：

```bash
echo '{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"root"}}' | ./shield-plugin-mysql
# 期望输出: {"status":"ready","web_port":xxxxx,"name":"MySQL Web Client","version":"0.1.0"}
```

### 端到端测试

```bash
# 1. 本地安装
shield plugin add mysql --from ./shield-plugin-mysql

# 2. 启动（连接到本地 MySQL）
shield mysql 127.0.0.1:3306 --db-user root --db-pass root

# 3. 验证 Web 界面可访问
# 4. 验证通过公网 URL 也能访问
```

## Web UI 开发建议

### 推荐技术栈

对于 Shield 插件的 Web 界面，推荐纯前端方案（无框架依赖）：

- **HTML + Vanilla JS** — 最轻量，适合简单界面
- **单个 HTML 文件** — 所有 CSS/JS 内联，通过 `embed.FS` 嵌入
- 避免使用 npm 构建工具链 — 增加复杂度但对插件场景收益有限

### 界面设计原则

1. **轻量** — 单文件 HTML，gzip 后控制在 50KB 以内
2. **自适应** — 支持窄屏（可能在嵌入式浏览器中打开）
3. **安全默认** — 默认只读模式，写操作需手动解锁
4. **响应式反馈** — 操作后给用户明确的反馈（加载状态、成功/失败提示）

## 现有插件参考

- **shield-plugin-mysql** — 项目内实现，位于 `plugins/mysql/`，可作为开发新插件的参考

## 下一步

- [插件系统概览](/plugins/)
- [MySQL 插件文档](/plugins/mysql)
- [命令参考](/reference/commands)
