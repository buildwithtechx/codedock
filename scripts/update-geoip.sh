#!/usr/bin/env bash
set -euo pipefail

: "${MAXMIND_ACCOUNT_ID:?set MAXMIND_ACCOUNT_ID from https://www.maxmind.com/en/accounts/current/license-key}"
: "${MAXMIND_LICENSE_KEY:?set MAXMIND_LICENSE_KEY from https://www.maxmind.com/en/accounts/current/license-key}"

EDITION="${GEOIP_EDITION:-GeoLite2-Country}"
DOWNLOAD_URL="https://download.maxmind.com/geoip/databases/${EDITION}/download?suffix=tar.gz"
TARGET_DIR="internal/engine/traffic/assets"
TARGET_FILE="${TARGET_DIR}/GeoLite2-Country.mmdb"
WORK_DIR="$(mktemp -d)"
trap 'rm -rf "${WORK_DIR}"' EXIT

curl -fsSL -u "${MAXMIND_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY}" -o "${WORK_DIR}/geoip.tar.gz" "${DOWNLOAD_URL}"
curl -fsSL -u "${MAXMIND_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY}" -o "${WORK_DIR}/geoip.tar.gz.sha256" "${DOWNLOAD_URL}.sha256"

EXPECTED="$(cut -d' ' -f1 "${WORK_DIR}/geoip.tar.gz.sha256")"
ACTUAL="$(sha256sum "${WORK_DIR}/geoip.tar.gz" | cut -d' ' -f1)"
if [ "${EXPECTED}" != "${ACTUAL}" ]; then
  echo "checksum mismatch for ${EDITION}" >&2
  exit 1
fi

tar -xzf "${WORK_DIR}/geoip.tar.gz" -C "${WORK_DIR}"
mkdir -p "${TARGET_DIR}"
find "${WORK_DIR}" -name "${EDITION}.mmdb" -exec cp {} "${TARGET_FILE}" \;
echo "Updated ${EDITION} at ${TARGET_FILE}"
