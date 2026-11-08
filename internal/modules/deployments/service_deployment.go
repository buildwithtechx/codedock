package deployments

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"

	"codedock/internal/platform/observability"
)

type AppServiceLookup interface {
	GetAppService(ctx context.Context, id string) (*AppServiceSummary, error)
	UpdateAppContainerID(ctx context.Context, id string, containerID string) error
}

type ProjectLookup interface {
	HasPermission(ctx context.Context, projectID, userID, role, action string) bool
}

type CodeSyncProvider interface {
	SyncCodebase(ctx context.Context, repoURL, branch, commitHash, targetDir string, logWriter io.Writer) error
}

type DeployEngine interface {
	BeginServiceOperation(ctx context.Context, serviceID string) (context.Context, func(), error)
	DeployApp(ctx context.Context, app *AppServiceSummary, sourceDir string, logWriter io.Writer) (string, error)
	StopApp(ctx context.Context, containerID string) error
}

type DeploymentService struct {
	repo           DeploymentRepository
	appLookup      AppServiceLookup
	projectLookup  ProjectLookup
	codeSync       CodeSyncProvider
	deployEngine   DeployEngine
	statsMonitor   *observability.StatsMonitor
	admissionGuard *BuildAdmissionGuard
	operationMu    sync.Mutex
	operations     map[string]context.CancelFunc
}

func NewDeploymentService(
	repo DeploymentRepository,
	appLookup AppServiceLookup,
	projectLookup ProjectLookup,
	codeSync CodeSyncProvider,
	deployEngine DeployEngine,
	statsMonitor *observability.StatsMonitor,
) *DeploymentService {
	return &DeploymentService{
		repo:           repo,
		appLookup:      appLookup,
		projectLookup:  projectLookup,
		codeSync:       codeSync,
		deployEngine:   deployEngine,
		statsMonitor:   statsMonitor,
		admissionGuard: NewBuildAdmissionGuard(),
		operations:     make(map[string]context.CancelFunc),
	}
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, d *Deployment) (*Deployment, error) {
	if d == nil || d.ServiceID == "" {
		return nil, errors.New("valid deployment with serviceId required")
	}
	if d.ImageRef == "" && s.appLookup != nil {
		if app, err := s.appLookup.GetAppService(ctx, d.ServiceID); err == nil && app != nil {
			d.ImageRef = app.ImageRef
		}
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Status == "" {
		d.Status = DeploymentStatusPending
	}
	now := time.Now().UTC()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id string) (*Deployment, error) {
	if id == "" {
		return nil, errors.New("deployment id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeploymentService) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*Deployment, int, error) {
	if serviceID == "" {
		return nil, 0, errors.New("service id required")
	}
	return s.repo.ListByService(ctx, serviceID, limit, offset)
}

func (s *DeploymentService) ListByOrganization(ctx context.Context, filter DeploymentListFilter) ([]DeploymentListItem, int, error) {
	if filter.OrganizationID == "" {
		return nil, 0, errors.New("organization id required")
	}
	return s.repo.ListByOrganization(ctx, filter)
}

func (s *DeploymentService) UpdateDeployment(ctx context.Context, d *Deployment) error {
	if d == nil || d.ID == "" {
		return errors.New("valid deployment required for update")
	}
	d.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, d)
}

func (s *DeploymentService) UpdateStatus(ctx context.Context, opts DeployStatusOpts) error {
	if opts.ID == "" {
		return errors.New("deployment id required")
	}
	return s.repo.UpdateStatus(ctx, opts.ID, opts.Status, opts.BuildLogs, opts.ContainerID, opts.ResourceIssue)
}

func (s *DeploymentService) GetAdmissionStatus() BuildAdmissionStatus {
	if s.admissionGuard == nil {
		return BuildAdmissionStatus{MaxConcurrentBuilds: DefaultMaxConcurrentBuilds}
	}
	return s.admissionGuard.Status()
}

func (s *DeploymentService) GetMetrics(ctx context.Context, appID string) (*observability.ContainerHealth, error) {
	if s.appLookup == nil {
		return nil, ErrAppServiceRequired
	}
	app, err := s.appLookup.GetAppService(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app.ContainerID == "" {
		return &observability.ContainerHealth{Status: observability.ContainerHealthStatusNotDeployed}, nil
	}
	if s.statsMonitor == nil {
		return &observability.ContainerHealth{Status: observability.ContainerHealthStatusRunning}, nil
	}
	return s.statsMonitor.GetHealth(ctx, app.ContainerID)
}

func (s *DeploymentService) TriggerDeployment(ctx context.Context, serviceID, environmentID, projectID, branch, trigger, buildLogs string) (string, error) {
	dep := &Deployment{
		ID:            uuid.NewString(),
		ServiceID:     serviceID,
		EnvironmentID: environmentID,
		ProjectID:     projectID,
		Status:        DeploymentStatusBuilding,
		Branch:        branch,
		Trigger:       trigger,
		BuildLogs:     buildLogs,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	created, err := s.CreateDeployment(ctx, dep)
	if err != nil {
		return "", err
	}
	s.ExecuteDeploymentAsync(created)
	return created.ID, nil
}
