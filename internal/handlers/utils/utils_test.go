package handlerutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/config"
)

func TestIsAllowedWebSocketOrigin(t *testing.T) {
	cfg := config.Get()
	cfg.Server.ExtraCORSOrigins = []string{"https://app.custom-domain.com", "http://external-sdk.local:3000"}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ws", nil)

	t.Run("allows exact match host", func(t *testing.T) {
		if !IsAllowedWebSocketOrigin(req, "http://localhost:8080") {
			t.Errorf("expected localhost:8080 to be allowed")
		}
	})

	t.Run("allows extra cors origins", func(t *testing.T) {
		if !IsAllowedWebSocketOrigin(req, "https://app.custom-domain.com") {
			t.Errorf("expected configured extra cors origin to be allowed")
		}
		if !IsAllowedWebSocketOrigin(req, "http://external-sdk.local:3000") {
			t.Errorf("expected configured extra cors origin to be allowed")
		}
	})

	t.Run("denies unknown origins", func(t *testing.T) {
		if IsAllowedWebSocketOrigin(req, "https://malicious.evil.com") {
			t.Errorf("expected malicious origin to be denied")
		}
	})

	t.Run("denies empty origin", func(t *testing.T) {
		if IsAllowedWebSocketOrigin(req, "") {
			t.Errorf("expected empty origin to be denied")
		}
	})
}
