package config

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"codedock.run/codedock/pkg/types"
)

func TestSelfHostedSecretsSurviveRestart(t *testing.T) {
	dataDir := t.TempDir()
	first := &types.Config{Server: types.ServerConfig{DataDir: dataDir}}
	if err := PrepareSelfHosted(first); err != nil {
		t.Fatal(err)
	}
	secrets := []string{first.Security.JWTSecret, first.Security.RefreshSecret, first.Telemetry.Salt}
	for _, secret := range secrets {
		decoded, err := hex.DecodeString(secret)
		if err != nil || len(decoded) != 32 {
			t.Fatal("expected a 256-bit secret")
		}
	}
	if secrets[0] == secrets[1] || secrets[0] == secrets[2] || secrets[1] == secrets[2] {
		t.Fatal("secrets must be independent")
	}
	if err := SaveSelfHostedOptions(first, "apps.example.com", "owner@example.com"); err != nil {
		t.Fatal(err)
	}
	second := &types.Config{Server: types.ServerConfig{DataDir: dataDir}}
	if err := PrepareSelfHosted(second); err != nil {
		t.Fatal(err)
	}
	if first.Security.JWTSecret != second.Security.JWTSecret || first.Security.RefreshSecret != second.Security.RefreshSecret || first.Telemetry.Salt != second.Telemetry.Salt {
		t.Fatal("restart changed authentication secrets")
	}
	if second.Domains.WildcardDomain != "apps.example.com" || second.Security.TLSEmail != "owner@example.com" {
		t.Fatal("setup options were not restored")
	}
	info, err := os.Stat(filepath.Join(dataDir, "self-hosted.json"))
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatal("configuration must be private")
	}
}

func TestSelfHostedExplicitSecretsTakePrecedence(t *testing.T) {
	cfg := &types.Config{Server: types.ServerConfig{DataDir: t.TempDir()}, Security: types.SecurityConfig{JWTSecret: "existing-jwt", RefreshSecret: "existing-refresh"}}
	if err := PrepareSelfHosted(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "existing-jwt" || cfg.Security.RefreshSecret != "existing-refresh" {
		t.Fatal("explicit secrets were replaced")
	}
	cfg.Security.JWTSecret = "override"
	if err := PrepareSelfHosted(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "override" {
		t.Fatal("explicit override was ignored")
	}
}

func TestSelfHostedCorruptConfigurationIsNotOverwritten(t *testing.T) {
	dataDir := t.TempDir()
	path := filepath.Join(dataDir, "self-hosted.json")
	if err := os.WriteFile(path, []byte("invalid"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := PrepareSelfHosted(&types.Config{Server: types.ServerConfig{DataDir: dataDir}}); err == nil {
		t.Fatal("corrupt configuration must fail")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "invalid" {
		t.Fatal("corrupt configuration was replaced")
	}
}

func TestCloudDoesNotGenerateSelfHostedSecrets(t *testing.T) {
	dataDir := t.TempDir()
	cfg := &types.Config{Server: types.ServerConfig{DataDir: dataDir}, Cloud: types.CloudConfig{Enabled: true}}
	if err := PrepareSelfHosted(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "" {
		t.Fatal("cloud credentials must be configured by the operator")
	}
	if _, err := os.Stat(filepath.Join(dataDir, "self-hosted.json")); !os.IsNotExist(err) {
		t.Fatal("cloud created self-hosted configuration")
	}
}

func TestSelfHostedOptionsCanBeCleared(t *testing.T) {
	cfg := &types.Config{Server: types.ServerConfig{DataDir: t.TempDir()}}
	if err := PrepareSelfHosted(cfg); err != nil {
		t.Fatal(err)
	}
	if err := SaveSelfHostedOptions(cfg, "apps.example.com", "owner@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := SaveSelfHostedOptions(cfg, "", ""); err != nil {
		t.Fatal(err)
	}
	restarted := &types.Config{Server: types.ServerConfig{DataDir: cfg.Server.DataDir}}
	if err := PrepareSelfHosted(restarted); err != nil {
		t.Fatal(err)
	}
	if restarted.Domains.WildcardDomain != "" || restarted.Security.TLSEmail != "" {
		t.Fatal("cleared routing options restored on restart")
	}
	if restarted.Security.JWTSecret != cfg.Security.JWTSecret {
		t.Fatal("clearing options changed the auth secret")
	}
}
