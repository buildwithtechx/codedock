package commands

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	authservices "codedock.run/codedock/internal/services/auth"
)

func setupTestDatabaseURL(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("CODEDOCK_TEST_PG_URL")
	if dsn == "" {
		t.Skip("CODEDOCK_TEST_PG_URL is not set")
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse test database url: %v", err)
	}
	parsed.Path = "/codedock_setup_test"
	target := parsed.String()
	parsed.Path = "/postgres"
	maintenance, err := sql.Open("pgx", parsed.String())
	if err != nil {
		t.Fatalf("open maintenance database: %v", err)
	}
	defer maintenance.Close()
	ctx := context.Background()
	if _, err := maintenance.ExecContext(ctx, `DROP DATABASE IF EXISTS "codedock_setup_test" WITH (FORCE)`); err != nil {
		t.Fatalf("drop setup test database: %v", err)
	}
	if _, err := maintenance.ExecContext(ctx, `CREATE DATABASE "codedock_setup_test"`); err != nil {
		t.Fatalf("create setup test database: %v", err)
	}
	t.Cleanup(func() {
		cleanup, err := sql.Open("pgx", parsed.String())
		if err != nil {
			return
		}
		defer cleanup.Close()
		_, _ = cleanup.ExecContext(context.Background(), `DROP DATABASE IF EXISTS "codedock_setup_test" WITH (FORCE)`)
	})
	return target
}

func TestFirstRunWithoutEnvironmentAndSavedSetup(t *testing.T) {
	cfg := config.Get()
	previous := *cfg
	t.Cleanup(func() { *cfg = previous })
	cfg.Cloud.Enabled = false
	cfg.Server.DataDir = t.TempDir()
	cfg.Database.URL = setupTestDatabaseURL(t)
	cfg.Security.JWTSecret = ""
	cfg.Security.RefreshSecret = ""
	cfg.Security.TLSEmail = ""
	cfg.Telemetry.Salt = ""
	cfg.Domains.WildcardDomain = ""
	_, db, vault := InitDataDir()
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	tokens, err := authservices.NewTokenService()
	if err != nil {
		t.Fatal(err)
	}
	token, err := tokens.GenerateToken(&models.User{ID: "owner", Email: "owner@example.com", Role: models.UserRoleOwner})
	if err != nil {
		t.Fatal(err)
	}
	if err := saveSetupOptions(db, vault, "apps.example.com", "owner@example.com"); err != nil {
		t.Fatal(err)
	}
	settings, err := repositories.NewSettingsRepo(db, vault).GetServerSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.DefaultWildcardDomain != "apps.example.com" {
		t.Fatal("wizard did not persist the app domain")
	}
	cfg.Security.JWTSecret = ""
	cfg.Security.RefreshSecret = ""
	if err := config.PrepareSelfHosted(cfg); err != nil {
		t.Fatal(err)
	}
	restarted, err := authservices.NewTokenService()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.ValidateToken(token); err != nil {
		t.Fatalf("token stopped working after restart: %v", err)
	}
	if cfg.Security.TLSEmail != "owner@example.com" {
		t.Fatal("certificate email was not restored")
	}
}

func TestSetupRejectsMalformedDomainNames(t *testing.T) {
	for _, domain := range []string{"a..b.com", "-apps.example.com", "apps-.example.com", "apps_example.com", "apps.example.com?", "apps.example.com#", "owner@apps.example.com", "localhost"} {
		if err := validateSetupDomain(domain); err == nil {
			t.Errorf("accepted malformed domain: %s", domain)
		}
	}
	for _, domain := range []string{"apps.example.com", "my-apps.example.com", "EXAMPLE.COM"} {
		if err := validateSetupDomain(domain); err != nil {
			t.Errorf("rejected %s: %v", domain, err)
		}
	}
}
