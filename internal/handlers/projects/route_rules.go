package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type RouteRuleHandler struct {
	service *systemservices.RouteRuleService
	appRepo repositories.AppServiceRepository
}

func NewRouteRuleHandler(service *systemservices.RouteRuleService, appRepo repositories.AppServiceRepository) *RouteRuleHandler {
	return &RouteRuleHandler{service: service, appRepo: appRepo}
}

func (h *RouteRuleHandler) List(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if _, err := h.appRepo.GetByID(c.Request().Context(), serviceID); err != nil {
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
	if _, err := h.appRepo.GetByID(c.Request().Context(), serviceID); err != nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	var req models.CreateRouteRuleRequest
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

	if _, err := h.appRepo.GetByID(c.Request().Context(), serviceID); err != nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	existing, err := h.service.GetByID(c.Request().Context(), ruleID)
	if err != nil || existing == nil || existing.ServiceID != serviceID {
		return utils.Error(c, http.StatusNotFound, "rule not found")
	}

	var req models.UpdateRouteRuleRequest
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

	if _, err := h.appRepo.GetByID(c.Request().Context(), serviceID); err != nil {
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
