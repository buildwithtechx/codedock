package notifications

import "github.com/labstack/echo/v4"

func RegisterRoutes(apiGroup, authGroup *echo.Group, notifH *NotificationHandler, settingsH *NotificationSettingsHandler, adminMw ...echo.MiddlewareFunc) {
	if settingsH == nil {
		settingsH = &NotificationSettingsHandler{}
	}
	if notifH == nil {
		notifH = &NotificationHandler{}
	}

	if authGroup != nil {
		authGroup.GET("/notifications", settingsH.GetNotificationSettings)
		authGroup.GET("/notifications/subscriptions", settingsH.ListSubscriptions)
		authGroup.PUT("/notifications/subscriptions", settingsH.UpsertSubscription)
		authGroup.GET("/notifications/defaults", settingsH.ListDefaults)
		authGroup.GET("/notifications/deliveries", settingsH.ListDeliveries)
		authGroup.GET("/notifications/devices", settingsH.ListDeviceTokens)
		authGroup.POST("/notifications/devices", settingsH.RegisterDeviceToken)
		authGroup.DELETE("/notifications/devices/:token", settingsH.DeleteDeviceToken)
	}

	if apiGroup != nil {
		apiGroup.PUT("/notifications", settingsH.UpdateNotificationSettings, adminMw...)
		apiGroup.PUT("/notifications/defaults", settingsH.UpsertDefault, adminMw...)
		apiGroup.POST("/settings/notifications/test", notifH.TestNotification, adminMw...)
	}
}
