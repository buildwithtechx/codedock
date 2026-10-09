package migrations

import (
	"context"
	"fmt"
	"sync"
	"time"

	"codedock/internal/models"
)

type promptAnswer struct {
	promptID string
	optionID string
	value    string
}

type runHandle struct {
	cancel context.CancelFunc
	answer chan promptAnswer
}

type runRegistry struct {
	mu      sync.Mutex
	handles map[string]*runHandle
}

func newRunRegistry() *runRegistry {
	return &runRegistry{handles: map[string]*runHandle{}}
}

func (r *runRegistry) track(id string) (context.Context, *runHandle) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	handle := &runHandle{cancel: cancel, answer: make(chan promptAnswer, 1)}
	r.handles[id] = handle
	return ctx, handle
}

func (r *runRegistry) release(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handles, id)
}

func (r *runRegistry) get(id string) *runHandle {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.handles[id]
}

func (r *runRegistry) cancel(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if handle, ok := r.handles[id]; ok {
		handle.cancel()
	}
}

type RunDetail struct {
	Run       *models.MigrationRun              `json:"run"`
	Selection models.MigrationSelection       `json:"selection"`
	Progress  []models.MigrationServiceProgress `json:"progress"`
	Prompt    *models.MigrationPrompt         `json:"prompt"`
}

func (s *Service) GetRun(ctx context.Context, userID, runID string) (*RunDetail, error) {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return nil, err
	}
	return &RunDetail{
		Run:       run,
		Selection: decodeSelection(run.Selection),
		Progress:  decodeProgress(run.Progress),
		Prompt:    decodePrompt(run.Prompt),
	}, nil
}

func (s *Service) ActiveRun(ctx context.Context, userID, orgID, sourceID string) (*RunDetail, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	run, err := s.migrations.ActiveRunForSource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if run == nil || run.OrganizationID != orgID {
		return nil, fmt.Errorf("no active migration")
	}
	return &RunDetail{
		Run:       run,
		Selection: decodeSelection(run.Selection),
		Progress:  decodeProgress(run.Progress),
		Prompt:    decodePrompt(run.Prompt),
	}, nil
}

func (s *Service) ListRuns(ctx context.Context, userID, orgID, sourceID string) ([]*models.MigrationRun, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if sourceID != "" {
		if _, err := s.orgSource(ctx, orgID, sourceID); err != nil {
			return nil, err
		}
		return s.migrations.ListRunsBySource(ctx, sourceID, 50)
	}
	return s.migrations.ListRunsByOrg(ctx, orgID, 50)
}

func (s *Service) DeleteRunRecord(ctx context.Context, userID, runID string) error {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return err
	}
	if err := s.requireOrg(ctx, userID, run.OrganizationID, true); err != nil {
		return err
	}
	return s.migrations.DeleteRun(ctx, runID)
}

func (s *Service) CancelRun(ctx context.Context, userID, runID string) error {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return err
	}
	if err := s.migrations.RequestCancel(ctx, run.ID); err != nil {
		return err
	}
	s.runs.cancel(run.ID)
	return nil
}

func (s *Service) RespondRun(ctx context.Context, userID, runID string, req models.RespondMigrationRequest) error {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return err
	}
	prompt := decodePrompt(run.Prompt)
	if prompt == nil {
		return fmt.Errorf("migration is not awaiting input")
	}
	if prompt.ID != req.PromptID {
		return fmt.Errorf("prompt expired; review the current prompt")
	}
	if time.Now().Unix() > prompt.ExpiresAt {
		return fmt.Errorf("prompt expired; resume the migration to continue")
	}
	valid := false
	for _, option := range prompt.Options {
		if option.ID == req.OptionID {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("unknown prompt option")
	}
	if req.OptionID == "rename" && req.Value == "" {
		return fmt.Errorf("rename requires a replacement name")
	}
	handle := s.runs.get(run.ID)
	if handle == nil {
		return fmt.Errorf("migration worker is not running; resume the migration to continue")
	}
	select {
	case handle.answer <- promptAnswer{promptID: req.PromptID, optionID: req.OptionID, value: req.Value}:
		return nil
	default:
		return fmt.Errorf("migration already has a pending answer")
	}
}

func (s *Service) Recover(ctx context.Context) error {
	rows, err := s.migrations.ListActive(ctx)
	if err != nil {
		return err
	}
	for _, run := range rows {
		if err := s.migrations.UpdateRunStatus(ctx, run.ID, models.MigrationStatusFailed, run.Phase, "interrupted by restart; resume to continue"); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) appendLog(ctx context.Context, runID, line string) {
	_ = s.migrations.AppendRunLogs(ctx, runID, line+"\n")
}

func (s *Service) setStatus(ctx context.Context, runID string, status models.MigrationStatus, phase, errMsg string) {
	_ = s.migrations.UpdateRunStatus(ctx, runID, status, phase, errMsg)
}
