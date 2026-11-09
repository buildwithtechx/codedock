package dns

type IngressDiscoverRequest struct {
	CredentialID    string  `json:"credentialId"`
	AccountID       *string `json:"accountId,omitempty"`
	IncludePolicies *bool   `json:"includePolicies,omitempty"`
}

type IngressConnectRequest struct {
	ServerID              string  `json:"serverId"`
	CredentialID          string  `json:"credentialId"`
	AccountID             string  `json:"accountId"`
	Name                  string  `json:"name"`
	TunnelID              *string `json:"tunnelId,omitempty"`
	Connector             string  `json:"connector"`
	EdgePort              *int    `json:"edgePort,omitempty"`
	DefaultForDomains     bool    `json:"defaultForDomains"`
	DefaultAccessPolicyID *string `json:"defaultAccessPolicyId,omitempty"`
}

type IngressUpdateRequest struct {
	DefaultForDomains     bool    `json:"defaultForDomains"`
	DefaultAccessPolicyID *string `json:"defaultAccessPolicyId"`
}

type IngressDomainRequest struct {
	DomainID       string  `json:"domainId,omitempty"`
	ConnectionID   *string `json:"connectionId"`
	AccessPolicyID *string `json:"accessPolicyId,omitempty"`
}
