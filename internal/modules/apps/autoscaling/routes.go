package autoscaling

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *AutoscalingHandler) {
	if h == nil {
		return
	}
	g.GET("/services/:serviceId/autoscaling", h.Get)
	g.PUT("/services/:serviceId/autoscaling", h.Update)
	g.GET("/services/:serviceId/serverless/code", h.GetServerlessCode)
	g.POST("/services/:serviceId/serverless/code", h.SaveServerlessCode)
}
