package system

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (h *ManagedHandler) StreamOperation(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	err = h.managedService.StreamOperation(c.Request().Context(), claims.UserID, c.Param("operationId"), func(op *models.Operation) error {
		encoded, err := json.Marshal(op)
		if err != nil {
			return err
		}
		if models.ManagedOperationComplete(op.Status) {
			writeStreamEvent(writer, flush, "done", string(encoded))
			return nil
		}
		writeStreamEvent(writer, flush, "progress", string(encoded))
		return nil
	})
	if err != nil {
		writeStreamEvent(writer, flush, "error", err.Error())
	}
	return nil
}
