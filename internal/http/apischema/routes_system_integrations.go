package apischema

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

func chatAppConnectionSchema() Schema {
	return Obj("Chat connection",
		RF("id", Str("Connection ID")),
		RF("userId", Str("User ID")),
		RF("provider", Str("Provider")),
		RF("transport", Str("Transport")),
		RF("purpose", Str("Purpose")),
		RF("name", Str("Name")),
		F("botUsername", Str("Bot username")),
		F("config", Any("Config")),
		RF("enabled", Bool("Enabled")),
		RF("status", Str("Status")),
		F("lastError", Str("Last error")),
		F("lastEventAt", DateTime("Last event timestamp")),
		RF("createdAt", DateTime("Created timestamp")),
		RF("updatedAt", DateTime("Updated timestamp")),
	)
}

func chatAppIdentitySchema() Schema {
	return Obj("Chat identity",
		RF("id", Str("Identity ID")),
		RF("connectionId", Str("Connection ID")),
		RF("userId", Str("User ID")),
		RF("provider", Str("Provider")),
		RF("externalUserId", Str("External user ID")),
		F("externalUsername", Str("External username")),
		RF("role", Str("Role")),
		F("permissions", Arr("Permissions", Str("Permission"))),
		RF("active", Bool("Active")),
		RF("pairedAt", DateTime("Paired timestamp")),
		RF("createdAt", DateTime("Created timestamp")),
		RF("updatedAt", DateTime("Updated timestamp")),
	)
}

func chatAppChallengeSchema() Schema {
	return Obj("Chat challenge",
		RF("id", Str("Challenge ID")),
		RF("connectionId", Str("Connection ID")),
		RF("provider", Str("Provider")),
		RF("externalUserId", Str("External user ID")),
		F("externalChannelId", Str("External channel ID")),
		RF("token", Str("Token")),
		RF("code", Str("Code")),
		RF("expiresAt", DateTime("Expires timestamp")),
		F("verifiedAt", DateTime("Verified timestamp")),
		RF("createdAt", DateTime("Created timestamp")),
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

func advisorySchema() Schema {
	return Obj("Security advisory",
		RF("id", Str("ID")),
		RF("revision", Int("Revision")),
		RF("title", Str("Title")),
		RF("severity", Str("Severity")),
		RF("message", Str("Message")),
		RF("url", Str("URL")),
		RF("action", Str("Action")),
	)
}

func projectSecuritySchema() Schema {
	return Obj("Project security",
		RF("projectId", Str("Project ID")),
		F("deploymentId", Str("Deployment ID")),
		F("checkedAt", DateTime("Checked at")),
		F("autoRefresh", Bool("Auto refresh")),
		F("findings", Arr("Findings", Obj("Finding",
			RF("advisoryId", Str("Advisory ID")),
			RF("title", Str("Title")),
			RF("severity", Str("Severity")),
			RF("action", Str("Action")),
			RF("confidence", Str("Confidence")),
			F("serviceName", Str("Service name")),
			F("evidence", Str("Evidence")),
			F("canRefresh", Bool("Can refresh")),
			F("runtimeUpgradeRequired", Bool("Runtime upgrade required")),
			F("imageChangeRequired", Bool("Image change required")),
		))),
	)
}

func advisoriesResponseSchema() Schema {
	return Obj("Advisories response",
		RF("revision", Int("Revision")),
		F("checkedAt", DateTime("Checked at")),
		F("expiresAt", DateTime("Expires at")),
		RF("stale", Bool("Stale")),
		RF("advisories", Arr("Advisories", advisorySchema())),
		RF("affectedProjects", Arr("Affected projects", projectSecuritySchema())),
	)
}
