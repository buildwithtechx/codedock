package monitoring

import (
	"context"
	"database/sql"
	"os"

	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/servers/managed"
	"codedock/internal/platform/observability"
	traffic "codedock/internal/platform/observability"
	"codedock/internal/platform/vault"
)

type statsHealthAdapter struct {
	monitor *observability.StatsMonitor
}

func (a *statsHealthAdapter) GetHealth(ctx context.Context, containerID string) (AttentionHealth, error) {
	health, err := a.monitor.GetHealth(ctx, containerID)
	if err != nil {
		return AttentionHealth{Status: "unknown"}, err
	}
	return AttentionHealth{
		Status:  string(health.Status),
		Running: health.Status == observability.ContainerHealthStatusRunning,
	}, nil
}

type attentionManagedAdapter struct {
	repo managed.ManagedRepository
}

func (a *attentionManagedAdapter) GetQuota(ctx context.Context, orgID string) (*ManagedQuota, error) {
	q, err := a.repo.GetQuota(ctx, orgID)
	if err != nil || q == nil {
		return nil, err
	}
	return &ManagedQuota{MaxServers: q.MaxServers, MaxMemoryGB: q.MaxMemoryGB}, nil
}

func (a *attentionManagedAdapter) ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error) {
	return a.repo.ListOrgServerIDs(ctx, orgID)
}

type attentionServerAdapter struct {
	repo servers.Repository
}

func (a *attentionServerAdapter) GetByID(ctx context.Context, id string) (*Server, error) {
	s, err := a.repo.GetByID(ctx, id)
	if err != nil || s == nil {
		return nil, err
	}
	return &Server{ID: s.ID, ServerType: s.ServerType}, nil
}

type trafficCollectorAdapter struct {
	repo TrafficRepository
}

func (t *trafficCollectorAdapter) RecordBatch(ctx context.Context, samples []observability.TrafficSample) error {
	mSamples := make([]TrafficSample, len(samples))
	for i, s := range samples {
		mSamples[i] = TrafficSample{
			Time:       s.Time,
			ProjectID:  s.ProjectID,
			Domain:     s.Domain,
			Path:       s.Path,
			Status:     s.Status,
			Bytes:      s.Bytes,
			DurationMs: s.DurationMs,
			ClientIP:   s.ClientIP,
			Country:    s.Country,
		}
	}
	return t.repo.RecordBatch(ctx, mSamples)
}

func (t *trafficCollectorAdapter) PathsEnabled(ctx context.Context, projectID string) (bool, error) {
	return t.repo.PathsEnabled(ctx, projectID)
}

func (t *trafficCollectorAdapter) DeleteBefore(ctx context.Context, cutoff string) error {
	return t.repo.DeleteBefore(ctx, cutoff)
}

func (t *trafficCollectorAdapter) RetentionCutoff(days int) string {
	return t.repo.RetentionCutoff(days)
}

type Services struct {
	AnalyticsService *Service
	AttentionService *AttentionManager
}

func InitAnalyticsServices(
	db *sql.DB,
	vlt *vault.Vault,
	accessLogPath string,
	deployRepo DeploymentRepository,
	projectRepo ProjectRepository,
	appRepo AppServiceRepository,
	domainRepo DomainLister,
	userRepo UserRepo,
	orgRepo OrganizationRepository,
	metricsService *MetricsService,
	statsMonitor *observability.StatsMonitor,
	deploymentService WorkloadDeployer,
	backupService *backupmodule.Service,
	migrationService MigrationRunner,
) *Services {
	trafficRepo := NewTrafficRepository(db)
	attentionRepo := NewAttentionRepository(db)
	managedRepo := &attentionManagedAdapter{repo: managed.NewManagedRepository(db, vlt)}
	serverRepo := &attentionServerAdapter{repo: servers.NewRepository(db, vlt)}

	usage := NewUsageReader(appRepo, metricsService, statsMonitor)
	analyticsService := NewService(trafficRepo, deployRepo, projectRepo, appRepo, domainRepo, userRepo, orgRepo, attentionRepo, usage)
	attentionService := NewAttentionManager(attentionRepo, deployRepo, appRepo, projectRepo, userRepo, orgRepo, backupService, migrationService, managedRepo, serverRepo, &statsHealthAdapter{monitor: statsMonitor}, deploymentService)

	if accessLogPath != "" {
		geoResolver, _ := traffic.NewGeoIPResolver(os.Getenv("GEOIP_DB_PATH"))
		collector := traffic.NewCollector(accessLogPath, &trafficCollectorAdapter{repo: trafficRepo}, analyticsService.AttributeDomain, geoResolver)
		go collector.Run(context.Background())
	}

	return &Services{
		AnalyticsService: analyticsService,
		AttentionService: attentionService,
	}
}
