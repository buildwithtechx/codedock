#!/usr/bin/env bash
set -eo pipefail

if [ -z "$1" ]; then
  echo "❌ Usage: ./scripts/restore.sh <path-to-backup.tar.gz>"
  exit 1
fi

BACKUP_FILE="$1"
CODEDOCK_DIR=${CODEDOCK_DIR:-/codedock}

if [ ! -f "${BACKUP_FILE}" ]; then
  echo "❌ Error: Backup file ${BACKUP_FILE} not found!"
  exit 1
fi

if [ ! -d "$CODEDOCK_DIR/data" ]; then
  if [ -d "./data" ]; then
    CODEDOCK_DIR="."
  else
    echo "❌ No Codedock data directory found at $CODEDOCK_DIR/data."
    exit 1
  fi
fi

DEST_DIR=$(realpath "$CODEDOCK_DIR")

check_paths() {
  local bad=0
  while IFS= read -r entry; do
    local resolved
    resolved=$(realpath -m "${DEST_DIR}/${entry}")
    if [[ "$entry" == *".."* ]] || [[ "$entry" == /* ]] || [[ "$resolved" != "${DEST_DIR}"/* && "$resolved" != "${DEST_DIR}" ]]; then
      echo "❌ Path traversal detected in archive: ${entry}"
      bad=1
    fi
  done < <(tar -tzf "${BACKUP_FILE}" 2>/dev/null)

  if [ "$bad" -ne 0 ]; then
    echo "❌ Aborting restore: archive contains unsafe paths."
    exit 1
  fi
}

echo "🔍 Validating archive paths..."
check_paths

echo "⚠️  Restoring Codedock state from ${BACKUP_FILE} into ${CODEDOCK_DIR}/data..."
tar --no-same-owner -xzf "${BACKUP_FILE}" -C "${CODEDOCK_DIR}"

echo "🔄 Restarting Codedock container to apply restored data..."
if command -v docker &> /dev/null && [ -f "$CODEDOCK_DIR/docker-compose.yml" ]; then
  docker compose -f "$CODEDOCK_DIR/docker-compose.yml" restart codedock
fi

echo "✅ Restore completed successfully! Codedock is back online."
