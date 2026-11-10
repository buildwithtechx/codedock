package databases

import "context"

// Repository defines database persistence and permission checks for managed databases.
type Repository interface {
	Create(ctx context.Context, d *Database) error
	GetByID(ctx context.Context, id string) (*Database, error)
	ListAll(ctx context.Context, search string) ([]*Database, error)
	ListByProject(ctx context.Context, projectID string) ([]*Database, error)
	Update(ctx context.Context, d *Database) error
	UpdateStatus(ctx context.Context, id string, status DatabaseStatus, containerID string) error
	Delete(ctx context.Context, id string) error
	GetNextAvailablePort(ctx context.Context) (int, error)
	ProjectExists(ctx context.Context, projectID string) (bool, error)
	IsProjectMemberOrOwner(ctx context.Context, projectID, userID, userRole string) bool
	GetCustomDNSResolvers(ctx context.Context) string
	LogAction(ctx context.Context, entry AuditLogEntry)
}

// Deployer manages the container runtime lifecycle for database instances.
type Deployer interface {
	SpinUp(ctx context.Context, dbConfig *Database) (string, error)
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string, deleteRecord func() error) error
	ImportData(ctx context.Context, dbConfig *Database, sourceURL string) error
}

// VolumeGate arbitrates exclusive access to shared storage volumes.
type VolumeGate interface {
	AcquireVolume(name string) (func(), error)
}
