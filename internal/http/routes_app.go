package http

import (
	"os"
	"path/filepath"

	"codedock/apps/dashboard"
	"codedock/internal/config"
	"codedock/internal/modules/actions"
	"codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/billing"
	"codedock/internal/modules/chatapps"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/dns/ingress"
	gitmod "codedock/internal/modules/git"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/operations"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/security"
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
	settings.RegisterRoutes(apiGroup, authGroup, s.settingsHandler, s.authGuard.RequireRole("admin"), s.appearanceHandler)
	ai.RegisterRoutes(authGroup, apiGroup, s.aiSettingsHandler, s.authGuard.RequireRole("admin"), s.aiRateLimiter.Middleware)
	notifsmod.RegisterRoutes(apiGroup, authGroup, s.notificationHandler, s.notifSettingsHandler, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerMiscRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.GET("/meta", s.getMeta)
	authGroup.GET("/docs/openapi.json", serveOpenAPISpec)

	// Security Advisories & Non-Invasive Runtime Audits
	security.RegisterRoutes(authGroup, s.securityHandler)

	// Chat Apps Integrations (Slack, Discord, Telegram, Teams)
	chatapps.RegisterRoutes(apiGroup, authGroup, s.chatAppHandler)

	dns.RegisterRoutes(authGroup, s.dnsHandler, s.authGuard.RequireRole("admin"))
	ingress.RegisterRoutes(authGroup, s.ingressHandler, s.authGuard.RequireRole("admin"))
	actions.RegisterRoutes(authGroup, s.authGuard.RequireRole("admin"), s.actionsHandler)
	support.RegisterRoutes(authGroup, s.supportHandler)
	tasks.RegisterRoutes(authGroup, s.scheduledTaskHandler)
	audit.RegisterRoutes(authGroup, s.auditLogHandler, s.authGuard.RequireRole("admin"))
	s.registerMcpAndAiRoutes(authGroup)
}

func (s *Server) registerBillingRoutes(apiGroup, authGroup *echo.Group) {
	billing.RegisterRoutes(apiGroup, authGroup, s.billingHandler, config.Get().Cloud.Enabled)
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
