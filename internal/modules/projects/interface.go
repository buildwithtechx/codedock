package projects

import (
	"context"

	"codedock/internal/modules/organizations"
)

// ProjectRepository defines persistence for projects and capacity settings.
type ProjectRepository interface {
	ListByOrganization(ctx context.Context, organizationID string, search string, limit, offset int) ([]ProjectConfig, int, error)
	ListAll(ctx context.Context, search string, limit, offset int) ([]ProjectConfig, int, error)
	Get(ctx context.Context, id string) (*ProjectConfig, error)
	GetByOrganization(ctx context.Context, id, organizationID string) (*ProjectConfig, error)
	CountByUser(ctx context.Context, userID string) (int, error)
	SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error)
	Create(ctx context.Context, p *ProjectConfig) error
	SetServer(ctx context.Context, projectID, serverID string) error
	Delete(ctx context.Context, id string) error
	UpdateRollbackSettings(ctx context.Context, projectID string, strategy *string, window *int) error
	UpdateCapacityMeasurement(ctx context.Context, projectID string, snapshotSizeBytes int64) error
}

// EnvRepository defines access to project environment variables.
type EnvRepository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}

// ProjectAppRepository defines persistence for application services bundled within projects.
type ProjectAppRepository interface {
	ListByOrganization(ctx context.Context, organizationID string) ([]*ProjectApp, error)
	GetByID(ctx context.Context, id string) (*ProjectApp, error)
	GetBySlug(ctx context.Context, organizationID, slug string) (*ProjectApp, error)
	Create(ctx context.Context, app *ProjectApp) error
	Update(ctx context.Context, app *ProjectApp) error
	Delete(ctx context.Context, id string) error
}

// ProjectSettingsRepository defines token creation and membership permissions for projects.
type ProjectSettingsRepository interface {
	CreateToken(ctx context.Context, t *ProjectToken, fullToken string) error
	ListTokensByProject(ctx context.Context, projectID string) ([]*ProjectToken, error)
	DeleteToken(ctx context.Context, id, projectID string) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
}

// OrganizationRepository checks user membership inside organizations.
type OrganizationRepository interface {
	GetMember(ctx context.Context, organizationID, userID string) (*organizations.OrganizationMember, error)
}

// UserOrgLookup resolves user organizations.
type UserOrgLookup interface {
	ListOrganizationsByUser(ctx context.Context, userID string) ([]*organizations.Organization, error)
	GetOrCreateDefaultOrganization(ctx context.Context, userID string) (*organizations.Organization, error)
}

// EnvironmentRepository retrieves environments configured for a project.
type EnvironmentRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]*EnvironmentConfig, error)
}
