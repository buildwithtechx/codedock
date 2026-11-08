package security

import "github.com/labstack/echo/v4"

// RegisterRoutes registers all security advisory endpoints onto an Echo group
func RegisterRoutes(authGroup *echo.Group, h *Handler) {
	authGroup.GET("/updates/advisories", h.GetAdvisories)
	authGroup.GET("/updates/:projectId/security", h.GetProjectSecurity)
	authGroup.POST("/updates/:projectId/security/check", h.CheckProject)
	authGroup.PUT("/updates/:projectId/security/policy", h.SetProjectPolicy)
	authGroup.POST("/updates/:projectId/security/apply", h.ApplyRemediation)
}
