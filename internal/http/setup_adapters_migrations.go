package http

import (
	"context"

	"codedock/internal/models"
	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/services/migrations"
	systemservices "codedock/internal/services/system"
)

type migrationLifecycleAdapter struct {
	service *migrations.Service
}

func (a *migrationLifecycleAdapter) ListSources(ctx context.Context, userID, orgID string) (any, error) {
	return a.service.ListSources(ctx, userID, orgID)
}

func (a *migrationLifecycleAdapter) TestSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (string, error) {
	return a.service.TestSource(ctx, userID, orgID, models.CreateMigrationSourceRequest{
		Name:          req.Name,
		SSHHost:       req.SSHHost,
		SSHPort:       req.SSHPort,
		SSHUser:       req.SSHUser,
		SSHAuthMethod: req.SSHAuthMethod,
		SSHKey:        req.SSHKey,
		SSHPassword:   req.SSHPassword,
		Fingerprint:   req.Fingerprint,
	})
}

func (a *migrationLifecycleAdapter) CreateSource(ctx context.Context, userID, orgID string, req modmigrations.CreateMigrationSourceRequest) (any, error) {
	return a.service.CreateSource(ctx, userID, orgID, models.CreateMigrationSourceRequest{
		Name:          req.Name,
		SSHHost:       req.SSHHost,
		SSHPort:       req.SSHPort,
		SSHUser:       req.SSHUser,
		SSHAuthMethod: req.SSHAuthMethod,
		SSHKey:        req.SSHKey,
		SSHPassword:   req.SSHPassword,
		Fingerprint:   req.Fingerprint,
	})
}

func (a *migrationLifecycleAdapter) DeleteSource(ctx context.Context, userID, orgID, sourceID string) error {
	return a.service.DeleteSource(ctx, userID, orgID, sourceID)
}

func (a *migrationLifecycleAdapter) Scan(ctx context.Context, userID, orgID, sourceID string) (any, error) {
	return a.service.ScanSource(ctx, userID, orgID, sourceID, nil)
}

func (a *migrationLifecycleAdapter) ScanStream(ctx context.Context, userID, orgID, sourceID string, send func(event string, data any) error) error {
	stack, err := a.service.ScanSource(ctx, userID, orgID, sourceID, func(step, detail string) {
		_ = send(step, detail)
	})
	if err != nil {
		return err
	}
	return send("RESULT", stack)
}

func (a *migrationLifecycleAdapter) RevealEnv(ctx context.Context, userID, orgID, sourceID, containerID string) (map[string]string, error) {
	return a.service.RevealSourceEnv(ctx, userID, orgID, models.RevealMigrationEnvRequest{
		SourceID:    sourceID,
		ContainerID: containerID,
	})
}

func (a *migrationLifecycleAdapter) Adopt(ctx context.Context, userID, orgID string, req modmigrations.AdoptMigrationRequest) (any, error) {
	return a.service.AdoptSource(ctx, userID, orgID, models.AdoptMigrationRequest{
		SourceID:     req.SourceID,
		ContainerIDs: req.ContainerIDs,
		ProjectName:  req.ProjectName,
		ImportEnv:    req.ImportEnv,
	})
}

func (a *migrationLifecycleAdapter) Reimport(ctx context.Context, userID, orgID string, req modmigrations.ReimportMigrationRequest) (any, error) {
	return a.service.ReimportSource(ctx, userID, orgID, models.ReimportMigrationRequest{
		SourceID:     req.SourceID,
		ContainerIDs: req.ContainerIDs,
	})
}

func (a *migrationLifecycleAdapter) RepoCompose(ctx context.Context, userID, orgID string, req modmigrations.RepoComposeMigrationRequest) (any, error) {
	return a.service.RepoCompose(ctx, userID, orgID, models.RepoComposeMigrationRequest{
		RepoURL:     req.RepoURL,
		Branch:      req.Branch,
		ComposePath: req.ComposePath,
	})
}

func (a *migrationLifecycleAdapter) Preview(ctx context.Context, userID, orgID string, req modmigrations.PreviewMigrationRequest) (any, error) {
	return a.service.PreviewMigration(ctx, userID, orgID, models.PreviewMigrationRequest{
		SourceID:       req.SourceID,
		ProjectID:      req.ProjectID,
		ContainerIDs:   req.ContainerIDs,
		TargetServerID: req.TargetServerID,
		Mode:           models.MigrationMode(req.Mode),
	})
}

func (a *migrationLifecycleAdapter) StartMigration(ctx context.Context, userID, orgID string, req modmigrations.StartMigrationRequest) (any, error) {
	return a.service.StartMigration(ctx, userID, orgID, models.StartMigrationRequest{
		SourceID:       req.SourceID,
		ContainerIDs:   req.ContainerIDs,
		ProjectName:    req.ProjectName,
		TargetServerID: req.TargetServerID,
		Mode:           models.MigrationMode(req.Mode),
		ImportEnv:      req.ImportEnv,
		KillOriginals:  req.KillOriginals,
	})
}

func (a *migrationLifecycleAdapter) StartProjectMove(ctx context.Context, userID, orgID string, req modmigrations.StartProjectMoveRequest) (any, error) {
	return a.service.StartProjectMove(ctx, userID, orgID, models.StartProjectMoveRequest{
		ProjectID:      req.ProjectID,
		ServiceIDs:     req.ServiceIDs,
		TargetServerID: req.TargetServerID,
		Mode:           models.MigrationMode(req.Mode),
		KillOriginals:  req.KillOriginals,
	})
}

func (a *migrationLifecycleAdapter) Active(ctx context.Context, userID, orgID string) (any, error) {
	return a.service.ActiveRun(ctx, userID, orgID, "")
}

func (a *migrationLifecycleAdapter) ListRuns(ctx context.Context, userID, orgID string) (any, error) {
	return a.service.ListRuns(ctx, userID, orgID, "")
}

func (a *migrationLifecycleAdapter) GetRun(ctx context.Context, userID, runID string) (any, error) {
	return a.service.GetRun(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) StreamRun(ctx context.Context, userID, runID string, send func(event string, data any) error) error {
	return nil
}

func (a *migrationLifecycleAdapter) Cutover(ctx context.Context, userID, runID string, req modmigrations.CutoverMigrationRequest) (any, error) {
	err := a.service.CutoverRun(ctx, userID, runID, models.CutoverMigrationRequest{
		Confirmation: req.Confirmation,
		Kill:         req.Kill,
	})
	return nil, err
}

func (a *migrationLifecycleAdapter) Cancel(ctx context.Context, userID, runID string) error {
	return a.service.CancelRun(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) Respond(ctx context.Context, userID, runID string, req modmigrations.RespondMigrationRequest) (any, error) {
	err := a.service.RespondRun(ctx, userID, runID, models.RespondMigrationRequest{
		PromptID: req.PromptID,
		OptionID: req.OptionID,
		Value:    req.Value,
	})
	return nil, err
}

func (a *migrationLifecycleAdapter) Resume(ctx context.Context, userID, runID string, req modmigrations.ResumeMigrationRequest) (any, error) {
	return a.service.ResumeRun(ctx, userID, runID, models.ResumeMigrationRequest{
		Overrides: req.Overrides,
		Skips:     req.Skips,
	})
}

func (a *migrationLifecycleAdapter) CleanupTarget(ctx context.Context, userID, runID string) error {
	return a.service.CleanupTargetData(ctx, userID, runID)
}

func (a *migrationLifecycleAdapter) DeleteRun(ctx context.Context, userID, runID string) error {
	return a.service.DeleteRunRecord(ctx, userID, runID)
}

type migrationBundleAdapter struct {
	service *systemservices.MigrationService
}

func (a *migrationBundleAdapter) Export(ctx context.Context, passphrase string) ([]byte, error) {
	return a.service.Export(ctx, passphrase)
}

func (a *migrationBundleAdapter) Import(ctx context.Context, data []byte, passphrase string) (any, error) {
	return a.service.Import(ctx, data, passphrase)
}

type migrationUserCounterAdapter struct {
	users interface {
		CountUsers(context.Context) (int, error)
	}
}

func (a *migrationUserCounterAdapter) CountUsers(ctx context.Context) (int, error) {
	return a.users.CountUsers(ctx)
}
