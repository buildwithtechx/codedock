package environments

import (
	"time"
)

type EnvironmentConfig struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"isDefault" db:"is_default"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type RouteProtection struct {
	ID              string    `json:"id" db:"id"`
	ProjectID       string    `json:"projectId" db:"project_id"`
	Hostname        string    `json:"hostname" db:"hostname"`
	Enabled         bool      `json:"enabled" db:"enabled"`
	Secret          *string   `json:"-" db:"secret"`
	Revision        string    `json:"revision" db:"revision"`
	AppliedRevision *string   `json:"appliedRevision,omitempty" db:"applied_revision"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type RouteProtectionPolicy struct {
	Version        string `json:"version"`
	Salt           string `json:"salt"`
	Verifier       string `json:"verifier"`
	SigningKey     string `json:"signingKey"`
	Iterations     int    `json:"iterations"`
	SessionSeconds int    `json:"sessionSeconds"`
}

type RouteProtectionStatus struct {
	Hostname string `json:"hostname"`
	Enabled  bool   `json:"enabled"`
	Status   string `json:"status"`
}

type RouteProtectionDraft struct {
	Hostname string `json:"hostname"`
	Action   string `json:"action"`
	Password string `json:"password"`
}

type RouteRule struct {
	ID        string        `json:"id" db:"id"`
	ServiceID string        `json:"serviceId" db:"service_id"`
	Name      string        `json:"name" db:"name"`
	Enabled   bool          `json:"enabled" db:"enabled"`
	RuleType  RouteRuleType `json:"ruleType" db:"rule_type"`
	SpecJSON  string        `json:"-" db:"spec_json"`
	CreatedAt time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time     `json:"updatedAt" db:"updated_at"`
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

type EnvExampleVariableSuggestion struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Label      string `json:"label"`
	SourcePath string `json:"sourcePath"`
}

type GitProviderConfig struct {
	AccessToken string `json:"accessToken"`
}
