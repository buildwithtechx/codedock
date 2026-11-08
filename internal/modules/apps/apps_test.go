package apps

import (
	"context"
	"testing"
)

type mockAppRepo struct {
	apps []*AppService
}

func (m *mockAppRepo) Create(ctx context.Context, svc *AppService) error {
	m.apps = append(m.apps, svc)
	return nil
}

func (m *mockAppRepo) GetByID(ctx context.Context, id string) (*AppService, error) {
	for _, a := range m.apps {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, ErrAppServiceNotFound
}

func (m *mockAppRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error) {
	var res []*AppService
	for _, a := range m.apps {
		if a.EnvironmentID == environmentID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *mockAppRepo) ListByProject(ctx context.Context, projectID string) ([]*AppService, error) {
	var res []*AppService
	for _, a := range m.apps {
		if a.ProjectID == projectID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *mockAppRepo) ListByOrganization(ctx context.Context, organizationID string, search string) ([]*AppService, error) {
	return m.apps, nil
}

func (m *mockAppRepo) ListAll(ctx context.Context, search string) ([]*AppService, error) {
	return m.apps, nil
}

func (m *mockAppRepo) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	return nil, nil
}

func (m *mockAppRepo) Update(ctx context.Context, svc *AppService) error {
	for i, a := range m.apps {
		if a.ID == svc.ID {
			m.apps[i] = svc
			return nil
		}
	}
	return nil
}

func (m *mockAppRepo) Delete(ctx context.Context, id string) error {
	for i, a := range m.apps {
		if a.ID == id {
			m.apps = append(m.apps[:i], m.apps[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockAppRepo) CreateWebhook(ctx context.Context, w *Webhook) error {
	return nil
}

func (m *mockAppRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*Webhook, error) {
	return nil, nil
}

func (m *mockAppRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	return nil
}

func (m *mockAppRepo) CreateLogDrain(ctx context.Context, d *LogDrain) error {
	return nil
}

func (m *mockAppRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*LogDrain, error) {
	return nil, nil
}

func (m *mockAppRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	return nil
}

func TestAppServiceCRUD(t *testing.T) {
	repo := &mockAppRepo{}
	svc := NewService(repo, nil, nil)
	ctx := context.Background()

	// Validation
	if _, err := svc.CreateAppService(ctx, nil); err == nil {
		t.Fatal("expected error on nil app service")
	}

	created, err := svc.CreateAppService(ctx, &AppService{
		ProjectID:     "proj-123",
		EnvironmentID: "env-123",
		Name:          "web-api",
		Branch:        "main",
		RuntimeMode:   RuntimeModeStandard,
	})
	if err != nil {
		t.Fatalf("unexpected error creating app service: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated UUID")
	}

	fetched, err := svc.GetAppService(ctx, created.ID)
	if err != nil || fetched.Name != "web-api" {
		t.Fatalf("failed to get app service: %v", err)
	}

	list, err := svc.ListByProject(ctx, "proj-123")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 service, got %d", len(list))
	}

	if err := svc.DeleteAppService(ctx, created.ID); err != nil {
		t.Fatalf("failed to delete service: %v", err)
	}
	listAfter, _ := svc.ListByProject(ctx, "proj-123")
	if len(listAfter) != 0 {
		t.Fatalf("expected 0 services after delete, got %d", len(listAfter))
	}
}

func TestAutoscalingPolicyDefaults(t *testing.T) {
	policy := DefaultAutoscalingPolicy("svc-abc")
	if policy.MinReplicas != 1 || policy.MaxReplicas != 5 {
		t.Fatalf("unexpected default replicas: min=%d, max=%d", policy.MinReplicas, policy.MaxReplicas)
	}
	if policy.ScaleUpCPU != 80 || policy.ScaleDownCPU != 20 {
		t.Fatalf("unexpected default CPU triggers: up=%v, down=%v", policy.ScaleUpCPU, policy.ScaleDownCPU)
	}
}
