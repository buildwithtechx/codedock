package actions

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, adminMw echo.MiddlewareFunc, h *ActionsHandler) {
	if h == nil {
		return
	}
	var admin []echo.MiddlewareFunc
	if adminMw != nil {
		admin = append(admin, adminMw)
	}

	authGroup.GET("/actions/workflows", h.ListWorkflows)
	authGroup.POST("/actions/workflows", h.CreateWorkflow, admin...)
	authGroup.POST("/actions/workflows/:id/dispatch", h.DispatchWorkflow)
	authGroup.DELETE("/actions/workflows/:id", h.DeleteWorkflow, admin...)

	authGroup.GET("/actions/runs", h.ListRuns)
	authGroup.GET("/actions/runs/:id", h.GetRun)
	authGroup.POST("/actions/runs/:id/retry", h.RetryRun)
	authGroup.POST("/actions/runs/:id/cancel", h.CancelRun)
	authGroup.DELETE("/actions/runs/:id", h.DeleteRun, admin...)

	authGroup.GET("/actions/runners", h.ListRunners)
	authGroup.POST("/actions/runners", h.CreateRunner, admin...)
	authGroup.POST("/actions/runners/inspect", h.InspectDestination, admin...)
	authGroup.DELETE("/actions/runners/:id", h.DeleteRunner, admin...)
}
