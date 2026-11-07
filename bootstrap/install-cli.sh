#!/usr/bin/env bash
set -eo pipefail

if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ] && [ "${BASH_SOURCE[0]}" != "-" ]; then
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  if [ -f "$SCRIPT_DIR/base.sh" ]; then
    source "$SCRIPT_DIR/base.sh"
  fi
fi

if ! declare -f detect_platform &>/dev/null && [ -f "/codedock/bootstrap/base.sh" ]; then
  source "/codedock/bootstrap/base.sh"
fi

if ! declare -f detect_platform &>/dev/null; then
  BOLD="\033[1m"; DIM="\033[2m"; GREEN="\033[0;32m"; YELLOW="\033[0;33m"; RED="\033[0;31m"; NC="\033[0m"
  detect_platform() {
    local os arch
    os="$(uname -s)"
    arch="$(uname -m)"
    case "$os" in
      Linux) os="linux" ;;
      Darwin) os="darwin" ;;
      *) echo -e "${RED}❌ Unsupported OS: $os. Install via npm (npm i -g codedock) or download a release instead.${NC}"; exit 1 ;;
    esac
    case "$arch" in
      x86_64|amd64) arch="amd64" ;;
      arm64|aarch64) arch="arm64" ;;
      *) echo -e "${RED}❌ Unsupported architecture: $arch${NC}"; exit 1 ;;
    esac
    echo "${os}_${arch}"
  }
fi

REPO="buildwithtechx/codedock"
BINARY="codedock"
INSTALL_DIR="/usr/local/bin"

install_via_go() {
  if command -v go &>/dev/null && command -v git &>/dev/null; then
    echo -e "${YELLOW}⚙️  Installing from source via Go...${NC}"
    local src_dir
    src_dir="$(mktemp -d)"
    if ! git clone --depth 1 "https://github.com/${REPO}.git" "$src_dir" &>/dev/null; then
      rm -rf "$src_dir"
      echo -e "${RED}❌ Source checkout failed.${NC}"
      return 1
    fi
    (cd "$src_dir" && go install ./cmd/codedock)
    rm -rf "$src_dir"
    local gobin
    gobin="$(go env GOBIN)"
    if [ -z "$gobin" ]; then
      gobin="$(go env GOPATH)/bin"
    fi
    local src_bin="$gobin/codedock"
    if [ ! -f "$src_bin" ]; then
      echo -e "${RED}❌ 'go install' executed but binary was not found at $src_bin.${NC}"
      return 1
    fi
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
      cp "$src_bin" "$INSTALL_DIR/$BINARY"
      echo -e "${GREEN}✅ Installed → $INSTALL_DIR/$BINARY${NC}"
    else
      LOCAL_BIN="$HOME/.local/bin"
      mkdir -p "$LOCAL_BIN"
      cp "$src_bin" "$LOCAL_BIN/$BINARY"
      echo -e "${GREEN}✅ Installed → $LOCAL_BIN/$BINARY${NC}"
      echo -e "   Ensure $LOCAL_BIN or $gobin is in your PATH."
    fi
    return 0
  fi
  return 1
}

echo -e "${BOLD}🛰️  codedock CLI Installer${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

PLATFORM=$(detect_platform)
echo -e "  Platform:  ${PLATFORM}"

TARGET_VERSION="${CODEDOCK_VERSION:-}"
if [ -z "$TARGET_VERSION" ]; then
  echo -e "  Version:   ${DIM}checking latest...${NC}"
  TARGET_VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | head -1 | cut -d'"' -f4 || echo "")
fi

if [ -z "$TARGET_VERSION" ]; then
  install_via_go || {
    echo -e "${RED}❌ Could not fetch latest release. Install Go or set CODEDOCK_VERSION.${NC}"
    exit 1
  }
  exit 0
fi

echo -e "  Version:   ${TARGET_VERSION}"
echo ""

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TARGET_VERSION}/codedock_${PLATFORM}.tar.gz"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo -e "${BOLD}⬇️  Downloading codedock ${TARGET_VERSION} (${PLATFORM})...${NC}"

if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/codedock.tar.gz" 2>/dev/null; then
  install_via_go || {
    echo -e "${RED}❌ Download failed and source install needs Go and git.${NC}"
    exit 1
  }
  exit 0
fi

tar -xzf "$TMP_DIR/codedock.tar.gz" -C "$TMP_DIR"
BINARY_PATH="$TMP_DIR/codedock"
if [ ! -f "$BINARY_PATH" ]; then
  BINARY_PATH=$(find "$TMP_DIR" -name "codedock" -type f | head -1)
fi

if [ -z "$BINARY_PATH" ]; then
  echo -e "${RED}❌ Could not find codedock binary in archive.${NC}"
  exit 1
fi

chmod +x "$BINARY_PATH"

if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
  mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY"
  echo -e "${GREEN}✅ Installed → $INSTALL_DIR/$BINARY${NC}"
else
  LOCAL_BIN="$HOME/.local/bin"
  mkdir -p "$LOCAL_BIN"
  mv "$BINARY_PATH" "$LOCAL_BIN/$BINARY"
  echo -e "${GREEN}✅ Installed → $LOCAL_BIN/$BINARY${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ codedock ${TARGET_VERSION} installed successfully!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
