package deployments

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		repo:         r,
		appRepo:      ar,
		projectRepo:  pr,
		deployer:     d,
		gitService:   gs,
		statsMonitor: sm,
		volumeRepo:   vr,
		sshManager:   smgr,
		operations:   make(map[string]context.CancelFunc),
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
	ID          string
	Status      models.DeploymentStatus
	BuildLogs   string
	ContainerID string
}

func (s *DeploymentService) UpdateStatus(ctx context.Context, opts DeployStatusOpts) error {
	if opts.ID == "" {
		return errors.New("deployment id required")
	}
	return s.repo.UpdateStatus(ctx, opts.ID, opts.Status, opts.BuildLogs, opts.ContainerID)
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

func detectAppIcon(sourceDir string) string {
	if b, err := os.ReadFile(filepath.Join(sourceDir, "package.json")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "\"next\"") {
			return "nextjs"
		}
		if strings.Contains(s, "\"nuxt\"") {
			return "nuxt"
		}
		if strings.Contains(s, "\"@sveltejs/kit\"") {
			return "sveltekit"
		}
		if strings.Contains(s, "\"@solidjs/start\"") {
			return "solidstart"
		}
		if strings.Contains(s, "\"@remix-run") {
			return "remix"
		}
		if strings.Contains(s, "\"astro\"") {
			return "astro"
		}
		if strings.Contains(s, "\"@tanstack/start\"") {
			return "tanstack"
		}
		if strings.Contains(s, "\"expo\"") {
			return "expo"
		}
		if strings.Contains(s, "\"react-native\"") {
			return "react-native"
		}
		if strings.Contains(s, "\"vite\"") {
			return "vite"
		}
		if strings.Contains(s, "\"vue\"") {
			return "vue"
		}
		if strings.Contains(s, "\"svelte\"") {
			return "svelte"
		}
		if strings.Contains(s, "\"@angular/core\"") {
			return "angular"
		}
		if strings.Contains(s, "\"@builder.io/qwik\"") {
			return "qwik"
		}
		if strings.Contains(s, "\"gatsby\"") {
			return "gatsby"
		}
		if strings.Contains(s, "\"@redwoodjs/") {
			return "redwoodjs"
		}
		if strings.Contains(s, "\"electron\"") {
			return "electron"
		}
		if strings.Contains(s, "\"@tauri-apps/") {
			return "tauri"
		}
		if strings.Contains(s, "\"hono\"") {
			return "hono"
		}
		if strings.Contains(s, "\"elysia\"") {
			return "elysia"
		}
		if strings.Contains(s, "\"@nestjs/core\"") {
			return "nestjs"
		}
		if strings.Contains(s, "\"fastify\"") {
			return "fastify"
		}
		if strings.Contains(s, "\"express\"") {
			return "express"
		}
		if strings.Contains(s, "\"koa\"") {
			return "koa"
		}
		if strings.Contains(s, "\"@adonisjs/") {
			return "adonisjs"
		}
		if strings.Contains(s, "\"@strapi/") {
			return "strapi"
		}
		if strings.Contains(s, "\"payload\"") {
			return "payload"
		}
		if strings.Contains(s, "\"@trpc/") {
			return "trpc"
		}
		if strings.Contains(s, "\"graphql\"") {
			return "graphql"
		}
		if strings.Contains(s, "\"react\"") {
			return "react"
		}
	}

	files := map[string]string{
		"next.config.js": "nextjs", "next.config.mjs": "nextjs",
		"vite.config.ts": "vite", "vite.config.js": "vite",
		"astro.config.mjs": "astro", "nuxt.config.ts": "nuxt",
		"nest-cli.json": "nestjs", "svelte.config.js": "sveltekit",
		"remix.config.js": "remix", "angular.json": "angular",
		"pom.xml": "java", "build.gradle": "java",
		"manage.py": "django", "Gemfile": "ruby",
		"composer.json": "php", "Cargo.toml": "rust",
		"go.mod": "golang", "artisan": "laravel",
	}

	for file, icon := range files {
		if _, err := os.Stat(filepath.Join(sourceDir, file)); err == nil {
			return icon
		}
	}

	if b, err := os.ReadFile(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "fastapi") {
			return "fastapi"
		}
		if strings.Contains(s, "flask") {
			return "flask"
		}
		if strings.Contains(s, "django") {
			return "django"
		}
		return "python"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "package.json")); err == nil {
		return "nodejs"
	}

	return "git"
}
