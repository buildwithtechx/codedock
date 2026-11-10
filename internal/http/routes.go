package http

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/modules/auth"
	"codedock/internal/modules/projects/canvas"
	modclusters "codedock/internal/modules/clusters"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	mailmodule "codedock/internal/modules/mail"
	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	serversmod "codedock/internal/modules/servers"
	systemmod "codedock/internal/modules/system"
	templatesmod "codedock/internal/modules/templates"
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
	s.registerMigrationRoutes(apiGroup, authGroup)
	s.registerAnalyticsRoutes(apiGroup, authGroup)
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

	environments.RegisterEdgeAuthRoutes(s.router, s.edgeAuthHandler)

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
		AdminAuth:         s.authGuard.RequireRole("admin"),
	})
}

func (s *Server) registerSystemRoutes(apiGroup, authGroup *echo.Group) {
	systemmod.RegisterRoutes(authGroup, apiGroup, systemmod.RoutesConfig{
		System:    s.systemHandler,
		Updater:   s.updaterHandler,
		AdminAuth: s.authGuard.RequireRole("admin"),
	})
}

func (s *Server) registerUserRoutes(apiGroup, authGroup *echo.Group) {
	ah := s.authHandler
	if ah == nil {
		ah = &auth.AuthHandler{}
	}
	users.RegisterRoutes(authGroup, s.userHandler, ah.AdminInviteUser, s.authGuard.RequireRole("admin"), s.otpRateLimiter.Middleware)
	users.RegisterOnboardingRoutes(apiGroup, s.onboardingHandler)
}

func (s *Server) registerProjectRoutes(apiGroup, authGroup *echo.Group) {
	if s.projectHandler != nil {
		projects.RegisterRoutes(authGroup, projects.RoutesConfig{
			Handler:         s.projectHandler,
			AppHandler:      s.projectAppHandler,
			SettingsHandler: s.projectSettingsHandler,
			ChecksHandler:   s.checksHandler,
		})
	}
	if s.composeStackHandler != nil {
		composemod.RegisterRoutes(authGroup, composemod.RoutesConfig{
			Handler: s.composeStackHandler,
		})
	}
	if s.oneClickHandler != nil {
		templatesmod.RegisterRoutes(authGroup, templatesmod.RoutesConfig{
			Handler: s.oneClickHandler,
		})
	}
	if s.canvasHandler != nil {
		canvas.RegisterRoutes(authGroup, s.canvasHandler)
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

func (s *Server) registerMigrationRoutes(apiGroup, authGroup *echo.Group) {
	modmigrations.RegisterRoutes(apiGroup, authGroup, modmigrations.RoutesConfig{
		LifecycleHandler: s.migrationLifecycleHandler,
		BundleHandler:    s.migrationHandler,
		Middleware: modmigrations.RouteMiddleware{
			OrgMember:   s.authGuard.RequireOrgRole("member"),
			OrgAdmin:    s.authGuard.RequireOrgRole("admin"),
			ServerRead:  s.authGuard.RequireScope("server:read"),
			ServerWrite: s.authGuard.RequireScope("server:write"),
			Admin:       s.authGuard.RequireRole("admin"),
		},
	})
}

func (s *Server) registerClusterRoutes(group *echo.Group) {
	modclusters.RegisterRoutes(group, modclusters.RoutesConfig{
		Handler:     s.clusterHandler,
		DataHandler: s.clusterDataHandler,
		Middleware: modclusters.RouteMiddleware{
			Admin:          s.authGuard.RequireRole("admin"),
			DatabaseManage: s.authGuard.RequireScope("database:manage"),
		},
	})
}

func (s *Server) registerAnalyticsRoutes(apiGroup, authGroup *echo.Group) {
	serviceAuth := s.RequireServiceRole("")
	monitoring.RegisterRoutes(authGroup, apiGroup, monitoring.RoutesConfig{
		Analytics:     s.analyticsHandler,
		Attention:     s.attentionHandler,
		Metrics:       s.metricsHandler,
		Logs:          s.logHandler,
		ServiceLogsWS: s.serviceLogsWSHandler,
		Example:       s.exampleHandler,
		Middleware: monitoring.RouteMiddleware{
			OrgMember:   s.authGuard.RequireOrgRole("member"),
			OrgAdmin:    s.authGuard.RequireOrgRole("admin"),
			ReadScope:   s.authGuard.RequireScope("server:read"),
			WriteScope:  s.authGuard.RequireScope("server:write"),
			ServiceAuth: serviceAuth,
		},
	})
}
