package networks

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(authGroup *echo.Group, h *ServerNetworksHandler, adminMw ...echo.MiddlewareFunc) {
	if h == nil {
		return
	}
	authGroup.GET("/system/networks", h.List)
	authGroup.GET("/system/networks/:id", h.Get)
	authGroup.POST("/system/networks", h.Create, adminMw...)
	authGroup.PUT("/system/networks/:id/access", h.ReviseAccess, adminMw...)
	authGroup.DELETE("/system/networks/:id", h.Delete, adminMw...)
	authGroup.GET("/system/networks/operations/:operationId/watchdog", h.GetWatchdog)
	authGroup.POST("/system/networks/operations/:operationId/watchdog/commit", h.CommitWatchdog, adminMw...)
	authGroup.POST("/system/networks/operations/:operationId/watchdog/rollback", h.RollbackWatchdog, adminMw...)
}
