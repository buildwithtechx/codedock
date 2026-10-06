package middleware

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type personalTestStore struct {
	token *models.PersonalAccessToken
	user  *models.User
	hash  string
}

func (s *personalTestStore) GetPATByHash(_ context.Context, hash string) (*models.PersonalAccessToken, error) {
	s.hash = hash
	return s.token, nil
}
func (s *personalTestStore) GetUserByID(context.Context, string) (*models.User, error) {
	return s.user, nil
}
func (s *personalTestStore) ProjectForResource(_ context.Context, kind, id string) (string, error) {
	if kind == "projects" {
		return id, nil
	}
	if id == "app" {
		return "allowed", nil
	}
	return "other", nil
}

func TestPersonalTokensEnforcePolicyAndOwnerStatus(t *testing.T) {
	for _, test := range []struct {
		name, method, path, id, body               string
		read, specific, expired, revoked, inactive bool
		want                                       int
	}{
		{name: "read authenticates", method: "GET", path: "/api/apps/:id", id: "app", read: true, want: 200},
		{name: "read denies write", method: "PUT", path: "/api/apps/:id", id: "app", read: true, want: 403},
		{name: "write authenticates", method: "PUT", path: "/api/apps/:id", id: "app", want: 200},
		{name: "scoped token identity", method: "GET", path: "/api/auth/me", specific: true, want: 200},
		{name: "selected project", method: "GET", path: "/api/apps/:id", id: "app", specific: true, want: 200},
		{name: "other project", method: "GET", path: "/api/apps/:id", id: "other", specific: true, want: 403},
		{name: "broad listing", method: "GET", path: "/api/apps", specific: true, want: 403},
		{name: "nested resource outside scope", method: "GET", path: "/api/projects/:projectId/apps/:id", id: "other", specific: true, want: 403},
		{name: "project token escalation", method: "POST", path: "/api/projects/:projectId/tokens", specific: true, want: 403},
		{name: "nested body outside scope", method: "PUT", path: "/api/apps/:id", id: "app", specific: true, body: `{"settings":{"targetProjectId":"other"}}`, want: 403},
		{name: "move outside scope", method: "PUT", path: "/api/apps/:id", id: "app", specific: true, body: `{"projectId":"other"}`, want: 403},
		{name: "token escalation", method: "POST", path: "/api/profile/tokens", want: 403},
		{name: "expiry", method: "GET", path: "/api/apps/:id", id: "app", expired: true, want: 401},
		{name: "revocation", method: "GET", path: "/api/apps/:id", id: "app", revoked: true, want: 401},
		{name: "inactive owner", method: "GET", path: "/api/apps/:id", id: "app", inactive: true, want: 401},
	} {
		t.Run(test.name, func(t *testing.T) {
			allowed := `["allowed"]`
			expires := time.Now().Add(time.Hour)
			store := &personalTestStore{token: &models.PersonalAccessToken{UserID: "owner", AccessLevel: "read_write", ProjectScope: "all", ExpiresAt: &expires}, user: &models.User{ID: "owner", IsActive: !test.inactive, Role: models.UserRoleMember}}
			if test.read {
				store.token.AccessLevel = "read"
			}
			if test.specific {
				store.token.ProjectScope = "specific"
				store.token.AllowedProjects = &allowed
			}
			if test.expired {
				expired := time.Now().Add(-time.Hour)
				store.token.ExpiresAt = &expired
			}
			if test.revoked {
				store.token = nil
			}
			guard := &AuthGuard{PersonalTokens: store, UserStatusProvider: store, PersonalResources: store}
			e := echo.New()
			e.Add(test.method, test.path, func(c echo.Context) error {
				claims := GetUserClaimsFromContext(c.Request().Context())
				if claims.UserID != "owner" {
					return fmt.Errorf("wrong owner")
				}
				return c.NoContent(200)
			}, guard.RequireAuth())
			url := strings.Replace(test.path, ":id", test.id, 1)
			url = strings.ReplaceAll(url, ":projectId", "allowed")
			req := httptest.NewRequest(test.method, url, strings.NewReader(test.body))
			req.Header.Set("Authorization", "Bearer vpt_secret")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != test.want {
				t.Fatalf("got %d: %s", rec.Code, rec.Body.String())
			}
			hash := sha256.Sum256([]byte("vpt_secret"))
			if store.hash != hex.EncodeToString(hash[:]) {
				t.Fatal("token lookup did not hash the bearer credential")
			}
		})
	}
}

func TestPersonalTokenCannotAccessRoleRestrictedEndpoints(t *testing.T) {
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	guard := &AuthGuard{}
	if _, err := guard.validatePersonalToken(c, "vpt_secret", true); err == nil {
		t.Fatal("personal token bypassed role restriction")
	}
}
