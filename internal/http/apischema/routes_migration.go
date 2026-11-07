package apischema

func migrationSourceSchema() Schema {
	return Obj("Migration source",
		RF("id", ID("Source ID")),
		RF("organizationId", Str("Organization ID")),
		RF("name", Str("Source name")),
		RF("sshHost", Str("SSH hostname")),
		RF("sshPort", Int("SSH port")),
		RF("sshUser", Str("SSH username")),
		RF("sshAuthMethod", Str("Auth method")),
		F("fingerprint", Str("Host fingerprint")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}

func migrationRunSchema() Schema {
	return Obj("Migration run",
		RF("id", ID("Run ID")),
		RF("organizationId", Str("Organization ID")),
		F("sourceId", Str("Source ID")),
		RF("sourceKind", Str("Source kind")),
		F("projectId", Str("Project ID")),
		F("targetServerId", Str("Target server")),
		RF("mode", StrEnum("Mode", "move", "copy")),
		RF("status", Str("Status")),
		F("phase", Str("Phase")),
		F("logs", Str("Logs")),
		F("error", Str("Error")),
		F("cancelRequested", Bool("Cancel requested")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}

func migrationPreviewSchema() Schema {
	return Obj("Migration preview",
		F("images", Arr("Images", Str("Image"))),
		F("volumes", Arr("Volumes", Str("Volume"))),
		F("services", Arr("Services", Obj("Service", RF("name", Str("Name")), RF("image", Str("Image")), F("ports", Arr("Ports", Str("Port"))), F("volumes", Arr("Volumes", Str("Volume"))), F("status", Str("Status"))))),
		F("conflicts", Arr("Conflicts", Obj("Conflict", RF("kind", Str("Kind")), RF("subject", Str("Subject")), F("detail", Str("Detail"))))),
		F("warnings", Arr("Warnings", Str("Warning"))),
		F("downtime", Str("Expected downtime")),
	)
}

func migrationPromptSchema() Schema {
	return Obj("Migration prompt",
		RF("id", ID("Prompt ID")),
		RF("kind", Str("Prompt kind")),
		RF("subject", Str("Subject")),
		F("detail", Str("Detail")),
		F("options", Arr("Options", Obj("Option", RF("id", Str("ID")), RF("label", Str("Label")), F("description", Str("Description"))))),
		F("expiresAt", Int64("Expiry unix time")),
	)
}

func migrationOperations() []Operation {
	migration := "migration"
	write := Scope("server:write")
	return []Operation{
		{Method: "GET", Path: "/api/organizations/:id/migration/sources", Summary: "List migration sources", Tags: []string{migration}, Auth: AuthOrg, Response: Arr("Sources", migrationSourceSchema())},
		{Method: "POST", Path: "/api/organizations/:id/migration/sources", Summary: "Add a migration source", Tags: []string{migration}, Auth: AuthOrgAdmin, Request: JSON(Obj("Add source", RF("name", Str("Source name")), RF("sshHost", Str("SSH hostname")), RF("sshPort", Int("SSH port")), RF("sshUser", Str("SSH username")), RF("sshAuthMethod", Str("Auth method")), F("sshKey", Str("SSH key")), F("sshPassword", Str("SSH password")), F("fingerprint", Str("Expected fingerprint")))), Response: migrationSourceSchema()},
		{Method: "DELETE", Path: "/api/organizations/:id/migration/sources/:sourceId", Summary: "Delete a migration source", Tags: []string{migration}, Auth: AuthOrgAdmin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/migration/sources/test", Summary: "Test a migration source", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Test source", RF("sshHost", Str("SSH hostname")), RF("sshPort", Int("SSH port")), RF("sshUser", Str("SSH username")), F("fingerprint", Str("Expected fingerprint")))), Response: Obj("Fingerprint", RF("fingerprint", Str("Host fingerprint")))},
		{Method: "POST", Path: "/api/organizations/:id/migration/scan", Summary: "Scan a source", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Scan", RF("sourceId", Str("Source ID")))), Response: Any("Source stack")},
		{Method: "GET", Path: "/api/organizations/:id/migration/scan/stream", Summary: "Stream a source scan", Tags: []string{migration}, Auth: write, Raw: true, Content: ContentSSE, Response: Str("Scan events"), Query: []Param{QueryParam("sourceId", "Source ID")}},
		{Method: "POST", Path: "/api/organizations/:id/migration/reveal-env", Summary: "Reveal container env", Tags: []string{migration}, Auth: AuthOrgAdmin, Request: JSON(Obj("Reveal", RF("sourceId", Str("Source ID")), RF("containerId", Str("Container ID")))), Response: Map("Environment", Str("Value"))},
		{Method: "POST", Path: "/api/organizations/:id/migration/adopt", Summary: "Adopt source containers", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Adopt", RF("sourceId", Str("Source ID")), RF("containerIds", Arr("Containers", Str("Container ID"))), F("projectName", Str("Project name")), F("importEnv", Bool("Import env")))), Response: Any("Adoption result")},
		{Method: "POST", Path: "/api/organizations/:id/migration/reimport", Summary: "Reimport source containers", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Reimport", RF("sourceId", Str("Source ID")), RF("containerIds", Arr("Containers", Str("Container ID"))))), Response: Any("Reimport result")},
		{Method: "POST", Path: "/api/organizations/:id/migration/repo-compose", Summary: "Scan a repo compose file", Tags: []string{migration}, Auth: AuthOrg, Request: JSON(Obj("Repo compose", RF("repoUrl", Str("Repository URL")), F("branch", Str("Branch")), F("composePath", Str("Compose path")))), Response: Any("Source stack")},
		{Method: "POST", Path: "/api/organizations/:id/migration/preview", Summary: "Preview a migration", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Preview", RF("sourceId", Str("Source ID")), RF("projectId", Str("Project ID")), RF("containerIds", Arr("Containers", Str("Container ID"))), F("targetServerId", Str("Target server")), RF("mode", StrEnum("Mode", "move", "copy")))), Response: migrationPreviewSchema()},
		{Method: "POST", Path: "/api/organizations/:id/migration/migrate", Summary: "Start a migration", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Migrate", RF("sourceId", Str("Source ID")), RF("containerIds", Arr("Containers", Str("Container ID"))), F("projectName", Str("Project name")), F("targetServerId", Str("Target server")), RF("mode", StrEnum("Mode", "move", "copy")), F("importEnv", Bool("Import env")), F("killOriginals", Bool("Kill originals")))), Response: migrationRunSchema()},
		{Method: "POST", Path: "/api/organizations/:id/migration/project", Summary: "Start a project move", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Move project", RF("projectId", Str("Project ID")), RF("targetServerId", Str("Target server")), RF("mode", StrEnum("Mode", "move", "copy")))), Response: migrationRunSchema()},
		{Method: "GET", Path: "/api/organizations/:id/migration/active", Summary: "Get the active migration", Tags: []string{migration}, Auth: Scope("server:read"), Response: migrationRunSchema()},
		{Method: "GET", Path: "/api/organizations/:id/migration/runs", Summary: "List migration runs", Tags: []string{migration}, Auth: Scope("server:read"), Response: Arr("Runs", migrationRunSchema())},
		{Method: "GET", Path: "/api/migrations/:runId", Summary: "Get a migration run", Tags: []string{migration}, Auth: Scope("server:read"), Response: migrationRunSchema()},
		{Method: "GET", Path: "/api/migrations/:runId/stream", Summary: "Stream a migration run", Tags: []string{migration}, Auth: Scope("server:read"), Raw: true, Content: ContentSSE, Response: Str("Run events")},
		{Method: "POST", Path: "/api/migrations/:runId/respond", Summary: "Answer a migration prompt", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Respond", RF("promptId", Str("Prompt ID")), RF("optionId", Str("Option ID")), F("value", Str("Value")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/migrations/:runId/cancel", Summary: "Cancel a migration run", Tags: []string{migration}, Auth: write, Response: EmptySchema()},
		{Method: "POST", Path: "/api/migrations/:runId/resume", Summary: "Resume a migration run", Tags: []string{migration}, Auth: write, Response: migrationRunSchema()},
		{Method: "POST", Path: "/api/migrations/:runId/cutover", Summary: "Cut over a migration", Tags: []string{migration}, Auth: write, Request: JSON(Obj("Cutover", RF("confirmation", Str("Confirmation token")), F("kill", Bool("Kill originals")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/migrations/:runId/cleanup-target", Summary: "Clean up a migration target", Tags: []string{migration}, Auth: write, Response: EmptySchema()},
		{Method: "DELETE", Path: "/api/migrations/:runId", Summary: "Delete a migration run", Tags: []string{migration}, Auth: write, Response: EmptySchema()},
	}
}
