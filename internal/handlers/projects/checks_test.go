package projects

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

func TestChecksHandler_AdmissionAndConfig(t *testing.T) {
	e := echo.New()
	handler := NewChecksHandler()

	// 1. Test GetAdmission
	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj-123/deployments/admission", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("proj-123")

	if err := handler.GetAdmission(c); err != nil {
		t.Fatalf("GetAdmission returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var admissionResp struct {
		Data *models.BuildAdmissionStatus `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &admissionResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !admissionResp.Data.CanAdmit {
		t.Fatalf("expected CanAdmit to be true, got false")
	}

	// 2. Test GetChecks default config
	req = httptest.NewRequest(http.MethodGet, "/api/projects/proj-123/checks", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("proj-123")

	if err := handler.GetChecks(c); err != nil {
		t.Fatalf("GetChecks returned error: %v", err)
	}
	var getResp struct {
		Data *models.ProjectChecksConfig `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !getResp.Data.ReportCommitStatus {
		t.Fatalf("expected ReportCommitStatus true, got false")
	}

	// 3. Test UpdateChecks
	updateBody := models.UpdateChecksConfigRequest{
		Enabled:              true,
		ReportCommitStatus:   true,
		RequirePassingChecks: true,
		TargetBranches:       []string{"production", "staging"},
		ContextName:          "custom/status",
	}
	raw, _ := json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPut, "/api/projects/proj-123/checks", bytes.NewReader(raw))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("proj-123")

	if err := handler.UpdateChecks(c); err != nil {
		t.Fatalf("UpdateChecks returned error: %v", err)
	}
	var updateResp struct {
		Data *models.ProjectChecksConfig `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !updateResp.Data.RequirePassingChecks {
		t.Fatalf("expected RequirePassingChecks to be true")
	}
	if updateResp.Data.ContextName != "custom/status" {
		t.Fatalf("expected custom/status, got %s", updateResp.Data.ContextName)
	}
}
