package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/handlers/system"
)

func TestBillingRoutesRespectDeploymentMode(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Cloud.Enabled
	t.Cleanup(func() { cfg.Cloud.Enabled = previous })
	for _, cloud := range []bool{false, true} {
		cfg.Cloud.Enabled = cloud
		router := echo.New()
		configureEchoMiddleware(router)
		server := &Server{router: router, billingHandler: system.NewBillingHandler(nil)}
		api := router.Group("/api")
		server.registerBillingRoutes(api, api.Group(""))
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
				request.Header.Set("Authorization", "Bearer test")
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
	}
}
