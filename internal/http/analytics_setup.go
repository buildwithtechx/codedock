package http

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"codedock/internal/engine/networking"
	"codedock/internal/engine/observability"
	"codedock/internal/engine/traffic"
	"codedock/internal/handlers/system"
	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/services/analytics"
	"codedock/internal/services/attention"
	backupservices "codedock/internal/services/backups"
	deploymentservices "codedock/internal/services/deployments"
	"codedock/internal/services/migrations"
	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"
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
	vault *utils.Vault,
	traefikManager *networking.TraefikManager,
	deployRepo *repositories.DeploymentRepo,
	projectRepo *repositories.ProjectRepo,
	appRepo *repositories.AppServiceRepo,
	domainRepo *repositories.DomainRepo,
	userRepo *repositories.UserRepo,
	orgRepo repositories.OrganizationRepository,
	metricsService *systemservices.MetricsService,
	statsMonitor *observability.StatsMonitor,
	deploymentService *deploymentservices.DeploymentService,
	backupService *backupservices.BackupService,
	migrationService *migrations.Service,
) {
	trafficRepo := repositories.NewTrafficRepository(db)
	attentionRepo := repositories.NewAttentionRepository(db)
	managedRepo := repositories.NewManagedRepository(db, vault)
	serverRepo := repositories.NewServerRepository(db, vault)

	usage := analytics.NewUsageReader(appRepo, metricsService, statsMonitor)
	analyticsService := analytics.NewService(trafficRepo, deployRepo, projectRepo, appRepo, domainRepo, userRepo, orgRepo, attentionRepo, usage)
	attentionService := attention.NewService(attentionRepo, deployRepo, appRepo, projectRepo, userRepo, orgRepo, backupService, migrationService, managedRepo, serverRepo, &statsHealthAdapter{monitor: statsMonitor}, deploymentService)

	server.analyticsHandler = system.NewAnalyticsHandler(analyticsService)
	server.attentionHandler = system.NewAttentionHandler(attentionService)

	if path := traefikManager.AccessLogPath(); path != "" {
		collector := traffic.NewCollector(path, trafficRepo, analyticsService.AttributeDomain)
		go collector.Run(context.Background())
	}
}

func (s *Server) registerAnalyticsRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole(models.MemberPermissionMember)
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
	authGroup.POST("/organizations/:id/attention/:issueId/ack", s.attentionHandler.Acknowledge, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/resolve", s.attentionHandler.Resolve, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/act", s.attentionHandler.Act, member, write)
}
