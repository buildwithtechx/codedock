package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type ManagedHandler struct {
	managedService systemservices.ManagedService
}

func NewManagedHandler(managedService systemservices.ManagedService) *ManagedHandler {
	return &ManagedHandler{managedService: managedService}
}

func managedClaims(c echo.Context) (*models.UserClaims, error) {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return userClaims, nil
}

func (h *ManagedHandler) Catalog(c echo.Context) error {
	if _, err := managedClaims(c); err != nil {
		return err
	}
	return utils.Success(c, "Managed catalog retrieved", h.managedService.Catalog())
}

func (h *ManagedHandler) ListCredentials(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	credentials, err := h.managedService.ListCredentials(c.Request().Context(), claims.UserID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if credentials == nil {
		credentials = []*models.ManagedCredential{}
	}
	return utils.Success(c, "Managed credentials retrieved", credentials)
}

func (h *ManagedHandler) CreateCredential(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	var req models.CreateManagedCredentialRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	credential, err := h.managedService.CreateCredential(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credential stored", credential)
}

func (h *ManagedHandler) DeleteCredential(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	if err := h.managedService.DeleteCredential(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("credentialId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credential deleted", nil)
}

func (h *ManagedHandler) GetQuota(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	quota, err := h.managedService.GetQuota(c.Request().Context(), claims.UserID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed quota retrieved", quota)
}

func (h *ManagedHandler) SetQuota(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	var req models.ManagedQuotaRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	quota, err := h.managedService.SetQuota(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed quota updated", quota)
}

func (h *ManagedHandler) ReviewProvision(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	var req models.ReviewManagedProvisionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	review, err := h.managedService.ReviewProvision(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed provision reviewed", review)
}

func (h *ManagedHandler) ReviewResize(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	var req models.ReviewManagedResizeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	review, err := h.managedService.ReviewResize(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("serverId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed resize reviewed", review)
}

func (h *ManagedHandler) ReviewDelete(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	review, err := h.managedService.ReviewDelete(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("serverId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed delete reviewed", review)
}

func (h *ManagedHandler) ApplyOperation(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	var req models.OperationApply
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Confirmation == "" {
		return utils.Error(c, http.StatusBadRequest, "confirmation is required")
	}
	if err := h.managedService.ApplyOperation(c.Request().Context(), claims.UserID, c.Param("operationId"), req.Confirmation); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed operation started", nil)
}

func (h *ManagedHandler) RefreshServer(c echo.Context) error {
	claims, err := managedClaims(c)
	if err != nil {
		return err
	}
	server, err := h.managedService.RefreshServer(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("serverId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	redactManagedServer(server)
	return utils.Success(c, "Managed server refreshed", server)
}

func redactManagedServer(server *models.Server) {
	if server == nil {
		return
	}
	server.WorkerToken = "********"
	server.SSHPassword = "********"
	server.SSHPrivateKey = "********"
	if server.SSHKey != "" {
		server.SSHKey = "********"
	}
}
