package http

import (
	"context"

	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/modules/users"
)

type migrationLifecycleAdapter struct {
	service modmigrations.LifecycleService
}

func (a *migrationLifecycleAdapter) ListSources(ctx context.Context, userID, orgID string) (any, error) {
	if a.service == nil {
		return []any{}, nil
	}
	return a.service.ListSources(ctx, userID, orgID)
}

func (a *migrationLifecycleAdapter) TestSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (string, error) {
	if a.service == nil {
		return "", nil
	}
	return a.service.TestSource(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) CreateSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.CreateSource(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) DeleteSource(ctx context.Context, userID, orgID, sourceID string) error {
	if a.service == nil {
		return nil
	}
	return a.service.DeleteSource(ctx, userID, orgID, sourceID)
}

func (a *migrationLifecycleAdapter) Scan(ctx context.Context, userID, orgID, sourceID string) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Scan(ctx, userID, orgID, sourceID)
}

func (a *migrationLifecycleAdapter) ScanStream(ctx context.Context, userID, orgID, sourceID string, send func(event string, data any) error) error {
	if a.service == nil {
		return nil
	}
	return a.service.ScanStream(ctx, userID, orgID, sourceID, send)
}

func (a *migrationLifecycleAdapter) RevealEnv(ctx context.Context, userID, orgID, sourceID, containerID string) (map[string]string, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.RevealEnv(ctx, userID, orgID, sourceID, containerID)
}

func (a *migrationLifecycleAdapter) Adopt(ctx context.Context, userID, orgID string, req modmigrations.AdoptMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Adopt(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) Reimport(ctx context.Context, userID, orgID string, req modmigrations.ReimportMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Reimport(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) RepoCompose(ctx context.Context, userID, orgID string, req modmigrations.RepoComposeMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.RepoCompose(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) Preview(ctx context.Context, userID, orgID string, req modmigrations.PreviewMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Preview(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) StartMigration(ctx context.Context, userID, orgID string, req modmigrations.StartMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.StartMigration(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) StartProjectMove(ctx context.Context, userID, orgID string, req modmigrations.StartProjectMoveRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.StartProjectMove(ctx, userID, orgID, req)
}

func (a *migrationLifecycleAdapter) Active(ctx context.Context, userID, orgID string) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Active(ctx, userID, orgID)
}

func (a *migrationLifecycleAdapter) ListRuns(ctx context.Context, userID, orgID string) (any, error) {
	if a.service == nil {
		return []any{}, nil
	}
	return a.service.ListRuns(ctx, userID, orgID)
}

func (a *migrationLifecycleAdapter) GetRun(ctx context.Context, userID, runID string) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.GetRun(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) StreamRun(ctx context.Context, userID, runID string, send func(event string, data any) error) error {
	if a.service == nil {
		return nil
	}
	return a.service.StreamRun(ctx, userID, runID, send)
}

func (a *migrationLifecycleAdapter) Cutover(ctx context.Context, userID, runID string, req modmigrations.CutoverMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Cutover(ctx, userID, runID, req)
}

func (a *migrationLifecycleAdapter) Cancel(ctx context.Context, userID, runID string) error {
	if a.service == nil {
		return nil
	}
	return a.service.Cancel(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) Respond(ctx context.Context, userID, runID string, req modmigrations.RespondMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Respond(ctx, userID, runID, req)
}

func (a *migrationLifecycleAdapter) Resume(ctx context.Context, userID, runID string, req modmigrations.ResumeMigrationRequest) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Resume(ctx, userID, runID, req)
}

func (a *migrationLifecycleAdapter) CleanupTarget(ctx context.Context, userID, runID string) error {
	if a.service == nil {
		return nil
	}
	return a.service.CleanupTarget(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) DeleteRun(ctx context.Context, userID, runID string) error {
	if a.service == nil {
		return nil
	}
	return a.service.DeleteRun(ctx, userID, runID)
}

type migrationBundleAdapter struct {
	service modmigrations.BundleService
}

func (a *migrationBundleAdapter) Export(ctx context.Context, passphrase string) ([]byte, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Export(ctx, passphrase)
}

func (a *migrationBundleAdapter) Import(ctx context.Context, data []byte, passphrase string) (any, error) {
	if a.service == nil {
		return nil, nil
	}
	return a.service.Import(ctx, data, passphrase)
}

type migrationUserCounterAdapter struct {
	users users.UserRepository
}

func (a *migrationUserCounterAdapter) CountUsers(ctx context.Context) (int, error) {
	if a.users == nil {
		return 0, nil
	}
	return a.users.CountUsers(ctx)
}
