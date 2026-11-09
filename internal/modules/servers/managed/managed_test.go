package managed

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type testClaims struct {
	UserID string
}

func TestManagedHandler_Catalog(t *testing.T) {
	e := echo.New()
	h := NewManagedHandler(nil)

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/managed/catalog", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if err := h.Catalog(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("authorized default catalog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/managed/catalog", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.Catalog(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestControlsHandler(t *testing.T) {
	e := echo.New()
	h := NewControlsHandler()

	t.Run("GetNetworkSettings default", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/srv-1/network-settings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.GetNetworkSettings(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("UpdateNetworkSettings", func(t *testing.T) {
		body := `{"internetAccess": false, "egress": ["10.0.0.0/8"]}`
		req := httptest.NewRequest(http.MethodPatch, "/servers/srv-1/network-settings", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.UpdateNetworkSettings(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetManagedInfo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/srv-1/managed/info", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.GetManagedInfo(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetManagedBootLogs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/servers/srv-1/managed/boot-logs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.GetManagedBootLogs(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetManagedSshStatus and SetManagedSsh", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/srv-1/managed/ssh", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.GetManagedSshStatus(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		body := `{"enabled": true}`
		reqSet := httptest.NewRequest(http.MethodPatch, "/servers/srv-1/managed/ssh", strings.NewReader(body))
		reqSet.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recSet := httptest.NewRecorder()
		cSet := e.NewContext(reqSet, recSet)
		cSet.SetParamNames("id")
		cSet.SetParamValues("srv-1")
		cSet.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.SetManagedSsh(cSet); err != nil {
			t.Fatal(err)
		}
		if recSet.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", recSet.Code)
		}
	})

	t.Run("GetManagedRuntimeStatus and Enable", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servers/srv-1/managed/runtime-api", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("srv-1")
		c.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.GetManagedRuntimeStatus(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		reqEnable := httptest.NewRequest(http.MethodPost, "/servers/srv-1/managed/runtime-api/enable", nil)
		recEnable := httptest.NewRecorder()
		cEnable := e.NewContext(reqEnable, recEnable)
		cEnable.SetParamNames("id")
		cEnable.SetParamValues("srv-1")
		cEnable.Set("user", &testClaims{UserID: "usr-1"})
		if err := h.EnableManagedRuntime(cEnable); err != nil {
			t.Fatal(err)
		}
		if recEnable.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", recEnable.Code)
		}
	})
}

func TestSerializeManagedModels(t *testing.T) {
	spec := ManagedServerSpec{
		Name:       "test-srv",
		ServerType: "cx22",
		Region:     "fsn1",
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal spec: %v", err)
	}
	var decoded ManagedServerSpec
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal spec: %v", err)
	}
	if decoded.Name != spec.Name {
		t.Errorf("expected name %s, got %s", spec.Name, decoded.Name)
	}
}
