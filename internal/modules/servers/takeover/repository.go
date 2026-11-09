package takeover

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type TakeoverRepository interface {
	Create(ctx context.Context, run *TakeoverRun) error
	GetByID(ctx context.Context, id string) (*TakeoverRun, error)
	UpdateStatus(ctx context.Context, id string, status TakeoverStatus, errMsg string) error
	UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error
	UpdateAdopted(ctx context.Context, id string, projectIDs []string) error
	ListByUser(ctx context.Context, userID string) ([]*TakeoverRun, error)
}

type postgresTakeoverRepository struct {
	db    *sql.DB
	vault Vault
}

func NewTakeoverRepository(db *sql.DB, vault Vault) TakeoverRepository {
	return &postgresTakeoverRepository{db: db, vault: vault}
}

func (r *postgresTakeoverRepository) encryptDiscovered(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *postgresTakeoverRepository) decryptDiscovered(value string) string {
	if value == "" || r.vault == nil {
		return value
	}
	if decrypted, err := r.vault.Decrypt(value); err == nil {
		return decrypted
	}
	return value
}

func (r *postgresTakeoverRepository) Create(ctx context.Context, run *TakeoverRun) error {
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		run.ID, run.UserID, run.SourceHost, run.SourcePlatform, run.Status,
		discoveredJSON, run.AdoptedProjectIDs, run.Error, run.CreatedAt, run.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create takeover run: %w", err)
	}
	return nil
}

func (r *postgresTakeoverRepository) GetByID(ctx context.Context, id string) (*TakeoverRun, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, source_host, source_platform, status, discovered_json, adopted_project_ids, error, created_at, updated_at
		FROM takeover_runs WHERE id = $1`, id)
	return r.scan(row)
}

func (r *postgresTakeoverRepository) UpdateStatus(ctx context.Context, id string, status TakeoverStatus, errMsg string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET status = $1, error = $2, updated_at = $3 WHERE id = $4`,
		status, errMsg, time.Now().UTC(), id,
	)
	return err
}

func (r *postgresTakeoverRepository) UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error {
	encrypted, err := r.encryptDiscovered(discoveredJSON)
	if err != nil {
		return fmt.Errorf("encrypt discovered stack: %w", err)
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET discovered_json = $1, updated_at = $2 WHERE id = $3`,
		encrypted, time.Now().UTC(), id,
	)
	return err
}

func (r *postgresTakeoverRepository) UpdateAdopted(ctx context.Context, id string, projectIDs []string) error {
	joined := ""
	for i, pid := range projectIDs {
		if i > 0 {
			joined += ","
		}
		joined += pid
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE takeover_runs SET adopted_project_ids = $1, updated_at = $2 WHERE id = $3`,
		joined, time.Now().UTC(), id,
	)
	return err
}

func (r *postgresTakeoverRepository) ListByUser(ctx context.Context, userID string) ([]*TakeoverRun, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, source_host, source_platform, status, discovered_json, adopted_project_ids, error, created_at, updated_at
		FROM takeover_runs WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list takeover runs: %w", err)
	}
	defer rows.Close()

	var runs []*TakeoverRun
	for rows.Next() {
		run, err := r.scanRows(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func (r *postgresTakeoverRepository) scan(row *sql.Row) (*TakeoverRun, error) {
	var run TakeoverRun
	var errStr sql.NullString
	err := row.Scan(
		&run.ID, &run.UserID, &run.SourceHost, &run.SourcePlatform, &run.Status,
		&run.DiscoveredJSON, &run.AdoptedProjectIDs, &errStr, &run.CreatedAt, &run.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan takeover run: %w", err)
	}
	if errStr.Valid {
		run.Error = errStr.String
	}
	run.DiscoveredJSON = r.decryptDiscovered(run.DiscoveredJSON)
	return &run, nil
}

func (r *postgresTakeoverRepository) scanRows(rows *sql.Rows) (*TakeoverRun, error) {
	var run TakeoverRun
	var errStr sql.NullString
	err := rows.Scan(
		&run.ID, &run.UserID, &run.SourceHost, &run.SourcePlatform, &run.Status,
		&run.DiscoveredJSON, &run.AdoptedProjectIDs, &errStr, &run.CreatedAt, &run.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan takeover run row: %w", err)
	}
	if errStr.Valid {
		run.Error = errStr.String
	}
	run.DiscoveredJSON = r.decryptDiscovered(run.DiscoveredJSON)
	return &run, nil
}
