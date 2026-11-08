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

		if mw.ProjectAuth != nil {
			authGroup.GET("/projects/:id/deployments", depHandler.ListProjectDeployments, mw.ProjectAuth)
		} else {
			authGroup.GET("/projects/:id/deployments", depHandler.ListProjectDeployments)
		}

		if mw.ServiceAuth != nil {
			authGroup.GET("/services/:serviceId/deployments", depHandler.ListServiceDeployments, mw.ServiceAuth)
			authGroup.GET("/services/:serviceId/previews", depHandler.ListPRPreviews, mw.ServiceAuth)
			authGroup.GET("/services/:serviceId/metrics", depHandler.GetMetrics, mw.ServiceAuth)
		} else {
			authGroup.GET("/services/:serviceId/deployments", depHandler.ListServiceDeployments)
			authGroup.GET("/services/:serviceId/previews", depHandler.ListPRPreviews)
			authGroup.GET("/services/:serviceId/metrics", depHandler.GetMetrics)
		}

		if mw.RequireScope != nil {
			authGroup.GET("/deployments/:id/logs", depHandler.GetLogs, mw.RequireScope("logs:read"))
		} else {
			authGroup.GET("/deployments/:id/logs", depHandler.GetLogs)
		}
		authGroup.GET("/deployments/:id/explain", depHandler.ExplainFailure)
	}

	// Lifecycle operations
	if authGroup != nil && lifecycleHandler != nil {
		authGroup.GET("/projects/:id/deployments/admission", lifecycleHandler.GetAdmission)

		if mw.ServiceAdmin != nil {
			authGroup.POST("/services/:serviceId/deploy", lifecycleHandler.Trigger, mw.ServiceAdmin)
		} else {
			authGroup.POST("/services/:serviceId/deploy", lifecycleHandler.Trigger)
		}

		if mw.ProjectAdmin != nil {
			authGroup.POST("/projects/:id/deploy", lifecycleHandler.TriggerProject, mw.ProjectAdmin)
		} else {
			authGroup.POST("/projects/:id/deploy", lifecycleHandler.TriggerProject)
		}

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
		if mw.ProjectAdmin != nil {
			authGroup.GET("/projects/:projectId/registries", registryHandler.List, mw.ProjectAdmin)
			authGroup.POST("/projects/:projectId/registries", registryHandler.Create, mw.ProjectAdmin)
			authGroup.DELETE("/projects/:projectId/registries/:id", registryHandler.Delete, mw.ProjectAdmin)
		} else {
			authGroup.GET("/projects/:projectId/registries", registryHandler.List)
			authGroup.POST("/projects/:projectId/registries", registryHandler.Create)
			authGroup.DELETE("/projects/:projectId/registries/:id", registryHandler.Delete)
		}
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
