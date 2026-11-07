package security

import (
	"strings"
	"testing"
)

func TestGenerateRuntimeProbeScript(t *testing.T) {
	script := GenerateRuntimeProbeScript([]string{"node", "bun", "php"})
	if !strings.Contains(script, "unset NODE_OPTIONS") {
		t.Fatal("expected unset dangerous options in script")
	}
	if !strings.Contains(script, "command -v node") || !strings.Contains(script, "CODEDOCK_RUNTIME_node") {
		t.Fatal("expected node check in script")
	}
	if !strings.Contains(script, "command -v bun") || !strings.Contains(script, "CODEDOCK_RUNTIME_bun") {
		t.Fatal("expected bun check in script")
	}
	if !strings.Contains(script, "command -v php") || !strings.Contains(script, "CODEDOCK_RUNTIME_php") {
		t.Fatal("expected php check in script")
	}
}

func TestParseRuntimeProbeOutput(t *testing.T) {
	raw := `
Some banner output
CODEDOCK_RUNTIME_node
v22.23.3
CODEDOCK_RUNTIME_php
PHP 8.2.30 (cli) (built: Dec 12 2024 10:00:00) (NTS)
Copyright (c) The PHP Group
Zend Engine v4.2.30, Copyright (c) Zend Technologies
`
	parsed := ParseRuntimeProbeOutput(raw)
	if parsed["node"] != "22.23.3" {
		t.Fatalf("expected node 22.23.3, got %q", parsed["node"])
	}
	if parsed["php"] != "8.2.30" {
		t.Fatalf("expected php 8.2.30, got %q", parsed["php"])
	}

	openshipRaw := `
OPENSHIP_RUNTIME_bun
1.4.3
`
	parsedOpenship := ParseRuntimeProbeOutput(openshipRaw)
	if parsedOpenship["bun"] != "1.4.3" {
		t.Fatalf("expected bun 1.4.3, got %q", parsedOpenship["bun"])
	}
}
