package servers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type TakeoverHandler struct {
	service *TakeoverService
	repo    TakeoverRepository
}

func NewTakeoverHandler(service *TakeoverService, repo TakeoverRepository) *TakeoverHandler {
	return &TakeoverHandler{service: service, repo: repo}
}

func (h *TakeoverHandler) Scan(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req TakeoverScanRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Host == "" {
		return utils.Error(c, http.StatusBadRequest, "host is required")
	}
	if req.SSHUser == "" {
		return utils.Error(c, http.StatusBadRequest, "sshUser is required")
	}
	if req.SSHKey == "" {
		return utils.Error(c, http.StatusBadRequest, "sshKey is required")
	}

	run := &TakeoverRun{
		UserID:         userID,
		SourceHost:     req.Host,
		SourcePlatform: req.Platform,
		Status:         TakeoverStatusScanning,
	}
	if run.SourcePlatform == "" {
		run.SourcePlatform = TakeoverPlatformDocker
	}
	if err := h.repo.Create(c.Request().Context(), run); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create run record")
	}

	stack, err := h.service.Scan(c.Request().Context(), req)
	if err != nil {
		_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, TakeoverStatusFailed, err.Error())
		return utils.Error(c, http.StatusBadGateway, "scan failed: "+err.Error())
	}

	stackJSON, _ := json.Marshal(stack)
	if err := h.repo.UpdateDiscovered(c.Request().Context(), run.ID, string(stackJSON)); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save discovered stack: "+err.Error())
	}

	return utils.Success(c, "Scan complete", map[string]interface{}{
		"runId": run.ID,
		"stack": stack,
	})
}

func (h *TakeoverHandler) Adopt(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req TakeoverAdoptRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.RunID == "" {
		return utils.Error(c, http.StatusBadRequest, "runId is required")
	}
	if req.ProjectName == "" {
		return utils.Error(c, http.StatusBadRequest, "projectName is required")
	}
	if len(req.ServiceNames) == 0 {
		return utils.Error(c, http.StatusBadRequest, "select at least one service")
	}

	run, err := h.repo.GetByID(c.Request().Context(), req.RunID)
	if err != nil || run == nil || run.UserID != userID {
		return utils.Error(c, http.StatusNotFound, "run not found")
	}

	stack, err := DeserializeStack(run.DiscoveredJSON)
	if err != nil || stack == nil {
		return utils.Error(c, http.StatusUnprocessableEntity, "no scanned data found — run scan first")
	}

	_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, TakeoverStatusAdopting, "")

	projectIDs, err := h.service.Adopt(c.Request().Context(), req, stack, userID)
	if err != nil {
		_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, TakeoverStatusFailed, err.Error())
		return utils.Error(c, http.StatusInternalServerError, "adopt failed: "+err.Error())
	}

	if err := h.repo.UpdateAdopted(c.Request().Context(), run.ID, projectIDs); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to record adopted projects: "+err.Error())
	}

	return utils.Success(c, "Services adopted successfully", map[string]interface{}{
		"projectIds": projectIDs,
	})
}

func (h *TakeoverHandler) ListRuns(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	runs, err := h.repo.ListByUser(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", runs)
}

func (h *TakeoverHandler) GetRun(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	run, err := h.repo.GetByID(c.Request().Context(), c.Param("id"))
	if err != nil || run == nil || run.UserID != userID {
		return utils.Error(c, http.StatusNotFound, "run not found")
	}
	return utils.Success(c, "Operation successful", run)
}
