package mail

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/utils"
)

type ProjectPermissionChecker interface {
	GetDefaultOrgID(ctx context.Context, userID string) (string, error)
	HasOrgPermission(ctx context.Context, orgID, userID, role, minPermission string) bool
}

type AuditLogger interface {
	LogAction(ctx context.Context, userID, action, resource, ip string, details map[string]string)
}

type Handler struct {
	mailService    *Service
	projectService ProjectPermissionChecker
	auditService   AuditLogger
}

func NewHandler(mailService *Service, projectService ProjectPermissionChecker, auditService AuditLogger) *Handler {
	if mailService != nil && mailService.repo != nil {
		if projectService == nil {
			if pc, ok := mailService.repo.(ProjectPermissionChecker); ok {
				projectService = pc
			}
		}
		if auditService == nil {
			if al, ok := mailService.repo.(AuditLogger); ok {
				auditService = al
			}
		}
	}
	return &Handler{
		mailService:    mailService,
		projectService: projectService,
		auditService:   auditService,
	}
}

func (h *Handler) ListServers(c echo.Context) error {
	orgID, err := h.resolveOrg(c, "", PermissionMember)
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
	orgID, err := h.resolveOrg(c, req.OrganizationID, PermissionAdmin)
	if err != nil {
		return err
	}
	req.OrganizationID = orgID
	server, err := h.mailService.CreateServer(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), utils.AuditActorID(c), "mail.server.create", server.ID, c.RealIP(), map[string]string{"domain": server.Domain, "maildUrl": server.MaildURL})
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
	server, err := h.scopedServer(c, PermissionAdmin)
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
		h.auditService.LogAction(c.Request().Context(), utils.AuditActorID(c), "mail.server.delete", server.ID, c.RealIP(), map[string]string{"domain": server.Domain})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) scopedServer(c echo.Context, minPermission string) (*MailServer, error) {
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

func (h *Handler) resolveOrg(c echo.Context, orgID string, minPermission string) (string, error) {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		_ = utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return "", echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if h.projectService != nil {
		if orgID == "" {
			defaultOrgID, err := h.projectService.GetDefaultOrgID(c.Request().Context(), user.UserID)
			if err != nil {
				_ = utils.Error(c, http.StatusInternalServerError, err.Error())
				return "", err
			}
			if defaultOrgID != "" {
				orgID = defaultOrgID
			}
		}
		if orgID == "" {
			_ = utils.Error(c, http.StatusBadRequest, "missing organization id")
			return "", echo.NewHTTPError(http.StatusBadRequest, "missing organization id")
		}
		if !h.projectService.HasOrgPermission(c.Request().Context(), orgID, user.UserID, string(user.Role), minPermission) {
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
