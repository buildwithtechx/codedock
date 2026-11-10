package templates

import "context"

// OneClickOperations defines the high-level application catalog operations.
type OneClickOperations interface {
	ListApps() any
	GetApp(id string) (any, error)
	AddCustomApp(tmpl any) (any, error)
	DeleteCustomApp(id string) error
	Review(ctx context.Context, req any) (any, error)
	Deploy(ctx context.Context, req any) (any, error)
}

// EnvironmentRepository retrieves environment configurations for a project.
type EnvironmentRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]*EnvironmentConfig, error)
}
