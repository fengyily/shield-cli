---
title: MySQL 插件 — 浏览器内数据库管理
description: Shield CLI MySQL 插件提供浏览器内的 Web SQL Client，支持数据库浏览、表结构查看、SQL 执行、结果排序、CSV 导出等功能。
head:
  - - meta
    - name: keywords
      content: Shield CLI MySQL, MySQL Web客户端, 数据库管理, Web SQL, MySQL浏览器, MariaDB, shield mysql
---

# MySQL 插件

通过 Shield CLI 在浏览器中打开完整的 MySQL / MariaDB 数据库管理界面。

![Shield MySQL 插件演示](/demo/demo-mysql.gif)

## 安装

```bash
shield plugin add mysql
```

验证安装：

```bash
shield plugin list
# NAME   VERSION  PROTOCOLS       INSTALLED
# mysql  v0.1.0   mysql, mariadb  2026-03-24T10:00:00+08:00
```

## 快速连接

```bash
# 连接本机 MySQL（默认端口 3306）
shield mysql

# 指定端口
shield mysql 3307

# 连接远程服务器
shield mysql 10.0.0.5

# 完整地址
shield mysql 10.0.0.5:3307
```

## 认证方式

### 交互式输入（推荐）

不传凭证参数时，会提示输入用户名、密码和数据库名：

```bash
shield mysql 127.0.0.1:3306

  🔐 Database credentials (press Enter to skip)

  Username [root]: root
  Password: ****
  Database (optional): mydb

  ✓ Connecting as root
    Database: mydb
```

- **Username** 默认 `root`，直接回车使用默认值
- **Password** 隐藏输入，不会回显
- **Database** 可选，回车跳过

### 命令行传参

```bash
# 使用数据库专用参数
shield mysql 127.0.0.1:3306 --db-user root --db-pass mypassword --db-name mydb

# 也兼容通用认证参数
shield mysql 127.0.0.1:3306 --username root --auth-pass mypassword
```

| 参数 | 别名 | 说明 |
|---|---|---|
| `--db-user` | `--username` | 数据库用户名 |
| `--db-pass` | `--auth-pass` | 数据库密码 |
| `--db-name` | — | 数据库名（可选） |
| `--readonly` | — | 强制只读模式，禁止写操作 |

## Web 管理界面

连接成功后，浏览器自动打开 Web SQL Client：

### 功能概览

| 功能 | 说明 |
|---|---|
| 数据库浏览 | 左侧栏列出所有数据库，点击切换 |
| 表列表 | 选中数据库后显示所有表，支持过滤和翻页 |
| 表结构查看 | 单击表名查看字段、类型、索引等信息 |
| SQL 执行 | 支持 SELECT、SHOW、DESCRIBE 等查询语句 |
| 结果排序 | 点击列头升序/降序排列 |
| CSV 导出 | 一键导出查询结果为 CSV 文件 |
| 复制结果 | 复制为 Tab 分隔文本，可直接粘贴到 Excel |
| 只读模式 | 默认开启，阻止 INSERT/UPDATE/DELETE 等写操作 |

### 只读模式

只读/读写模式完全由启动参数决定，页面上仅展示当前状态，远程用户无法更改。

**CLI 模式**：通过 `--readonly` 参数控制：

```bash
# 只读模式（推荐用于分享场景）
shield mysql 10.0.0.5:3306 --db-user root --readonly

# 读写模式（默认）
shield mysql 10.0.0.5:3306 --db-user root
```

**Web UI 模式**：在添加/编辑应用时，勾选 **Read-Only Mode** 复选框。

页面右上角显示当前模式：

- **🔒 Read-Only**（橙色标签）— 写操作被前后端双重拦截
- **🔓 Read-Write**（绿色标签）— 允许所有操作

只读模式下以下语句会被阻止：

```
INSERT, UPDATE, DELETE, DROP, ALTER, CREATE,
TRUNCATE, RENAME, REPLACE, GRANT, REVOKE
```

::: tip 安全建议
在公网暴露数据库管理界面时，建议开启只读模式，并配合 `--invisible` 隐身模式使用：
```bash
shield mysql 127.0.0.1:3306 --db-user readonly_user --readonly --invisible
```
:::

### 表过滤与翻页

当数据库包含大量表时：

- **过滤** — 表列表顶部的搜索框支持实时过滤，输入关键词即可筛选
- **翻页** — 超过 50 个表时自动分页，底部显示翻页控件

### 快捷键

| 快捷键 | 功能 |
|---|---|
| `Ctrl+Enter` / `Cmd+Enter` | 执行当前 SQL |
| `Tab` | 插入两个空格（不会切换焦点） |

### 数据库导航

1. 左侧栏显示所有数据库
2. 点击数据库名，切换到该数据库的表列表
3. 点击 **← Databases** 返回数据库列表
4. 单击表名 → 查看表结构
5. 双击表名 → 自动填入 `SELECT * FROM ... LIMIT 100` 并执行

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield mysql` | `127.0.0.1:3306` |
| `shield mysql 3307` | `127.0.0.1:3307` |
| `shield mysql 10.0.0.5` | `10.0.0.5:3306` |
| `shield mysql 10.0.0.5:3307` | `10.0.0.5:3307` |

`mariadb` 是 `mysql` 的别名，行为完全一致：

```bash
shield mariadb 127.0.0.1:3306 --db-user root
```

## 使用场景

### 场景一：远程调试生产数据库

在本地连接远程 MySQL，通过 Shield 暴露 Web 管理界面给团队成员：

```bash
shield mysql 10.0.0.5:3306 --db-user readonly --db-pass xxx --invisible
```

分享 Auth URL 给需要查看数据的同事，无需安装任何客户端。

### 场景二：Docker 中的 MySQL

直接连接 Docker 容器中的 MySQL：

```bash
# MySQL 容器监听在 3306
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8

# 通过 Shield 暴露
shield mysql 127.0.0.1:3306 --db-user root --db-pass root
```

### 场景三：内网数据库审计

以只读模式连接，让审计人员可以通过浏览器查看数据库：

```bash
shield mysql 192.168.1.100:3306 --db-user auditor --db-pass xxx --invisible
```

- 默认只读模式，无法执行写操作
- 隐身模式需要授权码才能访问
- 查询结果可导出 CSV 用于报告

### 场景四：临时分享数据库给外部合作伙伴

```bash
shield mysql 127.0.0.1:3306 --db-user report_user --db-pass xxx
```

对方只需浏览器即可访问，无需安装 MySQL 客户端、配置 VPN 或防火墙规则。

## API 接口

MySQL 插件在本地启动一个 HTTP 服务，提供以下 API：

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/info` | 服务器信息（版本、主机、用户） |
| GET | `/api/databases` | 数据库列表 |
| GET | `/api/tables?db=mydb` | 表列表 |
| GET | `/api/schema?db=mydb&table=users` | 表结构（字段、类型、索引） |
| POST | `/api/query` | 执行 SQL 查询 |

### 查询 API 示例

```bash
curl -X POST http://localhost:19876/api/query \
  -H 'Content-Type: application/json' \
  -d '{"sql": "SELECT * FROM users LIMIT 10", "db": "mydb"}'
```

响应：

```json
{
  "code": 200,
  "data": {
    "columns": ["id", "name", "email"],
    "rows": [
      {"id": 1, "name": "Alice", "email": "alice@example.com"}
    ],
    "count": 1,
    "duration": "2.3ms"
  }
}
```

## 故障排除

### 连接被拒绝

```
plugin error: cannot connect to MySQL at 127.0.0.1:3306: connection refused
```

检查 MySQL 服务是否正在运行：

```bash
# macOS
brew services list | grep mysql

# Linux
systemctl status mysql

# Docker
docker ps | grep mysql
```

### 认证失败

```
plugin error: Access denied for user 'root'@'172.17.0.1'
```

- 确认用户名和密码正确
- Docker 环境下注意 MySQL 的 `bind-address` 和用户权限配置
- 尝试创建允许远程访问的用户：
  ```sql
  CREATE USER 'shield'@'%' IDENTIFIED BY 'password';
  GRANT SELECT ON *.* TO 'shield'@'%';
  ```

### 插件未安装

```
unsupported protocol "mysql"
```

先安装插件：

```bash
shield plugin add mysql
```

## 下一步

- [插件系统概览](/plugins/)
- [插件开发指南](/plugins/development)
- [TCP 端口代理](/protocols/tcp-udp)（作为替代方案直接代理 3306 端口）
