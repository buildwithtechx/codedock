package projects

import (
	"context"
	"strings"
	"testing"

	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type oneClickTestEnvRepo struct {
	repositories.EnvironmentRepository
	envs []models.EnvironmentConfig
}

func (r *oneClickTestEnvRepo) ListByProject(_ context.Context, _ string) ([]models.EnvironmentConfig, error) {
	return r.envs, nil
}

type oneClickTestStackStore struct {
	ComposeStackStore
	stack *models.ComposeStack
}

func (s *oneClickTestStackStore) Get(_ context.Context, _, _ string) (*models.ComposeStack, error) {
	return s.stack, nil
}

func (s *oneClickTestStackStore) Save(_ context.Context, stack *models.ComposeStack, _ int) error {
	s.stack = stack
	return nil
}

func (s *oneClickTestStackStore) Claim(_ context.Context, _, _ string, _ int) error { return nil }

func (s *oneClickTestStackStore) Observe(_ context.Context, _, _, _ string, _ string) error {
	return nil
}

type oneClickTestStackRuntime struct{ ComposeStackRuntime }

func (r *oneClickTestStackRuntime) Canonical(_ context.Context, _ string, _ string, _ map[string]string) (string, error) {
	return `{"services":{"plausible":{"image":"plausible/analytics"},"db":{"image":"postgres:16-alpine"}}}`, nil
}

func (r *oneClickTestStackRuntime) ValidatePorts(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *oneClickTestStackRuntime) Apply(_ context.Context, _ *models.ComposeStack, progress func(string) error) error {
	return progress("READY")
}

func (r *oneClickTestStackRuntime) Results(_ context.Context, _ string) ([]models.ComposeServiceResult, error) {
	return nil, nil
}

func oneClickTestService(t *testing.T) (*OneClickService, *oneClickTestStackStore) {
	t.Helper()
	manager, err := compose.NewTemplateManager()
	if err != nil {
		t.Fatal(err)
	}
	stackStore := &oneClickTestStackStore{stack: &models.ComposeStack{Revision: 1}}
	stacks := NewComposeStackService(stackStore, &oneClickTestStackRuntime{})
	envRepo := &oneClickTestEnvRepo{envs: []models.EnvironmentConfig{{ID: "env-1", ProjectID: "project-1"}}}
	return NewOneClickService(manager, stacks, envRepo), stackStore
}

func TestOneClickListAppsExposesInstallMetadata(t *testing.T) {
	service, _ := oneClickTestService(t)
	apps := service.ListApps()
	if len(apps) == 0 {
		t.Fatal("expected catalogue apps")
	}
	for _, app := range apps {
		if app.ID == "" || app.Name == "" || app.DockerImage == "" || len(app.Services) == 0 {
			t.Errorf("app %q misses install metadata: %+v", app.ID, app)
		}
	}
}

func TestOneClickReviewMasksSecrets(t *testing.T) {
	service, _ := oneClickTestService(t)
	preview, err := service.ReviewInstall(context.Background(), models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: "project-1",
		Name:      "Automation",
		HostPort:  5678,
	})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Digest == "" || preview.Kind == "" || len(preview.Services) == 0 {
		t.Fatalf("incomplete preview: %+v", preview)
	}
	for _, service := range preview.Services {
		for _, entry := range service.Env {
			if strings.Contains(entry, "5678:") {
				t.Errorf("preview leaks host mapping in env: %q", entry)
			}
		}
	}
}

func TestOneClickInstallsSingleServiceAsStack(t *testing.T) {
	service, stackStore := oneClickTestService(t)
	preview, err := service.ReviewInstall(context.Background(), models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: "project-1",
		Name:      "Automation",
		Domain:    "n8n.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Kind != "stack" {
		t.Fatalf("expected stack preview, got %+v", preview)
	}
	result, err := service.InstallApp(context.Background(), models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: "project-1",
		Name:      "Automation",
		Domain:    "n8n.example.com",
		Digest:    preview.Digest,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Kind != "stack" || result.Stack == nil {
		t.Fatalf("expected stack result, got %+v", result)
	}
	if stackStore.stack == nil || stackStore.stack.EnvironmentID != "env-1" {
		t.Fatalf("expected saved stack, got %+v", stackStore.stack)
	}
}

func TestOneClickInstallsMultiServiceAsStack(t *testing.T) {
	service, stackStore := oneClickTestService(t)
	preview, err := service.ReviewInstall(context.Background(), models.InstallAppInput{
		AppID:     "plausible",
		ProjectID: "project-1",
		Name:      "Analytics",
	})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Kind != "stack" {
		t.Fatalf("expected stack preview, got %+v", preview)
	}
	result, err := service.InstallApp(context.Background(), models.InstallAppInput{
		AppID:     "plausible",
		ProjectID: "project-1",
		Name:      "Analytics",
		Digest:    preview.Digest,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Kind != "stack" || result.Stack == nil {
		t.Fatalf("expected stack result, got %+v", result)
	}
	if stackStore.stack == nil || !strings.Contains(stackStore.stack.Config, "plausible") {
		t.Fatalf("expected saved stack config, got %+v", stackStore.stack)
	}
}

func TestOneClickInstallRejectsChangedInputs(t *testing.T) {
	service, _ := oneClickTestService(t)
	preview, err := service.ReviewInstall(context.Background(), models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: "project-1",
		Name:      "Automation",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.InstallApp(context.Background(), models.InstallAppInput{
		AppID:     "n8n",
		ProjectID: "project-1",
		Name:      "Automation",
		HostPort:  9999,
		Digest:    preview.Digest,
	})
	if err == nil {
		t.Fatal("expected digest mismatch error")
	}
}
