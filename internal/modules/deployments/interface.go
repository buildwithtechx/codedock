package deployments

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-jwt/jwt/v5"
)

// DeploymentRepository defines persistence for container deployment histories and build artifacts.
type DeploymentRepository interface {
	Create(ctx context.Context, d *Deployment) error
	GetByID(ctx context.Context, id string) (*Deployment, error)
	ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*Deployment, int, error)
	ListByOrganization(ctx context.Context, filter DeploymentListFilter) ([]DeploymentListItem, int, error)
	Update(ctx context.Context, d *Deployment) error
	UpdateStatus(ctx context.Context, id string, status DeploymentStatus, buildLogs, containerID, resourceIssue string) error
	SetPinned(ctx context.Context, id string, pinned bool) error
	TogglePinned(ctx context.Context, id string) (bool, error)
	ListForRetention(ctx context.Context, projectID string) ([]*Deployment, error)
	SetArtifactRetainedAt(ctx context.Context, id string, at *time.Time) error
	SetRestartLoop(ctx context.Context, id string, finding string) error
	RecoverInterrupted(ctx context.Context) error
}

// RegistryRepository manages registry connection credentials.
type RegistryRepository interface {
	Create(ctx context.Context, registry *Registry) error
	ListByProject(ctx context.Context, projectID string) ([]*Registry, error)
	Get(ctx context.Context, id string) (*Registry, error)
	Delete(ctx context.Context, id string) error
}

// AIExplainer generates plain-language explanations of failed builds and runtime logs.
type AIExplainer interface {
	ExplainDeploymentFailure(ctx context.Context, deploymentID string) (any, error)
}

// PRPreviewLookup retrieves ephemeral preview metadata.
type PRPreviewLookup interface {
	ListByApp(ctx context.Context, appID string) (any, error)
}

// AuditLogger records deployment lifecycle transitions.
type AuditLogger interface {
	LogAction(ctx any, action, resource, ip string, details map[string]string)
}

// AppServiceLookup queries application service records.
type AppServiceLookup interface {
	GetAppService(ctx context.Context, id string) (*AppServiceSummary, error)
	UpdateAppContainerID(ctx context.Context, id string, containerID string) error
}

// ProjectLookup resolves project permissions and ownership.
type ProjectLookup interface {
	HasPermission(ctx context.Context, projectID, userID, role, action string) bool
}

// CodeSyncProvider fetches source trees from Git providers.
type CodeSyncProvider interface {
	SyncCodebase(ctx context.Context, repoURL, branch, commitHash, targetDir string, logWriter io.Writer) error
}

// DeployEngine compiles Dockerfiles or Nixpacks into runnable containers.
type DeployEngine interface {
	BeginServiceOperation(ctx context.Context, serviceID string) (context.Context, func(), error)
	DeployApp(ctx context.Context, app *AppServiceSummary, sourceDir string, logWriter io.Writer) (string, error)
	StopApp(ctx context.Context, containerID string) error
}

// TokenValidator parses and checks bearer terminal session tokens.
type TokenValidator interface {
	ValidateToken(token string) (jwt.MapClaims, error)
}

// UserStatusLookup verifies user active state.
type UserStatusLookup interface {
	IsUserActive(ctx context.Context, userID string) (bool, error)
}

// TerminalDockerTargets connects active terminal sessions to Docker engine daemons.
type TerminalDockerTargets interface {
	DockerForService(ctx context.Context, serviceID string) (*client.Client, func(), error)
}
