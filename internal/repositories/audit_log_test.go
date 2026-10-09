package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestAuditLogRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewAuditLogRepo(db)
	first := &models.AuditLog{ID: "log-1", UserID: "user-1", Action: "create", Resource: "project", Details: "created", IPAddress: "10.0.0.1"}
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("create first log: %v", err)
	}
	second := &models.AuditLog{ID: "log-2", UserID: "user-1", Action: "delete", Resource: "service"}
	if err := repo.Create(ctx, second); err != nil {
		t.Fatalf("create second log: %v", err)
	}
	logs, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("list logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
	for _, entry := range logs {
		if entry.CreatedAt == "" {
			t.Fatalf("expected created_at to roundtrip: %+v", entry)
		}
	}
	if logs[0].ID != "log-2" || logs[1].ID != "log-1" {
		t.Fatalf("expected newest-first order, got %q then %q", logs[0].ID, logs[1].ID)
	}
	if logs[1].Details != "created" || logs[1].IPAddress != "10.0.0.1" {
		t.Fatalf("optional columns did not roundtrip: %+v", logs[1])
	}
	page, err := repo.List(ctx, 1, 1)
	if err != nil {
		t.Fatalf("list page: %v", err)
	}
	if len(page) != 1 || page[0].ID != "log-1" {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestAuditLogFacetsAndFilter(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewAuditLogRepo(db)
	seed := []*models.AuditLog{
		{ID: "facet-1", UserID: "user-1", Action: "deployment.trigger", Resource: "svc-1"},
		{ID: "facet-2", UserID: "user-1", Action: "deployment.rollback", Resource: "svc-1"},
		{ID: "facet-3", UserID: "user-1", Action: "server.create", Resource: "srv-1"},
		{ID: "facet-4", UserID: "user-1", Action: "backup.trigger", Resource: "cfg-1"},
	}
	for _, entry := range seed {
		if err := repo.Create(ctx, entry); err != nil {
			t.Fatalf("seed %s: %v", entry.ID, err)
		}
	}
	counts, total, err := repo.Facets(ctx)
	if err != nil {
		t.Fatalf("facets: %v", err)
	}
	if total != len(seed) {
		t.Fatalf("expected total %d, got %d", len(seed), total)
	}
	if counts["deployments"] != 2 || counts["servers"] != 1 || counts["system"] != 1 {
		t.Fatalf("unexpected facet counts: %+v", counts)
	}
	filtered, err := repo.ListFiltered(ctx, models.CategoryPrefixes("deployments"), 10, 0)
	if err != nil {
		t.Fatalf("filtered list: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 deployment logs, got %d", len(filtered))
	}
}
