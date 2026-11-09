package http

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/modules/auth"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	mailmodule "codedock/internal/modules/mail"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	serversmod "codedock/internal/modules/servers"
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

func (s *Server) RequireServiceRole(minPermission organizations.MemberPermission) echo.MiddlewareFunc {
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
			if userClaims.Role == users.UserRoleAdmin || userClaims.Role == users.UserRoleOwner {
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
				if minPermission != "" && minPermission != organizations.MemberPermissionMember {
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
	apiGroup.POST("/system/restart", s.systemHandler.Restart, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/maintenance/cleanup", s.systemHandler.Cleanup, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/export", s.migrationHandler.Export, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/import", s.migrationHandler.Import, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/local-cli-session", s.systemHandler.CreateLocalCLISession)
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
	serversmod.RegisterRoutes(apiGroup, authGroup, serversmod.RoutesConfig{
		Handler:         s.serverHandler,
		ManagedHandler:  s.managedHandler,
		TakeoverHandler: s.takeoverHandler,
		NetworksHandler: s.serverNetworksHandler,
		WSHandler:       s.serverMetricsWSHandler,
		Middleware: serversmod.RouteMiddleware{
			ServerRead:  s.authGuard.RequireScope("server:read"),
			ServerWrite: s.authGuard.RequireScope("server:write"),
			Admin:       s.authGuard.RequireRole("admin"),
			OrgMember:   s.authGuard.RequireOrgRole(organizations.MemberPermissionMember),
			OrgAdmin:    s.authGuard.RequireOrgRole(organizations.MemberPermissionAdmin),
		},
	})
}

func (s *Server) registerOrganizationRoutes(authGroup *echo.Group) {
	organizations.RegisterRoutes(
		authGroup,
		s.orgHandler,
		s.authGuard.RequireOrgRole(organizations.MemberPermissionMember),
		s.authGuard.RequireOrgRole(organizations.MemberPermissionAdmin),
		s.authGuard.RequireOrgRole(organizations.MemberPermissionOwner),
	)
}

func (s *Server) registerDatabaseRoutes(authGroup *echo.Group) {
	databases.RegisterRoutes(authGroup, s.dbHandler, s.authGuard.RequireScope)
}

func (s *Server) registerMailRoutes(authGroup *echo.Group) {
	mailmodule.RegisterRoutes(authGroup, s.mailHandler)
}
