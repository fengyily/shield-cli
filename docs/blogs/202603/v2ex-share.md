# V2EX

**发布节点**: /go/programmer 或 /go/share

**发布地址**: https://www.v2ex.com/new

---

**标题:**

Shield CLI：开源工具，浏览器直接访问内网 SSH/RDP/VNC，不需要 VPN

**正文:**

做了一个开源命令行工具 Shield CLI，用来解决远程访问内网服务的问题。

和 ngrok、frp 的区别：Shield CLI 不只是端口转发，而是通过 HTML5 直接在浏览器中渲染 RDP 桌面、VNC 会话和 SSH 终端。不需要装任何客户端软件，打开浏览器就能用。

技术栈：Go 编写，单二进制文件部署，支持 macOS/Linux/Windows，提供 Docker 镜像和 Web 管理界面。Apache 2.0 开源协议。

典型场景：
- 在外面用浏览器直接 RDP 连回家里的 Windows 机器
- 给同事临时开个 SSH 访问，发个链接就行，不用配 VPN
- 手机浏览器也能操作远程终端

项目还在积极开发中，欢迎试用和反馈。

GitHub: https://github.com/fengyily/shield-cli
