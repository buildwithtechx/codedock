package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func TestServerRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewServerRepository(db, vault)
	now := time.Now().UTC().Truncate(time.Second)
	server := &models.Server{
		ID:            "srv-1",
		UserID:        "user-1",
		Name:          "web",
		IPAddress:     "10.0.0.1",
		IsLocal:       true,
		SSHHost:       "example.com",
		SSHPort:       22,
		SSHUser:       "root",
		SSHAuthMethod: "key",
		SSHKey:        "key-material",
		Status:        models.ServerStatusOnline,
		Provider:      "hetzner",
		Region:        "fsn1",
		WorkerToken:   "token-1",
		Metrics:       []byte(`{"cpu":1}`),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := repo.Create(ctx, server); err != nil {
		t.Fatalf("create server: %v", err)
	}
	loaded, err := repo.GetByID(ctx, "srv-1")
	if err != nil {
		t.Fatalf("get server: %v", err)
	}
	if loaded == nil || loaded.Name != "web" || !loaded.IsLocal || loaded.SSHKey != "key-material" || string(loaded.Metrics) != `{"cpu":1}` {
		t.Fatalf("unexpected server row: %+v", loaded)
	}
	if !loaded.CreatedAt.Truncate(time.Second).Equal(now) || !loaded.UpdatedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("server timestamps mismatch: %+v", loaded)
	}
	if loaded.LastSeenAt != nil {
		t.Fatalf("expected nil last seen, got %v", loaded.LastSeenAt)
	}
	byToken, err := repo.GetByToken(ctx, "token-1")
	if err != nil {
		t.Fatalf("get by token: %v", err)
	}
	if byToken == nil || byToken.ID != "srv-1" {
		t.Fatalf("token lookup returned wrong server: %+v", byToken)
	}
	listed, err := repo.ListByUser(ctx, "user-1")
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != "srv-1" {
		t.Fatalf("expected 1 server, got %d", len(listed))
	}
	server.Name = "web-renamed"
	server.UpdatedAt = time.Now().UTC()
	if err := repo.Update(ctx, server); err != nil {
		t.Fatalf("update server: %v", err)
	}
	renamed, err := repo.GetByID(ctx, "srv-1")
	if err != nil {
		t.Fatalf("get renamed server: %v", err)
	}
	if renamed == nil || renamed.Name != "web-renamed" {
		t.Fatalf("update did not persist: %+v", renamed)
	}
	if err := repo.UpdateStatus(ctx, "srv-1", models.ServerStatusOffline); err != nil {
		t.Fatalf("update status: %v", err)
	}
	offline, err := repo.GetByID(ctx, "srv-1")
	if err != nil {
		t.Fatalf("get offline server: %v", err)
	}
	if offline == nil || offline.Status != models.ServerStatusOffline {
		t.Fatalf("status update did not persist: %+v", offline)
	}
	if err := repo.UpdateMetrics(ctx, "srv-1", []byte(`{"cpu":2}`)); err != nil {
		t.Fatalf("update metrics: %v", err)
	}
	metered, err := repo.GetByID(ctx, "srv-1")
	if err != nil {
		t.Fatalf("get metered server: %v", err)
	}
	if metered == nil || string(metered.Metrics) != `{"cpu":2}` || metered.LastSeenAt == nil {
		t.Fatalf("metrics update did not persist: %+v", metered)
	}
	if err := repo.Delete(ctx, "srv-1"); err != nil {
		t.Fatalf("delete server: %v", err)
	}
	deleted, err := repo.GetByID(ctx, "srv-1")
	if err != nil {
		t.Fatalf("get deleted server: %v", err)
	}
	if deleted != nil {
		t.Fatal("expected deleted server to be gone")
	}
}
