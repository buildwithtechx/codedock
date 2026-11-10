package canvas

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *CanvasHandler) {
	if h == nil {
		return
	}
	g.GET("/canvas/projects", h.ListCanvasSummaries)
	g.GET("/projects/:id/summary", h.GetCanvasSummary)
	g.GET("/environments/:id/canvas", h.GetEnvironmentCanvas)
	g.PUT("/environments/:id/canvas", h.ApplyTopology)
	g.GET("/environments/:id/canvas/observe", h.ObserveResource)
}
