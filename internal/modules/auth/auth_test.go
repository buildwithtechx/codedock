package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type tokenCreationStore struct {
	UserRepository
	created *PersonalAccessToken
}

func (s *tokenCreationStore) CreatePAT(_ context.Context, pat *PersonalAccessToken) error {
	s.created = pat
	return nil
}

func TestPersonalTokenCreationMatchesUIContract(t *testing.T) {
	store := &tokenCreationStore{}
	service := NewUserService(store)
	token, plain, err := service.CreatePAT(context.Background(), "user", "automation", "read", "specific", []string{"project"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token.ExpiresAt != nil {
		t.Fatal("no expiration token acquired an expiry")
	}
	if len(plain) != 68 || plain[:4] != "vpt_" {
		t.Fatal("invalid bearer token")
	}
	digest := sha256.Sum256([]byte(plain))
	if store.created.TokenHash != hex.EncodeToString(digest[:]) {
		t.Fatal("credential was not hashed")
	}
	encoded, err := json.Marshal(CreatePATResponse{Token: token, Plain: plain})
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Plain string `json:"plain"`
		Token struct {
			AllowedProjects []string `json:"allowedProjects"`
		} `json:"token"`
	}
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatal(err)
	}
	if result.Plain != plain || len(result.Token.AllowedProjects) != 1 || result.Token.AllowedProjects[0] != "project" {
		t.Fatal("token response does not match UI")
	}
}

func TestPersonalTokenCreationRejectsInvalidPolicy(t *testing.T) {
	expired := time.Now().Add(-time.Minute)
	for _, test := range []struct {
		access, scope string
		projects      []string
		expires       *time.Time
	}{
		{access: "read", scope: "specific", projects: []string{" project "}},
		{access: "admin", scope: "all"},
		{access: "read", scope: "unknown"},
		{access: "read", scope: "specific"},
		{access: "read", scope: "specific", projects: []string{""}},
		{access: "read", scope: "all", expires: &expired},
	} {
		store := &tokenCreationStore{}
		service := NewUserService(store)
		if _, _, err := service.CreatePAT(context.Background(), "user", "token", test.access, test.scope, test.projects, test.expires); err == nil || store.created != nil {
			t.Fatal("invalid policy was persisted")
		}
	}
}

type mockUserRepoForAppearance struct {
	UserRepository
	appearance *AppearancePreferences
}

func (m *mockUserRepoForAppearance) GetAppearance(ctx context.Context, userID string) (*AppearancePreferences, error) {
	return m.appearance, nil
}

func (m *mockUserRepoForAppearance) SetAppearance(ctx context.Context, userID string, pref *AppearancePreferences) error {
	m.appearance = pref
	return nil
}

func TestAppearanceHandler(t *testing.T) {
	repo := &mockUserRepoForAppearance{}
	handler := NewAppearanceHandler(repo)
	e := echo.New()

	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", &UserClaims{
				UserID: "usr_123",
				Email:  "user@example.com",
			})
			return next(c)
		}
	}

	e.GET("/settings/appearance", handler.GetAppearance, authMiddleware)
	e.PUT("/settings/appearance", handler.SetAppearance, authMiddleware)

	t.Run("get appearance default when nil", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/settings/appearance", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]AppearancePreferences
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if resp["appearance"].ActiveTheme != "default" || resp["appearance"].Mode != "dark" {
			t.Fatalf("unexpected default appearance: %+v", resp["appearance"])
		}
	})

	t.Run("set valid appearance", func(t *testing.T) {
		payload := map[string]any{
			"appearance": map[string]any{
				"activeTheme": "dracula",
				"mode":        "dark",
				"fontFamily":  "Fira Code, monospace",
				"radius":      "12px",
				"tokens": map[string]string{
					"--accent": "#ff79c6",
				},
			},
		}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/settings/appearance", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		if repo.appearance == nil || repo.appearance.ActiveTheme != "dracula" {
			t.Fatalf("expected appearance saved in repo, got %+v", repo.appearance)
		}
	})

	t.Run("set invalid appearance with injection", func(t *testing.T) {
		payload := map[string]any{
			"appearance": map[string]any{
				"activeTheme": "dracula",
				"mode":        "dark",
				"tokens": map[string]string{
					"--bg": "red; background: blue",
				},
			},
		}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/settings/appearance", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 bad request, got %d: %s", rec.Code, rec.Body.String())
		}
	})
}
