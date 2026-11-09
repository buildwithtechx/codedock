package http

import (
	"context"
	"database/sql"
	"os"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/audit"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/servers"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/observability"
	traffic "codedock/internal/platform/observability"
	"codedock/internal/platform/vault"
)

type statsHealthAdapter struct {
	monitor *observability.StatsMonitor
}

func (a *statsHealthAdapter) GetHealth(ctx context.Context, containerID string) (monitoring.AttentionHealth, error) {
	health, err := a.monitor.GetHealth(ctx, containerID)
	if err != nil {
		return monitoring.AttentionHealth{Status: "unknown"}, err
	}
	return monitoring.AttentionHealth{
		Status:  string(health.Status),
		Running: health.Status == observability.ContainerHealthStatusRunning,
	}, nil
}

type attentionManagedAdapter struct {
	repo servers.ManagedRepository
}

func (a *attentionManagedAdapter) GetQuota(ctx context.Context, orgID string) (*monitoring.ManagedQuota, error) {
	q, err := a.repo.GetQuota(ctx, orgID)
	if err != nil || q == nil {
		return nil, err
	}
	return &monitoring.ManagedQuota{MaxServers: q.MaxServers, MaxMemoryGB: q.MaxMemoryGB}, nil
}

func (a *attentionManagedAdapter) ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error) {
	return a.repo.ListOrgServerIDs(ctx, orgID)
}

type attentionServerAdapter struct {
	repo servers.Repository
}

func (a *attentionServerAdapter) GetByID(ctx context.Context, id string) (*monitoring.Server, error) {
	s, err := a.repo.GetByID(ctx, id)
	if err != nil || s == nil {
		return nil, err
	}
	return &monitoring.Server{ID: s.ID, ServerType: s.ServerType}, nil
}

type trafficCollectorAdapter struct {
	repo monitoring.TrafficRepository
}

func (t *trafficCollectorAdapter) RecordBatch(ctx context.Context, samples []observability.TrafficSample) error {
	mSamples := make([]monitoring.TrafficSample, len(samples))
	for i, s := range samples {
		mSamples[i] = monitoring.TrafficSample{
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

func configureAnalytics(
	server *Server,
	db *sql.DB,
	vlt *vault.Vault,
	traefikManager *networking.TraefikManager,
	deployRepo monitoring.DeploymentRepository,
	projectRepo monitoring.ProjectRepository,
	appRepo monitoring.AppServiceRepository,
	domainRepo monitoring.DomainLister,
	userRepo monitoring.UserRepo,
	orgRepo monitoring.OrganizationRepository,
	metricsService *monitoring.MetricsService,
	statsMonitor *observability.StatsMonitor,
	deploymentService monitoring.WorkloadDeployer,
	backupService *backupmodule.Service,
	migrationService monitoring.MigrationRunner,
	auditService *audit.AuditService,
) {
	trafficRepo := monitoring.NewTrafficRepository(db)
	attentionRepo := monitoring.NewAttentionRepository(db)
	managedRepo := &attentionManagedAdapter{repo: servers.NewManagedRepository(db, vlt)}
	serverRepo := &attentionServerAdapter{repo: servers.NewRepository(db, vlt)}

	usage := monitoring.NewUsageReader(appRepo, metricsService, statsMonitor)
	analyticsService := monitoring.NewService(trafficRepo, deployRepo, projectRepo, appRepo, domainRepo, userRepo, orgRepo, attentionRepo, usage)
	attentionService := monitoring.NewAttentionManager(attentionRepo, deployRepo, appRepo, projectRepo, userRepo, orgRepo, backupService, migrationService, managedRepo, serverRepo, &statsHealthAdapter{monitor: statsMonitor}, deploymentService)

	server.analyticsHandler = monitoring.NewAnalyticsHandler(&analyticsServiceAdapter{svc: analyticsService})
	server.attentionHandler = monitoring.NewAttentionHandler(&attentionServiceAdapter{svc: attentionService})

	if path := traefikManager.AccessLogPath(); path != "" {
		geoResolver, _ := traffic.NewGeoIPResolver(os.Getenv("GEOIP_DB_PATH"))
		collector := traffic.NewCollector(path, &trafficCollectorAdapter{repo: trafficRepo}, analyticsService.AttributeDomain, geoResolver)
		go collector.Run(context.Background())
	}
}

func (s *Server) registerAnalyticsRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole("member")
	admin := s.authGuard.RequireOrgRole("admin")
	read := s.authGuard.RequireScope("server:read")
	write := s.authGuard.RequireScope("server:write")

	authGroup.GET("/organizations/:id/analytics/periods", s.analyticsHandler.Periods, member, read)
	authGroup.GET("/organizations/:id/analytics/dashboard", s.analyticsHandler.Dashboard, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/summary", s.analyticsHandler.Summary, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/overview", s.analyticsHandler.Overview, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/geo", s.analyticsHandler.Geo, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/live", s.analyticsHandler.Live, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/paths", s.analyticsHandler.GetPaths, member, read)
	authGroup.PUT("/organizations/:id/projects/:projectId/analytics/paths", s.analyticsHandler.SetPaths, member, write)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/deployments", s.analyticsHandler.DeploymentStats, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage", s.analyticsHandler.Usage, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/resources", s.analyticsHandler.Resources, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/history", s.analyticsHandler.UsageHistory, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/stream", s.analyticsHandler.UsageStream, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/container", s.analyticsHandler.Container, member, read)

	authGroup.GET("/organizations/:id/attention", s.attentionHandler.List, member, read)
	authGroup.POST("/organizations/:id/attention/evaluate", s.attentionHandler.Evaluate, member, write)
	authGroup.DELETE("/organizations/:id/attention/history", s.attentionHandler.ClearHistory, admin, write)
	authGroup.POST("/organizations/:id/attention/:issueId/ack", s.attentionHandler.Acknowledge, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/resolve", s.attentionHandler.Resolve, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/act", s.attentionHandler.Act, member, write)
}
