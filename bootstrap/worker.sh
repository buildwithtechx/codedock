#!/usr/bin/env bash
set -eo pipefail

if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ] && [ "${BASH_SOURCE[0]}" != "-" ]; then
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  if [ -f "$SCRIPT_DIR/base.sh" ]; then
    source "$SCRIPT_DIR/base.sh"
  fi
fi

if [ -z "${BOLD:-}" ] && [ -f "/codedock/bootstrap/base.sh" ]; then
  source "/codedock/bootstrap/base.sh"
fi

if [ -z "${BOLD:-}" ]; then
  BOLD="\033[1m"; DIM="\033[2m"; GREEN="\033[0;32m"; YELLOW="\033[0;33m"; RED="\033[0;31m"; NC="\033[0m"
  ensure_root() { [ "$EUID" -eq 0 ] || { echo -e "${RED}❌ Run as root.${NC}"; exit 1; }; }
  ensure_docker() {
    if ! command -v docker &>/dev/null; then
      if command -v apt-get &>/dev/null; then apt-get update -qq && apt-get install -y -qq docker.io 2>/dev/null || true; fi
      if ! command -v docker &>/dev/null; then curl -fsSL https://get.docker.com | sh; fi
      systemctl enable --now docker 2>/dev/null || true
    fi
  }
  setup_systemd_service() { [ -d /etc/systemd/system ] && echo "$2" > "/etc/systemd/system/$1.service" && systemctl daemon-reload && systemctl enable --now "$1.service"; }
fi

CODEDOCK_DIR="/opt/codedock"
DOWNLOAD_URL="https://github.com/buildwithtechx/codedock/releases/latest/download"
ENV_FILE="/etc/codedockw.env"

LICENSE_KEY=""
API_HOST="wss://api.codedock.dev"

while [[ $# -gt 0 ]]; do
  case $1 in
    --key) LICENSE_KEY="$2"; shift 2 ;;
    --host) API_HOST="$2"; shift 2 ;;
    *) echo -e "${RED}Unknown argument: $1${NC}"; exit 1 ;;
  esac
done

if [ -z "$LICENSE_KEY" ]; then
  echo -e "${RED}❌ Missing required argument: --key <LICENSE_KEY>${NC}"
  exit 1
fi

echo -e "${BOLD}🛰️  Installing Codedock Worker Daemon${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

ensure_root
ensure_docker

echo -e "${BOLD}⬇️  Downloading codedockw...${NC}"
mkdir -p "$CODEDOCK_DIR/bin"
ARCH=$(uname -m)
case $ARCH in
  x86_64) GOARCH="amd64" ;;
  aarch64|arm64) GOARCH="arm64" ;;
  *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

TAR_FILE="codedockw_linux_${GOARCH}.tar.gz"
TMP_TAR=$(mktemp)

curl -fsSL "${DOWNLOAD_URL}/${TAR_FILE}" -o "$TMP_TAR" || {
  echo -e "${RED}Failed to download worker binary.${NC}"
  rm -f "$TMP_TAR"
  exit 1
}

if [ -n "${CODEDOCK_WORKER_SHA256:-}" ]; then
  ACTUAL_SHA=$(sha256sum "$TMP_TAR" | awk '{print $1}')
  if [ "$ACTUAL_SHA" != "$CODEDOCK_WORKER_SHA256" ]; then
    echo -e "${RED}❌ Checksum mismatch for worker binary! Expected: ${CODEDOCK_WORKER_SHA256}, Got: ${ACTUAL_SHA}${NC}"
    rm -f "$TMP_TAR"
    exit 1
  fi
fi

tar --no-same-owner -xzf "$TMP_TAR" -C "$CODEDOCK_DIR/bin" codedockw
chmod +x "$CODEDOCK_DIR/bin/codedockw"
rm -f "$TMP_TAR"

echo -e "${BOLD}🔑 Writing credentials to ${ENV_FILE}...${NC}"
install -m 0600 /dev/null "$ENV_FILE"
cat > "$ENV_FILE" <<EOF
CODEDOCK_WORKER_TOKEN=${LICENSE_KEY}
CODEDOCK_CONTROL_PLANE=${API_HOST}
CODEDOCK_API_HOST=${API_HOST}
EOF
echo -e "  ${GREEN}✅ Credentials written.${NC}"

setup_systemd_service "codedockw" "[Unit]
Description=Codedock Worker Daemon
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
ExecStart=${CODEDOCK_DIR}/bin/codedockw
Restart=always
RestartSec=10
EnvironmentFile=${ENV_FILE}

[Install]
WantedBy=multi-user.target"

echo ""
echo -e "${GREEN}✅ Worker daemon installed and started!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
