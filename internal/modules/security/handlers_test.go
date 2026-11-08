package security

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil)
	if h == nil {
		t.Fatal("expected NewHandler to return non-nil instance")
	}
}

func TestHandlerUnauthorizedContext(t *testing.T) {
	e := echo.New()
	h := NewHandler(nil)

	// GetAdvisories without auth
	req := httptest.NewRequest(http.MethodGet, "/updates/advisories", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetAdvisories(c)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
	}

	// GetProjectSecurity without auth
	req = httptest.NewRequest(http.MethodGet, "/updates/proj-123/security", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("projectId")
	c.SetParamValues("proj-123")

	err = h.GetProjectSecurity(c)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestSetPolicyRequestBinding(t *testing.T) {
	body := []byte(`{"autoRefresh":true}`)
	var req SetPolicyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("failed to unmarshal SetPolicyRequest: %v", err)
	}
	if !req.AutoRefresh {
		t.Error("expected autoRefresh to be true")
	}
}

func TestApplyRemediationRequestBinding(t *testing.T) {
	body := []byte(`{
		"advisoryId": "adv-101",
		"deploymentID": "dep-202",
		"serviceId": "svc-303"
	}`)
	var req ApplyRemediationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("failed to unmarshal ApplyRemediationRequest: %v", err)
	}
	if req.AdvisoryID != "adv-101" {
		t.Errorf("expected advisoryId adv-101, got %s", req.AdvisoryID)
	}
	if req.ServiceID == nil || *req.ServiceID != "svc-303" {
		t.Errorf("expected serviceId svc-303, got %+v", req.ServiceID)
	}
}

func TestRegisterRoutes(t *testing.T) {
	e := echo.New()
	g := e.Group("/api")
	h := NewHandler(nil)

	RegisterRoutes(g, h)

	routes := e.Routes()
	if len(routes) != 5 {
		t.Fatalf("expected 5 routes registered, got %d", len(routes))
	}
}
