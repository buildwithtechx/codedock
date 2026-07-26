#!/usr/bin/env bash
set -eo pipefail

echo "🔍 Running pre-flight system diagnostics..."
echo "OS: $(uname -a)"
echo "Docker: $(docker --version 2>/dev/null || echo 'Not installed')"
echo "Memory: $(free -m | awk '/^Mem:/{print $2 " MB total, " $4 " MB free"}')"
echo "CPU Cores: $(nproc)"
echo "✅ Bootstrap verification completed."
