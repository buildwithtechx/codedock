package http

import (
	"net/http"
	"os"
	"path/filepath"

	"codedock/apps/dashboard"
	"codedock/internal/config"
	"codedock/internal/modules/actions"
	"codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/chatapps"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/dns/ingress"
	gitmod "codedock/internal/modules/git"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/operations"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/settings"
	"codedock/internal/modules/support"
	"codedock/internal/modules/tasks"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerAppRoutes(apiGroup, authGroup *echo.Group) {
	appsmod.RegisterRoutes(authGroup, appsmod.RoutesConfig{
		Handler:            s.appServiceHandler,
		VariableHandler:    s.appVarHandler,
		VolumeHandler:      s.appVolumeHandler,
		WebhookHandler:     s.appWebhookHandler,
		AutoscalingHandler: s.autoscalingHandler,
	})

	authGroup.GET("/projects/:id/runtime-clusters", s.clusterHandler.List, s.RequireProjectRole(organizations.MemberPermissionAdmin))
	apiGroup.GET("/services/:serviceId/logs", s.serviceLogsWSHandler.Handle)
}

func (s *Server) registerDeploymentRoutes(apiGroup, authGroup *echo.Group) {
	serviceAuth := s.RequireServiceRole("")

	deploymentsmod.RegisterRoutes(
		apiGroup,
		authGroup,
		s.deploymentHandler,
		s.deploymentLifecycleHandler,
		s.archiveHandler,
		s.registryHandler,
		s.terminalHandler,
		deploymentsmod.RouteMiddleware{
			Admin:              s.authGuard.RequireRole("admin"),
			ServiceAuth:        serviceAuth,
			ServiceAdmin:       s.RequireServiceRole(organizations.MemberPermissionAdmin),
			ProjectAuth:        s.RequireProjectRole(""),
			ProjectAdmin:       s.RequireProjectRole(organizations.MemberPermissionAdmin),
			RequireProjectRole: func(r string) echo.MiddlewareFunc { return s.RequireProjectRole(organizations.MemberPermission(r)) },
			RequireScope:       s.authGuard.RequireScope,
		},
	)

	gitmod.RegisterRoutes(
		apiGroup,
		authGroup,
		s.gitHandler,
		s.gitInspectionHandler,
		s.gitAppsHandler,
		s.webhookHandler,
		gitmod.RouteMiddleware{
			Admin:              s.authGuard.RequireRole("admin"),
			RequireProjectRole: func(r string) echo.MiddlewareFunc { return s.RequireProjectRole(organizations.MemberPermission(r)) },
		},
	)

	if s.metricsHandler != nil {
		authGroup.GET("/services/:serviceId/metrics/historical", s.metricsHandler.GetHistoricalMetrics, serviceAuth)
	}
	if s.logHandler != nil {
		authGroup.GET("/services/:serviceId/logs/historical", s.logHandler.GetHistoricalLogs, serviceAuth)
	}
}

func (s *Server) registerBackupRoutes(authGroup *echo.Group) {
	operations.RegisterRoutes(authGroup, s.operationHandler, s.authGuard.RequireScope("backup:read"), s.authGuard.RequireScope("backup:write"))
	backupmodule.RegisterRoutes(authGroup, s.backupHandler, backupmodule.RouteMiddleware{
		Admin:          s.authGuard.RequireRole("admin"),
		BackupWrite:    s.authGuard.RequireScope("backup:write"),
		DatabaseManage: s.authGuard.RequireScope("database:manage"),
	})
}

func (s *Server) registerSettingsRoutes(apiGroup, authGroup *echo.Group) {
	settings.RegisterRoutes(apiGroup, authGroup, s.settingsHandler, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/appearance", s.appearanceHandler.GetAppearance)
	authGroup.PUT("/settings/appearance", s.appearanceHandler.SetAppearance)
	ai.RegisterRoutes(authGroup, apiGroup, s.aiSettingsHandler, s.authGuard.RequireRole("admin"), s.aiRateLimiter.Middleware)
	notifsmod.RegisterRoutes(apiGroup, authGroup, s.notificationHandler, s.notifSettingsHandler, s.authGuard.RequireRole("admin"))
	apiGroup.GET("/settings/oauth/providers", s.oauthHandler.ListProviders, s.authGuard.RequireRole("admin"))
	apiGroup.PUT("/settings/oauth/providers", s.oauthHandler.SaveProvider, s.authGuard.RequireRole("admin"))
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

	authGroup.GET("/examples", s.exampleHandler.List)
	dns.RegisterRoutes(authGroup, s.dnsHandler, s.authGuard.RequireRole("admin"))
	ingress.RegisterRoutes(authGroup, s.ingressHandler, s.authGuard.RequireRole("admin"))
	actions.RegisterRoutes(authGroup, s.authGuard.RequireRole("admin"), s.actionsHandler)
	support.RegisterRoutes(authGroup, s.supportHandler)
	tasks.RegisterRoutes(authGroup, s.scheduledTaskHandler)
	audit.RegisterRoutes(authGroup, s.auditLogHandler, s.authGuard.RequireRole("admin"))
	s.registerMcpAndAiRoutes(authGroup)
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
