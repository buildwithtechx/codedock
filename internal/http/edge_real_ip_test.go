package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
)

func TestEchoIPExtractor(t *testing.T) {
	os.Setenv("CODEDOCK_EDGE_TRUSTED_PROXIES", "10.0.0.0/8, 192.168.1.1")
	defer os.Unsetenv("CODEDOCK_EDGE_TRUSTED_PROXIES")

	e := echo.New()
	configureEchoIPExtractor(e)

	e.GET("/ip", func(c echo.Context) error {
		return c.String(http.StatusOK, c.RealIP())
	})

	t.Run("untrusted peer ignores spoofed CF-Connecting-IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "203.0.113.195:12345"
		req.Header.Set("CF-Connecting-IP", "1.1.1.1")

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Body.String() != "203.0.113.195" {
			t.Errorf("expected 203.0.113.195, got %q", rec.Body.String())
		}
	})

	t.Run("trusted Cloudflare peer trusts CF-Connecting-IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "173.245.48.50:54321" // Within 173.245.48.0/20
		req.Header.Set("CF-Connecting-IP", "198.51.100.42")

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Body.String() != "198.51.100.42" {
			t.Errorf("expected 198.51.100.42, got %q", rec.Body.String())
		}
	})

	t.Run("trusted operator proxy trusts declared header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "10.0.0.5:54321" // Within 10.0.0.0/8
		req.Header.Set("CF-Connecting-IP", "203.0.113.99")

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Body.String() != "203.0.113.99" {
			t.Errorf("expected 203.0.113.99, got %q", rec.Body.String())
		}
	})
}

func TestEchoIPExtractorXForwardedFor(t *testing.T) {
	os.Setenv("CODEDOCK_EDGE_TRUSTED_PROXIES", "10.0.0.0/8")
	defer os.Unsetenv("CODEDOCK_EDGE_TRUSTED_PROXIES")

	cfg := config.Get()
	cfg.Traefik.RealIPHeader = "X-Forwarded-For"

	e := echo.New()
	configureEchoIPExtractor(e)

	e.GET("/ip", func(c echo.Context) error {
		return c.String(http.StatusOK, c.RealIP())
	})

	t.Run("walks XFF from right to find last untrusted hop", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = "10.0.0.1:443"
		req.Header.Set("X-Forwarded-For", "client.spoof, 198.51.100.77, 10.0.0.2")

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Body.String() != "198.51.100.77" {
			t.Errorf("expected 198.51.100.77, got %q", rec.Body.String())
		}
	})
}
