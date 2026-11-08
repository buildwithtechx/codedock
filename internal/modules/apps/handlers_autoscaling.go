package apps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AutoscalingHandler struct {
	service *AutoscalingService
}

func NewAutoscalingHandler(s *AutoscalingService) *AutoscalingHandler {
	return &AutoscalingHandler{service: s}
}

func (h *AutoscalingHandler) Get(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	policy, err := h.service.GetPolicy(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", policy)
}

func (h *AutoscalingHandler) Update(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var policy AutoscalingPolicy
	if err := c.Bind(&policy); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	policy.ServiceID = serviceID

	if err := h.service.UpdatePolicy(c.Request().Context(), &policy); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Autoscaling policy updated", policy)
}

func (h *AutoscalingHandler) SaveServerlessCode(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var req SaveServerlessCodeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Runtime == "" {
		return utils.Error(c, http.StatusBadRequest, "runtime is required")
	}

	code, err := h.service.SaveServerlessCode(c.Request().Context(), serviceID, req.Runtime, req.CodeContent)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Serverless code saved", code)
}

func (h *AutoscalingHandler) GetServerlessCode(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	code, err := h.service.GetServerlessCode(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", code)
}
