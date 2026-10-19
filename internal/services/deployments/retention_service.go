package deployments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/client"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

const (
	DefaultRollbackWindow      = 5
	MaxRollbackWindow          = 20
	BuildStorageGCFreeFraction = 0.20
)

func rollbackStoragePressure(totalBytes, freeBytes int64) bool {
	return totalBytes > 0 && freeBytes >= 0 && freeBytes <= totalBytes &&
		float64(freeBytes) < float64(totalBytes)*BuildStorageGCFreeFraction
}

func normalizeRollbackWindow(w *int) int {
	if w == nil {
		return DefaultRollbackWindow
	}
	if *w < 0 {
		return 0
	}
	if *w > MaxRollbackWindow {
		return MaxRollbackWindow
	}
	return *w
}

type RetentionService struct {
	projectRepo    repositories.ProjectRepository
	deploymentRepo repositories.DeploymentRepository
	dockerProvider func(ctx context.Context, projectID string) (*client.Client, func(), error)
	diskReader     func(ctx context.Context, project *models.ProjectConfig) (int64, int64, error)
}

func NewRetentionService(
	projectRepo repositories.ProjectRepository,
	deploymentRepo repositories.DeploymentRepository,
	dockerProvider func(ctx context.Context, projectID string) (*client.Client, func(), error),
) *RetentionService {
	return &RetentionService{
		projectRepo:    projectRepo,
		deploymentRepo: deploymentRepo,
		dockerProvider: dockerProvider,
		diskReader: func(_ context.Context, _ *models.ProjectConfig) (int64, int64, error) {
			return readHostDiskSpace()
		},
	}
}

func (s *RetentionService) SetDiskReader(reader func(ctx context.Context, project *models.ProjectConfig) (int64, int64, error)) {
	s.diskReader = reader
}

func (s *RetentionService) ResolveRollbackWindow(project *models.ProjectConfig) (int, string) {
	if project.RollbackWindow != nil {
		return normalizeRollbackWindow(project.RollbackWindow), "explicit"
	}
	return DefaultRollbackWindow, "instance-default"
}

func (s *RetentionService) GetRollbackCapacity(ctx context.Context, projectID string) (*models.RollbackCapacity, error) {
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}
	window, source := s.ResolveRollbackWindow(project)
	var explicit *int
	if project.RollbackWindow != nil {
		w := normalizeRollbackWindow(project.RollbackWindow)
		explicit = &w
	}
	var measuredStr *string
	if project.CapacityMeasuredAt != nil {
		str := project.CapacityMeasuredAt.Format(time.RFC3339)
		measuredStr = &str
	}

	strategy := project.DefaultRollbackStrategy
	if strategy == "" {
		strategy = "snapshot"
	}

	total, free, diskErr := s.diskReader(ctx, project)
	var totalPtr, freePtr *int64
	if diskErr == nil && total > 0 {
		totalPtr = &total
		freePtr = &free
	}

	return &models.RollbackCapacity{
		Window:             window,
		Source:             source,
		Explicit:           explicit,
		SnapshotSizeBytes:  project.SnapshotSizeBytes,
		MeasuredAt:         measuredStr,
		DiskFreeBytes:      freePtr,
		DiskTotalBytes:     totalPtr,
		MaxWindow:          MaxRollbackWindow,
		DiskBudgetFraction: BuildStorageGCFreeFraction,
		Strategy:           strategy,
	}, nil
}

func (s *RetentionService) UpdateRollbackSettings(ctx context.Context, projectID string, req models.UpdateRollbackSettingsRequest) error {
	var strategy *string
	if req.DefaultRollbackStrategy != nil {
		str := strings.ToLower(strings.TrimSpace(*req.DefaultRollbackStrategy))
		if str != "snapshot" && str != "git" {
			return fmt.Errorf("invalid rollback strategy (must be 'snapshot' or 'git')")
		}
		strategy = &str
	}
	var window *int
	if req.RollbackWindow != nil {
		w := normalizeRollbackWindow(req.RollbackWindow)
		window = &w
	}
	if err := s.projectRepo.UpdateRollbackSettings(ctx, projectID, strategy, window); err != nil {
		return err
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		_, _ = s.PruneProject(bgCtx, projectID)
	}()
	return nil
}

func (s *RetentionService) IsArtifactAvailable(ctx context.Context, projectID string, dep *models.Deployment) (bool, error) {
	if dep.ArtifactRetainedAt == nil || strings.TrimSpace(dep.ImageRef) == "" {
		return false, nil
	}
	if s.dockerProvider == nil {
		return true, nil
	}
	dockerCli, release, err := s.dockerProvider(ctx, projectID)
	if err != nil {
		return true, nil
	}
	defer release()

	_, _, err = dockerCli.ImageInspectWithRaw(ctx, dep.ImageRef)
	if err != nil {
		if client.IsErrNotFound(err) {
			_ = s.deploymentRepo.SetArtifactRetainedAt(ctx, dep.ID, nil)
			return false, nil
		}
		return false, err
	}
	return true, nil
}
