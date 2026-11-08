package http

import (
	"net/http"
	"os"
	"path/filepath"

	"codedock/apps/dashboard"
	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/modules/ai"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/chatapps"
	"codedock/internal/modules/dns"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerAppRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.POST("/projects/:id/repository-inspection", s.gitHandler.InspectRepository, s.RequireProjectRole(models.MemberPermissionAdmin))
	serviceAuthAdmin := s.RequireServiceRole(models.MemberPermissionAdmin)
	serviceAuthOwner := s.RequireServiceRole(models.MemberPermissionOwner)
	serviceAuth := s.RequireServiceRole("")

	authGroup.GET("/environments/:id/apps", s.appServiceHandler.ListByEnvironment)
	authGroup.POST("/environments/:id/apps", s.appServiceHandler.Create)
	authGroup.DELETE("/environments/:id", s.environmentHandler.Delete)
	authGroup.GET("/apps", s.appServiceHandler.ListByOrganization)
	authGroup.GET("/apps/:id", s.appServiceHandler.Get, serviceAuth)
	authGroup.GET("/apps/:id/runtime", s.appServiceHandler.GetRuntime, serviceAuth)
	authGroup.GET("/apps/:id/runtime/operation", s.appServiceHandler.RuntimeOperation, serviceAuthAdmin)
	authGroup.GET("/projects/:id/runtime-clusters", s.clusterHandler.List, s.RequireProjectRole(models.MemberPermissionAdmin))
	authGroup.GET("/projects/:id/runtime-updates", s.appServiceHandler.CheckProjectRuntimeUpdates, s.RequireProjectRole(""))
	authGroup.POST("/projects/:id/runtime-updates", s.appServiceHandler.ApplyProjectRuntimeUpdates, s.RequireProjectRole(models.MemberPermissionAdmin))
	authGroup.GET("/projects/:id/deployments/admission", s.checksHandler.GetAdmission)
	authGroup.GET("/projects/:id/checks", s.checksHandler.GetChecks)
	authGroup.PUT("/projects/:id/checks", s.checksHandler.UpdateChecks, s.RequireProjectRole(models.MemberPermissionAdmin))
	authGroup.POST("/apps/:id/runtime/review", s.appServiceHandler.ReviewRuntime, serviceAuthAdmin)
	authGroup.POST("/apps/:id/runtime/apply", s.appServiceHandler.ApplyRuntime, serviceAuthAdmin)
	authGroup.GET("/apps/:id/runtime/observe", s.appServiceHandler.ObserveRuntime, serviceAuth, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/apps/:id/runtime/logs", s.appServiceHandler.RuntimeLogs, serviceAuth, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/apps/:id/runtime/logs/stream", s.appServiceHandler.RuntimeLogStream, serviceAuth, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/apps/:id/runtime/graph", s.appServiceHandler.InstanceGraph, serviceAuth, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/apps/:id/runtime/resources", s.appServiceHandler.ResourceSummary, serviceAuth)
	authGroup.GET("/apps/:id/runtime/exec-terminal", s.appServiceHandler.RuntimeExecTerminal, serviceAuthAdmin, s.authGuard.RequireScope("terminal:write"))
	authGroup.POST("/apps/:id/runtime/exec", s.appServiceHandler.RuntimeExec, serviceAuthAdmin, s.authGuard.RequireScope("terminal:write"))
	authGroup.PUT("/apps/:id", s.appServiceHandler.Update, serviceAuthAdmin)
	authGroup.PUT("/apps/:id/build", s.appServiceHandler.UpdateBuildSettings, serviceAuthAdmin)
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
	appsGroup.POST("/:id/volumes/:volumeId/empty", s.appServiceHandler.EmptyVolume, serviceAuthAdmin)
	appsGroup.GET("/:id/runtime-updates", s.appServiceHandler.CheckAppRuntimeUpdates, serviceAuth)
	appsGroup.POST("/:id/runtime-updates", s.appServiceHandler.ApplyAppRuntimeUpdates, serviceAuthAdmin)
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

	authGroup.GET("/deployments", s.deploymentHandler.ListOrganizationDeployments)
	authGroup.GET("/deployments/:id", s.deploymentHandler.Get)
	authGroup.GET("/services/:serviceId/deployments", s.deploymentHandler.ListServiceDeployments, serviceAuth)
	authGroup.GET("/services/:serviceId/previews", s.deploymentHandler.ListPRPreviews, serviceAuth)
	authGroup.POST("/services/:serviceId/deploy", s.deploymentHandler.Trigger, serviceAuthAdmin)
	authGroup.POST("/deployments/:id/rollback", s.deploymentHandler.Rollback)
	authGroup.POST("/deployments/:id/retry", s.deploymentHandler.Retry)
	authGroup.POST("/deployments/:id/pin", s.deploymentHandler.TogglePin)
	authGroup.POST("/deployments/:id/cancel", s.deploymentHandler.Cancel)
	authGroup.GET("/deployments/:id/logs", s.deploymentHandler.GetLogs, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/deployments/:id/explain", s.deploymentHandler.ExplainFailure)
	authGroup.GET("/services/:serviceId/metrics", s.deploymentHandler.GetMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/metrics/historical", s.metricsHandler.GetHistoricalMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/autoscaling", s.autoscalingHandler.Get, s.RequireServiceRole(""))
	authGroup.PUT("/services/:serviceId/autoscaling", s.autoscalingHandler.Save, s.RequireServiceRole(models.MemberPermissionAdmin))
	authGroup.GET("/services/:serviceId/logs/historical", s.logHandler.GetHistoricalLogs, serviceAuth)
}

func (s *Server) registerBackupRoutes(authGroup *echo.Group) {
	authGroup.GET("/operations", s.operationHandler.List, s.authGuard.RequireScope("backup:read"))
	authGroup.GET("/operations/:operationId", s.operationHandler.Get, s.authGuard.RequireScope("backup:read"))
	authGroup.POST("/operations/:operationId/cancel", s.operationHandler.Cancel, s.authGuard.RequireScope("backup:write"))
	backupmodule.RegisterRoutes(authGroup, s.backupHandler, backupmodule.RouteMiddleware{
		Admin:          s.authGuard.RequireRole("admin"),
		BackupWrite:    s.authGuard.RequireScope("backup:write"),
		DatabaseManage: s.authGuard.RequireScope("database:manage"),
	})
}

func (s *Server) registerSettingsRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/settings", s.settingsHandler.GetSettings)
	apiGroup.PUT("/settings", s.settingsHandler.UpdateSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/appearance", s.appearanceHandler.GetAppearance)
	authGroup.PUT("/settings/appearance", s.appearanceHandler.SetAppearance)
	ai.RegisterRoutes(authGroup, apiGroup, s.aiSettingsHandler, s.authGuard.RequireRole("admin"), s.aiRateLimiter.Middleware)
	authGroup.GET("/notifications", s.notifSettingsHandler.GetNotificationSettings)
	apiGroup.PUT("/notifications", s.notifSettingsHandler.UpdateNotificationSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/notifications/subscriptions", s.notifSettingsHandler.ListSubscriptions)
	authGroup.PUT("/notifications/subscriptions", s.notifSettingsHandler.UpsertSubscription)
	authGroup.GET("/notifications/defaults", s.notifSettingsHandler.ListDefaults)
	apiGroup.PUT("/notifications/defaults", s.notifSettingsHandler.UpsertDefault, s.authGuard.RequireRole("admin"))
	authGroup.GET("/notifications/deliveries", s.notifSettingsHandler.ListDeliveries)
	authGroup.GET("/notifications/devices", s.notifSettingsHandler.ListDeviceTokens)
	authGroup.POST("/notifications/devices", s.notifSettingsHandler.RegisterDeviceToken)
	authGroup.DELETE("/notifications/devices/:token", s.notifSettingsHandler.DeleteDeviceToken)
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
	apiGroup.GET("/meta", s.getMeta)
	authGroup.GET("/docs/openapi.json", serveOpenAPISpec)

	// Security Advisories & Non-Invasive Runtime Audits
	authGroup.GET("/updates/advisories", s.securityHandler.GetAdvisories)
	authGroup.GET("/updates/:projectId/security", s.securityHandler.GetProjectSecurity)
	authGroup.POST("/updates/:projectId/security/check", s.securityHandler.CheckProject)
	authGroup.PUT("/updates/:projectId/security/policy", s.securityHandler.SetProjectPolicy)
	authGroup.POST("/updates/:projectId/security/apply", s.securityHandler.ApplyRemediation)

	// Chat Apps Integrations (Slack, Discord, Telegram, Teams)
	chatapps.RegisterRoutes(apiGroup, authGroup, s.chatAppHandler)

	authGroup.GET("/projects/:id/stacks", s.composeStackHandler.List)
	authGroup.GET("/projects/:id/stacks/:stackId/config", s.composeStackHandler.Config, s.authGuard.RequireScope("env:read"))
	authGroup.POST("/projects/:id/stacks/review", s.composeStackHandler.Review)
	authGroup.POST("/projects/:id/stacks", s.composeStackHandler.Save)
	authGroup.POST("/projects/:id/stacks/:stackId/deploy", s.composeStackHandler.Deploy)
	authGroup.POST("/projects/:id/stacks/:stackId/cancel", s.composeStackHandler.Cancel)
	authGroup.POST("/projects/:id/stacks/:stackId/restart", s.composeStackHandler.Restart)
	authGroup.POST("/projects/:id/stacks/:stackId/stop", s.composeStackHandler.Stop)
	authGroup.PUT("/environments/:id/canvas", s.canvasHandler.ApplyTopology)
	authGroup.GET("/environments/:id/canvas/observe", s.canvasHandler.ObserveResource, s.authGuard.RequireScope("logs:read"))
	authGroup.POST("/compose/deploy", s.composeHandler.Deploy)
	authGroup.POST("/compose/analyze", s.composeHandler.Analyze)
	authGroup.POST("/deploy/archive", s.archiveHandler.DeployArchive)
	authGroup.GET("/examples", s.exampleHandler.List)
	authGroup.GET("/one-click", s.oneClickHandler.List)
	authGroup.GET("/one-click/:id", s.oneClickHandler.Get)
	authGroup.POST("/one-click/review", s.oneClickHandler.Review)
	authGroup.POST("/one-click/deploy", s.oneClickHandler.Deploy)
	authGroup.POST("/one-click/custom", s.oneClickHandler.AddCustom)
	authGroup.DELETE("/one-click/custom/:id", s.oneClickHandler.DeleteCustom)
	dns.RegisterRoutes(authGroup, s.dnsHandler, s.authGuard.RequireRole("admin"))
	authGroup.GET("/actions/workflows", s.actionsHandler.ListWorkflows)
	authGroup.POST("/actions/workflows", s.actionsHandler.CreateWorkflow, s.authGuard.RequireRole("admin"))
	authGroup.POST("/actions/workflows/:id/dispatch", s.actionsHandler.DispatchWorkflow)
	authGroup.DELETE("/actions/workflows/:id", s.actionsHandler.DeleteWorkflow, s.authGuard.RequireRole("admin"))
	authGroup.GET("/actions/runs", s.actionsHandler.ListRuns)
	authGroup.GET("/actions/runs/:id", s.actionsHandler.GetRun)
	authGroup.POST("/actions/runs/:id/retry", s.actionsHandler.RetryRun)
	authGroup.POST("/actions/runs/:id/cancel", s.actionsHandler.CancelRun)
	authGroup.DELETE("/actions/runs/:id", s.actionsHandler.DeleteRun, s.authGuard.RequireRole("admin"))
	authGroup.GET("/actions/runners", s.actionsHandler.ListRunners)
	authGroup.POST("/actions/runners", s.actionsHandler.CreateRunner, s.authGuard.RequireRole("admin"))
	authGroup.POST("/actions/runners/inspect", s.actionsHandler.InspectDestination, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/actions/runners/:id", s.actionsHandler.DeleteRunner, s.authGuard.RequireRole("admin"))
	authGroup.GET("/dns/ingress", s.ingressHandler.List)
	authGroup.POST("/dns/ingress/discover", s.ingressHandler.Discover)
	authGroup.POST("/dns/ingress", s.ingressHandler.Connect, s.authGuard.RequireRole("admin"))
	authGroup.PUT("/dns/ingress/:id", s.ingressHandler.Update, s.authGuard.RequireRole("admin"))
	authGroup.POST("/dns/ingress/:id/retry", s.ingressHandler.Retry, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/dns/ingress/:id", s.ingressHandler.Disconnect, s.authGuard.RequireRole("admin"))
	authGroup.GET("/dns/ingress/domains/:domainId", s.ingressHandler.GetDomain)
	authGroup.PUT("/dns/ingress/domains/:domainId", s.ingressHandler.SetDomain, s.authGuard.RequireRole("admin"))
	authGroup.GET("/system/networks", s.serverNetworksHandler.List)
	authGroup.GET("/system/networks/:id", s.serverNetworksHandler.Get)
	authGroup.POST("/system/networks", s.serverNetworksHandler.Create, s.authGuard.RequireRole("admin"))
	authGroup.PUT("/system/networks/:id/access", s.serverNetworksHandler.ReviseAccess, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/system/networks/:id", s.serverNetworksHandler.Delete, s.authGuard.RequireRole("admin"))
	authGroup.GET("/system/networks/operations/:operationId/watchdog", s.serverNetworksHandler.GetWatchdog)
	authGroup.POST("/system/networks/operations/:operationId/watchdog/commit", s.serverNetworksHandler.CommitWatchdog, s.authGuard.RequireRole("admin"))
	authGroup.POST("/system/networks/operations/:operationId/watchdog/rollback", s.serverNetworksHandler.RollbackWatchdog, s.authGuard.RequireRole("admin"))
	authGroup.GET("/cloud/support/session", s.supportHandler.GetSession)
	authGroup.GET("/cloud/support/mine", s.supportHandler.ListTickets)
	authGroup.POST("/cloud/support/mine", s.supportHandler.CreateTicket)
	authGroup.GET("/cloud/support/mine/:id", s.supportHandler.GetTicket)
	authGroup.POST("/cloud/support/mine/:id/replies", s.supportHandler.ReplyTicket)
	authGroup.PATCH("/cloud/support/mine/:id", s.supportHandler.UpdateStatus)
	authGroup.GET("/scheduled-tasks", s.scheduledTaskHandler.ListProjectScheduledTasks)
	authGroup.POST("/scheduled-tasks", s.scheduledTaskHandler.Create)
	authGroup.GET("/scheduled-tasks/:id", s.scheduledTaskHandler.Get)
	authGroup.PUT("/scheduled-tasks/:id", s.scheduledTaskHandler.Update)
	authGroup.DELETE("/scheduled-tasks/:id", s.scheduledTaskHandler.Delete)
	authGroup.POST("/scheduled-tasks/:id/trigger", s.scheduledTaskHandler.Run)
	authGroup.GET("/jobs", s.scheduledTaskHandler.ListProjectScheduledTasks)
	authGroup.POST("/jobs", s.scheduledTaskHandler.Create)
	authGroup.GET("/jobs/:id", s.scheduledTaskHandler.Get)
	authGroup.PUT("/jobs/:id", s.scheduledTaskHandler.Update)
	authGroup.DELETE("/jobs/:id", s.scheduledTaskHandler.Delete)
	authGroup.POST("/jobs/:id/trigger", s.scheduledTaskHandler.Run)
	authGroup.POST("/git/connect", s.gitHandler.Connect)
	authGroup.GET("/git/status", s.gitHandler.Status)
	authGroup.DELETE("/git/connect/:provider", s.gitHandler.Disconnect)
	authGroup.GET("/git/repos", s.gitHandler.ListRepos)
	authGroup.POST("/git/tree", s.gitHandler.RepoTree)
	authGroup.GET("/git/tree", s.gitHandler.RepoTree)
	authGroup.POST("/git/inspect", s.gitHandler.InspectRepository)
	apiGroup.POST("/webhooks/git/services/:serviceId", s.webhookHandler.HandleServiceWebhook)
	apiGroup.POST("/webhooks/github/services/:serviceId", s.webhookHandler.HandleGitHubWebhook)
	apiGroup.POST("/webhooks/github", s.webhookHandler.HandleGitHubWebhook)
	authGroup.GET("/projects/:id/webhooks/deliveries", s.webhookHandler.ListProjectDeliveries)
	authGroup.GET("/canvas/projects", s.canvasHandler.ListCanvasSummaries)
	authGroup.GET("/projects/:id/summary", s.canvasHandler.GetCanvasSummary)
	authGroup.GET("/environments/:id/canvas", s.canvasHandler.GetEnvironmentCanvas)
	auth.RegisterAuditRoutes(authGroup, s.auditLogHandler, s.authGuard.RequireRole("admin"))
	s.registerMcpAndAiRoutes(authGroup)
	apiGroup.GET("/ws/terminal/:id", s.terminalHandler.HandleWebSocket)
	apiGroup.GET("/ws/services/:id/terminal", s.terminalHandler.HandleWebSocket)
}

func (s *Server) registerBillingRoutes(apiGroup, authGroup *echo.Group) {
	if !config.Get().Cloud.Enabled {
		unavailable := func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusNotFound, "billing is unavailable on self-hosted instances")
		}
		apiGroup.Any("/billing", unavailable)
		apiGroup.Any("/billing/*", unavailable)
		return
	}
	apiGroup.POST("/billing/webhook", s.billingHandler.Webhook)

	authGroup.GET("/billing/config", s.billingHandler.GetConfig)
	authGroup.POST("/billing/checkout", s.billingHandler.CreateCheckoutSession)
	authGroup.GET("/billing/pending", s.billingHandler.GetPendingCheckout)
	authGroup.POST("/billing/checkouts/:id/cancel", s.billingHandler.CancelCheckoutSession)
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
