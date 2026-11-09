package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/auth"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/settings"
	"codedock/internal/modules/users"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type UserClaims = users.UserClaims

type SettingsProvider interface {
	GetSettings(context.Context) (*settings.ServerSettings, error)
}

type ProjectTokenProvider interface {
	GetTokenByHash(ctx context.Context, tokenHash string) (*projects.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
}

type OrganizationMemberProvider interface {
	GetMember(ctx context.Context, orgID, userID string) (*organizations.OrganizationMember, error)
}

type UserStatusProvider interface {
	GetUserByID(ctx context.Context, id string) (*users.User, error)
}

type ProjectProvider interface {
	Get(ctx context.Context, id string) (*projects.ProjectConfig, error)
}

type AuthGuard struct {
	PersonalTokens     PersonalTokenProvider
	PersonalResources  PersonalResourceProvider
	TokenService       *auth.TokenService
	Settings           SettingsProvider
	ProjectTokens      ProjectTokenProvider
	OrgMembers         OrganizationMemberProvider
	ProjectRepo        ProjectProvider
	UserStatusProvider UserStatusProvider
}

func NewAuthGuard(ts *auth.TokenService, sp SettingsProvider, pt ProjectTokenProvider, orgMembers OrganizationMemberProvider, pr ProjectProvider, userStatus UserStatusProvider) *AuthGuard {
	return &AuthGuard{TokenService: ts, Settings: sp, ProjectTokens: pt, OrgMembers: orgMembers, ProjectRepo: pr, UserStatusProvider: userStatus}
}

func (g *AuthGuard) checkIPAllowlist(c echo.Context) error {
	if g.Settings == nil {
		return nil
	}
	settings, err := g.Settings.GetSettings(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load server security policy")
	}
	if settings == nil || strings.TrimSpace(settings.IPAllowlist) == "" {
		return nil
	}

	clientIP := c.RealIP()
	allowedIPs := strings.Split(settings.IPAllowlist, ",")
	for _, ip := range allowedIPs {
		if strings.TrimSpace(ip) == clientIP {
			return nil
		}
	}
	return echo.NewHTTPError(http.StatusForbidden, "IP address not allowed")
}

func (g *AuthGuard) validateAPIToken(c echo.Context, tokenStr string, denyAPITokens bool) (*users.UserClaims, error) {
	if denyAPITokens {
		return nil, echo.NewHTTPError(http.StatusForbidden, "API tokens cannot access role-restricted endpoints")
	}
	if g.ProjectTokens == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "API tokens not supported")
	}
	tokenHash := auth.HashAPIToken(tokenStr)
	token, err := g.ProjectTokens.GetTokenByHash(c.Request().Context(), tokenHash)
	if err != nil || token == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid API token")
	}
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "API token has expired")
	}

	_ = g.ProjectTokens.UpdateTokenLastUsed(c.Request().Context(), token.ID)
	c.Set("project_id", token.ProjectID)

	return &users.UserClaims{
		UserID:   "api-token:" + token.ID,
		Email:    "api-token@" + token.ProjectID,
		Role:     users.UserRole("token"),
		PlanType: "pro",
	}, nil
}

func (g *AuthGuard) validateJWT(c echo.Context, tokenStr string) (*users.UserClaims, error) {
	claimsMap, err := g.TokenService.ValidateAccessToken(tokenStr)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
	}

	userID := fmt.Sprintf("%v", claimsMap["sub"])
	if g.UserStatusProvider != nil {
		u, err := g.UserStatusProvider.GetUserByID(c.Request().Context(), userID)
		if err != nil || u == nil || !u.IsActive {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "user account is deactivated")
		}
	}

	totpEnabled, _ := claimsMap["totp_enabled"].(bool)
	return &users.UserClaims{
		UserID:      userID,
		Email:       fmt.Sprintf("%v", claimsMap["email"]),
		Role:        users.UserRole(fmt.Sprintf("%v", claimsMap["role"])),
		TOTPEnabled: totpEnabled,
	}, nil
}

func (g *AuthGuard) baseAuth(c echo.Context, denyAPITokens bool) (*users.UserClaims, error) {
	if err := g.checkIPAllowlist(c); err != nil {
		return nil, err
	}
	tokenStr := ExtractTokenFromRequest(c)
	if tokenStr == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
	}

	if strings.HasPrefix(tokenStr, "vpt_") {
		return g.validatePersonalToken(c, tokenStr, denyAPITokens)
	}
	if strings.HasPrefix(tokenStr, "vsl_tok_") {
		return g.validateAPIToken(c, tokenStr, denyAPITokens)
	}

	return g.validateJWT(c, tokenStr)
}
