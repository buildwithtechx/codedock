package servers

import (
	"context"
	"encoding/json"
	"fmt"
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
	RefreshServer(ctx context.Context, userID, orgID, serverID string) (*Server, error)
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
	credentials, err := h.service.ListCredentials(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if credentials == nil {
		credentials = []*ManagedCredential{}
	}
	return utils.Success(c, "Managed credentials retrieved", credentials)
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
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	credential, err := h.service.CreateCredential(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed credential stored", credential)
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
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
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
			MaxServers:     defaultManagedMaxServers,
			MaxMemoryGB:    defaultManagedMaxMemoryGB,
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

func (h *ManagedHandler) ReviewProvision(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ReviewManagedProvisionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewProvision(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed provision reviewed", review)
}

func (h *ManagedHandler) ReviewResize(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ReviewManagedResizeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewResize(c.Request().Context(), userID, c.Param("id"), c.Param("serverId"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed resize reviewed", review)
}

func (h *ManagedHandler) ReviewDelete(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	review, err := h.service.ReviewDelete(c.Request().Context(), userID, c.Param("id"), c.Param("serverId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed delete reviewed", review)
}

func (h *ManagedHandler) ApplyOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req struct {
		Confirmation string `json:"confirmation"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Confirmation == "" {
		return utils.Error(c, http.StatusBadRequest, "confirmation is required")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	if err := h.service.ApplyOperation(c.Request().Context(), userID, c.Param("operationId"), req.Confirmation); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed operation started", nil)
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
	if server != nil {
		server.WorkerToken = "********"
		server.SSHPassword = "********"
		server.SSHPrivateKey = "********"
		if server.SSHKey != "" {
			server.SSHKey = "********"
		}
	}
	return utils.Success(c, "Managed server refreshed", server)
}

func (h *ManagedHandler) GetOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	op, err := h.service.GetOperation(c.Request().Context(), userID, c.Param("operationId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Managed operation retrieved", op)
}

func (h *ManagedHandler) StreamOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "managed service not available")
	}
	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set(echo.HeaderConnection, "keep-alive")
	res.WriteHeader(http.StatusOK)
	res.Flush()

	return h.service.StreamOperation(c.Request().Context(), userID, c.Param("operationId"), func(op any) error {
		b, err := json.Marshal(op)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(res, "data: %s\n\n", b); err != nil {
			return err
		}
		res.Flush()
		return nil
	})
}
