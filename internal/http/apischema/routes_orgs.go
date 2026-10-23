package apischema

func operationSchema() Schema {
	return Obj("Durable operation",
		RF("id", ID("Operation ID")),
		F("projectId", Str("Project ID")),
		RF("kind", Str("Operation kind")),
		F("target", Str("Target description")),
		RF("status", StrEnum("Status", "pending", "running", "completed", "failed", "cancelled", "expired")),
		F("phase", Str("Current phase")),
		F("effects", Str("Effects summary")),
		F("error", Str("Error message")),
		F("logs", Str("Operation logs")),
		F("expiresAt", Int64("Expiry unix time")),
		F("updatedAt", Str("Last update")),
	)
}

func billingConfigSchema() Schema {
	return Obj("Billing configuration",
		F("publishableKey", Str("Stripe publishable key")),
		F("plans", Arr("Plans", Obj("Plan", RF("id", Str("Plan ID")), RF("name", Str("Plan name")), RF("price", Num("Monthly price")), F("priceId", Str("Stripe price ID")), F("features", Arr("Features", Str("Feature")))))),
	)
}

func notificationSettingsSchema() Schema {
	return Obj("Notification settings",
		F("discordWebhookUrl", Str("Discord webhook URL")),
		F("discordEnabled", Bool("Discord enabled")),
		F("slackWebhookUrl", Str("Slack webhook URL")),
		F("slackEnabled", Bool("Slack enabled")),
		F("telegramEnabled", Bool("Telegram enabled")),
		F("telegramChatId", Str("Telegram chat ID")),
		F("smtpHost", Str("SMTP host")),
		F("smtpPort", Int("SMTP port")),
		F("smtpUser", Str("SMTP username")),
		F("smtpFromAddress", Str("SMTP from address")),
		F("smtpEnabled", Bool("SMTP enabled")),
		F("resendEnabled", Bool("Resend enabled")),
		F("pushoverEnabled", Bool("Pushover enabled")),
		F("genericWebhookUrl", Str("Generic webhook URL")),
		F("genericWebhookEnabled", Bool("Generic webhook enabled")),
		F("notificationAlerts", Bool("Alert notifications")),
	)
}

func notificationSubscriptionSchema() Schema {
	return Obj("Notification subscription",
		RF("id", ID("Subscription ID")),
		RF("userId", Str("User ID")),
		RF("organizationId", Str("Organization ID")),
		RF("category", Str("Category")),
		RF("channels", Arr("Channels", Str("Channel"))),
		RF("enabled", Bool("Enabled")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}

func notificationDefaultSchema() Schema {
	return Obj("Notification default",
		RF("organizationId", Str("Organization ID")),
		RF("category", Str("Category")),
		RF("channels", Arr("Channels", Str("Channel"))),
		RF("enabled", Bool("Enabled")),
		F("updatedAt", Str("Last update")),
	)
}

func notificationPreferenceSchema() Schema {
	return Obj("Notification preference",
		RF("category", Str("Category")),
		RF("channels", Arr("Channels", Str("Channel"))),
		RF("enabled", Bool("Enabled")),
	)
}

func notificationDeliverySchema() Schema {
	return Obj("Notification delivery",
		RF("id", Str("Delivery ID")),
		RF("eventType", Str("Event type")),
		RF("channel", Str("Channel")),
		F("title", Str("Title")),
		RF("success", Bool("Delivered")),
		F("error", Str("Error")),
		F("createdAt", Str("Created at")),
	)
}

func organizationOperations() []Operation {
	orgs := "organizations"
	billing := "billing"
	ops := "operations"
	return []Operation{
		{Method: "POST", Path: "/api/organizations", Summary: "Create an organization", Tags: []string{orgs}, Auth: AuthUser, Raw: true, Code: 201, Request: JSON(Obj("Create organization", RF("name", Str("Organization name")))), Response: OrganizationSchema()},
		{Method: "GET", Path: "/api/organizations", Summary: "List my organizations", Tags: []string{orgs}, Auth: AuthUser, Raw: true, Response: Arr("Organizations", OrganizationSchema())},
		{Method: "GET", Path: "/api/organizations/:id", Summary: "Get an organization", Tags: []string{orgs}, Auth: AuthOrg, Raw: true, Response: OrganizationSchema()},
		{Method: "PUT", Path: "/api/organizations/:id", Summary: "Update an organization name", Tags: []string{orgs}, Auth: AuthOrgAdmin, Raw: true, Request: JSON(Obj("Update organization", RF("name", Str("Organization name")))), Response: OrganizationSchema()},
		{Method: "DELETE", Path: "/api/organizations/:id", Summary: "Delete an organization", Tags: []string{orgs}, Auth: AuthOrgOwner, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/organizations/:id/members", Summary: "List organization members", Tags: []string{orgs}, Auth: AuthOrg, Raw: true, Response: Arr("Members", MemberSchema())},
		{Method: "POST", Path: "/api/organizations/:id/members", Summary: "Invite a member", Tags: []string{orgs}, Auth: AuthOrgAdmin, Raw: true, Code: 201, Request: JSON(Obj("Invite member", RF("email", Str("Member email")), RF("permission", StrEnum("Permission", "admin", "member", "viewer")))), Response: MemberSchema()},
		{Method: "PUT", Path: "/api/organizations/:id/members/:userId", Summary: "Update a member permission", Tags: []string{orgs}, Auth: AuthOrgAdmin, Raw: true, Code: 204, Request: JSON(Obj("Update member", RF("permission", StrEnum("Permission", "admin", "member", "viewer")))), Response: EmptySchema()},
		{Method: "DELETE", Path: "/api/organizations/:id/members/:memberId", Summary: "Remove a member", Tags: []string{orgs}, Auth: AuthOrgAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/billing/config", Summary: "Get billing configuration and plans", Tags: []string{billing}, Auth: AuthUser, Raw: true, Response: billingConfigSchema()},
		{Method: "POST", Path: "/api/billing/checkout", Summary: "Create a checkout session", Tags: []string{billing}, Auth: AuthUser, Raw: true, Request: JSON(Obj("Checkout", RF("priceId", Str("Stripe price ID")), RF("successUrl", Str("Success URL")), RF("cancelUrl", Str("Cancel URL")))), Response: Obj("Checkout session", RF("url", Str("Checkout URL")))},
		{Method: "POST", Path: "/api/billing/webhook", Summary: "Handle a Stripe webhook", Tags: []string{billing}, Raw: true, Request: JSON(Any("Stripe event payload")), Response: EmptySchema()},
		{Method: "GET", Path: "/api/billing/pending", Summary: "Get active pending checkout", Tags: []string{billing}, Auth: AuthUser, Raw: true, Response: Obj("Pending checkout", F("pending", Any("Pending checkout session")))},
		{Method: "POST", Path: "/api/billing/checkouts/:id/cancel", Summary: "Cancel a pending checkout session", Tags: []string{billing}, Auth: AuthUser, Raw: true, Request: JSON(Obj("Cancel checkout", F("idempotencyKey", Str("Idempotency key")))), Response: Obj("Canceled checkout", RF("status", Str("Status")))},
		{Method: "GET", Path: "/api/operations", Summary: "List durable operations", Tags: []string{ops}, Auth: Scope("backup:read"), Response: Arr("Operations", operationSchema())},
		{Method: "GET", Path: "/api/operations/:operationId", Summary: "Get a durable operation", Tags: []string{ops}, Auth: Scope("backup:read"), Response: operationSchema()},
		{Method: "POST", Path: "/api/operations/:operationId/cancel", Summary: "Cancel a durable operation", Tags: []string{ops}, Auth: Scope("backup:write"), Response: EmptySchema()},
		{Method: "GET", Path: "/api/notifications", Summary: "Get notification settings", Tags: []string{orgs}, Auth: AuthUser, Response: notificationSettingsSchema()},
		{Method: "PUT", Path: "/api/notifications", Summary: "Update notification settings", Tags: []string{orgs}, Auth: AuthAdmin, Request: JSON(notificationSettingsSchema()), Response: notificationSettingsSchema()},
		{Method: "GET", Path: "/api/notifications/subscriptions", Summary: "List notification subscriptions", Tags: []string{orgs}, Auth: AuthUser, Response: Arr("Subscriptions", notificationSubscriptionSchema())},
		{Method: "PUT", Path: "/api/notifications/subscriptions", Summary: "Save a notification subscription", Tags: []string{orgs}, Auth: AuthUser, Request: JSON(notificationPreferenceSchema()), Response: notificationSubscriptionSchema()},
		{Method: "GET", Path: "/api/notifications/defaults", Summary: "List notification defaults", Tags: []string{orgs}, Auth: AuthUser, Response: Arr("Defaults", notificationDefaultSchema())},
		{Method: "GET", Path: "/api/notifications/deliveries", Summary: "List notification delivery attempts", Tags: []string{orgs}, Auth: AuthUser, Response: Arr("Deliveries", notificationDeliverySchema()), Query: []Param{QueryParamInt("limit", "Max entries")}},
		{Method: "PUT", Path: "/api/notifications/defaults", Summary: "Save a notification default", Tags: []string{orgs}, Auth: AuthAdmin, Request: JSON(notificationPreferenceSchema()), Response: notificationDefaultSchema()},
		{Method: "GET", Path: "/api/notifications/devices", Summary: "List registered device tokens", Tags: []string{orgs}, Auth: AuthUser, Response: Arr("Device tokens", deviceTokenSchema())},
		{Method: "POST", Path: "/api/notifications/devices", Summary: "Register a device token", Tags: []string{orgs}, Auth: AuthUser, Request: JSON(Obj("Register token", RF("token", Str("Device token")), F("platform", Str("Platform")))), Response: deviceTokenSchema()},
		{Method: "DELETE", Path: "/api/notifications/devices/:token", Summary: "Delete a registered device token", Tags: []string{orgs}, Auth: AuthUser, Response: EmptySchema()},
	}
}

func deviceTokenSchema() Schema {
	return Obj("Device token",
		RF("id", Str("Device token ID")),
		RF("userId", Str("User ID")),
		RF("token", Str("Push token")),
		RF("platform", Str("Platform")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}
