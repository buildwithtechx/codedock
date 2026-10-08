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

func systemOperations() []Operation {
	system := "system"
	settings := "settings"
	takeover := "takeover"
	dns := "dns"
	tasks := "tasks"
	mcp := "mcp"
	return []Operation{
		{Method: "GET", Path: "/api/system/public", Summary: "Get public settings", Tags: []string{system}, Response: settingsSchema()},
		{Method: "GET", Path: "/api/system/stats", Summary: "Get system stats", Tags: []string{system}, Auth: AuthUser, Response: systemStatsSchema()},
		{Method: "POST", Path: "/api/system/restart", Summary: "Restart the daemon", Tags: []string{system}, Auth: AuthAdmin, Response: StatusSchema("restarting")},
		{Method: "POST", Path: "/api/system/maintenance/cleanup", Summary: "Run maintenance cleanup", Tags: []string{system}, Auth: AuthAdmin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/system/export", Summary: "Export an instance bundle", Tags: []string{system}, Auth: AuthAdmin, Request: JSON(Obj("Export", RF("passphrase", Str("Bundle passphrase")))), Response: Any("Bundle manifest")},
		{Method: "POST", Path: "/api/system/import", Summary: "Import an instance bundle", Tags: []string{system}, Auth: AuthAdmin, Request: Form(Obj("Import", RF("bundle", Str("Bundle file")), RF("passphrase", Str("Bundle passphrase")))), Response: Any("Import manifest")},
		{Method: "POST", Path: "/api/system/setup/import", Summary: "Import a bundle during setup", Tags: []string{system}, Request: Form(Obj("Import", RF("bundle", Str("Bundle file")), RF("passphrase", Str("Bundle passphrase")))), Response: Any("Import manifest")},
		{Method: "POST", Path: "/api/system/takeover/scan", Summary: "Scan a host for adoption", Tags: []string{takeover}, Auth: AuthAdmin, Request: JSON(Obj("Scan", RF("host", Str("SSH hostname")), RF("sshUser", Str("SSH username")), RF("sshKey", Str("SSH key")), F("sshFingerprint", Str("Expected fingerprint")), RF("platform", Str("Platform")))), Response: Obj("Scan started", RF("runId", Str("Run ID")))},
		{Method: "POST", Path: "/api/system/takeover/adopt", Summary: "Adopt discovered containers", Tags: []string{takeover}, Auth: AuthAdmin, Request: JSON(Obj("Adopt", RF("runId", Str("Run ID")), RF("projectName", Str("Project name")), F("serviceNames", Arr("Services", Str("Service"))), F("importEnv", Bool("Import env")))), Response: Obj("Adopted", RF("projectIds", Arr("Projects", Str("Project ID"))))},
		{Method: "GET", Path: "/api/system/takeover/runs", Summary: "List takeover runs", Tags: []string{takeover}, Auth: AuthAdmin, Response: Arr("Runs", takeoverRunSchema())},
		{Method: "GET", Path: "/api/system/takeover/runs/:id", Summary: "Get a takeover run", Tags: []string{takeover}, Auth: AuthAdmin, Response: takeoverRunSchema()},
		{Method: "GET", Path: "/api/settings", Summary: "Get server settings", Tags: []string{settings}, Auth: AuthUser, Response: settingsSchema()},
		{Method: "PUT", Path: "/api/settings", Summary: "Update server settings", Tags: []string{settings}, Auth: AuthAdmin, Request: JSON(settingsSchema()), Response: settingsSchema()},
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
		{Method: "GET", Path: "/api/scheduled-tasks", Summary: "List scheduled tasks", Tags: []string{tasks}, Auth: AuthUser, Response: Arr("Tasks", scheduledTaskSchema())},
		{Method: "POST", Path: "/api/scheduled-tasks", Summary: "Create a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create task", RF("serviceId", Str("Service ID")), RF("name", Str("Task name")), RF("schedule", Str("Cron schedule")), RF("command", Str("Command")))), Response: scheduledTaskSchema()},
		{Method: "GET", Path: "/api/scheduled-tasks/:id", Summary: "Get a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Response: scheduledTaskSchema()},
		{Method: "DELETE", Path: "/api/scheduled-tasks/:id", Summary: "Delete a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/scheduled-tasks/:id/trigger", Summary: "Trigger a scheduled task", Tags: []string{tasks}, Auth: AuthUser, Response: StatusSchema("executed")},
		{Method: "GET", Path: "/api/mcp/sse", Summary: "Open the MCP event stream", Tags: []string{mcp}, Auth: AuthUser, Raw: true, Content: ContentSSE, Response: Str("MCP events")},
		{Method: "POST", Path: "/api/mcp/messages", Summary: "Send an MCP message", Tags: []string{mcp}, Auth: AuthUser, Raw: true, Request: JSON(Obj("MCP message", RF("jsonrpc", Str("JSON-RPC version")), F("id", Any("Message ID")), F("method", Str("Method")), F("params", Any("Params")))), Response: Obj("MCP response", RF("jsonrpc", Str("JSON-RPC version")), F("id", Any("Message ID")), F("result", Any("Result")), F("error", Obj("Error", RF("code", Int("Code")), RF("message", Str("Message")))))},
		{Method: "GET", Path: "/api/docs/openapi.json", Summary: "Download the OpenAPI document", Tags: []string{system}, Auth: AuthUser, Raw: true, Response: Any("OpenAPI 3.1 document")},
	}
}
