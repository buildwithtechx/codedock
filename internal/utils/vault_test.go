package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codedock/internal/config"
)

func TestVaultKeyOverrideSkipsKeyFile(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Security.VaultKey
	t.Cleanup(func() { cfg.Security.VaultKey = previous })
	cfg.Security.VaultKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	first, err := NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("vault with override: %v", err)
	}
	second, err := NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("second vault with override: %v", err)
	}
	sealed, err := first.Encrypt("shared secret")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	opened, err := second.Decrypt(sealed)
	if err != nil {
		t.Fatalf("cross-instance decrypt: %v", err)
	}
	if opened != "shared secret" {
		t.Fatal("override keys diverged between instances")
	}
}

func TestVaultKeyOverrideRejectsMalformedKeys(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Security.VaultKey
	t.Cleanup(func() { cfg.Security.VaultKey = previous })
	for _, key := range []string{"too-short", strings.Repeat("zz", 32), strings.Repeat("ab", 31)} {
		cfg.Security.VaultKey = key
		if _, err := NewVault(t.TempDir()); err == nil {
			t.Fatalf("accepted malformed vault key %q", key)
		}
	}
}

func TestVaultGeneratesKeyFileWithoutOverride(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Security.VaultKey
	t.Cleanup(func() { cfg.Security.VaultKey = previous })
	cfg.Security.VaultKey = ""
	dir := t.TempDir()
	vault, err := NewVault(dir)
	if err != nil {
		t.Fatalf("vault without override: %v", err)
	}
	key, err := os.ReadFile(filepath.Join(dir, ".vault_key"))
	if err != nil {
		t.Fatalf("key file was not created: %v", err)
	}
	if len(key) != 32 {
		t.Fatal("generated key has the wrong length")
	}
	sealed, err := vault.Encrypt("roundtrip")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	reloaded, err := NewVault(dir)
	if err != nil {
		t.Fatalf("reload vault: %v", err)
	}
	opened, err := reloaded.Decrypt(sealed)
	if err != nil || opened != "roundtrip" {
		t.Fatal("key file did not survive reload")
	}
}
