# Shield CLI - Windows Installer (PowerShell)
#
# Basic install (download binary only):
#   irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
#
# Install + register as service + open browser:
#   irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex; Install-ShieldService
#
# Or use install.bat for a double-click one-click experience.

$ErrorActionPreference = "Stop"

$Repo = "fengyily/shield-cli"
$Binary = "shield.exe"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else { "386" }

# Get latest version
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Version = $Release.tag_name -replace '^v', ''

$Filename = "shield-windows-${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/v${Version}/$Filename"

Write-Host ""
Write-Host "  Downloading Shield CLI v$Version for windows/$Arch..." -ForegroundColor Cyan
$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "shield-install-$(Get-Random)")
$ZipPath = Join-Path $TmpDir $Filename

try {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing
} catch {
    Write-Host "  GitHub download failed, trying mirror..." -ForegroundColor Yellow
    $MirrorUrl = "https://mirror.ghproxy.com/$Url"
    Invoke-WebRequest -Uri $MirrorUrl -OutFile $ZipPath -UseBasicParsing
}

Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install to Program Files if running as admin, otherwise user local
$IsAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($IsAdmin) {
    $InstallDir = Join-Path $env:ProgramFiles "ShieldCLI"
} else {
    $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\shield"
}

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Copy-Item (Join-Path $TmpDir $Binary) -Destination (Join-Path $InstallDir $Binary) -Force

# Add to PATH
if ($IsAdmin) {
    $SysPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($SysPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$SysPath;$InstallDir", "Machine")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "  Added $InstallDir to system PATH" -ForegroundColor Gray
    }
} else {
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "  Added $InstallDir to user PATH" -ForegroundColor Gray
    }
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host ""
Write-Host "  Shield CLI v$Version installed successfully!" -ForegroundColor Green
Write-Host "  Binary: $InstallDir\$Binary" -ForegroundColor Gray
Write-Host ""

# Helper function to also install as a service
function global:Install-ShieldService {
    param([int]$Port = 8181)

    $ShieldExe = Join-Path $InstallDir $Binary

    if (-not $IsAdmin) {
        Write-Host "  Service installation requires Administrator privileges." -ForegroundColor Red
        Write-Host "  Please run PowerShell as Administrator and try again." -ForegroundColor Yellow
        return
    }

    Write-Host "  Installing Shield as a Windows service (port $Port)..." -ForegroundColor Cyan

    # Remove existing service if present
    $existing = sc.exe query ShieldCLI 2>&1
    if ($LASTEXITCODE -eq 0) {
        sc.exe stop ShieldCLI 2>&1 | Out-Null
        Start-Sleep -Seconds 2
        sc.exe delete ShieldCLI 2>&1 | Out-Null
        Start-Sleep -Seconds 1
    }

    & $ShieldExe install --port $Port

    Write-Host ""
    Write-Host "  Shield Web UI is running at:" -ForegroundColor Green
    Write-Host "    http://localhost:$Port" -ForegroundColor Cyan
    Write-Host ""

    Start-Process "http://localhost:$Port"
}

Write-Host "  Quick start:" -ForegroundColor White
Write-Host "    shield start              # Web UI at http://localhost:8181"
Write-Host "    shield ssh 10.0.0.5       # SSH in browser"
Write-Host "    shield --help             # More options"
Write-Host ""

if ($IsAdmin) {
    Write-Host "  To also install as a service, run:" -ForegroundColor Yellow
    Write-Host "    Install-ShieldService" -ForegroundColor White
    Write-Host ""
}
