package projects

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler             *Handler
	AppHandler          *ProjectAppHandler
	SettingsHandler     *ProjectSettingsHandler
	ChecksHandler       *ChecksHandler
	CanvasHandler       *CanvasHandler
	ComposeStackHandler *ComposeStackHandler
	ComposeHandler      *ComposeHandler
	OneClickHandler     *OneClickHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler != nil {
		g.GET("/projects", cfg.Handler.ListProjects)
		g.POST("/projects", cfg.Handler.CreateProject)
		g.GET("/projects/:id", cfg.Handler.GetProject)
		g.DELETE("/projects/:id", cfg.Handler.DeleteProject)
	}

	if cfg.AppHandler != nil {
		g.GET("/project-apps", cfg.AppHandler.List)
		g.POST("/project-apps", cfg.AppHandler.Create)
		g.GET("/project-apps/:id", cfg.AppHandler.Get)
		g.PUT("/project-apps/:id", cfg.AppHandler.Update)
		g.DELETE("/project-apps/:id", cfg.AppHandler.Delete)
	}

	if cfg.SettingsHandler != nil {
		g.GET("/projects/:projectId/tokens", cfg.SettingsHandler.ListTokens)
		g.POST("/projects/:projectId/tokens", cfg.SettingsHandler.CreateToken)
		g.DELETE("/projects/:projectId/tokens/:id", cfg.SettingsHandler.DeleteToken)
		g.GET("/projects/:id/members", cfg.SettingsHandler.ListMembers)
		g.POST("/projects/:id/members", cfg.SettingsHandler.AddMember)
		g.DELETE("/projects/:id/members/:memberId", cfg.SettingsHandler.RemoveMember)
		g.GET("/projects/:id/rollback-capacity", cfg.SettingsHandler.GetRollbackCapacity)
		g.PATCH("/projects/:id/rollback-settings", cfg.SettingsHandler.UpdateRollbackSettings)
	}

	if cfg.ChecksHandler != nil {
		g.GET("/projects/:id/deployments/admission", cfg.ChecksHandler.GetAdmission)
		g.GET("/projects/:id/checks", cfg.ChecksHandler.GetChecks)
		g.PUT("/projects/:id/checks", cfg.ChecksHandler.UpdateChecks)
	}

	if cfg.CanvasHandler != nil {
		g.GET("/canvas/projects", cfg.CanvasHandler.ListCanvasSummaries)
		g.GET("/projects/:id/summary", cfg.CanvasHandler.GetCanvasSummary)
		g.GET("/environments/:id/canvas", cfg.CanvasHandler.GetEnvironmentCanvas)
		g.PUT("/environments/:id/canvas", cfg.CanvasHandler.ApplyTopology)
		g.GET("/environments/:id/canvas/observe", cfg.CanvasHandler.ObserveResource)
	}

	if cfg.ComposeStackHandler != nil {
		g.GET("/projects/:id/stacks", cfg.ComposeStackHandler.List)
		g.GET("/projects/:id/stacks/:stackId/config", cfg.ComposeStackHandler.Config)
		g.POST("/projects/:id/stacks/review", cfg.ComposeStackHandler.Review)
		g.POST("/projects/:id/stacks", cfg.ComposeStackHandler.Save)
		g.POST("/projects/:id/stacks/:stackId/deploy", cfg.ComposeStackHandler.Deploy)
		g.POST("/projects/:id/stacks/:stackId/cancel", cfg.ComposeStackHandler.Cancel)
		g.POST("/projects/:id/stacks/:stackId/restart", cfg.ComposeStackHandler.Restart)
		g.POST("/projects/:id/stacks/:stackId/stop", cfg.ComposeStackHandler.Stop)
	}

	if cfg.ComposeHandler != nil {
		g.POST("/compose/deploy", cfg.ComposeHandler.Deploy)
		g.POST("/compose/analyze", cfg.ComposeHandler.Analyze)
	}

	if cfg.OneClickHandler != nil {
		g.GET("/one-click", cfg.OneClickHandler.List)
		g.GET("/one-click/:id", cfg.OneClickHandler.Get)
		g.POST("/one-click/review", cfg.OneClickHandler.Review)
		g.POST("/one-click/deploy", cfg.OneClickHandler.Deploy)
		g.POST("/one-click/custom", cfg.OneClickHandler.AddCustom)
		g.DELETE("/one-click/custom/:id", cfg.OneClickHandler.DeleteCustom)
	}
}
