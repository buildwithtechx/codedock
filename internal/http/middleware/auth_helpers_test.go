package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestExtractTokenFromRequest(t *testing.T) {
	e := echo.New()

	t.Run("extracts from Authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer test-token-123")
		c := e.NewContext(req, httptest.NewRecorder())

		token := ExtractTokenFromRequest(c)
		if token != "test-token-123" {
			t.Fatalf("expected test-token-123, got: %q", token)
		}
	})

	t.Run("extracts from query parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?token=query-token-456", nil)
		c := e.NewContext(req, httptest.NewRecorder())

		token := ExtractTokenFromRequest(c)
		if token != "query-token-456" {
			t.Fatalf("expected query-token-456, got: %q", token)
		}
	})

	t.Run("extracts from Sec-WebSocket-Protocol comma pair", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Sec-WebSocket-Protocol", "bearer, ws-token-789")
		c := e.NewContext(req, httptest.NewRecorder())

		token := ExtractTokenFromRequest(c)
		if token != "ws-token-789" {
			t.Fatalf("expected ws-token-789, got: %q", token)
		}
	})

	t.Run("extracts from Sec-WebSocket-Protocol space separated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Sec-WebSocket-Protocol", "Bearer ws-space-token")
		c := e.NewContext(req, httptest.NewRecorder())

		token := ExtractTokenFromRequest(c)
		if token != "ws-space-token" {
			t.Fatalf("expected ws-space-token, got: %q", token)
		}
	})

	t.Run("extracts from codedock_token cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "codedock_token", Value: "cookie-token-000"})
		c := e.NewContext(req, httptest.NewRecorder())

		token := ExtractTokenFromRequest(c)
		if token != "cookie-token-000" {
			t.Fatalf("expected cookie-token-000, got: %q", token)
		}
	})
}
