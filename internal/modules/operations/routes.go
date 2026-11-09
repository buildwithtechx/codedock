package operations

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *OperationHandler, readMw, writeMw echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &OperationHandler{}
	}

	if readMw != nil {
		authGroup.GET("/operations", h.List, readMw)
		authGroup.GET("/operations/:operationId", h.Get, readMw)
	} else {
		authGroup.GET("/operations", h.List)
		authGroup.GET("/operations/:operationId", h.Get)
	}

	if writeMw != nil {
		authGroup.POST("/operations/:operationId/cancel", h.Cancel, writeMw)
	} else {
		authGroup.POST("/operations/:operationId/cancel", h.Cancel)
	}
}
