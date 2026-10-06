package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"testing"
)

func TestBackupPolicyPersistsScheduledOwnerAndProject(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	repo := NewBackupRepo(db, nil)
	ctx := context.Background()
	cfg := &models.BackupConfig{ID: "policy", Name: "Policy", OwnerID: "owner", ProjectID: "project", Schedule: "manual", BackupEnabled: true}
	if err := repo.CreateConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	saved, err := repo.GetConfigByID(ctx, cfg.ID)
	if err != nil {
		t.Fatal(err)
	}
	if saved.OwnerID != "owner" || saved.ProjectID != "project" {
		t.Fatal("scheduled identity lost")
	}
	saved.OwnerID = "replacement"
	if err := repo.UpdateConfig(ctx, saved); err != nil {
		t.Fatal(err)
	}
	all, err := repo.ListConfigs(ctx)
	if err != nil || len(all) != 1 || all[0].OwnerID != "replacement" || all[0].ProjectID != "project" {
		t.Fatal("updated scheduled identity lost", err)
	}
}
