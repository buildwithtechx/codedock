package projects

import (
	"codedock.run/codedock/internal/models"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AutoscalingHandler struct {
	service *projectservices.AutoscalingService
}

func NewAutoscalingHandler(service *projectservices.AutoscalingService) *AutoscalingHandler {
	return &AutoscalingHandler{service: service}
}
func (h *AutoscalingHandler) Get(c echo.Context) error {
	p, err := h.service.Get(c.Request().Context(), c.Param("serviceId"))
	if err != nil {
		if utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, "service not found")
		}
		return utils.Error(c, http.StatusInternalServerError, "failed to load autoscaling policy")
	}
	return utils.Success(c, "Autoscaling policy", p)
}
func (h *AutoscalingHandler) Save(c echo.Context) error {
	var p models.AutoscalingPolicy
	if err := c.Bind(&p); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid autoscaling policy")
	}
	p.ServiceID = c.Param("serviceId")
	if err := h.service.Save(c.Request().Context(), &p); err != nil {
		if utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, "service not found")
		}
		if utils.IsValidation(err) {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, "failed to save autoscaling policy")
	}
	saved, err := h.service.Get(c.Request().Context(), p.ServiceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "policy saved but could not reload it")
	}
	return utils.Success(c, "Autoscaling policy saved", saved)
}
