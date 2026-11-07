package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"codedock/internal/models"
	"github.com/labstack/echo/v4"
)

type PersonalTokenProvider interface {
	GetPATByHash(context.Context, string) (*models.PersonalAccessToken, error)
}
type PersonalResourceProvider interface {
	ProjectForResource(context.Context, string, string) (string, error)
}

func (g *AuthGuard) validatePersonalToken(c echo.Context, raw string, deny bool) (*models.UserClaims, error) {
	if deny {
		return nil, echo.NewHTTPError(http.StatusForbidden, "personal tokens cannot access role-restricted endpoints")
	}
	if g.PersonalTokens == nil || g.UserStatusProvider == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "personal tokens unavailable")
	}
	hash := sha256.Sum256([]byte(raw))
	pat, err := g.PersonalTokens.GetPATByHash(c.Request().Context(), hex.EncodeToString(hash[:]))
	if err != nil || pat == nil || (pat.ExpiresAt != nil && !pat.ExpiresAt.After(time.Now())) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid, expired or revoked personal token")
	}
	if (pat.AccessLevel != "read" && pat.AccessLevel != "read_write") || (pat.ProjectScope != "all" && pat.ProjectScope != "specific") {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid personal token policy")
	}
	user, err := g.UserStatusProvider.GetUserByID(c.Request().Context(), pat.UserID)
	if err != nil || user == nil || !user.IsActive {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "token owner unavailable or inactive")
	}
	c.Set("personal_token", pat)
	return &models.UserClaims{UserID: user.ID, Email: user.Email, Role: user.Role, PlanType: string(user.PlanType)}, nil
}

func (g *AuthGuard) authorizePersonalRequest(c echo.Context) error {
	pat, ok := c.Get("personal_token").(*models.PersonalAccessToken)
	if !ok {
		return nil
	}
	route := strings.TrimPrefix(c.Path(), "/api/")
	if strings.Contains(route, "/tokens") {
		return echo.NewHTTPError(http.StatusForbidden, "use a session for token management")
	}
	if strings.Contains(route, "terminal") {
		return echo.NewHTTPError(http.StatusForbidden, "use a session for interactive terminals")
	}
	if strings.HasPrefix(route, "auth/") || strings.HasPrefix(route, "profile") || strings.HasPrefix(route, "users") {
		if route != "auth/me" && route != "profile" {
			return echo.NewHTTPError(http.StatusForbidden, "use a session for account and token management")
		}
		if c.Request().Method != http.MethodGet {
			return echo.NewHTTPError(http.StatusForbidden, "use a session for account changes")
		}
	}
	if route == "auth/me" || route == "profile" {
		return nil
	}
	method := c.Request().Method
	if pat.AccessLevel == "read" && method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions {
		return echo.NewHTTPError(http.StatusForbidden, "personal token is read only")
	}
	if pat.ProjectScope == "all" {
		return nil
	}
	var allowed []string
	if pat.AllowedProjects == nil || json.Unmarshal([]byte(*pat.AllowedProjects), &allowed) != nil || len(allowed) == 0 || g.PersonalResources == nil {
		return echo.NewHTTPError(http.StatusForbidden, "invalid personal token project scope")
	}
	check := func(kind, id string) bool {
		project, err := g.PersonalResources.ProjectForResource(c.Request().Context(), kind, id)
		if err != nil {
			return false
		}
		for _, permitted := range allowed {
			if project == permitted {
				return true
			}
		}
		return false
	}
	parts := strings.Split(route, "/")
	scoped := false
	for index, part := range parts {
		if !strings.HasPrefix(part, ":") {
			continue
		}
		parameter := strings.TrimPrefix(part, ":")
		if (parameter == "table" || parameter == "key") && scoped {
			continue
		}
		if index == 0 {
			return echo.NewHTTPError(http.StatusForbidden, "personal token requires a project resource")
		}
		kind := parts[index-1]
		if kind == "records" {
			kind = "backup-records"
		}
		if !check(kind, c.Param(parameter)) {
			return echo.NewHTTPError(http.StatusForbidden, "personal token does not permit this project resource")
		}
		scoped = true
	}
	if !scoped {
		return echo.NewHTTPError(http.StatusForbidden, "personal token requires a project scoped route")
	}
	for key, values := range c.QueryParams() {
		if kind := personalReferenceKind(key); kind != "" {
			for _, value := range values {
				if value != "" && !check(kind, value) {
					return echo.NewHTTPError(http.StatusForbidden, "personal token does not permit referenced project")
				}
			}
		}
	}
	mediaType, _, _ := mime.ParseMediaType(strings.ToLower(c.Request().Header.Get("Content-Type")))
	if c.Request().Body != nil {
		body, err := io.ReadAll(io.LimitReader(c.Request().Body, 1024*1024+1))
		c.Request().Body = io.NopCloser(bytes.NewReader(body))
		if err != nil || len(body) > 1024*1024 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid scoped token request body")
		}
		if len(body) > 0 && mediaType != "application/json" {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, "specific-project personal tokens require JSON request bodies")
		}
		var refs any
		if len(body) > 0 && json.Unmarshal(body, &refs) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request JSON")
		}
		if !personalReferencesAllowed(refs, check) {
			return echo.NewHTTPError(http.StatusForbidden, "personal token does not permit referenced project")
		}

	}
	return nil
}

func personalReferenceKind(key string) string {
	for suffix, kind := range map[string]string{"projectid": "projects", "environmentid": "environments", "serviceid": "services", "databaseid": "databases"} {
		if strings.HasSuffix(strings.ReplaceAll(strings.ToLower(key), "_", ""), suffix) {
			return kind
		}
	}
	return ""
}
func personalReferencesAllowed(value any, check func(string, string) bool) bool {
	switch value := value.(type) {
	case map[string]any:
		for key, child := range value {
			if kind := personalReferenceKind(key); kind != "" {
				id, ok := child.(string)
				if !ok || (id != "" && !check(kind, id)) {
					return false
				}
			}
			if !personalReferencesAllowed(child, check) {
				return false
			}
		}
	case []any:
		for _, child := range value {
			if !personalReferencesAllowed(child, check) {
				return false
			}
		}
	}
	return true
}
