#!/usr/bin/env bash
# Codedock Worker Daemon Installer
# Usage: curl -sL https://get.codedock.dev | bash -s -- --key <LICENSE_KEY>

set -eo pipefail

CODEDOCK_DIR="/opt/codedock"
DOWNLOAD_URL="https://github.com/buildwithtechx/codedock/releases/latest/download"
ENV_FILE="/etc/codedockw.env"

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

verify_checksum() {
  local file="$1"
  local expected="$2"
  local label="$3"
  local var_name="$4"

  if [ -z "$expected" ]; then
    echo -e "  ${YELLOW}⚠️  No checksum provided for ${label} (${var_name:-checksum}). Skipping checksum check.${NC}"
    return 0
  fi

  if ! command -v sha256sum &>/dev/null; then
    echo -e "${RED}❌ sha256sum is not installed. Cannot verify ${label}.${NC}"
    rm -f "$file"
    exit 1
  fi

  local actual
  actual=$(sha256sum "$file" | awk '{print $1}')
  if [ "$actual" != "$expected" ]; then
    echo -e "${RED}❌ Checksum mismatch for ${label}!${NC}"
    echo -e "${RED}   Expected: ${expected}${NC}"
    echo -e "${RED}   Got:      ${actual}${NC}"
    rm -f "$file"
    exit 1
  fi
  echo -e "  ${GREEN}✅ Checksum verified for ${label}${NC}"
}

install_docker() {
  if ! command -v docker &> /dev/null; then
    echo -e "${BOLD}📦 Installing Docker...${NC}"
    if command -v apt-get &> /dev/null; then
      apt-get update -qq && apt-get install -y -qq docker.io docker-compose-plugin 2>/dev/null || true
    elif command -v dnf &> /dev/null; then
      dnf install -y -q docker docker-compose-plugin 2>/dev/null || true
    elif command -v yum &> /dev/null; then
      yum install -y -q docker 2>/dev/null || true
    fi

    if ! command -v docker &> /dev/null; then
      DOCKER_SCRIPT=$(mktemp)
      curl -fsSL https://get.docker.com -o "$DOCKER_SCRIPT"
      verify_checksum "$DOCKER_SCRIPT" "${DOCKER_INSTALLER_SHA256:-}" "Docker installer" "DOCKER_INSTALLER_SHA256"
      sh "$DOCKER_SCRIPT"
      rm -f "$DOCKER_SCRIPT"
    fi
    systemctl enable --now docker 2>/dev/null || true
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
  TMP_TAR=$(mktemp)

  curl -fsSL "${DOWNLOAD_URL}/${TAR_FILE}" -o "$TMP_TAR" || {
    echo -e "${RED}Failed to download worker binary.${NC}"
    rm -f "$TMP_TAR"
    exit 1
  }

  verify_checksum "$TMP_TAR" "${CODEDOCK_WORKER_SHA256:-}" "codedockw" "CODEDOCK_WORKER_SHA256"
  tar --no-same-owner -xzf "$TMP_TAR" -C "$CODEDOCK_DIR/bin" codedockw
  chmod +x "$CODEDOCK_DIR/bin/codedockw"
  rm -f "$TMP_TAR"
}

write_env_file() {
  echo -e "${BOLD}🔑 Writing credentials to ${ENV_FILE}...${NC}"
  install -m 0600 /dev/null "$ENV_FILE"
  cat > "$ENV_FILE" <<EOF
CODEDOCK_WORKER_TOKEN=${LICENSE_KEY}
CODEDOCK_CONTROL_PLANE=${API_HOST}
CODEDOCK_API_HOST=${API_HOST}
EOF
  echo -e "  ${GREEN}✅ Credentials written (readable by root only).${NC}"
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
ExecStart=${CODEDOCK_DIR}/bin/codedockw
Restart=always
RestartSec=10
EnvironmentFile=${ENV_FILE}

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
  write_env_file
  setup_systemd

  echo ""
  echo -e "${GREEN}✅ Worker daemon installed and started!${NC}"
  echo -e "   It will automatically connect to your control plane."
  echo -e "   Run ${BOLD}journalctl -u codedockw -f${NC} to view logs."
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

main "$@"
