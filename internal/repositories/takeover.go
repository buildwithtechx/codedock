package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
)

type TakeoverRepository interface {
	Create(ctx context.Context, run *models.TakeoverRun) error
	GetByID(ctx context.Context, id string) (*models.TakeoverRun, error)
	UpdateStatus(ctx context.Context, id string, status models.TakeoverStatus, errMsg string) error
	UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error
	UpdateAdopted(ctx context.Context, id string, projectIDs []string) error
	ListByUser(ctx context.Context, userID string) ([]*models.TakeoverRun, error)
}

type sqliteTakeoverRepository struct {
	db    *sql.DB
	vault Vault
}

func NewTakeoverRepository(db *sql.DB, vaults ...Vault) TakeoverRepository {
	var vault Vault
	if len(vaults) > 0 {
		vault = vaults[0]
	}
	return &sqliteTakeoverRepository{db: db, vault: vault}
}

func (r *sqliteTakeoverRepository) encryptDiscovered(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *sqliteTakeoverRepository) decryptDiscovered(value string) string {
	if value == "" || r.vault == nil {
		return value
	}
	if decrypted, err := r.vault.Decrypt(value); err == nil {
		return decrypted
	}
	return value
}

func (r *sqliteTakeoverRepository) Create(ctx context.Context, run *models.TakeoverRun) error {
	if run.ID == "" {
		run.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	run.CreatedAt = now
	run.UpdatedAt = now
	discoveredJSON, err := r.encryptDiscovered(run.DiscoveredJSON)
	if err != nil {
		return fmt.Errorf("encrypt discovered stack: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO takeover_runs (id, user_id, source_host, source_platform, status, discovered_json, adopted_project_ids, error, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		run.ID, run.UserID, run.SourceHost, run.SourcePlatform, run.Status,
		discoveredJSON, run.AdoptedProjectIDs, run.Error, run.CreatedAt, run.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create takeover run: %w", err)
	}
	return nil
}

func (r *sqliteTakeoverRepository) GetByID(ctx context.Context, id string) (*models.TakeoverRun, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, source_host, source_platform, status, discovered_json, adopted_project_ids, error, created_at, updated_at
		FROM takeover_runs WHERE id = ?`, id)
	return r.scan(row)
}

func (r *sqliteTakeoverRepository) UpdateStatus(ctx context.Context, id string, status models.TakeoverStatus, errMsg string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET status = ?, error = ?, updated_at = ? WHERE id = ?`,
		status, errMsg, time.Now().UTC(), id,
	)
	return err
}

func (r *sqliteTakeoverRepository) UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error {
	encrypted, err := r.encryptDiscovered(discoveredJSON)
	if err != nil {
		return fmt.Errorf("encrypt discovered stack: %w", err)
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET discovered_json = ?, status = ?, updated_at = ? WHERE id = ?`,
		encrypted, models.TakeoverStatusScanned, time.Now().UTC(), id,
	)
	return err
}

func (r *sqliteTakeoverRepository) UpdateAdopted(ctx context.Context, id string, projectIDs []string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET adopted_project_ids = ?, status = ?, updated_at = ? WHERE id = ?`,
		strings.Join(projectIDs, ","), models.TakeoverStatusDone, time.Now().UTC(), id,
	)
	return err
}

func (r *sqliteTakeoverRepository) ListByUser(ctx context.Context, userID string) ([]*models.TakeoverRun, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, source_host, source_platform, status, discovered_json, adopted_project_ids, error, created_at, updated_at
		FROM takeover_runs WHERE user_id = ? ORDER BY created_at DESC LIMIT 50`, userID)
	if err != nil {
		return nil, fmt.Errorf("list takeover runs: %w", err)
	}
	defer rows.Close()
	var runs []*models.TakeoverRun
	for rows.Next() {
		run, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("takeover runs iteration: %w", err)
	}
	return runs, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func (r *sqliteTakeoverRepository) scan(s scannable) (*models.TakeoverRun, error) {
	var run models.TakeoverRun
	err := s.Scan(
		&run.ID, &run.UserID, &run.SourceHost, &run.SourcePlatform,
		&run.Status, &run.DiscoveredJSON, &run.AdoptedProjectIDs,
		&run.Error, &run.CreatedAt, &run.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan takeover run: %w", err)
	}
	run.DiscoveredJSON = r.decryptDiscovered(run.DiscoveredJSON)
	return &run, nil
}
