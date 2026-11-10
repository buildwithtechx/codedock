package system_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/adapters"
	"codedock/internal/config"
	authservices "codedock/internal/modules/auth"
	"codedock/internal/modules/system"
	"codedock/internal/modules/users"
)

type mockUserRepoForCLI struct {
	admin *users.User
}

func (m *mockUserRepoForCLI) CreateUser(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForCLI) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetUserByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) GetUserByID(ctx context.Context, id string) (*users.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetUserTOTPSecret(ctx context.Context, userID string) (string, []string, error) {
	return "", nil, nil
}
func (m *mockUserRepoForCLI) ListUsers(ctx context.Context, limit, offset int) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForCLI) CountUsers(ctx context.Context) (int, error)         { return 1, nil }
func (m *mockUserRepoForCLI) UpdateUser(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForCLI) DeleteUser(ctx context.Context, id string) error     { return nil }
func (m *mockUserRepoForCLI) CreatePAT(ctx context.Context, pat *users.PersonalAccessToken) error {
	return nil
}
func (m *mockUserRepoForCLI) ListPATs(ctx context.Context, userID string) ([]*users.PersonalAccessToken, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) DeletePAT(ctx context.Context, id, userID string) error { return nil }
func (m *mockUserRepoForCLI) GetPATByHash(ctx context.Context, hash string) (*users.PersonalAccessToken, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) FindFoundingAdmin(ctx context.Context) (*users.User, error) {
	return m.admin, nil
}
func (m *mockUserRepoForCLI) GetAppearance(ctx context.Context, userID string) (*users.AppearancePreferences, error) {
	return nil, nil
}
func (m *mockUserRepoForCLI) SetAppearance(ctx context.Context, userID string, pref *users.AppearancePreferences) error {
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
		admin: &users.User{
			ID:    "admin-uuid-1",
			Email: "admin@codedock.local",
			Role:  users.UserRoleAdmin,
		},
	}

	handler := system.NewSystemHandler(
		&adapters.SystemServiceAdapter{Svc: system.NewSystemService()},
		&adapters.SystemAdminFinderAdapter{UserRepo: repo},
		&adapters.SystemTokenGenAdapter{TokenService: tokenSvc},
	)

	e := echo.New()

	t.Run("fails when header is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/cli-session", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateLocalCLISession(c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 unauthorized, got %v", err)
		}
	})

	t.Run("fails when token is incorrect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/cli-session", nil)
		req.Header.Set("X-Codedock-Internal-Token", "wrong-secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateLocalCLISession(c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 unauthorized, got %v", err)
		}
	})

	t.Run("succeeds with valid internal token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/cli-session", nil)
		req.Header.Set("X-Codedock-Internal-Token", "valid-operator-secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateLocalCLISession(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		var res struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if res.Token == "" {
			t.Fatal("expected non-empty token")
		}

		claims, err := tokenSvc.ValidateAccessToken(res.Token)
		if err != nil {
			t.Fatalf("generated token is invalid: %v", err)
		}
		if claims["sub"] != "admin-uuid-1" {
			t.Errorf("expected sub to be admin-uuid-1, got %v", claims["sub"])
		}
		if claims["role"] != "admin" {
			t.Errorf("expected role to be admin, got %v", claims["role"])
		}
	})
}
