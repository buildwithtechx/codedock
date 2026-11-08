package http

import (
	"context"
	"database/sql"
	"os"

	"github.com/labstack/echo/v4"

	"codedock/internal/handlers/system"
	"codedock/internal/models"
	"codedock/internal/modules/audit"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/observability"
	traffic "codedock/internal/platform/observability"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/analytics"
	"codedock/internal/services/attention"
	"codedock/internal/services/migrations"
	systemservices "codedock/internal/services/system"
)

type statsHealthAdapter struct {
	monitor *observability.StatsMonitor
}

func (a *statsHealthAdapter) GetHealth(ctx context.Context, containerID string) (attention.AttentionHealth, error) {
	health, err := a.monitor.GetHealth(ctx, containerID)
	if err != nil {
		return attention.AttentionHealth{Status: "unknown"}, err
	}
	return attention.AttentionHealth{
		Status:  string(health.Status),
		Running: health.Status == observability.ContainerHealthStatusRunning,
	}, nil
}

func configureAnalytics(
	server *Server,
	db *sql.DB,
	vlt *vault.Vault,
	traefikManager *networking.TraefikManager,
	deployRepo *repositories.DeploymentRepo,
	projectRepo *repositories.ProjectRepo,
	appRepo *repositories.AppServiceRepo,
	domainRepo *repositories.DomainRepo,
	userRepo *repositories.UserRepo,
	orgRepo repositories.OrganizationRepository,
	metricsService *systemservices.MetricsService,
	statsMonitor *observability.StatsMonitor,
	deploymentService attention.WorkloadDeployer,
	backupService *backupmodule.Service,
	migrationService *migrations.Service,
	auditService *audit.AuditService,
) {
	trafficRepo := repositories.NewTrafficRepository(db)
	attentionRepo := repositories.NewAttentionRepository(db)
	managedRepo := repositories.NewManagedRepository(db, vlt)
	serverRepo := repositories.NewServerRepository(db, vlt)

	usage := analytics.NewUsageReader(appRepo, metricsService, statsMonitor)
	analyticsService := analytics.NewService(trafficRepo, deployRepo, projectRepo, appRepo, domainRepo, userRepo, orgRepo, attentionRepo, usage)
	attentionService := attention.NewService(attentionRepo, deployRepo, appRepo, projectRepo, userRepo, orgRepo, backupService, migrationService, managedRepo, serverRepo, &statsHealthAdapter{monitor: statsMonitor}, deploymentService)

	server.analyticsHandler = system.NewAnalyticsHandler(analyticsService)
	server.attentionHandler = system.NewAttentionHandler(attentionService, auditService)

	if path := traefikManager.AccessLogPath(); path != "" {
		geoResolver, _ := traffic.NewGeoIPResolver(os.Getenv("GEOIP_DB_PATH"))
		collector := traffic.NewCollector(path, trafficRepo, analyticsService.AttributeDomain, geoResolver)
		go collector.Run(context.Background())
	}
}

func (s *Server) registerAnalyticsRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole(models.MemberPermissionMember)
	admin := s.authGuard.RequireOrgRole(models.MemberPermissionAdmin)
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
