package attention

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"codedock.run/codedock/internal/engine/dockerprobe"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/migrations"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/testdb"
)

type fakeHealth struct {
	states map[string]bool
}

func (f *fakeHealth) GetHealth(ctx context.Context, containerID string) (AttentionHealth, error) {
	if running, ok := f.states[containerID]; ok {
		if running {
			return AttentionHealth{Status: "running", Running: true}, nil
		}
		return AttentionHealth{Status: "stopped", Running: false}, nil
	}
	return AttentionHealth{Status: "unknown"}, fmt.Errorf("no such container")
}

type fakeBackups struct {
	records   []*models.BackupRecord
	triggered []string
}

func (f *fakeBackups) ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*models.BackupRecord, error) {
	return f.records, nil
}

func (f *fakeBackups) TriggerBackup(ctx context.Context, configID string) (*models.BackupRecord, error) {
	f.triggered = append(f.triggered, configID)
	return &models.BackupRecord{ID: "rec-new", BackupConfigID: configID}, nil
}

type fakeWorkloads struct {
	restarted  []string
	redeployed []string
}

func (f *fakeWorkloads) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	f.redeployed = append(f.redeployed, appID)
	return "container-new", nil
}

func (f *fakeWorkloads) RestartAppService(ctx context.Context, app *models.AppService) error {
	f.restarted = append(f.restarted, app.ID)
	return nil
}

type errRunner struct{}

func (errRunner) Run(ctx context.Context, cmd string) (string, error) {
	return "", fmt.Errorf("no docker in test")
}

func (errRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	return fmt.Errorf("no docker in test")
}

func setupAttentionTest(t *testing.T) (*Service, *fakeHealth, *fakeBackups, *fakeWorkloads, *repositories.DeploymentRepo, *repositories.AppServiceRepo, repositories.MigrationRepository, repositories.ManagedRepository, repositories.ServerRepository, string) {
	t.Helper()
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	ctx := context.Background()
	userRepo := repositories.NewUserRepo(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	envRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, envRepo)
	appRepo := repositories.NewAppServiceRepo(db)
	deployRepo := repositories.NewDeploymentRepo(db)
	attentionRepo := repositories.NewAttentionRepository(db)
	migrationRepo := repositories.NewMigrationRepository(db, nil)
	managedRepo := repositories.NewManagedRepository(db, nil)
	serverRepo := repositories.NewServerRepository(db, nil)
	for _, user := range []*models.User{
		{ID: "admin", Email: "admin@example.com", Role: models.UserRoleMember, IsActive: true},
		{ID: "member", Email: "member@example.com", Role: models.UserRoleMember, IsActive: true},
	} {
		if err := userRepo.CreateUser(ctx, user); err != nil {
			t.Fatalf("create user: %v", err)
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
			t.Fatalf("add member: %v", err)
		}
	}
	if err := projectRepo.Create(ctx, &models.ProjectConfig{ID: "proj-1", OrganizationID: "org-1", Name: "shop", CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create project: %v", err)
	}
	appService := projectservices.NewAppService(appRepo, repositories.NewServiceVarRepo(db), repositories.NewServiceVolumeRepo(db))
	migrationService := migrations.NewService(migrationRepo, projectRepo, appService, appRepo, repositories.NewServiceVolumeRepo(db), envRepo, serverRepo, userRepo, orgRepo, nil, projectservices.NewComposeParserService(), nil, nil, nil,
		func(*models.MigrationSource) dockerprobe.Runner { return errRunner{} },
		func(*models.Server) dockerprobe.Runner { return errRunner{} })
	health := &fakeHealth{states: map[string]bool{}}
	backups := &fakeBackups{}
	workloads := &fakeWorkloads{}
	service := NewService(attentionRepo, deployRepo, appRepo, projectRepo, userRepo, orgRepo, backups, migrationService, managedRepo, serverRepo, health, workloads)
	return service, health, backups, workloads, deployRepo, appRepo, migrationRepo, managedRepo, serverRepo, "org-1"
}

func seedAttentionService(t *testing.T, appRepo *repositories.AppServiceRepo, projectID, containerID string) string {
	t.Helper()
	svc := &models.AppService{ProjectID: projectID, Name: "web", Status: models.AppServiceStatusRunning, ContainerID: containerID}
	if err := appRepo.Create(context.Background(), svc); err != nil {
		t.Fatalf("create service: %v", err)
	}
	return svc.ID
}

func TestAttentionEvaluatesAndActs(t *testing.T) {
	service, health, backups, workloads, deployRepo, appRepo, _, _, _, orgID := setupAttentionTest(t)
	ctx := context.Background()
	svcID := seedAttentionService(t, appRepo, "proj-1", "c-bad")
	health.states["c-bad"] = false
	if err := deployRepo.Create(ctx, &models.Deployment{OrganizationID: orgID, ProjectID: "proj-1", ServiceID: svcID, Status: models.DeploymentStatusFailed}); err != nil {
		t.Fatalf("create deployment: %v", err)
	}
	backups.records = []*models.BackupRecord{{ID: "rec-1", BackupConfigID: "cfg-1", Status: "failed", Logs: "disk full"}}
	if err := service.Evaluate(ctx, "member", orgID); err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	issues, err := service.ListIssues(ctx, "member", orgID, "open")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	kinds := map[string]*models.AttentionIssue{}
	for _, issue := range issues {
		kinds[issue.Kind] = issue
	}
	if kinds["deployment-failed"] == nil || kinds["service-unhealthy"] == nil || kinds["backup-failed"] == nil {
		t.Fatalf("expected deployment/workload/backup issues, got %v", kinds)
	}
	outcome, err := service.Act(ctx, "member", orgID, kinds["service-unhealthy"].ID, models.AttentionActionRequest{})
	if err != nil {
		t.Fatalf("restart action: %v", err)
	}
	if outcome == "" || len(workloads.restarted) != 1 {
		t.Fatalf("expected restart, got %q %v", outcome, workloads.restarted)
	}
	if _, err := service.Act(ctx, "member", orgID, kinds["deployment-failed"].ID, models.AttentionActionRequest{}); err != nil {
		t.Fatalf("redeploy action: %v", err)
	}
	if _, err := service.Act(ctx, "member", orgID, kinds["backup-failed"].ID, models.AttentionActionRequest{}); err != nil {
		t.Fatalf("retry backup action: %v", err)
	}
	if len(backups.triggered) != 1 {
		t.Fatalf("expected backup retry, got %v", backups.triggered)
	}
	if err := service.Acknowledge(ctx, "member", orgID, kinds["deployment-failed"].ID); err != nil {
		t.Fatalf("ack: %v", err)
	}
	if err := service.Resolve(ctx, "member", orgID, kinds["deployment-failed"].ID); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if _, err := service.Act(ctx, "member", orgID, kinds["deployment-failed"].ID, models.AttentionActionRequest{}); err == nil {
		t.Fatal("expected action on resolved issue to fail")
	}
}

func TestAttentionMigrationAndQuota(t *testing.T) {
	service, _, _, _, _, _, migrationRepo, managedRepo, serverRepo, orgID := setupAttentionTest(t)
	ctx := context.Background()
	if err := migrationRepo.CreateRun(ctx, &models.MigrationRun{ID: "run-1", OrganizationID: orgID, UserID: "member", SourceKind: "external", Mode: models.MigrationModeMove, Status: models.MigrationStatusFailed, Phase: "TRANSFER", Error: "boom"}); err != nil {
		t.Fatalf("seed failed run: %v", err)
	}
	if err := managedRepo.SetQuota(ctx, &models.ManagedQuota{OrganizationID: orgID, MaxServers: 1, MaxMemoryGB: 32}); err != nil {
		t.Fatalf("seed quota: %v", err)
	}
	if err := serverRepo.Create(ctx, &models.Server{ID: "srv-1", UserID: "member", Name: "m1", Status: models.ServerStatusOnline, Provider: "hetzner", ServerType: "cx22"}); err != nil {
		t.Fatalf("seed server: %v", err)
	}
	if err := managedRepo.CreateCredential(ctx, &models.ManagedCredential{ID: "cred-1", OrganizationID: orgID, Provider: models.ManagedProviderHetzner, Label: "t", Token: "x"}); err != nil {
		t.Fatalf("seed credential: %v", err)
	}
	if err := managedRepo.SaveLink(ctx, &models.ManagedServerLink{ServerID: "srv-1", OrganizationID: orgID, CredentialID: "cred-1"}); err != nil {
		t.Fatalf("seed link: %v", err)
	}
	if err := service.Evaluate(ctx, "member", orgID); err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	issues, err := service.ListIssues(ctx, "member", orgID, "open")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	kinds := map[string]*models.AttentionIssue{}
	for _, issue := range issues {
		kinds[issue.Kind+":"+issue.Subject] = issue
	}
	if kinds["migration-failed:run-1"] == nil || kinds["quota-pressure:managed-servers"] == nil {
		t.Fatalf("expected migration and quota issues, got %v", kinds)
	}
	outcome, err := service.Act(ctx, "member", orgID, kinds["migration-failed:run-1"].ID, models.AttentionActionRequest{})
	if err != nil {
		t.Fatalf("resume action: %v", err)
	}
	if outcome == "" {
		t.Fatal("expected resume outcome")
	}
}
