package integration_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	cdkdb "codedock/internal/db"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/projects"
	templatesmod "codedock/internal/modules/templates"
	compose "codedock/internal/platform/docker"
	"codedock/internal/testdb"
	"codedock/internal/utils"
)

type oneclickIntegrationRuntime struct {
	composemod.ComposeStackRuntime
}

type oneClickEnvAdapter struct {
	repo *environments.EnvironmentRepo
}

func (a *oneClickEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*templatesmod.EnvironmentConfig, error) {
	envs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]*templatesmod.EnvironmentConfig, len(envs))
	for i, e := range envs {
		res[i] = &templatesmod.EnvironmentConfig{
			ID:        e.ID,
			ProjectID: e.ProjectID,
			Name:      e.Name,
			IsDefault: e.IsDefault,
		}
	}
	return res, nil
}

func (a *oneClickEnvAdapter) GetDefaultByProject(ctx context.Context, projectID string) (*templatesmod.EnvironmentConfig, error) {
	envs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.IsDefault {
			return &templatesmod.EnvironmentConfig{
				ID:        env.ID,
				ProjectID: env.ProjectID,
				Name:      env.Name,
				IsDefault: env.IsDefault,
			}, nil
		}
	}
	return nil, nil
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

func (r *oneclickIntegrationRuntime) Apply(_ context.Context, _ *composemod.ComposeStack, _ func(string) error) error {
	return nil
}

func (r *oneclickIntegrationRuntime) Results(_ context.Context, _ string) ([]composemod.ComposeServiceResult, error) {
	return nil, nil
}

func newOneClickIntegrationService(t *testing.T) (*templatesmod.OneClickService, *composemod.ComposeStackRepo, string) {
	db := testdb.Open(t)
	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	vlt, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-catalog', 'Catalog') ON CONFLICT (id) DO NOTHING`); err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	envRepo := environments.NewEnvironmentRepo(db)
	projectRepo := projects.NewProjectRepo(db, envRepo)
	project := &projects.ProjectConfig{ID: uuid.NewString(), OrganizationID: "org-catalog", Name: "Catalog Integration"}
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
	store := composemod.NewComposeStackRepo(db, vlt)
	stacks := composemod.NewComposeStackService(store, &oneclickIntegrationRuntime{})
	service := templatesmod.NewOneClickService(tmplMgr, stacks, &oneClickEnvAdapter{repo: envRepo})
	return service, store, project.ID
}

func TestOneClickInstallPersistsStack(t *testing.T) {
	service, store, projectID := newOneClickIntegrationService(t)
	ctx := context.Background()

	preview, err := service.ReviewInstall(ctx, templatesmod.InstallAppInput{
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

	installReview, err := service.ReviewInstall(ctx, templatesmod.InstallAppInput{
		AppID:     "n8n",
		ProjectID: projectID,
		Name:      "automation",
	})
	if err != nil {
		t.Fatalf("expected install review to succeed: %v", err)
	}
	result, err := service.InstallApp(ctx, templatesmod.InstallAppInput{
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

	_, err = service.InstallApp(ctx, templatesmod.InstallAppInput{
		AppID:     "n8n",
		ProjectID: projectID,
		Name:      "automation",
		Digest:    "stale-digest",
	})
	if err == nil || !strings.Contains(err.Error(), "review again") {
		t.Fatalf("expected stale digest rejection, got %v", err)
	}
}
