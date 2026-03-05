#!/bin/sh
set -e

REPO="DiyRex/outline-cli"
BINARY_NAME="outline"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

info() { printf "${CYAN}▸${NC} %s\n" "$1"; }
success() { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn() { printf "${YELLOW}!${NC} %s\n" "$1"; }
error() { printf "${RED}✗${NC} %s\n" "$1" >&2; exit 1; }

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)      error "Unsupported OS: $OS" ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)   ARCH="amd64" ;;
    aarch64|arm64)   ARCH="arm64" ;;
    *)               error "Unsupported architecture: $ARCH" ;;
esac

info "Detected platform: ${BOLD}${OS}/${ARCH}${NC}"

# Get latest release tag
info "Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
    error "Could not determine latest release. Check https://github.com/${REPO}/releases"
fi

info "Latest version: ${BOLD}${LATEST}${NC}"

# Download
BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY_FILE}"

info "Downloading ${BOLD}${BINARY_FILE}${NC}..."
TMP_DIR=$(mktemp -d)
curl -fsSL -o "${TMP_DIR}/${BINARY_NAME}" "$DOWNLOAD_URL" || error "Download failed. Check if the release exists at:\n  ${DOWNLOAD_URL}"

chmod +x "${TMP_DIR}/${BINARY_NAME}"

# Install
info "Installing to ${BOLD}${INSTALL_DIR}/${BINARY_NAME}${NC}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

rm -rf "$TMP_DIR"

# Verify
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    success "outline installed successfully!"
    printf "\n"
    printf "  ${BOLD}Get started:${NC}\n"
    printf "    outline config --url=\"https://your-instance.com\"\n"
    printf "    outline config --api-key=\"ol_api_your_key\"\n"
    printf "    outline status\n"
    printf "\n"
else
    warn "Installed but '${BINARY_NAME}' not found in PATH. Add ${INSTALL_DIR} to your PATH."
fi
