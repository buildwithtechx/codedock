package monitoring

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AttentionHandler struct {
	service AttentionService
}

func NewAttentionHandler(service AttentionService) *AttentionHandler {
	return &AttentionHandler{service: service}
}

func (h *AttentionHandler) List(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	issues, err := h.service.List(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Attention issues retrieved", issues)
}

func (h *AttentionHandler) Evaluate(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	result, err := h.service.Evaluate(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Evaluation complete", result)
}

func (h *AttentionHandler) ClearHistory(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	if err := h.service.ClearHistory(c.Request().Context(), userID, c.Param("id")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "History cleared", nil)
}

func (h *AttentionHandler) Acknowledge(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	if err := h.service.Acknowledge(c.Request().Context(), userID, c.Param("id"), c.Param("issueId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Issue acknowledged", nil)
}

func (h *AttentionHandler) Resolve(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	if err := h.service.Resolve(c.Request().Context(), userID, c.Param("id"), c.Param("issueId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Issue resolved", nil)
}

func (h *AttentionHandler) Act(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	var req AttentionActionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.service.Act(c.Request().Context(), userID, c.Param("id"), c.Param("issueId"), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Action executed", nil)
}
