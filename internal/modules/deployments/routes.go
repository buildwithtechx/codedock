package deployments

import (
	"github.com/labstack/echo/v4"
)

type RouteMiddleware struct {
	Admin              echo.MiddlewareFunc
	ServiceAuth        echo.MiddlewareFunc
	ServiceAdmin       echo.MiddlewareFunc
	ProjectAuth        echo.MiddlewareFunc
	ProjectAdmin       echo.MiddlewareFunc
	RequireProjectRole func(role string) echo.MiddlewareFunc
	RequireScope       func(scope string) echo.MiddlewareFunc
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
	depHandler *DeploymentHandler,
	lifecycleHandler *LifecycleHandler,
	archiveHandler *ArchiveHandler,
	registryHandler *RegistryHandler,
	terminalHandler *TerminalHandler,
	mw RouteMiddleware,
) {
	// Deployments CRUD & inspection
	if authGroup != nil && depHandler != nil {
		authGroup.GET("/deployments", depHandler.ListOrganizationDeployments)
		authGroup.GET("/deployments/:id", depHandler.Get)
		authGroup.GET("/projects/:id/deployments", depHandler.ListProjectDeployments, withMw(mw.ProjectAuth)...)
		authGroup.GET("/services/:serviceId/deployments", depHandler.ListServiceDeployments, withMw(mw.ServiceAuth)...)
		authGroup.GET("/services/:serviceId/previews", depHandler.ListPRPreviews, withMw(mw.ServiceAuth)...)
		authGroup.GET("/services/:serviceId/metrics", depHandler.GetMetrics, withMw(mw.ServiceAuth)...)

		var logsMw []echo.MiddlewareFunc
		if mw.RequireScope != nil {
			logsMw = append(logsMw, mw.RequireScope("logs:read"))
		}
		authGroup.GET("/deployments/:id/logs", depHandler.GetLogs, logsMw...)
		authGroup.GET("/deployments/:id/explain", depHandler.ExplainFailure)
	}

	// Lifecycle operations
	if authGroup != nil && lifecycleHandler != nil {
		authGroup.GET("/projects/:id/deployments/admission", lifecycleHandler.GetAdmission)
		authGroup.POST("/services/:serviceId/deploy", lifecycleHandler.Trigger, withMw(mw.ServiceAdmin)...)
		authGroup.POST("/projects/:id/deploy", lifecycleHandler.TriggerProject, withMw(mw.ProjectAdmin)...)
		authGroup.POST("/deployments/:id/rollback", lifecycleHandler.Rollback)
		authGroup.POST("/deployments/:id/retry", lifecycleHandler.Retry)
		authGroup.POST("/deployments/:id/pin", lifecycleHandler.TogglePin)
		authGroup.POST("/deployments/:id/cancel", lifecycleHandler.Cancel)
	}

	// Archive deploy
	if authGroup != nil && archiveHandler != nil {
		authGroup.POST("/deploy/archive", archiveHandler.DeployArchive)
	}

	// Registries
	if authGroup != nil && registryHandler != nil {
		authGroup.GET("/projects/:projectId/registries", registryHandler.List, withMw(mw.ProjectAdmin)...)
		authGroup.POST("/projects/:projectId/registries", registryHandler.Create, withMw(mw.ProjectAdmin)...)
		authGroup.DELETE("/projects/:projectId/registries/:id", registryHandler.Delete, withMw(mw.ProjectAdmin)...)
	}

	// Terminal WebSockets
	if terminalHandler != nil {
		targetGroup := apiGroup
		if targetGroup == nil {
			targetGroup = authGroup
		}
		if targetGroup != nil {
			targetGroup.GET("/ws/terminal/:id", terminalHandler.HandleWebSocket)
			targetGroup.GET("/ws/services/:id/terminal", terminalHandler.HandleWebSocket)
		}
	}
}
