package ai

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(authGroup *echo.Group, apiGroup *echo.Group, h *Handler, adminMiddleware echo.MiddlewareFunc, diagnoseMiddleware ...echo.MiddlewareFunc) {
	if h == nil {
		h = &Handler{}
	}

	if authGroup != nil {
		authGroup.GET("/ai", h.GetAISettings)
		if len(diagnoseMiddleware) > 0 {
			authGroup.POST("/ai/diagnose", h.DiagnoseLogs, diagnoseMiddleware...)
		} else {
			authGroup.POST("/ai/diagnose", h.DiagnoseLogs)
		}
	}

	if apiGroup != nil {
		if adminMiddleware != nil {
			apiGroup.PUT("/ai", h.UpdateAISettings, adminMiddleware)
		} else {
			apiGroup.PUT("/ai", h.UpdateAISettings)
		}
	}
}
