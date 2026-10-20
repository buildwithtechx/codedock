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

func TestActionsHandler_Lifecycle(t *testing.T) {
	e := echo.New()
	handler := NewActionsHandler()

	// 1. List initial workflows (should be empty by default)
	req := httptest.NewRequest(http.MethodGet, "/actions/workflows", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.ListWorkflows(c); err != nil {
		t.Fatalf("unexpected ListWorkflows error: %v", err)
	}
	var wfResp struct {
		Data []*models.ActionWorkflow `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &wfResp); err != nil {
		t.Fatalf("failed to decode workflows: %v", err)
	}
	if len(wfResp.Data) != 0 {
		t.Fatalf("expected 0 workflows initially, got %d", len(wfResp.Data))
	}

	// 2. Create a new workflow
	createBody, _ := json.Marshal(models.CreateWorkflowRequest{
		Name:        "Deploy Production",
		Description: "Staged deployment pipeline",
		Trigger:     "manual",
		YAML:        "name: Deploy\non: workflow_dispatch\n",
	})
	req = httptest.NewRequest(http.MethodPost, "/actions/workflows", bytes.NewReader(createBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.CreateWorkflow(c); err != nil {
		t.Fatalf("unexpected CreateWorkflow error: %v", err)
	}
	var createdWf struct {
		Data models.ActionWorkflow `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &createdWf)
	if createdWf.Data.Name != "Deploy Production" {
		t.Fatalf("expected name 'Deploy Production', got %q", createdWf.Data.Name)
	}

	// 2b. Dispatch workflow to create a run
	req = httptest.NewRequest(http.MethodPost, "/actions/workflows/"+createdWf.Data.ID+"/dispatch", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(createdWf.Data.ID)
	if err := handler.DispatchWorkflow(c); err != nil {
		t.Fatalf("unexpected DispatchWorkflow error: %v", err)
	}

	// 3. List and get runs
	req = httptest.NewRequest(http.MethodGet, "/actions/runs", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.ListRuns(c); err != nil {
		t.Fatalf("unexpected ListRuns error: %v", err)
	}
	var runsResp struct {
		Data []*models.ActionRun `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &runsResp)
	if len(runsResp.Data) == 0 {
		t.Fatalf("expected dispatched run")
	}
	runID := runsResp.Data[0].ID

	// 4. Retry run
	req = httptest.NewRequest(http.MethodPost, "/actions/runs/"+runID+"/retry", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(runID)
	if err := handler.RetryRun(c); err != nil {
		t.Fatalf("unexpected RetryRun error: %v", err)
	}
	var retryResp struct {
		Data models.ActionRun `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &retryResp)
	if retryResp.Data.Status != "running" {
		t.Fatalf("expected status 'running', got %q", retryResp.Data.Status)
	}
	newRunID := retryResp.Data.ID

	// 5. Cancel run
	req = httptest.NewRequest(http.MethodPost, "/actions/runs/"+newRunID+"/cancel", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(newRunID)
	if err := handler.CancelRun(c); err != nil {
		t.Fatalf("unexpected CancelRun error: %v", err)
	}
	var cancelResp struct {
		Data models.ActionRun `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &cancelResp)
	if cancelResp.Data.Status != "cancelled" {
		t.Fatalf("expected status 'cancelled', got %q", cancelResp.Data.Status)
	}

	// 6. Inspect destination capabilities
	inspectBody, _ := json.Marshal(models.InspectDestinationRequest{ServerID: "srv-prod-1"})
	req = httptest.NewRequest(http.MethodPost, "/actions/runners/inspect", bytes.NewReader(inspectBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.InspectDestination(c); err != nil {
		t.Fatalf("unexpected InspectDestination error: %v", err)
	}

	// 7. Register a new runner
	runnerBody, _ := json.Marshal(models.CreateRunnerRequest{
		ServerID: "srv-prod-1",
		Name:     "ARM64 Native Runner",
		Mode:     "native",
		Labels:   []string{"arm64", "linux"},
		Enabled:  true,
	})
	req = httptest.NewRequest(http.MethodPost, "/actions/runners", bytes.NewReader(runnerBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.CreateRunner(c); err != nil {
		t.Fatalf("unexpected CreateRunner error: %v", err)
	}
	var runnerResp struct {
		Data models.ActionRunner `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &runnerResp)

	// 8. Reject duplicate runner on same server
	req = httptest.NewRequest(http.MethodPost, "/actions/runners", bytes.NewReader(runnerBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if err := handler.CreateRunner(c); err == nil {
		t.Fatalf("expected conflict on duplicate runner, got nil")
	}

	// 9. Delete workflow, runner, and run
	req = httptest.NewRequest(http.MethodDelete, "/actions/workflows/"+createdWf.Data.ID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(createdWf.Data.ID)
	if err := handler.DeleteWorkflow(c); err != nil {
		t.Fatalf("unexpected DeleteWorkflow error: %v", err)
	}

	req = httptest.NewRequest(http.MethodDelete, "/actions/runners/"+runnerResp.Data.ID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(runnerResp.Data.ID)
	if err := handler.DeleteRunner(c); err != nil {
		t.Fatalf("unexpected DeleteRunner error: %v", err)
	}

	req = httptest.NewRequest(http.MethodDelete, "/actions/runs/"+runID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(runID)
	if err := handler.DeleteRun(c); err != nil {
		t.Fatalf("unexpected DeleteRun error: %v", err)
	}
}
