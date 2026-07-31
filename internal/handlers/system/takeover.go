package system

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type TakeoverHandler struct {
	scanner *systemservices.TakeoverScanner
	adopter *systemservices.TakeoverAdopter
	repo    repositories.TakeoverRepository
}

func NewTakeoverHandler(
	scanner *systemservices.TakeoverScanner,
	adopter *systemservices.TakeoverAdopter,
	repo repositories.TakeoverRepository,
) *TakeoverHandler {
	return &TakeoverHandler{scanner: scanner, adopter: adopter, repo: repo}
}

func (h *TakeoverHandler) Scan(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req models.TakeoverScanRequest
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

	run := &models.TakeoverRun{
		UserID:         userClaims.UserID,
		SourceHost:     req.Host,
		SourcePlatform: req.Platform,
		Status:         models.TakeoverStatusScanning,
	}
	if run.SourcePlatform == "" {
		run.SourcePlatform = models.TakeoverPlatformDocker
	}
	if err := h.repo.Create(c.Request().Context(), run); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create run record")
	}

	stack, err := h.scanner.Scan(c.Request().Context(), req)
	if err != nil {
		_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, models.TakeoverStatusFailed, err.Error())
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
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req models.TakeoverAdoptRequest
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
	if err != nil || run == nil || run.UserID != userClaims.UserID {
		return utils.Error(c, http.StatusNotFound, "run not found")
	}

	stack, err := systemservices.DeserializeStack(run.DiscoveredJSON)
	if err != nil || stack == nil {
		return utils.Error(c, http.StatusUnprocessableEntity, "no scanned data found — run scan first")
	}

	_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, models.TakeoverStatusAdopting, "")

	projectIDs, err := h.adopter.Adopt(c.Request().Context(), req, stack, userClaims.UserID)
	if err != nil {
		_ = h.repo.UpdateStatus(c.Request().Context(), run.ID, models.TakeoverStatusFailed, err.Error())
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
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	runs, err := h.repo.ListByUser(c.Request().Context(), userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", runs)
}

func (h *TakeoverHandler) GetRun(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	run, err := h.repo.GetByID(c.Request().Context(), c.Param("id"))
	if err != nil || run == nil || run.UserID != userClaims.UserID {
		return utils.Error(c, http.StatusNotFound, "run not found")
	}
	return utils.Success(c, "Operation successful", run)
}
