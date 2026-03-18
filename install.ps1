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

Write-Host "Downloading Shield CLI v$Version for windows/$Arch..."
$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "shield-install-$(Get-Random)")
$ZipPath = Join-Path $TmpDir $Filename

Invoke-WebRequest -Uri $Url -OutFile $ZipPath
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install to user's local bin
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\shield"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Copy-Item (Join-Path $TmpDir $Binary) -Destination (Join-Path $InstallDir $Binary) -Force

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path = "$env:Path;$InstallDir"
    Write-Host "Added $InstallDir to PATH"
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host ""
Write-Host "Shield CLI v$Version installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage:"
Write-Host "  shield ssh                  # 127.0.0.1:22"
Write-Host "  shield ssh 2222             # 127.0.0.1:2222"
Write-Host "  shield http 3000            # 127.0.0.1:3000"
Write-Host "  shield --help               # More options"
Write-Host ""
Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
