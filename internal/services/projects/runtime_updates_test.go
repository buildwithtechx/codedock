package projects

import (
	"context"
	"testing"

	"codedock/internal/models"
)

type mockAppRepoForUpdates struct {
	apps map[string]*models.AppService
}

func (m *mockAppRepoForUpdates) Create(ctx context.Context, app *models.AppService) error {
	m.apps[app.ID] = app
	return nil
}
func (m *mockAppRepoForUpdates) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	return m.apps[id], nil
}
func (m *mockAppRepoForUpdates) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	var list []*models.AppService
	for _, a := range m.apps {
		if a.ProjectID == projectID {
			list = append(list, a)
		}
	}
	return list, nil
}
func (m *mockAppRepoForUpdates) ListByEnvironment(ctx context.Context, envID string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) ListByOrganization(ctx context.Context, orgID string, search string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) ListAll(ctx context.Context, search string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) Update(ctx context.Context, app *models.AppService) error {
	m.apps[app.ID] = app
	return nil
}
func (m *mockAppRepoForUpdates) Delete(ctx context.Context, id string) error {
	delete(m.apps, id)
	return nil
}
func (m *mockAppRepoForUpdates) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	return nil
}
func (m *mockAppRepoForUpdates) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	return nil
}
func (m *mockAppRepoForUpdates) CreateLogDrain(ctx context.Context, d *models.LogDrain) error {
	return nil
}
func (m *mockAppRepoForUpdates) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	return nil, nil
}
func (m *mockAppRepoForUpdates) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	return nil
}

func TestCheckRuntimeUpdates(t *testing.T) {
	repo := &mockAppRepoForUpdates{apps: make(map[string]*models.AppService)}
	repo.apps["svc-1"] = &models.AppService{
		ID:        "svc-1",
		ProjectID: "proj-1",
		Name:      "web-service",
		ImageRef:  "node:20-alpine",
	}
	repo.apps["svc-2"] = &models.AppService{
		ID:        "svc-2",
		ProjectID: "proj-1",
		Name:      "db-service",
		ImageRef:  "postgres@sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	svc := NewAppService(repo, nil, nil)
	ctx := context.Background()

	result, err := svc.CheckRuntimeUpdates(ctx, "proj-1", "")
	if err != nil {
		t.Fatalf("CheckRuntimeUpdates error: %v", err)
	}

	if len(result.Images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(result.Images))
	}

	for _, img := range result.Images {
		if *img.ServiceID == "svc-1" {
			if img.Status != models.RuntimeImageStatusUpdateAvailable {
				t.Errorf("expected svc-1 to have status update_available, got %v", img.Status)
			}
		}
		if *img.ServiceID == "svc-2" {
			if img.Status != models.RuntimeImageStatusPinned {
				t.Errorf("expected svc-2 to have status pinned, got %v", img.Status)
			}
		}
	}
}
