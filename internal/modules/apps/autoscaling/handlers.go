package autoscaling

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

func (h *AutoscalingHandler) getServiceID(c echo.Context) string {
	id := c.Param("serviceId")
	if id == "" {
		id = c.Param("id")
	}
	return id
}

func (h *AutoscalingHandler) Get(c echo.Context) error {
	serviceID := h.getServiceID(c)
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
	serviceID := h.getServiceID(c)
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
	return utils.Success(c, "Operation successful", policy)
}

func (h *AutoscalingHandler) GetServerlessCode(c echo.Context) error {
	serviceID := h.getServiceID(c)
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	code, err := h.service.GetServerlessCode(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "serverless code not found")
	}
	return utils.Success(c, "Operation successful", code)
}

func (h *AutoscalingHandler) SaveServerlessCode(c echo.Context) error {
	serviceID := h.getServiceID(c)
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var payload struct {
		Runtime     string `json:"runtime"`
		CodeContent string `json:"codeContent"`
	}
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	code, err := h.service.SaveServerlessCode(c.Request().Context(), serviceID, payload.Runtime, payload.CodeContent)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", code)
}
