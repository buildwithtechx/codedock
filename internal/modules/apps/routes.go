package apps

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler            *Handler
	VariableHandler    *VariableHandler
	VolumeHandler      *VolumeHandler
	WebhookHandler     *WebhookHandler
	AutoscalingHandler *AutoscalingHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler != nil {
		g.GET("/apps", cfg.Handler.ListAll)
		g.GET("/projects/:id/apps", cfg.Handler.ListByProject)
		g.GET("/environments/:id/apps", cfg.Handler.ListByEnvironment)
		g.POST("/environments/:id/apps", cfg.Handler.Create)
		g.GET("/apps/:id", cfg.Handler.Get)
		g.PUT("/apps/:id", cfg.Handler.Update)
		g.PATCH("/apps/:id", cfg.Handler.Update)
		g.DELETE("/apps/:id", cfg.Handler.Delete)
	}

	if cfg.VariableHandler != nil {
		g.GET("/apps/:id/vars", cfg.VariableHandler.List)
		g.POST("/apps/:id/vars", cfg.VariableHandler.Create)
		g.PUT("/apps/:id/vars/:varId", cfg.VariableHandler.Update)
		g.PATCH("/apps/:id/vars/:varId", cfg.VariableHandler.Update)
		g.DELETE("/apps/:id/vars/:varId", cfg.VariableHandler.Delete)
	}

	if cfg.VolumeHandler != nil {
		g.GET("/apps/:id/volumes", cfg.VolumeHandler.List)
		g.POST("/apps/:id/volumes", cfg.VolumeHandler.Create)
		g.DELETE("/apps/:id/volumes/:volumeId", cfg.VolumeHandler.Delete)
		g.POST("/apps/:id/volumes/:volumeId/empty", cfg.VolumeHandler.EmptyVolume)
	}

	if cfg.WebhookHandler != nil {
		g.GET("/apps/:id/webhooks", cfg.WebhookHandler.ListWebhooks)
		g.POST("/apps/:id/webhooks", cfg.WebhookHandler.CreateWebhook)
		g.DELETE("/apps/:id/webhooks/:webhookId", cfg.WebhookHandler.DeleteWebhook)

		g.GET("/apps/:id/drains", cfg.WebhookHandler.ListLogDrains)
		g.POST("/apps/:id/drains", cfg.WebhookHandler.CreateLogDrain)
		g.DELETE("/apps/:id/drains/:drainId", cfg.WebhookHandler.DeleteLogDrain)
	}

	if cfg.AutoscalingHandler != nil {
		g.GET("/apps/:id/autoscaling", cfg.AutoscalingHandler.Get)
		g.PUT("/apps/:id/autoscaling", cfg.AutoscalingHandler.Update)
		g.GET("/apps/:id/serverless", cfg.AutoscalingHandler.GetServerlessCode)
		g.POST("/apps/:id/serverless", cfg.AutoscalingHandler.SaveServerlessCode)
	}
}
