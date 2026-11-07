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

func TestServerNetworksHandler_LifecycleAndWatchdog(t *testing.T) {
	e := echo.New()
	handler := NewServerNetworksHandler()

	// 1. Create network with 120s rollback watchdog
	createBody, _ := json.Marshal(models.CreateNetworkRequest{
		Name:             "Alpha Mesh",
		Mode:             "wireguard",
		SubnetCIDR:       "10.99.0.0/16",
		VerificationPort: 51821,
		Members: []models.NetworkMember{
			{ServerID: "srv-1", PrivateIP: "10.99.0.1", InterfaceName: "wg0", Status: "ready", LatencyMS: 1.2},
			{ServerID: "srv-2", PrivateIP: "10.99.0.2", InterfaceName: "wg0", Status: "ready", LatencyMS: 2.4},
		},
		AccessPolicy: models.NetworkAccessPolicy{
			Mode: "full_mesh",
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/system/networks", bytes.NewReader(createBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.Create(c); err != nil {
		t.Fatalf("unexpected Create error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var createResp struct {
		Status string               `json:"status"`
		Data   models.ServerNetwork `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create network response: %v", err)
	}
	net := createResp.Data
	if net.ID == "" {
		t.Fatalf("expected non-empty network ID")
	}
	if net.Watchdog == nil || !net.Watchdog.Armed {
		t.Fatalf("expected armed watchdog on new network")
	}
	if net.Watchdog.TimeoutSeconds != 120 {
		t.Fatalf("expected 120s timeout, got %d", net.Watchdog.TimeoutSeconds)
	}
	opID := net.Watchdog.OperationID

	// 2. Get watchdog status
	req = httptest.NewRequest(http.MethodGet, "/system/networks/operations/"+opID+"/watchdog", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operationId")
	c.SetParamValues(opID)
	if err := handler.GetWatchdog(c); err != nil {
		t.Fatalf("unexpected GetWatchdog error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. Commit watchdog
	req = httptest.NewRequest(http.MethodPost, "/system/networks/operations/"+opID+"/watchdog/commit", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operationId")
	c.SetParamValues(opID)
	if err := handler.CommitWatchdog(c); err != nil {
		t.Fatalf("unexpected CommitWatchdog error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 4. Verify network status is verified
	req = httptest.NewRequest(http.MethodGet, "/system/networks/"+net.ID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(net.ID)
	if err := handler.Get(c); err != nil {
		t.Fatalf("unexpected Get error: %v", err)
	}
	var getResp struct {
		Status string               `json:"status"`
		Data   models.ServerNetwork `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}
	if getResp.Data.Status != "verified" {
		t.Fatalf("expected status 'verified', got %q", getResp.Data.Status)
	}
	if getResp.Data.Watchdog.Status != "committed" {
		t.Fatalf("expected watchdog 'committed', got %q", getResp.Data.Watchdog.Status)
	}
}
