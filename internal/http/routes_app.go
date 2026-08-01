package http

import (
	"os"
	"path/filepath"

	"codedock.run/codedock/apps/dashboard"
	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerAppRoutes(apiGroup, authGroup *echo.Group) {
	serviceAuthAdmin := s.RequireServiceRole(models.MemberPermissionAdmin)
	serviceAuthOwner := s.RequireServiceRole(models.MemberPermissionOwner)
	serviceAuth := s.RequireServiceRole("")

	authGroup.GET("/environments/:id/apps", s.appServiceHandler.ListByEnvironment)
	authGroup.POST("/environments/:id/apps", s.appServiceHandler.Create)
	authGroup.DELETE("/environments/:id", s.environmentHandler.Delete)
	authGroup.GET("/apps", s.appServiceHandler.ListByOrganization)
	authGroup.GET("/apps/:id", s.appServiceHandler.Get, serviceAuth)
	authGroup.PUT("/apps/:id", s.appServiceHandler.Update, serviceAuthAdmin)
	authGroup.DELETE("/apps/:id", s.appServiceHandler.Delete, serviceAuthOwner)
	authGroup.POST("/apps/:id/stop", s.appServiceHandler.StopService, serviceAuthAdmin)
	authGroup.POST("/apps/:id/redeploy", s.appServiceHandler.RedeployService, serviceAuthAdmin)
	authGroup.POST("/apps/:id/restart", s.appServiceHandler.RestartService, serviceAuthAdmin)
	appsGroup := authGroup.Group("/apps")

	appsGroup.GET("/:id/webhooks", s.appServiceHandler.ListWebhooks, serviceAuth)
	appsGroup.POST("/:id/webhooks", s.appServiceHandler.CreateWebhook, serviceAuthAdmin)
	appsGroup.DELETE("/:id/webhooks/:webhookId", s.appServiceHandler.DeleteWebhook, serviceAuthAdmin)
	appsGroup.GET("/:id/volumes", s.appServiceHandler.ListVolumes, serviceAuth)
	appsGroup.POST("/:id/volumes", s.appServiceHandler.CreateVolume, serviceAuthAdmin)
	appsGroup.DELETE("/:id/volumes/:volumeId", s.appServiceHandler.DeleteVolume, serviceAuthAdmin)
	appsGroup.GET("/:id/log-drains", s.appServiceHandler.ListLogDrains, serviceAuth)
	appsGroup.POST("/:id/log-drains", s.appServiceHandler.CreateLogDrain, serviceAuthAdmin)
	appsGroup.DELETE("/:id/log-drains/:drainId", s.appServiceHandler.DeleteLogDrain, serviceAuthAdmin)

	authGroup.GET("/services/:serviceId/variables", s.serviceVarHandler.List, serviceAuth)
	authGroup.GET("/services/:serviceId/env-suggestions", s.serviceVarHandler.Suggest, serviceAuth)
	authGroup.POST("/services/:serviceId/variables", s.serviceVarHandler.Create, serviceAuthAdmin)
	authGroup.PUT("/services/:serviceId/variables/:id", s.serviceVarHandler.Update, serviceAuthAdmin)
	authGroup.DELETE("/services/:serviceId/variables/:id", s.serviceVarHandler.Delete, serviceAuthAdmin)

	authGroup.GET("/services/:serviceId/serverless/code", s.serverlessHandler.GetCode, serviceAuth)
	authGroup.POST("/services/:serviceId/serverless/code", s.serverlessHandler.SaveCode, serviceAuthAdmin)

	apiGroup.GET("/services/:serviceId/logs", s.serviceLogsWSHandler.Handle)
}

func (s *Server) registerDeploymentRoutes(authGroup *echo.Group) {
	serviceAuthAdmin := s.RequireServiceRole(models.MemberPermissionAdmin)
	serviceAuth := s.RequireServiceRole("")

	authGroup.GET("/services/:serviceId/deployments", s.deploymentHandler.ListServiceDeployments, serviceAuth)
	authGroup.GET("/services/:serviceId/previews", s.deploymentHandler.ListPRPreviews, serviceAuth)
	authGroup.POST("/services/:serviceId/deploy", s.deploymentHandler.Trigger, serviceAuthAdmin)
	authGroup.POST("/deployments/:id/rollback", s.deploymentHandler.Rollback)
	authGroup.GET("/deployments/:id/logs", s.deploymentHandler.GetLogs, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/deployments/:id/explain", s.deploymentHandler.ExplainFailure)
	authGroup.GET("/services/:serviceId/metrics", s.deploymentHandler.GetMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/metrics/historical", s.metricsHandler.GetHistoricalMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/logs/historical", s.logHandler.GetHistoricalLogs, serviceAuth)
}

func (s *Server) registerBackupRoutes(authGroup *echo.Group) {
	authGroup.GET("/backups", s.backupHandler.List)
	authGroup.POST("/backups", s.backupHandler.Create)
	authGroup.GET("/backups/:id", s.backupHandler.Get)
	authGroup.PUT("/backups/:id", s.backupHandler.Update, s.authGuard.RequireScope("backup:write"))
	authGroup.DELETE("/backups/:id", s.backupHandler.Delete)
	authGroup.POST("/backups/:id/trigger", s.backupHandler.Trigger)
	authGroup.POST("/backups/:id/restore", s.backupHandler.Restore)
	authGroup.GET("/backups/:id/records", s.backupHandler.ListRecords)
	authGroup.GET("/backups/:id/records/:recordId/download", s.backupHandler.DownloadRecord)
	authGroup.DELETE("/backups/:id/records/:recordId", s.backupHandler.DeleteRecord)
	authGroup.GET("/s3-destinations", s.backupHandler.ListS3Destinations, s.authGuard.RequireRole("admin"))
	authGroup.POST("/s3-destinations", s.backupHandler.CreateS3Destination, s.authGuard.RequireRole("admin"))
	authGroup.PUT("/s3-destinations/:id", s.backupHandler.UpdateS3Destination, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/s3-destinations/:id", s.backupHandler.DeleteS3Destination, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerSettingsRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/settings", s.settingsHandler.GetSettings)
	apiGroup.PUT("/settings", s.settingsHandler.UpdateSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/ai", s.aiSettingsHandler.GetAISettings)
	authGroup.POST("/ai/diagnose", s.aiSettingsHandler.DiagnoseLogs, s.aiRateLimiter.Middleware)
	apiGroup.PUT("/ai", s.aiSettingsHandler.UpdateAISettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/notifications", s.notifSettingsHandler.GetNotificationSettings)
	apiGroup.PUT("/notifications", s.notifSettingsHandler.UpdateNotificationSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/updates/status", s.updaterHandler.GetUpdateStatus)
	apiGroup.POST("/settings/updates/check", s.updaterHandler.CheckUpdate, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/updates/deploy", s.updaterHandler.DeployUpdate, s.authGuard.RequireRole("admin"))
	apiGroup.GET("/settings/oauth/providers", s.oauthHandler.ListProviders, s.authGuard.RequireRole("admin"))
	apiGroup.PUT("/settings/oauth/providers", s.oauthHandler.SaveProvider, s.authGuard.RequireRole("admin"))

	apiGroup.POST("/settings/git_apps/github/manifest-callback", s.gitAppsHandler.ExchangeGithubManifestCode, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/git_apps/github", s.gitAppsHandler.ListGithubApps)
	authGroup.GET("/settings/git_apps/github/:id", s.gitAppsHandler.GetGithubApp)
	apiGroup.PUT("/settings/git_apps/github", s.gitAppsHandler.SaveGithubApp, s.authGuard.RequireRole("admin"))
	apiGroup.DELETE("/settings/git_apps/github/:id", s.gitAppsHandler.DeleteGithubApp, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/notifications/test", s.notificationHandler.TestNotification, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerMiscRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.POST("/compose/deploy", s.composeHandler.Deploy)
	authGroup.POST("/compose/analyze", s.composeHandler.Analyze)
	authGroup.POST("/deploy/archive", s.archiveHandler.DeployArchive)
	authGroup.GET("/examples", s.exampleHandler.List)
	authGroup.GET("/one-click", s.oneClickHandler.List)
	authGroup.POST("/one-click/deploy", s.oneClickHandler.Deploy)
	authGroup.POST("/dns", s.dnsHandler.Create, s.authGuard.RequireRole("admin"))
	authGroup.GET("/dns", s.dnsHandler.List, s.authGuard.RequireRole("admin"))
	authGroup.PUT("/dns/:id", s.dnsHandler.Update, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/dns/:id", s.dnsHandler.Delete, s.authGuard.RequireRole("admin"))
	authGroup.GET("/scheduled-tasks", s.scheduledTaskHandler.ListProjectScheduledTasks)
	authGroup.POST("/scheduled-tasks", s.scheduledTaskHandler.Create)
	authGroup.GET("/scheduled-tasks/:id", s.scheduledTaskHandler.Get)
	authGroup.DELETE("/scheduled-tasks/:id", s.scheduledTaskHandler.Delete)
	authGroup.POST("/scheduled-tasks/:id/trigger", s.scheduledTaskHandler.Run)
	authGroup.POST("/git/connect", s.gitHandler.Connect)
	authGroup.GET("/git/status", s.gitHandler.Status)
	authGroup.DELETE("/git/connect/:provider", s.gitHandler.Disconnect)
	authGroup.GET("/git/repos", s.gitHandler.ListRepos)
	apiGroup.POST("/webhooks/git/services/:serviceId", s.webhookHandler.HandleServiceWebhook)
	apiGroup.POST("/webhooks/github/services/:serviceId", s.webhookHandler.HandleGitHubWebhook)
	authGroup.GET("/canvas/projects", s.canvasHandler.ListCanvasSummaries)
	authGroup.GET("/projects/:id/summary", s.canvasHandler.GetCanvasSummary)
	authGroup.GET("/environments/:id/canvas", s.canvasHandler.GetEnvironmentCanvas)
	authGroup.GET("/audit-logs", s.auditLogHandler.List, s.authGuard.RequireRole("admin"))
	authGroup.GET("/mcp/sse", s.HandleMCPSSE)
	authGroup.POST("/mcp/messages", s.HandleMCPMessage)
	apiGroup.GET("/ws/terminal/:id", s.terminalHandler.HandleWebSocket)
	apiGroup.GET("/ws/services/:id/terminal", s.terminalHandler.HandleWebSocket)
}

func (s *Server) registerBillingRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.POST("/billing/webhook", s.billingHandler.Webhook)

	authGroup.GET("/billing/config", s.billingHandler.GetConfig)
	authGroup.POST("/billing/checkout", s.billingHandler.CreateCheckoutSession)
}

func (s *Server) setupSPAFallback() {
	staticDir := config.Get().Server.StaticDir

	if staticDir != "" {
		if stat, err := os.Stat(staticDir); err == nil && stat.IsDir() {
			s.router.GET("/*", func(c echo.Context) error {
				path := filepath.Join(staticDir, filepath.Clean(c.Request().URL.Path))
				if _, err := os.Stat(path); os.IsNotExist(err) {
					return c.File(filepath.Join(staticDir, "index.html"))
				}
				return c.File(path)
			})
			return
		}
	}

	dashboard.RegisterHandlers(s.router)
}
