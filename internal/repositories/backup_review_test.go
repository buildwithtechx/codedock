package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"testing"
)

func TestBackupScheduleStatusPersistsOnUpdate(t *testing.T) {
	for _, password := range []string{"", "secret"} {
		t.Run("password="+password, func(t *testing.T) {
			db := openPGTestDB(t)
			repo := NewBackupRepo(db, nil)
			ctx := context.Background()
			cfg := &models.BackupConfig{Name: "manual", Schedule: "manual", BackupEnabled: true, Status: models.BackupConfigStatusInactive}
			if err := repo.CreateConfig(ctx, cfg); err != nil {
				t.Fatal(err)
			}
			cfg.Schedule = "0 2 * * *"
			cfg.Status = models.BackupConfigStatusActive
			cfg.DbPassword = password
			if err := repo.UpdateConfig(ctx, cfg); err != nil {
				t.Fatal(err)
			}
			restarted := NewBackupRepo(db, nil)
			configs, err := restarted.ListAllActiveConfigs(ctx)
			if err != nil || len(configs) != 1 || configs[0].ID != cfg.ID || configs[0].Schedule != cfg.Schedule {
				t.Fatalf("schedule lost after reload: %+v %v", configs, err)
			}
		})
	}
}

func TestBackupAccessScopePrecedesLimit(t *testing.T) {
	db := openPGTestDB(t)
	repo := NewBackupRepo(db, nil)
	ctx := context.Background()
	for _, id := range []string{"allowed", "other"} {
		if err := repo.CreateConfig(ctx, &models.BackupConfig{ID: id, Name: id, Schedule: "manual"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.CreateRecord(ctx, &models.BackupRecord{ID: "older", BackupConfigID: "allowed", StartedAt: "2025-01-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	for i := range 60 {
		if err := repo.CreateRecord(ctx, &models.BackupRecord{ID: fmt.Sprint(i), BackupConfigID: "other", StartedAt: "2026-01-01T00:00:00Z"}); err != nil {
			t.Fatal(err)
		}
	}
	records, err := repo.ListRecordsByConfigs(ctx, []string{"allowed"}, 50)
	if err != nil || len(records) != 1 || records[0].ID != "older" {
		t.Fatalf("accessible backup hidden by global window: %+v %v", records, err)
	}
	records, err = repo.ListRecordsByConfigs(ctx, nil, 50)
	if err != nil || len(records) != 0 {
		t.Fatalf("empty access scope leaked records: %+v %v", records, err)
	}
}

func TestBackupListingExcludesDeletedConfigsBeforeLimit(t *testing.T) {
	db := openPGTestDB(t)
	repo := NewBackupRepo(db, nil)
	ctx := context.Background()
	for _, id := range []string{"retained", "deleted"} {
		if err := repo.CreateConfig(ctx, &models.BackupConfig{ID: id, Name: id, Schedule: "manual"}); err != nil {
			t.Fatal(err)
		}
		startedAt := "2025-01-01T00:00:00Z"
		if id == "deleted" {
			startedAt = "2026-01-01T00:00:00Z"
		}
		if err := repo.CreateRecord(ctx, &models.BackupRecord{ID: id, BackupConfigID: id, StartedAt: startedAt, Status: models.BackupRecordStatusCompleted}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec("ALTER TABLE backup_records DISABLE TRIGGER ALL"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec("ALTER TABLE backup_records ENABLE TRIGGER ALL")
	})
	if err := repo.DeleteConfig(ctx, "deleted"); err != nil {
		t.Fatal(err)
	}
	records, err := repo.ListAllRecords(ctx, 1)
	if err != nil || len(records) != 1 || records[0].BackupConfigID != "retained" {
		t.Fatalf("deleted configuration leaked: %+v %v", records, err)
	}
}
