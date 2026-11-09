package vars

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *VariableHandler) {
	if h == nil {
		return
	}
	g.GET("/services/:serviceId/variables", h.List)
	g.POST("/services/:serviceId/variables", h.Create)
	g.PUT("/services/:serviceId/variables/:id", h.Update)
	g.DELETE("/services/:serviceId/variables/:id", h.Delete)
	g.GET("/services/:serviceId/env-suggestions", h.Suggest)
}
