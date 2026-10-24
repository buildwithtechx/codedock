package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetMeta(t *testing.T) {
	e := echo.New()
	s := &Server{}

	req := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := s.getMeta(c); err != nil {
		t.Fatalf("getMeta returned unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res struct {
		Name         string          `json:"name"`
		Version      string          `json:"version"`
		Capabilities map[string]bool `json:"capabilities"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Name != "codedock" {
		t.Errorf("expected name codedock, got %s", res.Name)
	}
	if res.Version != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %s", res.Version)
	}
	if !res.Capabilities["sseLogs"] {
		t.Errorf("expected sseLogs capability to be true")
	}
	if !res.Capabilities["push"] {
		t.Errorf("expected push capability to be true")
	}
	if !res.Capabilities["actions"] {
		t.Errorf("expected actions capability to be true")
	}
	if !res.Capabilities["apps"] {
		t.Errorf("expected apps capability to be true")
	}
}
