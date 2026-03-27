---
title: 只是想查个数据，不想装 phpMyAdmin
description: 日常查数据、看表结构、跑几条 SQL，不需要几百 MB 的管理平台。Shield 插件 Docker 镜像只有 7-9MB，一行命令启动，打开浏览器就能用。
date: 2026-03-26
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: phpMyAdmin替代, pgAdmin替代, 轻量数据库工具, Docker数据库管理, MySQL Web查询, PostgreSQL Web查询, shield-mysql, shield-postgres
---

# 只是想查个数据，不想装 phpMyAdmin

> 日常工作中经常碰到一个场景：需要看下数据库里的数据，但手头没有趁手的工具。装 phpMyAdmin 太重，装 pgAdmin 更重。其实大多数时候只是想查几条数据、看看表结构，不需要那么大的东西。

---

## 日常查数据的烦恼

你有没有遇到过这些情况：

- 新到一台服务器，想看下数据库有什么表，发现没有任何管理工具
- 开发环境 Docker 跑了个 MySQL，想看下数据，又不想为了查数据再装一堆东西
- 同事问你某个表的数据，你得先打开 Navicat、连上数据库、找到表、截图发过去
- 树莓派上跑了个小项目，想看下 PostgreSQL 里的数据，pgAdmin 装不动（800MB 镜像 + 300MB 内存）

这些场景有个共同点：**需求很简单，就是查个数据，但工具很重**。

---

## 一行命令，打开浏览器就能查

Shield 的 MySQL 和 PostgreSQL 插件做成了 Docker 镜像，专门解决这个问题：

```bash
# MySQL
docker run -d -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_PASS=mypass \
  -p 8080:8080 fengyily/shield-mysql

# PostgreSQL
docker run -d -e DB_HOST=host.docker.internal -e DB_USER=postgres -e DB_PASS=mypass \
  -p 8080:8080 fengyily/shield-postgres
```

打开 `http://localhost:8080` ，浏览器里直接查数据。

不需要配置文件，不需要 PHP 环境，不需要注册服务器。环境变量传入数据库地址和凭证，启动即连接。

---

## 有多轻？

| | phpMyAdmin | pgAdmin | Shield MySQL | Shield PostgreSQL |
|---|---|---|---|---|
| 镜像大小 | ~400 MB | ~800 MB | **~7 MB** | **~9 MB** |
| 启动时间 | 5-10 秒 | 10-20 秒 | **< 1 秒** | **< 1 秒** |
| 运行内存 | ~200 MB | ~300 MB | **~10 MB** | **~10 MB** |
| 运行依赖 | PHP + Apache | Python + Flask | **无** | **无** |

7MB 的镜像，几秒钟下载完，启动不到 1 秒。这个体积是因为 Go 编译成静态二进制，基础镜像用的 `scratch`（空镜像），里面只有一个可执行文件和 CA 证书，没有操作系统、没有运行时、没有包管理器。

在树莓派、NAS、低配 VPS 这种资源有限的环境上，差别尤其明显——pgAdmin 跑不动的地方，这个轻轻松松。

---

## 能干什么

定位很明确：**日常查数据够用的浏览器工具**，不是全功能数据库管理平台。

具体来说：

- 浏览数据库、表、字段 — 左侧树形结构，点击切换
- 执行 SQL — 多标签编辑器，`Ctrl+Enter` 执行
- 查看表结构 — 字段类型、默认值、索引
- 结果导出 — CSV 导出，单元格复制，列排序
- 行级操作 — 插入、编辑、删除
- 建表/删表 — 可视化操作
- 只读模式 — 设置 `DB_READONLY=true`，写操作全部拦截

不能干什么（这些请用专业工具）：

- 存储过程 / 触发器编辑
- 数据库迁移
- 慢查询分析
- 多实例统一管理

---

## 几个典型用法

### 开发环境随手带一个

Docker Compose 里加两行，数据库和管理工具一起起来：

```yaml
services:
  mysql:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: root
    ports:
      - "3306:3306"

  mysql-web:
    image: fengyily/shield-mysql
    environment:
      DB_HOST: mysql
      DB_USER: root
      DB_PASS: root
    ports:
      - "8080:8080"
    depends_on:
      - mysql
```

`docker compose up -d`，`http://localhost:8080` 查数据。开发环境重建也是秒级。

### 临时用完就删

服务器上排查问题，不想留工具：

```bash
docker run -d --name tmp-viewer --network host \
  -e DB_HOST=127.0.0.1 -e DB_USER=readonly -e DB_PASS=xxx \
  -e DB_READONLY=true \
  fengyily/shield-mysql

# 查完删掉
docker rm -f tmp-viewer
```

### 低配设备

树莓派、NAS、1 核 512MB 的 VPS：

```bash
docker run -d -e DB_HOST=192.168.1.100 -e DB_USER=postgres -e DB_PASS=mypass \
  -p 8080:8080 fengyily/shield-postgres
```

9MB 镜像 + 10MB 运行内存，不抢资源。

### 同时看多个库

不同端口开不同实例：

```bash
# 开发库
docker run -d -e DB_HOST=10.0.0.1 -e DB_USER=root -e DB_PASS=dev \
  -p 8080:8080 fengyily/shield-mysql

# 测试库（只读）
docker run -d -e DB_HOST=10.0.0.2 -e DB_USER=root -e DB_PASS=test \
  -e DB_READONLY=true \
  -p 8081:8080 fengyily/shield-mysql

# PG 生产库（只读）
docker run -d -e DB_HOST=10.0.0.3 -e DB_USER=readonly -e DB_PASS=prod \
  -e DB_READONLY=true \
  -p 8082:8080 fengyily/shield-postgres
```

---

## 网络说明

Docker 容器访问宿主机数据库，根据系统选择方式：

```bash
# Linux：用 host 网络，最简单
docker run -d --network host \
  -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_PASS=mypass \
  fengyily/shield-mysql

# macOS / Windows Docker Desktop：用 host.docker.internal
docker run -d \
  -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_PASS=mypass \
  -p 8080:8080 fengyily/shield-mysql

# 数据库也是容器：放同一网络，用容器名
docker run -d --network my-net \
  -e DB_HOST=my-mysql -e DB_USER=root -e DB_PASS=mypass \
  -p 8080:8080 fengyily/shield-mysql
```

---

## 环境变量

| 变量 | MySQL 默认 | PostgreSQL 默认 | 说明 |
|---|---|---|---|
| `DB_HOST` | `127.0.0.1` | `127.0.0.1` | 数据库地址 |
| `DB_PORT` | `3306` | `5432` | 端口 |
| `DB_USER` | `root` | `postgres` | 用户名 |
| `DB_PASS` | — | — | 密码 |
| `DB_NAME` | — | `postgres` | 默认数据库 |
| `DB_READONLY` | `false` | `false` | 只读模式 |
| `WEB_PORT` | `8080` | `8080` | Web 端口 |

---

## 如果需要远程访问

上面说的都是本地或内网场景。如果你需要把数据库管理界面分享给外部同事——比如让外包团队查一下数据，或者在家访问公司内网的数据库——可以配合 [Shield CLI](https://github.com/fengyily/shield-cli) 使用：

```bash
shield plugin add mysql
shield mysql 10.0.0.5:3306 --db-user root --readonly
```

Shield CLI 会自动建立加密隧道，生成一个公网 URL，对方浏览器打开就能查。用的是同一套 Web 界面，只是多了远程访问能力。用完断开，链接即失效。

详细用法见：[MySQL 插件介绍](/blogs/blog-shield-cli-plugin-mysql) · [PostgreSQL 插件介绍](/blogs/blog-shield-cli-plugin-postgres)

---

## 试一下

```bash
# MySQL
docker run -d -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_PASS=root \
  -p 8080:8080 fengyily/shield-mysql

# PostgreSQL
docker run -d -e DB_HOST=host.docker.internal -e DB_USER=postgres -e DB_PASS=postgres \
  -p 8081:8080 fengyily/shield-postgres
```

Docker Hub：[fengyily/shield-mysql](https://hub.docker.com/r/fengyily/shield-mysql) · [fengyily/shield-postgres](https://hub.docker.com/r/fengyily/shield-postgres)

GitHub：[shield-cli](https://github.com/fengyily/shield-cli) · [shield-plugins](https://github.com/fengyily/shield-plugins)
