package deployments

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"codedock/internal/engine/deploy"
	"codedock/internal/engine/observability"
	"codedock/internal/engine/ssh"
	"codedock/internal/models"
	"codedock/internal/repositories"
)

type DeploymentService struct {
	LocalDocker             *client.Client
	Databases               repositories.DatabaseRepository
	RuntimeKinds            dockerRuntimeStore
	Servers                 repositories.ServerRepository
	HostGate                deploy.VolumeOperations
	BeforeProjectDeployment func(context.Context, string, string) error
	Runtime                 TargetRuntime
	BeforeDeployment        func(context.Context, string) error
	repo                    repositories.DeploymentRepository
	appRepo                 repositories.AppServiceRepository
	projectRepo             repositories.ProjectRepository
	deployer                *deploy.Deployer
	gitService              *GitService
	statsMonitor            *observability.StatsMonitor
	volumeRepo              repositories.ServiceVolumeRepository
	sshManager              *ssh.SSHManager
	admissionGuard          *BuildAdmissionGuard
	operationMu             sync.Mutex
	operations              map[string]context.CancelFunc
}

func NewDeploymentService(
	r repositories.DeploymentRepository,
	ar repositories.AppServiceRepository,
	pr repositories.ProjectRepository,
	d *deploy.Deployer,
	gs *GitService,
	sm *observability.StatsMonitor,
	vr repositories.ServiceVolumeRepository,
	smgr *ssh.SSHManager,
) *DeploymentService {
	return &DeploymentService{
		repo:           r,
		appRepo:        ar,
		projectRepo:    pr,
		deployer:       d,
		gitService:     gs,
		statsMonitor:   sm,
		volumeRepo:     vr,
		sshManager:     smgr,
		admissionGuard: NewBuildAdmissionGuard(),
		operations:     make(map[string]context.CancelFunc),
	}
}

func (s *DeploymentService) GetService(ctx context.Context, id string) (*models.AppService, error) {
	if s.appRepo == nil {
		return nil, errors.New("app repo not available")
	}
	return s.appRepo.GetByID(ctx, id)
}

func (s *DeploymentService) UpdateService(ctx context.Context, app *models.AppService) error {
	if s.appRepo == nil {
		return errors.New("app repo not available")
	}
	return s.appRepo.Update(ctx, app)
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error) {
	if d == nil || d.ServiceID == "" {
		return nil, errors.New("valid deployment with serviceId required")
	}
	if d.ImageRef == "" && s.appRepo != nil {
		if app, err := s.appRepo.GetByID(ctx, d.ServiceID); err == nil && app != nil {
			d.ImageRef = app.ImageRef
		}
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Status == "" {
		d.Status = models.DeploymentStatusPending
	}
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	if id == "" {
		return nil, errors.New("deployment id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeploymentService) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error) {
	if serviceID == "" {
		return nil, 0, errors.New("service id required")
	}
	return s.repo.ListByService(ctx, serviceID, limit, offset)
}

func (s *DeploymentService) ListByOrganization(ctx context.Context, filter models.DeploymentListFilter) ([]models.DeploymentListItem, int, error) {
	if filter.OrganizationID == "" {
		return nil, 0, errors.New("organization id required")
	}
	return s.repo.ListByOrganization(ctx, filter)
}

func (s *DeploymentService) UpdateDeployment(ctx context.Context, d *models.Deployment) error {
	if d == nil || d.ID == "" {
		return errors.New("valid deployment required for update")
	}
	d.UpdatedAt = time.Now()
	return s.repo.Update(ctx, d)
}

type DeployStatusOpts struct {
	ID            string
	Status        models.DeploymentStatus
	BuildLogs     string
	ContainerID   string
	ResourceIssue string
}

func (s *DeploymentService) UpdateStatus(ctx context.Context, opts DeployStatusOpts) error {
	if opts.ID == "" {
		return errors.New("deployment id required")
	}
	return s.repo.UpdateStatus(ctx, opts.ID, opts.Status, opts.BuildLogs, opts.ContainerID, opts.ResourceIssue)
}

func (s *DeploymentService) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	if s.deployer == nil || s.appRepo == nil {
		return "", errors.New("deployer or app repo not available")
	}
	operation, release, err := s.deployer.BeginServiceOperation(ctx, appID)
	if err != nil {
		return "", err
	}
	defer release()
	ctx = operation
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return "", err
	}
	if s.BeforeProjectDeployment != nil {
		if err := s.BeforeProjectDeployment(ctx, app.ProjectID, app.ID); err != nil {
			return "", fmt.Errorf("pre-deployment backup: %w", err)
		}
	} else if s.BeforeDeployment != nil {
		if err := s.BeforeDeployment(ctx, app.ID); err != nil {
			return "", fmt.Errorf("pre-deployment backup: %w", err)
		}
	}
	containerID, err := s.deployTarget(ctx, app, sourceDir, logWriter)
	if err == nil && containerID != "" {
		app.ContainerID = containerID
		if updateErr := s.appRepo.Update(ctx, app); updateErr != nil {
			return "", fmt.Errorf("save deployed application: %w", updateErr)
		}
	}
	return containerID, err
}

func (s *DeploymentService) GetMetrics(ctx context.Context, appID string) (*observability.ContainerHealth, error) {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, appID)
		if err != nil {
			return nil, err
		}
		if handles {
			observed, err := s.Runtime.Observe(ctx, appID)
			if err != nil {
				return nil, err
			}
			if observed.Desired > 0 && !observed.MetricsAvailable {
				return nil, fmt.Errorf("cluster metrics unavailable: %s", observed.MetricsError)
			}
			health := &observability.ContainerHealth{Status: observability.ContainerHealthStatusStopped}
			if observed.Available > 0 {
				health.Status = observability.ContainerHealthStatusRunning
			}
			for _, pod := range observed.Pods {
				health.CPUUsagePercentage += pod.CPU * 100
				health.MemoryUsageBytes += pod.MemoryBytes
			}
			return health, nil
		}
	}
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app.ContainerID == "" {
		return &observability.ContainerHealth{Status: observability.ContainerHealthStatusNotDeployed}, nil
	}
	if s.statsMonitor == nil {
		return nil, errors.New("stats monitor not available")
	}
	target, release, err := s.DockerForService(ctx, appID)
	if err != nil {
		return nil, err
	}
	defer release()
	return observability.NewStatsMonitor(target).GetHealth(ctx, app.ContainerID)
}
