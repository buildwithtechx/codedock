package auth

import (
	"github.com/labstack/echo/v4"
)

type AuthRouteDeps struct {
	AuthHandler       *AuthHandler
	UserHandler       *UserHandler
	OAuthHandler      *OAuthHandler
	DeviceFlowHandler *DeviceFlowHandler
	PasskeyHandler    *PasskeyHandler
	BootstrapCSRF     echo.HandlerFunc
	AuthRateLimit     echo.MiddlewareFunc
	OTPRateLimit      echo.MiddlewareFunc
}

func RegisterAuthRoutes(apiGroup, authGroup *echo.Group, deps AuthRouteDeps) {
	if deps.AuthHandler == nil {
		deps.AuthHandler = &AuthHandler{}
	}
	if deps.UserHandler == nil {
		deps.UserHandler = &UserHandler{}
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
	authGroup.GET("/auth/me", deps.UserHandler.GetProfile)

	apiGroup.GET("/auth/oauth/providers/enabled", deps.OAuthHandler.ListEnabledProviders)
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

func RegisterUserRoutes(authGroup *echo.Group, userHandler *UserHandler, authHandler *AuthHandler, requireAdmin, otpRateLimit echo.MiddlewareFunc) {
	if userHandler == nil {
		userHandler = &UserHandler{}
	}
	if authHandler == nil {
		authHandler = &AuthHandler{}
	}
	authGroup.GET("/users", userHandler.ListUsers, requireAdmin)
	authGroup.POST("/users/invite", authHandler.AdminInviteUser, requireAdmin)
	authGroup.DELETE("/users/:id", userHandler.DeleteUser, requireAdmin)
	authGroup.GET("/profile", userHandler.GetProfile)
	authGroup.PUT("/profile", userHandler.UpdateProfile)
	authGroup.POST("/profile/email/request", userHandler.RequestEmailChange)
	authGroup.POST("/profile/email/verify", userHandler.VerifyEmailChange, otpRateLimit)
	authGroup.PUT("/profile/password", userHandler.ChangePassword)
	authGroup.GET("/profile/tokens", userHandler.ListPATs)
	authGroup.POST("/profile/tokens", userHandler.CreatePAT)
	authGroup.DELETE("/profile/tokens/:id", userHandler.DeletePAT)
}

func RegisterOrganizationRoutes(authGroup *echo.Group, orgHandler *OrganizationHandler, requireMember, requireAdmin, requireOwner echo.MiddlewareFunc) {
	if orgHandler == nil {
		orgHandler = &OrganizationHandler{}
	}
	authGroup.GET("/organizations", orgHandler.List)
	authGroup.POST("/organizations", orgHandler.Create)
	authGroup.GET("/organizations/:id", orgHandler.Get, requireMember)
	authGroup.PUT("/organizations/:id", orgHandler.Update, requireAdmin)
	authGroup.DELETE("/organizations/:id", orgHandler.Delete, requireOwner)
	authGroup.GET("/organizations/:id/members", orgHandler.ListMembers, requireMember)
	authGroup.POST("/organizations/:id/members", orgHandler.InviteMember, requireAdmin)
	authGroup.PUT("/organizations/:id/members/:userId", orgHandler.UpdateMember, requireAdmin)
	authGroup.DELETE("/organizations/:id/members/:memberId", orgHandler.RemoveMember, requireAdmin)
}

func RegisterAuditRoutes(authGroup *echo.Group, auditHandler *AuditLogHandler, requireAdmin echo.MiddlewareFunc) {
	if auditHandler == nil {
		auditHandler = &AuditLogHandler{}
	}
	authGroup.GET("/audit-logs", auditHandler.List, requireAdmin)
	authGroup.GET("/audit-logs/facets", auditHandler.Facets, requireAdmin)
	authGroup.GET("/audit-logs/settings", auditHandler.GetSettings, requireAdmin)
	authGroup.PUT("/audit-logs/settings", auditHandler.UpdateSettings, requireAdmin)
}
