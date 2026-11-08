package git

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type InspectionHandler struct {
	service *ProviderService
}

func NewInspectionHandler(s *ProviderService) *InspectionHandler {
	return &InspectionHandler{service: s}
}

func (h *InspectionHandler) RepoTree(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	var req RepoTreeRequest
	_ = c.Bind(&req)
	if req.RepositoryURL == "" {
		req.RepositoryURL = c.QueryParam("repositoryUrl")
		req.Branch = c.QueryParam("branch")
	}
	if req.RepositoryURL == "" {
		return utils.Error(c, http.StatusBadRequest, "repositoryUrl is required")
	}
	tree, err := h.service.GetRepoTree(c.Request().Context(), userID, req.RepositoryURL, req.Branch)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Operation successful", tree)
}

func (h *InspectionHandler) InspectRepository(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req RepositoryInspectionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid inspection request")
	}
	result, err := h.service.InspectRepository(c.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(c, http.StatusUnprocessableEntity, err.Error())
	}
	return utils.Success(c, "Repository inspected", result)
}
