package managed

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ManagedService interface {
	Catalog() ManagedCatalog
	Offers(ctx context.Context, userID, orgID, credentialID string) (*ManagedOffers, error)
	GetOperation(ctx context.Context, userID, operationID string) (any, error)
	StreamOperation(ctx context.Context, userID, operationID string, emit func(any) error) error
	CreateCredential(ctx context.Context, userID, orgID string, req CreateManagedCredentialRequest) (*ManagedCredential, error)
	ListCredentials(ctx context.Context, userID, orgID string) ([]*ManagedCredential, error)
	DeleteCredential(ctx context.Context, userID, orgID, credentialID string) error
	GetQuota(ctx context.Context, userID, orgID string) (*ManagedQuota, error)
	SetQuota(ctx context.Context, userID, orgID string, req ManagedQuotaRequest) (*ManagedQuota, error)
	ReviewProvision(ctx context.Context, userID, orgID string, req ReviewManagedProvisionRequest) (*ManagedReview, error)
	ReviewResize(ctx context.Context, userID, orgID, serverID string, req ReviewManagedResizeRequest) (*ManagedReview, error)
	ReviewDelete(ctx context.Context, userID, orgID, serverID string) (*ManagedReview, error)
	ApplyOperation(ctx context.Context, userID, operationID, confirmation string) error
	RefreshServer(ctx context.Context, userID, orgID, serverID string) (any, error)
}

type ManagedHandler struct {
	service ManagedService
}

func NewManagedHandler(service ManagedService) *ManagedHandler {
	return &ManagedHandler{service: service}
}

func (h *ManagedHandler) Catalog(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Success(c, "Managed catalog retrieved", ManagedCatalog{
			Provider: ManagedProviderHetzner,
			Tiers:    []ManagedTier{},
			Regions:  []string{},
			Images:   []string{},
		})
	}
	return utils.Success(c, "Managed catalog retrieved", h.service.Catalog())
}

func (h *ManagedHandler) ListCredentials(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Success(c, "Managed credentials retrieved", []*ManagedCredential{})
	}
	creds, err := h.service.ListCredentials(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credentials retrieved", creds)
}

func (h *ManagedHandler) CreateCredential(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CreateManagedCredentialRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Label == "" {
		return utils.Error(c, http.StatusBadRequest, "label is required")
	}
	if req.Token == "" {
		return utils.Error(c, http.StatusBadRequest, "token is required")
	}
	if req.Provider == "" {
		req.Provider = ManagedProviderHetzner
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	cred, err := h.service.CreateCredential(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credential created", cred)
}

func (h *ManagedHandler) DeleteCredential(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	if err := h.service.DeleteCredential(c.Request().Context(), userID, c.Param("id"), c.Param("credentialId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credential deleted", nil)
}

func (h *ManagedHandler) Offers(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Success(c, "Managed offers retrieved", ManagedOffers{
			Provider:  ManagedProviderHetzner,
			Locations: []ManagedLocation{},
			Offers:    []ManagedOffer{},
			Images:    []string{},
		})
	}
	offers, err := h.service.Offers(c.Request().Context(), userID, c.Param("id"), c.Param("credentialId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed offers retrieved", offers)
}

func (h *ManagedHandler) GetQuota(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Success(c, "Managed quota retrieved", &ManagedQuota{
			OrganizationID: c.Param("id"),
			MaxServers:     5,
			MaxMemoryGB:    32,
		})
	}
	quota, err := h.service.GetQuota(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed quota retrieved", quota)
}

func (h *ManagedHandler) SetQuota(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ManagedQuotaRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	quota, err := h.service.SetQuota(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed quota updated", quota)
}

func (h *ManagedHandler) RefreshServer(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	server, err := h.service.RefreshServer(c.Request().Context(), userID, c.Param("id"), c.Param("serverId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed server refreshed", server)
}
