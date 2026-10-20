package apischema

func systemStatsSchema() Schema {
	return Obj("System stats",
		F("cpu", Obj("CPU", F("percent", Num("Percent")), F("cores", Int("Cores")))),
		F("memory", Obj("Memory", F("totalMb", Int64("Total MB")), F("usedMb", Int64("Used MB")), F("percent", Num("Percent")))),
		F("disk", Obj("Disk", F("totalGb", Int64("Total GB")), F("usedGb", Int64("Used GB")), F("percent", Num("Percent")))),
		F("uptimeSeconds", Int64("Uptime seconds")),
		F("processes", Int("Processes")),
	)
}

func settingsSchema() Schema {
	return Obj("Server settings",
		F("registrationEnabled", Bool("Registration enabled")),
		F("customDnsResolvers", Str("Custom DNS resolvers")),
		F("dnsValidationEnabled", Bool("DNS validation")),
		F("ipAllowlist", Str("IP allowlist")),
		F("mcpServerEnabled", Bool("MCP server enabled")),
		F("defaultWildcardDomain", Str("Wildcard domain")),
		F("dashboardDomain", Str("Dashboard domain")),
		F("siteName", Str("Site name")),
		F("autoUpdateEnabled", Bool("Auto update")),
		F("telemetryEnabled", Bool("Telemetry")),
		F("concurrentBuilds", Int("Concurrent builds")),
		F("deploymentTimeout", Int("Deployment timeout")),
		F("serverTimezone", Str("Timezone")),
		F("currentVersion", Str("Current version")),
		F("latestVersion", Str("Latest version")),
	)
}

func appearanceSchema() Schema {
	return Obj("Appearance preferences",
		F("activeTheme", Str("Active theme name")),
		F("mode", Str("Color mode")),
		F("fontFamily", Str("Font family token")),
		F("radius", Str("Border radius token")),
		F("tokens", Map("Theme CSS tokens", Str("CSS property value"))),
		F("customThemes", Arr("Custom themes", Obj("Custom theme",
			RF("id", Str("Theme ID")),
			RF("name", Str("Theme name")),
			F("mode", Str("Color mode")),
			F("tokens", Map("Theme CSS tokens", Str("CSS property value"))),
		))),
	)
}

func aiSettingsSchema() Schema {
	return Obj("AI settings",
		F("defaultProvider", Str("Default provider")),
		F("openAIModel", Str("OpenAI model")),
		F("anthropicModel", Str("Anthropic model")),
		F("googleModel", Str("Google model")),
		F("mistralModel", Str("Mistral model")),
		F("groqModel", Str("Groq model")),
		F("deepSeekModel", Str("DeepSeek model")),
		F("xaiModel", Str("xAI model")),
		F("moonshotModel", Str("Moonshot model")),
	)
}

func takeoverRunSchema() Schema {
	props := []Property{
		RF("id", ID("Run ID")),
		RF("sourceHost", Str("Source host")),
		F("sourcePlatform", Str("Source platform")),
		RF("status", Str("Status")),
		F("error", Str("Error")),
	}
	return Obj("Takeover run", append(props, Timestamps()...)...)
}

func dnsRecordSchema() Schema {
	props := []Property{
		RF("id", ID("Record ID")),
		RF("domainName", Str("Domain name")),
		RF("recordType", Str("Record type")),
		RF("recordName", Str("Record name")),
		RF("recordValue", Str("Record value")),
		F("ttl", Int("TTL")),
	}
	return Obj("DNS record", append(props, Timestamps()...)...)
}

func scheduledTaskSchema() Schema {
	props := []Property{
		RF("id", ID("Task ID")),
		RF("serviceId", Str("Service ID")),
		RF("name", Str("Task name")),
		RF("schedule", Str("Cron schedule")),
		RF("command", Str("Command")),
		F("status", Str("Status")),
		F("lastRunAt", DateTime("Last run")),
		F("lastOutput", Str("Last output")),
	}
	return Obj("Scheduled task", append(props, Timestamps()...)...)
}

func githubAppSchema() Schema {
	return Obj("GitHub app",
		RF("id", ID("App ID")),
		RF("name", Str("App name")),
		F("appId", Str("GitHub app ID")),
		F("installationId", Str("Installation ID")),
		F("clientId", Str("Client ID")),
		F("isPublic", Bool("Public")),
	)
}

func ingressConnectionSchema() Schema {
	return Obj("Ingress connection",
		RF("id", Str("Connection ID")),
		RF("serverId", Str("Server ID")),
		RF("credentialId", Str("Credential ID")),
		RF("accountId", Str("Account ID")),
		F("tunnelId", Str("Tunnel ID")),
		RF("name", Str("Tunnel name")),
		RF("connector", Str("Connector type")),
		RF("edgePort", Int("Edge port")),
		RF("managedTunnel", Bool("Managed tunnel")),
		RF("defaultForDomains", Bool("Default for domains")),
		F("defaultAccessPolicyId", Str("Default access policy ID")),
		RF("status", Str("Status")),
		F("error", Str("Error")),
		RF("createdAt", Str("Creation timestamp")),
		RF("updatedAt", Str("Update timestamp")),
	)
}

func ingressHostnameSchema() Schema {
	return Obj("Ingress hostname",
		RF("id", Str("Hostname ID")),
		F("domainId", Str("Domain ID")),
		RF("connectionId", Str("Connection ID")),
		RF("hostname", Str("Hostname")),
		F("accessPolicyId", Str("Access policy ID")),
		RF("status", Str("Status")),
		RF("phase", Str("Phase")),
		F("error", Str("Error")),
		RF("updatedAt", Str("Update timestamp")),
	)
}

func ingressDiscoverySchema() Schema {
	return Obj("Ingress discovery",
		RF("accounts", Arr("Accounts", Obj("Account", RF("id", Str("ID")), RF("name", Str("Name"))))),
		RF("tunnels", Arr("Tunnels", Obj("Tunnel", RF("id", Str("ID")), RF("name", Str("Name")), RF("remote", Bool("Remote")), RF("status", Str("Status"))))),
		RF("policies", Arr("Policies", Obj("Policy", RF("id", Str("ID")), RF("name", Str("Name")), RF("decision", Str("Decision"))))),
		F("policiesError", Str("Policies error")),
	)
}

func networkWatchdogSchema() Schema {
	return Obj("Network rollback watchdog",
		RF("operationId", Str("Operation ID")),
		RF("armed", Bool("Armed")),
		RF("timeoutSeconds", Int("Timeout seconds")),
		RF("remainingSeconds", Int("Remaining seconds")),
		RF("expiresAt", Str("Expiration timestamp")),
		RF("canCommit", Bool("Can commit")),
		RF("canRollback", Bool("Can rollback")),
		RF("status", Str("Status")),
	)
}

func serverNetworkSchema() Schema {
	return Obj("Server network",
		RF("id", Str("Network ID")),
		F("organizationId", Str("Organization ID")),
		RF("name", Str("Network name")),
		RF("mode", Str("Mode")),
		RF("subnetCidr", Str("Subnet CIDR")),
		RF("mtu", Int("MTU")),
		RF("verificationPort", Int("Verification port")),
		RF("status", Str("Status")),
		RF("revision", Int("Revision")),
		F("members", Arr("Members", Obj("Member", RF("serverId", Str("Server ID")), F("serverName", Str("Server name")), RF("privateIp", Str("Private IP")), RF("interfaceName", Str("Interface name")), F("status", Str("Status")), F("latencyMs", Num("Latency ms"))))),
		F("accessPolicy", Obj("Access policy", RF("mode", Str("Mode")), F("allowedPairs", Arr("Pairs", Arr("Pair", Str("Server ID")))))),
		F("watchdog", networkWatchdogSchema()),
		RF("createdAt", Str("Creation timestamp")),
		RF("updatedAt", Str("Update timestamp")),
	)
}

func actionWorkflowSchema() Schema {
	return Obj("Action workflow",
		RF("id", Str("Workflow ID")),
		F("projectId", Str("Project ID")),
		RF("name", Str("Workflow name")),
		F("description", Str("Description")),
		RF("trigger", Str("Trigger")),
		F("yaml", Str("Workflow YAML")),
		RF("enabled", Bool("Enabled")),
		RF("createdAt", Str("Created timestamp")),
		RF("updatedAt", Str("Updated timestamp")),
	)
}

func actionStepSchema() Schema {
	return Obj("Action step",
		RF("id", Str("Step ID")),
		RF("name", Str("Step name")),
		RF("status", Str("Step status")),
		F("durationSeconds", Int("Duration in seconds")),
	)
}

func actionJobSchema() Schema {
	return Obj("Action job",
		RF("id", Str("Job ID")),
		RF("name", Str("Job name")),
		RF("status", Str("Job status")),
		F("runnerId", Str("Runner ID")),
		F("durationSeconds", Int("Duration in seconds")),
		RF("steps", Arr("Steps", actionStepSchema())),
	)
}

func actionRunSchema() Schema {
	return Obj("Action run",
		RF("id", Str("Run ID")),
		RF("workflowId", Str("Workflow ID")),
		RF("name", Str("Run name")),
		RF("number", Int("Run number")),
		RF("status", Str("Run status")),
		RF("triggeredBy", Str("Triggered by")),
		F("runnerId", Str("Runner ID")),
		F("logs", Str("Execution logs")),
		F("jobs", Arr("Jobs", actionJobSchema())),
		F("startedAt", Str("Started timestamp")),
		F("finishedAt", Str("Finished timestamp")),
		RF("durationSeconds", Int("Duration in seconds")),
		RF("createdAt", Str("Created timestamp")),
		RF("updatedAt", Str("Updated timestamp")),
	)
}

func actionRunnerSchema() Schema {
	return Obj("Action runner",
		RF("id", Str("Runner ID")),
		RF("name", Str("Runner name")),
		RF("mode", Str("Runner mode")),
		F("image", Str("Container image")),
		RF("status", Str("Runner status")),
		RF("labels", Arr("Labels", Str("Label"))),
		RF("enabled", Bool("Enabled")),
		F("lastHeartbeatAt", Str("Last heartbeat timestamp")),
		RF("createdAt", Str("Created timestamp")),
		RF("updatedAt", Str("Updated timestamp")),
	)
}

func supportSessionSchema() Schema {
	return Obj("Support session",
		F("account", Obj("Account",
			RF("key", Str("Account key")),
			RF("id", Str("Account ID")),
			F("name", Str("Account name")),
			RF("email", Str("Account email")),
		)),
	)
}

func supportTicketSchema() Schema {
	return Obj("Support ticket",
		RF("id", Str("Ticket ID")),
		RF("subject", Str("Ticket subject")),
		RF("category", Str("Category")),
		RF("status", Str("Status")),
		RF("createdAt", Str("Created timestamp")),
		RF("updatedAt", Str("Updated timestamp")),
	)
}

func supportDetailSchema() Schema {
	return Obj("Support ticket detail",
		RF("ticket", Obj("Ticket",
			RF("id", Str("Ticket ID")),
			RF("subject", Str("Ticket subject")),
			RF("category", Str("Category")),
			RF("status", Str("Status")),
			RF("message", Str("Initial message")),
			RF("createdAt", Str("Created timestamp")),
			RF("updatedAt", Str("Updated timestamp")),
		)),
		RF("messages", Arr("Messages", Obj("Message",
			RF("id", Str("Message ID")),
			RF("author", Str("Author")),
			RF("body", Str("Message body")),
			RF("createdAt", Str("Created timestamp")),
		))),
	)
}

func supportReceiptSchema() Schema {
	return Obj("Support ticket receipt",
		RF("id", Str("Ticket ID")),
		RF("createdAt", Str("Created timestamp")),
	)
}

func systemOperations() []Operation {
	system := "system"
	settings := "settings"
	takeover := "takeover"
	dns := "dns"
	tasks := "tasks"
	mcp := "mcp"
	networks := "networks"
	actions := "actions"
	support := "support"
	return []Operation{
		{Method: "GET", Path: "/api/meta", Summary: "Get instance metadata and capabilities", Tags: []string{system}, Response: Obj("Metadata", RF("name", Str("Name")), RF("version", Str("Version")), RF("capabilities", Obj("Capabilities", RF("sseLogs", Bool("SSE logs")), RF("push", Bool("Push")), RF("actions", Bool("Actions")), RF("apps", Bool("Apps")))))},
		{Method: "GET", Path: "/api/system/public", Summary: "Get public settings", Tags: []string{system}, Response: settingsSchema()},
		{Method: "GET", Path: "/api/system/stats", Summary: "Get system stats", Tags: []string{system}, Auth: AuthUser, Response: systemStatsSchema()},
		{Method: "GET", Path: "/api/system/monitor/stream", Summary: "Stream server monitoring metrics", Tags: []string{system}, Auth: Scope("server:read"), Raw: true, Content: ContentSSE, Response: Str("SSE server monitor stream")},
		{Method: "POST", Path: "/api/system/restart", Summary: "Restart the daemon", Tags: []string{system}, Auth: AuthAdmin, Response: StatusSchema("restarting")},
		{Method: "POST", Path: "/api/system/maintenance/cleanup", Summary: "Run maintenance cleanup", Tags: []string{system}, Auth: AuthAdmin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/system/export", Summary: "Export an instance bundle", Tags: []string{system}, Auth: AuthAdmin, Request: JSON(Obj("Export", RF("passphrase", Str("Bundle passphrase")))), Response: Any("Bundle manifest")},
		{Method: "POST", Path: "/api/system/import", Summary: "Import an instance bundle", Tags: []string{system}, Auth: AuthAdmin, Request: Form(Obj("Import", RF("bundle", Str("Bundle file")), RF("passphrase", Str("Bundle passphrase")))), Response: Any("Import manifest")},
		{Method: "POST", Path: "/api/system/setup/import", Summary: "Import a bundle during setup", Tags: []string{system}, Request: Form(Obj("Import", RF("bundle", Str("Bundle file")), RF("passphrase", Str("Bundle passphrase")))), Response: Any("Import manifest")},
		{Method: "POST", Path: "/api/system/local-cli-session", Summary: "Create temporary CLI session from local operator token", Tags: []string{system}, Response: Obj("CLI session", RF("token", Str("Session token")), RF("expiresAt", Str("Expiration timestamp")))},
		{Method: "POST", Path: "/api/system/takeover/scan", Summary: "Scan a host for adoption", Tags: []string{takeover}, Auth: AuthAdmin, Request: JSON(Obj("Scan", RF("host", Str("SSH hostname")), RF("sshUser", Str("SSH username")), RF("sshKey", Str("SSH key")), F("sshFingerprint", Str("Expected fingerprint")), RF("platform", Str("Platform")))), Response: Obj("Scan started", RF("runId", Str("Run ID")))},
		{Method: "POST", Path: "/api/system/takeover/adopt", Summary: "Adopt discovered containers", Tags: []string{takeover}, Auth: AuthAdmin, Request: JSON(Obj("Adopt", RF("runId", Str("Run ID")), RF("projectName", Str("Project name")), F("serviceNames", Arr("Services", Str("Service"))), F("importEnv", Bool("Import env")))), Response: Obj("Adopted", RF("projectIds", Arr("Projects", Str("Project ID"))))},
		{Method: "GET", Path: "/api/system/takeover/runs", Summary: "List takeover runs", Tags: []string{takeover}, Auth: AuthAdmin, Response: Arr("Runs", takeoverRunSchema())},
		{Method: "GET", Path: "/api/system/takeover/runs/:id", Summary: "Get a takeover run", Tags: []string{takeover}, Auth: AuthAdmin, Response: takeoverRunSchema()},
		{Method: "GET", Path: "/api/settings", Summary: "Get server settings", Tags: []string{settings}, Auth: AuthUser, Response: settingsSchema()},
		{Method: "PUT", Path: "/api/settings", Summary: "Update server settings", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(settingsSchema()), Response: settingsSchema()},
		{Method: "GET", Path: "/api/settings/appearance", Summary: "Get user appearance preferences", Tags: []string{settings}, Auth: AuthUser, Response: Obj("User appearance", F("appearance", appearanceSchema()))},
		{Method: "PUT", Path: "/api/settings/appearance", Summary: "Update user appearance preferences", Tags: []string{settings}, Auth: AuthUser, Request: JSON(Obj("User appearance input", F("appearance", appearanceSchema()))), Response: Obj("User appearance", F("appearance", appearanceSchema()))},
		{Method: "GET", Path: "/api/ai", Summary: "Get AI settings", Tags: []string{settings}, Auth: AuthUser, Response: aiSettingsSchema()},
		{Method: "PUT", Path: "/api/ai", Summary: "Update AI settings", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(aiSettingsSchema()), Response: aiSettingsSchema()},
		{Method: "POST", Path: "/api/ai/diagnose", Summary: "Diagnose logs with AI", Tags: []string{settings}, Auth: AuthUser, Raw: true, Content: ContentText, Request: JSON(Obj("Diagnose", RF("prompt", Str("Diagnosis prompt")))), Response: Str("Diagnosis")},
		{Method: "POST", Path: "/api/settings/notifications/test", Summary: "Send a test notification", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(Obj("Test", RF("provider", Str("Provider")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/settings/updates/status", Summary: "Get update status", Tags: []string{settings}, Auth: AuthUser, Response: Obj("Update status", F("currentVersion", Str("Current version")), F("latestVersion", Str("Latest version")), F("lastUpdateCheck", Str("Last check")))},
		{Method: "POST", Path: "/api/settings/updates/check", Summary: "Check for updates", Tags: []string{settings}, Auth: AuthAdmin, Response: Obj("Update status", F("currentVersion", Str("Current version")), F("latestVersion", Str("Latest version")))},
		{Method: "POST", Path: "/api/settings/updates/deploy", Summary: "Deploy an update", Tags: []string{settings}, Auth: AuthAdmin, Code: 202, Response: EmptySchema()},
		{Method: "GET", Path: "/api/settings/git_apps/github", Summary: "List GitHub apps", Tags: []string{settings}, Auth: AuthUser, Response: Arr("Apps", githubAppSchema())},
		{Method: "GET", Path: "/api/settings/git_apps/github/:id", Summary: "Get a GitHub app", Tags: []string{settings}, Auth: AuthUser, Response: githubAppSchema()},
		{Method: "PUT", Path: "/api/settings/git_apps/github", Summary: "Save a GitHub app", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(githubAppSchema()), Response: githubAppSchema()},
		{Method: "DELETE", Path: "/api/settings/git_apps/github/:id", Summary: "Delete a GitHub app", Tags: []string{settings}, Auth: AuthAdmin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/settings/git_apps/github/manifest-callback", Summary: "Exchange a manifest code", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(Obj("Exchange", RF("code", Str("Manifest code")))), Response: githubAppSchema()},
		{Method: "GET", Path: "/api/dns", Summary: "List DNS records", Tags: []string{dns}, Auth: AuthAdmin, Response: Arr("Records", dnsRecordSchema())},
		{Method: "POST", Path: "/api/dns", Summary: "Create a DNS record", Tags: []string{dns}, Auth: AuthAdmin, Code: 201, Request: JSON(Obj("Create record", RF("domainName", Str("Domain name")), RF("recordType", Str("Record type")), RF("recordName", Str("Record name")), RF("recordValue", Str("Record value")), F("ttl", Int("TTL")))), Response: dnsRecordSchema()},
		{Method: "PUT", Path: "/api/dns/:id", Summary: "Update a DNS record", Tags: []string{dns}, Auth: AuthAdmin, Request: JSON(Obj("Update record", RF("recordType", Str("Record type")), RF("recordName", Str("Record name")), RF("recordValue", Str("Record value")), F("ttl", Int("TTL")))), Response: dnsRecordSchema()},
		{Method: "DELETE", Path: "/api/dns/:id", Summary: "Delete a DNS record", Tags: []string{dns}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/dns/ingress", Summary: "List ingress connections", Tags: []string{dns}, Auth: AuthUser, Response: Arr("Connections", ingressConnectionSchema())},
		{Method: "POST", Path: "/api/dns/ingress/discover", Summary: "Discover Cloudflare resources", Tags: []string{dns}, Auth: AuthUser, Request: JSON(Obj("Discover", RF("credentialId", Str("Credential ID")), F("accountId", Str("Account ID")))), Response: ingressDiscoverySchema()},
		{Method: "POST", Path: "/api/dns/ingress", Summary: "Connect ingress tunnel", Tags: []string{dns}, Auth: AuthAdmin, Request: JSON(Obj("Connect", RF("serverId", Str("Server ID")), RF("name", Str("Tunnel name")), RF("connector", Str("Connector")))), Response: ingressConnectionSchema()},
		{Method: "PUT", Path: "/api/dns/ingress/:id", Summary: "Update ingress connection preferences", Tags: []string{dns}, Auth: AuthAdmin, Request: JSON(Obj("Update", RF("defaultForDomains", Bool("Default for domains")))), Response: ingressConnectionSchema()},
		{Method: "POST", Path: "/api/dns/ingress/:id/retry", Summary: "Retry ingress reconciliation", Tags: []string{dns}, Auth: AuthAdmin, Response: ingressConnectionSchema()},
		{Method: "DELETE", Path: "/api/dns/ingress/:id", Summary: "Disconnect ingress tunnel", Tags: []string{dns}, Auth: AuthAdmin, Response: Obj("Deleted", RF("success", Bool("Success")))},
		{Method: "GET", Path: "/api/dns/ingress/domains/:domainId", Summary: "Get domain ingress", Tags: []string{dns}, Auth: AuthUser, Response: ingressHostnameSchema()},
		{Method: "PUT", Path: "/api/dns/ingress/domains/:domainId", Summary: "Set domain ingress", Tags: []string{dns}, Auth: AuthAdmin, Request: JSON(Obj("Set domain ingress", F("connectionId", Str("Connection ID")))), Response: ingressHostnameSchema()},
		{Method: "GET", Path: "/api/system/networks", Summary: "List server networks", Tags: []string{networks}, Auth: AuthUser, Response: Arr("Networks", serverNetworkSchema())},
		{Method: "GET", Path: "/api/system/networks/:id", Summary: "Get server network", Tags: []string{networks}, Auth: AuthUser, Response: serverNetworkSchema()},
		{Method: "POST", Path: "/api/system/networks", Summary: "Create server network", Tags: []string{networks}, Auth: AuthAdmin, Request: JSON(Obj("Create network", RF("name", Str("Network name")), RF("mode", Str("Mode")))), Response: serverNetworkSchema()},
		{Method: "PUT", Path: "/api/system/networks/:id/access", Summary: "Revise network access policy", Tags: []string{networks}, Auth: AuthAdmin, Request: JSON(Obj("Revise access", RF("accessPolicy", Obj("Access policy", RF("mode", Str("Mode")))))), Response: serverNetworkSchema()},
		{Method: "DELETE", Path: "/api/system/networks/:id", Summary: "Delete server network", Tags: []string{networks}, Auth: AuthAdmin, Response: Obj("Deleted", RF("success", Bool("Success")))},
		{Method: "GET", Path: "/api/system/networks/operations/:operationId/watchdog", Summary: "Get network rollback watchdog", Tags: []string{networks}, Auth: AuthUser, Response: networkWatchdogSchema()},
		{Method: "POST", Path: "/api/system/networks/operations/:operationId/watchdog/commit", Summary: "Commit network rollback watchdog", Tags: []string{networks}, Auth: AuthAdmin, Response: networkWatchdogSchema()},
		{Method: "POST", Path: "/api/system/networks/operations/:operationId/watchdog/rollback", Summary: "Rollback network watchdog changes", Tags: []string{networks}, Auth: AuthAdmin, Response: networkWatchdogSchema()},
		{Method: "GET", Path: "/api/scheduled-tasks", Summary: "List scheduled tasks", Tags: []string{tasks}, Auth: AuthUser, Response: Arr("Tasks", scheduledTaskSchema())},
		{Method: "POST", Path: "/api/scheduled-tasks", Summary: "Create a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create task", RF("serviceId", Str("Service ID")), RF("name", Str("Task name")), RF("schedule", Str("Cron schedule")), RF("command", Str("Command")))), Response: scheduledTaskSchema()},
		{Method: "GET", Path: "/api/scheduled-tasks/:id", Summary: "Get a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Response: scheduledTaskSchema()},
		{Method: "PUT", Path: "/api/scheduled-tasks/:id", Summary: "Update a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Request: JSON(Obj("Update task", F("name", Str("Task name")), F("schedule", Str("Cron schedule")), F("command", Str("Command")), F("status", Str("Status")))), Response: scheduledTaskSchema()},
		{Method: "DELETE", Path: "/api/scheduled-tasks/:id", Summary: "Delete a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/scheduled-tasks/:id/trigger", Summary: "Trigger a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Response: StatusSchema("executed")},
		{Method: "GET", Path: "/api/mcp/sse", Summary: "Open the MCP event stream", Tags: []string{mcp}, Auth: AuthUser, Raw: true, Content: ContentSSE, Response: Str("MCP events")},
		{Method: "GET", Path: "/api/actions/workflows", Summary: "List action workflows", Tags: []string{actions}, Auth: AuthUser, Response: Arr("Workflows", actionWorkflowSchema())},
		{Method: "POST", Path: "/api/actions/workflows", Summary: "Create action workflow", Tags: []string{actions}, Auth: AuthAdmin, Code: 200, Request: JSON(Obj("Create workflow", F("projectId", Str("Project ID")), RF("name", Str("Workflow name")), F("description", Str("Description")), F("trigger", Str("Trigger")), F("yaml", Str("YAML")))), Response: actionWorkflowSchema()},
		{Method: "DELETE", Path: "/api/actions/workflows/:id", Summary: "Delete action workflow", Tags: []string{actions}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/actions/workflows/:id/dispatch", Summary: "Dispatch action workflow", Tags: []string{actions}, Auth: AuthUser, Code: 200, Request: JSON(Obj("Dispatch workflow", F("inputs", Any("Inputs")))), Response: actionRunSchema()},
		{Method: "GET", Path: "/api/actions/runs", Summary: "List action runs", Tags: []string{actions}, Auth: AuthUser, Response: Arr("Runs", actionRunSchema())},
		{Method: "GET", Path: "/api/actions/runs/:id", Summary: "Get action run", Tags: []string{actions}, Auth: AuthUser, Response: actionRunSchema()},
		{Method: "DELETE", Path: "/api/actions/runs/:id", Summary: "Delete action run", Tags: []string{actions}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/actions/runs/:id/retry", Summary: "Retry action run", Tags: []string{actions}, Auth: AuthUser, Response: actionRunSchema()},
		{Method: "POST", Path: "/api/actions/runs/:id/cancel", Summary: "Cancel action run", Tags: []string{actions}, Auth: AuthUser, Response: actionRunSchema()},
		{Method: "GET", Path: "/api/actions/runners", Summary: "List action runners", Tags: []string{actions}, Auth: AuthUser, Response: Arr("Runners", actionRunnerSchema())},
		{Method: "POST", Path: "/api/actions/runners", Summary: "Register action runner", Tags: []string{actions}, Auth: AuthAdmin, Code: 200, Request: JSON(Obj("Create runner", F("serverId", Str("Server ID")), RF("name", Str("Runner name")), F("mode", Str("Mode")), F("image", Str("Image")), F("labels", Arr("Labels", Str("Label"))), RF("enabled", Bool("Enabled")))), Response: actionRunnerSchema()},
		{Method: "POST", Path: "/api/actions/runners/inspect", Summary: "Inspect destination capabilities", Tags: []string{actions}, Auth: AuthAdmin, Code: 200, Request: JSON(Obj("Inspect destination", RF("serverId", Str("Server ID")))), Response: Obj("Capabilities", RF("os", Str("OS")), RF("architecture", Str("Architecture")), RF("docker", Bool("Docker")), RF("git", Bool("Git")), RF("node", Bool("Node")))},
		{Method: "DELETE", Path: "/api/actions/runners/:id", Summary: "Delete action runner", Tags: []string{actions}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/cloud/support/session", Summary: "Get support session", Tags: []string{support}, Auth: AuthUser, Response: supportSessionSchema()},
		{Method: "GET", Path: "/api/cloud/support/mine", Summary: "List support tickets", Tags: []string{support}, Auth: AuthUser, Response: Obj("Support list", RF("tickets", Arr("Tickets", supportTicketSchema())), F("nextCursor", Str("Next cursor")))},
		{Method: "POST", Path: "/api/cloud/support/mine", Summary: "Create support ticket", Tags: []string{support}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create ticket", RF("requestId", Str("Request ID")), RF("subject", Str("Subject")), RF("message", Str("Message")), RF("category", Str("Category")))), Response: supportReceiptSchema()},
		{Method: "GET", Path: "/api/cloud/support/mine/:id", Summary: "Get support ticket", Tags: []string{support}, Auth: AuthUser, Response: supportDetailSchema()},
		{Method: "POST", Path: "/api/cloud/support/mine/:id/replies", Summary: "Reply to support ticket", Tags: []string{support}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Reply ticket", RF("requestId", Str("Request ID")), RF("message", Str("Message")))), Response: supportDetailSchema()},
		{Method: "PATCH", Path: "/api/cloud/support/mine/:id", Summary: "Update support ticket status", Tags: []string{support}, Auth: AuthUser, Request: JSON(Obj("Update status", RF("status", Str("Status")))), Response: supportDetailSchema()},
		{Method: "POST", Path: "/api/mcp/messages", Summary: "Send an MCP message", Tags: []string{mcp}, Auth: AuthUser, Raw: true, Request: JSON(Obj("MCP message", RF("jsonrpc", Str("JSON-RPC version")), F("id", Any("Message ID")), F("method", Str("Method")), F("params", Any("Params")))), Response: Obj("MCP response", RF("jsonrpc", Str("JSON-RPC version")), F("id", Any("Message ID")), F("result", Any("Result")), F("error", Obj("Error", RF("code", Int("Code")), RF("message", Str("Message")))))},
		{Method: "GET", Path: "/api/docs/openapi.json", Summary: "Download the OpenAPI document", Tags: []string{system}, Auth: AuthUser, Raw: true, Response: Any("OpenAPI 3.1 document")},
	}
}
