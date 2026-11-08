package mail

type CreateMailServerRequest struct {
	OrganizationID string `json:"organizationId"`
	ServerID       string `json:"serverId,omitempty"`
	Domain         string `json:"domain"`
	MaildURL       string `json:"maildUrl"`
	WebmailURL     string `json:"webmailUrl,omitempty"`
	APIToken       string `json:"apiToken"`
}

type CreateMailDomainRequest struct {
	Name string `json:"name"`
}

type CreateMailboxRequest struct {
	LocalPart string `json:"localPart"`
}

type CreateDKIMRequest struct {
	Selector  string `json:"selector"`
	Algorithm string `json:"algorithm"`
}

type ProbeMailCertRequest struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Fail2banActionRequest struct {
	IP   string `json:"ip"`
	Jail string `json:"jail,omitempty"`
}

type CreateMailAliasRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
