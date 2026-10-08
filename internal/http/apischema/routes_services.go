package apischema

func variableSchema() Schema {
	props := []Property{
		RF("id", ID("Variable ID")),
		RF("serviceId", Str("Service ID")),
		RF("projectId", Str("Project ID")),
		RF("environmentId", Str("Environment ID")),
		RF("key", Str("Variable key")),
		F("value", Str("Variable value")),
		F("isSecret", Bool("Secret variable")),
	}
	return Obj("Service variable", append(props, Timestamps()...)...)
}

func routeRuleSchema() Schema {
	props := []Property{
		RF("id", ID("Rule ID")),
		RF("serviceId", Str("Service ID")),
		RF("name", Str("Rule name")),
		F("enabled", Bool("Whether the rule is enabled")),
		RF("ruleType", Str("Rule type")),
	}
	return Obj("Route rule", append(props, Timestamps()...)...)
}

func autoscalingSchema() Schema {
	return Obj("Autoscaling policy",
		F("supported", Bool("Whether autoscaling is supported")),
		F("serviceId", Str("Service ID")),
		F("enabled", Bool("Whether autoscaling is enabled")),
		F("minReplicas", Int("Minimum replicas")),
		F("maxReplicas", Int("Maximum replicas")),
		F("scaleUpCpu", Num("Scale-up CPU threshold")),
		F("scaleDownCpu", Num("Scale-down CPU threshold")),
		F("cooldownSeconds", Int("Cooldown seconds")),
		F("lastDecision", Str("Last scaling decision")),
		F("lastEvaluatedAt", Str("Last evaluation time")),
	)
}

func serverlessCodeSchema() Schema {
	return Obj("Serverless function code",
		RF("id", ID("Function ID")),
		RF("serviceId", Str("Service ID")),
		RF("runtime", Str("Runtime")),
		RF("codeContent", Str("Function source")),
	)
}

func domainVerifySchema() Schema {
	return Obj("Domain verification",
		RF("domainId", Str("Domain ID")),
		RF("domainName", Str("Domain name")),
		RF("verified", Bool("Verified")),
		RF("status", Str("Status")),
		F("resolvedIp", Str("Resolved IP")),
		F("serverIp", Str("Server IP")),
		F("message", Str("Message")),
	)
}

func serviceOperations() []Operation {
	variables := "variables"
	routing := "routing"
	scaling := "scaling"
	serverless := "serverless"
	domains := "domains"
	return []Operation{
		{Method: "GET", Path: "/api/services/:serviceId/variables", Summary: "List service variables", Tags: []string{variables}, Auth: ServiceRole(""), Response: Arr("Variables", variableSchema())},
		{Method: "POST", Path: "/api/services/:serviceId/variables", Summary: "Create a service variable", Tags: []string{variables}, Auth: ServiceRole("admin"), Code: 201, Request: JSON(Obj("Create variable", RF("key", Str("Variable key")), RF("value", Str("Variable value")), F("isSecret", Bool("Secret variable")))), Response: variableSchema()},
		{Method: "PUT", Path: "/api/services/:serviceId/variables/:id", Summary: "Update a service variable", Tags: []string{variables}, Auth: ServiceRole("admin"), Request: JSON(Obj("Update variable", RF("key", Str("Variable key")), RF("value", Str("Variable value")), F("isSecret", Bool("Secret variable")))), Response: variableSchema()},
		{Method: "DELETE", Path: "/api/services/:serviceId/variables/:id", Summary: "Delete a service variable", Tags: []string{variables}, Auth: ServiceRole("admin"), Response: EmptySchema()},
		{Method: "GET", Path: "/api/services/:serviceId/env-suggestions", Summary: "Suggest environment variables", Tags: []string{variables}, Auth: ServiceRole(""), Response: Arr("Suggestions", Any("Suggestion"))},
		{Method: "GET", Path: "/api/services/:serviceId/route-rules", Summary: "List route rules", Tags: []string{routing}, Auth: ServiceRole(""), Response: Arr("Rules", routeRuleSchema())},
		{Method: "POST", Path: "/api/services/:serviceId/route-rules", Summary: "Create a route rule", Tags: []string{routing}, Auth: ServiceRole("admin"), Raw: true, Code: 201, Request: JSON(Obj("Create rule", RF("name", Str("Rule name")), RF("ruleType", Str("Rule type")), RF("spec", Any("Rule spec")), F("enabled", Bool("Enabled")))), Response: Obj("Created rule", RF("success", Bool("Success")), RF("rule", routeRuleSchema()))},
		{Method: "PATCH", Path: "/api/services/:serviceId/route-rules/:ruleId", Summary: "Update a route rule", Tags: []string{routing}, Auth: ServiceRole("admin"), Request: JSON(Obj("Update rule", F("name", Str("Rule name")), F("enabled", Bool("Enabled")), F("spec", Any("Rule spec")))), Response: routeRuleSchema()},
		{Method: "DELETE", Path: "/api/services/:serviceId/route-rules/:ruleId", Summary: "Delete a route rule", Tags: []string{routing}, Auth: ServiceRole("admin"), Response: EmptySchema()},
		{Method: "GET", Path: "/api/services/:serviceId/autoscaling", Summary: "Get the autoscaling policy", Tags: []string{scaling}, Auth: ServiceRole(""), Response: autoscalingSchema()},
		{Method: "PUT", Path: "/api/services/:serviceId/autoscaling", Summary: "Save the autoscaling policy", Tags: []string{scaling}, Auth: ServiceRole("admin"), Request: JSON(autoscalingSchema()), Response: autoscalingSchema()},
		{Method: "GET", Path: "/api/services/:serviceId/serverless/code", Summary: "Get serverless function code", Tags: []string{serverless}, Auth: ServiceRole(""), Response: Obj("Function", RF("code", serverlessCodeSchema()))},
		{Method: "POST", Path: "/api/services/:serviceId/serverless/code", Summary: "Save serverless function code", Tags: []string{serverless}, Auth: ServiceRole("admin"), Request: JSON(Obj("Save code", RF("runtime", Str("Runtime")), RF("codeContent", Str("Function source")))), Response: Obj("Function", RF("code", serverlessCodeSchema()))},
		{Method: "GET", Path: "/api/services/:id/domains", Summary: "List service domains", Tags: []string{domains}, Auth: AuthUser, Response: Arr("Domains", DomainSchema())},
		{Method: "POST", Path: "/api/services/:id/domains", Summary: "Attach a domain", Tags: []string{domains}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Attach domain", RF("domainName", Str("Domain name")), F("redirectTo", Str("Redirect target")), F("pathPrefix", Str("Path prefix")))), Response: DomainSchema()},
		{Method: "GET", Path: "/api/domains", Summary: "List all domains", Tags: []string{domains}, Auth: AuthUser, Response: Arr("Domains", DomainSchema())},
		{Method: "GET", Path: "/api/domains/:id/verify", Summary: "Verify a domain", Tags: []string{domains}, Auth: AuthUser, Response: domainVerifySchema()},
		{Method: "POST", Path: "/api/domains/:id/verify", Summary: "Verify a domain", Tags: []string{domains}, Auth: AuthUser, Response: domainVerifySchema()},
		{Method: "DELETE", Path: "/api/domains/:id", Summary: "Delete a domain", Tags: []string{domains}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
	}
}
