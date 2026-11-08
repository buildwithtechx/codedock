package dns

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *Handler, middleware ...echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &Handler{}
	}

	authGroup.POST("/dns", h.Create, middleware...)
	authGroup.GET("/dns", h.List, middleware...)
	authGroup.PUT("/dns/:id", h.Update, middleware...)
	authGroup.DELETE("/dns/:id", h.Delete, middleware...)
}

