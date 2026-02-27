#!/bin/bash
set -e

# prevent-llm-delete installer
# Usage: curl -fsSL https://raw.githubusercontent.com/yourusername/prevent-llm-delete/main/install.sh | bash

VERSION="1.0.0"
REPO="yourusername/prevent-llm-delete"  # Update with your GitHub username
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="prevent-llm-delete"

echo "🔒 prevent-llm-delete installer v${VERSION}"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)
    OS="linux"
    ;;
  darwin)
    OS="darwin"
    ;;
  *)
    echo "❌ Unsupported operating system: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

PLATFORM="${OS}-${ARCH}"
echo "📦 Detected platform: ${PLATFORM}"

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/prevent-llm-delete-${PLATFORM}.tar.gz"

echo "⬇️  Downloading prevent-llm-delete..."

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download binary
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$DOWNLOAD_URL" -o prevent-llm-delete.tar.gz
elif command -v wget >/dev/null 2>&1; then
  wget -q "$DOWNLOAD_URL" -O prevent-llm-delete.tar.gz
else
  echo "❌ Neither curl nor wget found. Please install one of them."
  exit 1
fi

# Extract
echo "📂 Extracting..."
tar -xzf prevent-llm-delete.tar.gz

# Install
echo "📦 Installing to ${INSTALL_DIR}..."

if [ -w "$INSTALL_DIR" ]; then
  mv "prevent-llm-delete-${PLATFORM}" "${INSTALL_DIR}/${BINARY_NAME}"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
  echo "   (requires sudo)"
  sudo mv "prevent-llm-delete-${PLATFORM}" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
cd - >/dev/null
rm -rf "$TMP_DIR"

echo "✅ prevent-llm-delete installed successfully!"
echo ""
echo "🚀 Get started:"
echo "   prevent-llm-delete install    # Install the protection"
echo "   prevent-llm-delete status     # Check status"
echo "   prevent-llm-delete help       # Show help"
echo ""
echo "📚 Documentation: https://github.com/${REPO}"
