package adapters

import (
	"context"

	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/modules/users"
)

type MigrationLifecycleAdapter struct {
	Service modmigrations.LifecycleService
}

func (a *MigrationLifecycleAdapter) ListSources(ctx context.Context, userID, orgID string) (any, error) {
	if a.Service == nil {
		return []any{}, nil
	}
	return a.Service.ListSources(ctx, userID, orgID)
}

func (a *MigrationLifecycleAdapter) TestSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (string, error) {
	if a.Service == nil {
		return "", nil
	}
	return a.Service.TestSource(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) CreateSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.CreateSource(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) DeleteSource(ctx context.Context, userID, orgID, sourceID string) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.DeleteSource(ctx, userID, orgID, sourceID)
}

func (a *MigrationLifecycleAdapter) Scan(ctx context.Context, userID, orgID, sourceID string) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Scan(ctx, userID, orgID, sourceID)
}

func (a *MigrationLifecycleAdapter) ScanStream(ctx context.Context, userID, orgID, sourceID string, send func(event string, data any) error) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.ScanStream(ctx, userID, orgID, sourceID, send)
}

func (a *MigrationLifecycleAdapter) RevealEnv(ctx context.Context, userID, orgID, sourceID, containerID string) (map[string]string, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.RevealEnv(ctx, userID, orgID, sourceID, containerID)
}

func (a *MigrationLifecycleAdapter) Adopt(ctx context.Context, userID, orgID string, req modmigrations.AdoptMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Adopt(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) Reimport(ctx context.Context, userID, orgID string, req modmigrations.ReimportMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Reimport(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) RepoCompose(ctx context.Context, userID, orgID string, req modmigrations.RepoComposeMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.RepoCompose(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) Preview(ctx context.Context, userID, orgID string, req modmigrations.PreviewMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Preview(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) StartMigration(ctx context.Context, userID, orgID string, req modmigrations.StartMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.StartMigration(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) StartProjectMove(ctx context.Context, userID, orgID string, req modmigrations.StartProjectMoveRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.StartProjectMove(ctx, userID, orgID, req)
}

func (a *MigrationLifecycleAdapter) Active(ctx context.Context, userID, orgID string) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Active(ctx, userID, orgID)
}

func (a *MigrationLifecycleAdapter) ListRuns(ctx context.Context, userID, orgID string) (any, error) {
	if a.Service == nil {
		return []any{}, nil
	}
	return a.Service.ListRuns(ctx, userID, orgID)
}

func (a *MigrationLifecycleAdapter) GetRun(ctx context.Context, userID, runID string) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.GetRun(ctx, userID, runID)
}

func (a *MigrationLifecycleAdapter) StreamRun(ctx context.Context, userID, runID string, send func(event string, data any) error) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.StreamRun(ctx, userID, runID, send)
}

func (a *MigrationLifecycleAdapter) Cutover(ctx context.Context, userID, runID string, req modmigrations.CutoverMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Cutover(ctx, userID, runID, req)
}

func (a *MigrationLifecycleAdapter) Cancel(ctx context.Context, userID, runID string) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.Cancel(ctx, userID, runID)
}

func (a *MigrationLifecycleAdapter) Respond(ctx context.Context, userID, runID string, req modmigrations.RespondMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Respond(ctx, userID, runID, req)
}

func (a *MigrationLifecycleAdapter) Resume(ctx context.Context, userID, runID string, req modmigrations.ResumeMigrationRequest) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Resume(ctx, userID, runID, req)
}

func (a *MigrationLifecycleAdapter) CleanupTarget(ctx context.Context, userID, runID string) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.CleanupTarget(ctx, userID, runID)
}

func (a *MigrationLifecycleAdapter) DeleteRun(ctx context.Context, userID, runID string) error {
	if a.Service == nil {
		return nil
	}
	return a.Service.DeleteRun(ctx, userID, runID)
}

type MigrationBundleAdapter struct {
	Service modmigrations.BundleService
}

func (a *MigrationBundleAdapter) Export(ctx context.Context, passphrase string) ([]byte, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Export(ctx, passphrase)
}

func (a *MigrationBundleAdapter) Import(ctx context.Context, data []byte, passphrase string) (any, error) {
	if a.Service == nil {
		return nil, nil
	}
	return a.Service.Import(ctx, data, passphrase)
}

type MigrationUserCounterAdapter struct {
	Users users.UserRepository
}

func (a *MigrationUserCounterAdapter) CountUsers(ctx context.Context) (int, error) {
	if a.Users == nil {
		return 0, nil
	}
	return a.Users.CountUsers(ctx)
}
