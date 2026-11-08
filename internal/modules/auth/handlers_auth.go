package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/telemetry"
	"codedock/internal/utils"
)

type AuthHandler struct {
	authService  *AuthService
	auditService *AuditService
}

func NewAuthHandler(s *AuthService, auditService *AuditService) *AuthHandler {
	return &AuthHandler{authService: s, auditService: auditService}
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var payload ForgotPasswordRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	originUrl := c.Request().Header.Get("Origin")
	if originUrl == "" {
		originUrl = c.Request().Header.Get("Referer")
	}
	if originUrl == "" {
		originUrl = "http://localhost:3000"
	}

	err := h.authService.ForgotPassword(c.Request().Context(), payload.Email, originUrl)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "If an account with that email exists, a password reset link has been sent.", nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var payload ResetPasswordRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	err := h.authService.ResetPassword(c.Request().Context(), payload.Token, payload.NewPassword)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Password reset successful", nil)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var payload RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	originUrl := c.Request().Header.Get("Origin")
	if originUrl == "" {
		originUrl = c.Request().Header.Get("Referer")
	}
	if originUrl == "" {
		originUrl = "http://localhost:3000"
	}

	u, token, refreshToken, err := h.authService.Register(c.Request().Context(), payload.Name, payload.Email, payload.Password, originUrl)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	setAuthCookie(c, token)
	setRefreshCookie(c, refreshToken)
	telemetry.Track(u.Email, "user_signed_up", map[string]any{
		"email": u.Email,
		"name":  u.Name,
	})
	h.auditService.LogAction(c.Request().Context(), AuditActionOpts{
		UserID:    u.ID,
		Action:    "auth.signup",
		Resource:  u.ID,
		IPAddress: c.RealIP(),
	})
	return utils.Success(c, "Registration successful", AuthResponse{
		User:         u,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var payload AuthRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, refreshToken, err := h.authService.Login(c.Request().Context(), payload.Email, payload.Password, payload.TOTPCode)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	setAuthCookie(c, token)
	setRefreshCookie(c, refreshToken)
	telemetry.Track(u.Email, "user_logged_in", map[string]any{
		"email": u.Email,
	})
	h.auditService.LogAction(c.Request().Context(), AuditActionOpts{
		UserID:    u.ID,
		Action:    "auth.login",
		Resource:  u.ID,
		IPAddress: c.RealIP(),
	})
	return utils.Success(c, "Login successful", AuthResponse{
		User:         u,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var payload RefreshRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	if payload.RefreshToken == "" {
		if cookie, err := c.Cookie("codedock_refresh_token"); err == nil {
			payload.RefreshToken = strings.TrimSpace(cookie.Value)
		}
	}
	if payload.RefreshToken == "" {
		return utils.Error(c, http.StatusBadRequest, "missing refresh token")
	}
	u, token, newRefreshToken, err := h.authService.RefreshToken(c.Request().Context(), payload.RefreshToken)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	setAuthCookie(c, token)
	setRefreshCookie(c, newRefreshToken)
	return utils.Success(c, "Token refreshed successfully", AuthResponse{
		User:         u,
		Token:        token,
		RefreshToken: newRefreshToken,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	clearAuthCookie(c)
	clearRefreshCookie(c)
	return utils.Success(c, "Logged out successfully", nil)
}

func (h *AuthHandler) AdminInviteUser(c echo.Context) error {
	var req AdminInviteUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "email is required")
	}
	if req.Role == "" {
		req.Role = string(UserRoleMember)
	}

	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		origin = c.Request().Header.Get("Referer")
	}
	if origin == "" {
		origin = "http://localhost:3000"
	}

	u, err := h.authService.InviteUser(c.Request().Context(), req.Email, UserRole(req.Role), origin)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	u.PasswordHash = ""
	return utils.Created(c, "User invited", u)
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	var payload ResendVerificationEmailRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	originUrl := c.Request().Header.Get("Origin")
	if originUrl == "" {
		originUrl = c.Request().Header.Get("Referer")
	}
	if originUrl == "" {
		originUrl = "http://localhost:3000"
	}

	err := h.authService.SendEmailVerification(c.Request().Context(), payload.Email, originUrl)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "If your email is not verified, a verification link has been sent.", nil)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var payload VerifyEmailRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	err := h.authService.VerifyEmail(c.Request().Context(), payload.Token)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Email verified successfully", nil)
}
