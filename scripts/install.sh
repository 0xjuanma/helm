#!/bin/sh
set -e

REPO="0xjuanma/helm"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

LATEST=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$LATEST" ]; then
    echo "Failed to fetch latest release"
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/helm-${OS}-${ARCH}"

echo "Downloading helm $LATEST for $OS/$ARCH..."
curl -sL "$URL" -o helm
chmod +x helm

if [ -w "$INSTALL_DIR" ]; then
    mv helm "$INSTALL_DIR/helm"
else
    sudo mv helm "$INSTALL_DIR/helm"
fi

echo "Installed helm to $INSTALL_DIR/helm"

