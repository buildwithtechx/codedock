package repositories

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/operations"
	"codedock.run/codedock/internal/utils"
	"context"
	"strings"
	"testing"
	"time"
)

func TestOperationsEncryptReviewsAndRejectConflicts(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewOperationRepo(db, vault)
	service := operations.NewService(repo)
	ctx := context.Background()
	review, err := service.Review(ctx, "owner", "project", "restore", "database", "fixture-secret", "version-1", "Replace data")
	if err != nil {
		t.Fatal(err)
	}
	var encrypted string
	if err := db.QueryRow(`SELECT encrypted_payload FROM operations WHERE id=?`, review.Operation.ID).Scan(&encrypted); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(encrypted, "fixture-secret") {
		t.Fatal("operation payload stored in plain text")
	}
	ran := make(chan struct{}, 1)
	runner := func(context.Context, *models.Operation, func(string, string) error) error {
		ran <- struct{}{}
		return nil
	}
	for _, test := range []struct{ user, token, snapshot string }{{"other", review.Confirmation, "version-1"}, {"owner", "wrong", "version-1"}, {"owner", review.Confirmation, "version-2"}} {
		if err := service.Apply(ctx, review.Operation.ID, test.user, test.token, test.snapshot, runner); err == nil {
			t.Fatal("invalid confirmation was accepted")
		}
	}
	if err := repo.Claim(ctx, review.Operation.ID, "version-1"); err != nil {
		t.Fatal(err)
	}
	second, err := service.Review(ctx, "owner", "project", "restore", "database", "payload", "version-1", "Replace data")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Claim(ctx, second.Operation.ID, "version-1"); err == nil {
		t.Fatal("conflicting target operations both started")
	}
	if err := repo.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Get(ctx, review.Operation.ID)
	if err != nil || loaded.Status != "INTERRUPTED" || loaded.Payload != "fixture-secret" {
		t.Fatal("restart recovery or decryption failed", err)
	}
	expired, err := service.Review(ctx, "owner", "project", "restore", "other", "payload", "version-1", "Replace data")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE operations SET expires_at=? WHERE id=?`, time.Now().Add(-time.Minute).Unix(), expired.Operation.ID); err != nil {
		t.Fatal(err)
	}
	if err := service.Apply(ctx, expired.Operation.ID, "owner", expired.Confirmation, "version-1", runner); err == nil {
		t.Fatal("expired review was applied")
	}
	select {
	case <-ran:
		t.Fatal("rejected operation changed target state")
	default:
	}
}
