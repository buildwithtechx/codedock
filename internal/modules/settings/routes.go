package settings

import "github.com/labstack/echo/v4"

type AppearanceHandler interface {
	GetAppearance(c echo.Context) error
	SetAppearance(c echo.Context) error
}

func RegisterRoutes(apiGroup, authGroup *echo.Group, handler *SettingsHandler, requireAdmin echo.MiddlewareFunc, appearance ...AppearanceHandler) {
	if handler == nil {
		handler = &SettingsHandler{}
	}
	if authGroup != nil {
		authGroup.GET("/settings", handler.GetSettings)
		for _, app := range appearance {
			if app != nil {
				authGroup.GET("/settings/appearance", app.GetAppearance)
				authGroup.PUT("/settings/appearance", app.SetAppearance)
			}
		}
	}
	if apiGroup != nil {
		apiGroup.GET("/system/public", handler.GetPublicSettings)
		apiGroup.PUT("/settings", handler.UpdateSettings, requireAdmin)
	}
}
