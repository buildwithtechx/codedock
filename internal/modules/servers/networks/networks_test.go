package networks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestServerNetworksHandler(t *testing.T) {
	e := echo.New()
	h := NewServerNetworksHandler()

	t.Run("List empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/system/networks", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if err := h.List(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	var createdNetID string
	var operationID string

	t.Run("Create network", func(t *testing.T) {
		body := `{"name": "wg-mesh", "mode": "wireguard", "subnetCidr": "10.50.0.0/16"}`
		req := httptest.NewRequest(http.MethodPost, "/system/networks", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if err := h.Create(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		var resp struct {
			Data ServerNetwork `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode resp: %v", err)
		}
		if resp.Data.Name != "wg-mesh" {
			t.Errorf("expected name wg-mesh, got %s", resp.Data.Name)
		}
		createdNetID = resp.Data.ID
		if resp.Data.Watchdog != nil {
			operationID = resp.Data.Watchdog.OperationID
		}
	})

	t.Run("Get network", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/system/networks/"+createdNetID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(createdNetID)
		if err := h.Get(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ReviseAccess", func(t *testing.T) {
		body := `{"accessPolicy": {"allowAll": true}}`
		req := httptest.NewRequest(http.MethodPut, "/system/networks/"+createdNetID+"/access", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(createdNetID)
		if err := h.ReviseAccess(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetWatchdog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/system/networks/operations/"+operationID+"/watchdog", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("operationId")
		c.SetParamValues(operationID)
		if err := h.GetWatchdog(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("CommitWatchdog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/networks/operations/"+operationID+"/watchdog/commit", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("operationId")
		c.SetParamValues(operationID)
		if err := h.CommitWatchdog(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("Delete network", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/system/networks/"+createdNetID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(createdNetID)
		if err := h.Delete(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
