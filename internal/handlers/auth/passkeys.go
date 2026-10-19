package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/labstack/echo/v4"

	handlerutils "codedock/internal/handlers/utils"
	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	"codedock/internal/telemetry"
	"codedock/internal/utils"
)

type PasskeyHandler struct {
	passkeyService *authservices.PasskeyService
}

func NewPasskeyHandler(s *authservices.PasskeyService) *PasskeyHandler {
	return &PasskeyHandler{passkeyService: s}
}

type verifyRegistrationPayload struct {
	Name     string          `json:"name"`
	Response json.RawMessage `json:"response"`
}

type verifyAuthenticationPayload struct {
	TOTPCode string          `json:"totpCode"`
	Response json.RawMessage `json:"response"`
}

func (h *PasskeyHandler) BeginRegistration(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}

	options, err := h.passkeyService.BeginRegistration(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, options)
}

func (h *PasskeyHandler) FinishRegistration(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}

	rawBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to read request body")
	}

	var payload verifyRegistrationPayload
	_ = json.Unmarshal(rawBody, &payload)

	passkeyName := payload.Name
	if passkeyName == "" {
		passkeyName = c.QueryParam("name")
	}

	var credReader io.Reader = bytes.NewReader(rawBody)
	if len(payload.Response) > 0 {
		credReader = bytes.NewReader(payload.Response)
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(credReader)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid passkey creation data: "+err.Error())
	}

	result, err := h.passkeyService.FinishRegistration(
		c.Request().Context(),
		userID,
		passkeyName,
		parsedResponse,
		c.RealIP(),
	)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Passkey registered successfully", result)
}

func (h *PasskeyHandler) BeginAuthentication(c echo.Context) error {
	options, err := h.passkeyService.BeginAuthentication(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, options)
}

func (h *PasskeyHandler) FinishAuthentication(c echo.Context) error {
	rawBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to read request body")
	}

	var payload verifyAuthenticationPayload
	_ = json.Unmarshal(rawBody, &payload)

	totpCode := payload.TOTPCode
	if totpCode == "" {
		totpCode = c.Request().Header.Get("X-TOTP-Code")
	}
	if totpCode == "" {
		totpCode = c.QueryParam("totpCode")
	}

	var assertionReader io.Reader = bytes.NewReader(rawBody)
	if len(payload.Response) > 0 {
		assertionReader = bytes.NewReader(payload.Response)
	}

	parsedAssertion, err := protocol.ParseCredentialRequestResponseBody(assertionReader)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid passkey assertion data: "+err.Error())
	}

	u, token, refreshToken, err := h.passkeyService.FinishAuthentication(
		c.Request().Context(),
		parsedAssertion,
		totpCode,
		c.RealIP(),
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "2fa code required") {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success":           false,
				"error":             "2FA code required",
				"twoFactorRedirect": true,
			})
		}
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}

	handlerutils.SetAuthCookie(c, token)
	handlerutils.SetRefreshCookie(c, refreshToken)
	telemetry.Track(u.Email, "user_logged_in_passkey", map[string]any{
		"email": u.Email,
	})

	return utils.Success(c, "Login successful", models.AuthResponse{
		User:         u,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *PasskeyHandler) List(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}

	list, err := h.passkeyService.ListUserPasskeys(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Passkeys retrieved", list)
}

func (h *PasskeyHandler) Rename(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "passkey id required")
	}

	var req models.PasskeyRenameRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if err := h.passkeyService.RenamePasskey(c.Request().Context(), id, userID, req.Name, c.RealIP()); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Passkey renamed successfully", nil)
}

func (h *PasskeyHandler) Delete(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "passkey id required")
	}

	if err := h.passkeyService.DeletePasskey(c.Request().Context(), id, userID, c.RealIP()); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Passkey deleted successfully", nil)
}
