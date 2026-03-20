---
title: 自定义服务器 — 使用私有 Shield 网关
description: 配置 Shield CLI 连接到私有服务器，适合有合规要求或私有部署需求的企业用户。
head:
  - - meta
    - name: keywords
      content: Shield CLI 自定义服务器, 私有服务器, 私有部署, 企业部署, --server 参数
---

# 自定义服务器

默认情况下，Shield CLI 连接到 `https://console.yishield.com/raas` 公共服务。如果你部署了私有服务端，可以通过 `--server` 参数指向自己的服务器。

## 使用方法

```bash
shield ssh 10.0.0.5 --server https://your-server.com/raas
```

## 适用场景

- 企业内部部署了私有 Shield 服务端
- 需要数据完全不经过公网
- 合规要求不允许使用外部 SaaS 服务

## 注意事项

- 私有服务端需要与 Shield CLI 版本兼容
- 确保 Shield CLI 所在机器可以访问自定义服务器地址
