package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
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
