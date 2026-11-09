package networking

type RouteRuleType string

const (
	RouteRuleTypeRateLimit   RouteRuleType = "rate_limit"
	RouteRuleTypeIPAllowlist RouteRuleType = "ip_allowlist"
	RouteRuleTypeHeaders     RouteRuleType = "headers"
)

type RouteRule struct {
	ID        string        `json:"id"`
	ServiceID string        `json:"serviceId"`
	Name      string        `json:"name"`
	RuleType  RouteRuleType `json:"ruleType"`
	SpecJSON  string        `json:"specJson"`
	Priority  int           `json:"priority"`
	Enabled   bool          `json:"enabled"`
}

type RateLimitSpec struct {
	Average int    `json:"average"`
	Burst   int    `json:"burst"`
	Period  string `json:"period"`
}

type IPListSpec struct {
	CIDRs []string `json:"cidrs"`
}

type HeadersSpec struct {
	Set    map[string]string `json:"set,omitempty"`
	Remove []string          `json:"remove,omitempty"`
}
