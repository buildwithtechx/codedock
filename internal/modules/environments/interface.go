package environments

import "context"

// EnvironmentRepository persists project environments (production, staging, preview).
type EnvironmentRepository interface {
	Get(ctx context.Context, id string) (*EnvironmentConfig, error)
	ListByProject(ctx context.Context, projectID string) ([]EnvironmentConfig, error)
	Create(ctx context.Context, env *EnvironmentConfig) error
	Delete(ctx context.Context, id string) error
}

// ProjectVarRepository accesses project-scoped environment variables.
type ProjectVarRepository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}

// RouteProtectionRepository manages edge authentication, password protection, and bypass tokens.
type RouteProtectionRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]RouteProtection, error)
	Get(ctx context.Context, projectID, hostname string) (*RouteProtection, error)
	GetByHostname(ctx context.Context, hostname string) (*RouteProtection, error)
	Save(ctx context.Context, rp *RouteProtection) error
	MarkApplied(ctx context.Context, projectID, hostname, revision string) error
	Delete(ctx context.Context, projectID, hostname string) error
}

// RouteRuleRepository manages HTTP routing rules and path redirections.
type RouteRuleRepository interface {
	Create(ctx context.Context, rule *RouteRule) error
	GetByID(ctx context.Context, id string) (*RouteRule, error)
	ListByService(ctx context.Context, serviceID string) ([]*RouteRule, error)
	Update(ctx context.Context, id string, name *string, enabled *bool, specJSON *string) error
	Delete(ctx context.Context, id string) error
}
