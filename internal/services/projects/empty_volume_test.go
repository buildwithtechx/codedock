package projects

import (
	"context"
	"strings"
	"testing"
	"time"

	"codedock/internal/models"
)

type mockVolumeRepo struct {
	volumes map[string]*models.ServiceVolume
}

func newMockVolumeRepo() *mockVolumeRepo {
	return &mockVolumeRepo{volumes: make(map[string]*models.ServiceVolume)}
}

func (m *mockVolumeRepo) Create(ctx context.Context, v *models.ServiceVolume) error {
	m.volumes[v.ID] = v
	return nil
}

func (m *mockVolumeRepo) GetByID(ctx context.Context, id string) (*models.ServiceVolume, error) {
	return m.volumes[id], nil
}

func (m *mockVolumeRepo) ListByService(ctx context.Context, serviceID string) ([]models.ServiceVolume, error) {
	var list []models.ServiceVolume
	for _, v := range m.volumes {
		if v.ServiceID == serviceID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (m *mockVolumeRepo) Update(ctx context.Context, v *models.ServiceVolume) error {
	m.volumes[v.ID] = v
	return nil
}

func (m *mockVolumeRepo) Delete(ctx context.Context, id string) error {
	delete(m.volumes, id)
	return nil
}

func TestAppServiceEmptyVolume(t *testing.T) {
	repo := newMockVolumeRepo()
	appSvc := NewAppService(nil, nil, repo)
	ctx := context.Background()

	original := &models.ServiceVolume{
		ID:            "vol-123",
		ServiceID:     "svc-456",
		HostPath:      "my-postgres-volume",
		ContainerPath: "/var/lib/postgresql/data",
		CreatedAt:     time.Now(),
	}
	if err := repo.Create(ctx, original); err != nil {
		t.Fatal(err)
	}

	updated, retained, err := appSvc.EmptyVolume(ctx, "svc-456", "vol-123")
	if err != nil {
		t.Fatalf("EmptyVolume failed: %v", err)
	}

	if retained != "my-postgres-volume" {
		t.Errorf("expected retained volume to be 'my-postgres-volume', got '%s'", retained)
	}

	if !strings.HasPrefix(updated.HostPath, "codedock-fresh-") {
		t.Errorf("expected new host path to start with 'codedock-fresh-', got '%s'", updated.HostPath)
	}

	if updated.ContainerPath != "/var/lib/postgresql/data" {
		t.Errorf("container path should remain unchanged, got '%s'", updated.ContainerPath)
	}
}
