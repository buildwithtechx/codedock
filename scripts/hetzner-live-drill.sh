#!/usr/bin/env bash
set -euo pipefail

: "${CODEDOCK_URL:?set CODEDOCK_URL to the daemon api base, e.g. http://localhost:8080/api}"
: "${CODEDOCK_TOKEN:?set CODEDOCK_TOKEN to a login JWT with server scopes}"
: "${CODEDOCK_ORG_ID:?set CODEDOCK_ORG_ID to the drill organization id}"
: "${HETZNER_TOKEN:?set HETZNER_TOKEN to a disposable Hetzner Cloud project token}"

REGION="${DRILL_REGION:-fsn1}"
TYPE="${DRILL_TYPE:-cx23}"
RESIZE_TYPE="${DRILL_RESIZE_TYPE:-cx33}"
IMAGE="${DRILL_IMAGE:-ubuntu-24.04}"
STAMP="$(date +%s)"
LABEL="live-drill-${STAMP}"
NAME="drill-${STAMP}"

api() {
  curl -fsSL -H "Authorization: Bearer ${CODEDOCK_TOKEN}" -H "Content-Type: application/json" "$@"
}

hcloud() {
  curl -fsSL -H "Authorization: Bearer ${HETZNER_TOKEN}" "$@"
}

wait_operation() {
  local id="$1" timeout="$2" waited=0
  while [ "${waited}" -lt "${timeout}" ]; do
    local body status
    body="$(api "${CODEDOCK_URL}/managed-operations/${id}")"
    status="$(echo "${body}" | jq -r .data.status)"
    case "${status}" in
      COMPLETED) echo "${body}"; return 0 ;;
      FAILED|INTERRUPTED)
        echo "operation ${id} ended as ${status}" >&2
        echo "${body}" | jq -r .data.logs >&2
        return 1 ;;
    esac
    sleep 5
    waited=$((waited + 5))
  done
  echo "operation ${id} timed out after ${timeout}s" >&2
  return 1
}

echo "creating drill credential ${LABEL}"
CRED_ID="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/credentials" \
  -d "{\"provider\":\"hetzner\",\"label\":\"${LABEL}\",\"token\":\"${HETZNER_TOKEN}\"}" | jq -r .data.id)"
echo "credential ${CRED_ID}"

echo "checking live offers for ${TYPE} in ${REGION}"
OFFERS="$(api "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/credentials/${CRED_ID}/offers")"
echo "${OFFERS}" | jq -e --arg region "${REGION}" --arg type "${TYPE}" \
  '([.data.locations[].name] | index($region)) != null and (([.data.offers[] | select(.tier.name == $type and (.prices[$region] != null))] | length) == 1)' >/dev/null

echo "reviewing provision ${NAME}"
REVIEW="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/servers/review" \
  -d "{\"credentialId\":\"${CRED_ID}\",\"name\":\"${NAME}\",\"region\":\"${REGION}\",\"serverType\":\"${TYPE}\",\"image\":\"${IMAGE}\"}")"
OP_ID="$(echo "${REVIEW}" | jq -r .data.review.operation.id)"
CONFIRM="$(echo "${REVIEW}" | jq -r .data.review.confirmation)"
SERVER_ID="$(echo "${REVIEW}" | jq -r .data.plan.serverId)"

echo "applying provision ${OP_ID}"
START="$(date +%s)"
api -X POST "${CODEDOCK_URL}/managed-operations/${OP_ID}/apply" \
  -d "{\"confirmation\":\"${CONFIRM}\"}" >/dev/null
wait_operation "${OP_ID}" 600 >/dev/null
END="$(date +%s)"
echo "provision completed in $((END - START))s"

SERVER="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/servers/${SERVER_ID}/refresh")"
echo "${SERVER}" | jq -e '.data.status == "online" and (.data.ipAddress != "")' >/dev/null
EXTERNAL_ID="$(echo "${SERVER}" | jq -r .data.externalId)"
echo "server online at $(echo "${SERVER}" | jq -r .data.ipAddress) (hetzner ${EXTERNAL_ID})"

echo "reviewing resize to ${RESIZE_TYPE}"
RESIZE="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/servers/${SERVER_ID}/resize/review" \
  -d "{\"serverType\":\"${RESIZE_TYPE}\"}")"
RESIZE_OP="$(echo "${RESIZE}" | jq -r .data.review.operation.id)"
RESIZE_CONFIRM="$(echo "${RESIZE}" | jq -r .data.review.confirmation)"
api -X POST "${CODEDOCK_URL}/managed-operations/${RESIZE_OP}/apply" \
  -d "{\"confirmation\":\"${RESIZE_CONFIRM}\"}" >/dev/null
wait_operation "${RESIZE_OP}" 300 >/dev/null
AFTER_RESIZE="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/servers/${SERVER_ID}/refresh")"
echo "${AFTER_RESIZE}" | jq -e --arg id "${EXTERNAL_ID}" --arg type "${RESIZE_TYPE}" \
  '.data.externalId == $id and .data.serverType == $type and .data.status == "online"' >/dev/null
echo "resize kept hetzner server ${EXTERNAL_ID} on type ${RESIZE_TYPE}"

echo "reviewing delete"
DELETE_REVIEW="$(api -X POST "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/servers/${SERVER_ID}/delete/review")"
DELETE_OP="$(echo "${DELETE_REVIEW}" | jq -r .data.review.operation.id)"
DELETE_CONFIRM="$(echo "${DELETE_REVIEW}" | jq -r .data.review.confirmation)"
api -X POST "${CODEDOCK_URL}/managed-operations/${DELETE_OP}/apply" \
  -d "{\"confirmation\":\"${DELETE_CONFIRM}\"}" >/dev/null
wait_operation "${DELETE_OP}" 300 >/dev/null

LEFTOVER_SERVERS="$(hcloud "https://api.hetzner.cloud/v1/servers?label_selector=codedock-server%3D%3D${SERVER_ID}" | jq -r '.servers | length')"
LEFTOVER_KEYS="$(hcloud "https://api.hetzner.cloud/v1/ssh_keys" | jq --arg id "${SERVER_ID}" -r '[.ssh_keys[] | select(.labels["codedock-server"] == $id or (.name | startswith("codedock-m-")))] | length')"
if [ "${LEFTOVER_SERVERS}" != "0" ]; then
  echo "dangling hetzner servers remain for ${SERVER_ID}" >&2
  exit 1
fi
echo "hetzner project clean (servers: 0, codedock keys present: ${LEFTOVER_KEYS})"

api -X DELETE "${CODEDOCK_URL}/organizations/${CODEDOCK_ORG_ID}/managed/credentials/${CRED_ID}" >/dev/null
echo "drill credential removed; live drill passed"
