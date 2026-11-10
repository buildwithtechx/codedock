package migrations

import (
	"context"
)

type LifecycleService interface {
	ListSources(ctx context.Context, userID, orgID string) (any, error)
	TestSource(ctx context.Context, userID, orgID string, req CreateMigrationSourceRequest) (string, error)
	CreateSource(ctx context.Context, userID, orgID string, req CreateMigrationSourceRequest) (any, error)
	DeleteSource(ctx context.Context, userID, orgID, sourceID string) error
	Scan(ctx context.Context, userID, orgID, sourceID string) (any, error)
	ScanStream(ctx context.Context, userID, orgID, sourceID string, send func(event string, data any) error) error
	RevealEnv(ctx context.Context, userID, orgID, sourceID, containerID string) (map[string]string, error)
	Adopt(ctx context.Context, userID, orgID string, req AdoptMigrationRequest) (any, error)
	Reimport(ctx context.Context, userID, orgID string, req ReimportMigrationRequest) (any, error)
	RepoCompose(ctx context.Context, userID, orgID string, req RepoComposeMigrationRequest) (any, error)
	Preview(ctx context.Context, userID, orgID string, req PreviewMigrationRequest) (any, error)
	StartMigration(ctx context.Context, userID, orgID string, req StartMigrationRequest) (any, error)
	StartProjectMove(ctx context.Context, userID, orgID string, req StartProjectMoveRequest) (any, error)
	Active(ctx context.Context, userID, orgID string) (any, error)
	ListRuns(ctx context.Context, userID, orgID string) (any, error)
	GetRun(ctx context.Context, userID, runID string) (any, error)
	StreamRun(ctx context.Context, userID, runID string, send func(event string, data any) error) error
	Cutover(ctx context.Context, userID, runID string, req CutoverMigrationRequest) (any, error)
	Cancel(ctx context.Context, userID, runID string) error
	Respond(ctx context.Context, userID, runID string, req RespondMigrationRequest) (any, error)
	Resume(ctx context.Context, userID, runID string, req ResumeMigrationRequest) (any, error)
	CleanupTarget(ctx context.Context, userID, runID string) error
	DeleteRun(ctx context.Context, userID, runID string) error
}

type BundleService interface {
	Export(ctx context.Context, passphrase string) ([]byte, error)
	Import(ctx context.Context, data []byte, passphrase string) (any, error)
}

type UserCounter interface {
	CountUsers(ctx context.Context) (int, error)
}

// MigrationRepository defines persistence for migration source systems, runs, and logs.
type MigrationRepository interface {
	CreateSource(ctx context.Context, source *MigrationSource) error
	GetSource(ctx context.Context, id string) (*MigrationSource, error)
	ListSourcesByOrg(ctx context.Context, orgID string) ([]*MigrationSource, error)
	DeleteSource(ctx context.Context, id string) error
	CreateRun(ctx context.Context, run *MigrationRun) error
	GetRun(ctx context.Context, id string) (*MigrationRun, error)
	UpdateRun(ctx context.Context, run *MigrationRun) error
	UpdateRunStatus(ctx context.Context, id string, status MigrationStatus, phase, errMsg string) error
	AppendRunLogs(ctx context.Context, id, logs string) error
	SetRunPrompt(ctx context.Context, id, prompt string) error
	RequestCancel(ctx context.Context, id string) error
	ListRunsByOrg(ctx context.Context, orgID string, limit int) ([]*MigrationRun, error)
	ListRunsBySource(ctx context.Context, sourceID string, limit int) ([]*MigrationRun, error)
	ActiveRunForSource(ctx context.Context, sourceID string) (*MigrationRun, error)
	ListActive(ctx context.Context) ([]*MigrationRun, error)
	DeleteRun(ctx context.Context, id string) error
}
