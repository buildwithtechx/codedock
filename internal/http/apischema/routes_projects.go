package apischema

func registrySchema() Schema {
	props := []Property{
		RF("id", ID("Registry ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Registry name")),
		RF("registryUrl", Str("Registry URL")),
		F("username", Str("Username")),
	}
	return Obj("Container registry", append(props, Timestamps()...)...)
}

func projectTokenSchema() Schema {
	return Obj("Project token",
		RF("id", ID("Token ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Token name")),
		F("prefix", Str("Identifying prefix")),
		F("role", Str("Token role")),
		F("scopes", Arr("Scopes", Str("Scope"))),
		F("expiresAt", DateTime("Expiry time")),
		F("createdAt", DateTime("Creation time")),
	)
}

func canvasSummarySchema() Schema {
	return Obj("Project canvas summary",
		RF("id", ID("Project ID")),
		RF("name", Str("Project name")),
		F("description", Str("Description")),
		F("environmentsCount", Int("Environment count")),
		F("appsCount", Int("App count")),
		F("databasesCount", Int("Database count")),
		F("totalServices", Int("Total services")),
		F("onlineServices", Int("Online services")),
	)
}

func canvasSchema() Schema {
	return Obj("Environment canvas",
		RF("id", ID("Canvas ID")),
		RF("name", Str("Canvas name")),
		F("apps", Arr("Apps", AppServiceSchema())),
		F("databases", Arr("Databases", DatabaseSchema())),
		F("nodes", Arr("Nodes", Obj("Node", RF("id", Str("Node ID")), RF("type", Str("Node type")), F("data", Map("Node data", Any("Value")))))),
		F("edges", Arr("Edges", Obj("Edge", RF("id", Str("Edge ID")), RF("source", Str("Source")), RF("target", Str("Target")), F("kind", Str("Kind")), F("label", Str("Label"))))),
		F("revision", Str("Canvas revision")),
	)
}

func composeStackSchema() Schema {
	return Obj("Compose stack",
		RF("id", ID("Stack ID")),
		RF("projectId", Str("Project ID")),
		RF("environmentId", Str("Environment ID")),
		RF("name", Str("Stack name")),
		F("revision", Int("Revision")),
		F("status", Str("Status")),
		F("error", Str("Error message")),
		F("results", Str("Service results JSON")),
		F("updatedAt", Str("Last update")),
	)
}

func composeReviewSchema() Schema {
	return Obj("Compose review",
		RF("config", Str("Resolved configuration")),
		RF("digest", Str("Review digest")),
		RF("services", Arr("Services", Str("Service"))),
		RF("effects", Arr("Effects", Str("Effect"))),
	)
}

func composeStackRequestSchema() Schema {
	return Obj("Compose stack request",
		RF("id", Str("Stack UUID")),
		RF("environmentId", Str("Environment ID")),
		RF("name", Str("Stack name")),
		RF("content", Str("Compose document")),
		F("variables", Map("Variables", Str("Value"))),
		F("revision", Int("Revision")),
		F("digest", Str("Review digest")),
		F("repositoryUrl", Str("Repository URL")),
		F("branch", Str("Branch")),
		F("rootDirectory", Str("Root directory")),
	)
}

func projectAppSchema() Schema {
	props := []Property{
		RF("id", ID("Project app ID")),
		RF("organizationId", Str("Organization ID")),
		RF("name", Str("App name")),
		F("slug", Str("Slug")),
		F("gitProvider", Str("Git provider")),
		F("gitOwner", Str("Repository owner")),
		F("gitRepo", Str("Repository name")),
		F("gitUrl", Str("Repository URL")),
	}
	return Obj("Project app", append(props, Timestamps()...)...)
}

func projectOperations() []Operation {
	projects := "projects"
	canvas := "canvas"
	stacks := "stacks"
	return []Operation{
		{Method: "GET", Path: "/api/projects", Summary: "List projects", Tags: []string{projects}, Auth: AuthUser, Response: PaginatedSchema(Arr("Projects", ProjectSchema()))},
		{Method: "POST", Path: "/api/projects", Summary: "Create a project", Tags: []string{projects}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create project", RF("name", Str("Project name")), F("description", Str("Description")), F("organizationId", Str("Organization ID")), F("serverId", Str("Server ID")), F("gitProvider", Str("Git provider")), F("gitOwner", Str("Repository owner")), F("gitRepo", Str("Repository name")), F("gitBranch", Str("Branch")), F("gitUrl", Str("Repository URL")), F("environmentName", Str("Environment name")))), Response: ProjectSchema()},
		{Method: "GET", Path: "/api/projects/:id", Summary: "Get a project", Tags: []string{projects}, Auth: ProjectRole("member"), Response: ProjectSchema()},
		{Method: "DELETE", Path: "/api/projects/:id", Summary: "Delete a project", Tags: []string{projects}, Auth: ProjectRole("owner"), Response: StatusSchema("deleted")},
		{Method: "GET", Path: "/api/projects/:id/env", Summary: "Get project variables", Tags: []string{projects}, Auth: Scope("env:read"), Response: Map("Variables", Str("Value"))},
		{Method: "POST", Path: "/api/projects/:id/env", Summary: "Set project variables", Tags: []string{projects}, Auth: Scope("env:write"), Request: JSON(Obj("Set variables", F("variables", Map("Variables", Str("Value"))))), Response: Map("Variables", Str("Value"))},
		{Method: "PUT", Path: "/api/projects/:id/env", Summary: "Replace project variables", Tags: []string{projects}, Auth: Scope("env:write"), Request: JSON(Obj("Set variables", F("variables", Map("Variables", Str("Value"))))), Response: Map("Variables", Str("Value"))},
		{Method: "GET", Path: "/api/projects/:id/environments", Summary: "List project environments", Tags: []string{projects}, Auth: ProjectRole("member"), Response: Arr("Environments", EnvironmentSchema())},
		{Method: "POST", Path: "/api/projects/:id/environments", Summary: "Create an environment", Tags: []string{projects}, Auth: ProjectRole("admin"), Code: 201, Request: JSON(Obj("Create environment", RF("name", Str("Environment name")), F("isDefault", Bool("Default environment")))), Response: EnvironmentSchema()},
		{Method: "DELETE", Path: "/api/environments/:id", Summary: "Delete an environment", Tags: []string{projects}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:projectId/tokens", Summary: "List project tokens", Tags: []string{projects}, Auth: Scope("env:read"), Response: Arr("Tokens", projectTokenSchema())},
		{Method: "POST", Path: "/api/projects/:projectId/tokens", Summary: "Create a project token", Tags: []string{projects}, Auth: Scope("env:write"), Code: 201, Request: JSON(Obj("Create token", RF("name", Str("Token name")), RF("role", Str("Token role")), F("environmentId", Str("Environment ID")), F("scopes", Arr("Scopes", Str("Scope"))), F("ipAllowlist", Arr("IP allowlist", Str("CIDR"))), F("expiresAt", DateTime("Expiry time")))), Response: Obj("Created token", RF("token", Str("One-time token")), RF("projectToken", projectTokenSchema()))},
		{Method: "DELETE", Path: "/api/projects/:projectId/tokens/:id", Summary: "Delete a project token", Tags: []string{projects}, Auth: Scope("env:write"), Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:projectId/registries", Summary: "List container registries", Tags: []string{projects}, Auth: ProjectRole("admin"), Response: Arr("Registries", registrySchema())},
		{Method: "POST", Path: "/api/projects/:projectId/registries", Summary: "Add a container registry", Tags: []string{projects}, Auth: ProjectRole("admin"), Request: JSON(Obj("Add registry", RF("name", Str("Registry name")), RF("registryUrl", Str("Registry URL")), F("username", Str("Username")), F("passwordToken", Str("Password or token")))), Response: registrySchema()},
		{Method: "DELETE", Path: "/api/projects/:projectId/registries/:id", Summary: "Delete a container registry", Tags: []string{projects}, Auth: ProjectRole("admin"), Response: Obj("Deleted", RF("success", Bool("Success")))},
		{Method: "GET", Path: "/api/projects/:id/summary", Summary: "Get the project canvas summary", Tags: []string{canvas}, Auth: AuthUser, Response: canvasSummarySchema()},
		{Method: "GET", Path: "/api/canvas/projects", Summary: "List canvas summaries", Tags: []string{canvas}, Auth: AuthUser, Response: Arr("Summaries", canvasSummarySchema())},
		{Method: "GET", Path: "/api/environments/:id/canvas", Summary: "Get the environment canvas", Tags: []string{canvas}, Auth: AuthUser, Response: canvasSchema()},
		{Method: "PUT", Path: "/api/environments/:id/canvas", Summary: "Apply canvas topology", Tags: []string{canvas}, Auth: AuthUser, Request: JSON(Obj("Apply topology", RF("revision", Str("Canvas revision")), RF("dependencies", Arr("Dependencies", Obj("Edge", RF("source", Str("Source")), RF("target", Str("Target")), F("kind", Str("Kind"))))))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/environments/:id/canvas/observe", Summary: "Observe a canvas resource", Tags: []string{canvas}, Auth: Scope("logs:read"), Response: Any("Resource observation"), Query: []Param{QueryParam("resourceId", "Resource ID"), QueryParam("kind", "Resource kind")}},
		{Method: "GET", Path: "/api/project-apps", Summary: "List project apps", Tags: []string{projects}, Auth: AuthUser, Response: Arr("Apps", projectAppSchema())},
		{Method: "POST", Path: "/api/project-apps", Summary: "Create a project app", Tags: []string{projects}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create project app", RF("name", Str("App name")), F("slug", Str("Slug")), F("gitProvider", Str("Git provider")), F("gitOwner", Str("Repository owner")), F("gitRepo", Str("Repository name")), F("gitUrl", Str("Repository URL")))), Response: projectAppSchema()},
		{Method: "GET", Path: "/api/project-apps/:id", Summary: "Get a project app", Tags: []string{projects}, Auth: AuthUser, Response: projectAppSchema()},
		{Method: "PUT", Path: "/api/project-apps/:id", Summary: "Update a project app", Tags: []string{projects}, Auth: AuthUser, Request: JSON(Obj("Update project app", F("name", Str("App name")), F("slug", Str("Slug")), F("favicon", Str("Favicon URL")))), Response: projectAppSchema()},
		{Method: "DELETE", Path: "/api/project-apps/:id", Summary: "Delete a project app", Tags: []string{projects}, Auth: AuthUser, Response: StatusSchema("deleted")},
		{Method: "GET", Path: "/api/projects/:id/stacks", Summary: "List compose stacks", Tags: []string{stacks}, Auth: AuthUser, Response: Arr("Stacks", composeStackSchema())},
		{Method: "POST", Path: "/api/projects/:id/stacks/review", Summary: "Review a compose stack", Tags: []string{stacks}, Auth: AuthUser, Request: JSON(composeStackRequestSchema()), Response: composeReviewSchema()},
		{Method: "POST", Path: "/api/projects/:id/stacks", Summary: "Save a compose stack", Tags: []string{stacks}, Auth: AuthUser, Request: JSON(composeStackRequestSchema()), Response: composeStackSchema()},
		{Method: "POST", Path: "/api/projects/:id/stacks/:stackId/deploy", Summary: "Deploy a compose stack", Tags: []string{stacks}, Auth: AuthUser, Code: 202, Response: EmptySchema()},
		{Method: "POST", Path: "/api/projects/:id/stacks/:stackId/cancel", Summary: "Cancel a stack deployment", Tags: []string{stacks}, Auth: AuthUser, Code: 202, Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:id/stacks/:stackId/config", Summary: "Get a stack config", Tags: []string{stacks}, Auth: Scope("env:read"), Response: Obj("Stack config", RF("stack", composeStackSchema()), F("config", Str("Resolved config")))},
	}
}
