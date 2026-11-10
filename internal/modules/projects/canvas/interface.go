package canvas

import "context"

// Vault provides decryption/encryption for environment secrets.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// EnvironmentRepository queries environment configs for canvas views.
type EnvironmentRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]*EnvironmentConfig, error)
	GetDefaultByProject(ctx context.Context, projectID string) (*EnvironmentConfig, error)
}

// ClusterDataLister resolves cluster instances by project.
type ClusterDataLister interface {
	ListByProject(ctx context.Context, project string) ([]ClusterData, error)
}

// CanvasRepository persists and queries canvas topology summaries.
type CanvasRepository interface {
	ListCanvasSummaries(ctx context.Context, organizationID string) ([]CanvasSummary, error)
	GetCanvasSummary(ctx context.Context, id string) (*CanvasSummary, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error)
}
