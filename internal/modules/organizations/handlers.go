package organizations

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/audit"
	"codedock/internal/utils"
)

type OrganizationHandler struct {
	orgService   *OrganizationService
	auditService *audit.AuditService
}

func NewOrganizationHandler(orgService *OrganizationService, auditService *audit.AuditService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService, auditService: auditService}
}

func (h *OrganizationHandler) Create(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	org, err := h.orgService.CreateOrganization(c.Request().Context(), userID, req.Name)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    userID,
		Action:    "organization.create",
		Resource:  org.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"name": req.Name},
	})
	return c.JSON(http.StatusCreated, org)
}

func (h *OrganizationHandler) List(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	orgs, err := h.orgService.ListOrganizationsByUser(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, orgs)
}

func (h *OrganizationHandler) Get(c echo.Context) error {
	id := c.Param("id")
	org, err := h.orgService.GetOrganization(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if org == nil {
		return utils.Error(c, http.StatusNotFound, "organization not found")
	}
	return c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) Update(c echo.Context) error {
	id := c.Param("id")
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	org, err := h.orgService.UpdateOrganization(c.Request().Context(), userID, id, req.Name)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    userID,
		Action:    "organization.update",
		Resource:  id,
		Details:   map[string]string{"name": req.Name},
		IPAddress: c.RealIP(),
	})
	return c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.orgService.DeleteOrganization(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    auditActorID(c),
		Action:    "organization.delete",
		Resource:  id,
		IPAddress: c.RealIP(),
	})
	return c.NoContent(http.StatusNoContent)
}

func (h *OrganizationHandler) ListMembers(c echo.Context) error {
	orgID := c.Param("id")
	members, err := h.orgService.ListMembers(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, members)
}

func (h *OrganizationHandler) InviteMember(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	orgID := c.Param("id")
	var req InviteOrganizationMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	member, err := h.orgService.InviteMember(c.Request().Context(), userID, orgID, req.Email, req.Permission)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	member.InviteLink = "/invite/" + member.ID
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    userID,
		Action:    "member.invite",
		Resource:  orgID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"email": req.Email, "permission": string(req.Permission), "delivery": req.Delivery},
	})
	return c.JSON(http.StatusCreated, member)
}

func (h *OrganizationHandler) UpdateMember(c echo.Context) error {
	requesterID := extractUserID(c)
	if requesterID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	orgID := c.Param("id")
	userID := c.Param("userId")
	var req UpdateOrganizationMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if err := h.orgService.UpdateMemberPermission(c.Request().Context(), requesterID, orgID, userID, req.Permission); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    requesterID,
		Action:    "member.role",
		Resource:  orgID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"member": userID, "permission": string(req.Permission)},
	})
	return c.NoContent(http.StatusNoContent)
}

func (h *OrganizationHandler) RemoveMember(c echo.Context) error {
	memberID := c.Param("memberId")
	if err := h.orgService.RemoveMember(c.Request().Context(), memberID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), audit.AuditActionOpts{
		UserID:    auditActorID(c),
		Action:    "member.remove",
		Resource:  c.Param("id"),
		IPAddress: c.RealIP(),
		Details:   map[string]string{"member": memberID},
	})
	return c.NoContent(http.StatusNoContent)
}
