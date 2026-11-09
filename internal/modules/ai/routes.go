package ai

import (
	"github.com/labstack/echo/v4"
)

func withMw(mws ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	var active []echo.MiddlewareFunc
	for _, m := range mws {
		if m != nil {
			active = append(active, m)
		}
	}
	return active
}

func RegisterRoutes(authGroup *echo.Group, apiGroup *echo.Group, h *Handler, adminMiddleware echo.MiddlewareFunc, diagnoseMiddleware ...echo.MiddlewareFunc) {
	if h == nil {
		h = &Handler{}
	}

	if authGroup != nil {
		authGroup.GET("/ai", h.GetAISettings)
		authGroup.POST("/ai/diagnose", h.DiagnoseLogs, withMw(diagnoseMiddleware...)...)
	}

	if apiGroup != nil {
		apiGroup.PUT("/ai", h.UpdateAISettings, withMw(adminMiddleware)...)
	}
}

type AssistantRoutesConfig struct {
	ListResources echo.HandlerFunc
	Chat          echo.HandlerFunc
	ConfirmAction echo.HandlerFunc
}

func RegisterAssistantRoutes(groups []*echo.Group, cfg AssistantRoutesConfig) {
	for _, grp := range groups {
		if grp == nil {
			continue
		}
		if cfg.ListResources != nil {
			grp.GET("/ai/resources", cfg.ListResources)
		}
		if cfg.Chat != nil {
			grp.POST("/ai/chat", cfg.Chat)
		}
		if cfg.ConfirmAction != nil {
			grp.POST("/ai/actions/confirm", cfg.ConfirmAction)
		}
	}
}
