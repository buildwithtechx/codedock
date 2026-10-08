package repositories

import (
	"context"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

func TestTakeoverRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	users := NewUserRepo(db)
	user := &models.User{Email: "takeover@example.com", Name: "Takeover", PasswordHash: "hash", Role: models.UserRoleOwner, IsActive: true}
	if err := users.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	repo := NewTakeoverRepository(db)
	run := &models.TakeoverRun{
		UserID:         user.ID,
		SourceHost:     "10.0.0.2",
		SourcePlatform: models.TakeoverPlatformDocker,
		Status:         models.TakeoverStatusScanning,
		DiscoveredJSON: `{"containers":[]}`,
	}
	if err := repo.Create(ctx, run); err != nil {
		t.Fatalf("create run: %v", err)
	}
	if run.ID == "" {
		t.Fatal("expected generated run id")
	}
	loaded, err := repo.GetByID(ctx, run.ID)
	if err != nil {
		t.Fatalf("get run: %v", err)
	}
	if loaded.UserID != user.ID || loaded.SourceHost != "10.0.0.2" || loaded.DiscoveredJSON != `{"containers":[]}` {
		t.Fatalf("unexpected run row: %+v", loaded)
	}
	if loaded.CreatedAt.IsZero() || loaded.UpdatedAt.IsZero() {
		t.Fatal("expected run timestamps to roundtrip")
	}
	if err := repo.UpdateStatus(ctx, run.ID, models.TakeoverStatusFailed, "boom"); err != nil {
		t.Fatalf("update status: %v", err)
	}
	failed, err := repo.GetByID(ctx, run.ID)
	if err != nil {
		t.Fatalf("get failed run: %v", err)
	}
	if failed.Status != models.TakeoverStatusFailed || failed.Error != "boom" {
		t.Fatalf("status update did not persist: %+v", failed)
	}
	if err := repo.UpdateDiscovered(ctx, run.ID, `{"containers":[{"name":"web"}]}`); err != nil {
		t.Fatalf("update discovered: %v", err)
	}
	scanned, err := repo.GetByID(ctx, run.ID)
	if err != nil {
		t.Fatalf("get scanned run: %v", err)
	}
	if scanned.Status != models.TakeoverStatusScanned || scanned.DiscoveredJSON != `{"containers":[{"name":"web"}]}` {
		t.Fatalf("discovered update did not persist: %+v", scanned)
	}
	if err := repo.UpdateAdopted(ctx, run.ID, []string{"p1", "p2"}); err != nil {
		t.Fatalf("update adopted: %v", err)
	}
	done, err := repo.GetByID(ctx, run.ID)
	if err != nil {
		t.Fatalf("get done run: %v", err)
	}
	if done.Status != models.TakeoverStatusDone || done.AdoptedProjectIDs != "p1,p2" {
		t.Fatalf("adopted update did not persist: %+v", done)
	}
	if done.UpdatedAt.Truncate(time.Second).Before(loaded.CreatedAt.Truncate(time.Second)) {
		t.Fatal("expected updated_at to advance")
	}
	runs, err := repo.ListByUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(runs) != 1 || runs[0].ID != run.ID {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
}
