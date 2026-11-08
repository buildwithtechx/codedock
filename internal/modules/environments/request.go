package environments

type CreateEnvironmentRequest struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type RouteProtectionInput struct {
	Hostname string `json:"hostname"`
	Action   string `json:"action"` // "protect", "revoke", "retry", "unlock"
	Password string `json:"password,omitempty"`
}

type CreateRouteRuleRequest struct {
	Name     string        `json:"name"`
	RuleType RouteRuleType `json:"ruleType"`
	Spec     any           `json:"spec"`
	Enabled  *bool         `json:"enabled,omitempty"`
}

type UpdateRouteRuleRequest struct {
	Name     *string        `json:"name,omitempty"`
	Enabled  *bool          `json:"enabled,omitempty"`
	Spec     any            `json:"spec,omitempty"`
	RuleType *RouteRuleType `json:"ruleType,omitempty"`
}

type EdgeLoginRequest struct {
	Password string `json:"password" form:"password"`
	Host     string `json:"host,omitempty" form:"host"`
	Path     string `json:"path,omitempty" form:"path"`
}

type SetEnvironmentVarRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
