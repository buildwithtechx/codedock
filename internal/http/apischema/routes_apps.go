package apischema

func serviceVolumeSchema() Schema {
	props := []Property{
		RF("id", ID("Volume ID")),
		RF("serviceId", Str("Service ID")),
		RF("hostPath", Str("Host path")),
		RF("containerPath", Str("Container path")),
	}
	return Obj("Service volume", append(props, Timestamps()...)...)
}

func logDrainSchema() Schema {
	props := []Property{
		RF("id", ID("Drain ID")),
		RF("serviceId", Str("Service ID")),
		RF("projectId", Str("Project ID")),
		RF("drainType", Str("Drain type")),
		RF("endpointUrl", Str("Endpoint URL")),
	}
	return Obj("Log drain", append(props, Timestamps()...)...)
}

func webhookSchema() Schema {
	props := []Property{
		RF("id", ID("Webhook ID")),
		RF("serviceId", Str("Service ID")),
		RF("url", Str("Webhook URL")),
		F("eventTypes", Arr("Events", Str("Event"))),
		F("includePrEnvironments", Bool("Include PR environments")),
	}
	return Obj("Webhook", append(props, Timestamps()...)...)
}

func runtimeTargetSchema() Schema {
	return Obj("Runtime target",
		RF("kind", StrEnum("Target kind", "docker", "cluster", "native")),
		F("clusterId", Str("Cluster ID")),
		F("nodeIds", Arr("Node IDs", Str("Node"))),
		F("imageRepository", Str("Image repository")),
		F("registryId", Str("Registry ID")),
		F("volumes", Arr("Volumes", Obj("Runtime volume", RF("name", Str("Name")), RF("mountPath", Str("Mount path")), F("sizeGiB", Int("Size GiB")), F("storageClass", Str("Storage class")), F("shared", Bool("Shared"))))),
	)
}

func serviceRuntimeSchema() Schema {
	return Obj("Service runtime",
		F("serviceId", Str("Service ID")),
		F("projectId", Str("Project ID")),
		F("target", runtimeTargetSchema()),
		F("revision", Int("Revision")),
		F("status", Str("Status")),
		F("error", Str("Error message")),
		F("updatedAt", Str("Last update")),
	)
}

func runtimePodSchema() Schema {
	return Obj("Runtime pod",
		RF("name", Str("Pod name")),
		F("node", Str("Node")),
		F("phase", Str("Phase")),
		F("ready", Bool("Ready")),
		F("restarts", Int("Restarts")),
		F("cpu", Num("CPU cores")),
		F("memoryBytes", Int64("Memory bytes")),
	)
}

func appOperations() []Operation {
	apps := "apps"
	runtime := "runtime"
	return []Operation{
		{Method: "GET", Path: "/api/apps", Summary: "List organization services", Tags: []string{apps}, Auth: AuthUser, Response: Arr("Services", AppServiceSchema())},
		{Method: "GET", Path: "/api/projects/:id/apps", Summary: "List project services", Tags: []string{apps}, Auth: ProjectRole("member"), Response: Arr("Services", AppServiceSchema())},
		{Method: "GET", Path: "/api/projects/:id/services", Summary: "List project services", Tags: []string{apps}, Auth: ProjectRole("member"), Response: Arr("Services", AppServiceSchema())},
		{Method: "GET", Path: "/api/environments/:id/apps", Summary: "List environment services", Tags: []string{apps}, Auth: AuthUser, Response: Arr("Services", AppServiceSchema())},
		{Method: "POST", Path: "/api/environments/:id/apps", Summary: "Create a service", Tags: []string{apps}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create service", RF("projectId", Str("Project ID")), RF("name", Str("Service name")), RF("repositoryUrl", Str("Repository URL")), F("branch", Str("Branch")), F("runtimeMode", Str("Runtime mode")), F("buildCommand", Str("Build command")), F("startCommand", Str("Start command")), F("internalPort", Int("Container port")), F("domain", Str("Domain")))), Response: AppServiceSchema()},
		{Method: "GET", Path: "/api/apps/:id", Summary: "Get a service", Tags: []string{apps}, Auth: ServiceRole(""), Response: AppServiceSchema()},
		{Method: "PUT", Path: "/api/apps/:id", Summary: "Update a service", Tags: []string{apps}, Auth: ServiceRole("admin"), Request: JSON(Obj("Update service", RF("name", Str("Service name")), RF("repositoryUrl", Str("Repository URL")), F("branch", Str("Branch")), F("runtimeMode", Str("Runtime mode")), F("buildCommand", Str("Build command")), F("startCommand", Str("Start command")), F("internalPort", Int("Container port")), F("domain", Str("Domain")), F("status", Str("Status")))), Response: AppServiceSchema()},
		{Method: "DELETE", Path: "/api/apps/:id", Summary: "Delete a service", Tags: []string{apps}, Auth: ServiceRole("owner"), Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/apps/:id/redeploy", Summary: "Redeploy a service", Tags: []string{apps}, Auth: ServiceRole("admin"), Code: 202, Response: DeploymentSchema()},
		{Method: "POST", Path: "/api/apps/:id/restart", Summary: "Restart a service", Tags: []string{apps}, Auth: ServiceRole("admin"), Response: AppServiceSchema()},
		{Method: "POST", Path: "/api/apps/:id/stop", Summary: "Stop a service", Tags: []string{apps}, Auth: ServiceRole("admin"), Response: AppServiceSchema()},
		{Method: "GET", Path: "/api/apps/:id/volumes", Summary: "List service volumes", Tags: []string{apps}, Auth: ServiceRole(""), Response: Arr("Volumes", serviceVolumeSchema())},
		{Method: "POST", Path: "/api/apps/:id/volumes", Summary: "Add a service volume", Tags: []string{apps}, Auth: ServiceRole("admin"), Request: JSON(Obj("Add volume", RF("hostPath", Str("Host path")), RF("containerPath", Str("Container path")))), Response: serviceVolumeSchema()},
		{Method: "DELETE", Path: "/api/apps/:id/volumes/:volumeId", Summary: "Delete a service volume", Tags: []string{apps}, Auth: ServiceRole("admin"), Response: EmptySchema()},
		{Method: "GET", Path: "/api/apps/:id/log-drains", Summary: "List log drains", Tags: []string{apps}, Auth: ServiceRole(""), Response: Arr("Drains", logDrainSchema())},
		{Method: "POST", Path: "/api/apps/:id/log-drains", Summary: "Create a log drain", Tags: []string{apps}, Auth: ServiceRole("admin"), Code: 201, Request: JSON(Obj("Create drain", RF("drainType", Str("Drain type")), RF("endpointUrl", Str("Endpoint URL")), F("authToken", Str("Auth token")))), Response: logDrainSchema()},
		{Method: "DELETE", Path: "/api/apps/:id/log-drains/:drainId", Summary: "Delete a log drain", Tags: []string{apps}, Auth: ServiceRole("admin"), Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/apps/:id/webhooks", Summary: "List webhooks", Tags: []string{apps}, Auth: ServiceRole(""), Response: Arr("Webhooks", webhookSchema())},
		{Method: "POST", Path: "/api/apps/:id/webhooks", Summary: "Create a webhook", Tags: []string{apps}, Auth: ServiceRole("admin"), Code: 201, Request: JSON(Obj("Create webhook", RF("url", Str("Webhook URL")), F("eventTypes", Arr("Events", Str("Event"))), F("includePrEnvironments", Bool("Include PR environments")))), Response: webhookSchema()},
		{Method: "DELETE", Path: "/api/apps/:id/webhooks/:webhookId", Summary: "Delete a webhook", Tags: []string{apps}, Auth: ServiceRole("admin"), Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/apps/:id/runtime", Summary: "Get the service runtime", Tags: []string{runtime}, Auth: ServiceRole(""), Response: serviceRuntimeSchema()},
		{Method: "POST", Path: "/api/apps/:id/runtime/review", Summary: "Review a runtime change", Tags: []string{runtime}, Auth: ServiceRole("admin"), Request: JSON(Obj("Review runtime", RF("target", runtimeTargetSchema()), F("revision", Int("Revision")))), Response: Any("Runtime review")},
		{Method: "POST", Path: "/api/apps/:id/runtime/apply", Summary: "Apply a runtime change", Tags: []string{runtime}, Auth: ServiceRole("admin"), Code: 202, Request: JSON(Obj("Apply runtime", RF("operationId", Str("Operation ID")), RF("confirmation", Str("Confirmation token")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/apps/:id/runtime/operation", Summary: "Get the runtime operation", Tags: []string{runtime}, Auth: ServiceRole("admin"), Response: operationSchema()},
		{Method: "GET", Path: "/api/apps/:id/runtime/observe", Summary: "Observe the runtime workload", Tags: []string{runtime}, Auth: Scope("logs:read"), Response: Obj("Workload observation", RF("kind", Str("Kind")), RF("status", Str("Status")), F("desired", Int("Desired replicas")), F("available", Int("Available replicas")), F("pods", Arr("Pods", runtimePodSchema())))},
		{Method: "GET", Path: "/api/apps/:id/runtime/logs", Summary: "Get runtime logs", Tags: []string{runtime}, Auth: Scope("logs:read"), Response: Obj("Runtime logs", RF("logs", Str("Log output")), F("metrics", Map("Metrics", Any("Value")))), Query: []Param{QueryParam("pod", "Pod name"), QueryParamInt("tail", "Tail lines")}},
		{Method: "GET", Path: "/api/apps/:id/runtime/logs/stream", Summary: "Stream runtime logs", Tags: []string{runtime}, Auth: Scope("logs:read"), Raw: true, Content: ContentSSE, Response: Str("Log events"), Query: []Param{QueryParam("pod", "Pod name")}},
		{Method: "GET", Path: "/api/apps/:id/runtime/graph", Summary: "Get the instance graph", Tags: []string{runtime}, Auth: Scope("logs:read"), Response: Any("Instance graph")},
		{Method: "GET", Path: "/api/apps/:id/runtime/resources", Summary: "Get the resource summary", Tags: []string{runtime}, Auth: ServiceRole(""), Response: Any("Resource summary")},
		{Method: "POST", Path: "/api/apps/:id/runtime/exec", Summary: "Execute a runtime command", Tags: []string{runtime}, Auth: Scope("terminal:write"), Request: JSON(Obj("Execute", RF("command", Arr("Command", Str("Argument"))), F("pod", Str("Pod name")))), Response: Any("Execution result")},
		{Method: "GET", Path: "/api/apps/:id/runtime/exec-terminal", Summary: "Open a runtime terminal socket", Tags: []string{runtime}, Auth: Scope("terminal:write"), Raw: true, Response: Str("WebSocket terminal session")},
	}
}
