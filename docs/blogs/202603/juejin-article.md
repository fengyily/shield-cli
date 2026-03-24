# 掘金

**发布地址**: https://juejin.cn/editor/drafts/new

**标签建议**: Go、开源、远程桌面、SSH、运维工具

---

**标题:**

开源实践：用 Go 实现浏览器直连内网 RDP/SSH/VNC，告别 VPN 和客户端

**正文:**

分享一个我最近在做的开源项目 **Shield CLI**，一个用 Go 写的命令行工具，主要解决远程访问内网服务的痛点。

## 它做了什么

Shield CLI 在内网机器和外部之间建立加密隧道，并通过 HTML5 技术将 RDP 桌面、VNC 会话、SSH 终端直接渲染在浏览器中。访问端不需要安装任何软件——打开浏览器，输入地址，就能看到完整的远程桌面或终端。

## 和 ngrok / frp 有什么不同

ngrok 和 frp 本质上是端口转发工具，暴露端口后你还是需要对应的客户端（RDP 客户端、VNC Viewer 等）来连接。Shield CLI 多做了一步：协议层面的 HTML5 渲染，把 RDP/VNC/SSH 协议翻译成浏览器可以直接呈现的内容。

## 技术细节

- 语言：Go，编译为单二进制文件
- 平台：macOS / Linux / Windows
- 部署：支持 Docker，也支持 apt/yum 包管理器安装
- 管理：内置 Web UI 仪表盘，可视化管理隧道
- 协议：SSH、RDP、VNC、HTTP
- 许可证：Apache 2.0

## 实际使用场景

- 远程办公：浏览器直接操作公司内网机器的桌面
- 运维：手机浏览器 SSH 到服务器处理紧急问题
- 临时协作：给外部人员发一个链接即可访问，用完即关，不用配置 VPN 账号

项目还在持续迭代中，欢迎 Star、提 Issue 和 PR。

GitHub: https://github.com/fengyily/shield-cli
