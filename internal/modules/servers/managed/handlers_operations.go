package managed

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *ManagedHandler) ReviewProvision(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ReviewManagedProvisionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.CredentialID == "" {
		return utils.Error(c, http.StatusBadRequest, "credentialId is required")
	}
	if req.Region == "" {
		return utils.Error(c, http.StatusBadRequest, "region is required")
	}
	if req.ServerType == "" {
		return utils.Error(c, http.StatusBadRequest, "serverType is required")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewProvision(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed provision reviewed", review)
}

func (h *ManagedHandler) ReviewResize(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ReviewManagedResizeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.ServerType == "" {
		return utils.Error(c, http.StatusBadRequest, "serverType is required")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewResize(c.Request().Context(), userID, c.Param("id"), c.Param("serverId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed resize reviewed", review)
}

func (h *ManagedHandler) ReviewDelete(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewDelete(c.Request().Context(), userID, c.Param("id"), c.Param("serverId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed delete reviewed", review)
}

func (h *ManagedHandler) ApplyOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req struct {
		Confirmation string `json:"confirmation"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Confirmation == "" {
		return utils.Error(c, http.StatusBadRequest, "confirmation is required")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	if err := h.service.ApplyOperation(c.Request().Context(), userID, c.Param("operationId"), req.Confirmation); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed operation started", nil)
}

func (h *ManagedHandler) GetOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	op, err := h.service.GetOperation(c.Request().Context(), userID, c.Param("operationId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed operation retrieved", op)
}

func (h *ManagedHandler) StreamOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set(echo.HeaderConnection, "keep-alive")
	res.WriteHeader(http.StatusOK)
	res.Flush()

	return h.service.StreamOperation(c.Request().Context(), userID, c.Param("operationId"), func(op any) error {
		b, err := json.Marshal(op)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(res, "data: %s\n\n", b); err != nil {
			return err
		}
		res.Flush()
		return nil
	})
}
