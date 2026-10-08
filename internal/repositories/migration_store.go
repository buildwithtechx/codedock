package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"codedock.run/codedock/internal/models"
)

type MigrationRepository interface {
	CreateSource(ctx context.Context, source *models.MigrationSource) error
	GetSource(ctx context.Context, id string) (*models.MigrationSource, error)
	ListSourcesByOrg(ctx context.Context, orgID string) ([]*models.MigrationSource, error)
	DeleteSource(ctx context.Context, id string) error
	CreateRun(ctx context.Context, run *models.MigrationRun) error
	GetRun(ctx context.Context, id string) (*models.MigrationRun, error)
	UpdateRun(ctx context.Context, run *models.MigrationRun) error
	UpdateRunStatus(ctx context.Context, id string, status models.MigrationStatus, phase, errMsg string) error
	AppendRunLogs(ctx context.Context, id, logs string) error
	SetRunPrompt(ctx context.Context, id, prompt string) error
	RequestCancel(ctx context.Context, id string) error
	ListRunsByOrg(ctx context.Context, orgID string, limit int) ([]*models.MigrationRun, error)
	ListRunsBySource(ctx context.Context, sourceID string, limit int) ([]*models.MigrationRun, error)
	ActiveRunForSource(ctx context.Context, sourceID string) (*models.MigrationRun, error)
	ListActive(ctx context.Context) ([]*models.MigrationRun, error)
	DeleteRun(ctx context.Context, id string) error
}

type postgresMigrationRepository struct {
	db    *sql.DB
	vault Vault
}

func NewMigrationRepository(db *sql.DB, vault Vault) MigrationRepository {
	return &postgresMigrationRepository{db: db, vault: vault}
}

const migrationSourceColumns = `id, organization_id, name, ssh_host, ssh_port, ssh_user, ssh_auth_method, ssh_key, ssh_password, fingerprint, created_at, updated_at`

func (r *postgresMigrationRepository) CreateSource(ctx context.Context, source *models.MigrationSource) error {
	key, err := r.encryptSecret(source.SSHKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt source key: %w", err)
	}
	password, err := r.encryptSecret(source.SSHPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt source password: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO migration_sources (id, organization_id, name, ssh_host, ssh_port, ssh_user, ssh_auth_method, ssh_key, ssh_password, fingerprint, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, source.ID, source.OrganizationID, source.Name, source.SSHHost, source.SSHPort, source.SSHUser, source.SSHAuthMethod, key, password, source.Fingerprint)
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) GetSource(ctx context.Context, id string) (*models.MigrationSource, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+migrationSourceColumns+` FROM migration_sources WHERE id = $1`, id)
	var source models.MigrationSource
	if err := row.Scan(&source.ID, &source.OrganizationID, &source.Name, &source.SSHHost, &source.SSHPort, &source.SSHUser, &source.SSHAuthMethod, &source.SSHKey, &source.SSHPassword, &source.Fingerprint, &source.CreatedAt, &source.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan migration source: %w", err)
	}
	r.decryptSource(&source)
	return &source, nil
}

func (r *postgresMigrationRepository) ListSourcesByOrg(ctx context.Context, orgID string) ([]*models.MigrationSource, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+migrationSourceColumns+` FROM migration_sources WHERE organization_id = $1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list migration sources: %w", err)
	}
	defer rows.Close()
	var sources []*models.MigrationSource
	for rows.Next() {
		var source models.MigrationSource
		if err := rows.Scan(&source.ID, &source.OrganizationID, &source.Name, &source.SSHHost, &source.SSHPort, &source.SSHUser, &source.SSHAuthMethod, &source.SSHKey, &source.SSHPassword, &source.Fingerprint, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration source: %w", err)
		}
		source.SSHKey = ""
		source.SSHPassword = ""
		sources = append(sources, &source)
	}
	return sources, rows.Err()
}

func (r *postgresMigrationRepository) DeleteSource(ctx context.Context, id string) error {
	var active int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM migration_runs WHERE source_id = $1 AND status NOT IN ('completed','failed','cancelled')`, id).Scan(&active)
	if err != nil {
		return fmt.Errorf("failed to check source runs: %w", err)
	}
	if active > 0 {
		return fmt.Errorf("source still has %d active run(s)", active)
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM migration_sources WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete migration source: %w", err)
	}
	return nil
}

const migrationRunColumns = `id, organization_id, user_id, source_id, source_kind, project_id, target_server_id, mode, status, phase, selection, progress, logs, prompt, token_hash, error, cancel_requested, created_at, updated_at`

func (r *postgresMigrationRepository) CreateRun(ctx context.Context, run *models.MigrationRun) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO migration_runs (id, organization_id, user_id, source_id, source_kind, project_id, target_server_id, mode, status, phase, selection, progress, logs, prompt, token_hash, error, cancel_requested, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, run.ID, run.OrganizationID, run.UserID, run.SourceID, run.SourceKind, run.ProjectID, run.TargetServerID, string(run.Mode), string(run.Status), run.Phase, run.Selection, run.Progress, run.Logs, run.Prompt, run.TokenHash, run.Error, boolToInt(run.CancelRequested))
	if err != nil {
		return fmt.Errorf("failed to create migration run: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) GetRun(ctx context.Context, id string) (*models.MigrationRun, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE id = $1`, id)
	return r.scanRun(row)
}

func (r *postgresMigrationRepository) UpdateRun(ctx context.Context, run *models.MigrationRun) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE migration_runs SET project_id = $1, target_server_id = $2, mode = $3, status = $4, phase = $5, selection = $6, progress = $7, logs = $8, prompt = $9, token_hash = $10, error = $11, cancel_requested = $12, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
	`, run.ProjectID, run.TargetServerID, string(run.Mode), string(run.Status), run.Phase, run.Selection, run.Progress, run.Logs, run.Prompt, run.TokenHash, run.Error, boolToInt(run.CancelRequested), run.ID)
	if err != nil {
		return fmt.Errorf("failed to update migration run: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) UpdateRunStatus(ctx context.Context, id string, status models.MigrationStatus, phase, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET status = $1, phase = $2, error = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`, string(status), phase, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update migration run status: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) AppendRunLogs(ctx context.Context, id, logs string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET logs = right(logs || $1, 65536), updated_at = CURRENT_TIMESTAMP WHERE id = $2`, logs, id)
	if err != nil {
		return fmt.Errorf("failed to append migration logs: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) SetRunPrompt(ctx context.Context, id, prompt string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET prompt = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, prompt, id)
	if err != nil {
		return fmt.Errorf("failed to set migration prompt: %w", err)
	}
	return nil
}

func (r *postgresMigrationRepository) RequestCancel(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET cancel_requested = 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status NOT IN ('completed','failed','cancelled','awaiting_cutover')`, id)
	if err != nil {
		return fmt.Errorf("failed to request migration cancel: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("migration cannot be cancelled in its current state")
	}
	return nil
}

func (r *postgresMigrationRepository) ListRunsByOrg(ctx context.Context, orgID string, limit int) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2`, orgID, limit)
}

func (r *postgresMigrationRepository) ListRunsBySource(ctx context.Context, sourceID string, limit int) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE source_id = $1 ORDER BY created_at DESC LIMIT $2`, sourceID, limit)
}

func (r *postgresMigrationRepository) ActiveRunForSource(ctx context.Context, sourceID string) (*models.MigrationRun, error) {
	rows, err := r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE source_id = $1 AND status NOT IN ('completed','failed','cancelled') ORDER BY created_at DESC LIMIT 1`, sourceID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func (r *postgresMigrationRepository) ListActive(ctx context.Context) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE status NOT IN ('completed','failed','cancelled','awaiting_cutover') ORDER BY created_at DESC`)
}

func (r *postgresMigrationRepository) DeleteRun(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM migration_runs WHERE id = $1 AND status IN ('completed','failed','cancelled')`, id)
	if err != nil {
		return fmt.Errorf("failed to delete migration run: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("only terminal migration runs can be deleted")
	}
	return nil
}

func (r *postgresMigrationRepository) listRuns(ctx context.Context, query string, args ...any) ([]*models.MigrationRun, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list migration runs: %w", err)
	}
	defer rows.Close()
	var runs []*models.MigrationRun
	for rows.Next() {
		var run models.MigrationRun
		var mode, status string
		var cancelRequested int
		if err := rows.Scan(&run.ID, &run.OrganizationID, &run.UserID, &run.SourceID, &run.SourceKind, &run.ProjectID, &run.TargetServerID, &mode, &status, &run.Phase, &run.Selection, &run.Progress, &run.Logs, &run.Prompt, &run.TokenHash, &run.Error, &cancelRequested, &run.CreatedAt, &run.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration run: %w", err)
		}
		run.Mode = models.MigrationMode(mode)
		run.Status = models.MigrationStatus(status)
		run.CancelRequested = cancelRequested != 0
		runs = append(runs, &run)
	}
	return runs, rows.Err()
}

func (r *postgresMigrationRepository) scanRun(row *sql.Row) (*models.MigrationRun, error) {
	var run models.MigrationRun
	var mode, status string
	var cancelRequested int
	if err := row.Scan(&run.ID, &run.OrganizationID, &run.UserID, &run.SourceID, &run.SourceKind, &run.ProjectID, &run.TargetServerID, &mode, &status, &run.Phase, &run.Selection, &run.Progress, &run.Logs, &run.Prompt, &run.TokenHash, &run.Error, &cancelRequested, &run.CreatedAt, &run.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan migration run: %w", err)
	}
	run.Mode = models.MigrationMode(mode)
	run.Status = models.MigrationStatus(status)
	run.CancelRequested = cancelRequested != 0
	return &run, nil
}

func (r *postgresMigrationRepository) encryptSecret(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *postgresMigrationRepository) decryptSource(source *models.MigrationSource) {
	if r.vault == nil {
		return
	}
	if source.SSHKey != "" {
		if plain, err := r.vault.Decrypt(source.SSHKey); err == nil {
			source.SSHKey = plain
		}
	}
	if source.SSHPassword != "" {
		if plain, err := r.vault.Decrypt(source.SSHPassword); err == nil {
			source.SSHPassword = plain
		}
	}
}
