package system

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

func TestIngressHandler_Lifecycle(t *testing.T) {
	e := echo.New()
	handler := NewIngressHandler()

	// 1. Discover
	discoverBody, _ := json.Marshal(models.IngressDiscoverRequest{
		CredentialID: "cf-cred-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/dns/ingress/discover", bytes.NewReader(discoverBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.Discover(c); err != nil {
		t.Fatalf("unexpected Discover error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 2. Connect
	connectBody, _ := json.Marshal(models.IngressConnectRequest{
		ServerID:          "srv-alpha",
		CredentialID:      "cf-cred-1",
		AccountID:         "cf-acc-primary",
		Name:              "Alpha Tunnel",
		Connector:         "managed",
		DefaultForDomains: true,
	})
	req = httptest.NewRequest(http.MethodPost, "/dns/ingress", bytes.NewReader(connectBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.Connect(c); err != nil {
		t.Fatalf("unexpected Connect error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var connectResp struct {
		Status string                   `json:"status"`
		Data   models.IngressConnection `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &connectResp); err != nil {
		t.Fatalf("failed to parse connect response: %v", err)
	}
	connID := connectResp.Data.ID
	if connID == "" {
		t.Fatalf("expected non-empty connection ID")
	}
	if connectResp.Data.ServerID != "srv-alpha" {
		t.Fatalf("expected serverId 'srv-alpha', got %q", connectResp.Data.ServerID)
	}
	if !connectResp.Data.DefaultForDomains {
		t.Fatalf("expected defaultForDomains to be true")
	}

	// 3. List
	req = httptest.NewRequest(http.MethodGet, "/dns/ingress?serverId=srv-alpha", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.List(c); err != nil {
		t.Fatalf("unexpected List error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 4. Update
	updateBody, _ := json.Marshal(models.IngressUpdateRequest{
		DefaultForDomains: false,
	})
	req = httptest.NewRequest(http.MethodPut, "/dns/ingress/"+connID, bytes.NewReader(updateBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(connID)
	if err := handler.Update(c); err != nil {
		t.Fatalf("unexpected Update error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 5. Retry
	req = httptest.NewRequest(http.MethodPost, "/dns/ingress/"+connID+"/retry", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(connID)
	if err := handler.Retry(c); err != nil {
		t.Fatalf("unexpected Retry error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 6. Domain mapping
	domainBody, _ := json.Marshal(models.IngressDomainRequest{
		ConnectionID: &connID,
	})
	req = httptest.NewRequest(http.MethodPut, "/dns/ingress/domains/app.codedock.dev", bytes.NewReader(domainBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("domainId")
	c.SetParamValues("app.codedock.dev")
	if err := handler.SetDomain(c); err != nil {
		t.Fatalf("unexpected SetDomain error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 7. Get Domain mapping
	req = httptest.NewRequest(http.MethodGet, "/dns/ingress/domains/app.codedock.dev", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("domainId")
	c.SetParamValues("app.codedock.dev")
	if err := handler.GetDomain(c); err != nil {
		t.Fatalf("unexpected GetDomain error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 8. Disconnect
	req = httptest.NewRequest(http.MethodDelete, "/dns/ingress/"+connID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(connID)
	if err := handler.Disconnect(c); err != nil {
		t.Fatalf("unexpected Disconnect error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}
