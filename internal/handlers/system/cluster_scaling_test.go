package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

func TestClusterHandler_ScalingLifecycle(t *testing.T) {
	e := echo.New()
	handler := NewClusterHandler(nil, nil)

	clusterID := "cls-prod-1"

	// 1. Initial scaling status is not_enabled
	req := httptest.NewRequest(http.MethodGet, "/system/compute-clusters/"+clusterID+"/scaling", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("clusterId")
	c.SetParamValues(clusterID)
	if err := handler.GetScaling(c); err != nil {
		t.Fatalf("unexpected GetScaling error: %v", err)
	}
	var getResp struct {
		Status string                 `json:"status"`
		Data   *models.ClusterScaling `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse get response: %v", err)
	}
	if getResp.Data.Status != "not_enabled" {
		t.Fatalf("expected status 'not_enabled', got %q", getResp.Data.Status)
	}

	// 2. Setup scaling
	req = httptest.NewRequest(http.MethodPost, "/system/compute-clusters/"+clusterID+"/scaling/setup", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("clusterId")
	c.SetParamValues(clusterID)
	if err := handler.SetupScaling(c); err != nil {
		t.Fatalf("unexpected SetupScaling error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. Retry scaling
	req = httptest.NewRequest(http.MethodPost, "/system/compute-clusters/"+clusterID+"/scaling/retry", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("clusterId")
	c.SetParamValues(clusterID)
	if err := handler.RetryScaling(c); err != nil {
		t.Fatalf("unexpected RetryScaling error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 4. Remove scaling
	req = httptest.NewRequest(http.MethodPost, "/system/compute-clusters/"+clusterID+"/scaling/remove", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("clusterId")
	c.SetParamValues(clusterID)
	if err := handler.RemoveScaling(c); err != nil {
		t.Fatalf("unexpected RemoveScaling error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 5. Verify status is removed
	req = httptest.NewRequest(http.MethodGet, "/system/compute-clusters/"+clusterID+"/scaling", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("clusterId")
	c.SetParamValues(clusterID)
	if err := handler.GetScaling(c); err != nil {
		t.Fatalf("unexpected GetScaling error: %v", err)
	}
	var finalResp struct {
		Status string                 `json:"status"`
		Data   *models.ClusterScaling `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &finalResp); err != nil {
		t.Fatalf("failed to parse get response: %v", err)
	}
	if finalResp.Data.Status != "removed" {
		t.Fatalf("expected status 'removed', got %q", finalResp.Data.Status)
	}
}
