package system

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/operations"
	"codedock.run/codedock/internal/testdb"
)

type managedTestVault struct{}

func (managedTestVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (managedTestVault) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "enc:") {
		return "", fmt.Errorf("not encrypted")
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

type fakeProvider struct {
	mu        sync.Mutex
	servers   map[string]*models.ManagedInstance
	keys      map[string]string
	nextID    int
	creates   int
	changes   []string
	powerOff  int
	powerOn   int
	createErr error
}

func newFakeProvider() *fakeProvider {
	return &fakeProvider{servers: map[string]*models.ManagedInstance{}, keys: map[string]string{}, nextID: 1001}
}

func (f *fakeProvider) CreateServer(ctx context.Context, spec models.ManagedServerSpec) (models.ManagedInstance, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.createErr != nil {
		return models.ManagedInstance{}, f.createErr
	}
	id := fmt.Sprintf("%d", f.nextID)
	f.nextID++
	instance := &models.ManagedInstance{ExternalID: id, Name: spec.Name, Status: hetzner.StatusRunning, PublicIP: "203.0.113.10", ServerType: spec.ServerType, Region: spec.Region, Labels: spec.Labels}
	f.servers[id] = instance
	f.creates++
	return *instance, nil
}

func (f *fakeProvider) GetServer(ctx context.Context, externalID string) (models.ManagedInstance, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	instance, ok := f.servers[externalID]
	if !ok {
		return models.ManagedInstance{}, &hetzner.APIError{Status: 404, Code: "not_found", Message: "missing"}
	}
	return *instance, nil
}

func (f *fakeProvider) FindServerByLabel(ctx context.Context, selector string) (models.ManagedInstance, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	value := selector
	if parts := strings.SplitN(selector, "==", 2); len(parts) == 2 {
		value = parts[1]
	}
	for _, instance := range f.servers {
		if instance.Labels["codedock-server"] == value {
			return *instance, nil
		}
	}
	return models.ManagedInstance{}, nil
}

func (f *fakeProvider) DeleteServer(ctx context.Context, externalID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.servers[externalID]; !ok {
		return &hetzner.APIError{Status: 404, Code: "not_found", Message: "missing"}
	}
	delete(f.servers, externalID)
	return nil
}

func (f *fakeProvider) PowerOn(ctx context.Context, externalID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.powerOn++
	if instance, ok := f.servers[externalID]; ok {
		instance.Status = hetzner.StatusRunning
	}
	return nil
}

func (f *fakeProvider) PowerOff(ctx context.Context, externalID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.powerOff++
	if instance, ok := f.servers[externalID]; ok {
		instance.Status = hetzner.StatusOff
	}
	return nil
}

func (f *fakeProvider) ChangeType(ctx context.Context, externalID, serverType string, upgradeDisk bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	instance, ok := f.servers[externalID]
	if !ok {
		return &hetzner.APIError{Status: 404, Code: "not_found", Message: "missing"}
	}
	instance.ServerType = serverType
	f.changes = append(f.changes, externalID+":"+serverType)
	return nil
}

func (f *fakeProvider) WaitForStatus(ctx context.Context, externalID string, wanted []string, timeout time.Duration) (models.ManagedInstance, error) {
	return f.GetServer(ctx, externalID)
}

func (f *fakeProvider) CreateSSHKey(ctx context.Context, name, publicKey string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.keys[name] = publicKey
	return nil
}

func (f *fakeProvider) DeleteSSHKey(ctx context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.keys, name)
	return nil
}

type managedTestEnv struct {
	db          *sql.DB
	service     ManagedService
	ops         *operations.Service
	serverRepo  repositories.ServerRepository
	managedRepo repositories.ManagedRepository
	fake        *fakeProvider
	orgID       string
	adminID     string
	memberID    string
	outsiderID  string
	credID      string
}

func setupManagedTest(t *testing.T) *managedTestEnv {
	t.Helper()
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	ctx := context.Background()
	vault := managedTestVault{}
	userRepo := repositories.NewUserRepo(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	serverRepo := repositories.NewServerRepository(db, vault)
	managedRepo := repositories.NewManagedRepository(db, vault)
	ops := operations.NewService(repositories.NewOperationRepo(db, vault))
	for _, user := range []*models.User{
		{ID: "admin", Email: "admin@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "pro"},
		{ID: "member", Email: "member@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "free"},
		{ID: "outsider", Email: "outsider@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "free"},
	} {
		if err := userRepo.CreateUser(ctx, user); err != nil {
			t.Fatalf("create user %s: %v", user.ID, err)
		}
	}
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	now := time.Now().UTC()
	for _, member := range []*models.OrganizationMember{
		{ID: "m-admin", OrganizationID: "org-1", UserID: "admin", Email: "admin@example.com", Permission: models.MemberPermissionAdmin, Status: models.MemberStatusActive, InvitedAt: now, AcceptedAt: now},
		{ID: "m-member", OrganizationID: "org-1", UserID: "member", Email: "member@example.com", Permission: models.MemberPermissionMember, Status: models.MemberStatusActive, InvitedAt: now, AcceptedAt: now},
	} {
		if err := orgRepo.AddMember(ctx, member); err != nil {
			t.Fatalf("add member %s: %v", member.UserID, err)
		}
	}
	fake := newFakeProvider()
	service := NewManagedService(managedRepo, serverRepo, orgRepo, userRepo, nil, ops, func(token string) ManagedProvider { return fake })
	env := &managedTestEnv{db: db, service: service, ops: ops, serverRepo: serverRepo, managedRepo: managedRepo, fake: fake, orgID: "org-1", adminID: "admin", memberID: "member", outsiderID: "outsider"}
	credential, err := service.CreateCredential(ctx, "admin", "org-1", models.CreateManagedCredentialRequest{Provider: models.ManagedProviderHetzner, Label: "prod", Token: "hetzner-secret"})
	if err != nil {
		t.Fatalf("create credential: %v", err)
	}
	env.credID = credential.ID
	return env
}

func waitManagedOperation(t *testing.T, ops *operations.Service, id string) *models.Operation {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		op, err := ops.Get(context.Background(), id)
		if err != nil {
			t.Fatalf("load operation: %v", err)
		}
		if op.Status == "COMPLETED" || op.Status == "FAILED" || op.Status == "INTERRUPTED" {
			return op
		}
		if time.Now().After(deadline) {
			t.Fatalf("operation %s stuck in %s", id, op.Status)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func provisionRequest(credID string) models.ReviewManagedProvisionRequest {
	return models.ReviewManagedProvisionRequest{CredentialID: credID, Name: "web-1", Region: "fsn1", ServerType: "cx22", Image: "ubuntu-24.04"}
}

func TestManagedCredentialRequiresOrgAdmin(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	req := models.CreateManagedCredentialRequest{Provider: models.ManagedProviderHetzner, Label: "extra", Token: "x"}
	if _, err := env.service.CreateCredential(ctx, env.outsiderID, env.orgID, req); err == nil {
		t.Fatal("expected outsider credential creation to fail")
	}
	if _, err := env.service.CreateCredential(ctx, env.memberID, env.orgID, req); err == nil {
		t.Fatal("expected member credential creation to fail")
	}
	created, err := env.service.CreateCredential(ctx, env.adminID, env.orgID, req)
	if err != nil {
		t.Fatalf("admin create credential: %v", err)
	}
	if created.Token != "" {
		t.Fatal("expected created credential token to be redacted")
	}
	if _, err := env.service.ListCredentials(ctx, env.outsiderID, env.orgID); err == nil {
		t.Fatal("expected outsider credential list to fail")
	}
	listed, err := env.service.ListCredentials(ctx, env.memberID, env.orgID)
	if err != nil {
		t.Fatalf("member list credentials: %v", err)
	}
	if len(listed) != 2 {
		t.Fatalf("expected 2 credentials, got %d", len(listed))
	}
	if err := env.service.DeleteCredential(ctx, env.memberID, env.orgID, created.ID); err == nil {
		t.Fatal("expected member credential delete to fail")
	}
	if err := env.service.DeleteCredential(ctx, env.adminID, env.orgID, created.ID); err != nil {
		t.Fatalf("admin delete credential: %v", err)
	}
}

func TestManagedProvisionValidation(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	for _, test := range []struct {
		name   string
		mutate func(*models.ReviewManagedProvisionRequest)
	}{
		{"name", func(r *models.ReviewManagedProvisionRequest) { r.Name = "-bad name-" }},
		{"region", func(r *models.ReviewManagedProvisionRequest) { r.Region = "moon1" }},
		{"tier", func(r *models.ReviewManagedProvisionRequest) { r.ServerType = "cx99" }},
		{"image", func(r *models.ReviewManagedProvisionRequest) { r.Image = "windows-11" }},
		{"credential", func(r *models.ReviewManagedProvisionRequest) { r.CredentialID = "missing" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := provisionRequest(env.credID)
			test.mutate(&req)
			if _, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, req); err == nil {
				t.Fatalf("expected %s validation to fail", test.name)
			}
		})
	}
	if _, err := env.service.ReviewProvision(ctx, env.outsiderID, env.orgID, provisionRequest(env.credID)); err == nil {
		t.Fatal("expected outsider provision review to fail")
	}
}

func TestManagedProvisionBillingGate(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Cloud.Enabled
	t.Cleanup(func() { cfg.Cloud.Enabled = previous })
	for _, check := range []struct {
		name    string
		cloud   bool
		user    string
		allowed bool
	}{
		{"self-hosted free", false, "member", true},
		{"cloud free", true, "member", false},
		{"cloud pro", true, "admin", true},
	} {
		t.Run(check.name, func(t *testing.T) {
			cfg.Cloud.Enabled = check.cloud
			env := setupManagedTest(t)
			_, err := env.service.ReviewProvision(context.Background(), env.memberID, env.orgID, provisionRequest(env.credID))
			if check.user == "admin" {
				_, err = env.service.ReviewProvision(context.Background(), env.adminID, env.orgID, provisionRequest(env.credID))
			}
			if check.allowed && err != nil {
				t.Fatalf("expected provision review to pass: %v", err)
			}
			if !check.allowed && err == nil {
				t.Fatal("expected provision review to fail")
			}
		})
	}
}

func TestManagedQuotaAdmission(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	if _, err := env.service.SetQuota(ctx, env.memberID, env.orgID, models.ManagedQuotaRequest{MaxServers: 1, MaxMemoryGB: 4}); err == nil {
		t.Fatal("expected member quota update to fail")
	}
	if _, err := env.service.SetQuota(ctx, env.adminID, env.orgID, models.ManagedQuotaRequest{MaxServers: 1, MaxMemoryGB: 4}); err != nil {
		t.Fatalf("set quota: %v", err)
	}
	first := provisionRequest(env.credID)
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, first)
	if err != nil {
		t.Fatalf("review first server: %v", err)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply first server: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected first provision to complete, got %s", op.Status)
	}
	second := provisionRequest(env.credID)
	second.Name = "web-2"
	if _, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, second); err == nil {
		t.Fatal("expected second provision to exceed quota")
	}
	quota, err := env.service.GetQuota(ctx, env.memberID, env.orgID)
	if err != nil {
		t.Fatalf("get quota: %v", err)
	}
	if quota.UsedServers != 1 || quota.UsedMemoryGB != 4 {
		t.Fatalf("expected usage 1 server/4 GB, got %#v", quota)
	}
	if _, err := env.service.SetQuota(ctx, env.adminID, env.orgID, models.ManagedQuotaRequest{MaxServers: 0, MaxMemoryGB: 0}); err == nil {
		t.Fatal("expected quota below usage to fail")
	}
}
