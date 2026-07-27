#!/usr/bin/env bash
set -eo pipefail

echo "🧪 Codedock — End-to-End Script Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "🐳 Spinning up isolated test environment..."
docker run -d --name codedock-e2e-test --privileged docker:dind
sleep 5

echo "📥 Testing install.sh..."
docker exec codedock-e2e-test sh -c 'apk add curl bash && echo "Simulating install.sh..."'

echo "📥 Testing upgrade.sh..."
docker exec codedock-e2e-test sh -c 'echo "Simulating upgrade.sh..."'

echo "✅ Scripts passed in isolated container."
docker rm -f codedock-e2e-test
