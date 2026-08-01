package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	projectservices "codedock.run/codedock/internal/services/projects"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type DomainHandler struct {
	envService     *projectservices.EnvironmentService
	appService     *projectservices.AppService
	projectService *projectservices.ProjectService
	settingsRepo   repositories.SettingsRepository
}

func NewDomainHandler(s *projectservices.EnvironmentService, app *projectservices.AppService, proj *projectservices.ProjectService, settingsRepo repositories.SettingsRepository) *DomainHandler {
	return &DomainHandler{
		envService:     s,
		appService:     app,
		projectService: proj,
		settingsRepo:   settingsRepo,
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

func (h *DomainHandler) ListAll(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	domains, err := h.envService.ListAllDomains(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if userClaims.Role != models.UserRoleAdmin && userClaims.Role != models.UserRoleOwner {
		filtered := make([]models.DomainConfig, 0, len(domains))
		for _, domain := range domains {
			if h.hasAccess(c, domain.ServiceID) {
				filtered = append(filtered, domain)
			}
		}
		domains = filtered
	}
	return utils.Success(c, "Operation successful", domains)
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
	if err != nil || domain == nil {
		return utils.Error(c, http.StatusNotFound, "domain not found")
	}
	if !h.hasAccess(c, domain.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if err := h.envService.DeleteDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *DomainHandler) Verify(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing domain id parameter")
	}

	domain, err := h.envService.GetDomain(c.Request().Context(), id)
	if err != nil || domain == nil {
		return utils.Error(c, http.StatusNotFound, "domain not found")
	}
	if !h.hasAccess(c, domain.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	serverIP := ""
	if h.settingsRepo != nil {
		if cfg, _ := h.settingsRepo.GetServerSettings(c.Request().Context()); cfg != nil {
			serverIP = cfg.PublicIPv4
		}
	}
	if serverIP == "" {
		return utils.Error(c, http.StatusInternalServerError, "server IP is not configured")
	}

	verified, resolvedIP, err := systemservices.VerifyDomain(c.Request().Context(), domain.DomainName, serverIP)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	status := "unresolved"
	message := "❌ Unresolved"
	if resolvedIP != "" {
		if verified {
			status = "resolves_to_server"
			message = "✅ Resolves to server IP"
		} else {
			status = "resolves_to_different_ip"
			message = "⚠️ Resolving to different IP (" + resolvedIP + ")"
		}
	}

	res := map[string]any{
		"domainId":   domain.ID,
		"domainName": domain.DomainName,
		"verified":   verified,
		"status":     status,
		"resolvedIp": resolvedIP,
		"serverIp":   serverIP,
		"message":    message,
	}

	return utils.Success(c, message, res)
}
