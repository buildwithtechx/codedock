package migrations

import (
	"context"
	"fmt"
	"testing"

	"codedock.run/codedock/internal/models"
)

func startExternal(t *testing.T, env *migrationTestEnv, mode models.MigrationMode, kill bool) *StartMigrationResult {
	t.Helper()
	result, err := env.service.StartMigration(context.Background(), env.memberID, env.orgID, models.StartMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID, fakeDbID},
		ProjectName: "shop", Mode: mode, ImportEnv: true, KillOriginals: kill,
	})
	if err != nil {
		t.Fatalf("start migration: %v", err)
	}
	if result.Confirmation == "" {
		t.Fatal("expected cutover confirmation token")
	}
	return result
}

func TestMigrationPipelineCutover(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	result := startExternal(t, env, models.MigrationModeMove, false)
	detail := waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusAwaitingCutover)
	if len(detail.Progress) != 2 {
		t.Fatalf("expected 2 service progresses, got %d", len(detail.Progress))
	}
	for _, entry := range detail.Progress {
		if !entry.Verified || entry.ServiceID == "" {
			t.Fatalf("expected verified service, got %#v", entry)
		}
	}
	if err := env.service.CutoverRun(ctx, env.memberID, result.RunID, models.CutoverMigrationRequest{Confirmation: result.Confirmation, Kill: true}); err == nil {
		t.Fatal("expected member cutover to fail")
	}
	if err := env.service.CutoverRun(ctx, env.adminID, result.RunID, models.CutoverMigrationRequest{Confirmation: "wrong", Kill: true}); err == nil {
		t.Fatal("expected wrong confirmation to fail")
	}
	if err := env.service.CutoverRun(ctx, env.adminID, result.RunID, models.CutoverMigrationRequest{Confirmation: result.Confirmation, Kill: true}); err != nil {
		t.Fatalf("cutover: %v", err)
	}
	final := waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusCompleted)
	if len(final.Run.Logs) == 0 {
		t.Fatal("expected migration logs")
	}
	if len(env.source.stopped) != 2 || len(env.source.removed) != 2 {
		t.Fatalf("expected originals stopped and removed, got %+v", env.source)
	}
	if env.target.loaded != 2 {
		t.Fatalf("expected 2 images loaded on target, got %d", env.target.loaded)
	}
}

func TestMigrationPipelineCopy(t *testing.T) {
	env := setupMigrationTest(t)
	result := startExternal(t, env, models.MigrationModeCopy, false)
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusCompleted)
	if len(env.source.stopped) != 0 || len(env.source.removed) != 0 {
		t.Fatal("expected copy to leave originals running")
	}
}

func TestMigrationVolumePromptRespond(t *testing.T) {
	env := setupMigrationTest(t)
	env.target.volumes["webdata"] = "index.html\n"
	result, err := env.service.StartMigration(context.Background(), env.memberID, env.orgID, models.StartMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID},
		ProjectName: "shop", Mode: models.MigrationModeMove,
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	detail := waitRunPrompt(t, env, env.memberID, result.RunID)
	if detail.Prompt.Kind != models.MigrationPromptVolumeConflict {
		t.Fatalf("expected volume prompt, got %#v", detail.Prompt)
	}
	if err := env.service.RespondRun(context.Background(), env.memberID, result.RunID, models.RespondMigrationRequest{PromptID: "stale", OptionID: "overwrite"}); err == nil {
		t.Fatal("expected stale prompt id to fail")
	}
	if err := env.service.RespondRun(context.Background(), env.memberID, result.RunID, models.RespondMigrationRequest{PromptID: detail.Prompt.ID, OptionID: "overwrite"}); err != nil {
		t.Fatalf("respond: %v", err)
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusAwaitingCutover)
}

func TestMigrationCancelRollsBack(t *testing.T) {
	env := setupMigrationTest(t)
	env.target.volumes["webdata"] = "index.html\n"
	result, err := env.service.StartMigration(context.Background(), env.memberID, env.orgID, models.StartMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID},
		ProjectName: "shop", Mode: models.MigrationModeMove,
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	waitRunPrompt(t, env, env.memberID, result.RunID)
	if err := env.service.CancelRun(context.Background(), env.memberID, result.RunID); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusCancelled)
}

func TestMigrationResumeAfterFailure(t *testing.T) {
	env := setupMigrationTest(t)
	env.deployer.failNext = fmt.Errorf("registry unavailable")
	result := startExternal(t, env, models.MigrationModeMove, false)
	failed := waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusFailed)
	if failed.Run.Error == "" {
		t.Fatal("expected failure message")
	}
	resumed, err := env.service.ResumeRun(context.Background(), env.memberID, result.RunID, models.ResumeMigrationRequest{})
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Confirmation == "" || resumed.Confirmation == result.Confirmation {
		t.Fatal("expected fresh cutover confirmation on resume")
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusAwaitingCutover)
}

func TestMigrationCleanupTarget(t *testing.T) {
	env := setupMigrationTest(t)
	env.deployer.failNext = fmt.Errorf("registry unavailable")
	result := startExternal(t, env, models.MigrationModeMove, false)
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusFailed)
	if len(env.target.volumes) == 0 {
		t.Fatal("expected copied target volumes")
	}
	if err := env.service.CleanupTargetData(context.Background(), env.memberID, result.RunID); err == nil {
		t.Fatal("expected member cleanup to fail")
	}
	if err := env.service.CleanupTargetData(context.Background(), env.adminID, result.RunID); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if len(env.target.volumes) != 0 {
		t.Fatalf("expected target volumes removed, got %v", env.target.volumes)
	}
	if _, err := env.service.ResumeRun(context.Background(), env.memberID, result.RunID, models.ResumeMigrationRequest{}); err != nil {
		t.Fatalf("resume after cleanup: %v", err)
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusAwaitingCutover)
	if err := env.service.CleanupTargetData(context.Background(), env.adminID, result.RunID); err == nil {
		t.Fatal("expected cleanup of live run to fail")
	}
}

func TestProjectMoveLifecycle(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	adopted, err := env.service.AdoptSource(ctx, env.memberID, env.orgID, models.AdoptMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID}, ProjectName: "shop",
	})
	if err != nil {
		t.Fatalf("adopt: %v", err)
	}
	serviceID := adopted.Services[0].ServiceID
	app, err := env.service.appRepo.GetByID(ctx, serviceID)
	if err != nil {
		t.Fatalf("load service: %v", err)
	}
	app.ContainerID = fakeWebID
	if err := env.service.appRepo.Update(ctx, app); err != nil {
		t.Fatalf("stage source container: %v", err)
	}
	moved, err := env.service.StartProjectMove(ctx, env.memberID, env.orgID, models.StartProjectMoveRequest{
		ProjectID: adopted.ProjectID, ServiceIDs: []string{serviceID}, Mode: models.MigrationModeMove,
	})
	if err != nil {
		t.Fatalf("start move: %v", err)
	}
	waitRunStatus(t, env, env.memberID, moved.RunID, models.MigrationStatusAwaitingCutover)
	if err := env.service.CutoverRun(ctx, env.adminID, moved.RunID, models.CutoverMigrationRequest{Confirmation: moved.Confirmation}); err != nil {
		t.Fatalf("cutover retain: %v", err)
	}
	waitRunStatus(t, env, env.memberID, moved.RunID, models.MigrationStatusCompleted)
	if len(env.target.stopped) != 1 || len(env.target.removed) != 0 {
		t.Fatalf("expected original stopped and retained, got %+v", env.target)
	}
}

func TestMigrationDeleteGuards(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	result := startExternal(t, env, models.MigrationModeMove, false)
	if err := env.service.DeleteRunRecord(ctx, env.adminID, result.RunID); err == nil {
		t.Fatal("expected active run delete to fail")
	}
	if err := env.service.DeleteSource(ctx, env.adminID, env.orgID, env.sourceID); err == nil {
		t.Fatal("expected source delete with active run to fail")
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusAwaitingCutover)
	if err := env.service.CutoverRun(ctx, env.adminID, result.RunID, models.CutoverMigrationRequest{Confirmation: result.Confirmation, Kill: true}); err != nil {
		t.Fatalf("cutover: %v", err)
	}
	waitRunStatus(t, env, env.memberID, result.RunID, models.MigrationStatusCompleted)
	if err := env.service.DeleteRunRecord(ctx, env.adminID, result.RunID); err != nil {
		t.Fatalf("delete terminal run: %v", err)
	}
	if err := env.service.DeleteSource(ctx, env.adminID, env.orgID, env.sourceID); err != nil {
		t.Fatalf("delete source: %v", err)
	}
}
