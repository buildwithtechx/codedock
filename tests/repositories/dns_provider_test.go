package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestDNSProviderSettingsRoundTrip(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewSettingsRepo(db)
	settings, err := repo.GetServerSettings(context.Background())
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}
	settings.CloudflareAPIToken = "cloudflare-token"
	settings.NamecheapAPIUser = "namecheap-user"
	settings.NamecheapAPIKey = "namecheap-key"
	settings.SpaceshipAPIKey = "spaceship-key"
	settings.SpaceshipAPISecret = "spaceship-secret"
	if err := repo.UpdateServerSettings(context.Background(), settings); err != nil {
		t.Fatalf("update settings: %v", err)
	}

	updated, err := repo.GetServerSettings(context.Background())
	if err != nil {
		t.Fatalf("reload settings: %v", err)
	}
	if updated.CloudflareAPIToken != "cloudflare-token" || updated.NamecheapAPIKey != "namecheap-key" || updated.SpaceshipAPISecret != "spaceship-secret" {
		t.Fatalf("provider credentials did not round-trip: %+v", updated)
	}
}

func TestDNSRecordRepositoryLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewDNSRepo(db)
	record := &models.DNSRecord{
		DomainName:  "app.example.com",
		RecordType:  "A",
		RecordName:  "app",
		RecordValue: "192.0.2.10",
	}
	if err := repo.Create(context.Background(), record); err != nil {
		t.Fatalf("create dns record: %v", err)
	}
	list, err := repo.ListByDomain(context.Background(), record.DomainName)
	if err != nil || len(list) != 1 {
		t.Fatalf("list dns records: count=%d err=%v", len(list), err)
	}
	if err := repo.Delete(context.Background(), record.ID); err != nil {
		t.Fatalf("delete dns record: %v", err)
	}
	list, err = repo.ListByDomain(context.Background(), record.DomainName)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected dns record deletion, count=%d err=%v", len(list), err)
	}
}
