package domains

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/dns"
	"codedock/internal/utils"
)

type DomainAccessChecker interface {
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
	AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error)
}

type ServiceProjectLookup interface {
	GetProjectID(ctx context.Context, serviceID string) (string, error)
	ListServiceIDsByProject(ctx context.Context, projectID string) ([]string, error)
}

type ServerSettingsProvider interface {
	GetPublicIPv4(ctx context.Context) (string, error)
}

type AuditLogger interface {
	Log(ctx context.Context, userID, action, resource, ip string, details map[string]string)
}

type DomainHandler struct {
	domainService *DomainService
	accessChecker DomainAccessChecker
	serviceLookup ServiceProjectLookup
	settings      ServerSettingsProvider
	auditLogger   AuditLogger
}

func NewDomainHandler(ds *DomainService, ac DomainAccessChecker, sl ServiceProjectLookup, sp ServerSettingsProvider, al AuditLogger) *DomainHandler {
	return &DomainHandler{
		domainService: ds,
		accessChecker: ac,
		serviceLookup: sl,
		settings:      sp,
		auditLogger:   al,
	}
}

func (h *DomainHandler) hasAccess(ctx echo.Context, serviceID string) bool {
	userID, role, ok := utils.ExtractUser(ctx)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	if serviceID == "" || h.serviceLookup == nil || h.accessChecker == nil {
		return false
	}
	projectID, err := h.serviceLookup.GetProjectID(ctx.Request().Context(), serviceID)
	if err != nil || projectID == "" {
		return false
	}
	return h.accessChecker.HasPermission(ctx.Request().Context(), projectID, userID, role, "")
}

func (h *DomainHandler) ListAll(c echo.Context) error {
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	domains, err := h.domainService.ListAllDomains(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if role != "admin" && role != "owner" && h.accessChecker != nil {
		accessibleServices, err := h.accessChecker.AccessibleServiceIDs(c.Request().Context(), userID, role)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		filtered := make([]DomainConfig, 0, len(domains))
		for _, domain := range domains {
			if _, allowed := accessibleServices[domain.ServiceID]; allowed {
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

	domains, err := h.domainService.ListDomainsByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

func (h *DomainHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}

	domains, err := h.domainService.ListAllDomains(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.serviceLookup == nil {
		return utils.Success(c, "Operation successful", []DomainConfig{})
	}

	serviceIDs, err := h.serviceLookup.ListServiceIDsByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	serviceMap := make(map[string]bool, len(serviceIDs))
	for _, id := range serviceIDs {
		serviceMap[id] = true
	}

	projectDomains := make([]DomainConfig, 0)
	for _, d := range domains {
		if serviceMap[d.ServiceID] {
			projectDomains = append(projectDomains, d)
		}
	}
	return utils.Success(c, "Operation successful", projectDomains)
}

func (h *DomainHandler) Create(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id parameter")
	}
	if !h.hasAccess(c, serviceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}
	_, role, ok := utils.ExtractUser(c)
	if !ok || (role != "admin" && role != "owner") {
		return utils.Error(c, http.StatusForbidden, "DNS provider-backed domains require administrator access")
	}

	var req CreateDomainRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.DomainName == "" {
		return utils.Error(c, http.StatusBadRequest, "domainName is required")
	}
	d := DomainConfig{
		ServiceID:  serviceID,
		DomainName: req.DomainName,
		RedirectTo: req.RedirectTo,
		PathPrefix: req.PathPrefix,
	}
	created, err := h.domainService.CreateDomain(c.Request().Context(), &d)
	if err != nil {
		if utils.IsValidation(err) {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger.Log(c.Request().Context(), utils.AuditActorID(c), "domain.create", created.ID, c.RealIP(), map[string]string{"domain": req.DomainName})
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *DomainHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	domain, err := h.domainService.GetDomain(c.Request().Context(), id)
	if err != nil {
		if !utils.IsNotFound(err) {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Error(c, http.StatusNotFound, "domain not found")
	}
	if domain == nil {
		return utils.Error(c, http.StatusNotFound, "domain not found")
	}
	if !h.hasAccess(c, domain.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if err := h.domainService.DeleteDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger.Log(c.Request().Context(), utils.AuditActorID(c), "domain.delete", id, c.RealIP(), map[string]string{"domain": domain.DomainName})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *DomainHandler) Verify(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing domain id parameter")
	}

	domain, err := h.domainService.GetDomain(c.Request().Context(), id)
	if err != nil {
		if utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, "domain not found")
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if domain == nil {
		return utils.Error(c, http.StatusNotFound, "domain not found")
	}
	if !h.hasAccess(c, domain.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	serverIP := ""
	if h.settings != nil {
		serverIP, _ = h.settings.GetPublicIPv4(c.Request().Context())
	}
	if serverIP == "" {
		return utils.Success(c, "Server IP is not configured", DomainVerifyResult{
			DomainID:   domain.ID,
			DomainName: domain.DomainName,
			Verified:   false,
			Status:     "server_ip_not_configured",
			Message:    "Server IP is not configured",
		})
	}

	verified, resolvedIP, err := dns.VerifyDomain(c.Request().Context(), domain.DomainName, serverIP)
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

	res := DomainVerifyResult{
		DomainID:   domain.ID,
		DomainName: domain.DomainName,
		Verified:   verified,
		Status:     status,
		ResolvedIP: resolvedIP,
		ServerIP:   serverIP,
		Message:    message,
	}

	return utils.Success(c, message, res)
}
