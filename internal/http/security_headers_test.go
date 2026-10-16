package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestFramingProtectionHeaders(t *testing.T) {
	e := echo.New()
	configureEchoMiddleware(e)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if xfo := rec.Header().Get("X-Frame-Options"); xfo != "DENY" {
		t.Errorf("expected X-Frame-Options: DENY, got %q", xfo)
	}

	if csp := rec.Header().Get("Content-Security-Policy"); csp != "frame-ancestors 'none'" {
		t.Errorf("expected Content-Security-Policy: frame-ancestors 'none', got %q", csp)
	}
}
