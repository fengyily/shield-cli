#!/bin/sh
set -e

REPO="fengyily/shield-cli"
INSTALL_DIR="/usr/local/bin"
BINARY="shield"

# Parse arguments
USE_PKG_MANAGER=""
for arg in "$@"; do
    case "$arg" in
        --apt)  USE_PKG_MANAGER="apt" ;;
        --yum)  USE_PKG_MANAGER="yum" ;;
        --dnf)  USE_PKG_MANAGER="dnf" ;;
    esac
done

# If --apt/--yum/--dnf specified, delegate to setup-repo.sh
if [ -n "$USE_PKG_MANAGER" ]; then
    echo "Setting up Shield CLI via package manager (${USE_PKG_MANAGER})..."
    curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/scripts/setup-repo.sh" | sudo bash
    exit 0
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    i386|i686)     ARCH="386" ;;
    *)             echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest version"
    exit 1
fi

FILENAME="shield-${OS}-${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

echo "Downloading Shield CLI v${VERSION} for ${OS}/${ARCH}..."
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$URL" -o "${TMPDIR}/${FILENAME}"
tar xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

echo "Installing to ${INSTALL_DIR}/${BINARY}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
    sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi
chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "Shield CLI v${VERSION} installed successfully!"
echo ""
echo "Usage:"
echo "  shield start                # Web UI at http://localhost:8181"
echo "  shield ssh                  # 127.0.0.1:22"
echo "  shield ssh 2222             # 127.0.0.1:2222"
echo "  shield http 3000            # 127.0.0.1:3000"
echo "  shield --help               # More options"
echo ""
echo "Tip: On Linux, you can also install via package manager:"
echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh -s -- --apt"
echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh -s -- --yum"
