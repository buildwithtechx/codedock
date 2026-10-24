package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/handlers/system"
	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	systemservices "codedock/internal/services/system"
)

type mockUserRepoForCLI struct {
	admin *models.User
}

func (m *mockUserRepoForCLI) CreateUser(ctx context.Context, u *models.User) error { return nil }
func (m *mockUserRepoForCLI) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetUserByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetUserTOTPSecret(ctx context.Context, userID string) (string, []string, error) {
	return "", nil, nil
}
func (m *mockUserRepoForCLI) ListUsers(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForCLI) CountUsers(ctx context.Context) (int, error)          { return 1, nil }
func (m *mockUserRepoForCLI) UpdateUser(ctx context.Context, u *models.User) error { return nil }
func (m *mockUserRepoForCLI) DeleteUser(ctx context.Context, id string) error      { return nil }
func (m *mockUserRepoForCLI) CreatePAT(ctx context.Context, pat *models.PersonalAccessToken) error {
	return nil
}
func (m *mockUserRepoForCLI) ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) DeletePAT(ctx context.Context, id, userID string) error { return nil }
func (m *mockUserRepoForCLI) FindFoundingAdmin(ctx context.Context) (*models.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetAppearance(ctx context.Context, userID string) (*models.AppearancePreferences, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) SetAppearance(ctx context.Context, userID string, pref *models.AppearancePreferences) error {
	return nil
}

func TestLocalCLISession(t *testing.T) {
	cfg := config.Get()
	cfg.Security.JWTSecret = "jwt-secret-testing-1234567890"
	cfg.Cloud.Enabled = false

	os.Setenv("CODEDOCK_INTERNAL_TOKEN", "valid-operator-secret")
	defer os.Unsetenv("CODEDOCK_INTERNAL_TOKEN")

	tokenSvc, err := authservices.NewTokenService()
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}

	repo := &mockUserRepoForCLI{
		admin: &models.User{
			ID:    "admin-uuid-1",
			Email: "admin@codedock.local",
			Role:  models.UserRoleAdmin,
		},
	}

	handler := system.NewSystemHandler(systemservices.NewSystemService(), repo, tokenSvc)
	e := echo.New()

	t.Run("rejects missing internal token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/system/local-cli-session", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.CreateLocalCLISession(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("rejects browser Origin header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/system/local-cli-session", nil)
		req.Header.Set("X-Internal-Token", "valid-operator-secret")
		req.Header.Set("Origin", "http://attacker.com")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.CreateLocalCLISession(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden, got %d", rec.Code)
		}
	})

	t.Run("mints session for valid operator credential", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/system/local-cli-session", nil)
		req.Header.Set("X-Internal-Token", "valid-operator-secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.CreateLocalCLISession(c); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		var resp struct {
			Token     string `json:"token"`
			ExpiresAt string `json:"expiresAt"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Token == "" {
			t.Errorf("expected non-empty token")
		}
		exp, err := time.Parse(time.RFC3339, resp.ExpiresAt)
		if err != nil || exp.Before(time.Now()) {
			t.Errorf("expected future expiration timestamp, got %s", resp.ExpiresAt)
		}
	})
}
