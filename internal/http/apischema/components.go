package apischema

func ID(desc string) Schema {
	return Str(desc)
}

func Timestamps() []Property {
	return []Property{F("createdAt", DateTime("Creation time")), F("updatedAt", DateTime("Last update time"))}
}

func UserSchema() Schema {
	props := []Property{
		RF("id", ID("User ID")),
		RF("email", Str("Email address")),
		RF("name", Str("Display name")),
		RF("role", StrEnum("Instance role", "owner", "admin", "member", "viewer")),
		F("isActive", Bool("Whether the account is active")),
		F("emailVerified", Bool("Whether the email is verified")),
		F("totpEnabled", Bool("Whether two-factor auth is enabled")),
		F("planType", Str("Billing plan")),
		F("projectsCount", Int("Owned project count")),
		F("servicesCount", Int("Owned service count")),
	}
	return Obj("User", append(props, Timestamps()...)...)
}

func AuthResponseSchema() Schema {
	return Obj("Authenticated session",
		RF("user", UserSchema()),
		RF("token", Str("Access JWT")),
		F("refreshToken", Str("Refresh token")),
	)
}

func PATSchema() Schema {
	props := []Property{
		RF("id", ID("Token ID")),
		RF("name", Str("Token name")),
		F("prefix", Str("Identifying prefix")),
		F("accessLevel", Str("Access level")),
		F("projectScope", Str("Project scope")),
		F("expiresAt", DateTime("Expiry time")),
	}
	return Obj("Personal access token", append(props, Timestamps()...)...)
}

func OrganizationSchema() Schema {
	props := []Property{RF("id", ID("Organization ID")), RF("name", Str("Organization name"))}
	return Obj("Organization", append(props, Timestamps()...)...)
}

func MemberSchema() Schema {
	return Obj("Organization membership",
		RF("id", ID("Membership ID")),
		RF("organizationId", Str("Organization ID")),
		F("userId", Str("User ID")),
		RF("email", Str("Member email")),
		RF("permission", StrEnum("Permission", "owner", "admin", "member", "viewer")),
		RF("status", StrEnum("Invite status", "invited", "active", "revoked")),
		F("invitedAt", DateTime("Invite time")),
		F("acceptedAt", DateTime("Acceptance time")),
	)
}

func ProjectSchema() Schema {
	props := []Property{
		RF("id", ID("Project ID")),
		F("organizationId", Str("Owning organization")),
		F("serverId", Str("Assigned server")),
		RF("name", Str("Project name")),
		F("slug", Str("URL slug")),
		F("description", Str("Description")),
		F("status", Str("Project status")),
		F("gitProvider", Str("Git provider")),
		F("gitOwner", Str("Repository owner")),
		F("gitRepo", Str("Repository name")),
		F("gitBranch", Str("Default branch")),
	}
	return Obj("Project", append(props, Timestamps()...)...)
}

func EnvironmentSchema() Schema {
	props := []Property{
		RF("id", ID("Environment ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Environment name")),
		F("isDefault", Bool("Whether this is the default environment")),
	}
	return Obj("Environment", append(props, Timestamps()...)...)
}

func AppServiceSchema() Schema {
	props := []Property{
		RF("id", ID("Service ID")),
		RF("projectId", Str("Project ID")),
		RF("environmentId", Str("Environment ID")),
		RF("name", Str("Service name")),
		F("repositoryUrl", Str("Source repository")),
		F("imageRef", Str("Container image reference")),
		F("branch", Str("Deploy branch")),
		F("runtimeMode", Str("Runtime mode")),
		F("buildCommand", Str("Build command")),
		F("startCommand", Str("Start command")),
		F("internalPort", Int("Container port")),
		F("domain", Str("Assigned domain")),
		F("status", Str("Service status")),
		F("containerId", Str("Running container ID")),
		F("enablePrPreviews", Bool("PR previews enabled")),
		F("maintenanceMode", Bool("Maintenance mode")),
	}
	return Obj("Application service", append(props, Timestamps()...)...)
}

func DeploymentSchema() Schema {
	props := []Property{
		RF("id", ID("Deployment ID")),
		RF("serviceId", Str("Service ID")),
		RF("projectId", Str("Project ID")),
		RF("environmentId", Str("Environment ID")),
		RF("status", Str("Deployment status")),
		F("branch", Str("Branch")),
		F("commitHash", Str("Commit hash")),
		F("commitMessage", Str("Commit message")),
		F("trigger", Str("Trigger source")),
		F("containerId", Str("Container ID")),
		F("finishedAt", DateTime("Finish time")),
	}
	return Obj("Deployment", append(props, Timestamps()...)...)
}

func DatabaseSchema() Schema {
	props := []Property{
		RF("id", ID("Database ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Database name")),
		RF("engine", Str("Database engine")),
		F("version", Str("Engine version")),
		F("username", Str("Username")),
		F("databaseName", Str("Database name on server")),
		F("internalHost", Str("Internal hostname")),
		F("internalDns", Str("Internal DNS")),
		F("externalDns", Str("External DNS")),
		F("port", Int("Container port")),
		F("externalPort", Int("Published port")),
		F("status", Str("Database status")),
		F("containerId", Str("Container ID")),
	}
	return Obj("Database", append(props, Timestamps()...)...)
}

func ServerSchema() Schema {
	props := []Property{
		RF("id", ID("Server ID")),
		RF("name", Str("Server name")),
		F("organizationId", Str("Organization ID")),
		F("ipAddress", Str("IP address")),
		F("isLocal", Bool("Local server")),
		F("sshHost", Str("SSH hostname")),
		F("sshPort", Int("SSH port")),
		F("sshUser", Str("SSH username")),
		F("sshAuthMethod", Str("SSH auth method")),
		F("isControlPlane", Bool("Control plane")),
		F("status", Str("Connection status")),
		F("provider", Str("Provider")),
		F("externalId", Str("External ID")),
		F("region", Str("Region")),
		F("serverType", Str("Server type")),
		F("lastSeenAt", DateTime("Last successful contact")),
	}
	return Obj("Server", append(props, Timestamps()...)...)
}

func DomainSchema() Schema {
	props := []Property{
		RF("id", ID("Domain ID")),
		RF("serviceId", Str("Service ID")),
		RF("domainName", Str("Domain name")),
		F("redirectTo", Str("Redirect target")),
		F("isCustom", Bool("Custom domain")),
		F("sslStatus", Str("Certificate status")),
		F("verified", Bool("DNS verified")),
	}
	return Obj("Domain", append(props, Timestamps()...)...)
}

func StatusSchema(status string) Schema {
	return Obj("Status acknowledgement", RF("status", StrEnum("Result", status)))
}

func OkSchema() Schema {
	return Obj("Acknowledgement", RF("ok", Bool("Success")))
}

func EmptySchema() Schema {
	return Obj("Empty result", F("result", Str("Result detail")))
}
