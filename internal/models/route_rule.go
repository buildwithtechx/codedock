package models

import "time"

type RouteRuleType string

const (
	RouteRuleTypeRateLimit   RouteRuleType = "rate_limit"
	RouteRuleTypeIPAllowlist RouteRuleType = "ip_allowlist"
	RouteRuleTypeIPBlocklist RouteRuleType = "ip_blocklist"
	RouteRuleTypeHeaders     RouteRuleType = "headers"
)

type RouteRule struct {
	ID        string        `json:"id"`
	ServiceID string        `json:"serviceId"`
	Name      string        `json:"name"`
	Enabled   bool          `json:"enabled"`
	RuleType  RouteRuleType `json:"ruleType"`
	SpecJSON  string        `json:"-"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
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

type CreateRouteRuleRequest struct {
	Name     string        `json:"name"`
	RuleType RouteRuleType `json:"ruleType"`
	Spec     any           `json:"spec"`
	Enabled  *bool         `json:"enabled,omitempty"`
}

type UpdateRouteRuleRequest struct {
	Name    *string `json:"name,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
	Spec    any     `json:"spec,omitempty"`
}
