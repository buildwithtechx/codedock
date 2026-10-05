package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
)

func TestBillingRoutesRespectDeploymentMode(t *testing.T) {
	cfg := config.Get()
	previous := *cfg
	t.Cleanup(func() { *cfg = previous })
	cfg.Security.JWTSecret = "billing-test-secret"
	tokens, err := authservices.NewTokenService()
	if err != nil {
		t.Fatal(err)
	}
	token, err := tokens.GenerateToken(&models.User{ID: "owner", Role: models.UserRoleOwner})
	if err != nil {
		t.Fatal(err)
	}
	guard := middleware.NewAuthGuard(tokens, nil, nil, nil, nil, nil)
	for _, cloud := range []bool{false, true} {
		cfg.Cloud.Enabled = cloud
		router := echo.New()
		configureEchoMiddleware(router)
		server := &Server{router: router, billingHandler: system.NewBillingHandler(nil)}
		api := router.Group("/api")
		server.registerBillingRoutes(api, api.Group("", guard.RequireAuth()))
		router.GET("/*", func(c echo.Context) error { return c.String(200, "dashboard") })
		for _, check := range []struct {
			method, path string
			cloudStatus  int
		}{
			{http.MethodGet, "/api/billing/config", 200},
			{http.MethodPost, "/api/billing/checkout", 401},
			{http.MethodPost, "/api/billing/webhook", 400},
		} {
			request := httptest.NewRequest(check.method, check.path, nil)
			if !cloud || check.path != "/api/billing/webhook" {
				credential := "invalid"
				if check.path == "/api/billing/config" {
					credential = token
				}
				request.Header.Set("Authorization", "Bearer "+credential)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			expected := http.StatusNotFound
			if cloud {
				expected = check.cloudStatus
			}
			if response.Code != expected {
				t.Fatalf("cloud=%v %s %s: expected %d, got %d", cloud, check.method, check.path, expected, response.Code)
			}
		}
		if cloud {
			request := httptest.NewRequest(http.MethodGet, "/api/billing/config", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("unauthenticated billing config exposed: %d", response.Code)
			}
		}
	}
}
