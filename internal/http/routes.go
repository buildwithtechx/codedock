package http

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/modules/auth"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	mailmodule "codedock/internal/modules/mail"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/users"
	"codedock/internal/utils"
)

func (s *Server) registerRoutes() {
	apiGroup := s.router.Group("/api")

	authGroup := apiGroup.Group("")
	authGroup.Use(s.authGuard.RequireAuth())

	s.registerAuthRoutes(apiGroup, authGroup)
	s.registerSystemRoutes(apiGroup, authGroup)
	s.registerUserRoutes(apiGroup, authGroup)
	s.registerProjectRoutes(apiGroup, authGroup)
	s.registerOrganizationRoutes(authGroup)
	s.registerServerRoutes(apiGroup, authGroup)
	s.registerManagedRoutes(authGroup)
	s.registerMigrationRoutes(authGroup)
	s.registerAnalyticsRoutes(authGroup)
	s.registerDatabaseRoutes(authGroup)
	s.registerAppRoutes(apiGroup, authGroup)
	s.registerDeploymentRoutes(apiGroup, authGroup)
	s.registerBackupRoutes(authGroup)
	s.registerClusterRoutes(authGroup)
	s.registerSettingsRoutes(apiGroup, authGroup)
	s.registerMiscRoutes(apiGroup, authGroup)
	s.registerBillingRoutes(apiGroup, authGroup)
	s.registerMailRoutes(authGroup)

	s.router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	s.router.GET("/health/env", s.settingsHandler.GetHealthEnv)
	apiGroup.GET("/health/env", s.settingsHandler.GetHealthEnv)

	s.router.GET("/api/v1/edge/auth", s.edgeAuthHandler.Verify)
	s.router.GET("/_codedock/access", s.edgeAuthHandler.LoginForm)
	s.router.POST("/_codedock/access", s.edgeAuthHandler.Login)

	s.setupSPAFallback()
}

func (s *Server) RequireServiceRole(minPermission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			serviceID := c.Param("serviceId")
			if serviceID == "" {
				serviceID = c.Param("id")
			}
			if serviceID == "" {
				return next(c)
			}

			userClaims := GetUserClaimsFromContext(c.Request().Context())
			if userClaims == nil {
				return utils.Error(c, 401, "unauthorized")
			}
			if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
				return next(c)
			}

			svc, err := s.appService.GetAppService(c.Request().Context(), serviceID)
			if err != nil || svc == nil {
				return utils.Error(c, 404, "service not found")
			}

			if userClaims.Role == "api" {
				if c.Get("project_id") != svc.ProjectID {
					return utils.Error(c, 403, "api token not authorized for this project")
				}
				if minPermission != "" && minPermission != models.MemberPermissionMember {
					return utils.Error(c, 403, "api tokens cannot perform admin/owner actions")
				}
				return next(c)
			}

			if !s.projectService.HasPermission(c.Request().Context(), svc.ProjectID, userClaims.UserID, userClaims.Role, minPermission) {
				return utils.Error(c, 403, "insufficient project permissions")
			}
			return next(c)
		}
	}
}

func (s *Server) registerAuthRoutes(apiGroup, authGroup *echo.Group) {
	uh := s.userHandler
	if uh == nil {
		uh = &users.UserHandler{}
	}
	auth.RegisterAuthRoutes(apiGroup, authGroup, auth.AuthRouteDeps{
		AuthHandler:       s.authHandler,
		GetProfile:        uh.GetProfile,
		OAuthHandler:      s.oauthHandler,
		DeviceFlowHandler: s.deviceFlowHandler,
		PasskeyHandler:    s.passkeyHandler,
		BootstrapCSRF:     bootstrapCSRF,
		AuthRateLimit:     s.authRateLimiter.Middleware,
		OTPRateLimit:      s.otpRateLimiter.Middleware,
	})
}

func (s *Server) registerSystemRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.GET("/github-stars", s.getGithubStars)
	apiGroup.GET("/system/public", s.settingsHandler.GetPublicSettings)
	apiGroup.GET("/system/setup-status", s.onboardingHandler.SetupStatus)
	apiGroup.POST("/system/setup/import", s.migrationHandler.ImportDuringSetup)
	authGroup.GET("/system/stats", s.systemHandler.GetStats)
	authGroup.GET("/system/monitor/stream", s.serverHandler.StreamMonitor, s.authGuard.RequireScope("server:read"))
	apiGroup.POST("/system/restart", s.systemHandler.Restart, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/maintenance/cleanup", s.systemHandler.Cleanup, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/export", s.migrationHandler.Export, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/import", s.migrationHandler.Import, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/takeover/scan", s.takeoverHandler.Scan, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/takeover/adopt", s.takeoverHandler.Adopt, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/local-cli-session", s.systemHandler.CreateLocalCLISession)
	authGroup.GET("/system/takeover/runs", s.takeoverHandler.ListRuns, s.authGuard.RequireRole("admin"))
	authGroup.GET("/system/takeover/runs/:id", s.takeoverHandler.GetRun, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerUserRoutes(apiGroup, authGroup *echo.Group) {
	ah := s.authHandler
	if ah == nil {
		ah = &auth.AuthHandler{}
	}
	users.RegisterRoutes(authGroup, s.userHandler, ah.AdminInviteUser, s.authGuard.RequireRole("admin"), s.otpRateLimiter.Middleware)
}

func (s *Server) registerProjectRoutes(apiGroup, authGroup *echo.Group) {
	if s.projectHandler != nil {
		projects.RegisterRoutes(authGroup, projects.RoutesConfig{
			Handler:             s.projectHandler,
			AppHandler:          s.projectAppHandler,
			SettingsHandler:     s.projectSettingsHandler,
			ChecksHandler:       s.checksHandler,
			CanvasHandler:       s.canvasHandler,
			ComposeStackHandler: s.composeStackHandler,
			ComposeHandler:      s.composeHandler,
			OneClickHandler:     s.oneClickHandler,
		})
	}
	if s.environmentHandler != nil {
		environments.RegisterRoutes(authGroup, environments.RoutesConfig{
			EnvHandler:        s.environmentHandler,
			ProtectionHandler: s.routeProtectionHandler,
			RouteRuleHandler:  s.routeRuleHandler,
		})
	}
	if s.domainHandler != nil {
		domains.RegisterRoutes(authGroup, domains.RoutesConfig{
			Handler: s.domainHandler,
		})
	}
}

func (s *Server) registerServerRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/servers", s.serverHandler.List, s.authGuard.RequireScope("server:read"))
	authGroup.GET("/servers/:id", s.serverHandler.Get, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers", s.serverHandler.Create, s.authGuard.RequireScope("server:write"))
	authGroup.PATCH("/servers/:id", s.serverHandler.Update, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/servers/test-ssh", s.serverHandler.TestSSH, s.authGuard.RequireScope("server:write"))
	authGroup.DELETE("/servers/:id", s.serverHandler.Delete, s.authGuard.RequireScope("server:write"))
	authGroup.GET("/servers/:id/components", s.serverHandler.GetComponents, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers/:id/scan-ports", s.serverHandler.ScanPorts, s.authGuard.RequireScope("server:read"))
	authGroup.GET("/servers/:id/rate-limit", s.serverHandler.GetRateLimit, s.authGuard.RequireScope("server:read"))
	authGroup.PUT("/servers/:id/rate-limit", s.serverHandler.UpdateRateLimit, s.authGuard.RequireScope("server:write"))
	authGroup.GET("/servers/:id/network-settings", s.serverHandler.GetNetworkSettings, s.authGuard.RequireScope("server:read"))
	authGroup.PATCH("/servers/:id/network-settings", s.serverHandler.UpdateNetworkSettings, s.authGuard.RequireScope("server:write"))
	authGroup.GET("/servers/:id/managed/info", s.serverHandler.GetManagedInfo, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers/:id/managed/boot-logs", s.serverHandler.GetManagedBootLogs, s.authGuard.RequireScope("server:read"))
	authGroup.GET("/servers/:id/managed/ssh", s.serverHandler.GetManagedSshStatus, s.authGuard.RequireScope("server:read"))
	authGroup.PATCH("/servers/:id/managed/ssh", s.serverHandler.SetManagedSsh, s.authGuard.RequireScope("server:write"))
	authGroup.PUT("/servers/:id/managed/ssh/key", s.serverHandler.SetManagedSshKey, s.authGuard.RequireScope("server:write"))
	authGroup.PUT("/servers/:id/managed/ssh/password", s.serverHandler.SetManagedSshPassword, s.authGuard.RequireScope("server:write"))
	authGroup.GET("/servers/:id/managed/runtime-api", s.serverHandler.GetManagedRuntimeStatus, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers/:id/managed/runtime-api/enable", s.serverHandler.EnableManagedRuntime, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/servers/:id/managed/runtime-api/credential", s.serverHandler.GetManagedRuntimeCredential, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/servers/:id/managed/runtime-api/rotate", s.serverHandler.RotateManagedRuntimeCredential, s.authGuard.RequireScope("server:write"))
	authGroup.GET("/servers/:id/managed/workloads", s.serverHandler.GetManagedWorkloads, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers/:id/managed/workloads", s.serverHandler.CreateManagedWorkload, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/servers/:id/managed/workloads/control", s.serverHandler.ControlManagedWorkload, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/servers/:id/managed/workloads/logs", s.serverHandler.GetManagedWorkloadLogs, s.authGuard.RequireScope("server:read"))
	apiGroup.GET("/ws/servers/:serverId/metrics", s.serverMetricsWSHandler.Handle)
}

func (s *Server) registerOrganizationRoutes(authGroup *echo.Group) {
	organizations.RegisterRoutes(
		authGroup,
		s.orgHandler,
		s.authGuard.RequireOrgRole(models.MemberPermissionMember),
		s.authGuard.RequireOrgRole(models.MemberPermissionAdmin),
		s.authGuard.RequireOrgRole(models.MemberPermissionOwner),
	)
}

func (s *Server) registerDatabaseRoutes(authGroup *echo.Group) {
	databases.RegisterRoutes(authGroup, s.dbHandler, s.authGuard.RequireScope)
}

func (s *Server) registerMailRoutes(authGroup *echo.Group) {
	mailmodule.RegisterRoutes(authGroup, s.mailHandler)
}
