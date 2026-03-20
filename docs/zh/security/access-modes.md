---
title: 访问模式 — 可见模式与隐身模式
description: Shield CLI 提供可见模式（公开 URL）和隐身模式（需授权码）两种访问模式，控制 Access URL 的可见性和安全级别。
head:
  - - meta
    - name: keywords
      content: Shield CLI 访问模式, 可见模式, 隐身模式, 访问控制, 授权码, 隧道安全
---

# 访问模式

Shield CLI 提供两种访问模式，控制 Access URL 的可见性和安全级别。

## 可见模式（默认）

Access URL 任何人都可以直接访问，适合内部团队使用或信任的场景。

```bash
# 默认就是可见模式
shield ssh 10.0.0.5

# 指定接入节点（如选择香港节点）
shield ssh 10.0.0.5 --visable=HK
```

### 特点

- Access URL 可直接分享，打开即用
- 适合临时分享、内部协作
- 可通过 `--visable=<节点名>` 选择就近的接入节点

## 隐身模式

Access URL 需要额外的授权码才能访问，适合敏感服务。

```bash
shield ssh 10.0.0.5 --invisible
```

### 特点

- 除 Access URL 外，还会生成一个 Auth URL
- 访问者需要先通过 Auth URL 输入授权码
- 适合对外暴露敏感内网服务

## 如何选择

| 场景 | 推荐模式 |
|---|---|
| 内部团队协作 | 可见模式 |
| 给同事临时分享 | 可见模式 |
| 暴露生产环境服务 | 隐身模式 |
| 给外部客户演示 | 隐身模式 |
