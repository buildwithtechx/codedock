package integration_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/testdb"
	"codedock.run/codedock/internal/utils"
)

type oneclickIntegrationRuntime struct {
	projectservices.ComposeStackRuntime
}

func (r *oneclickIntegrationRuntime) Canonical(_ context.Context, _ string, source string, _ map[string]string) (string, error) {
	rendered, err := json.Marshal(map[string]any{"services": map[string]any{"stack": map[string]any{"source": source}}})
	if err != nil {
		return "", err
	}
	return string(rendered), nil
}

func (r *oneclickIntegrationRuntime) ValidatePorts(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *oneclickIntegrationRuntime) Apply(_ context.Context, _ *models.ComposeStack, _ func(string) error) error {
	return nil
}

func (r *oneclickIntegrationRuntime) Results(_ context.Context, _ string) ([]models.ComposeServiceResult, error) {
	return nil, nil
}

func newOneClickIntegrationService(t *testing.T) (*projectservices.OneClickService, *repositories.ComposeStackRepo, string) {
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	vlt, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-catalog', 'Catalog') ON CONFLICT (id) DO NOTHING`); err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	envRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, envRepo)
	project := &models.ProjectConfig{ID: uuid.NewString(), OrganizationID: "org-catalog", Name: "Catalog Integration"}
	if err := projectRepo.Create(context.Background(), project); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	envs, err := envRepo.ListByProject(context.Background(), project.ID)
	if err != nil || len(envs) == 0 {
		t.Fatalf("expected default environment, got %v %v", envs, err)
	}
	tmplMgr, err := compose.NewTemplateManager()
	if err != nil {
		t.Fatalf("failed to load templates: %v", err)
	}
	store := repositories.NewComposeStackRepo(db, vlt)
	stacks := projectservices.NewComposeStackService(store, &oneclickIntegrationRuntime{})
	service := projectservices.NewOneClickService(tmplMgr, stacks, envRepo)
	return service, store, project.ID
}

func TestOneClickInstallPersistsStack(t *testing.T) {
	service, store, projectID := newOneClickIntegrationService(t)
	ctx := context.Background()

	preview, err := service.ReviewInstall(ctx, models.InstallAppInput{
		AppID:       "plausible",
		ProjectID:   projectID,
		Name:        "stats",
		Environment: map[string]string{"BASE_URL": "https://stats.example.com"},
	})
	if err != nil {
		t.Fatalf("expected review to succeed: %v", err)
	}
	if !strings.Contains(preview.ComposeYAML, "https://stats.example.com") {
		t.Fatalf("expected resolved input in document:\n%s", preview.ComposeYAML)
	}
	if !strings.Contains(preview.ComposeYAML, "healthcheck:") || !strings.Contains(preview.ComposeYAML, "pg_isready") {
		t.Fatalf("expected rendered healthcheck in document:\n%s", preview.ComposeYAML)
	}

	installReview, err := service.ReviewInstall(ctx, models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: projectID,
		Name:      "automation",
	})
	if err != nil {
		t.Fatalf("expected install review to succeed: %v", err)
	}
	result, err := service.InstallApp(ctx, models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: projectID,
		Name:      "automation",
		Digest:    installReview.Digest,
	})
	if err != nil {
		t.Fatalf("expected digest-pinned install to succeed: %v", err)
	}
	if result.Stack == nil || result.Stack.ID == "" {
		t.Fatalf("expected persisted stack, got %+v", result)
	}
	stored, err := store.Get(ctx, projectID, result.Stack.ID)
	if err != nil {
		t.Fatalf("expected stack row in postgres: %v", err)
	}
	if !strings.Contains(stored.Config, "n8nio/n8n") {
		t.Fatalf("expected rendered image in stored stack, got:\n%s", stored.Config)
	}

	_, err = service.InstallApp(ctx, models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: projectID,
		Name:      "automation",
		Digest:    "stale-digest",
	})
	if err == nil || !strings.Contains(err.Error(), "review again") {
		t.Fatalf("expected stale digest rejection, got %v", err)
	}
}
