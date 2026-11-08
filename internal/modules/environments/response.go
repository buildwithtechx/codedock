package environments

type EnvironmentResponse struct {
	Environment *EnvironmentConfig `json:"environment"`
}

type EnvironmentListResponse struct {
	Environments []EnvironmentConfig `json:"environments"`
}

type RouteRuleResponse struct {
	ID        string        `json:"id"`
	ServiceID string        `json:"serviceId"`
	Name      string        `json:"name"`
	Enabled   bool          `json:"enabled"`
	RuleType  RouteRuleType `json:"ruleType"`
	Spec      any           `json:"spec"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
}

type RouteRuleListResponse struct {
	Rules []*RouteRuleResponse `json:"rules"`
}

type RouteProtectionListResponse struct {
	Status []RouteProtectionStatus `json:"status"`
}

type EdgeAuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	Hostname      string `json:"hostname,omitempty"`
	Message       string `json:"message,omitempty"`
}
