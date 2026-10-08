package apischema

func managedTierSchema() Schema {
	return Obj("Managed tier",
		RF("name", Str("Tier name")),
		RF("cores", Int("vCPUs")),
		RF("memoryGB", Int("Memory GB")),
		RF("diskGB", Int("Disk GB")),
		RF("monthlyPrice", Num("Monthly price")),
	)
}

func managedCatalogSchema() Schema {
	return Obj("Managed catalog",
		RF("provider", Str("Provider")),
		RF("tiers", Arr("Tiers", managedTierSchema())),
		RF("regions", Arr("Regions", Str("Region"))),
		RF("images", Arr("Images", Str("Image"))),
	)
}

func managedCredentialSchema() Schema {
	return Obj("Managed credential",
		RF("id", ID("Credential ID")),
		RF("organizationId", Str("Organization ID")),
		RF("provider", Str("Provider")),
		RF("label", Str("Label")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}

func managedQuotaSchema() Schema {
	return Obj("Managed quota",
		RF("organizationId", Str("Organization ID")),
		RF("maxServers", Int("Max servers")),
		RF("maxMemoryGB", Int("Max memory GB")),
		F("usedServers", Int("Used servers")),
		F("usedMemoryGB", Int("Used memory GB")),
		F("updatedAt", Str("Last update")),
	)
}

func managedReviewSchema() Schema {
	return Obj("Managed review",
		RF("operationId", Str("Operation ID")),
		RF("confirmation", Str("Confirmation token")),
		F("monthlyPrice", Num("Monthly price")),
		F("action", Str("Action")),
		F("name", Str("Server name")),
		F("region", Str("Region")),
		F("serverType", Str("Server type")),
	)
}

func serverOperations() []Operation {
	servers := "servers"
	managed := "managed"
	read := Scope("server:read")
	write := Scope("server:write")
	return []Operation{
		{Method: "GET", Path: "/api/servers", Summary: "List servers", Tags: []string{servers}, Auth: read, Response: Arr("Servers", ServerSchema())},
		{Method: "POST", Path: "/api/servers", Summary: "Add a server", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Add server", RF("name", Str("Server name")), F("ipAddress", Str("IP address")), F("isLocal", Bool("Local server")), F("sshHost", Str("SSH hostname")), F("sshPort", Int("SSH port")), F("sshUser", Str("SSH username")), F("sshAuthMethod", Str("Auth method")), F("sshKey", Str("SSH key")), F("sshPrivateKey", Str("SSH private key")), F("sshPassword", Str("SSH password")), F("sshJumpHost", Str("Jump host")))), Response: ServerSchema()},
		{Method: "GET", Path: "/api/servers/:id", Summary: "Get a server", Tags: []string{servers}, Auth: read, Response: ServerSchema()},
		{Method: "PATCH", Path: "/api/servers/:id", Summary: "Update a server", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Update server", F("name", Str("Server name")), F("ipAddress", Str("IP address")), F("sshHost", Str("SSH hostname")), F("sshPort", Int("SSH port")), F("sshUser", Str("SSH username")))), Response: ServerSchema()},
		{Method: "DELETE", Path: "/api/servers/:id", Summary: "Delete a server", Tags: []string{servers}, Auth: write, Response: EmptySchema()},
		{Method: "POST", Path: "/api/servers/test-ssh", Summary: "Test an SSH connection", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Test SSH", RF("sshHost", Str("SSH hostname")), RF("sshPort", Int("SSH port")), RF("sshUser", Str("SSH username")), F("sshKey", Str("SSH key")), F("sshPrivateKey", Str("SSH private key")), F("sshPassword", Str("SSH password")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/ws/servers/:serverId/metrics", Summary: "Stream server metrics", Tags: []string{servers}, Auth: AuthUser, Raw: true, Response: Str("WebSocket metrics session")},
		{Method: "GET", Path: "/api/services/:serviceId/logs", Summary: "Stream service logs", Tags: []string{servers}, Auth: AuthUser, Raw: true, Response: Str("WebSocket log session")},
		{Method: "GET", Path: "/api/ws/services/:id/terminal", Summary: "Open a service terminal", Tags: []string{servers}, Auth: AuthUser, Raw: true, Response: Str("WebSocket terminal session")},
		{Method: "GET", Path: "/api/ws/terminal/:id", Summary: "Open a terminal", Tags: []string{servers}, Auth: AuthUser, Raw: true, Response: Str("WebSocket terminal session")},
		{Method: "GET", Path: "/api/managed/catalog", Summary: "Get the managed server catalog", Tags: []string{managed}, Auth: AuthUser, Response: managedCatalogSchema()},
		{Method: "GET", Path: "/api/organizations/:id/managed/credentials", Summary: "List provider credentials", Tags: []string{managed}, Auth: AuthOrg, Response: Arr("Credentials", managedCredentialSchema())},
		{Method: "POST", Path: "/api/organizations/:id/managed/credentials", Summary: "Add a provider credential", Tags: []string{managed}, Auth: AuthOrgAdmin, Request: JSON(Obj("Add credential", RF("provider", Str("Provider")), RF("label", Str("Label")), RF("token", Str("API token")))), Response: managedCredentialSchema()},
		{Method: "DELETE", Path: "/api/organizations/:id/managed/credentials/:credentialId", Summary: "Delete a provider credential", Tags: []string{managed}, Auth: AuthOrgAdmin, Response: EmptySchema()},
		{Method: "GET", Path: "/api/organizations/:id/managed/quota", Summary: "Get the server quota", Tags: []string{managed}, Auth: AuthOrg, Response: managedQuotaSchema()},
		{Method: "PUT", Path: "/api/organizations/:id/managed/quota", Summary: "Set the server quota", Tags: []string{managed}, Auth: AuthOrgAdmin, Request: JSON(Obj("Set quota", RF("maxServers", Int("Max servers")), RF("maxMemoryGB", Int("Max memory GB")))), Response: managedQuotaSchema()},
		{Method: "POST", Path: "/api/organizations/:id/managed/servers/review", Summary: "Review a server provision", Tags: []string{managed}, Auth: write, Request: JSON(Obj("Review provision", RF("credentialId", Str("Credential ID")), RF("name", Str("Server name")), RF("region", Str("Region")), RF("serverType", Str("Server type")), RF("image", Str("Image")), F("serverId", Str("Existing server ID")))), Response: managedReviewSchema()},
		{Method: "POST", Path: "/api/organizations/:id/managed/servers/:serverId/resize/review", Summary: "Review a server resize", Tags: []string{managed}, Auth: AuthOrgAdmin, Request: JSON(Obj("Review resize", RF("serverType", Str("Target server type")))), Response: managedReviewSchema()},
		{Method: "POST", Path: "/api/organizations/:id/managed/servers/:serverId/delete/review", Summary: "Review a server deletion", Tags: []string{managed}, Auth: AuthOrgAdmin, Response: managedReviewSchema()},
		{Method: "POST", Path: "/api/managed-operations/:operationId/apply", Summary: "Apply a managed operation", Tags: []string{managed}, Auth: write, Request: JSON(operationApplySchema()), Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/managed/servers/:serverId/refresh", Summary: "Refresh a managed server", Tags: []string{managed}, Auth: read, Response: ServerSchema()},
	}
}
