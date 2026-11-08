package deployments

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRegistryRepo struct {
	regs map[string]*Registry
}

func newMockRegistryRepo() *mockRegistryRepo {
	return &mockRegistryRepo{regs: make(map[string]*Registry)}
}

func (m *mockRegistryRepo) Create(ctx context.Context, reg *Registry) error {
	m.regs[reg.ID] = reg
	return nil
}

func (m *mockRegistryRepo) ListByProject(ctx context.Context, projectID string) ([]*Registry, error) {
	var list []*Registry
	for _, r := range m.regs {
		if r.ProjectID == projectID {
			list = append(list, r)
		}
	}
	return list, nil
}

func (m *mockRegistryRepo) Get(ctx context.Context, id string) (*Registry, error) {
	r, ok := m.regs[id]
	if !ok {
		return nil, ErrRegistryNotFound
	}
	return r, nil
}

func (m *mockRegistryRepo) Delete(ctx context.Context, id string) error {
	delete(m.regs, id)
	return nil
}

func TestRollbackAndRetry(t *testing.T) {
	repo := newMockDeploymentRepo()
	svc := NewDeploymentService(repo, nil, nil, nil, nil, nil)

	_ = repo.Create(context.Background(), &Deployment{
		ID:         "dep-success",
		ServiceID:  "svc-1",
		ProjectID:  "proj-1",
		Status:     DeploymentStatusReady,
		CommitHash: "c123456",
		ImageRef:   "myimage:latest",
		Branch:     "main",
	})

	// Rollback
	rollbackDep, err := svc.Rollback(context.Background(), "dep-success", "testing rollback")
	require.NoError(t, err)
	assert.Equal(t, "svc-1", rollbackDep.ServiceID)
	assert.Equal(t, "myimage:latest", rollbackDep.ImageRef)
	assert.Equal(t, "Rollback", rollbackDep.Trigger)

	// Retry
	retryDep, err := svc.Retry(context.Background(), "dep-success")
	require.NoError(t, err)
	assert.Equal(t, "retry", retryDep.Trigger)
	assert.Equal(t, "myimage:latest", retryDep.ImageRef)
}

func TestRetentionService_Validation(t *testing.T) {
	repo := newMockDeploymentRepo()
	now := time.Now()
	_ = repo.Create(context.Background(), &Deployment{
		ID:                 "dep-retained",
		ServiceID:          "svc-1",
		ProjectID:          "proj-1",
		Status:             DeploymentStatusReady,
		ImageRef:           "img:1",
		ArtifactRetainedAt: &now,
	})

	retService := NewRetentionService(repo)
	dep, _ := repo.GetByID(context.Background(), "dep-retained")
	val := retService.ValidateRollback(context.Background(), dep, nil)
	assert.True(t, val.Eligible)
	assert.Equal(t, "img:1", val.RollbackImage)

	// Ineligible status
	failedDep := &Deployment{
		ID:     "dep-failed",
		Status: DeploymentStatusFailed,
	}
	valFailed := retService.ValidateRollback(context.Background(), failedDep, nil)
	assert.False(t, valFailed.Eligible)
	assert.Equal(t, RollbackErrorNotReady, valFailed.ErrorCode)
}

func TestRegistryService_CRUD(t *testing.T) {
	repo := newMockRegistryRepo()
	svc := NewRegistryService(repo)

	reg, err := svc.CreateRegistry(context.Background(), "proj-1", CreateRegistryRequest{
		Name:          "Docker Hub",
		RegistryURL:   "https://index.docker.io/v1/",
		Username:      "dockuser",
		PasswordToken: "secrettoken",
	})
	require.NoError(t, err)
	assert.Equal(t, "Docker Hub", reg.Name)
	assert.Empty(t, reg.PasswordToken) // Cleared for safety

	list, err := svc.ListRegistries(context.Background(), "proj-1")
	require.NoError(t, err)
	assert.Len(t, list, 1)

	err = svc.DeleteRegistry(context.Background(), reg.ID)
	require.NoError(t, err)

	listAfter, _ := svc.ListRegistries(context.Background(), "proj-1")
	assert.Empty(t, listAfter)
}
