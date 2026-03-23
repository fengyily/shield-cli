#!/bin/bash
# Shield CLI - Linux Package Repository Setup Script
# Usage: curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash
set -e

REPO_BASE="https://fengyily.github.io/linux-repo"

detect_pkg_manager() {
    if command -v apt-get &>/dev/null; then
        echo "apt"
    elif command -v dnf &>/dev/null; then
        echo "dnf"
    elif command -v yum &>/dev/null; then
        echo "yum"
    else
        echo "unknown"
    fi
}

setup_apt() {
    echo "Setting up APT repository for Shield CLI..."

    # Add repository source
    cat > /etc/apt/sources.list.d/shield-cli.list <<EOF
deb [trusted=yes] ${REPO_BASE}/apt stable main
EOF

    # Update and install
    apt-get update -o Dir::Etc::sourcelist="sources.list.d/shield-cli.list" \
                   -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"
    apt-get install -y shield-cli

    echo ""
    echo "Shield CLI installed successfully!"
    echo "Run 'shield --help' to get started."
}

setup_yum() {
    local mgr="$1"
    echo "Setting up YUM/DNF repository for Shield CLI..."

    cat > /etc/yum.repos.d/shield-cli.repo <<EOF
[shield-cli]
name=Shield CLI Repository
baseurl=${REPO_BASE}/yum
enabled=1
gpgcheck=0
EOF

    ${mgr} install -y shield-cli

    echo ""
    echo "Shield CLI installed successfully!"
    echo "Run 'shield --help' to get started."
}

# --- Main ---
if [ "$(id -u)" -ne 0 ]; then
    echo "Error: This script must be run as root (use sudo)."
    exit 1
fi

PKG_MANAGER=$(detect_pkg_manager)

case "$PKG_MANAGER" in
    apt)
        setup_apt
        ;;
    dnf)
        setup_yum "dnf"
        ;;
    yum)
        setup_yum "yum"
        ;;
    *)
        echo "Error: Unsupported package manager."
        echo "Supported: apt (Debian/Ubuntu), yum/dnf (RHEL/CentOS/Fedora)"
        echo ""
        echo "You can install manually from: https://github.com/fengyily/shield-cli/releases"
        exit 1
        ;;
esac
