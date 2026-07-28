package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func (g *AuthGuard) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyAuth(c); err != nil {
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
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if userClaims.Role == "api" {
		scopes, ok := c.Get("api_scopes").([]string)
		if !ok {
			return utils.Error(c, http.StatusForbidden, "insufficient scopes")
		}
		hasScope := false
		for _, s := range scopes {
			if s == requiredScope || s == "admin" || s == "*" {
				hasScope = true
				break
			}
		}
		if !hasScope {
			return utils.Error(c, http.StatusForbidden, "missing required scope: "+requiredScope)
		}
	}
	return nil
}

func (g *AuthGuard) RequireProjectRole(minPermission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyProjectRole(c, minPermission); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyProjectRole(c echo.Context, minPermission models.MemberPermission) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return nil
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}

	if userClaims.Role == "api" {
		if c.Get("project_id") != projectID {
			return utils.Error(c, http.StatusForbidden, "api token not authorized for this project")
		}
		if minPermission != "" && minPermission != models.MemberPermissionMember {
			return utils.Error(c, http.StatusForbidden, "api tokens cannot perform administrative actions")
		}
		return nil
	}

	if g.OrgMembers == nil || g.ProjectRepo == nil {
		return utils.Error(c, http.StatusInternalServerError, "project or organization members provider not configured")
	}

	project, err := g.ProjectRepo.Get(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}

	if project.OrganizationID == "" {
		return utils.Error(c, http.StatusForbidden, "project does not belong to any organization")
	}

	member, err := g.OrgMembers.GetMember(c.Request().Context(), project.OrganizationID, userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to verify organization membership")
	}
	if member == nil || member.Status != models.MemberStatusAccepted {
		return utils.Error(c, http.StatusForbidden, "you do not have access to this project's organization")
	}

	if minPermission != "" && member.Permission != minPermission && member.Permission != models.MemberPermissionAdmin && member.Permission != models.MemberPermissionOwner {
		return utils.Error(c, http.StatusForbidden, "insufficient organization permissions")
	}

	return nil
}

func (g *AuthGuard) RequireRole(requiredRole models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyRole(c, requiredRole); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyRole(c echo.Context, requiredRole models.UserRole) error {
	userClaims, err := g.baseAuth(c, true)
	if err != nil {
		return err
	}
	if userClaims.Role != requiredRole && userClaims.Role != models.UserRoleAdmin && userClaims.Role != models.UserRoleOwner {
		return utils.Error(c, http.StatusForbidden, "insufficient instance permissions")
	}
	c.Set("user", userClaims)
	ctx := context.WithValue(c.Request().Context(), userClaimsKey, userClaims)
	c.SetRequest(c.Request().WithContext(ctx))
	return nil
}

func GetUserClaimsFromContext(ctx context.Context) *models.UserClaims {
	if c, ok := ctx.Value(userClaimsKey).(*models.UserClaims); ok {
		return c
	}
	return nil
}

func (g *AuthGuard) RequireOrgRole(minPermission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.verifyOrgRole(c, minPermission); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) verifyOrgRole(c echo.Context, minPermission models.MemberPermission) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	if userClaims.Role == "admin" {
		return nil
	}

	orgID := c.Param("orgId")
	if orgID == "" {
		orgID = c.Param("id")
	}
	if orgID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing organization id")
	}

	if userClaims.Role == "api" {
		projectID, _ := c.Get("project_id").(string)
		if projectID != "" && g.ProjectRepo != nil {
			project, err := g.ProjectRepo.Get(c.Request().Context(), projectID)
			if err == nil && project != nil && project.OrganizationID == orgID {
				if minPermission != "" && minPermission != models.MemberPermissionMember {
					return utils.Error(c, http.StatusForbidden, "api tokens cannot perform administrative organization actions")
				}
				return nil
			}
		}
		return utils.Error(c, http.StatusForbidden, "api token not authorized for this organization")
	}

	if g.OrgMembers == nil {
		return utils.Error(c, http.StatusInternalServerError, "organization members provider not configured")
	}

	member, err := g.OrgMembers.GetMember(c.Request().Context(), orgID, userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to verify organization membership")
	}
	if member == nil || member.Status != models.MemberStatusAccepted {
		return utils.Error(c, http.StatusForbidden, "you do not have access to this organization")
	}

	if minPermission != "" && member.Permission != minPermission && member.Permission != models.MemberPermissionAdmin && member.Permission != models.MemberPermissionOwner {
		return utils.Error(c, http.StatusForbidden, "insufficient organization permissions")
	}

	return nil
}
