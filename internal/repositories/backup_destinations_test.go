package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestSFTPDestinationRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed project chain: %v", err)
		}
	}
	repo := NewSFTPDestinationRepo(db)
	dest := &models.SFTPDestination{
		OrganizationID: "org-one",
		ProjectID:      "project-one",
		Name:           "sftp-one",
		Host:           "sftp.example.com",
		Username:       "uploader",
		Password:       "secret",
	}
	if err := repo.Create(ctx, dest); err != nil {
		t.Fatalf("create destination: %v", err)
	}
	if dest.Port != 22 {
		t.Fatalf("expected default port 22, got %d", dest.Port)
	}
	stored, err := repo.Get(ctx, dest.ID)
	if err != nil {
		t.Fatalf("get destination: %v", err)
	}
	if stored.Host != "sftp.example.com" || stored.Password != "secret" {
		t.Fatalf("unexpected destination row: %+v", stored)
	}
	list, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list destinations: %v", err)
	}
	if len(list) != 1 || list[0].Password != "" {
		t.Fatalf("expected credentials redacted in listing: %+v", list)
	}
	if err := repo.Delete(ctx, dest.ID); err != nil {
		t.Fatalf("delete destination: %v", err)
	}
	if _, err := repo.Get(ctx, dest.ID); err == nil {
		t.Fatal("expected deleted destination to be gone")
	}
}

func TestPolicyBatchRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed project chain: %v", err)
		}
	}
	repo := NewPolicyBatchRepo(db)
	batch := &models.BackupPolicyBatch{ProjectID: "project-one", Name: "nightly", Schedule: "0 2 * * *"}
	if err := repo.Create(ctx, batch); err != nil {
		t.Fatalf("create batch: %v", err)
	}
	stored, err := repo.Get(ctx, batch.ID)
	if err != nil {
		t.Fatalf("get batch: %v", err)
	}
	if stored.Name != "nightly" || stored.Status != "active" {
		t.Fatalf("unexpected batch row: %+v", stored)
	}
	list, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list batches: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(list))
	}
	backups := NewBackupRepo(db, nil)
	cfg := &models.BackupConfig{ID: "member", Name: "member", Schedule: "manual", ParentBatchID: batch.ID}
	if err := backups.CreateConfig(ctx, cfg); err != nil {
		t.Fatalf("create member config: %v", err)
	}
	members, err := repo.ListMembers(ctx, batch.ID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != 1 || members[0].ID != "member" {
		t.Fatalf("expected 1 member, got %+v", members)
	}
}

func TestRedisSnapshotRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO cluster_databases (id, organization_id, project_id, name, engine, version, port, username, database_name, secret_encrypted) VALUES ('db-one', 'org-one', 'project-one', 'cache', 'redis', '7', 6379, 'default', 'cache', 'enc')`,
		`INSERT INTO clusters (id, project_id, organization_id, name, version, nodes_json, encrypted_token, updated_at) VALUES ('cluster-one', 'project-one', 'org-one', 'cluster', '1', '[]', 'enc', '2026-01-01T00:00:00Z')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed cluster database: %v", err)
		}
	}
	repo := NewRedisSnapshotRepo(db)
	snapshot := &models.RedisSnapshot{DatabaseID: "db-one", ProjectID: "project-one", ClusterID: "cluster-one", S3Key: "snapshots/a.rdb", SizeBytes: 1024, Status: "completed"}
	if err := repo.Create(ctx, snapshot); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}
	stored, err := repo.Get(ctx, snapshot.ID)
	if err != nil {
		t.Fatalf("get snapshot: %v", err)
	}
	if stored.S3Key != "snapshots/a.rdb" || stored.SizeBytes != 1024 {
		t.Fatalf("unexpected snapshot row: %+v", stored)
	}
	list, err := repo.List(ctx, "db-one")
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(list))
	}
}

func TestS3DestinationRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewS3DestinationRepo(db, databaseTestVault{})
	first := &models.S3Destination{Name: "first", Endpoint: "https://s3.example.com", Bucket: "one", AccessKeyID: "ak", SecretAccessKey: "sk", IsDefault: true}
	if err := repo.CreateS3Destination(ctx, first); err != nil {
		t.Fatalf("create first: %v", err)
	}
	second := &models.S3Destination{Name: "second", Endpoint: "https://s3.example.com", Bucket: "two", IsDefault: true}
	if err := repo.CreateS3Destination(ctx, second); err != nil {
		t.Fatalf("create second: %v", err)
	}
	stored, err := repo.GetS3Destination(ctx, first.ID)
	if err != nil {
		t.Fatalf("get first: %v", err)
	}
	if stored.IsDefault {
		t.Fatal("first destination should have lost default status")
	}
	if stored.SecretAccessKey != "sk" {
		t.Fatalf("expected decrypted secret, got %q", stored.SecretAccessKey)
	}
	if err := repo.SetDefaultDestination(ctx, first.ID); err != nil {
		t.Fatalf("set default: %v", err)
	}
	restored, err := repo.GetS3Destination(ctx, first.ID)
	if err != nil {
		t.Fatalf("get restored default: %v", err)
	}
	if !restored.IsDefault {
		t.Fatal("expected default status to move back")
	}
	if err := repo.RecordVerificationResult(ctx, first.ID, true, ""); err != nil {
		t.Fatalf("record verification: %v", err)
	}
	verified, err := repo.GetS3Destination(ctx, first.ID)
	if err != nil {
		t.Fatalf("get verified: %v", err)
	}
	if verified.LastVerifiedAt == nil {
		t.Fatal("expected verification timestamp")
	}
	second.Name = "renamed"
	second.IsDefault = false
	if err := repo.UpdateS3Destination(ctx, second); err != nil {
		t.Fatalf("update: %v", err)
	}
	list, err := repo.ListS3Destinations(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 || list[0].ID != first.ID {
		t.Fatalf("expected default destination first: %+v", list)
	}
	if err := repo.DeleteS3Destination(ctx, second.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
