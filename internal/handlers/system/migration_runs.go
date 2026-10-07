package system

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func (h *MigrationLifecycleHandler) StartMigration(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.StartMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.migrations.StartMigration(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration started", result)
}

func (h *MigrationLifecycleHandler) StartProjectMove(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.StartProjectMoveRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.migrations.StartProjectMove(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Project move started", result)
}

func (h *MigrationLifecycleHandler) GetRun(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	detail, err := h.migrations.GetRun(c.Request().Context(), claims.UserID, c.Param("runId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Migration retrieved", detail)
}

func (h *MigrationLifecycleHandler) StreamRun(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
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
			detail, err := h.migrations.GetRun(c.Request().Context(), claims.UserID, runID)
			if err != nil {
				writeStreamEvent(writer, flush, "ERROR", err.Error())
				return nil
			}
			if detail.Run.Phase != lastPhase {
				lastPhase = detail.Run.Phase
				encoded, _ := json.Marshal(detail)
				writeStreamEvent(writer, flush, "STATUS", string(encoded))
			}
			if len(detail.Run.Logs) > lastLogs {
				writeStreamEvent(writer, flush, "LOGS", detail.Run.Logs[lastLogs:])
				lastLogs = len(detail.Run.Logs)
			}
			switch detail.Run.Status {
			case models.MigrationStatusCompleted, models.MigrationStatusFailed, models.MigrationStatusCancelled:
				encoded, _ := json.Marshal(detail)
				writeStreamEvent(writer, flush, "RESULT", string(encoded))
				return nil
			}
		}
	}
}

func (h *MigrationLifecycleHandler) Cutover(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.CutoverMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Confirmation == "" {
		return utils.Error(c, http.StatusBadRequest, "confirmation is required")
	}
	if err := h.migrations.CutoverRun(c.Request().Context(), claims.UserID, c.Param("runId"), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Cutover confirmed", nil)
}

func (h *MigrationLifecycleHandler) Cancel(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	if err := h.migrations.CancelRun(c.Request().Context(), claims.UserID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration cancel requested", nil)
}

func (h *MigrationLifecycleHandler) Respond(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.RespondMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if err := h.migrations.RespondRun(c.Request().Context(), claims.UserID, c.Param("runId"), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration prompt answered", nil)
}

func (h *MigrationLifecycleHandler) Resume(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.ResumeMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.migrations.ResumeRun(c.Request().Context(), claims.UserID, c.Param("runId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration resumed", result)
}

func (h *MigrationLifecycleHandler) CleanupTarget(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	if err := h.migrations.CleanupTargetData(c.Request().Context(), claims.UserID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Target data cleaned", nil)
}

func (h *MigrationLifecycleHandler) DeleteRun(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	if err := h.migrations.DeleteRunRecord(c.Request().Context(), claims.UserID, c.Param("runId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration record deleted", nil)
}

func (h *MigrationLifecycleHandler) Active(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	detail, err := h.migrations.ActiveRun(c.Request().Context(), claims.UserID, c.Param("id"), c.QueryParam("sourceId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Active migration retrieved", detail)
}

func (h *MigrationLifecycleHandler) ListRuns(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	runs, err := h.migrations.ListRuns(c.Request().Context(), claims.UserID, c.Param("id"), c.QueryParam("sourceId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if runs == nil {
		runs = []*models.MigrationRun{}
	}
	return utils.Success(c, "Migrations retrieved", runs)
}
