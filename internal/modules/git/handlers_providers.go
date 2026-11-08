package git

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ProvidersHandler struct {
	service *ProviderService
}

func NewProvidersHandler(s *ProviderService) *ProvidersHandler {
	return &ProvidersHandler{service: s}
}

func (h *ProvidersHandler) Connect(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req GitConnectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	gp, err := h.service.ConnectProvider(c.Request().Context(), userID, &req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Created(c, "Created successfully", gp)
}

func (h *ProvidersHandler) Status(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	status, err := h.service.GetConnectedProviders(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", status)
}

func (h *ProvidersHandler) Disconnect(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	provider := c.Param("provider")
	if provider == "" {
		return utils.Error(c, http.StatusBadRequest, "missing provider parameter")
	}
	if err := h.service.DisconnectProvider(c.Request().Context(), userID, provider); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", DisconnectResponse{Status: "disconnected"})
}

func (h *ProvidersHandler) ListRepos(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	provider := c.QueryParam("provider")
	repos, err := h.service.ListRepositories(c.Request().Context(), userID, provider)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Operation successful", repos)
}
