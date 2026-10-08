package repositories

import (
	"context"
	"strings"
	"testing"

	"codedock.run/codedock/internal/models"
)

type migrationTestVault struct{}

func (migrationTestVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (migrationTestVault) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "enc:") {
		return "", context.DeadlineExceeded
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

func TestMigrationSourceEncryptedAtRest(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewMigrationRepository(db, migrationTestVault{})
	source := &models.MigrationSource{ID: "src-1", OrganizationID: "org-1", Name: "legacy", SSHHost: "10.0.0.9", SSHPort: 22, SSHUser: "root", SSHAuthMethod: "key", SSHKey: "private-key", Fingerprint: "SHA256:abc"}
	if err := repo.CreateSource(ctx, source); err != nil {
		t.Fatalf("create source: %v", err)
	}
	var stored string
	if err := db.QueryRowContext(ctx, `SELECT ssh_key FROM migration_sources WHERE id = $1`, "src-1").Scan(&stored); err != nil {
		t.Fatalf("read stored key: %v", err)
	}
	if stored == "" || stored == "private-key" {
		t.Fatalf("expected encrypted key at rest, got %q", stored)
	}
	loaded, err := repo.GetSource(ctx, "src-1")
	if err != nil {
		t.Fatalf("get source: %v", err)
	}
	if loaded.SSHKey != "private-key" {
		t.Fatalf("expected decrypted key, got %q", loaded.SSHKey)
	}
	listed, err := repo.ListSourcesByOrg(ctx, "org-1")
	if err != nil {
		t.Fatalf("list sources: %v", err)
	}
	if len(listed) != 1 || listed[0].SSHKey != "" || listed[0].SSHPassword != "" {
		t.Fatalf("expected redacted list, got %#v", listed)
	}
}

func TestMigrationRunLifecycleGuards(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewMigrationRepository(db, migrationTestVault{})
	if err := repo.CreateSource(ctx, &models.MigrationSource{ID: "src-1", OrganizationID: "org-1", Name: "legacy", SSHHost: "h", SSHUser: "root", Fingerprint: "fp"}); err != nil {
		t.Fatalf("create source: %v", err)
	}
	run := &models.MigrationRun{ID: "run-1", OrganizationID: "org-1", UserID: "u", SourceID: "src-1", SourceKind: "external", Mode: models.MigrationModeMove, Status: models.MigrationStatusTransferring}
	if err := repo.CreateRun(ctx, run); err != nil {
		t.Fatalf("create run: %v", err)
	}
	if err := repo.DeleteSource(ctx, "src-1"); err == nil {
		t.Fatal("expected source delete to fail with active run")
	}
	if err := repo.DeleteRun(ctx, "run-1"); err == nil {
		t.Fatal("expected active run delete to fail")
	}
	if err := repo.RequestCancel(ctx, "run-1"); err != nil {
		t.Fatalf("request cancel: %v", err)
	}
	active, err := repo.ActiveRunForSource(ctx, "src-1")
	if err != nil {
		t.Fatalf("active run: %v", err)
	}
	if active == nil || !active.CancelRequested {
		t.Fatalf("expected cancellable active run, got %#v", active)
	}
	if err := repo.UpdateRunStatus(ctx, "run-1", models.MigrationStatusCancelled, "DONE", ""); err != nil {
		t.Fatalf("mark cancelled: %v", err)
	}
	if err := repo.AppendRunLogs(ctx, "run-1", "hello"); err != nil {
		t.Fatalf("append logs: %v", err)
	}
	loaded, err := repo.GetRun(ctx, "run-1")
	if err != nil {
		t.Fatalf("get run: %v", err)
	}
	if loaded.Logs != "hello" || loaded.Status != models.MigrationStatusCancelled {
		t.Fatalf("expected cancelled run with logs, got %#v", loaded)
	}
	actives, err := repo.ListActive(ctx)
	if err != nil {
		t.Fatalf("list active: %v", err)
	}
	if len(actives) != 0 {
		t.Fatalf("expected no active runs, got %d", len(actives))
	}
	if err := repo.DeleteRun(ctx, "run-1"); err != nil {
		t.Fatalf("delete terminal run: %v", err)
	}
	if err := repo.DeleteSource(ctx, "src-1"); err != nil {
		t.Fatalf("delete source: %v", err)
	}
}
