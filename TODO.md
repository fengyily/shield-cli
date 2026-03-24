# Shield CLI - TODO List

## 产品定位升级

> 详见 [docs/positioning.md](docs/positioning.md)

- [ ] **Phase 1: 调整措辞**
  - [ ] 重写 README，弱化"隧道"，强化"浏览器访问内部服务"
  - [ ] 调整项目 tagline / GitHub description
  - [ ] 文档站首页按使用场景重新组织

- [ ] **Phase 2: 插件扩展，支撑叙事**（P0 插件完成后）
  - [ ] 以新定位对外宣传

- [ ] **Phase 3: 平台化演进**（长期）
  - [ ] 权限控制 (RBAC)
  - [ ] 审计日志
  - [ ] 团队共享（隧道多人接入）
  - [ ] 服务目录（Web UI 展示可用服务）
  - [ ] 会话录制（SSH/RDP 回放）

---

## 插件开发

### P0 - High Priority

- [ ] **Redis Web Client Plugin**
  - Protocols: `redis`
  - Default Port: 6379
  - Features: Key 浏览/搜索、命令执行、内存监控、TTL 管理
  - Repo: `shield-plugin-redis`

- [ ] **SFTP/FTP File Browser Plugin**
  - Protocols: `sftp`, `ftp`
  - Default Port: 22 (SFTP) / 21 (FTP)
  - Features: Web 文件浏览器、拖拽上传下载、文件预览、权限管理
  - Repo: `shield-plugin-filebrowser`

## P1 - Medium Priority

- [ ] **Kafka Web Client Plugin**
  - Protocols: `kafka`
  - Default Port: 9092
  - Features: Topic 列表、消息预览/生产、Consumer Group 状态、Partition 信息
  - Repo: `shield-plugin-kafka`

- [ ] **Elasticsearch Web Client Plugin**
  - Protocols: `elasticsearch`, `es`, `opensearch`
  - Default Port: 9200
  - Features: 索引浏览、查询 DSL 编辑器、集群健康状态、Mapping 查看
  - Repo: `shield-plugin-elasticsearch`

## P2 - Lower Priority

- [ ] **Docker Web Manager Plugin**
  - Protocols: `docker`
  - Default Port: 2375
  - Features: 容器列表、日志查看、exec 终端、镜像管理
  - Repo: `shield-plugin-docker`

- [ ] **MongoDB Web Client Plugin**
  - Protocols: `mongodb`, `mongo`
  - Default Port: 27017
  - Features: Collection 浏览、文档查询/编辑、索引管理、聚合管道
  - Repo: `shield-plugin-mongodb`

- [ ] **etcd Web Browser Plugin**
  - Protocols: `etcd`
  - Default Port: 2379
  - Features: KV 浏览/编辑、Lease 管理、Watch 监听、集群状态
  - Repo: `shield-plugin-etcd`

- [ ] **MQTT Web Client Plugin**
  - Protocols: `mqtt`
  - Default Port: 1883
  - Features: Topic 订阅/发布、消息实时查看、QoS 设置
  - Repo: `shield-plugin-mqtt`

- [ ] **RabbitMQ Web Client Plugin**
  - Protocols: `rabbitmq`, `amqp`
  - Default Port: 5672
  - Features: Queue/Exchange 管理、消息投递/消费、Binding 配置
  - Repo: `shield-plugin-rabbitmq`

- [ ] **LDAP/AD Browser Plugin**
  - Protocols: `ldap`, `ldaps`
  - Default Port: 389 (LDAP) / 636 (LDAPS)
  - Features: 目录树浏览、用户/组搜索、属性编辑
  - Repo: `shield-plugin-ldap`

- [ ] **MinIO/S3 Web Browser Plugin**
  - Protocols: `s3`, `minio`
  - Default Port: 9000
  - Features: Bucket 浏览、文件上传下载、预签名 URL、Bucket Policy
  - Repo: `shield-plugin-s3`
