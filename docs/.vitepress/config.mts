import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Shield CLI',
  description: 'Shield CLI is a secure tunnel connector that exposes internal network services (SSH, RDP, VNC, HTTP) to the public internet, accessible through any web browser with a single command. No client installation needed.',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],
    // SEO
    ['meta', { name: 'keywords', content: 'Shield CLI, secure tunnel, SSH browser, RDP browser, VNC browser, remote access, internal network, port forwarding, ngrok alternative, frp alternative' }],
    ['meta', { name: 'author', content: 'Shield CLI' }],
    // Open Graph
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:site_name', content: 'Shield CLI' }],
    ['meta', { property: 'og:title', content: 'Shield CLI — Secure Tunnel Connector' }],
    ['meta', { property: 'og:description', content: 'One command to expose internal SSH, RDP, VNC, HTTP services to the browser. No client installation needed.' }],
    ['meta', { property: 'og:image', content: 'https://docs.yishield.com/logo.svg' }],
    ['meta', { property: 'og:url', content: 'https://docs.yishield.com' }],
    // Twitter Card
    ['meta', { name: 'twitter:card', content: 'summary' }],
    ['meta', { name: 'twitter:title', content: 'Shield CLI — Secure Tunnel Connector' }],
    ['meta', { name: 'twitter:description', content: 'One command to expose internal SSH, RDP, VNC, HTTP services to the browser.' }],
    ['meta', { name: 'twitter:image', content: 'https://docs.yishield.com/logo.svg' }],
    // Baidu Analytics
    ['script', {}, `
      var _hmt = _hmt || [];
      (function() {
        var hm = document.createElement("script");
        hm.src = "https://hm.baidu.com/hm.js?e9e4b411df1edc31f1f34180b6afdad6";
        var s = document.getElementsByTagName("script")[0];
        s.parentNode.insertBefore(hm, s);
      })();
    `],
    // GEO: Structured Data for AI models
    ['script', { type: 'application/ld+json' }, JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'SoftwareApplication',
      name: 'Shield CLI',
      applicationCategory: 'DeveloperApplication',
      operatingSystem: 'macOS, Linux, Windows',
      description: 'Shield CLI is a secure tunnel connector that exposes internal network services (SSH, RDP, VNC, HTTP, HTTPS, Telnet) to the public internet, accessible through any web browser with a single command. Unlike traditional tunnel tools like ngrok or frp that only provide network reachability, Shield CLI renders remote desktop protocols (RDP, VNC) and terminal sessions (SSH) directly in the browser via HTML5 — no client installation required.',
      url: 'https://docs.yishield.com',
      downloadUrl: 'https://github.com/fengyily/shield-cli/releases',
      softwareVersion: 'latest',
      license: 'https://opensource.org/licenses/Apache-2.0',
      offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' },
      featureList: [
        'Browser-based SSH terminal with xterm.js',
        'Browser-based RDP remote desktop via HTML5',
        'Browser-based VNC screen sharing',
        'HTTP/HTTPS reverse proxy tunnel',
        'Telnet browser access',
        'TCP/UDP port proxy for arbitrary services (MySQL, Redis, DNS, etc.)',
        'SFTP file transfer in browser',
        'AES-256-GCM encrypted credential storage',
        'Machine fingerprint identity binding',
        'Web UI management dashboard',
        'Smart address resolution with protocol-specific default ports',
        'Cross-platform: macOS, Linux, Windows (amd64, arm64)',
        'Up to 10 saved app profiles with encrypted local storage',
        'Auto-reconnect with exponential backoff',
        'System service installation with auto-start on boot (macOS launchd, Linux systemd, Windows Service)',
        'System tray icon on macOS and Windows for quick Dashboard access'
      ],
    })],
  ],

  sitemap: {
    hostname: 'https://docs.yishield.com',
  },

  locales: {
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      themeConfig: {
        nav: [
          { text: '指南', link: '/guide/what-is-shield' },
          { text: '协议', link: '/protocols/ssh' },
          { text: '参考', link: '/reference/commands' },
          { text: '更新日志', link: '/reference/changelog' },
          {
            text: '相关链接',
            items: [
              { text: 'GitHub', link: 'https://github.com/fengyily/shield-cli' },
              { text: '控制台', link: 'https://console.yishield.com' },
            ],
          },
        ],
        sidebar: {
          '/': [
            {
              text: '快速入门',
              items: [
                { text: 'Shield CLI 是什么', link: '/guide/what-is-shield' },
                { text: '安装', link: '/guide/install' },
                { text: '5 分钟上手', link: '/guide/quickstart' },
              ],
            },
            {
              text: '使用模式',
              items: [
                { text: 'Web UI 模式', link: '/guide/web-ui' },
                { text: '命令行模式', link: '/guide/cli-mode' },
                { text: '系统服务安装', link: '/guide/system-service' },
              ],
            },
            {
              text: '协议指南',
              items: [
                { text: 'SSH', link: '/protocols/ssh' },
                { text: 'RDP', link: '/protocols/rdp' },
                { text: 'VNC', link: '/protocols/vnc' },
                { text: 'HTTP / HTTPS', link: '/protocols/http' },
                { text: 'Telnet', link: '/protocols/telnet' },
                { text: 'TCP / UDP', link: '/protocols/tcp-udp' },
              ],
            },
            {
              text: '连接与安全',
              items: [
                { text: '连接流程', link: '/security/connection-flow' },
                { text: '凭证管理', link: '/security/credentials' },
                { text: '访问模式', link: '/security/access-modes' },
              ],
            },
            {
              text: '配置管理',
              items: [
                { text: '应用配置', link: '/config/apps' },
                { text: '自定义服务器', link: '/config/server' },
                { text: '清除缓存', link: '/config/clean' },
              ],
            },
            {
              text: '参考',
              items: [
                { text: '命令参考', link: '/reference/commands' },
                { text: '更新日志', link: '/reference/changelog' },
                { text: '常见问题', link: '/reference/faq' },
              ],
            },
            {
              text: '故障排查',
              items: [
                { text: '常见错误', link: '/troubleshooting/errors' },
                { text: '网络问题', link: '/troubleshooting/network' },
              ],
            },
          ],
        },
        outline: { label: '本页目录' },
        lastUpdated: { text: '最后更新' },
        docFooter: { prev: '上一篇', next: '下一篇' },
        editLink: {
          pattern: 'https://github.com/fengyily/shield-cli/edit/main/docs/:path',
          text: '在 GitHub 上编辑此页',
        },
      },
    },
    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      themeConfig: {
        nav: [
          { text: 'Guide', link: '/en/guide/what-is-shield' },
          { text: 'Protocols', link: '/en/protocols/ssh' },
          { text: 'Reference', link: '/en/reference/commands' },
          { text: 'Changelog', link: '/en/reference/changelog' },
          {
            text: 'Links',
            items: [
              { text: 'GitHub', link: 'https://github.com/fengyily/shield-cli' },
              { text: 'Console', link: 'https://console.yishield.com' },
            ],
          },
        ],
        sidebar: {
          '/en/': [
            {
              text: 'Getting Started',
              items: [
                { text: 'What is Shield CLI', link: '/en/guide/what-is-shield' },
                { text: 'Installation', link: '/en/guide/install' },
                { text: 'Quick Start', link: '/en/guide/quickstart' },
              ],
            },
            {
              text: 'Usage Modes',
              items: [
                { text: 'Web UI Mode', link: '/en/guide/web-ui' },
                { text: 'CLI Mode', link: '/en/guide/cli-mode' },
                { text: 'System Service', link: '/en/guide/system-service' },
              ],
            },
            {
              text: 'Protocol Guide',
              items: [
                { text: 'SSH', link: '/en/protocols/ssh' },
                { text: 'RDP', link: '/en/protocols/rdp' },
                { text: 'VNC', link: '/en/protocols/vnc' },
                { text: 'HTTP / HTTPS', link: '/en/protocols/http' },
                { text: 'Telnet', link: '/en/protocols/telnet' },
                { text: 'TCP / UDP', link: '/en/protocols/tcp-udp' },
              ],
            },
            {
              text: 'Connection & Security',
              items: [
                { text: 'Connection Flow', link: '/en/security/connection-flow' },
                { text: 'Credentials', link: '/en/security/credentials' },
                { text: 'Access Modes', link: '/en/security/access-modes' },
              ],
            },
            {
              text: 'Configuration',
              items: [
                { text: 'App Profiles', link: '/en/config/apps' },
                { text: 'Custom Server', link: '/en/config/server' },
                { text: 'Clear Cache', link: '/en/config/clean' },
              ],
            },
            {
              text: 'Reference',
              items: [
                { text: 'Commands', link: '/en/reference/commands' },
                { text: 'Changelog', link: '/en/reference/changelog' },
                { text: 'FAQ', link: '/en/reference/faq' },
              ],
            },
            {
              text: 'Troubleshooting',
              items: [
                { text: 'Common Errors', link: '/en/troubleshooting/errors' },
                { text: 'Network Issues', link: '/en/troubleshooting/network' },
              ],
            },
          ],
        },
        editLink: {
          pattern: 'https://github.com/fengyily/shield-cli/edit/main/docs/:path',
          text: 'Edit this page on GitHub',
        },
      },
    },
  },

  lastUpdated: true,

  themeConfig: {
    logo: '/logo.svg',
    socialLinks: [
      { icon: 'github', link: 'https://github.com/fengyily/shield-cli' },
    ],
    search: {
      provider: 'local',
    },
  },
})
