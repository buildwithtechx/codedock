package http

import (
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) RequireProjectRole(permission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := middleware.GetUserClaimsFromContext(c.Request().Context())
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			project := c.Param("id")
			if selected, ok := c.Get("project_id").(string); ok && selected != "" && selected != project {
				return echo.NewHTTPError(http.StatusForbidden, "token does not have access to this project")
			}
			if !s.projectService.HasPermission(c.Request().Context(), project, user.UserID, models.UserRole(user.Role), permission) {
				return echo.NewHTTPError(http.StatusForbidden, "project permission required")
			}
			return next(c)
		}
	}
}
