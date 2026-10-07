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

type sqliteMigrationRepository struct {
	db    *sql.DB
	vault Vault
}

func NewMigrationRepository(db *sql.DB, vault Vault) MigrationRepository {
	return &sqliteMigrationRepository{db: db, vault: vault}
}

const migrationSourceColumns = `id, organization_id, name, ssh_host, ssh_port, ssh_user, ssh_auth_method, ssh_key, ssh_password, fingerprint, created_at, updated_at`

func (r *sqliteMigrationRepository) CreateSource(ctx context.Context, source *models.MigrationSource) error {
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, source.ID, source.OrganizationID, source.Name, source.SSHHost, source.SSHPort, source.SSHUser, source.SSHAuthMethod, key, password, source.Fingerprint)
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) GetSource(ctx context.Context, id string) (*models.MigrationSource, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+migrationSourceColumns+` FROM migration_sources WHERE id = ?`, id)
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

func (r *sqliteMigrationRepository) ListSourcesByOrg(ctx context.Context, orgID string) ([]*models.MigrationSource, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+migrationSourceColumns+` FROM migration_sources WHERE organization_id = ? ORDER BY created_at DESC`, orgID)
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

func (r *sqliteMigrationRepository) DeleteSource(ctx context.Context, id string) error {
	var active int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM migration_runs WHERE source_id = ? AND status NOT IN ('completed','failed','cancelled')`, id).Scan(&active)
	if err != nil {
		return fmt.Errorf("failed to check source runs: %w", err)
	}
	if active > 0 {
		return fmt.Errorf("source still has %d active run(s)", active)
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM migration_sources WHERE id = ?`, id); err != nil {
		return fmt.Errorf("failed to delete migration source: %w", err)
	}
	return nil
}

const migrationRunColumns = `id, organization_id, user_id, source_id, source_kind, project_id, target_server_id, mode, status, phase, selection, progress, logs, prompt, token_hash, error, cancel_requested, created_at, updated_at`

func (r *sqliteMigrationRepository) CreateRun(ctx context.Context, run *models.MigrationRun) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO migration_runs (id, organization_id, user_id, source_id, source_kind, project_id, target_server_id, mode, status, phase, selection, progress, logs, prompt, token_hash, error, cancel_requested, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, run.ID, run.OrganizationID, run.UserID, run.SourceID, run.SourceKind, run.ProjectID, run.TargetServerID, string(run.Mode), string(run.Status), run.Phase, run.Selection, run.Progress, run.Logs, run.Prompt, run.TokenHash, run.Error, run.CancelRequested)
	if err != nil {
		return fmt.Errorf("failed to create migration run: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) GetRun(ctx context.Context, id string) (*models.MigrationRun, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE id = ?`, id)
	return r.scanRun(row)
}

func (r *sqliteMigrationRepository) UpdateRun(ctx context.Context, run *models.MigrationRun) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE migration_runs SET project_id = ?, target_server_id = ?, mode = ?, status = ?, phase = ?, selection = ?, progress = ?, logs = ?, prompt = ?, token_hash = ?, error = ?, cancel_requested = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, run.ProjectID, run.TargetServerID, string(run.Mode), string(run.Status), run.Phase, run.Selection, run.Progress, run.Logs, run.Prompt, run.TokenHash, run.Error, run.CancelRequested, run.ID)
	if err != nil {
		return fmt.Errorf("failed to update migration run: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) UpdateRunStatus(ctx context.Context, id string, status models.MigrationStatus, phase, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET status = ?, phase = ?, error = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, string(status), phase, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update migration run status: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) AppendRunLogs(ctx context.Context, id, logs string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET logs = substr(logs || ?, -65536), updated_at = CURRENT_TIMESTAMP WHERE id = ?`, logs, id)
	if err != nil {
		return fmt.Errorf("failed to append migration logs: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) SetRunPrompt(ctx context.Context, id, prompt string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET prompt = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, prompt, id)
	if err != nil {
		return fmt.Errorf("failed to set migration prompt: %w", err)
	}
	return nil
}

func (r *sqliteMigrationRepository) RequestCancel(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE migration_runs SET cancel_requested = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND status NOT IN ('completed','failed','cancelled','awaiting_cutover')`, id)
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

func (r *sqliteMigrationRepository) ListRunsByOrg(ctx context.Context, orgID string, limit int) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE organization_id = ? ORDER BY created_at DESC LIMIT ?`, orgID, limit)
}

func (r *sqliteMigrationRepository) ListRunsBySource(ctx context.Context, sourceID string, limit int) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE source_id = ? ORDER BY created_at DESC LIMIT ?`, sourceID, limit)
}

func (r *sqliteMigrationRepository) ActiveRunForSource(ctx context.Context, sourceID string) (*models.MigrationRun, error) {
	rows, err := r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE source_id = ? AND status NOT IN ('completed','failed','cancelled') ORDER BY created_at DESC LIMIT 1`, sourceID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func (r *sqliteMigrationRepository) ListActive(ctx context.Context) ([]*models.MigrationRun, error) {
	return r.listRuns(ctx, `SELECT `+migrationRunColumns+` FROM migration_runs WHERE status NOT IN ('completed','failed','cancelled','awaiting_cutover') ORDER BY created_at DESC`)
}

func (r *sqliteMigrationRepository) DeleteRun(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM migration_runs WHERE id = ? AND status IN ('completed','failed','cancelled')`, id)
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

func (r *sqliteMigrationRepository) listRuns(ctx context.Context, query string, args ...any) ([]*models.MigrationRun, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list migration runs: %w", err)
	}
	defer rows.Close()
	var runs []*models.MigrationRun
	for rows.Next() {
		var run models.MigrationRun
		var mode, status string
		if err := rows.Scan(&run.ID, &run.OrganizationID, &run.UserID, &run.SourceID, &run.SourceKind, &run.ProjectID, &run.TargetServerID, &mode, &status, &run.Phase, &run.Selection, &run.Progress, &run.Logs, &run.Prompt, &run.TokenHash, &run.Error, &run.CancelRequested, &run.CreatedAt, &run.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration run: %w", err)
		}
		run.Mode = models.MigrationMode(mode)
		run.Status = models.MigrationStatus(status)
		runs = append(runs, &run)
	}
	return runs, rows.Err()
}

func (r *sqliteMigrationRepository) scanRun(row *sql.Row) (*models.MigrationRun, error) {
	var run models.MigrationRun
	var mode, status string
	if err := row.Scan(&run.ID, &run.OrganizationID, &run.UserID, &run.SourceID, &run.SourceKind, &run.ProjectID, &run.TargetServerID, &mode, &status, &run.Phase, &run.Selection, &run.Progress, &run.Logs, &run.Prompt, &run.TokenHash, &run.Error, &run.CancelRequested, &run.CreatedAt, &run.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan migration run: %w", err)
	}
	run.Mode = models.MigrationMode(mode)
	run.Status = models.MigrationStatus(status)
	return &run, nil
}

func (r *sqliteMigrationRepository) encryptSecret(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *sqliteMigrationRepository) decryptSource(source *models.MigrationSource) {
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
