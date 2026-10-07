package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/attention"
	"codedock.run/codedock/internal/utils"
)

type AttentionHandler struct {
	attention *attention.Service
}

func NewAttentionHandler(service *attention.Service) *AttentionHandler {
	return &AttentionHandler{attention: service}
}

func (h *AttentionHandler) List(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	issues, err := h.attention.ListIssues(c.Request().Context(), claims.UserID, c.Param("id"), c.QueryParam("status"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if issues == nil {
		issues = []*models.AttentionIssue{}
	}
	return utils.Success(c, "Attention feed retrieved", issues)
}

func (h *AttentionHandler) Evaluate(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.attention.Evaluate(c.Request().Context(), claims.UserID, c.Param("id")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Attention evaluated", nil)
}

func (h *AttentionHandler) Acknowledge(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.attention.Acknowledge(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("issueId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Issue acknowledged", nil)
}

func (h *AttentionHandler) Resolve(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.attention.Resolve(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("issueId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Issue resolved", nil)
}

func (h *AttentionHandler) Act(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	var req models.AttentionActionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	outcome, err := h.attention.Act(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("issueId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Remediation executed", map[string]string{"outcome": outcome})
}
