package build

import (
	"testing"

	"codedock/internal/models"
)

func TestVersionBuildEnv(t *testing.T) {
	env := versionBuildEnv("PHP (composer)", "8.4")
	if env["NIXPACKS_PHP_VERSION"] != "8.4" || env["RAILPACK_PHP_VERSION"] != "8.4" {
		t.Fatalf("unexpected env: %#v", env)
	}
	if env := versionBuildEnv("Node.js (npm/yarn/pnpm)", "22.11.0"); env["NIXPACKS_NODE_VERSION"] != "22.11.0" {
		t.Fatalf("unexpected node env: %#v", env)
	}
	if env := versionBuildEnv("PHP (composer)", ""); env != nil {
		t.Fatalf("empty version should produce no env, got %#v", env)
	}
	if env := versionBuildEnv("Static / Universal HTML", "8.4"); env != nil {
		t.Fatalf("unknown stack should produce no env, got %#v", env)
	}
}

func TestAppRuntimeVersion(t *testing.T) {
	if got := appRuntimeVersion(BuildOptions{}); got != "" {
		t.Fatalf("nil AppConfig should be empty, got %q", got)
	}
	got := appRuntimeVersion(BuildOptions{AppConfig: &models.AppService{RuntimeVersion: "  8.4 "}})
	if got != "8.4" {
		t.Fatalf("expected trimmed version, got %q", got)
	}
}
