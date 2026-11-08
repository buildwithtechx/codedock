package settings

import "github.com/labstack/echo/v4"

func RegisterRoutes(apiGroup, authGroup *echo.Group, handler *SettingsHandler, requireAdmin echo.MiddlewareFunc) {
	if handler == nil {
		handler = &SettingsHandler{}
	}
	authGroup.GET("/settings", handler.GetSettings)
	apiGroup.PUT("/settings", handler.UpdateSettings, requireAdmin)
}
