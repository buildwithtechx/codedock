package volumes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type VolumeHandler struct {
	service VolumeService
	access  AccessChecker
}

func NewVolumeHandler(s VolumeService, ac AccessChecker) *VolumeHandler {
	return &VolumeHandler{service: s, access: ac}
}

func (h *VolumeHandler) List(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	volumes, err := h.service.ListVolumes(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", volumes)
}

func (h *VolumeHandler) Create(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var req CreateVolumeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.HostPath == "" || req.ContainerPath == "" {
		return utils.Error(c, http.StatusBadRequest, "hostPath and containerPath are required")
	}

	vol := &ServiceVolume{
		ServiceID:     serviceID,
		HostPath:      req.HostPath,
		ContainerPath: req.ContainerPath,
	}
	created, err := h.service.CreateVolume(c.Request().Context(), vol)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *VolumeHandler) Delete(c echo.Context) error {
	id := c.Param("volumeId")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing volume id")
	}

	if err := h.service.DeleteVolume(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *VolumeHandler) EmptyVolume(c echo.Context) error {
	serviceID := c.Param("id")
	volumeID := c.Param("volumeId")
	if serviceID == "" || volumeID == "" {
		return utils.Error(c, http.StatusBadRequest, "service id and volume id are required")
	}

	vol, prevPath, err := h.service.EmptyVolume(c.Request().Context(), serviceID, volumeID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Volume emptied successfully", map[string]any{
		"volume":           vol,
		"previousHostPath": prevPath,
	})
}
