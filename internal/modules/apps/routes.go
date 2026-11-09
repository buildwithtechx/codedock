package apps

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/modules/apps/autoscaling"
	"codedock/internal/modules/apps/vars"
	"codedock/internal/modules/apps/volumes"
)

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
		g.GET("/projects/:id/services", cfg.Handler.ListByProject)
		g.GET("/environments/:id/apps", cfg.Handler.ListByEnvironment)
		g.POST("/environments/:id/apps", cfg.Handler.Create)
		g.GET("/apps/:id", cfg.Handler.Get)
		g.PUT("/apps/:id", cfg.Handler.Update)
		g.DELETE("/apps/:id", cfg.Handler.Delete)

		// Lifecycle & Build
		g.POST("/apps/:id/stop", cfg.Handler.StopService)
		g.POST("/apps/:id/restart", cfg.Handler.RestartService)
		g.POST("/apps/:id/redeploy", cfg.Handler.RedeployService)
		g.PUT("/apps/:id/build", cfg.Handler.UpdateBuildSettings)

		// Runtime updates
		g.GET("/projects/:id/runtime-updates", cfg.Handler.CheckProjectRuntimeUpdates)
		g.POST("/projects/:id/runtime-updates", cfg.Handler.ApplyProjectRuntimeUpdates)
		g.GET("/apps/:id/runtime-updates", cfg.Handler.CheckAppRuntimeUpdates)
		g.POST("/apps/:id/runtime-updates", cfg.Handler.ApplyAppRuntimeUpdates)

		// Runtime operations & observation
		g.GET("/apps/:id/runtime", cfg.Handler.GetRuntime)
		g.GET("/apps/:id/runtime/operation", cfg.Handler.RuntimeOperation)
		g.POST("/apps/:id/runtime/review", cfg.Handler.ReviewRuntime)
		g.POST("/apps/:id/runtime/apply", cfg.Handler.ApplyRuntime)
		g.GET("/apps/:id/runtime/observe", cfg.Handler.ObserveRuntime)
		g.GET("/apps/:id/runtime/logs", cfg.Handler.RuntimeLogs)
		g.GET("/apps/:id/runtime/logs/stream", cfg.Handler.RuntimeLogStream)
		g.GET("/apps/:id/runtime/graph", cfg.Handler.InstanceGraph)
		g.GET("/apps/:id/runtime/resources", cfg.Handler.ResourceSummary)
		g.GET("/apps/:id/runtime/exec-terminal", cfg.Handler.RuntimeExecTerminal)
		g.POST("/apps/:id/runtime/exec", cfg.Handler.RuntimeExec)
	}

	if cfg.VariableHandler != nil {
		vars.RegisterRoutes(g, cfg.VariableHandler)
	}

	if cfg.VolumeHandler != nil {
		volumes.RegisterRoutes(g, cfg.VolumeHandler)
	}

	if cfg.WebhookHandler != nil {
		g.GET("/apps/:id/webhooks", cfg.WebhookHandler.ListWebhooks)
		g.POST("/apps/:id/webhooks", cfg.WebhookHandler.CreateWebhook)
		g.DELETE("/apps/:id/webhooks/:webhookId", cfg.WebhookHandler.DeleteWebhook)

		g.GET("/apps/:id/log-drains", cfg.WebhookHandler.ListLogDrains)
		g.POST("/apps/:id/log-drains", cfg.WebhookHandler.CreateLogDrain)
		g.DELETE("/apps/:id/log-drains/:drainId", cfg.WebhookHandler.DeleteLogDrain)
	}

	if cfg.AutoscalingHandler != nil {
		autoscaling.RegisterRoutes(g, cfg.AutoscalingHandler)
	}
}
