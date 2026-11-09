package takeover

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the takeover endpoints.
func RegisterRoutes(authGroup, apiGroup *echo.Group, h *TakeoverHandler, adminMw ...echo.MiddlewareFunc) {
	if h == nil {
		return
	}
	apiGroup.POST("/system/takeover/scan", h.Scan, adminMw...)
	apiGroup.POST("/system/takeover/adopt", h.Adopt, adminMw...)
	authGroup.GET("/system/takeover/runs", h.ListRuns, adminMw...)
	authGroup.GET("/system/takeover/runs/:id", h.GetRun, adminMw...)
}
