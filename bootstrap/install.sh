#!/usr/bin/env bash
# Codedock 1-Click Installer
# Usage: curl -fsSL https://get.codedock.run | sh
set -eo pipefail

RELEASE=${CODEDOCK_VERSION:-1.0.0}
CODEDOCK_DIR=/codedock
REPO_URL="https://raw.githubusercontent.com/buildwithtechx/codedock/main"
COMPOSE_URL="$REPO_URL/docker-compose.yml"
CTL_URL="$REPO_URL/bootstrap/codedockd"

COMPOSE_SHA256="${CODEDOCK_COMPOSE_SHA256:-}"
CTL_SHA256="${CODEDOCK_CTL_SHA256:-}"

setup_colors() {
  BOLD="\033[1m"
  DIM="\033[2m"
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

check_system() {
  echo -e "${BOLD}🔍 System check${NC}"
  TOTAL_RAM=$(free -m | awk '/^Mem:/{print $2}')
  if [ "$TOTAL_RAM" -lt 1024 ]; then
    echo -e "  ${YELLOW}⚠️  ${TOTAL_RAM}MB RAM detected. 1GB+ recommended for Docker builds.${NC}"
    SWAP=$(free -m | awk '/^Swap:/{print $2}')
    if [ "$SWAP" -eq 0 ]; then
      echo -e "  ${YELLOW}   No swap configured. Docker builds may fail on low memory.${NC}"
      echo -e "  ${YELLOW}   Consider adding swap:${NC}"
      echo -e "  ${DIM}    fallocate -l 2G /swapfile && chmod 600 /swapfile && mkswap /swapfile && swapon /swapfile${NC}"
    else
      echo -e "  ${GREEN}✅ ${SWAP}MB swap detected.${NC}"
    fi
  else
    echo -e "  ${GREEN}✅ ${TOTAL_RAM}MB RAM${NC}"
  fi
}

check_ports() {
  echo -e "${BOLD}🔌 Port check${NC}"
  for PORT in 80 443 8080; do
    if ss -tlnp "sport = :$PORT" 2>/dev/null | grep -q .; then
      echo -e "  ${YELLOW}⚠️  Port $PORT is already in use. Codedock needs it.${NC}"
    else
      echo -e "  ${GREEN}✅ Port $PORT available${NC}"
    fi
  done
}

detect_server_ip() {
  SERVER_IP=$(curl -4fsS ifconfig.me 2>/dev/null || echo "your-server-ip")
  echo -e "  ${GREEN}✅ Server IP: ${SERVER_IP}${NC}"
  echo ""
}

verify_checksum() {
  local file="$1"
  local expected="$2"
  local label="$3"

  if [ -z "$expected" ]; then
    echo -e "${RED}❌ No checksum provided for ${label}. Set ${label^^}_SHA256 to continue.${NC}"
    echo -e "${YELLOW}   Tip: compute it with: sha256sum <file>${NC}"
    rm -f "$file"
    exit 1
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
      verify_checksum "$DOCKER_SCRIPT" "${DOCKER_INSTALLER_SHA256:-}" "Docker installer"
      sh "$DOCKER_SCRIPT"
      rm -f "$DOCKER_SCRIPT"
    fi
    systemctl enable --now docker 2>/dev/null || true
  fi

  if ! docker info &> /dev/null; then
    echo -e "${YELLOW}⏳ Waiting for Docker...${NC}"
    sleep 3
  fi

  if ! docker compose version &> /dev/null; then
    echo -e "${YELLOW}📦 Installing Docker Compose plugin...${NC}"
    apt-get update -qq && apt-get install -y -qq docker-compose-plugin 2>/dev/null || \
      echo -e "${YELLOW}⚠️  Could not install compose plugin. Run: apt-get install docker-compose-plugin${NC}"
  fi

  echo -e "${GREEN}✅ Docker ready${NC}"
}

setup_directories() {
  mkdir -p "$CODEDOCK_DIR"/data/{backups,traefik,builds,storage}
}

fetch_config_files() {
  echo -e "${BOLD}⬇️  Fetching configuration files...${NC}"

  local compose_tmp ctl_tmp
  compose_tmp=$(mktemp)
  ctl_tmp=$(mktemp)

  curl -fsSL "$COMPOSE_URL" -o "$compose_tmp"
  verify_checksum "$compose_tmp" "$COMPOSE_SHA256" "docker-compose.yml"
  mv "$compose_tmp" "$CODEDOCK_DIR/docker-compose.yml"

  curl -fsSL "$CTL_URL" -o "$ctl_tmp"
  verify_checksum "$ctl_tmp" "$CTL_SHA256" "codedockd"
  mv "$ctl_tmp" "$CODEDOCK_DIR/codedockd"
  chmod +x "$CODEDOCK_DIR/codedockd"

  ln -sf "$CODEDOCK_DIR/codedockd" /usr/local/bin/codedockd
  ln -sf "$CODEDOCK_DIR/codedockd" /usr/local/bin/codedockctl
}

generate_env_file() {
  if [ ! -f "$CODEDOCK_DIR/.env" ]; then
    echo -e "${BOLD}🔑 Generating .env file...${NC}"

    JWT_SECRET=$(head -c 32 /dev/urandom | base64 | tr -dc 'a-zA-Z0-9' | head -c 32)
    REFRESH_SECRET=$(head -c 32 /dev/urandom | base64 | tr -dc 'a-zA-Z0-9' | head -c 32)

    TLS_EMAIL=""
    WILDCARD_DOMAIN=""
    DASHBOARD_URL=""

    install -m 0600 /dev/null "$CODEDOCK_DIR/.env"

    cat > "$CODEDOCK_DIR/.env" <<ENV
# Codedock Environment Configuration
# Auto-generated by installer. Edit if needed, then restart: docker compose restart

# Core
PORT=8080
HOST=0.0.0.0
CODEDOCK_DATA_DIR=/codedock/data
CODEDOCK_HOST_IP="${SERVER_IP}"

# Security (required - change this in production!)
CODEDOCK_JWT_SECRET="${JWT_SECRET}"
CODEDOCK_REFRESH_SECRET="${REFRESH_SECRET}"

# Let's Encrypt SSL (for automatic HTTPS on custom domains)
CODEDOCK_TLS_EMAIL="${TLS_EMAIL}"

# Wildcard domain — apps get https://myapp.yourdomain.com
# Leave empty — apps get myapp.IP.sslip.io (no DNS needed)
CODEDOCK_WILDCARD_DOMAIN="${WILDCARD_DOMAIN}"

# Magic Domain for local/IP routing (options: sslip.io, traefik.me, nip.io)
CODEDOCK_MAGIC_DOMAIN="sslip.io"

# Docker
DOCKER_SOCKET_PATH=/var/run/docker.sock

# Dashboard URL — used in notification links and redirects
CODEDOCK_DASHBOARD_URL="${DASHBOARD_URL:-http://${SERVER_IP}:8080}"

# Docker Deployments
CODEDOCK_RUNTIME_NETWORK=codedock-network
DEPLOY_HOST_PORT_START=4100
DEPLOY_HOST_PORT_END=4999
DEPLOY_DRY_RUN=false

# Updates
CODEDOCK_UPDATE_URL=https://api.github.com/repos/buildwithtechx/codedock/releases/latest
CODEDOCK_DOWNLOAD_URL=https://github.com/buildwithtechx/codedock/releases
ENV
    chmod 0600 "$CODEDOCK_DIR/.env"
    echo -e "  ${GREEN}✅ .env written with 0600 permissions (readable by root only).${NC}"
  fi
}

start_codedock() {
  echo -e "${BOLD}🐳 Pulling codedock:${RELEASE}...${NC}"
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" pull
  echo -e "${BOLD}🚀 Starting Codedock...${NC}"
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" up -d
}

check_health() {
  echo -e "${BOLD}⏳ Waiting for Codedock to become healthy...${NC}"
  for i in $(seq 1 30); do
    if curl -sfS "http://localhost:8080/healthz" > /dev/null 2>&1; then
      echo -e "${GREEN}✅ Codedock is healthy and accepting requests.${NC}"
      break
    fi
    if [ "$i" -eq 30 ]; then
      echo -e "${YELLOW}⚠️  Codedock is starting but not yet responding. Run 'codedockd logs -f' to check.${NC}"
    fi
    sleep 2
  done
}

setup_systemd() {
  if command -v systemctl &> /dev/null; then
    cat > /etc/systemd/system/codedock.service <<'SERVICE'
[Unit]
Description=Codedock – Self-hosted PaaS
After=docker.service
Requires=docker.service

[Service]
Restart=always
RestartSec=10
WorkingDirectory=/codedock
ExecStart=/usr/bin/docker compose -f /codedock/docker-compose.yml up codedock
ExecStop=/usr/bin/docker compose -f /codedock/docker-compose.yml stop codedock

[Install]
WantedBy=multi-user.target
SERVICE
    systemctl daemon-reload
    systemctl enable --now codedock.service
  fi
}

print_summary() {
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo -e "${GREEN}✅  Codedock v${RELEASE} installed successfully!${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""
  echo -e "  ${BOLD}📍 Dashboard:${NC}  http://${SERVER_IP}:8080"
  echo -e "  ${BOLD}📖 Docs:${NC}       https://docs.codedock.run"
  echo -e "  ${BOLD}🛠️  Server CLI:${NC}  codedockd --help (runs commands inside the container)"
  echo -e "  ${DIM}Install the remote CLI on your laptop: curl -fsSL https://get.codedock.run/cli | sh${NC}"
  echo ""

  echo -e "  ${BOLD}📍 Dashboard:${NC}  ${DASHBOARD_URL:-http://${SERVER_IP}:8080}"
  if [ -n "$WILDCARD_DOMAIN" ]; then
    echo -e "  ${BOLD}🌐 Apps:${NC}       https://myapp.${WILDCARD_DOMAIN}"
    echo ""
    echo -e "  ${YELLOW}⚠️  DNS: *.${WILDCARD_DOMAIN}  A  ${SERVER_IP}${NC}"
  else
    echo -e "  ${BOLD}🌐 Apps:${NC}       myapp.${SERVER_IP}.sslip.io (no DNS needed)"
  fi

  if [ -n "$TLS_EMAIL" ]; then
    echo -e "  ${GREEN}🔒 SSL:       Enabled for ${TLS_EMAIL}${NC}"
  else
    echo -e "  ${YELLOW}🔒 SSL:       Not configured${NC} (add CODEDOCK_TLS_EMAIL in .env later)"
  fi

  TOTAL_RAM_MB=$(free -m | awk '/^Mem:/{print $2}')
  echo -e "  ${DIM}📊 System:     ${TOTAL_RAM_MB}MB RAM, $(nproc) CPU cores${NC}"
  echo ""
  echo -e "  ${BOLD}Next steps:${NC}"
  echo -e "  1. Open the dashboard and follow the setup wizard."
  echo ""
  echo -e "  ${DIM}Run 'codedockd status' to check daemon health.${NC}"
  echo -e "  ${DIM}Run 'codedockd logs -f' to follow daemon logs.${NC}"
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

main() {
  setup_colors
  echo -e "${BOLD}🛰️  Codedock — Installing v${RELEASE}${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  check_root
  check_system
  check_ports
  detect_server_ip
  install_docker
  setup_directories
  fetch_config_files
  generate_env_file
  start_codedock
  check_health
  setup_systemd
  print_summary
}

main
