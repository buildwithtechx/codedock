package deployments

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type RegistryHandler struct {
	service *RegistryService
}

func NewRegistryHandler(service *RegistryService) *RegistryHandler {
	return &RegistryHandler{service: service}
}

func (h *RegistryHandler) List(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId parameter")
	}

	registries, err := h.service.ListRegistries(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Registries retrieved", registries)
}

func (h *RegistryHandler) Create(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId parameter")
	}

	var req CreateRegistryRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	reg, err := h.service.CreateRegistry(c.Request().Context(), projectID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Created(c, "Registry created", reg)
}

func (h *RegistryHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	if err := h.service.DeleteRegistry(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Registry deleted", MessageResponse{Message: "deleted"})
}
