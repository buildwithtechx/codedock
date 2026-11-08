package environments

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ServiceChecker interface {
	CheckServiceExists(ctx context.Context, serviceID string) bool
}

type RouteRuleHandler struct {
	service        *RouteRuleService
	serviceChecker ServiceChecker
}

func NewRouteRuleHandler(service *RouteRuleService, checker ServiceChecker) *RouteRuleHandler {
	return &RouteRuleHandler{
		service:        service,
		serviceChecker: checker,
	}
}

func (h *RouteRuleHandler) List(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if h.serviceChecker != nil && !h.serviceChecker.CheckServiceExists(c.Request().Context(), serviceID) {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	rules, err := h.service.List(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", rules)
}

func (h *RouteRuleHandler) Create(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if h.serviceChecker != nil && !h.serviceChecker.CheckServiceExists(c.Request().Context(), serviceID) {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	var req CreateRouteRuleRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "name is required")
	}
	if req.RuleType == "" {
		return utils.Error(c, http.StatusBadRequest, "ruleType is required")
	}

	rule, err := h.service.Create(c.Request().Context(), serviceID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"success": true, "rule": rule})
}

func (h *RouteRuleHandler) Update(c echo.Context) error {
	serviceID := c.Param("serviceId")
	ruleID := c.Param("ruleId")

	if h.serviceChecker != nil && !h.serviceChecker.CheckServiceExists(c.Request().Context(), serviceID) {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	existing, err := h.service.GetByID(c.Request().Context(), ruleID)
	if err != nil || existing == nil || existing.ServiceID != serviceID {
		return utils.Error(c, http.StatusNotFound, "rule not found")
	}

	var req UpdateRouteRuleRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	rule, err := h.service.Update(c.Request().Context(), ruleID, req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Rule updated", rule)
}

func (h *RouteRuleHandler) Delete(c echo.Context) error {
	serviceID := c.Param("serviceId")
	ruleID := c.Param("ruleId")

	if h.serviceChecker != nil && !h.serviceChecker.CheckServiceExists(c.Request().Context(), serviceID) {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	existing, err := h.service.GetByID(c.Request().Context(), ruleID)
	if err != nil || existing == nil || existing.ServiceID != serviceID {
		return utils.Error(c, http.StatusNotFound, "rule not found")
	}

	if err := h.service.Delete(c.Request().Context(), ruleID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Rule deleted", nil)
}
