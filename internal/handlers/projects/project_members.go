package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type addMemberPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *ProjectHandler) ListMembers(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	members, err := h.projectService.ListProjectMembers(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Members fetched successfully", members)
}

func (h *ProjectHandler) AddMember(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	var req addMemberPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "email is required")
	}
	role := models.MemberPermission(req.Role)
	if role == "" {
		role = models.MemberPermissionMember
	}
	member, err := h.projectService.AddProjectMember(c.Request().Context(), projectID, req.Email, role)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Member added successfully", member)
}

func (h *ProjectHandler) RemoveMember(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	memberID := c.Param("memberId")
	if err := h.projectService.RemoveProjectMember(c.Request().Context(), projectID, memberID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
