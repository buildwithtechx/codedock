package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"testing"
	"time"
)

func TestBackupProtectionAndRetentionClaimsAreExclusive(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	repo := NewBackupRepo(db, nil)
	ctx := context.Background()
	cfg := &models.BackupConfig{ID: "policy", Name: "Policy", BackupEnabled: true, Schedule: "manual", PreDeployment: true}
	if err := repo.CreateConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateRecord(ctx, &models.BackupRecord{ID: "record", BackupConfigID: cfg.ID, Status: models.BackupRecordStatusCompleted}); err != nil {
		t.Fatal(err)
	}
	if err := repo.ProtectRecord(ctx, "record", time.Now().Add(time.Hour).Unix()); err != nil {
		t.Fatal(err)
	}
	if claimed, err := repo.ClaimRecordExpiry(ctx, "record"); err != nil || claimed {
		t.Fatal("protected record claimed for retention", err)
	}
	if err := repo.DeleteRecord(ctx, "record"); err == nil {
		t.Fatal("direct deletion bypassed protection")
	}
	if claimed, err := repo.ClaimRecordDeletion(ctx, "record"); err != nil || claimed {
		t.Fatal("manual deletion claimed protected record", err)
	}
	if err := repo.DeleteConfig(ctx, cfg.ID); err == nil {
		t.Fatal("deleting a policy bypassed record protection")
	}
	loaded, err := repo.GetConfigByID(ctx, cfg.ID)
	if err != nil || !loaded.PreDeployment {
		t.Fatal("pre-deployment policy flag lost", err)
	}
	if err := repo.ProtectRecord(ctx, "record", 0); err != nil {
		t.Fatal(err)
	}
	if claimed, err := repo.ClaimRecordExpiry(ctx, "record"); err != nil || !claimed {
		t.Fatal("unprotected record could not be claimed", err)
	}
	if err := repo.UpdateRecord(ctx, &models.BackupRecord{ID: "record", Status: models.BackupRecordStatusCompleted}); err == nil {
		t.Fatal("archive verification resurrected an expiry claim")
	}
	if err := repo.ProtectRecord(ctx, "record", time.Now().Add(time.Hour).Unix()); err == nil {
		t.Fatal("protection was added after deletion began")
	}
}
