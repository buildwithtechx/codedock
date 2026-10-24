package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *ManagedHandler) Offers(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	offers, err := h.managedService.Offers(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("credentialId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed offers retrieved", offers)
}

func (h *ManagedHandler) GetOperation(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	op, err := h.managedService.GetOperation(c.Request().Context(), claims.UserID, c.Param("operationId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Managed operation observed", op)
}
