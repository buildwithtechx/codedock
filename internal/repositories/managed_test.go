package repositories

import (
	"context"
	"strings"
	"testing"

	"codedock/internal/models"
)

type managedTestVault struct{}

func (managedTestVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (managedTestVault) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "enc:") {
		return "", context.DeadlineExceeded
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

func TestManagedCredentialEncryptedAtRest(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewManagedRepository(db, managedTestVault{})
	credential := &models.ManagedCredential{ID: "cred-1", OrganizationID: "org-1", Provider: models.ManagedProviderHetzner, Label: "prod", Token: "hetzner-secret"}
	if err := repo.CreateCredential(ctx, credential); err != nil {
		t.Fatalf("create credential: %v", err)
	}
	var stored string
	if err := db.QueryRowContext(ctx, `SELECT encrypted_token FROM managed_providers WHERE id = $1`, "cred-1").Scan(&stored); err != nil {
		t.Fatalf("read stored token: %v", err)
	}
	if stored == "" || stored == "hetzner-secret" {
		t.Fatalf("expected encrypted token at rest, got %q", stored)
	}
	loaded, err := repo.GetCredential(ctx, "cred-1")
	if err != nil {
		t.Fatalf("get credential: %v", err)
	}
	if loaded.Token != "hetzner-secret" {
		t.Fatalf("expected decrypted token, got %q", loaded.Token)
	}
	listed, err := repo.ListCredentialsByOrg(ctx, "org-1")
	if err != nil {
		t.Fatalf("list credentials: %v", err)
	}
	if len(listed) != 1 || listed[0].Token != "" {
		t.Fatalf("expected redacted credential list, got %#v", listed)
	}
}

func TestManagedCredentialDeleteBlockedWhileLinked(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewManagedRepository(db, managedTestVault{})
	if err := repo.CreateCredential(ctx, &models.ManagedCredential{ID: "cred-1", OrganizationID: "org-1", Provider: models.ManagedProviderHetzner, Label: "prod", Token: "x"}); err != nil {
		t.Fatalf("create credential: %v", err)
	}
	serverRepo := NewServerRepository(db, managedTestVault{})
	if err := serverRepo.Create(ctx, &models.Server{ID: "srv-1", UserID: "u-1", Name: "web", Status: models.ServerStatusProvisioning, Provider: "hetzner"}); err != nil {
		t.Fatalf("create server: %v", err)
	}
	if err := repo.SaveLink(ctx, &models.ManagedServerLink{ServerID: "srv-1", OrganizationID: "org-1", CredentialID: "cred-1", SSHKeyName: "k"}); err != nil {
		t.Fatalf("save link: %v", err)
	}
	if err := repo.DeleteCredential(ctx, "cred-1"); err == nil {
		t.Fatal("expected delete to fail while credential backs a server")
	}
	provisioning, err := repo.ListProvisioningServerIDs(ctx)
	if err != nil {
		t.Fatalf("list provisioning: %v", err)
	}
	if len(provisioning) != 1 || provisioning[0] != "srv-1" {
		t.Fatalf("expected srv-1 provisioning, got %v", provisioning)
	}
	if err := repo.DeleteLink(ctx, "srv-1"); err != nil {
		t.Fatalf("delete link: %v", err)
	}
	if err := repo.DeleteCredential(ctx, "cred-1"); err != nil {
		t.Fatalf("delete credential after unlink: %v", err)
	}
}

func TestManagedQuotaDefaultsAndUpdate(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-9", Name: "Quota"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewManagedRepository(db, managedTestVault{})
	quota, err := repo.GetQuota(ctx, "org-9")
	if err != nil {
		t.Fatalf("get quota: %v", err)
	}
	if quota.MaxServers != 5 || quota.MaxMemoryGB != 32 {
		t.Fatalf("expected default quota 5/32, got %#v", quota)
	}
	quota.MaxServers = 2
	quota.MaxMemoryGB = 8
	if err := repo.SetQuota(ctx, quota); err != nil {
		t.Fatalf("set quota: %v", err)
	}
	loaded, err := repo.GetQuota(ctx, "org-9")
	if err != nil {
		t.Fatalf("reload quota: %v", err)
	}
	if loaded.MaxServers != 2 || loaded.MaxMemoryGB != 8 {
		t.Fatalf("expected updated quota 2/8, got %#v", loaded)
	}
}
