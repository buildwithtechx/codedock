package http

import (
	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
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
	s.registerDatabaseRoutes(authGroup)
	s.registerAppRoutes(apiGroup, authGroup)
	s.registerDeploymentRoutes(authGroup)
	s.registerBackupRoutes(authGroup)
	s.registerSettingsRoutes(apiGroup, authGroup)
	s.registerMiscRoutes(apiGroup, authGroup)
	s.registerBillingRoutes(apiGroup, authGroup)

	s.router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

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
			if userClaims.Role == "admin" {
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
	apiGroup.POST("/auth/signup", s.authHandler.Register, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/signin", s.authHandler.Login, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/refresh", s.authHandler.Refresh)
	apiGroup.POST("/auth/forgot-password", s.authHandler.ForgotPassword, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/reset-password", s.authHandler.ResetPassword, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/email/resend", s.authHandler.ResendVerificationEmail, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/email/verify", s.authHandler.VerifyEmail, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/logout", s.authHandler.Logout)
	authGroup.GET("/auth/me", s.userHandler.GetProfile)

	apiGroup.GET("/auth/oauth/providers/enabled", s.oauthHandler.ListEnabledProviders)
	apiGroup.GET("/auth/oauth/:provider", s.oauthHandler.OAuthRedirect)
	apiGroup.GET("/auth/oauth/:provider/callback", s.oauthHandler.OAuthCallback)
	authGroup.POST("/auth/2fa/setup", s.oauthHandler.Setup2FA)
	authGroup.POST("/auth/2fa/verify", s.oauthHandler.Verify2FA)
	authGroup.POST("/auth/2fa/disable", s.oauthHandler.Disable2FA, s.otpRateLimiter.Middleware)
}

func (s *Server) registerSystemRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.GET("/system/public", s.settingsHandler.GetPublicSettings)
	apiGroup.GET("/system/setup-status", s.onboardingHandler.SetupStatus)
	apiGroup.POST("/system/setup", s.onboardingHandler.Setup)
	authGroup.GET("/system/stats", s.systemHandler.GetStats)
	apiGroup.POST("/system/restart", s.systemHandler.Restart, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/maintenance/cleanup", s.systemHandler.Cleanup, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/export", s.migrationHandler.Export, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/import", s.migrationHandler.Import, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerUserRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/users", s.userHandler.ListUsers, s.authGuard.RequireRole("admin"))
	authGroup.POST("/users/invite", s.authHandler.AdminInviteUser, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/users/:id", s.userHandler.DeleteUser, s.authGuard.RequireRole("admin"))
	authGroup.GET("/profile", s.userHandler.GetProfile)
	authGroup.PUT("/profile", s.userHandler.UpdateProfile)
	authGroup.POST("/profile/email/request", s.userHandler.RequestEmailChange)
	authGroup.POST("/profile/email/verify", s.userHandler.VerifyEmailChange, s.otpRateLimiter.Middleware)
	authGroup.PUT("/profile/password", s.userHandler.ChangePassword)
	authGroup.GET("/profile/tokens", s.userHandler.ListPATs)
	authGroup.POST("/profile/tokens", s.userHandler.CreatePAT)
	authGroup.DELETE("/profile/tokens/:id", s.userHandler.DeletePAT)
}

func (s *Server) registerProjectRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/projects", s.projectHandler.ListProjects)
	authGroup.POST("/projects", s.projectHandler.CreateProject)

	projectAuth := s.authGuard.RequireProjectRole("")
	projectAuthAdmin := s.authGuard.RequireProjectRole(models.MemberPermissionAdmin)
	projectAuthOwner := s.authGuard.RequireProjectRole(models.MemberPermissionOwner)

	authGroup.GET("/projects/:id", s.projectHandler.GetProject, projectAuth)
	authGroup.DELETE("/projects/:id", s.projectHandler.DeleteProject, projectAuthOwner)

	authGroup.GET("/services/:id/domains", s.domainHandler.ListByService)
	authGroup.POST("/services/:id/domains", s.domainHandler.Create)
	authGroup.DELETE("/domains/:id", s.domainHandler.Delete)

	authGroup.GET("/projects/:id/env", s.projectEnvHandler.GetVars, projectAuth, s.authGuard.RequireScope("env:read"))
	authGroup.PUT("/projects/:id/env", s.projectEnvHandler.SetVars, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.POST("/projects/:id/environments", s.environmentHandler.Create, projectAuthAdmin)
	authGroup.GET("/projects/:id/environments", s.environmentHandler.ListByProject, projectAuth)
	authGroup.GET("/projects/:id/apps", s.appServiceHandler.ListByProject, projectAuth)

	authGroup.GET("/projects/:projectId/tokens", s.projectSettingsHandler.ListTokens, projectAuthAdmin, s.authGuard.RequireScope("env:read"))
	authGroup.POST("/projects/:projectId/tokens", s.projectSettingsHandler.CreateToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.DELETE("/projects/:projectId/tokens/:id", s.projectSettingsHandler.DeleteToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))

	authGroup.GET("/projects/:projectId/registries", s.registryHandler.List, projectAuthAdmin)
	authGroup.POST("/projects/:projectId/registries", s.registryHandler.Create, projectAuthAdmin)
	authGroup.DELETE("/projects/:projectId/registries/:id", s.registryHandler.Delete, projectAuthAdmin)

}

func (s *Server) registerServerRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/servers", s.serverHandler.List, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers", s.serverHandler.Create, s.authGuard.RequireScope("server:write"))
	apiGroup.GET("/ws/servers/:serverId/metrics", s.serverMetricsWSHandler.Handle)
}

func (s *Server) registerOrganizationRoutes(authGroup *echo.Group) {
	authGroup.GET("/organizations", s.orgHandler.List)
	authGroup.POST("/organizations", s.orgHandler.Create)
	authGroup.GET("/organizations/:id", s.orgHandler.Get, s.authGuard.RequireOrgRole(models.MemberPermissionMember))
	authGroup.DELETE("/organizations/:id", s.orgHandler.Delete, s.authGuard.RequireOrgRole(models.MemberPermissionOwner))
	authGroup.GET("/organizations/:id/members", s.orgHandler.ListMembers, s.authGuard.RequireOrgRole(models.MemberPermissionMember))
	authGroup.POST("/organizations/:id/members", s.orgHandler.InviteMember, s.authGuard.RequireOrgRole(models.MemberPermissionAdmin))
	authGroup.PUT("/organizations/:id/members/:userId", s.orgHandler.UpdateMember, s.authGuard.RequireOrgRole(models.MemberPermissionAdmin))
	authGroup.DELETE("/organizations/:id/members/:memberId", s.orgHandler.RemoveMember, s.authGuard.RequireOrgRole(models.MemberPermissionAdmin))
}

func (s *Server) registerDatabaseRoutes(authGroup *echo.Group) {
	authGroup.GET("/databases", s.dbHandler.ListDatabases, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases", s.dbHandler.CreateDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/databases/:id", s.dbHandler.GetDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.PUT("/databases/:id", s.dbHandler.UpdateDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.DELETE("/databases/:id", s.dbHandler.DeleteDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/start", s.dbHandler.StartDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/stop", s.dbHandler.StopDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/restart", s.dbHandler.RestartDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/query", s.dbHandler.QueryDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/import", s.dbHandler.ImportData, s.authGuard.RequireScope("database:manage"))

	authGroup.GET("/databases/:id/schemas", s.dbHandler.GetSchemas, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/databases/:id/data/:table", s.dbHandler.GetTableData, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/data/:table", s.dbHandler.InsertTableRow, s.authGuard.RequireScope("database:manage"))
	authGroup.PUT("/databases/:id/data/:table", s.dbHandler.UpdateTableRow, s.authGuard.RequireScope("database:manage"))
	authGroup.DELETE("/databases/:id/data/:table", s.dbHandler.DeleteTableRow, s.authGuard.RequireScope("database:manage"))
}
