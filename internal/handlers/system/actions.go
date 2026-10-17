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
	h := &ActionsHandler{
		workflows: make(map[string]*models.ActionWorkflow),
		runs:      make(map[string]*models.ActionRun),
		runners:   make(map[string]*models.ActionRunner),
	}

	// Seed default runner
	defaultRunner := &models.ActionRunner{
		ID:            "runner-ubuntu-default",
		Name:          "Ubuntu 22.04 Default Runner",
		Mode:          "container",
		Image:         "catthehacker/ubuntu:act-22.04",
		Status:        "online",
		Labels:        []string{"ubuntu-latest", "docker", "x86_64"},
		Enabled:       true,
		LastHeartbeat: time.Now().UTC().Format(time.RFC3339),
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	}
	h.runners[defaultRunner.ID] = defaultRunner

	// Seed default workflow
	defaultWorkflow := &models.ActionWorkflow{
		ID:          "wf-build-test",
		Name:        "Build and Verify",
		Description: "Run automated tests and container image build",
		Trigger:     "push",
		YAML:        "name: Build and Verify\non: [push]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: npm test\n",
		Enabled:     true,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	h.workflows[defaultWorkflow.ID] = defaultWorkflow

	// Seed default run
	h.runCount = 1
	nowIso := time.Now().UTC().Format(time.RFC3339)
	defaultRun := &models.ActionRun{
		ID: "run-1", WorkflowID: defaultWorkflow.ID, Name: defaultWorkflow.Name,
		Number: 1, Status: "ready", TriggeredBy: "admin", RunnerID: defaultRunner.ID,
		Logs: "[Setup] Initializing runner environment\n[Test] Tests passed\n[Complete] Finished in 14s",
		Jobs: []models.ActionJob{{
			ID: "job-build", Name: "build", Status: "ready", RunnerID: defaultRunner.ID, DurationS: 14,
			Steps: []models.ActionStep{
				{ID: "step-1", Name: "Set up runner", Status: "ready", DurationS: 2},
				{ID: "step-2", Name: "Checkout code", Status: "ready", DurationS: 3},
				{ID: "step-3", Name: "Run tests", Status: "ready", DurationS: 9},
			},
		}},
		StartedAt: nowIso, FinishedAt: nowIso, DurationS: 14, CreatedAt: nowIso, UpdatedAt: nowIso,
	}
	h.runs[defaultRun.ID] = defaultRun

	return h
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
