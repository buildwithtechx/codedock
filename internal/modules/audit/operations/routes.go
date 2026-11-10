package operations

import "github.com/labstack/echo/v4"

func withMw(mws ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	var active []echo.MiddlewareFunc
	for _, m := range mws {
		if m != nil {
			active = append(active, m)
		}
	}
	return active
}

func RegisterRoutes(authGroup *echo.Group, h *OperationHandler, readMw, writeMw echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &OperationHandler{}
	}

	authGroup.GET("/operations", h.List, withMw(readMw)...)
	authGroup.GET("/operations/:operationId", h.Get, withMw(readMw)...)
	authGroup.POST("/operations/:operationId/cancel", h.Cancel, withMw(writeMw)...)
}
