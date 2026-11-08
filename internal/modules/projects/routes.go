package projects

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler         *Handler
	AppHandler      *ProjectAppHandler
	SettingsHandler *ProjectSettingsHandler
	ChecksHandler   *ChecksHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler != nil {
		g.GET("/projects", cfg.Handler.ListProjects)
		g.POST("/projects", cfg.Handler.CreateProject)
		g.GET("/projects/:id", cfg.Handler.GetProject)
		g.DELETE("/projects/:id", cfg.Handler.DeleteProject)

		g.GET("/projects/:id/members", cfg.Handler.ListMembers)
		g.POST("/projects/:id/members", cfg.Handler.AddMember)
		g.DELETE("/projects/:id/members/:memberId", cfg.Handler.RemoveMember)
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
	}

	if cfg.ChecksHandler != nil {
		g.GET("/projects/:id/deployments/admission", cfg.ChecksHandler.GetAdmission)
		g.GET("/projects/:id/checks", cfg.ChecksHandler.GetChecks)
		g.PUT("/projects/:id/checks", cfg.ChecksHandler.UpdateChecks)
	}
}
