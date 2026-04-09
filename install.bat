@echo off
chcp 65001 >nul 2>&1
setlocal EnableDelayedExpansion

:: ================================================
::  Shield CLI - Windows One-Click Installer
::  Double-click to install Shield as a service
::  and open the Web UI in your browser.
:: ================================================

title Shield CLI Installer

:: ----------------------------
:: Step 1: Request Admin rights
:: ----------------------------
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [Shield CLI] Requesting administrator privileges...
    powershell -Command "Start-Process cmd -ArgumentList '/c \"\"%~f0\"\"' -Verb RunAs"
    exit /b
)

:: ----------------------------
:: Step 2: Configuration
:: ----------------------------
set "REPO=fengyily/shield-cli"
set "INSTALL_DIR=%ProgramFiles%\ShieldCLI"
set "BINARY=shield.exe"
set "PORT=8181"

echo.
echo   ============================================
echo     Shield CLI - One-Click Installer
echo   ============================================
echo.

:: ----------------------------
:: Step 3: Detect Architecture
:: ----------------------------
set "ARCH=amd64"
if /i "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"
if /i "%PROCESSOR_ARCHITECTURE%"=="x86" (
    if /i "%PROCESSOR_ARCHITEW6432%"=="" (
        echo   [ERROR] 32-bit Windows is not supported.
        echo   Please use a 64-bit version of Windows.
        echo.
        pause
        exit /b 1
    )
)
:: WoW64: 32-bit cmd on 64-bit OS
if /i "%PROCESSOR_ARCHITEW6432%"=="AMD64" set "ARCH=amd64"
if /i "%PROCESSOR_ARCHITEW6432%"=="ARM64" set "ARCH=arm64"

echo   Architecture: windows/%ARCH%
echo.

:: ----------------------------
:: Step 4: Get latest version
:: ----------------------------
echo   [1/5] Fetching latest version...

for /f "delims=" %%V in ('powershell -NoProfile -Command "(Invoke-RestMethod -Uri 'https://api.github.com/repos/%REPO%/releases/latest').tag_name -replace '^v',''"') do set "VERSION=%%V"

if "%VERSION%"=="" (
    echo.
    echo   [ERROR] Failed to fetch latest version.
    echo   Please check your network connection.
    echo.
    pause
    exit /b 1
)

echo   Latest version: v%VERSION%

:: ----------------------------
:: Step 5: Download
:: ----------------------------
set "FILENAME=shield-windows-%ARCH%.zip"
set "URL=https://github.com/%REPO%/releases/download/v%VERSION%/%FILENAME%"
set "TMPDIR=%TEMP%\shield-install-%RANDOM%"

mkdir "%TMPDIR%" >nul 2>&1

echo   [2/5] Downloading Shield CLI v%VERSION%...

powershell -NoProfile -Command ^
    "$ProgressPreference='SilentlyContinue'; " ^
    "try { Invoke-WebRequest -Uri '%URL%' -OutFile '%TMPDIR%\%FILENAME%' -UseBasicParsing } " ^
    "catch { " ^
    "  Write-Host ''; " ^
    "  Write-Host '  [INFO] GitHub download slow? Trying China mirror...'; " ^
    "  $mirrorUrl = 'https://mirror.ghproxy.com/%URL%'; " ^
    "  try { Invoke-WebRequest -Uri $mirrorUrl -OutFile '%TMPDIR%\%FILENAME%' -UseBasicParsing } " ^
    "  catch { Write-Host '  [ERROR] Download failed.'; exit 1 } " ^
    "}"

if not exist "%TMPDIR%\%FILENAME%" (
    echo.
    echo   [ERROR] Download failed.
    echo   URL: %URL%
    echo.
    rd /s /q "%TMPDIR%" >nul 2>&1
    pause
    exit /b 1
)

:: ----------------------------
:: Step 6: Extract
:: ----------------------------
echo   [3/5] Extracting...

powershell -NoProfile -Command "Expand-Archive -Path '%TMPDIR%\%FILENAME%' -DestinationPath '%TMPDIR%\extracted' -Force"

:: ----------------------------
:: Step 7: Install binary
:: ----------------------------
echo   [4/5] Installing to %INSTALL_DIR%...

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /y "%TMPDIR%\extracted\%BINARY%" "%INSTALL_DIR%\%BINARY%" >nul

:: Add to system PATH
powershell -NoProfile -Command ^
    "$path = [Environment]::GetEnvironmentVariable('Path','Machine'); " ^
    "if ($path -notlike '*%INSTALL_DIR%*') { " ^
    "  [Environment]::SetEnvironmentVariable('Path', \"$path;%INSTALL_DIR%\", 'Machine'); " ^
    "  Write-Host '  Added to system PATH' " ^
    "}"

:: Cleanup temp files
rd /s /q "%TMPDIR%" >nul 2>&1

:: ----------------------------
:: Step 8: Install as service
:: ----------------------------
echo   [5/5] Installing system service...

:: Stop and remove existing service if present
sc query ShieldCLI >nul 2>&1 && (
    sc stop ShieldCLI >nul 2>&1
    timeout /t 2 /nobreak >nul
    sc delete ShieldCLI >nul 2>&1
    timeout /t 1 /nobreak >nul
)

"%INSTALL_DIR%\%BINARY%" install --port %PORT%

if %errorLevel% neq 0 (
    echo.
    echo   [WARNING] Service installation had issues.
    echo   You can manually run: shield install --port %PORT%
    echo.
)

:: ----------------------------
:: Step 9: Open browser
:: ----------------------------
echo.
echo   ============================================
echo     Shield CLI v%VERSION% installed!
echo   ============================================
echo.
echo     Web UI: http://localhost:%PORT%
echo     Binary: %INSTALL_DIR%\%BINARY%
echo.
echo     Commands:
echo       shield start          - Start Web UI
echo       shield install        - Install as service
echo       shield uninstall      - Remove service
echo       shield --help         - More options
echo.

:: Open Web UI in default browser
start "" "http://localhost:%PORT%"

echo   Press any key to close this window...
pause >nul
