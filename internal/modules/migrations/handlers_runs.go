package migrations

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *LifecycleHandler) StartMigration(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req StartMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.service.StartMigration(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration started", result)
}

func (h *LifecycleHandler) StartProjectMove(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req StartProjectMoveRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.service.StartProjectMove(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Project move started", result)
}

func (h *LifecycleHandler) GetRun(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	detail, err := h.service.GetRun(c.Request().Context(), userID, c.Param("runId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Migration retrieved", detail)
}

func (h *LifecycleHandler) StreamRun(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	runID := c.Param("runId")
	lastPhase := ""
	lastLogs := 0
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			detail, err := h.service.GetRun(c.Request().Context(), userID, runID)
			if err != nil {
				writeStreamEvent(writer, flush, "ERROR", err.Error())
				return nil
			}
			runMap, ok := detail.(map[string]any)
			if !ok {
				encoded, _ := json.Marshal(detail)
				writeStreamEvent(writer, flush, "STATUS", string(encoded))
				continue
			}
			runObj, _ := runMap["run"].(map[string]any)
			phase, _ := runObj["phase"].(string)
			status, _ := runObj["status"].(string)
			logs, _ := runObj["logs"].(string)

			if phase != lastPhase {
				lastPhase = phase
				encoded, _ := json.Marshal(detail)
				writeStreamEvent(writer, flush, "STATUS", string(encoded))
			}
			if len(logs) > lastLogs {
				writeStreamEvent(writer, flush, "LOGS", logs[lastLogs:])
				lastLogs = len(logs)
			}
			switch MigrationStatus(status) {
			case MigrationStatusCompleted, MigrationStatusFailed, MigrationStatusCancelled:
				encoded, _ := json.Marshal(detail)
				writeStreamEvent(writer, flush, "RESULT", string(encoded))
				return nil
			}
		}
	}
}

func (h *LifecycleHandler) Cutover(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CutoverMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Confirmation == "" {
		return utils.Error(c, http.StatusBadRequest, "confirmation is required")
	}
	if _, err := h.service.Cutover(c.Request().Context(), userID, c.Param("runId"), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Cutover confirmed", nil)
}

func (h *LifecycleHandler) Cancel(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if err := h.service.Cancel(c.Request().Context(), userID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration cancel requested", nil)
}

func (h *LifecycleHandler) Respond(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req RespondMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if _, err := h.service.Respond(c.Request().Context(), userID, c.Param("runId"), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration prompt answered", nil)
}

func (h *LifecycleHandler) Resume(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ResumeMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.service.Resume(c.Request().Context(), userID, c.Param("runId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration resumed", result)
}

func (h *LifecycleHandler) CleanupTarget(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if err := h.service.CleanupTarget(c.Request().Context(), userID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Target data cleaned", nil)
}

func (h *LifecycleHandler) DeleteRun(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if err := h.service.DeleteRun(c.Request().Context(), userID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration record deleted", nil)
}

func (h *LifecycleHandler) Active(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	detail, err := h.service.Active(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Active migration retrieved", detail)
}

func (h *LifecycleHandler) ListRuns(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	runs, err := h.service.ListRuns(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if runs == nil {
		runs = []*MigrationRun{}
	}
	return utils.Success(c, "Migrations retrieved", runs)
}
