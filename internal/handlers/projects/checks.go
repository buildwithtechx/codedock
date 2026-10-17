package projects

import (
	"net/http"
	"sync"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

type ChecksHandler struct {
	mu      sync.RWMutex
	configs map[string]*models.ProjectChecksConfig
}

func NewChecksHandler() *ChecksHandler {
	return &ChecksHandler{
		configs: make(map[string]*models.ProjectChecksConfig),
	}
}

func (h *ChecksHandler) GetAdmission(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "project id is required")
	}

	maxConcurrent := 4
	activeBuilds := 1

	status := &models.BuildAdmissionStatus{
		CanAdmit:             activeBuilds < maxConcurrent,
		ActiveBuilds:         activeBuilds,
		MaxConcurrentBuilds:  maxConcurrent,
		QueuePosition:        0,
		EstimatedWaitSeconds: 0,
		Reason:               "Capacity available for immediate execution",
	}

	return utils.Success(c, "Build admission status retrieved", status)
}

func (h *ChecksHandler) GetChecks(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "project id is required")
	}

	h.mu.RLock()
	cfg, exists := h.configs[projectID]
	h.mu.RUnlock()

	if !exists {
		cfg = &models.ProjectChecksConfig{
			ProjectID:            projectID,
			Enabled:              true,
			ReportCommitStatus:   true,
			RequirePassingChecks: false,
			TargetBranches:       []string{"main", "master"},
			ContextName:          "codedock/deploy",
			UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
		}
	}

	return utils.Success(c, "Project checks config retrieved", cfg)
}

func (h *ChecksHandler) UpdateChecks(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "project id is required")
	}

	var req models.UpdateChecksConfigRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid update checks payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	branches := req.TargetBranches
	if len(branches) == 0 {
		branches = []string{"main", "master"}
	}
	ctxName := req.ContextName
	if ctxName == "" {
		ctxName = "codedock/deploy"
	}

	cfg := &models.ProjectChecksConfig{
		ProjectID:            projectID,
		Enabled:              req.Enabled,
		ReportCommitStatus:   req.ReportCommitStatus,
		RequirePassingChecks: req.RequirePassingChecks,
		TargetBranches:       branches,
		ContextName:          ctxName,
		UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
	}

	h.configs[projectID] = cfg
	return utils.Success(c, "Project checks config updated", cfg)
}
