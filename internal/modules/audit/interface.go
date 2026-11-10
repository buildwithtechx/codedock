package audit

import "context"

// AuditLogRepository defines database operations for recording and searching audit trails.
type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, limit, offset int) ([]AuditLog, error)
	ListFiltered(ctx context.Context, filter AuditLogFilter, limit, offset int) ([]AuditLog, error)
	Facets(ctx context.Context, filter AuditLogFilter) (map[string]int, int, error)
	ListActors(ctx context.Context) ([]AuditActor, error)
	GetSettings(ctx context.Context, orgID string) (*AuditSettings, error)
	UpdateSettings(ctx context.Context, settings *AuditSettings) error
}

// AuditResourceSearcher resolves resource identifiers across heterogeneous modules.
type AuditResourceSearcher interface {
	SearchResourceIDs(ctx context.Context, term string) []string
}

// NameSearchable represents any module resource capable of being filtered by string label.
type NameSearchable interface {
	SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error)
}
