package app

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/adapters"
	cdkhttp "codedock/internal/http"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/settings"
)

func (a *App) wireRoutes(s *appServices) {
	apiGroup := a.router.Group("/api")
	authGroup := apiGroup.Group("")
	authGroup.Use(a.authGuard.RequireAuth())

	a.registerAuthRoutes(apiGroup, authGroup, s)
	a.registerSystemRoutes(apiGroup, authGroup, s)
	a.registerUserRoutes(apiGroup, authGroup, s)
	a.registerProjectRoutes(apiGroup, authGroup, s)
	a.registerOrganizationRoutes(authGroup, s)
	a.registerServerRoutes(apiGroup, authGroup, s)
	a.registerMigrationRoutes(apiGroup, authGroup, s)
	a.registerAnalyticsRoutes(apiGroup, authGroup, s)
	a.registerDatabaseRoutes(authGroup, s)
	a.registerAppRoutes(apiGroup, authGroup, s)
	a.registerDeploymentRoutes(apiGroup, authGroup, s)
	a.registerBackupRoutes(authGroup, s)
	a.registerClusterRoutes(authGroup, s)
	a.registerSettingsRoutes(apiGroup, authGroup, s)
	a.registerMiscRoutes(apiGroup, authGroup, s)
	a.registerBillingRoutes(apiGroup, authGroup, s)
	a.registerMailRoutes(authGroup, s)

	// Health and environment endpoints
	settingsHandler := settings.NewSettingsHandler(s.settingsService)
	a.router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	a.router.GET("/health/env", settingsHandler.GetHealthEnv)
	apiGroup.GET("/health/env", settingsHandler.GetHealthEnv)

	// Edge Authentication for routes
	protectionService := environments.NewRouteProtectionService(s.routeProtectionRepo, &adapters.ProjectLookupAdapter{ProjectService: a.projectService})
	edgeAuthHandler := environments.NewEdgeAuthHandler(protectionService)
	environments.RegisterEdgeAuthRoutes(a.router, edgeAuthHandler)

	// Dashboard Single Page Application fallback
	cdkhttp.SetupSPAFallback(a.router)
}
