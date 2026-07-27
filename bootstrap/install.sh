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

RELEASE=${CODEDOCK_VERSION:-1.0.0}
CODEDOCK_DIR=/codedock
REPO_URL="https://raw.githubusercontent.com/buildwithtechx/codedock/main"
CTL_URL="$REPO_URL/bootstrap/codedockd"
CTL_SHA256="${CODEDOCK_CTL_SHA256:-}"

echo -e "${BOLD}🛰️  Codedock — Installing v${RELEASE}${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

ensure_root

echo -e "${BOLD}🔍 System check${NC}"
TOTAL_RAM=$(free -m | awk '/^Mem:/{print $2}')
if [ "$TOTAL_RAM" -lt 1024 ]; then
  echo -e "  ${YELLOW}⚠️  ${TOTAL_RAM}MB RAM detected. 1GB+ recommended for Docker builds.${NC}"
  SWAP=$(free -m | awk '/^Swap:/{print $2}')
  if [ "$SWAP" -eq 0 ]; then
    echo -e "  ${YELLOW}   No swap configured. Docker builds may fail on low memory.${NC}"
  else
    echo -e "  ${GREEN}✅ ${SWAP}MB swap detected.${NC}"
  fi
else
  echo -e "  ${GREEN}✅ ${TOTAL_RAM}MB RAM${NC}"
fi

echo -e "${BOLD}🔌 Port check${NC}"
for PORT in 80 443 8080; do
  if ss -tlnp "sport = :$PORT" 2>/dev/null | grep -q .; then
    echo -e "  ${YELLOW}⚠️  Port $PORT is already in use. Codedock needs it.${NC}"
  else
    echo -e "  ${GREEN}✅ Port $PORT available${NC}"
  fi
done

SERVER_IP=$(curl -4fsS ifconfig.me 2>/dev/null || echo "your-server-ip")
echo -e "  ${GREEN}✅ Server IP: ${SERVER_IP}${NC}"
echo ""

ensure_docker

mkdir -p "$CODEDOCK_DIR"/data/{backups,traefik,builds,storage}

echo -e "${BOLD}⬇️  Fetching configuration files...${NC}"
CTL_TMP=$(mktemp)
curl -fsSL "$CTL_URL" -o "$CTL_TMP"
if [ -n "$CTL_SHA256" ]; then
  ACTUAL_SHA=$(sha256sum "$CTL_TMP" | awk '{print $1}')
  if [ "$ACTUAL_SHA" != "$CTL_SHA256" ]; then
    echo -e "${RED}❌ Checksum mismatch for codedockd! Expected: ${CTL_SHA256}, Got: ${ACTUAL_SHA}${NC}"
    rm -f "$CTL_TMP"
    exit 1
  fi
fi
mv "$CTL_TMP" "$CODEDOCK_DIR/codedockd"
chmod +x "$CODEDOCK_DIR/codedockd"
ln -sf "$CODEDOCK_DIR/codedockd" /usr/local/bin/codedockd
ln -sf "$CODEDOCK_DIR/codedockd" /usr/local/bin/codedockctl

if [ ! -f "$CODEDOCK_DIR/.env" ]; then
  echo -e "${BOLD}🔑 Generating .env file...${NC}"
  JWT_SECRET=$(head -c 32 /dev/urandom | base64 | tr -dc 'a-zA-Z0-9' | head -c 32)
  REFRESH_SECRET=$(head -c 32 /dev/urandom | base64 | tr -dc 'a-zA-Z0-9' | head -c 32)

  install -m 0600 /dev/null "$CODEDOCK_DIR/.env"
  cat > "$CODEDOCK_DIR/.env" <<ENV
PORT=8080
HOST=0.0.0.0
CODEDOCK_DATA_DIR=/codedock/data
CODEDOCK_HOST_IP="${SERVER_IP}"
CODEDOCK_JWT_SECRET="${JWT_SECRET}"
CODEDOCK_REFRESH_SECRET="${REFRESH_SECRET}"
CODEDOCK_TLS_EMAIL=""
CODEDOCK_WILDCARD_DOMAIN=""
CODEDOCK_MAGIC_DOMAIN="sslip.io"
DOCKER_SOCKET_PATH=/var/run/docker.sock
CODEDOCK_DASHBOARD_URL="http://${SERVER_IP}:8080"
CODEDOCK_RUNTIME_NETWORK=codedock-network
DEPLOY_HOST_PORT_START=4100
DEPLOY_HOST_PORT_END=4999
DEPLOY_DRY_RUN=false
CODEDOCK_UPDATE_URL=https://api.github.com/repos/buildwithtechx/codedock/releases/latest
CODEDOCK_DOWNLOAD_URL=https://github.com/buildwithtechx/codedock/releases
ENV
  chmod 0600 "$CODEDOCK_DIR/.env"
  echo -e "  ${GREEN}✅ .env written with 0600 permissions.${NC}"
fi

echo -e "${BOLD}🐳 Pulling codedock:${RELEASE}...${NC}"
docker pull "ghcr.io/buildwithtechx/codedock:${RELEASE}"

echo -e "${BOLD}🔧 Creating Docker network...${NC}"
docker network create codedock-network 2>/dev/null || true

echo -e "${BOLD}🚀 Starting Codedock...${NC}"
docker stop codedock-control-plane 2>/dev/null || true
docker rm codedock-control-plane 2>/dev/null || true
docker run -d \
  --name codedock-control-plane \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 80:80 \
  -p 443:443 \
  --env-file "$CODEDOCK_DIR/.env" \
  -e CODEDOCK_DATA_DIR=/codedock/data \
  -v codedock_data:/codedock/data \
  -v "${DOCKER_SOCKET_PATH:-/var/run/docker.sock}":/var/run/docker.sock:ro \
  --network codedock-network \
  --label "traefik.enable=true" \
  --label "traefik.http.routers.dashboard.rule=Host(\`${CODEDOCK_DOMAIN:-codedock.local}\`)" \
  --label "traefik.http.services.dashboard.loadbalancer.server.port=8080" \
  ghcr.io/buildwithtechx/codedock:"${RELEASE}"

echo -e "${BOLD}⏳ Waiting for Codedock health check...${NC}"
for i in $(seq 1 30); do
  if curl -sfS "http://localhost:8080/healthz" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Codedock is healthy.${NC}"
    break
  fi
  sleep 2
done

docker stop codedock-control-plane 2>/dev/null || true
setup_systemd_service "codedock" "[Unit]
Description=Codedock – Self-hosted PaaS
After=docker.service
Requires=docker.service

[Service]
Restart=always
RestartSec=10
WorkingDirectory=/codedock
ExecStart=/usr/bin/docker start -a codedock-control-plane
ExecStop=/usr/bin/docker stop codedock-control-plane

[Install]
WantedBy=multi-user.target"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Codedock v${RELEASE} installed successfully!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "  ${BOLD}Dashboard:${NC}  http://${SERVER_IP}:8080"
echo -e "  ${BOLD}Server CLI:${NC} codedockd status"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
