package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type OrganizationHandler struct {
	orgService *services.OrganizationService
}

func NewOrganizationHandler(orgService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

func (h *OrganizationHandler) Create(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req models.CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	org, err := h.orgService.CreateOrganization(c.Request().Context(), userClaims.UserID, req.Name)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, org)
}

func (h *OrganizationHandler) List(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	orgs, err := h.orgService.ListOrganizationsByUser(c.Request().Context(), userClaims.UserID)
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

func (h *OrganizationHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.orgService.DeleteOrganization(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
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
	orgID := c.Param("id")
	var req struct {
		Email      string                  `json:"email"`
		Permission models.MemberPermission `json:"permission"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	member, err := h.orgService.InviteMember(c.Request().Context(), orgID, req.Email, req.Permission)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, member)
}

func (h *OrganizationHandler) UpdateMember(c echo.Context) error {
	orgID := c.Param("id")
	userID := c.Param("userId")
	var req struct {
		Permission models.MemberPermission `json:"permission"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if err := h.orgService.UpdateMemberPermission(c.Request().Context(), orgID, userID, req.Permission); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *OrganizationHandler) RemoveMember(c echo.Context) error {
	memberID := c.Param("memberId")
	if err := h.orgService.RemoveMember(c.Request().Context(), memberID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
