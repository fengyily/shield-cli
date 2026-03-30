"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const child_process = __importStar(require("child_process"));
// 执行 shield-cli 命令的函数
function executeCommand(command) {
    return new Promise((resolve, reject) => {
        child_process.exec(command, (error, stdout, stderr) => {
            if (error) {
                resolve({
                    stdout,
                    stderr,
                    code: error.code || 1
                });
            }
            else {
                resolve({
                    stdout,
                    stderr,
                    code: 0
                });
            }
        });
    });
}
// 检查 shield-cli 是否安装
async function checkShieldCli() {
    const result = await executeCommand('shield --version');
    return result.code === 0;
}
// 启动 shield 服务
async function startShieldService() {
    const result = await executeCommand('shield start');
    if (result.code === 0) {
        vscode.window.showInformationMessage('Shield service started successfully');
    }
    else {
        vscode.window.showErrorMessage(`Failed to start Shield service: ${result.stderr}`);
    }
}
// 停止 shield 服务
async function stopShieldService() {
    const result = await executeCommand('shield stop');
    if (result.code === 0) {
        vscode.window.showInformationMessage('Shield service stopped successfully');
    }
    else {
        vscode.window.showErrorMessage(`Failed to stop Shield service: ${result.stderr}`);
    }
}
// 打开 shield web 界面
function openShieldWebUI() {
    vscode.commands.executeCommand('vscode.open', vscode.Uri.parse('http://localhost:8181'));
}
// 连接 SSH
async function connectSSH() {
    const host = await vscode.window.showInputBox({ prompt: 'Enter host (e.g., 127.0.0.1:22)' });
    if (host) {
        const result = await executeCommand(`shield ssh ${host}`);
        if (result.code === 0) {
            vscode.window.showInformationMessage('SSH connection established');
        }
        else {
            vscode.window.showErrorMessage(`Failed to connect SSH: ${result.stderr}`);
        }
    }
}
// 连接 RDP
async function connectRDP() {
    const host = await vscode.window.showInputBox({ prompt: 'Enter host (e.g., 127.0.0.1:3389)' });
    if (host) {
        const result = await executeCommand(`shield rdp ${host}`);
        if (result.code === 0) {
            vscode.window.showInformationMessage('RDP connection established');
        }
        else {
            vscode.window.showErrorMessage(`Failed to connect RDP: ${result.stderr}`);
        }
    }
}
// 连接 HTTP
async function connectHTTP() {
    const port = await vscode.window.showInputBox({ prompt: 'Enter port (e.g., 3000)' });
    if (port) {
        const result = await executeCommand(`shield http ${port}`);
        if (result.code === 0) {
            vscode.window.showInformationMessage('HTTP connection established');
        }
        else {
            vscode.window.showErrorMessage(`Failed to connect HTTP: ${result.stderr}`);
        }
    }
}
function activate(context) {
    console.log('Shield CLI VS Code extension activated');
    // 检查 shield-cli 是否安装
    checkShieldCli().then(installed => {
        if (!installed) {
            vscode.window.showWarningMessage('Shield CLI is not installed. Please install it first.');
        }
    });
    // 注册命令
    context.subscriptions.push(vscode.commands.registerCommand('shield.startService', startShieldService), vscode.commands.registerCommand('shield.stopService', stopShieldService), vscode.commands.registerCommand('shield.openWebUI', openShieldWebUI), vscode.commands.registerCommand('shield.connectSSH', connectSSH), vscode.commands.registerCommand('shield.connectRDP', connectRDP), vscode.commands.registerCommand('shield.connectHTTP', connectHTTP));
    // 创建侧边栏
    const sidebarProvider = new ShieldSidebarProvider(context.extensionUri);
    context.subscriptions.push(vscode.window.registerWebviewViewProvider('shield-sidebar', sidebarProvider));
}
function deactivate() {
    console.log('Shield CLI VS Code extension deactivated');
}
// 侧边栏提供器
class ShieldSidebarProvider {
    constructor(_extensionUri) {
        this._extensionUri = _extensionUri;
    }
    resolveWebviewView(webviewView, context, _token) {
        this._view = webviewView;
        webviewView.webview.options = {
            enableScripts: true,
            localResourceRoots: [this._extensionUri]
        };
        webviewView.webview.html = this._getHtmlForWebview(webviewView.webview);
        webviewView.webview.onDidReceiveMessage(async (data) => {
            switch (data.type) {
                case 'startService':
                    await startShieldService();
                    break;
                case 'stopService':
                    await stopShieldService();
                    break;
                case 'openWebUI':
                    openShieldWebUI();
                    break;
                case 'connectSSH':
                    await connectSSH();
                    break;
                case 'connectRDP':
                    await connectRDP();
                    break;
                case 'connectHTTP':
                    await connectHTTP();
                    break;
            }
        });
    }
    _getHtmlForWebview(webview) {
        const nonce = getNonce();
        return `
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Shield CLI</title>
        <style>
          body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            padding: 16px;
            background-color: #f5f5f5;
          }
          .container {
            max-width: 100%;
          }
          h1 {
            font-size: 18px;
            margin-bottom: 16px;
            color: #333;
          }
          .button {
            display: block;
            width: 100%;
            padding: 10px;
            margin-bottom: 10px;
            border: none;
            border-radius: 4px;
            background-color: #007acc;
            color: white;
            font-size: 14px;
            cursor: pointer;
          }
          .button:hover {
            background-color: #005a9e;
          }
          .button-secondary {
            background-color: #e0e0e0;
            color: #333;
          }
          .button-secondary:hover {
            background-color: #d0d0d0;
          }
          .section {
            margin-bottom: 20px;
            padding: 16px;
            background-color: white;
            border-radius: 4px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
          }
          h2 {
            font-size: 14px;
            margin-bottom: 12px;
            color: #666;
          }
        </style>
      </head>
      <body>
        <div class="container">
          <h1>Shield CLI</h1>
          
          <div class="section">
            <h2>Service Management</h2>
            <button class="button" onclick="vscode.postMessage({ type: 'startService' })">Start Service</button>
            <button class="button button-secondary" onclick="vscode.postMessage({ type: 'stopService' })">Stop Service</button>
            <button class="button button-secondary" onclick="vscode.postMessage({ type: 'openWebUI' })">Open Web UI</button>
          </div>
          
          <div class="section">
            <h2>Quick Connect</h2>
            <button class="button" onclick="vscode.postMessage({ type: 'connectSSH' })">Connect SSH</button>
            <button class="button" onclick="vscode.postMessage({ type: 'connectRDP' })">Connect RDP</button>
            <button class="button" onclick="vscode.postMessage({ type: 'connectHTTP' })">Connect HTTP</button>
          </div>
        </div>
        
        <script nonce="${nonce}">
          const vscode = acquireVsCodeApi();
        </script>
      </body>
      </html>
    `;
    }
}
// 生成随机 nonce
function getNonce() {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
}
//# sourceMappingURL=extension.js.map