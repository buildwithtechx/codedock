package auth

import (
	"github.com/labstack/echo/v4"
)

type AuthRouteDeps struct {
	AuthHandler       *AuthHandler
	GetProfile        echo.HandlerFunc
	OAuthHandler      *OAuthHandler
	DeviceFlowHandler *DeviceFlowHandler
	PasskeyHandler    *PasskeyHandler
	BootstrapCSRF     echo.HandlerFunc
	AuthRateLimit     echo.MiddlewareFunc
	OTPRateLimit      echo.MiddlewareFunc
	AdminAuth         echo.MiddlewareFunc
}

func withMw(mws ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	var active []echo.MiddlewareFunc
	for _, m := range mws {
		if m != nil {
			active = append(active, m)
		}
	}
	return active
}

func RegisterAuthRoutes(apiGroup, authGroup *echo.Group, deps AuthRouteDeps) {
	if deps.AuthHandler == nil {
		deps.AuthHandler = &AuthHandler{}
	}
	if deps.OAuthHandler == nil {
		deps.OAuthHandler = &OAuthHandler{}
	}
	if deps.DeviceFlowHandler == nil {
		deps.DeviceFlowHandler = &DeviceFlowHandler{}
	}
	if deps.PasskeyHandler == nil {
		deps.PasskeyHandler = &PasskeyHandler{}
	}

	apiGroup.GET("/auth/csrf", deps.BootstrapCSRF)
	apiGroup.POST("/auth/signup", deps.AuthHandler.Register, deps.AuthRateLimit)
	apiGroup.POST("/auth/signin", deps.AuthHandler.Login, deps.AuthRateLimit)
	apiGroup.POST("/auth/refresh", deps.AuthHandler.Refresh)
	apiGroup.POST("/auth/forgot-password", deps.AuthHandler.ForgotPassword, deps.AuthRateLimit)
	apiGroup.POST("/auth/reset-password", deps.AuthHandler.ResetPassword, deps.AuthRateLimit)
	apiGroup.POST("/auth/email/resend", deps.AuthHandler.ResendVerificationEmail, deps.AuthRateLimit)
	apiGroup.POST("/auth/email/verify", deps.AuthHandler.VerifyEmail, deps.AuthRateLimit)
	apiGroup.POST("/auth/logout", deps.AuthHandler.Logout)
	if deps.GetProfile != nil {
		authGroup.GET("/auth/me", deps.GetProfile)
	}

	apiGroup.GET("/auth/oauth/providers/enabled", deps.OAuthHandler.ListEnabledProviders)
	apiGroup.GET("/settings/oauth/providers", deps.OAuthHandler.ListProviders, withMw(deps.AdminAuth)...)
	apiGroup.PUT("/settings/oauth/providers", deps.OAuthHandler.SaveProvider, withMw(deps.AdminAuth)...)
	apiGroup.GET("/auth/oauth/:provider", deps.OAuthHandler.OAuthRedirect)
	apiGroup.GET("/auth/oauth/:provider/callback", deps.OAuthHandler.OAuthCallback)
	authGroup.POST("/auth/2fa/setup", deps.OAuthHandler.Setup2FA)
	authGroup.POST("/auth/2fa/verify", deps.OAuthHandler.Verify2FA)
	authGroup.POST("/auth/2fa/disable", deps.OAuthHandler.Disable2FA, deps.OTPRateLimit)

	apiGroup.POST("/auth/device/code", deps.DeviceFlowHandler.CreateDeviceCode)
	apiGroup.POST("/auth/device/token", deps.DeviceFlowHandler.PollDeviceToken)
	apiGroup.GET("/auth/device/verify", deps.DeviceFlowHandler.VerifyUserCode)
	authGroup.POST("/auth/device/authorize", deps.DeviceFlowHandler.AuthorizeDevice)

	apiGroup.POST("/auth/passkeys/authenticate/options", deps.PasskeyHandler.BeginAuthentication, deps.AuthRateLimit)
	apiGroup.POST("/auth/passkeys/authenticate/verify", deps.PasskeyHandler.FinishAuthentication, deps.AuthRateLimit)
	authGroup.POST("/auth/passkeys/register/options", deps.PasskeyHandler.BeginRegistration)
	authGroup.POST("/auth/passkeys/register/verify", deps.PasskeyHandler.FinishRegistration)
	authGroup.GET("/auth/passkeys", deps.PasskeyHandler.List)
	authGroup.PUT("/auth/passkeys/:id", deps.PasskeyHandler.Rename)
	authGroup.DELETE("/auth/passkeys/:id", deps.PasskeyHandler.Delete)
}
