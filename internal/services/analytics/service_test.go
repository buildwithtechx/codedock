package analytics

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/testdb"
)

func setupAnalyticsTest(t *testing.T) (*Service, *repositories.DeploymentRepo, *repositories.EnvironmentRepo, string, string) {
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
	trafficRepo := repositories.NewTrafficRepository(db)
	attentionRepo := repositories.NewAttentionRepository(db)
	for _, user := range []*models.User{
		{ID: "admin", Email: "admin@example.com", Role: models.UserRoleMember, IsActive: true},
		{ID: "member", Email: "member@example.com", Role: models.UserRoleMember, IsActive: true},
		{ID: "outsider", Email: "outsider@example.com", Role: models.UserRoleMember, IsActive: true},
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
	project := &models.ProjectConfig{ID: "proj-1", OrganizationID: "org-1", Name: "shop", CreatedAt: now, UpdatedAt: now}
	if err := projectRepo.Create(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	service := NewService(trafficRepo, deployRepo, projectRepo, appRepo, nil, userRepo, orgRepo, attentionRepo, nil)
	return service, deployRepo, envRepo, "org-1", "proj-1"
}

func TestAnalyticsReadsSeededTraffic(t *testing.T) {
	service, _, _, orgID, projectID := setupAnalyticsTest(t)
	ctx := context.Background()
	minute := time.Now().UTC().Truncate(time.Minute).Format("2006-01-02T15:04")
	from := time.Now().UTC().Add(-time.Hour).Format("2006-01-02T15:04")
	to := time.Now().UTC().Add(time.Hour).Format("2006-01-02T15:04")
	_ = minute
	samples := []models.TrafficSample{}
	for i := 0; i < 10; i++ {
		status := 200
		if i >= 8 {
			status = 500
		}
		samples = append(samples, models.TrafficSample{Time: time.Now().UTC().Format(time.RFC3339), ProjectID: projectID, Domain: "shop.example.com", Path: "/cart", Status: status, Bytes: 100, DurationMs: 20, ClientIP: "10.0.0.2"})
	}
	if err := service.traffic.RecordBatch(ctx, samples); err != nil {
		t.Fatalf("record: %v", err)
	}
	summary, err := service.Summary(ctx, "member", orgID, projectID, "", from, to)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Requests != 10 || summary.ErrorRate != 0.2 {
		t.Fatalf("expected 10 requests/0.2 errors, got %#v", summary)
	}
	overview, err := service.Overview(ctx, "member", orgID, projectID, "", from, to, "5")
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if len(overview.Series) == 0 || len(overview.Statuses) != 2 || len(overview.TopPaths) != 1 {
		t.Fatalf("expected overview parts, got %#v", overview)
	}
	geo, err := service.Geo(ctx, "member", orgID, projectID, time.Now().UTC().Format("2006-01-02"), time.Now().UTC().Format("2006-01-02"))
	if err != nil {
		t.Fatalf("geo: %v", err)
	}
	if len(geo.Countries) != 1 || geo.Countries[0].Requests != 10 {
		t.Fatalf("expected geo rollup, got %#v", geo)
	}
	if _, err := service.Summary(ctx, "outsider", orgID, projectID, "", from, to); err == nil {
		t.Fatal("expected outsider to be denied")
	}
	if _, err := service.Summary(ctx, "member", orgID, "missing", "", from, to); err == nil {
		t.Fatal("expected missing project to fail")
	}
}

func TestAnalyticsDeploymentStatsAndDashboard(t *testing.T) {
	service, deployRepo, envRepo, orgID, projectID := setupAnalyticsTest(t)
	ctx := context.Background()
	appRepo := service.apps
	envs, err := envRepo.ListByProject(ctx, projectID)
	if err != nil || len(envs) == 0 {
		t.Fatalf("default environment: %v", err)
	}
	svc := &models.AppService{ProjectID: projectID, EnvironmentID: envs[0].ID, Name: "web", Status: models.AppServiceStatusRunning}
	if err := appRepo.Create(ctx, svc); err != nil {
		t.Fatalf("create service: %v", err)
	}
	finished := time.Now().UTC()
	for _, status := range []models.DeploymentStatus{models.DeploymentStatusActive, models.DeploymentStatusActive, models.DeploymentStatusFailed} {
		record := &models.Deployment{OrganizationID: orgID, ProjectID: projectID, ServiceID: svc.ID, Status: status, FinishedAt: &finished}
		if err := deployRepo.Create(ctx, record); err != nil {
			t.Fatalf("create deployment: %v", err)
		}
	}
	stats, err := service.DeploymentStats(ctx, "member", orgID, projectID, 30)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Total != 3 || stats.Succeeded != 2 || stats.Failed != 1 {
		t.Fatalf("expected 3/2/1 deployment stats, got %#v", stats)
	}
	dashboard, err := service.Dashboard(ctx, "member", orgID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if dashboard.Deployments7d != 3 {
		t.Fatalf("expected 3 recent deployments, got %#v", dashboard)
	}
}

func TestAnalyticsPathsToggleAdmin(t *testing.T) {
	service, _, _, orgID, projectID := setupAnalyticsTest(t)
	ctx := context.Background()
	if err := service.SetPathsEnabled(ctx, "member", orgID, projectID, true); err == nil {
		t.Fatal("expected member toggle to fail")
	}
	if err := service.SetPathsEnabled(ctx, "admin", orgID, projectID, true); err != nil {
		t.Fatalf("admin toggle: %v", err)
	}
	enabled, err := service.PathsEnabled(ctx, "member", orgID, projectID)
	if err != nil || !enabled {
		t.Fatalf("expected paths enabled, got %v %v", enabled, err)
	}
}
