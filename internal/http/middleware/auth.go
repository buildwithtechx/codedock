package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/auth"
	"codedock.run/codedock/internal/utils"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type SettingsProvider interface {
	GetSettings(context.Context) (*models.ServerSettings, error)
}

type ProjectTokenProvider interface {
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
}

type OrganizationMemberProvider interface {
	GetMember(ctx context.Context, orgID, userID string) (*models.OrganizationMember, error)
}

type AuthGuard struct {
	TokenService  *auth.TokenService
	Settings      SettingsProvider
	ProjectTokens ProjectTokenProvider
	OrgMembers    OrganizationMemberProvider
	ProjectRepo   repositories.ProjectRepository
}

func NewAuthGuard(ts *auth.TokenService, sp SettingsProvider, pt ProjectTokenProvider, orgMembers OrganizationMemberProvider, pr repositories.ProjectRepository) *AuthGuard {
	return &AuthGuard{TokenService: ts, Settings: sp, ProjectTokens: pt, OrgMembers: orgMembers, ProjectRepo: pr}
}

func (g *AuthGuard) checkIPAllowlist(c echo.Context) error {
	if g.Settings == nil {
		return nil
	}
	settings, _ := g.Settings.GetSettings(c.Request().Context())
	if settings == nil || strings.TrimSpace(settings.IPAllowlist) == "" {
		return nil
	}
	clientIP := c.RealIP()
	if !IsIPAllowed(clientIP, settings.IPAllowlist) {
		return utils.Error(c, http.StatusForbidden, fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP))
	}
	return nil
}

func (g *AuthGuard) validateAPIToken(c echo.Context, tokenStr string, denyAPITokens bool) (*models.UserClaims, error) {
	if denyAPITokens {
		return nil, utils.Error(c, http.StatusForbidden, "API tokens cannot access role-restricted endpoints")
	}
	if g.ProjectTokens == nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "API tokens not supported")
	}
	pt, err := g.ProjectTokens.GetTokenByHash(c.Request().Context(), tokenStr)
	if err != nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "invalid or revoked API token")
	}
	if pt.ExpiresAt != nil && pt.ExpiresAt.Before(time.Now()) {
		return nil, utils.Error(c, http.StatusUnauthorized, "API token has expired")
	}
	if len(pt.IPAllowlist) > 0 {
		if !IsIPAllowed(c.RealIP(), strings.Join(pt.IPAllowlist, ",")) {
			return nil, utils.Error(c, http.StatusForbidden, "IP address not allowed for this API token")
		}
	}
	_ = g.ProjectTokens.UpdateTokenLastUsed(c.Request().Context(), pt.ID)

	c.Set("api_scopes", pt.Scopes)
	c.Set("project_id", pt.ProjectID)
	c.Set("environment_id", pt.EnvironmentID)

	return &models.UserClaims{
		UserID: "api-token-" + pt.ID,
		Email:  "api@" + pt.ProjectID + ".codedock.local",
		Role:   "api",
	}, nil
}

func (g *AuthGuard) validateJWT(c echo.Context, tokenStr string) (*models.UserClaims, error) {
	claimsMap, err := g.TokenService.ValidateToken(tokenStr)
	if err != nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
	}

	totpEnabled, _ := claimsMap["totpEnabled"].(bool)
	return &models.UserClaims{
		UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
		Email:       fmt.Sprintf("%v", claimsMap["email"]),
		Role:        models.UserRole(fmt.Sprintf("%v", claimsMap["role"])),
		TOTPEnabled: totpEnabled,
	}, nil
}

func (g *AuthGuard) baseAuth(c echo.Context, denyAPITokens bool) (*models.UserClaims, error) {
	if err := g.checkIPAllowlist(c); err != nil {
		return nil, err
	}
	tokenStr := ExtractTokenFromRequest(c)
	if tokenStr == "" {
		return nil, utils.Error(c, http.StatusUnauthorized, "missing authentication token")
	}

	if strings.HasPrefix(tokenStr, "vsl_tok_") {
		return g.validateAPIToken(c, tokenStr, denyAPITokens)
	}

	return g.validateJWT(c, tokenStr)
}
