package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/organizations"
	"codedock/internal/modules/users"
)

func (g *AuthGuard) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyAuth(c); err != nil {
				return err
			}
			if err := g.authorizePersonalRequest(c); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyAuth(c echo.Context) error {
	userClaims, err := g.baseAuth(c, false)
	if err != nil {
		return err
	}
	c.Set("user", userClaims)
	ctx := context.WithValue(c.Request().Context(), userClaimsKey, userClaims)
	c.SetRequest(c.Request().WithContext(ctx))
	return nil
}

func (g *AuthGuard) RequireScope(requiredScope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyScope(c, requiredScope); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyScope(c echo.Context, requiredScope string) error {
	userClaims, ok := c.Get("user").(*users.UserClaims)
	if !ok || userClaims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if userClaims.Role == "api" {
		scopes, ok := c.Get("api_scopes").([]string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "insufficient scopes")
		}
		if hasMatchingScope(scopes, requiredScope) {
			return nil
		}
		return echo.NewHTTPError(http.StatusForbidden, "missing required scope: "+requiredScope)
	}
	return nil
}

func hasMatchingScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == required || s == "admin" || s == "*" {
			return true
		}
		switch required {
		case "domains:read":
			if s == "domains:write" || s == "domains:delete" || s == "domain:manage" {
				return true
			}
		case "domains:write":
			if s == "domains:delete" || s == "domain:manage" {
				return true
			}
		case "deployments:read":
			if s == "deployments:trigger" || s == "deploy:manage" {
				return true
			}
		case "deployments:trigger":
			if s == "deploy:manage" {
				return true
			}
		}
	}
	return false
}

func (g *AuthGuard) RequireProjectRole(minPermission organizations.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyProjectRole(c, minPermission); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyProjectRole(c echo.Context, minPermission organizations.MemberPermission) error {
	userClaims, ok := c.Get("user").(*users.UserClaims)
	if !ok || userClaims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if userClaims.Role == users.UserRoleAdmin || userClaims.Role == users.UserRoleOwner {
		return nil
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing project id")
	}

	if userClaims.Role == "api" {
		if c.Get("project_id") != projectID {
			return echo.NewHTTPError(http.StatusForbidden, "api token not authorized for this project")
		}
		if minPermission != "" && minPermission != organizations.MemberPermissionMember {
			return echo.NewHTTPError(http.StatusForbidden, "api tokens cannot perform administrative actions")
		}
		return nil
	}

	if g.OrgMembers == nil || g.ProjectRepo == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "project or organization members provider not configured")
	}

	project, err := g.ProjectRepo.Get(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	}

	if project.OrganizationID == "" {
		return echo.NewHTTPError(http.StatusForbidden, "project does not belong to any organization")
	}

	member, err := g.OrgMembers.GetMember(c.Request().Context(), project.OrganizationID, userClaims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify organization membership")
	}
	if member == nil || member.Status != organizations.MemberStatusAccepted {
		return echo.NewHTTPError(http.StatusForbidden, "you do not have access to this project's organization")
	}

	if minPermission != "" && member.Permission != minPermission && member.Permission != organizations.MemberPermissionAdmin && member.Permission != organizations.MemberPermissionOwner {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient organization permissions")
	}

	return nil
}

func (g *AuthGuard) RequireRole(requiredRole users.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyRole(c, requiredRole); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyRole(c echo.Context, requiredRole users.UserRole) error {
	userClaims, err := g.baseAuth(c, true)
	if err != nil {
		return err
	}
	if userClaims.Role != requiredRole && userClaims.Role != users.UserRoleAdmin && userClaims.Role != users.UserRoleOwner {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient instance permissions")
	}
	c.Set("user", userClaims)
	ctx := context.WithValue(c.Request().Context(), userClaimsKey, userClaims)
	c.SetRequest(c.Request().WithContext(ctx))
	return nil
}

func GetUserClaimsFromContext(ctx context.Context) *users.UserClaims {
	if c, ok := ctx.Value(userClaimsKey).(*users.UserClaims); ok {
		return c
	}
	return nil
}

func (g *AuthGuard) RequireOrgRole(minPermission organizations.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyOrgRole(c, minPermission); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyOrgRole(c echo.Context, minPermission organizations.MemberPermission) error {
	userClaims, ok := c.Get("user").(*users.UserClaims)
	if !ok || userClaims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if userClaims.Role == users.UserRoleAdmin || userClaims.Role == users.UserRoleOwner {
		return nil
	}

	orgID := c.Param("orgId")
	if orgID == "" {
		orgID = c.Param("id")
	}
	if orgID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing organization id")
	}

	if userClaims.Role == "api" {
		projectID, _ := c.Get("project_id").(string)
		if projectID != "" && g.ProjectRepo != nil {
			project, err := g.ProjectRepo.Get(c.Request().Context(), projectID)
			if err == nil && project != nil && project.OrganizationID == orgID {
				if minPermission != "" && minPermission != organizations.MemberPermissionMember {
					return echo.NewHTTPError(http.StatusForbidden, "api tokens cannot perform administrative organization actions")
				}
				return nil
			}
		}
		return echo.NewHTTPError(http.StatusForbidden, "api token not authorized for this organization")
	}

	if g.OrgMembers == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "organization members provider not configured")
	}

	member, err := g.OrgMembers.GetMember(c.Request().Context(), orgID, userClaims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify organization membership")
	}
	if member == nil || member.Status != organizations.MemberStatusAccepted {
		return echo.NewHTTPError(http.StatusForbidden, "you do not have access to this organization")
	}

	if minPermission != "" && member.Permission != minPermission && member.Permission != organizations.MemberPermissionAdmin && member.Permission != organizations.MemberPermissionOwner {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient organization permissions")
	}

	return nil
}
