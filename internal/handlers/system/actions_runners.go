package system

import (
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

func (h *ActionsHandler) ListRunners(c echo.Context) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []*models.ActionRunner
	for _, runner := range h.runners {
		list = append(list, runner)
	}
	if list == nil {
		list = make([]*models.ActionRunner, 0)
	}
	return utils.Success(c, "Runners retrieved", list)
}

func (h *ActionsHandler) InspectDestination(c echo.Context) error {
	var req models.InspectDestinationRequest
	if err := c.Bind(&req); err != nil || req.ServerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "serverId is required")
	}

	dist := "ubuntu"
	ver := "22.04"
	capabilities := &models.ActionCapabilities{
		OS:                 "linux",
		Architecture:       "x64",
		Docker:             true,
		DockerArchitecture: "x64",
		DockerPlatforms:    []string{"linux/amd64"},
		Git:                true,
		Node:               true,
		Distribution:       &dist,
		Version:            &ver,
	}

	return utils.Success(c, "Destination capabilities inspected", capabilities)
}

func (h *ActionsHandler) CreateRunner(c echo.Context) error {
	var req models.CreateRunnerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid runner payload")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "runner name is required")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if req.ServerID != "" {
		for _, existing := range h.runners {
			if existing.ServerID != nil && *existing.ServerID == req.ServerID {
				return echo.NewHTTPError(http.StatusConflict, "This server already has an Actions runner. Edit that runner instead.")
			}
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	mode := req.Config.Mode
	if mode == "" {
		mode = req.Mode
	}
	if mode == "" {
		mode = "container"
	}
	req.Config.Mode = mode

	var imageStr string
	if req.Config.Image != nil && *req.Config.Image != "" {
		imageStr = *req.Config.Image
	} else if req.Image != "" {
		imageStr = req.Image
		req.Config.Image = &req.Image
	}

	if mode == "container" && imageStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Container runners need a working Docker engine and a runner image")
	}

	kind := req.Kind
	if kind == "" {
		if req.CloudPoolID != "" {
			kind = "cloud"
		} else {
			kind = "server"
		}
	}

	labels := []string{"self-hosted", "linux", "x64"}
	for _, l := range req.Config.Labels {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !contains(labels, trimmed) {
			labels = append(labels, trimmed)
		}
	}
	for _, l := range req.Labels {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !contains(labels, trimmed) {
			labels = append(labels, trimmed)
		}
	}

	var serverIDPtr *string
	if req.ServerID != "" {
		s := req.ServerID
		serverIDPtr = &s
	}
	var cloudPoolIDPtr *string
	if req.CloudPoolID != "" {
		cp := req.CloudPoolID
		cloudPoolIDPtr = &cp
	}

	dist := "ubuntu"
	ver := "22.04"
	cap := &models.ActionCapabilities{
		OS:                 "linux",
		Architecture:       "x64",
		Docker:             true,
		DockerArchitecture: "x64",
		DockerPlatforms:    []string{"linux/amd64"},
		Git:                true,
		Node:               true,
		Distribution:       &dist,
		Version:            &ver,
	}

	runner := &models.ActionRunner{
		ID:            "runner-" + hex.EncodeToString(generateActionBytes(6)),
		Name:          req.Name,
		ServerID:      serverIDPtr,
		CloudPoolID:   cloudPoolIDPtr,
		Kind:          kind,
		Mode:          mode,
		Image:         imageStr,
		Status:        "online",
		Labels:        labels,
		Config:        req.Config,
		Capabilities:  cap,
		Enabled:       req.Enabled,
		CheckedAt:     now,
		LastHeartbeat: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	h.runners[runner.ID] = runner
	return utils.Success(c, "Action runner registered", runner)
}

func (h *ActionsHandler) DeleteRunner(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.runners[id]; !exists {
		return echo.NewHTTPError(http.StatusNotFound, "runner not found")
	}
	delete(h.runners, id)
	return utils.Success(c, "Runner deleted", map[string]string{"id": id})
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
