#!/usr/bin/env bash
set -eo pipefail

BOLD="\033[1m"
DIM="\033[2m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
NC="\033[0m"

ensure_root() {
  if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Please run as root (or with sudo).${NC}"
    exit 1
  fi
}

detect_platform() {
  local os arch
  os="$(uname -s)"
  arch="$(uname -m)"
  case "$os" in
    Linux)  os="linux" ;;
    Darwin) os="darwin" ;;
    *) echo -e "${RED}❌ Unsupported OS: $os${NC}"; exit 1 ;;
  esac
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    armv7l) arch="arm" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $arch${NC}"; exit 1 ;;
  esac
  echo "${os}_${arch}"
}

ensure_docker() {
  if ! command -v docker &> /dev/null; then
    echo -e "${BOLD}📦 Installing Docker...${NC}"
    if command -v apt-get &> /dev/null; then
      apt-get update -qq && apt-get install -y -qq docker.io 2>/dev/null || true
    elif command -v dnf &> /dev/null; then
      dnf install -y -q docker 2>/dev/null || true
    elif command -v yum &> /dev/null; then
      yum install -y -q docker 2>/dev/null || true
    fi
    if ! command -v docker &> /dev/null; then
      local docker_tmp
      docker_tmp=$(mktemp)
      curl -fsSL https://get.docker.com -o "$docker_tmp"
      sh "$docker_tmp"
      rm -f "$docker_tmp"
    fi
    systemctl enable --now docker 2>/dev/null || true
  fi

  if ! docker info &> /dev/null; then
    echo -e "${YELLOW}⏳ Waiting for Docker daemon to start...${NC}"
    systemctl start docker 2>/dev/null || true
    for _ in $(seq 1 10); do
      if docker info &> /dev/null; then break; fi
      sleep 1
    done
  fi

  if ! docker info &> /dev/null; then
    echo -e "${RED}❌ Docker daemon is not running or unreachable.${NC}"
    exit 1
  fi

  echo -e "${GREEN}✅ Docker ready${NC}"
}

setup_systemd_service() {
  local service_name="$1"
  local service_content="$2"
  if command -v systemctl &> /dev/null; then
    echo "$service_content" > "/etc/systemd/system/${service_name}.service"
    systemctl daemon-reload
    systemctl enable --now "${service_name}.service"
  fi
}
