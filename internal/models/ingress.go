package models

type IngressConnection struct {
	ID                    string  `json:"id" db:"id"`
	ServerID              string  `json:"serverId" db:"server_id"`
	CredentialID          string  `json:"credentialId" db:"credential_id"`
	AccountID             string  `json:"accountId" db:"account_id"`
	TunnelID              *string `json:"tunnelId" db:"tunnel_id"`
	Name                  string  `json:"name" db:"name"`
	Connector             string  `json:"connector" db:"connector"`
	EdgePort              int     `json:"edgePort" db:"edge_port"`
	ManagedTunnel         bool    `json:"managedTunnel" db:"managed_tunnel"`
	DefaultForDomains     bool    `json:"defaultForDomains" db:"default_for_domains"`
	DefaultAccessPolicyID *string `json:"defaultAccessPolicyId" db:"default_access_policy_id"`
	Status                string  `json:"status" db:"status"`
	Error                 *string `json:"error" db:"error"`
	CreatedAt             string  `json:"createdAt" db:"created_at"`
	UpdatedAt             string  `json:"updatedAt" db:"updated_at"`
}

type IngressHostname struct {
	ID             string  `json:"id" db:"id"`
	DomainID       *string `json:"domainId" db:"domain_id"`
	ConnectionID   string  `json:"connectionId" db:"connection_id"`
	Hostname       string  `json:"hostname" db:"hostname"`
	AccessPolicyID *string `json:"accessPolicyId" db:"access_policy_id"`
	Status         string  `json:"status" db:"status"`
	Phase          string  `json:"phase" db:"phase"`
	Error          *string `json:"error" db:"error"`
	UpdatedAt      string  `json:"updatedAt" db:"updated_at"`
}

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

type IngressDiscoveryAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IngressDiscoveryTunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Remote bool   `json:"remote"`
	Status string `json:"status"`
}

type IngressDiscoveryPolicy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Decision string `json:"decision"`
}

type IngressDiscoveryResult struct {
	Accounts      []IngressDiscoveryAccount `json:"accounts"`
	Tunnels       []IngressDiscoveryTunnel  `json:"tunnels"`
	Policies      []IngressDiscoveryPolicy  `json:"policies"`
	PoliciesError *string                   `json:"policiesError,omitempty"`
}
