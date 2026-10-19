package deployments

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

type mockProjectRepo struct {
	project *models.ProjectConfig
}

func (m *mockProjectRepo) ListByOrganization(_ context.Context, _ string, _ string, _, _ int) ([]models.ProjectConfig, int, error) {
	return nil, 0, nil
}
func (m *mockProjectRepo) ListAll(_ context.Context, _ string, _, _ int) ([]models.ProjectConfig, int, error) {
	return nil, 0, nil
}
func (m *mockProjectRepo) Get(_ context.Context, _ string) (*models.ProjectConfig, error) {
	return m.project, nil
}
func (m *mockProjectRepo) GetByOrganization(_ context.Context, _, _ string) (*models.ProjectConfig, error) {
	return m.project, nil
}
func (m *mockProjectRepo) CountByUser(_ context.Context, _ string) (int, error) {
	return 0, nil
}
func (m *mockProjectRepo) SearchIDsByName(_ context.Context, _ string, _ int) ([]string, error) {
	return nil, nil
}
func (m *mockProjectRepo) Create(_ context.Context, _ *models.ProjectConfig) error {
	return nil
}
func (m *mockProjectRepo) SetServer(_ context.Context, _, _ string) error {
	return nil
}
func (m *mockProjectRepo) Delete(_ context.Context, _ string) error {
	return nil
}
func (m *mockProjectRepo) UpdateRollbackSettings(_ context.Context, _ string, strategy *string, window *int) error {
	if strategy != nil {
		m.project.DefaultRollbackStrategy = *strategy
	}
	m.project.RollbackWindow = window
	return nil
}
func (m *mockProjectRepo) UpdateCapacityMeasurement(_ context.Context, _ string, snapshotSizeBytes int64) error {
	m.project.SnapshotSizeBytes = &snapshotSizeBytes
	return nil
}

type mockDeploymentRepo struct {
	deployments []*models.Deployment
	purged      []string
}

func (m *mockDeploymentRepo) Create(_ context.Context, _ *models.Deployment) error {
	return nil
}
func (m *mockDeploymentRepo) GetByID(_ context.Context, id string) (*models.Deployment, error) {
	for _, d := range m.deployments {
		if d.ID == id {
			return d, nil
		}
	}
	return nil, nil
}
func (m *mockDeploymentRepo) ListByService(_ context.Context, _ string, _, _ int) ([]*models.Deployment, int, error) {
	return m.deployments, len(m.deployments), nil
}
func (m *mockDeploymentRepo) ListByOrganization(_ context.Context, _ models.DeploymentListFilter) ([]models.DeploymentListItem, int, error) {
	return nil, 0, nil
}
func (m *mockDeploymentRepo) Update(_ context.Context, _ *models.Deployment) error {
	return nil
}
func (m *mockDeploymentRepo) UpdateStatus(_ context.Context, _ string, _ models.DeploymentStatus, _, _ string) error {
	return nil
}
func (m *mockDeploymentRepo) ListForRetention(_ context.Context, _ string) ([]*models.Deployment, error) {
	return m.deployments, nil
}
func (m *mockDeploymentRepo) SetArtifactRetainedAt(_ context.Context, id string, at *time.Time) error {
	for _, d := range m.deployments {
		if d.ID == id {
			d.ArtifactRetainedAt = at
			if at == nil {
				m.purged = append(m.purged, id)
			}
			return nil
		}
	}
	return nil
}
func (m *mockDeploymentRepo) SetPinned(_ context.Context, id string, pinned bool) error {
	for _, d := range m.deployments {
		if d.ID == id {
			d.Pinned = pinned
			return nil
		}
	}
	return nil
}
func (m *mockDeploymentRepo) TogglePinned(_ context.Context, id string) (bool, error) {
	for _, d := range m.deployments {
		if d.ID == id {
			d.Pinned = !d.Pinned
			return d.Pinned, nil
		}
	}
	return false, nil
}

func intPtr(i int) *int { return &i }

func makeDep(id string, status models.DeploymentStatus, imageRef, commit string, pinned bool, retained bool) *models.Deployment {
	var at *time.Time
	if retained {
		t := time.Now().UTC()
		at = &t
	}
	return &models.Deployment{
		ID:                 id,
		ProjectID:          "p1",
		Status:             status,
		ImageRef:           imageRef,
		CommitHash:         commit,
		Pinned:             pinned,
		ArtifactRetainedAt: at,
	}
}

func TestResolveRollbackWindow(t *testing.T) {
	svc := NewRetentionService(&mockProjectRepo{}, &mockDeploymentRepo{}, nil)

	pAuto := &models.ProjectConfig{RollbackWindow: nil}
	w, src := svc.ResolveRollbackWindow(pAuto)
	if w != 5 || src != "instance-default" {
		t.Fatalf("expected 5/instance-default, got %d/%s", w, src)
	}

	pExplicit := &models.ProjectConfig{RollbackWindow: intPtr(3)}
	w, src = svc.ResolveRollbackWindow(pExplicit)
	if w != 3 || src != "explicit" {
		t.Fatalf("expected 3/explicit, got %d/%s", w, src)
	}

	pNegative := &models.ProjectConfig{RollbackWindow: intPtr(-2)}
	w, src = svc.ResolveRollbackWindow(pNegative)
	if w != 0 || src != "explicit" {
		t.Fatalf("expected 0/explicit, got %d/%s", w, src)
	}

	pExcess := &models.ProjectConfig{RollbackWindow: intPtr(50)}
	w, src = svc.ResolveRollbackWindow(pExcess)
	if w != 20 || src != "explicit" {
		t.Fatalf("expected 20/explicit, got %d/%s", w, src)
	}
}

func TestRollbackStoragePressure(t *testing.T) {
	if rollbackStoragePressure(100, 30) {
		t.Errorf("30%% free space should not trigger pressure")
	}
	if !rollbackStoragePressure(100, 15) {
		t.Errorf("15%% free space should trigger pressure (<20%%)")
	}
	if rollbackStoragePressure(0, 0) || rollbackStoragePressure(-10, -5) {
		t.Errorf("invalid disk stats should not trigger pressure")
	}
}

func TestComputeRetentionInventory(t *testing.T) {
	p := &models.ProjectConfig{ID: "p1", RollbackWindow: intPtr(2)}
	deps := []*models.Deployment{
		makeDep("d5", models.DeploymentStatusReady, "img5", "c5", false, true), // active / latest ready
		makeDep("d4", models.DeploymentStatusReady, "img4", "c4", false, true), // window slot 1
		makeDep("d3", models.DeploymentStatusReady, "img3", "c3", false, true), // window slot 2
		makeDep("d2", models.DeploymentStatusReady, "img2", "c2", true, true),  // pinned (exempt from budget)
		makeDep("d1", models.DeploymentStatusReady, "img1", "c1", false, true), // overflow -> purge candidate
		makeDep("d0", models.DeploymentStatusReady, "img0", "c0", false, true), // overflow -> purge candidate
	}

	svc := NewRetentionService(&mockProjectRepo{project: p}, &mockDeploymentRepo{deployments: deps}, nil)
	retained, overflow, images, err := svc.ComputeRetentionInventory(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(retained) != 4 {
		t.Errorf("expected 4 retained (d5, d4, d3, d2 pinned), got %d", len(retained))
	}
	if len(overflow) != 2 {
		t.Errorf("expected 2 overflow (d1, d0), got %d", len(overflow))
	}
	if _, ok := images["img2"]; !ok {
		t.Errorf("expected pinned image img2 in retained images")
	}
}

func TestPruneProjectRetainsSharedImage(t *testing.T) {
	p := &models.ProjectConfig{ID: "p1", RollbackWindow: intPtr(1)}
	deps := []*models.Deployment{
		makeDep("d3", models.DeploymentStatusReady, "img-shared", "c3", false, true), // active
		makeDep("d2", models.DeploymentStatusReady, "img2", "c2", false, true),       // window slot 1
		makeDep("d1", models.DeploymentStatusReady, "img-shared", "c1", false, true), // overflow, but shares image with d3!
	}

	dRepo := &mockDeploymentRepo{deployments: deps}
	svc := NewRetentionService(&mockProjectRepo{project: p}, dRepo, nil)
	svc.SetDiskReader(func(_ context.Context, _ *models.ProjectConfig) (int64, int64, error) {
		return 100, 50, nil // healthy disk
	})

	summary, err := svc.PruneProject(context.Background(), "p1")
	if err != nil {
		t.Fatalf("prune error: %v", err)
	}
	if summary.Purged != 1 {
		t.Errorf("expected 1 purged overflow record, got %d", summary.Purged)
	}
	if len(dRepo.purged) != 1 || dRepo.purged[0] != "d1" {
		t.Errorf("expected d1 artifact marker cleared, got: %v", dRepo.purged)
	}
}

func TestReclaimStoragePressureReclaimsOldestReproducible(t *testing.T) {
	p := &models.ProjectConfig{ID: "p1", RollbackWindow: nil} // automatic window
	deps := []*models.Deployment{
		makeDep("d4", models.DeploymentStatusReady, "img4", "c4", false, true), // active
		makeDep("d3", models.DeploymentStatusReady, "img3", "c3", false, true), // reproducible
		makeDep("d2", models.DeploymentStatusReady, "img2", "", false, true),   // non-reproducible (no commit, not sha)
		makeDep("d1", models.DeploymentStatusReady, "img1", "c1", false, true), // reproducible, oldest
	}

	diskReads := 0
	svc := NewRetentionService(&mockProjectRepo{project: p}, &mockDeploymentRepo{deployments: deps}, nil)
	svc.SetDiskReader(func(_ context.Context, _ *models.ProjectConfig) (int64, int64, error) {
		diskReads++
		if diskReads == 1 {
			return 100, 10, nil // 10% free -> storage pressure!
		}
		return 100, 30, nil // 30% free -> relieved after 1 reclaim!
	})

	summary, err := svc.ReclaimStoragePressure(context.Background(), "p1")
	if err != nil {
		t.Fatalf("reclaim error: %v", err)
	}
	if summary.Purged != 1 {
		t.Errorf("expected 1 release reclaimed, got %d", summary.Purged)
	}
	if deps[3].ArtifactRetainedAt != nil {
		t.Errorf("expected d1 (oldest reproducible) to have artifactRetainedAt cleared")
	}
	if deps[2].ArtifactRetainedAt == nil {
		t.Errorf("expected d2 (non-reproducible) to NOT be reclaimed")
	}
}
