package monitoring

type CreateWorkflowRequest struct {
	ProjectID   string `json:"projectId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Trigger     string `json:"trigger,omitempty"`
	YAML        string `json:"yaml,omitempty"`
}

type InspectDestinationRequest struct {
	ServerID string `json:"serverId"`
}

type CreateRunnerRequest struct {
	ServerID    string             `json:"serverId,omitempty"`
	CloudPoolID string             `json:"cloudPoolId,omitempty"`
	Kind        string             `json:"kind,omitempty"`
	Name        string             `json:"name"`
	Mode        string             `json:"mode,omitempty"`
	Image       string             `json:"image,omitempty"`
	Labels      []string           `json:"labels,omitempty"`
	Enabled     bool               `json:"enabled"`
	Config      ActionRunnerConfig `json:"config,omitempty"`
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

type UpsertSubscriptionRequest struct {
	Category string   `json:"category"`
	Channels []string `json:"channels"`
	Enabled  bool     `json:"enabled"`
}

type UpdateNotificationSettingsRequest struct {
	NotificationSettings
}

type TestNotificationRequest struct {
	Provider string `json:"provider"`
}

type CloudSupportCustomerInput struct {
	RequestID string               `json:"requestId"`
	Subject   string               `json:"subject"`
	Message   string               `json:"message"`
	Category  CloudSupportCategory `json:"category"`
}

type CloudSupportCustomerReply struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

type CloudSupportCustomerStatus struct {
	Status CloudSupportStatus `json:"status"`
}

type AttentionActionRequest struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type PathsToggleRequest struct {
	Enabled bool `json:"enabled"`
}

type CheckoutRequest struct {
	Plan string `json:"plan"`
}
