package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type mockUserRepoForAppearance struct {
	repositories.UserRepository
	appearance *models.AppearancePreferences
}

func (m *mockUserRepoForAppearance) GetAppearance(ctx context.Context, userID string) (*models.AppearancePreferences, error) {
	return m.appearance, nil
}

func (m *mockUserRepoForAppearance) SetAppearance(ctx context.Context, userID string, pref *models.AppearancePreferences) error {
	m.appearance = pref
	return nil
}

func TestAppearanceHandler(t *testing.T) {
	repo := &mockUserRepoForAppearance{}
	handler := NewAppearanceHandler(repo)
	e := echo.New()

	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", &models.UserClaims{
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

		var resp map[string]models.AppearancePreferences
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
