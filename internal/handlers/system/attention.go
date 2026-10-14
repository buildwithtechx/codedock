package system

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/services/attention"
	authservices "codedock/internal/services/auth"
	"codedock/internal/utils"
)

type AttentionHandler struct {
	attention *attention.Service
	audit     *authservices.AuditService
}

func NewAttentionHandler(service *attention.Service, audit *authservices.AuditService) *AttentionHandler {
	return &AttentionHandler{attention: service, audit: audit}
}

func (h *AttentionHandler) List(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	orgID := c.Param("id")
	issues, err := h.attention.ListIssues(c.Request().Context(), claims.UserID, orgID, c.QueryParam("status"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if issues == nil {
		issues = []*models.AttentionIssue{}
	}
	canClear, err := h.attention.CanClearHistory(c.Request().Context(), claims.UserID, orgID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Attention feed retrieved", models.AttentionFeed{Issues: issues, CanClearHistory: canClear})
}

func (h *AttentionHandler) ClearHistory(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	orgID := c.Param("id")
	deleted, err := h.attention.ClearHistory(c.Request().Context(), claims.UserID, orgID)
	if err != nil {
		if errors.Is(err, attention.ErrClearHistoryForbidden) {
			return utils.Error(c, http.StatusForbidden, err.Error())
		}
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if h.audit != nil {
		h.audit.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "attention.history.clear",
			Resource:  orgID,
			IPAddress: c.RealIP(),
			Details:   map[string]any{"deleted": deleted},
		})
	}
	return utils.Success(c, "Monitoring history cleared", map[string]int64{"deleted": deleted})
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
