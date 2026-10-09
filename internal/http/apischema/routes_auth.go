package apischema

func oauthProviderSchema() Schema {
	return Obj("OAuth provider configuration",
		F("id", ID("Provider ID")),
		RF("providerName", Str("Provider name")),
		F("enabled", Bool("Whether the provider is enabled")),
		F("clientId", Str("OAuth client ID")),
		F("redirectUri", Str("Callback URI")),
		F("baseUrl", Str("Custom base URL")),
		F("tenant", Str("Tenant ID")),
	)
}

func auditLogSchema() Schema {
	return Obj("Audit log entry",
		RF("id", ID("Entry ID")),
		RF("userId", Str("Actor user ID")),
		RF("action", Str("Action performed")),
		F("resource", Str("Affected resource")),
		F("details", Str("Details")),
		F("ipAddress", Str("Client IP")),
		F("createdAt", Str("Creation time")),
		F("category", Str("Derived category")),
	)
}

func auditFacetsSchema() Schema {
	return Obj("Audit facets",
		RF("total", Int("Total entries")),
		RF("categories", Arr("Categories", Obj("Category count",
			RF("id", Str("Category ID")),
			RF("label", Str("Category label")),
			RF("description", Str("Category description")),
			RF("count", Int("Entry count")),
		))),
	)
}

func authOperations() []Operation {
	auth := "auth"
	return []Operation{
		{Method: "POST", Path: "/api/auth/signin", Summary: "Sign in with email and password", Tags: []string{auth}, Request: JSON(Obj("Sign-in credentials", RF("email", Str("Email address")), RF("password", Str("Password")), F("totpCode", Str("TOTP code when 2FA is enabled")))), Response: AuthResponseSchema()},
		{Method: "POST", Path: "/api/auth/signup", Summary: "Register a new account", Tags: []string{auth}, Request: JSON(Obj("Registration", RF("name", Str("Display name")), RF("email", Str("Email address")), RF("password", Str("Password")))), Response: AuthResponseSchema()},
		{Method: "POST", Path: "/api/auth/refresh", Summary: "Refresh an access token", Tags: []string{auth}, Request: JSON(Obj("Refresh", RF("refreshToken", Str("Refresh token")))), Response: AuthResponseSchema()},
		{Method: "POST", Path: "/api/auth/logout", Summary: "Invalidate the current session", Tags: []string{auth}, Request: JSON(Obj("Empty", F("refreshToken", Str("Refresh token")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/auth/forgot-password", Summary: "Request a password reset link", Tags: []string{auth}, Request: JSON(Obj("Forgot password", RF("email", Str("Email address")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/auth/reset-password", Summary: "Reset a password with a token", Tags: []string{auth}, Request: JSON(Obj("Reset password", RF("token", Str("Reset token")), RF("newPassword", Str("New password")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/auth/email/resend", Summary: "Resend the verification email", Tags: []string{auth}, Request: JSON(Obj("Resend verification", RF("email", Str("Email address")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/auth/email/verify", Summary: "Verify an email with a token", Tags: []string{auth}, Request: JSON(Obj("Verify email", RF("token", Str("Verification token")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/auth/csrf", Summary: "Bootstrap a CSRF token", Tags: []string{auth}, Response: Obj("CSRF token", RF("token", Str("CSRF token")))},
		{Method: "GET", Path: "/api/auth/oauth/:provider", Summary: "Redirect to an OAuth provider", Tags: []string{auth}, Raw: true, Code: 307, Response: EmptySchema()},
		{Method: "GET", Path: "/api/auth/oauth/:provider/callback", Summary: "Handle an OAuth callback", Tags: []string{auth}, Raw: true, Code: 307, Response: EmptySchema(), Query: []Param{QueryParam("code", "Authorization code"), QueryParam("state", "State token")}},
		{Method: "GET", Path: "/api/auth/oauth/providers/enabled", Summary: "List enabled OAuth providers", Tags: []string{auth}, Response: Arr("Providers", oauthProviderSchema())},
		{Method: "POST", Path: "/api/auth/2fa/setup", Summary: "Start two-factor setup", Tags: []string{auth}, Auth: AuthUser, Response: Obj("2FA setup", RF("qrCodeUri", Str("QR code URI")), RF("recoveryCodes", Arr("Recovery codes", Str("Code"))))},
		{Method: "POST", Path: "/api/auth/2fa/verify", Summary: "Verify and enable two-factor", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Verify 2FA", RF("passcode", Str("TOTP passcode")))), Response: StatusSchema("totp_enabled")},
		{Method: "POST", Path: "/api/auth/2fa/disable", Summary: "Disable two-factor", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Disable 2FA", RF("passcode", Str("TOTP passcode")))), Response: StatusSchema("totp_disabled")},
		{Method: "GET", Path: "/api/auth/me", Summary: "Get the current user", Tags: []string{auth}, Auth: AuthUser, Response: UserSchema()},
		{Method: "GET", Path: "/api/users", Summary: "List users", Tags: []string{auth}, Auth: AuthAdmin, Response: PaginatedSchema(Arr("Users", UserSchema())), Query: []Param{QueryParamInt("page", "Page number"), QueryParamInt("limit", "Page size")}},
		{Method: "DELETE", Path: "/api/users/:id", Summary: "Delete a user", Tags: []string{auth}, Auth: AuthAdmin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/users/invite", Summary: "Invite a user", Tags: []string{auth}, Auth: AuthAdmin, Code: 201, Request: JSON(Obj("Invite", RF("email", Str("Email address")), RF("role", Str("Instance role")))), Response: UserSchema()},
		{Method: "GET", Path: "/api/profile", Summary: "Get the user profile", Tags: []string{auth}, Auth: AuthUser, Response: UserSchema()},
		{Method: "PUT", Path: "/api/profile", Summary: "Update the user profile", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Profile update", F("name", Str("Display name")), F("role", Str("Role")))), Response: UserSchema()},
		{Method: "PUT", Path: "/api/profile/password", Summary: "Change the password", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Password change", RF("oldPassword", Str("Current password")), RF("newPassword", Str("New password")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/profile/email/request", Summary: "Request an email change", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Email change", RF("newEmail", Str("New email address")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/profile/email/verify", Summary: "Verify an email change", Tags: []string{auth}, Auth: AuthUser, Request: JSON(Obj("Verify email change", RF("otp", Str("One-time code")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/profile/tokens", Summary: "List personal access tokens", Tags: []string{auth}, Auth: AuthUser, Response: Arr("Tokens", PATSchema())},
		{Method: "POST", Path: "/api/profile/tokens", Summary: "Create a personal access token", Tags: []string{auth}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create token", RF("name", Str("Token name")), F("accessLevel", Str("Access level")), F("projectScope", Str("Project scope")), F("allowedProjects", Arr("Allowed projects", Str("Project ID"))), F("expiresAt", DateTime("Expiry time")))), Response: Obj("Created token", RF("token", PATSchema()), RF("plain", Str("One-time plain token")))},
		{Method: "DELETE", Path: "/api/profile/tokens/:id", Summary: "Delete a personal access token", Tags: []string{auth}, Auth: AuthUser, Response: EmptySchema()},
		{Method: "GET", Path: "/api/settings/oauth/providers", Summary: "List OAuth provider configs", Tags: []string{auth}, Auth: AuthAdmin, Response: Arr("Providers", oauthProviderSchema())},
		{Method: "PUT", Path: "/api/settings/oauth/providers", Summary: "Save an OAuth provider config", Tags: []string{auth}, Auth: AuthAdmin, Request: JSON(oauthProviderSchema()), Response: oauthProviderSchema()},
		{Method: "GET", Path: "/api/system/setup-status", Summary: "Check whether setup is required", Tags: []string{auth}, Response: Obj("Setup status", RF("setupRequired", Bool("Whether initial setup is required")))},
		{Method: "GET", Path: "/api/audit-logs", Summary: "List audit log entries", Tags: []string{auth}, Auth: AuthAdmin, Response: Arr("Entries", auditLogSchema()), Query: []Param{QueryParamInt("limit", "Max entries"), QueryParamInt("offset", "Offset"), QueryParam("category", "Category ID filter")}},
		{Method: "GET", Path: "/api/audit-logs/facets", Summary: "Count audit log entries per category", Tags: []string{auth}, Auth: AuthAdmin, Response: auditFacetsSchema()},
	}
}
