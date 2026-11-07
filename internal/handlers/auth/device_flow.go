package auth

import (
	"errors"
	"net/http"
	"strings"

	handlerutils "codedock/internal/handlers/utils"
	authservices "codedock/internal/services/auth"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

type DeviceFlowHandler struct {
	deviceFlowService *authservices.DeviceFlowService
	userService       *authservices.UserService
	auditService      *authservices.AuditService
}

func NewDeviceFlowHandler(
	deviceFlowService *authservices.DeviceFlowService,
	userService *authservices.UserService,
	auditService *authservices.AuditService,
) *DeviceFlowHandler {
	return &DeviceFlowHandler{
		deviceFlowService: deviceFlowService,
		userService:       userService,
		auditService:      auditService,
	}
}

type CreateDeviceCodeRequest struct {
	ClientName string `json:"clientName"`
}

type PollDeviceTokenRequest struct {
	DeviceCode string `json:"deviceCode"`
}

type AuthorizeDeviceRequest struct {
	UserCode string `json:"userCode"`
	Action   string `json:"action"` // "approve" or "deny"
}

func (h *DeviceFlowHandler) CreateDeviceCode(c echo.Context) error {
	var payload CreateDeviceCodeRequest
	_ = c.Bind(&payload)

	clientName := strings.TrimSpace(payload.ClientName)
	if clientName == "" {
		clientName = "Codedock Client"
	}

	session, err := h.deviceFlowService.CreateSession(clientName)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to initiate device flow")
	}

	return utils.Success(c, "Device code generated", map[string]any{
		"deviceCode":              session.DeviceCode,
		"userCode":                session.UserCode,
		"verificationUri":         "/auth/device",
		"verificationUriComplete": "/auth/device?code=" + session.UserCode,
		"expiresIn":               900,
		"interval":                2,
	})
}

func (h *DeviceFlowHandler) PollDeviceToken(c echo.Context) error {
	var payload PollDeviceTokenRequest
	if err := c.Bind(&payload); err != nil || strings.TrimSpace(payload.DeviceCode) == "" {
		return utils.Error(c, http.StatusBadRequest, "deviceCode is required")
	}

	session, err := h.deviceFlowService.GetByDeviceCode(payload.DeviceCode)
	if err != nil {
		if errors.Is(err, authservices.ErrDeviceSessionExpired) {
			return utils.Error(c, http.StatusBadRequest, "device session has expired")
		}
		return utils.Error(c, http.StatusNotFound, "device session not found")
	}

	switch session.Status {
	case authservices.DeviceStatusApproved:
		return utils.Success(c, "Device authorized", map[string]any{
			"status": "approved",
			"token":  session.Token,
			"email":  session.UserEmail,
			"name":   session.UserName,
		})
	case authservices.DeviceStatusDenied:
		return utils.Success(c, "Authorization declined", map[string]any{
			"status":  "denied",
			"message": "authorization request was declined by the user",
		})
	default:
		return utils.Success(c, "Authorization pending", map[string]any{
			"status": "pending",
		})
	}
}

func (h *DeviceFlowHandler) VerifyUserCode(c echo.Context) error {
	code := strings.TrimSpace(c.QueryParam("code"))
	if code == "" {
		return utils.Error(c, http.StatusBadRequest, "code parameter is required")
	}

	session, err := h.deviceFlowService.GetByUserCode(code)
	if err != nil {
		if errors.Is(err, authservices.ErrDeviceSessionExpired) {
			return utils.Error(c, http.StatusBadRequest, "code has expired")
		}
		return utils.Error(c, http.StatusNotFound, "invalid code")
	}

	return utils.Success(c, "Code valid", map[string]any{
		"valid":      true,
		"userCode":   session.UserCode,
		"clientName": session.ClientName,
		"status":     session.Status,
	})
}

func (h *DeviceFlowHandler) AuthorizeDevice(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "authentication required")
	}

	var payload AuthorizeDeviceRequest
	if err := c.Bind(&payload); err != nil || strings.TrimSpace(payload.UserCode) == "" {
		return utils.Error(c, http.StatusBadRequest, "userCode is required")
	}

	session, err := h.deviceFlowService.GetByUserCode(payload.UserCode)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid or expired user code")
	}

	action := strings.ToLower(strings.TrimSpace(payload.Action))
	if action != "approve" && action != "deny" {
		return utils.Error(c, http.StatusBadRequest, "action must be 'approve' or 'deny'")
	}

	if action == "deny" {
		if _, err := h.deviceFlowService.Authorize(payload.UserCode, "", "", "", "", false); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Success(c, "Authorization declined", map[string]any{"status": "denied"})
	}

	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return utils.Error(c, http.StatusNotFound, "user not found")
	}

	tokenName := session.ClientName + " (Device Auth)"
	pat, rawToken, err := h.userService.CreatePAT(
		c.Request().Context(),
		userID,
		tokenName,
		"admin",
		"all",
		nil,
		nil,
	)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to mint token for device")
	}

	authorizedSession, err := h.deviceFlowService.Authorize(
		payload.UserCode,
		userID,
		u.Email,
		u.Name,
		rawToken,
		true,
	)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    userID,
		Action:    "token.device_authorized",
		Resource:  pat.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"clientName": session.ClientName},
	})

	return utils.Success(c, "Device authorized successfully", map[string]any{
		"status":     "approved",
		"clientName": authorizedSession.ClientName,
		"token":      rawToken,
		"email":      u.Email,
	})
}
