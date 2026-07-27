package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	handlerutils "codedock.run/codedock/internal/handlers/utils"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
)

type OAuthHandler struct {
	oauthService *authservices.OAuthService
}

func NewOAuthHandler(s *authservices.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: s}
}

type Verify2FARequest struct {
	Passcode string `json:"passcode"`
}

func (h *OAuthHandler) ListProviders(c echo.Context) error {
	providers, err := h.oauthService.ListProviders(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for i := range providers {
		providers[i].ClientSecret = ""
	}
	return utils.Success(c, "Operation successful", providers)
}

func (h *OAuthHandler) ListEnabledProviders(c echo.Context) error {
	providers, err := h.oauthService.ListEnabledProviders(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if providers == nil {
		providers = []models.OAuthProviderConfig{}
	}
	for i := range providers {
		providers[i].ClientSecret = ""
	}
	return utils.Success(c, "Operation successful", providers)
}

func (h *OAuthHandler) SaveProvider(c echo.Context) error {
	var p models.OAuthProviderConfig
	if err := c.Bind(&p); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.oauthService.SaveProvider(c.Request().Context(), &p); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", p)
}

func (h *OAuthHandler) OAuthRedirect(c echo.Context) error {
	providerName := strings.TrimPrefix(c.Request().URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}
	p, err := h.oauthService.GetProvider(c.Request().Context(), providerName)
	if err != nil || p == nil {
		return utils.Error(c, http.StatusNotFound, "oauth provider not found or not enabled: "+providerName)
	}
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to generate secure state token")
	}
	state := hex.EncodeToString(stateBytes)
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state_" + providerName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().TLS != nil || strings.HasPrefix(c.Request().Header.Get("X-Forwarded-Proto"), "https"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})
	authURL, err := authservices.GetAuthorizationURL(p, state)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *OAuthHandler) OAuthCallback(c echo.Context) error {
	providerName := strings.TrimPrefix(c.Request().URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")
	code := c.QueryParam("code")
	if code == "" {
		return utils.Error(c, http.StatusBadRequest, "missing authorization code parameter")
	}
	stateParam := c.QueryParam("state")
	cookie, err := c.Cookie("oauth_state_" + providerName)
	if err != nil || cookie.Value == "" || cookie.Value != stateParam {
		return utils.Error(c, http.StatusUnauthorized, "invalid or expired oauth state parameter")
	}
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state_" + providerName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	token, refreshToken, _, err := h.oauthService.HandleCallback(c.Request().Context(), providerName, code)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	handlerutils.SetAuthCookie(c, token)
	c.SetCookie(&http.Cookie{
		Name:     "codedock_refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *OAuthHandler) Setup2FA(c echo.Context) error {
	claims := handlerutils.ExtractClaims(c)
	if claims == nil || claims.UserID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	res, err := h.oauthService.Setup2FA(c.Request().Context(), claims.UserID, claims.Email)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", res)
}

func (h *OAuthHandler) Verify2FA(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload Verify2FARequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "missing 6-digit passcode")
	}
	if err := h.oauthService.Verify2FA(c.Request().Context(), userID, payload.Passcode); err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "totp_enabled"})
}

func (h *OAuthHandler) Disable2FA(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload Verify2FARequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "missing passcode")
	}
	if err := h.oauthService.Validate2FA(c.Request().Context(), userID, payload.Passcode); err != nil {
		if errors.Is(err, authservices.ErrInvalidPasscode) {
			return utils.Error(c, http.StatusUnauthorized, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, "failed to validate 2fa")
	}
	if err := h.oauthService.Disable2FA(c.Request().Context(), userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "totp_disabled"})
}
