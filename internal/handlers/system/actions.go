package system

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

type ActionsHandler struct {
	mu        sync.RWMutex
	workflows map[string]*models.ActionWorkflow
	runs      map[string]*models.ActionRun
	runners   map[string]*models.ActionRunner
	runCount  int
}

func NewActionsHandler() *ActionsHandler {
	return &ActionsHandler{
		workflows: make(map[string]*models.ActionWorkflow),
		runs:      make(map[string]*models.ActionRun),
		runners:   make(map[string]*models.ActionRunner),
	}
}

func (h *ActionsHandler) ListWorkflows(c echo.Context) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []*models.ActionWorkflow
	for _, wf := range h.workflows {
		list = append(list, wf)
	}
	if list == nil {
		list = make([]*models.ActionWorkflow, 0)
	}
	return utils.Success(c, "Workflows retrieved", list)
}

func (h *ActionsHandler) CreateWorkflow(c echo.Context) error {
	var req models.CreateWorkflowRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create workflow payload")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "workflow name is required")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	wf := &models.ActionWorkflow{
		ID:          "wf-" + hex.EncodeToString(generateActionBytes(6)),
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Trigger:     req.Trigger,
		YAML:        req.YAML,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if wf.Trigger == "" {
		wf.Trigger = "push"
	}

	h.workflows[wf.ID] = wf
	return utils.Success(c, "Workflow created", wf)
}

func (h *ActionsHandler) DispatchWorkflow(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	wf, exists := h.workflows[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "workflow not found")
	}

	h.runCount++
	now := time.Now().UTC().Format(time.RFC3339)
	newRun := &models.ActionRun{
		ID:          "run-" + hex.EncodeToString(generateActionBytes(6)),
		WorkflowID:  wf.ID,
		Name:        wf.Name,
		Number:      h.runCount,
		Status:      "running",
		TriggeredBy: "manual",
		Logs:        "[Dispatch] Workflow triggered manually\n[Worker] Picked up by runner\n",
		StartedAt:   now,
		DurationS:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	h.runs[newRun.ID] = newRun
	return utils.Success(c, "Workflow dispatched", newRun)
}

func (h *ActionsHandler) ListRuns(c echo.Context) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []*models.ActionRun
	for _, r := range h.runs {
		list = append(list, r)
	}
	if list == nil {
		list = make([]*models.ActionRun, 0)
	}
	return utils.Success(c, "Action runs retrieved", list)
}

func (h *ActionsHandler) GetRun(c echo.Context) error {
	id := c.Param("id")
	h.mu.RLock()
	defer h.mu.RUnlock()

	run, exists := h.runs[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "action run not found")
	}
	return utils.Success(c, "Action run retrieved", run)
}

func (h *ActionsHandler) RetryRun(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	run, exists := h.runs[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "action run not found")
	}

	h.runCount++
	now := time.Now().UTC().Format(time.RFC3339)
	newRun := &models.ActionRun{
		ID:          "run-" + hex.EncodeToString(generateActionBytes(6)),
		WorkflowID:  run.WorkflowID,
		Name:        run.Name,
		Number:      h.runCount,
		Status:      "running",
		TriggeredBy: "retry",
		RunnerID:    run.RunnerID,
		Logs:        "[Retry] Re-dispatching workflow job\n[Worker] Picked up by runner\n",
		StartedAt:   now,
		DurationS:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	h.runs[newRun.ID] = newRun
	return utils.Success(c, "Action run retried", newRun)
}

func (h *ActionsHandler) CancelRun(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	run, exists := h.runs[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "action run not found")
	}

	run.Status = "cancelled"
	run.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	run.UpdatedAt = run.FinishedAt
	run.Logs += "\n[Abort] Run cancelled by operator"

	return utils.Success(c, "Action run cancelled", run)
}

func (h *ActionsHandler) DeleteWorkflow(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.workflows[id]; !exists {
		return echo.NewHTTPError(http.StatusNotFound, "workflow not found")
	}
	delete(h.workflows, id)
	return utils.Success(c, "Workflow deleted", map[string]string{"id": id})
}

func (h *ActionsHandler) DeleteRun(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.runs[id]; !exists {
		return echo.NewHTTPError(http.StatusNotFound, "action run not found")
	}
	delete(h.runs, id)
	return utils.Success(c, "Action run deleted", map[string]string{"id": id})
}

func generateActionBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
