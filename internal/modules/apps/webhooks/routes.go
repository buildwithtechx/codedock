package webhooks

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *WebhookHandler) {
	if h == nil {
		return
	}
	g.GET("/apps/:id/webhooks", h.ListWebhooks)
	g.POST("/apps/:id/webhooks", h.CreateWebhook)
	g.DELETE("/apps/:id/webhooks/:webhookId", h.DeleteWebhook)

	g.GET("/apps/:id/log-drains", h.ListLogDrains)
	g.POST("/apps/:id/log-drains", h.CreateLogDrain)
	g.DELETE("/apps/:id/log-drains/:drainId", h.DeleteLogDrain)
}
