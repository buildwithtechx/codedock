#!/usr/bin/env bash
# Codedock Worker Daemon Installer
# Usage: curl -sL https://get.codedock.dev | bash -s -- --key <LICENSE_KEY>

set -eo pipefail

RELEASE="latest"
CODEDOCK_DIR="/opt/codedock"
DOWNLOAD_URL="https://github.com/buildwithtechx/codedock/releases/latest/download"

setup_colors() {
  BOLD="\033[1m"
  GREEN="\033[0;32m"
  YELLOW="\033[0;33m"
  RED="\033[0;31m"
  NC="\033[0m"
}

check_root() {
  if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Please run as root (or with sudo).${NC}"
    exit 1
  fi
}

install_docker() {
  if ! command -v docker &> /dev/null; then
    echo -e "${BOLD}📦 Installing Docker...${NC}"
    curl -fsSL https://get.docker.com | sh
    systemctl enable --now docker
  fi
}

parse_args() {
  while [[ $# -gt 0 ]]; do
    case $1 in
      --key)
        LICENSE_KEY="$2"
        shift 2
        ;;
      --host)
        API_HOST="$2"
        shift 2
        ;;
      *)
        echo -e "${RED}Unknown argument: $1${NC}"
        exit 1
        ;;
    esac
  done

  if [ -z "$LICENSE_KEY" ]; then
    echo -e "${RED}❌ Missing required argument: --key <LICENSE_KEY>${NC}"
    exit 1
  fi
  
  if [ -z "$API_HOST" ]; then
    API_HOST="wss://api.codedock.dev"
  fi
}

download_worker() {
  echo -e "${BOLD}⬇️  Downloading codedockw...${NC}"
  mkdir -p "$CODEDOCK_DIR/bin"
  
  ARCH=$(uname -m)
  case $ARCH in
    x86_64) GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
  esac

  TAR_FILE="codedockw_linux_${GOARCH}.tar.gz"
  
  curl -fsSL "${DOWNLOAD_URL}/${TAR_FILE}" -o "/tmp/${TAR_FILE}" || {
    echo -e "${RED}Failed to download worker binary.${NC}"
    exit 1
  }
  
  tar -xzf "/tmp/${TAR_FILE}" -C "$CODEDOCK_DIR/bin"
  chmod +x "$CODEDOCK_DIR/bin/codedockw"
  rm "/tmp/${TAR_FILE}"
}

setup_systemd() {
  echo -e "${BOLD}⚙️  Setting up systemd service...${NC}"
  
  cat > /etc/systemd/system/codedockw.service <<EOF
[Unit]
Description=Codedock Worker Daemon
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
ExecStart=$CODEDOCK_DIR/bin/codedockw
Restart=always
RestartSec=10
Environment="CODEDOCK_WORKER_TOKEN=$LICENSE_KEY"
Environment="CODEDOCK_CONTROL_PLANE=$API_HOST"

[Install]
WantedBy=multi-user.target
EOF

  systemctl daemon-reload
  systemctl enable --now codedockw.service
}

main() {
  setup_colors
  parse_args "$@"
  
  echo -e "${BOLD}🛰️  Installing Codedock Worker Daemon${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  check_root
  install_docker
  download_worker
  setup_systemd
  
  echo ""
  echo -e "${GREEN}✅ Worker daemon installed and started!${NC}"
  echo -e "   It will automatically connect to your control plane."
  echo -e "   Run ${BOLD}journalctl -u codedockw -f${NC} to view logs."
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

main "$@"
