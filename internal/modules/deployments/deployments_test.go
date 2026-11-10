package deployments

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDeploymentRepo struct {
	mu   sync.RWMutex
	deps map[string]*Deployment
}

func newMockDeploymentRepo() *mockDeploymentRepo {
	return &mockDeploymentRepo{deps: make(map[string]*Deployment)}
}

func (m *mockDeploymentRepo) Create(ctx context.Context, d *Deployment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deps[d.ID] = d
	return nil
}

func (m *mockDeploymentRepo) GetByID(ctx context.Context, id string) (*Deployment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.deps[id]
	if !ok {
		return nil, ErrDeploymentNotFound
	}
	return d, nil
}

func (m *mockDeploymentRepo) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*Deployment, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*Deployment
	for _, d := range m.deps {
		if d.ServiceID == serviceID {
			list = append(list, d)
		}
	}
	return list, len(list), nil
}

func (m *mockDeploymentRepo) ListByOrganization(ctx context.Context, filter DeploymentListFilter) ([]DeploymentListItem, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []DeploymentListItem
	for _, d := range m.deps {
		list = append(list, DeploymentListItem{
			ID:        d.ID,
			ServiceID: d.ServiceID,
			ProjectID: d.ProjectID,
			Status:    d.Status,
		})
	}
	return list, len(list), nil
}

func (m *mockDeploymentRepo) Update(ctx context.Context, d *Deployment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deps[d.ID] = d
	return nil
}

func (m *mockDeploymentRepo) UpdateStatus(ctx context.Context, id string, status DeploymentStatus, buildLogs, containerID, resourceIssue string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deps[id]
	if !ok {
		return ErrDeploymentNotFound
	}
	d.Status = status
	d.BuildLogs = buildLogs
	d.ContainerID = containerID
	d.ResourceIssue = resourceIssue
	return nil
}

func (m *mockDeploymentRepo) SetPinned(ctx context.Context, id string, pinned bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deps[id]
	if !ok {
		return ErrDeploymentNotFound
	}
	d.Pinned = pinned
	return nil
}

func (m *mockDeploymentRepo) TogglePinned(ctx context.Context, id string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deps[id]
	if !ok {
		return false, ErrDeploymentNotFound
	}
	d.Pinned = !d.Pinned
	return d.Pinned, nil
}

func (m *mockDeploymentRepo) ListForRetention(ctx context.Context, projectID string) ([]*Deployment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*Deployment
	for _, d := range m.deps {
		if d.ProjectID == projectID {
			list = append(list, d)
		}
	}
	return list, nil
}

func (m *mockDeploymentRepo) SetArtifactRetainedAt(ctx context.Context, id string, at *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deps[id]
	if !ok {
		return ErrDeploymentNotFound
	}
	d.ArtifactRetainedAt = at
	return nil
}

func (m *mockDeploymentRepo) SetRestartLoop(ctx context.Context, id string, finding string) error {
	return nil
}

func (m *mockDeploymentRepo) RecoverInterrupted(ctx context.Context) error {
	return nil
}

func TestDeploymentService_CRUD(t *testing.T) {
	repo := newMockDeploymentRepo()
	svc := NewDeploymentService(repo, nil, nil, nil, nil, nil)

	dep := &Deployment{
		ID:        "dep-1",
		ServiceID: "svc-1",
		ProjectID: "proj-1",
		Status:    DeploymentStatusPending,
		Branch:    "main",
	}

	created, err := svc.CreateDeployment(context.Background(), dep)
	require.NoError(t, err)
	assert.Equal(t, "dep-1", created.ID)

	got, err := svc.GetDeployment(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.Equal(t, "svc-1", got.ServiceID)

	// Toggle pin
	pinned, err := svc.TogglePin(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.True(t, pinned)

	pinnedAgain, err := svc.TogglePin(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.False(t, pinnedAgain)

	// Update status
	err = svc.UpdateStatus(context.Background(), DeployStatusOpts{
		ID:        "dep-1",
		Status:    DeploymentStatusReady,
		BuildLogs: "Done.",
	})
	require.NoError(t, err)

	gotUpdated, _ := svc.GetDeployment(context.Background(), "dep-1")
	assert.Equal(t, DeploymentStatusReady, gotUpdated.Status)
}

func TestBuildAdmissionGuard(t *testing.T) {
	guard := NewBuildAdmissionGuard()

	// Acquire service 1
	err := guard.Acquire("svc-1", "")
	require.NoError(t, err)

	// Acquire service 1 again should fail (limit 1)
	err = guard.Acquire("svc-1", "")
	assert.ErrorIs(t, err, ErrAdmissionLimitExceeded)

	// Service 2 can acquire
	err = guard.Acquire("svc-2", "")
	require.NoError(t, err)

	// Release service 1
	guard.Release("svc-1")
	err = guard.Acquire("svc-1", "")
	require.NoError(t, err)

	status := guard.Status()
	assert.Equal(t, 2, status.ActiveBuilds)
}

func TestRoutes_RegisterDeployments(t *testing.T) {
	e := echo.New()
	apiGroup := e.Group("/api")
	authGroup := e.Group("/auth")

	depRepo := newMockDeploymentRepo()
	svc := NewDeploymentService(depRepo, nil, nil, nil, nil, nil)
	regRepo := newMockRegistryRepo()
	regSvc := NewRegistryService(regRepo)
	archSvc := NewArchiveService(svc)

	dHandler := NewDeploymentHandler(svc, nil, nil, nil, nil)
	lHandler := NewLifecycleHandler(svc, nil, nil)
	aHandler := NewArchiveHandler(archSvc, nil)
	rHandler := NewRegistryHandler(regSvc)
	tHandler := NewTerminalHandler(nil, nil, nil, nil, nil)

	RegisterRoutes(apiGroup, authGroup, dHandler, lHandler, aHandler, rHandler, tHandler, RouteMiddleware{})

	routes := e.Routes()
	assert.NotEmpty(t, routes)
}
