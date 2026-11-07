package mail

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	authservices "codedock/internal/services/auth"
	mailservices "codedock/internal/services/mail"
	projectservices "codedock/internal/services/projects"

	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
)

type MailHandler struct {
	mailService    *mailservices.MailService
	projectService *projectservices.ProjectService
	auditService   *authservices.AuditService
}

func NewMailHandler(mailService *mailservices.MailService, projectService *projectservices.ProjectService, auditService *authservices.AuditService) *MailHandler {
	return &MailHandler{mailService: mailService, projectService: projectService, auditService: auditService}
}

func (h *MailHandler) ListServers(c echo.Context) error {
	orgID, err := h.resolveOrg(c, "", models.MemberPermissionMember)
	if err != nil {
		return err
	}
	servers, err := h.mailService.ListServers(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Mail servers fetched successfully", servers)
}

func (h *MailHandler) CreateServer(c echo.Context) error {
	var req models.CreateMailServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	orgID, err := h.resolveOrg(c, req.OrganizationID, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	req.OrganizationID = orgID
	server, err := h.mailService.CreateServer(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "mail.server.create",
		Resource:  server.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"domain": server.Domain, "maildUrl": server.MaildURL},
	})
	return utils.Created(c, "Mail server registered successfully", server)
}

func (h *MailHandler) GetServer(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	return utils.Success(c, "Mail server fetched successfully", server)
}

func (h *MailHandler) DeleteServer(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	if err := h.mailService.DeleteServer(c.Request().Context(), server.ID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "mail.server.delete",
		Resource:  server.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"domain": server.Domain},
	})
	return c.NoContent(http.StatusNoContent)
}

func (h *MailHandler) ListDomains(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	domains, err := h.mailService.ListDomains(c.Request().Context(), server.ID)
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "Mail domains fetched successfully", domains)
}

func (h *MailHandler) CreateDomain(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	var req models.CreateMailDomainRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	domain, err := h.mailService.CreateDomain(c.Request().Context(), server.ID, req.Name)
	if err != nil {
		return h.mapMailError(c, err)
	}
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "mail.domain.create",
		Resource:  domain.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"domain": domain.Name, "server": server.ID},
	})
	return utils.Created(c, "Mail domain created successfully", domain)
}

func (h *MailHandler) DeleteDomain(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	domainID := c.Param("domainId")
	if domainID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing domain id")
	}
	if err := h.mailService.DeleteDomain(c.Request().Context(), server.ID, domainID); err != nil {
		return h.mapMailError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *MailHandler) ListMailboxes(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	mailboxes, err := h.mailService.ListMailboxes(c.Request().Context(), server.ID, c.Param("domainId"))
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "Mailboxes fetched successfully", mailboxes)
}

func (h *MailHandler) CreateMailbox(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	var req models.CreateMailboxRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	credentials, err := h.mailService.CreateMailbox(c.Request().Context(), server.ID, c.Param("domainId"), req.LocalPart)
	if err != nil {
		return h.mapMailError(c, err)
	}
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "mail.mailbox.create",
		Resource:  credentials.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"address": credentials.Address, "server": server.ID},
	})
	return utils.Created(c, "Mailbox created successfully", credentials)
}

func (h *MailHandler) DeleteMailbox(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	if err := h.mailService.DeleteMailbox(c.Request().Context(), server.ID, c.Param("domainId"), c.Param("mailboxId")); err != nil {
		return h.mapMailError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *MailHandler) CreateDKIM(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	var req models.CreateDKIMRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	key, err := h.mailService.CreateDKIM(c.Request().Context(), server.ID, c.Param("domainId"), req.Selector, req.Algorithm)
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Created(c, "DKIM key created successfully", key)
}

func (h *MailHandler) GetDNS(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	response, err := h.mailService.GetDNS(c.Request().Context(), server.ID, c.Param("domainId"))
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "DNS records fetched successfully", response)
}

func (h *MailHandler) GetChecks(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	report, err := h.mailService.GetChecks(c.Request().Context(), server.ID, c.Param("domainId"), c.QueryParam("ip"))
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "Mail checks fetched successfully", report)
}

func (h *MailHandler) ListCerts(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil {
		return err
	}
	certs, err := h.mailService.ListCerts(c.Request().Context(), server.ID)
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "Certificates fetched successfully", certs)
}

func (h *MailHandler) ProbeCert(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil {
		return err
	}
	var req models.ProbeMailCertRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	cert, err := h.mailService.ProbeCert(c.Request().Context(), server.ID, req.Host, req.Port)
	if err != nil {
		return h.mapMailError(c, err)
	}
	return utils.Success(c, "Certificate probed successfully", cert)
}

func (h *MailHandler) scopedServer(c echo.Context, minPermission models.MemberPermission) (*models.MailServer, error) {
	id := c.Param("id")
	if id == "" {
		return nil, utils.Error(c, http.StatusBadRequest, "missing mail server id")
	}
	server, err := h.mailService.GetServer(c.Request().Context(), id)
	if err != nil {
		return nil, h.mapMailError(c, err)
	}
	if _, err := h.resolveOrg(c, server.OrganizationID, minPermission); err != nil {
		return nil, err
	}
	return server, nil
}

func (h *MailHandler) resolveOrg(c echo.Context, orgID string, minPermission models.MemberPermission) (string, error) {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return "", utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if orgID == "" {
		org, err := h.projectService.GetOrCreateDefaultOrganization(c.Request().Context(), user.UserID)
		if err != nil {
			return "", utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if org != nil {
			orgID = org.ID
		}
	}
	if orgID == "" {
		return "", utils.Error(c, http.StatusBadRequest, "missing organization id")
	}
	if !h.projectService.HasOrgPermission(c.Request().Context(), orgID, user.UserID, user.Role, minPermission) {
		return "", utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
	}
	return orgID, nil
}

func (h *MailHandler) mapMailError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, mailservices.ErrMailNotFound):
		return utils.Error(c, http.StatusNotFound, "not found")
	case errors.Is(err, mailservices.ErrMailExists):
		return utils.Error(c, http.StatusConflict, "already exists")
	default:
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
}
