package git

import (
	"github.com/labstack/echo/v4"
)

type RouteMiddleware struct {
	Admin              echo.MiddlewareFunc
	RequireProjectRole func(role string) echo.MiddlewareFunc
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
		if mw.RequireProjectRole != nil {
			authGroup.POST("/projects/:id/repository-inspection", inspectionHandler.InspectRepository, mw.RequireProjectRole("admin"))
		} else {
			authGroup.POST("/projects/:id/repository-inspection", inspectionHandler.InspectRepository)
		}
	}

	// GitHub Apps
	if appsHandler != nil {
		if apiGroup != nil {
			if mw.Admin != nil {
				apiGroup.POST("/settings/git_apps/github/manifest-callback", appsHandler.ExchangeGithubManifestCode, mw.Admin)
				apiGroup.PUT("/settings/git_apps/github", appsHandler.SaveGithubApp, mw.Admin)
				apiGroup.DELETE("/settings/git_apps/github/:id", appsHandler.DeleteGithubApp, mw.Admin)
			} else {
				apiGroup.POST("/settings/git_apps/github/manifest-callback", appsHandler.ExchangeGithubManifestCode)
				apiGroup.PUT("/settings/git_apps/github", appsHandler.SaveGithubApp)
				apiGroup.DELETE("/settings/git_apps/github/:id", appsHandler.DeleteGithubApp)
			}
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
