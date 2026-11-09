package volumes

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *VolumeHandler) {
	if h == nil {
		return
	}
	g.GET("/apps/:id/volumes", h.List)
	g.POST("/apps/:id/volumes", h.Create)
	g.DELETE("/apps/:id/volumes/:volumeId", h.Delete)
	g.POST("/apps/:id/volumes/:volumeId/empty", h.EmptyVolume)
}
