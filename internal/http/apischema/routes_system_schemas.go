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
