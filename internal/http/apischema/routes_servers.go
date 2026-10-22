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

func managedLocationSchema() Schema {
	return Obj("Managed location",
		RF("name", Str("Location name")),
		RF("city", Str("City")),
		RF("country", Str("Country")),
	)
}

func managedOffersSchema() Schema {
	return Obj("Managed offers",
		RF("provider", Str("Provider")),
		RF("locations", Arr("Locations", managedLocationSchema())),
		RF("offers", Arr("Offers", Obj("Offer",
			RF("tier", managedTierSchema()),
			RF("prices", Any("Per-location monthly prices")),
		))),
		RF("images", Arr("Images", Str("Image"))),
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

func serverComponentSchema() Schema {
	return Obj("Server component",
		RF("name", Str("Component name")),
		RF("label", Str("Display label")),
		RF("description", Str("Description")),
		RF("installable", Bool("Installable")),
		RF("installed", Bool("Installed")),
		F("version", Str("Version")),
		RF("healthy", Bool("Healthy")),
		RF("message", Str("Status message")),
	)
}

func serverPortScanSchema() Schema {
	return Obj("Port scan result",
		RF("scanned", Bool("Scan completed")),
		RF("serverId", Str("Server ID")),
		RF("listeners", Arr("Listeners", Obj("Listener",
			RF("port", Int("Port")),
			RF("protocol", Str("Protocol")),
			RF("state", Str("State")),
			RF("exposed", Bool("Publicly exposed")),
			RF("process", Str("Process")),
		))),
	)
}

func serverRateLimitSchema() Schema {
	return Obj("Server rate limit",
		RF("rps", Int("Requests per second")),
		RF("burst", Int("Burst size")),
		RF("whitelist", Arr("Whitelist", Str("IP or CIDR"))),
	)
}

func serverNetworkSettingsSchema() Schema {
	return Obj("Server network settings",
		F("internetAccess", Bool("Internet access")),
		F("ingressPorts", Arr("Ingress ports", Int("Port"))),
		F("ingressAll", Bool("Ingress all")),
		F("egress", Arr("Egress", Str("Egress rule"))),
		F("privateIp", Str("Private IP")),
		F("outboundIp", Str("Outbound IP")),
		F("outboundMode", Str("Outbound mode")),
		F("revision", Str("Revision")),
	)
}

func managedServerInfoSchema() Schema {
	return Obj("Managed server info",
		RF("workspaceId", Str("Workspace ID")),
		RF("image", Str("Runner image")),
		RF("state", Str("State")),
		F("mode", Str("Mode")),
		F("operatingSystem", Str("Operating system")),
		F("restartPolicy", Str("Restart policy")),
		RF("resources", Obj("Resources",
			RF("cpuCores", Num("CPU cores")),
			RF("memoryMb", Int("Memory MB")),
			RF("diskMb", Int("Disk MB")),
		)),
	)
}

func managedLogsSchema() Schema {
	return Obj("Managed logs",
		RF("logs", Str("Log output")),
		RF("truncated", Bool("Truncated")),
	)
}

func managedSshStatusSchema() Schema {
	return Obj("Managed SSH status",
		F("enabled", Bool("SSH enabled")),
		F("keyConfigured", Bool("Key configured")),
		F("passwordConfigured", Bool("Password configured")),
		RF("requiresIdentityAccess", Bool("Requires identity access")),
		F("connection", Obj("Connection",
			RF("user", Str("SSH user")),
			RF("host", Str("SSH host")),
			RF("bastion", Str("Bastion host")),
			RF("command", Str("SSH command")),
		)),
	)
}

func managedRuntimeStatusSchema() Schema {
	return Obj("Managed runtime status",
		F("enabled", Bool("Runtime enabled")),
		F("running", Bool("Runtime running")),
	)
}

func managedRuntimeCredentialSchema() Schema {
	return Obj("Managed runtime credential",
		RF("endpoint", Str("Runtime endpoint")),
		RF("token", Str("Runtime token")),
		RF("revision", Str("Revision")),
		F("expiresAt", Str("Expiry time")),
	)
}

func managedWorkloadSchema() Schema {
	return Obj("Managed workload",
		RF("id", ID("Workload ID")),
		RF("name", Str("Workload name")),
		RF("state", Str("State")),
		F("restartPolicy", Str("Restart policy")),
		RF("source", Str("Source")),
		F("projectId", Str("Project ID")),
		RF("manageable", Bool("Manageable")),
	)
}

func serverOperations() []Operation {
	servers := "servers"
	managed := "managed"
	read := Scope("server:read")
	write := Scope("server:write")
	return []Operation{
		{Method: "GET", Path: "/api/servers", Summary: "List servers", Tags: []string{servers}, Auth: read, Response: Arr("Servers", ServerSchema()), Query: []Param{QueryParam("search", "Search text")}},
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
		{Method: "GET", Path: "/api/organizations/:id/managed/credentials/:credentialId/offers", Summary: "Get live provider offers", Tags: []string{managed}, Auth: read, Response: managedOffersSchema()},
		{Method: "GET", Path: "/api/managed-operations/:operationId", Summary: "Get a managed operation", Tags: []string{managed}, Auth: read, Response: operationSchema()},
		{Method: "GET", Path: "/api/managed-operations/:operationId/stream", Summary: "Stream a managed operation", Tags: []string{managed}, Auth: read, Raw: true, Response: Str("SSE operation session")},
		{Method: "POST", Path: "/api/managed-operations/:operationId/apply", Summary: "Apply a managed operation", Tags: []string{managed}, Auth: write, Request: JSON(operationApplySchema()), Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/managed/servers/:serverId/refresh", Summary: "Refresh a managed server", Tags: []string{managed}, Auth: read, Response: ServerSchema()},
		{Method: "GET", Path: "/api/servers/:id/components", Summary: "List server components", Tags: []string{servers}, Auth: read, Response: Arr("Components", serverComponentSchema())},
		{Method: "GET", Path: "/api/servers/:id/swarm", Summary: "Inspect Docker Swarm cluster status", Tags: []string{servers}, Auth: read, Response: Obj("Swarm inspection", RF("state", Str("Swarm state")), F("role", Str("Swarm role")), F("nodeId", Str("Node ID")), F("clusterId", Str("Cluster ID")), F("nodeCount", Int("Node count")), F("managerCount", Int("Manager count")), RF("observedAt", Str("Observed at")))},
		{Method: "POST", Path: "/api/servers/:id/scan-ports", Summary: "Scan server ports", Tags: []string{servers}, Auth: read, Response: serverPortScanSchema()},
		{Method: "GET", Path: "/api/servers/:id/rate-limit", Summary: "Get the server rate limit", Tags: []string{servers}, Auth: read, Response: serverRateLimitSchema()},
		{Method: "PUT", Path: "/api/servers/:id/rate-limit", Summary: "Update the server rate limit", Tags: []string{servers}, Auth: write, Request: JSON(serverRateLimitSchema()), Response: serverRateLimitSchema()},
		{Method: "GET", Path: "/api/servers/:id/network-settings", Summary: "Get server network settings", Tags: []string{servers}, Auth: read, Response: serverNetworkSettingsSchema()},
		{Method: "PATCH", Path: "/api/servers/:id/network-settings", Summary: "Update server network settings", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Update network settings", RF("internetAccess", Bool("Internet access")), RF("expectedInternetAccess", Bool("Expected internet access")), F("egress", Arr("Egress", Str("Egress rule"))), F("expectedRevision", Str("Expected revision")), RF("confirm", Bool("Confirm")))), Response: serverNetworkSettingsSchema()},
		{Method: "GET", Path: "/api/servers/:id/managed/info", Summary: "Get managed server info", Tags: []string{servers}, Auth: read, Response: managedServerInfoSchema()},
		{Method: "POST", Path: "/api/servers/:id/managed/boot-logs", Summary: "Get managed boot logs", Tags: []string{servers}, Auth: read, Request: JSON(Obj("Boot logs", F("tail", Int("Tail lines")))), Response: managedLogsSchema()},
		{Method: "GET", Path: "/api/servers/:id/managed/ssh", Summary: "Get managed SSH status", Tags: []string{servers}, Auth: read, Response: managedSshStatusSchema()},
		{Method: "PATCH", Path: "/api/servers/:id/managed/ssh", Summary: "Update managed SSH access", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Update SSH", RF("enabled", Bool("SSH enabled")), RF("expectedEnabled", Bool("Expected enabled")), RF("confirm", Bool("Confirm")))), Response: Obj("SSH result", RF("status", managedSshStatusSchema()), F("initialPassword", Str("Initial password")))},
		{Method: "PUT", Path: "/api/servers/:id/managed/ssh/key", Summary: "Set the managed SSH key", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Set SSH key", RF("publicKey", Str("Public key")), RF("confirm", Bool("Confirm")))), Response: Obj("SSH key result", RF("ok", Bool("Updated")))},
		{Method: "PUT", Path: "/api/servers/:id/managed/ssh/password", Summary: "Set the managed SSH password", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Set SSH password", RF("password", Str("Password")), RF("confirm", Bool("Confirm")))), Response: Obj("SSH password result", RF("ok", Bool("Updated")))},
		{Method: "GET", Path: "/api/servers/:id/managed/runtime-api", Summary: "Get managed runtime status", Tags: []string{servers}, Auth: read, Response: managedRuntimeStatusSchema()},
		{Method: "POST", Path: "/api/servers/:id/managed/runtime-api/enable", Summary: "Enable the managed runtime", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Enable runtime", RF("confirm", Bool("Confirm")))), Response: managedRuntimeStatusSchema()},
		{Method: "POST", Path: "/api/servers/:id/managed/runtime-api/credential", Summary: "Get a managed runtime credential", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Runtime credential", RF("confirm", Bool("Confirm")))), Response: managedRuntimeCredentialSchema()},
		{Method: "POST", Path: "/api/servers/:id/managed/runtime-api/rotate", Summary: "Rotate the managed runtime credential", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Rotate credential", RF("confirm", Bool("Confirm")))), Response: managedRuntimeCredentialSchema()},
		{Method: "GET", Path: "/api/servers/:id/managed/workloads", Summary: "List managed workloads", Tags: []string{servers}, Auth: read, Response: Obj("Workloads", RF("workloads", Arr("Workloads", managedWorkloadSchema())), RF("truncated", Bool("Truncated")))},
		{Method: "POST", Path: "/api/servers/:id/managed/workloads", Summary: "Create a managed workload", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Create workload", RF("name", Str("Workload name")), RF("command", Str("Command")), F("workingDirectory", Str("Working directory")), F("environment", Arr("Environment", Str("Variable"))), F("restartPolicy", Str("Restart policy")), RF("confirm", Bool("Confirm")), RF("idempotencyKey", Str("Idempotency key")))), Response: managedWorkloadSchema()},
		{Method: "POST", Path: "/api/servers/:id/managed/workloads/control", Summary: "Control a managed workload", Tags: []string{servers}, Auth: write, Request: JSON(Obj("Control workload", RF("workloadId", Str("Workload ID")), RF("action", StrEnum("Action", "start", "stop", "restart", "delete")), RF("confirm", Bool("Confirm")))), Response: Obj("Control result", RF("ok", Bool("Accepted")), F("workload", managedWorkloadSchema()))},
		{Method: "POST", Path: "/api/servers/:id/managed/workloads/logs", Summary: "Get managed workload logs", Tags: []string{servers}, Auth: read, Request: JSON(Obj("Workload logs", RF("workloadId", Str("Workload ID")), F("tail", Int("Tail lines")))), Response: managedLogsSchema()},
	}
}
