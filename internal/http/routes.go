package http

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/modules/auth"
	"codedock/internal/modules/databases"
	mailmodule "codedock/internal/modules/mail"
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
	s.registerDeploymentRoutes(authGroup)
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
	auth.RegisterAuthRoutes(apiGroup, authGroup, auth.AuthRouteDeps{
		AuthHandler:       s.authHandler,
		UserHandler:       s.userHandler,
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
	auth.RegisterUserRoutes(authGroup, s.userHandler, s.authHandler, s.authGuard.RequireRole("admin"), s.otpRateLimiter.Middleware)
}

func (s *Server) registerProjectRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/projects", s.projectHandler.ListProjects)
	authGroup.POST("/projects", s.projectHandler.CreateProject)

	projectAuth := s.authGuard.RequireProjectRole("")
	projectAuthAdmin := s.authGuard.RequireProjectRole(models.MemberPermissionAdmin)
	projectAuthOwner := s.authGuard.RequireProjectRole(models.MemberPermissionOwner)

	authGroup.GET("/projects/:id", s.projectHandler.GetProject, projectAuth)
	authGroup.DELETE("/projects/:id", s.projectHandler.DeleteProject, projectAuthOwner)
	authGroup.GET("/projects/:id/members", s.projectHandler.ListMembers, projectAuth)
	authGroup.POST("/projects/:id/members", s.projectHandler.AddMember, projectAuthAdmin)
	authGroup.DELETE("/projects/:id/members/:memberId", s.projectHandler.RemoveMember, projectAuthAdmin)

	authGroup.GET("/domains", s.domainHandler.ListAll)
	authGroup.GET("/projects/:id/domains", s.domainHandler.ListByProject, projectAuth)
	authGroup.GET("/services/:id/domains", s.domainHandler.ListByService)
	authGroup.POST("/services/:id/domains", s.domainHandler.Create)
	authGroup.DELETE("/domains/:id", s.domainHandler.Delete)
	authGroup.POST("/domains/:id/verify", s.domainHandler.Verify)
	authGroup.GET("/domains/:id/verify", s.domainHandler.Verify)

	authGroup.GET("/projects/:id/route-protection", s.routeProtectionHandler.List, projectAuth)
	authGroup.POST("/projects/:id/route-protection", s.routeProtectionHandler.Change, projectAuth)

	authGroup.GET("/projects/:id/env", s.projectEnvHandler.GetVars, projectAuth, s.authGuard.RequireScope("env:read"))
	authGroup.PUT("/projects/:id/env", s.projectEnvHandler.SetVars, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.POST("/projects/:id/env", s.projectEnvHandler.SetVars, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.POST("/projects/:id/env/migrate-commands", s.projectEnvHandler.MigrateCommands, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.POST("/projects/:id/environments", s.environmentHandler.Create, projectAuthAdmin)
	authGroup.GET("/projects/:id/environments", s.environmentHandler.ListByProject, projectAuth)
	authGroup.GET("/projects/:id/apps", s.appServiceHandler.ListByProject, projectAuth)
	authGroup.GET("/projects/:id/services", s.appServiceHandler.ListByProject, projectAuth)
	authGroup.GET("/projects/:id/deployments", s.deploymentHandler.ListProjectDeployments, projectAuth)
	authGroup.POST("/projects/:id/deploy", s.deploymentHandler.TriggerProject, projectAuthAdmin)
	authGroup.GET("/projects/:id/rollback-capacity", s.rollbackSettingsHandler.GetRollbackCapacity, projectAuth)
	authGroup.PATCH("/projects/:id/rollback-settings", s.rollbackSettingsHandler.UpdateRollbackSettings, projectAuthAdmin)

	authGroup.GET("/project-apps", s.projectAppHandler.List)
	authGroup.POST("/project-apps", s.projectAppHandler.Create)
	authGroup.GET("/project-apps/:id", s.projectAppHandler.Get)
	authGroup.PUT("/project-apps/:id", s.projectAppHandler.Update)
	authGroup.DELETE("/project-apps/:id", s.projectAppHandler.Delete)

	authGroup.GET("/projects/:projectId/tokens", s.projectSettingsHandler.ListTokens, projectAuthAdmin, s.authGuard.RequireScope("env:read"))
	authGroup.POST("/projects/:projectId/tokens", s.projectSettingsHandler.CreateToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.DELETE("/projects/:projectId/tokens/:id", s.projectSettingsHandler.DeleteToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))

	authGroup.GET("/projects/:projectId/registries", s.registryHandler.List, projectAuthAdmin)
	authGroup.POST("/projects/:projectId/registries", s.registryHandler.Create, projectAuthAdmin)
	authGroup.DELETE("/projects/:projectId/registries/:id", s.registryHandler.Delete, projectAuthAdmin)

	authGroup.GET("/services/:serviceId/route-rules", s.routeRuleHandler.List, s.RequireServiceRole(""))
	authGroup.POST("/services/:serviceId/route-rules", s.routeRuleHandler.Create, s.RequireServiceRole(models.MemberPermissionAdmin))
	authGroup.PATCH("/services/:serviceId/route-rules/:ruleId", s.routeRuleHandler.Update, s.RequireServiceRole(models.MemberPermissionAdmin))
	authGroup.DELETE("/services/:serviceId/route-rules/:ruleId", s.routeRuleHandler.Delete, s.RequireServiceRole(models.MemberPermissionAdmin))
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
	auth.RegisterOrganizationRoutes(
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
