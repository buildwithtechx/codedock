package mail

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

type Handler struct {
	mailService    *Service
	projectService *projectservices.ProjectService
	auditService   *authservices.AuditService
}

func NewHandler(mailService *Service, projectService *projectservices.ProjectService, auditService *authservices.AuditService) *Handler {
	return &Handler{
		mailService:    mailService,
		projectService: projectService,
		auditService:   auditService,
	}
}

func (h *Handler) ListServers(c echo.Context) error {
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

func (h *Handler) CreateServer(c echo.Context) error {
	var req CreateMailServerRequest
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
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.server.create",
			Resource:  server.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"domain": server.Domain, "maildUrl": server.MaildURL},
		})
	}
	return utils.Created(c, "Mail server registered successfully", server)
}

func (h *Handler) GetServer(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	return utils.Success(c, "Mail server fetched successfully", server)
}

func (h *Handler) DeleteServer(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	if err := h.mailService.DeleteServer(c.Request().Context(), server.ID); err != nil {
		if errors.Is(err, ErrMailNotFound) {
			return utils.Error(c, http.StatusNotFound, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.server.delete",
			Resource:  server.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"domain": server.Domain},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListDomains(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	domains, err := h.mailService.ListDomains(c.Request().Context(), server.ID)
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Mail domains fetched successfully", domains)
}

func (h *Handler) CreateDomain(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req CreateMailDomainRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	domain, err := h.mailService.CreateDomain(c.Request().Context(), server.ID, req.Name)
	if err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.domain.create",
			Resource:  domain.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"domain": domain.Name, "server": server.ID},
		})
	}
	return utils.Created(c, "Mail domain created successfully", domain)
}

func (h *Handler) DeleteDomain(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	domainID := c.Param("domainId")
	if domainID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing domain id")
	}
	if err := h.mailService.DeleteDomain(c.Request().Context(), server.ID, domainID); err != nil {
		return h.handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListMailboxes(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	mailboxes, err := h.mailService.ListMailboxes(c.Request().Context(), server.ID, c.Param("domainId"))
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Mailboxes fetched successfully", mailboxes)
}

func (h *Handler) CreateMailbox(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req CreateMailboxRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	credentials, err := h.mailService.CreateMailbox(c.Request().Context(), server.ID, c.Param("domainId"), req.LocalPart)
	if err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.mailbox.create",
			Resource:  credentials.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"address": credentials.Address, "server": server.ID},
		})
	}
	return utils.Created(c, "Mailbox created successfully", credentials)
}

func (h *Handler) DeleteMailbox(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	if err := h.mailService.DeleteMailbox(c.Request().Context(), server.ID, c.Param("domainId"), c.Param("mailboxId")); err != nil {
		return h.handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) CreateDKIM(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req CreateDKIMRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	key, err := h.mailService.CreateDKIM(c.Request().Context(), server.ID, c.Param("domainId"), req.Selector, req.Algorithm)
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Created(c, "DKIM key created successfully", key)
}

func (h *Handler) GetDNS(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	response, err := h.mailService.GetDNS(c.Request().Context(), server.ID, c.Param("domainId"))
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "DNS records fetched successfully", response)
}

func (h *Handler) GetChecks(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	report, err := h.mailService.GetChecks(c.Request().Context(), server.ID, c.Param("domainId"), c.QueryParam("ip"))
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Mail checks fetched successfully", report)
}

func (h *Handler) ListCerts(c echo.Context) error {
	server, err := h.scopedServer(c, "")
	if err != nil || server == nil {
		return err
	}
	certs, err := h.mailService.ListCerts(c.Request().Context(), server.ID)
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Mail certs fetched successfully", certs)
}

func (h *Handler) ProbeCert(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req ProbeMailCertRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	cert, err := h.mailService.ProbeCert(c.Request().Context(), server.ID, req.Host, req.Port)
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Mail cert probe completed", cert)
}

func (h *Handler) GetFail2banStatus(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionMember)
	if err != nil || server == nil {
		return err
	}
	resp, err := h.mailService.GetFail2banStatus(c.Request().Context(), server.ID, c.QueryParam("jail"))
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Fail2ban status fetched successfully", resp)
}

func (h *Handler) BanIP(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req Fail2banActionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.IP == "" {
		return utils.Error(c, http.StatusBadRequest, "ip is required")
	}
	resp, err := h.mailService.BanIP(c.Request().Context(), server.ID, req.IP, req.Jail)
	if err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.fail2ban.ban",
			Resource:  server.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"ip": req.IP, "jail": req.Jail},
		})
	}
	return utils.Success(c, "IP banned successfully", resp)
}

func (h *Handler) UnbanIP(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	var req Fail2banActionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.IP == "" {
		return utils.Error(c, http.StatusBadRequest, "ip is required")
	}
	resp, err := h.mailService.UnbanIP(c.Request().Context(), server.ID, req.IP, req.Jail)
	if err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.fail2ban.unban",
			Resource:  server.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"ip": req.IP, "jail": req.Jail},
		})
	}
	return utils.Success(c, "IP unbanned successfully", resp)
}

func (h *Handler) ListAliases(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionMember)
	if err != nil || server == nil {
		return err
	}
	domainID := c.Param("domainId")
	if domainID == "" {
		return utils.Error(c, http.StatusBadRequest, "domainId is required")
	}
	aliases, err := h.mailService.ListAliases(c.Request().Context(), server.ID, domainID)
	if err != nil {
		return h.handleError(c, err)
	}
	return utils.Success(c, "Aliases fetched successfully", aliases)
}

func (h *Handler) CreateAlias(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	domainID := c.Param("domainId")
	if domainID == "" {
		return utils.Error(c, http.StatusBadRequest, "domainId is required")
	}
	var req CreateMailAliasRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	alias, err := h.mailService.CreateAlias(c.Request().Context(), server.ID, domainID, req)
	if err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.alias.create",
			Resource:  alias.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"source": alias.Source, "destination": alias.Destination},
		})
	}
	return utils.Success(c, "Alias created successfully", alias)
}

func (h *Handler) DeleteAlias(c echo.Context) error {
	server, err := h.scopedServer(c, models.MemberPermissionAdmin)
	if err != nil || server == nil {
		return err
	}
	domainID := c.Param("domainId")
	aliasID := c.Param("aliasId")
	if domainID == "" || aliasID == "" {
		return utils.Error(c, http.StatusBadRequest, "domainId and aliasId are required")
	}
	if err := h.mailService.DeleteAlias(c.Request().Context(), server.ID, domainID, aliasID); err != nil {
		return h.handleError(c, err)
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "mail.alias.delete",
			Resource:  aliasID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"domainId": domainID, "aliasId": aliasID},
		})
	}
	return utils.Success(c, "Alias deleted successfully", nil)
}

func (h *Handler) scopedServer(c echo.Context, minPermission models.MemberPermission) (*MailServer, error) {
	id := c.Param("id")
	if id == "" {
		_ = utils.Error(c, http.StatusBadRequest, "missing mail server id")
		return nil, echo.NewHTTPError(http.StatusBadRequest, "missing mail server id")
	}
	server, err := h.mailService.GetServer(c.Request().Context(), id)
	if err != nil {
		_ = h.handleError(c, err)
		return nil, err
	}
	if _, err := h.resolveOrg(c, server.OrganizationID, minPermission); err != nil {
		return nil, err
	}
	return server, nil
}

func (h *Handler) resolveOrg(c echo.Context, orgID string, minPermission models.MemberPermission) (string, error) {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		_ = utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return "", echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if h.projectService != nil {
		if orgID == "" {
			org, err := h.projectService.GetOrCreateDefaultOrganization(c.Request().Context(), user.UserID)
			if err != nil {
				_ = utils.Error(c, http.StatusInternalServerError, err.Error())
				return "", err
			}
			if org != nil {
				orgID = org.ID
			}
		}
		if orgID == "" {
			_ = utils.Error(c, http.StatusBadRequest, "missing organization id")
			return "", echo.NewHTTPError(http.StatusBadRequest, "missing organization id")
		}
		if !h.projectService.HasOrgPermission(c.Request().Context(), orgID, user.UserID, user.Role, minPermission) {
			_ = utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
			return "", echo.NewHTTPError(http.StatusForbidden, "insufficient permissions for this organization")
		}
	}
	return orgID, nil
}

func (h *Handler) handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrMailNotFound):
		return utils.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrMailExists):
		return utils.Error(c, http.StatusConflict, err.Error())
	default:
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
}
