package git

import (
	"github.com/labstack/echo/v4"
)

type RouteMiddleware struct {
	Admin              echo.MiddlewareFunc
	RequireProjectRole func(role string) echo.MiddlewareFunc
}

func withMw(mws ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	var active []echo.MiddlewareFunc
	for _, m := range mws {
		if m != nil {
			active = append(active, m)
		}
	}
	return active
}

func RegisterRoutes(
	apiGroup *echo.Group,
	authGroup *echo.Group,
	providersHandler *ProvidersHandler,
	inspectionHandler *InspectionHandler,
	appsHandler *AppsHandler,
	webhookHandler *WebhookHandler,
	mw RouteMiddleware,
) {
	// Git provider connections
	if authGroup != nil && providersHandler != nil {
		authGroup.POST("/git/connect", providersHandler.Connect)
		authGroup.GET("/git/status", providersHandler.Status)
		authGroup.DELETE("/git/connect/:provider", providersHandler.Disconnect)
		authGroup.GET("/git/repos", providersHandler.ListRepos)
	}

	// Repository inspection & tree
	if authGroup != nil && inspectionHandler != nil {
		authGroup.POST("/git/tree", inspectionHandler.RepoTree)
		authGroup.GET("/git/tree", inspectionHandler.RepoTree)
		authGroup.POST("/git/inspect", inspectionHandler.InspectRepository)
		var inspMw []echo.MiddlewareFunc
		if mw.RequireProjectRole != nil {
			inspMw = append(inspMw, mw.RequireProjectRole("admin"))
		}
		authGroup.POST("/projects/:id/repository-inspection", inspectionHandler.InspectRepository, inspMw...)
	}

	// GitHub Apps
	if appsHandler != nil {
		if apiGroup != nil {
			admin := withMw(mw.Admin)
			apiGroup.POST("/settings/git_apps/github/manifest-callback", appsHandler.ExchangeGithubManifestCode, admin...)
			apiGroup.PUT("/settings/git_apps/github", appsHandler.SaveGithubApp, admin...)
			apiGroup.DELETE("/settings/git_apps/github/:id", appsHandler.DeleteGithubApp, admin...)
		}
		if authGroup != nil {
			authGroup.GET("/settings/git_apps/github", appsHandler.ListGithubApps)
			authGroup.GET("/settings/git_apps/github/:id", appsHandler.GetGithubApp)
		}
	}

	// Webhooks
	if webhookHandler != nil {
		if apiGroup != nil {
			apiGroup.POST("/webhooks/git/services/:serviceId", webhookHandler.HandleServiceWebhook)
			apiGroup.POST("/webhooks/github/services/:serviceId", webhookHandler.HandleGitHubWebhook)
			apiGroup.POST("/webhooks/github", webhookHandler.HandleGitHubWebhook)
		}
		if authGroup != nil {
			authGroup.GET("/projects/:id/webhooks/deliveries", webhookHandler.ListProjectDeliveries)
		}
	}
}
