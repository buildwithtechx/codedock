package apischema

func mailServerSchema() Schema {
	return Obj("Mail server",
		RF("id", ID("Mail server ID")),
		RF("organizationId", Str("Organization ID")),
		F("serverId", Str("Codedock server ID")),
		RF("domain", Str("Primary mail domain")),
		RF("maildUrl", Str("maild base URL")),
		F("webmailUrl", Str("Webmail URL")),
		F("createdAt", DateTime("Creation time")),
		F("updatedAt", DateTime("Update time")),
	)
}

func mailDomainSchema() Schema {
	return Obj("Mail domain",
		RF("id", ID("Domain ID")),
		RF("name", Str("Domain name")),
		F("createdAt", DateTime("Creation time")),
		F("updatedAt", DateTime("Update time")),
	)
}

func mailboxSchema() Schema {
	return Obj("Mailbox",
		RF("id", ID("Mailbox ID")),
		RF("domainId", Str("Domain ID")),
		RF("localPart", Str("Local part")),
		RF("address", Str("Email address")),
		F("createdAt", DateTime("Creation time")),
		F("updatedAt", DateTime("Update time")),
	)
}

func mailboxCredentialsSchema() Schema {
	return Obj("Mailbox credentials",
		RF("id", ID("Mailbox ID")),
		RF("address", Str("Email address")),
		F("password", Str("One-time mailbox password")),
	)
}

func dkimKeySchema() Schema {
	return Obj("DKIM key",
		RF("id", ID("DKIM key ID")),
		RF("domainId", Str("Domain ID")),
		RF("selector", Str("Selector")),
		RF("algorithm", Str("Algorithm")),
		F("publicKey", Str("Public key")),
		F("dnsRecord", Str("DNS record value")),
		F("createdAt", DateTime("Creation time")),
	)
}

func mailDNSRecordSchema() Schema {
	return Obj("Mail DNS record",
		RF("type", Str("Record type")),
		RF("name", Str("Record name")),
		RF("value", Str("Record value")),
		F("source", Str("Record source")),
	)
}

func mailCheckSchema() Schema {
	return Obj("Mail check",
		RF("name", Str("Check name")),
		RF("passed", Bool("Check passed")),
		F("detail", Str("Check detail")),
	)
}

func mailCertSchema() Schema {
	return Obj("Mail certificate",
		RF("id", ID("Certificate ID")),
		RF("hostname", Str("Hostname")),
		RF("expiresAt", DateTime("Expiry time")),
		F("createdAt", DateTime("Creation time")),
		F("updatedAt", DateTime("Update time")),
	)
}

func mailOperations() []Operation {
	mail := "mail"
	return []Operation{
		{Method: "GET", Path: "/api/mail/servers", Summary: "List mail servers", Tags: []string{mail}, Auth: AuthUser, Response: Arr("Mail servers", mailServerSchema())},
		{Method: "POST", Path: "/api/mail/servers", Summary: "Register a mail server", Tags: []string{mail}, Auth: AuthOrgAdmin, Code: 201, Request: JSON(Obj("Register mail server", F("organizationId", Str("Organization ID")), F("serverId", Str("Codedock server ID")), RF("domain", Str("Primary mail domain")), RF("maildUrl", Str("maild base URL")), F("webmailUrl", Str("Webmail URL")), RF("apiToken", Str("maild API token")))), Response: mailServerSchema()},
		{Method: "GET", Path: "/api/mail/servers/:id", Summary: "Get a mail server", Tags: []string{mail}, Auth: AuthUser, Response: mailServerSchema()},
		{Method: "DELETE", Path: "/api/mail/servers/:id", Summary: "Remove a mail server", Tags: []string{mail}, Auth: AuthOrgAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/mail/servers/:id/domains", Summary: "List mail domains", Tags: []string{mail}, Auth: AuthUser, Response: Arr("Mail domains", mailDomainSchema())},
		{Method: "POST", Path: "/api/mail/servers/:id/domains", Summary: "Create a mail domain", Tags: []string{mail}, Auth: AuthOrgAdmin, Code: 201, Request: JSON(Obj("Create mail domain", RF("name", Str("Domain name")))), Response: mailDomainSchema()},
		{Method: "DELETE", Path: "/api/mail/servers/:id/domains/:domainId", Summary: "Delete a mail domain", Tags: []string{mail}, Auth: AuthOrgAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/mail/servers/:id/domains/:domainId/mailboxes", Summary: "List mailboxes", Tags: []string{mail}, Auth: AuthUser, Response: Arr("Mailboxes", mailboxSchema())},
		{Method: "POST", Path: "/api/mail/servers/:id/domains/:domainId/mailboxes", Summary: "Create a mailbox", Tags: []string{mail}, Auth: AuthOrgAdmin, Code: 201, Request: JSON(Obj("Create mailbox", RF("localPart", Str("Local part")))), Response: mailboxCredentialsSchema()},
		{Method: "DELETE", Path: "/api/mail/servers/:id/domains/:domainId/mailboxes/:mailboxId", Summary: "Delete a mailbox", Tags: []string{mail}, Auth: AuthOrgAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/mail/servers/:id/domains/:domainId/dkim", Summary: "Create a DKIM key", Tags: []string{mail}, Auth: AuthOrgAdmin, Code: 201, Request: JSON(Obj("Create DKIM key", F("selector", Str("Selector")), F("algorithm", Str("Algorithm")))), Response: dkimKeySchema()},
		{Method: "GET", Path: "/api/mail/servers/:id/domains/:domainId/dns", Summary: "Get DNS instructions", Tags: []string{mail}, Auth: AuthUser, Response: Obj("DNS instructions", RF("domain", Str("Domain name")), RF("records", Arr("Records", mailDNSRecordSchema())), F("engineError", Str("Engine error")))},
		{Method: "GET", Path: "/api/mail/servers/:id/domains/:domainId/checks", Summary: "Run deliverability checks", Tags: []string{mail}, Auth: AuthUser, Response: Obj("Check report", RF("domain", Str("Domain name")), RF("checks", Arr("Checks", mailCheckSchema())))},
		{Method: "GET", Path: "/api/mail/servers/:id/certs", Summary: "List tracked certificates", Tags: []string{mail}, Auth: AuthUser, Response: Arr("Certificates", mailCertSchema())},
		{Method: "POST", Path: "/api/mail/servers/:id/certs/probe", Summary: "Probe a certificate", Tags: []string{mail}, Auth: AuthOrgAdmin, Request: JSON(Obj("Probe certificate", RF("host", Str("Hostname")), F("port", Int("Port")))), Response: mailCertSchema()},
	}
}
