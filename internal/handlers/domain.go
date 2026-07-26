package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type DomainHandler struct {
	envService     *services.EnvironmentService
	appService     *services.AppService
	projectService *services.ProjectService
}

func NewDomainHandler(s *services.EnvironmentService, app *services.AppService, proj *services.ProjectService) *DomainHandler {
	return &DomainHandler{
		envService:     s,
		appService:     app,
		projectService: proj,
	}
}

func (h *DomainHandler) hasAccess(ctx echo.Context, serviceID string) bool {
	userClaims, ok := ctx.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return false
	}
	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return true
	}

	if serviceID == "" {
		return false
	}

	svc, err := h.appService.GetAppService(ctx.Request().Context(), serviceID)
	if err != nil || svc == nil || svc.ProjectID == "" {
		return false
	}

	hasPerm := h.projectService.HasPermission(ctx.Request().Context(), svc.ProjectID, userClaims.UserID, userClaims.Role, "")
	return hasPerm
}

func (h *DomainHandler) ListByService(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id parameter")
	}

	if !h.hasAccess(c, serviceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	domains, err := h.envService.ListDomainsByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

func (h *DomainHandler) Create(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id parameter")
	}

	if !h.hasAccess(c, serviceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	var d models.DomainConfig
	if err := c.Bind(&d); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	d.ServiceID = serviceID
	if d.DomainName == "" {
		return utils.Error(c, http.StatusBadRequest, "domainName is required")
	}
	created, err := h.envService.CreateDomain(c.Request().Context(), &d)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *DomainHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	domain, err := h.envService.GetDomain(c.Request().Context(), id)
	if err == nil && domain != nil {
		if !h.hasAccess(c, domain.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	if err := h.envService.DeleteDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
