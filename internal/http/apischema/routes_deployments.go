package apischema

func deploymentItemSchema() Schema {
	props := []Property{
		RF("id", ID("Deployment ID")),
		RF("serviceId", Str("Service ID")),
		F("serviceName", Str("Service name")),
		RF("projectId", Str("Project ID")),
		RF("status", Str("Deployment status")),
		F("branch", Str("Branch")),
		F("commitHash", Str("Commit hash")),
		F("commitMessage", Str("Commit message")),
		F("trigger", Str("Trigger source")),
		F("finishedAt", DateTime("Finish time")),
	}
	return Obj("Deployment list item", append(props, Timestamps()...)...)
}

func gitProviderSchema() Schema {
	props := []Property{
		RF("id", ID("Connection ID")),
		RF("provider", Str("Git provider")),
		F("accountName", Str("Account name")),
	}
	return Obj("Git provider connection", append(props, Timestamps()...)...)
}

func gitRepoSchema() Schema {
	return Obj("Git repository",
		F("id", Int64("Repository ID")),
		RF("name", Str("Repository name")),
		F("fullName", Str("Full name")),
		F("private", Bool("Private")),
		F("cloneUrl", Str("Clone URL")),
		F("htmlUrl", Str("Web URL")),
		F("defaultBranch", Str("Default branch")),
	)
}

func metricPointSchema() Schema {
	return Obj("Metric point",
		RF("timestamp", Str("Timestamp")),
		F("cpuPercent", Num("CPU percent")),
		F("memoryMB", Num("Memory MB")),
		F("networkRxKB", Num("Network RX KB")),
		F("networkTxKB", Num("Network TX KB")),
	)
}

func oneClickAppSchema() Schema {
	return Obj("One-click app",
		RF("id", Str("App ID")),
		RF("name", Str("App name")),
		RF("description", Str("Description")),
		F("icon", Str("Icon")),
		F("category", Str("Category")),
		F("dockerImage", Str("Primary image")),
		F("defaultPort", Int("Default port")),
		F("services", Arr("Services", Str("Service"))),
		F("volumes", Arr("Volumes", Str("Volume"))),
		F("envVariables", Arr("Variables", Obj("Variable", RF("key", Str("Key")), RF("label", Str("Label")), F("defaultValue", Str("Default")), RF("secret", Bool("Secret"))))),
	)
}

func installPreviewSchema() Schema {
	return Obj("Install preview",
		RF("appId", Str("App ID")),
		RF("name", Str("Install name")),
		RF("services", Arr("Services", Obj("Service", RF("service", Str("Service")), RF("image", Str("Image")), F("env", Arr("Env", Str("Entry"))), F("ports", Arr("Ports", Str("Mapping")))))),
		RF("volumes", Arr("Volumes", Obj("Volume", RF("service", Str("Service")), RF("name", Str("Name")), RF("target", Str("Target"))))),
		F("generatedSecrets", Arr("Generated secrets", Str("Variable"))),
		RF("composeYaml", Str("Rendered compose, secrets masked")),
		RF("digest", Str("Review digest")),
		RF("kind", Str("Install kind")),
	)
}

func installInputSchema() Schema {
	return Obj("Install input",
		RF("appId", Str("App ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Install name")),
		F("environmentId", Str("Environment ID")),
		F("secrets", Map("Secrets", Str("Value"))),
		F("environment", Map("Environment", Str("Value"))),
		F("hostPort", Int("Published host port")),
		F("domain", Str("Public domain")),
		F("digest", Str("Review digest")),
	)
}

func deploymentOperations() []Operation {
	deployments := "deployments"
	git := "git"
	compose := "compose"
	catalog := "catalog"
	return []Operation{
		{Method: "GET", Path: "/api/deployments", Summary: "List organization deployments", Tags: []string{deployments}, Auth: AuthUser, Response: PaginatedSchema(Arr("Deployments", deploymentItemSchema())), Query: []Param{QueryParam("status", "Filter by status"), QueryParam("search", "Search text"), QueryParamInt("page", "Page"), QueryParamInt("limit", "Page size")}},
		{Method: "GET", Path: "/api/projects/:id/deployments", Summary: "List project deployments", Tags: []string{deployments}, Auth: ProjectRole("member"), Response: PaginatedSchema(Arr("Deployments", deploymentItemSchema()))},
		{Method: "GET", Path: "/api/services/:serviceId/deployments", Summary: "List service deployments", Tags: []string{deployments}, Auth: ServiceRole(""), Response: PaginatedSchema(Arr("Deployments", deploymentItemSchema()))},
		{Method: "GET", Path: "/api/deployments/:id/logs", Summary: "Get deployment logs", Tags: []string{deployments}, Auth: Scope("logs:read"), Response: Map("Logs", Str("Log output"))},
		{Method: "GET", Path: "/api/deployments/:id/explain", Summary: "Explain a deployment failure", Tags: []string{deployments}, Auth: AuthUser, Response: Any("Failure explanation")},
		{Method: "POST", Path: "/api/deployments/:id/cancel", Summary: "Cancel a deployment", Tags: []string{deployments}, Auth: AuthUser, Code: 202, Response: EmptySchema()},
		{Method: "POST", Path: "/api/deployments/:id/rollback", Summary: "Roll back a deployment", Tags: []string{deployments}, Auth: AuthUser, Code: 202, Response: DeploymentSchema()},
		{Method: "POST", Path: "/api/services/:serviceId/deploy", Summary: "Trigger a service deployment", Tags: []string{deployments}, Auth: ServiceRole("admin"), Code: 202, Request: JSON(Obj("Trigger", F("branch", Str("Branch override")))), Response: DeploymentSchema()},
		{Method: "POST", Path: "/api/projects/:id/deploy", Summary: "Trigger a project deployment", Tags: []string{deployments}, Auth: ProjectRole("admin"), Code: 202, Response: Any("Triggered deployments")},
		{Method: "GET", Path: "/api/services/:serviceId/metrics", Summary: "Get service metrics", Tags: []string{deployments}, Auth: ServiceRole(""), Response: Arr("Metrics", metricPointSchema())},
		{Method: "GET", Path: "/api/services/:serviceId/metrics/historical", Summary: "Get historical metrics", Tags: []string{deployments}, Auth: ServiceRole(""), Response: Any("Historical metrics"), Query: []Param{QueryParam("from", "Start time"), QueryParam("to", "End time")}},
		{Method: "GET", Path: "/api/services/:serviceId/logs/historical", Summary: "Get historical logs", Tags: []string{deployments}, Auth: ServiceRole(""), Response: Any("Historical logs"), Query: []Param{QueryParam("from", "Start time"), QueryParam("to", "End time"), QueryParamInt("tail", "Tail lines")}},
		{Method: "GET", Path: "/api/services/:serviceId/previews", Summary: "List PR previews", Tags: []string{deployments}, Auth: ServiceRole(""), Response: Arr("Previews", DeploymentSchema())},
		{Method: "POST", Path: "/api/git/connect", Summary: "Connect a git provider", Tags: []string{git}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Connect", RF("provider", StrEnum("Provider", "github", "gitlab", "bitbucket")), RF("accessToken", Str("Access token")), F("accountName", Str("Account name")))), Response: gitProviderSchema()},
		{Method: "GET", Path: "/api/git/status", Summary: "Get git connection status", Tags: []string{git}, Auth: AuthUser, Response: Any("Connection status")},
		{Method: "GET", Path: "/api/git/repos", Summary: "List git repositories", Tags: []string{git}, Auth: AuthUser, Response: Arr("Repositories", gitRepoSchema()), Query: []Param{QueryParam("provider", "Provider filter"), QueryParam("search", "Search text")}},
		{Method: "DELETE", Path: "/api/git/connect/:provider", Summary: "Disconnect a git provider", Tags: []string{git}, Auth: AuthUser, Response: StatusSchema("disconnected")},
		{Method: "POST", Path: "/api/projects/:id/repository-inspection", Summary: "Inspect a repository", Tags: []string{git}, Auth: ProjectRole("admin"), Request: JSON(Obj("Inspect", RF("repositoryUrl", Str("Repository URL")), F("branch", Str("Branch")), F("rootDirectory", Str("Root directory")))), Response: Obj("Inspection", F("framework", Str("Framework")), F("packageManager", Str("Package manager")), F("installCommand", Str("Install command")), F("buildCommand", Str("Build command")), F("startCommand", Str("Start command")), F("internalPort", Int("Internal port")), F("staticOutput", Str("Static output")), F("buildEngine", Str("Build engine")), F("dockerfilePath", Str("Dockerfile path")), F("warnings", Arr("Warnings", Str("Warning"))))},
		{Method: "POST", Path: "/api/webhooks/git/services/:serviceId", Summary: "Handle a generic git webhook", Tags: []string{git}, Raw: true, Code: 202, Request: JSON(Any("Webhook payload")), Response: Str("Service name")},
		{Method: "POST", Path: "/api/webhooks/github/services/:serviceId", Summary: "Handle a GitHub webhook", Tags: []string{git}, Raw: true, Code: 202, Request: JSON(Any("GitHub event")), Response: Str("Service branch")},
		{Method: "POST", Path: "/api/compose/analyze", Summary: "Analyze a compose document", Tags: []string{compose}, Auth: AuthUser, Raw: true, Code: 410, Request: JSON(Obj("Analyze", RF("composeContent", Str("Compose content")), F("projectId", Str("Project ID")))), Response: errorSchema()},
		{Method: "POST", Path: "/api/compose/deploy", Summary: "Deploy a compose file", Tags: []string{compose}, Auth: AuthUser, Raw: true, Code: 410, Request: Form(Obj("Upload", RF("file", Str("Compose file")), F("projectId", Str("Project ID")))), Response: errorSchema()},
		{Method: "POST", Path: "/api/deploy/archive", Summary: "Deploy an archive", Tags: []string{deployments}, Auth: AuthUser, Request: Form(Obj("Upload", RF("file", Str("Archive file")), F("projectId", Str("Project ID")), F("name", Str("App name")))), Response: Any("Deploy result")},
		{Method: "GET", Path: "/api/examples", Summary: "List example apps", Tags: []string{catalog}, Auth: AuthUser, Response: Arr("Examples", Obj("Example", RF("id", Str("ID")), RF("name", Str("Name")), RF("description", Str("Description")), RF("repo", Str("Repository")), F("icon", Str("Icon"))))},
		{Method: "GET", Path: "/api/one-click", Summary: "List one-click apps", Tags: []string{catalog}, Auth: AuthUser, Response: Arr("Apps", oneClickAppSchema())},
		{Method: "GET", Path: "/api/one-click/:id", Summary: "Get a one-click app", Tags: []string{catalog}, Auth: AuthUser, Response: oneClickAppSchema()},
		{Method: "POST", Path: "/api/one-click/review", Summary: "Preview an app install", Tags: []string{catalog}, Auth: AuthUser, Request: JSON(installInputSchema()), Response: installPreviewSchema()},
		{Method: "POST", Path: "/api/one-click/deploy", Summary: "Install a one-click app", Tags: []string{catalog}, Auth: AuthUser, Request: JSON(installInputSchema()), Response: Obj("Install result", RF("kind", Str("Install kind")), F("stack", composeStackSchema()))},
	}
}
