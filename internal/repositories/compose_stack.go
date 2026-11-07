package repositories

import (
	"codedock/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ComposeStackRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewComposeStackRepo(db *sql.DB, vault Vault) *ComposeStackRepo {
	return &ComposeStackRepo{db: sqlx.NewDb(db, "pgx"), vault: vault}
}

func (r *ComposeStackRepo) List(ctx context.Context, projectID string) ([]models.ComposeStack, error) {
	stacks := []models.ComposeStack{}
	err := r.db.SelectContext(ctx, &stacks, `SELECT id, project_id, environment_id, name, revision, status, error, results, updated_at FROM compose_stacks WHERE project_id = $1 ORDER BY name`, projectID)
	return stacks, err
}

func (r *ComposeStackRepo) Get(ctx context.Context, projectID, id string) (*models.ComposeStack, error) {
	var stack models.ComposeStack
	if err := r.db.GetContext(ctx, &stack, `SELECT * FROM compose_stacks WHERE project_id = $1 AND id = $2`, projectID, id); err != nil {
		return nil, fmt.Errorf("load stack: %w", err)
	}
	config, err := r.vault.Decrypt(stack.Config)
	if err != nil {
		return nil, fmt.Errorf("decrypt stack: %w", err)
	}
	stack.Config = config
	return &stack, nil
}

func (r *ComposeStackRepo) Save(ctx context.Context, stack *models.ComposeStack, previousRevision int) error {
	encrypted, err := r.vault.Encrypt(stack.Config)
	if err != nil {
		return fmt.Errorf("encrypt stack: %w", err)
	}
	stack.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	stack.Revision = previousRevision + 1
	var result sql.Result
	if previousRevision == 0 {
		result, err = r.db.ExecContext(ctx, `INSERT INTO compose_stacks(id,project_id,environment_id,name,encrypted_config,revision,updated_at) VALUES($1,$2,$3,$4,$5,1,$6)`, stack.ID, stack.ProjectID, stack.EnvironmentID, stack.Name, encrypted, stack.UpdatedAt)
	} else {
		result, err = r.db.ExecContext(ctx, `UPDATE compose_stacks SET name=$1, encrypted_config=$2, revision=revision+1, updated_at=$3, status='saved', error='' WHERE id=$4 AND project_id=$5 AND environment_id=$6 AND revision=$7 AND status NOT IN ('PREPARING','BUILDING','STARTING','READINESS')`, stack.Name, encrypted, stack.UpdatedAt, stack.ID, stack.ProjectID, stack.EnvironmentID, previousRevision)
	}
	if err != nil {
		return fmt.Errorf("save stack: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("stack changed or is deploying; reload and review again")
	}
	return nil
}

func (r *ComposeStackRepo) Claim(ctx context.Context, projectID, id string, revision int) error {
	result, err := r.db.ExecContext(ctx, `UPDATE compose_stacks SET status='PREPARING', error='', updated_at=$1 WHERE project_id=$2 AND id=$3 AND revision=$4 AND status NOT IN ('PREPARING','BUILDING','STARTING','READINESS')`, time.Now().UTC().Format(time.RFC3339Nano), projectID, id, revision)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("stack already deploying or unavailable")
	}
	return nil
}

func (r *ComposeStackRepo) Observe(ctx context.Context, id, status, message, results string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE compose_stacks SET status=$1, error=$2, results=$3, updated_at=$4 WHERE id=$5`, status, message, results, time.Now().UTC().Format(time.RFC3339Nano), id)
	return err
}

func (r *ComposeStackRepo) Recover(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE compose_stacks SET status='INTERRUPTED', error='Daemon restarted during deployment. Inspect service states before retrying; activation may be partial.' WHERE status IN ('PREPARING','BUILDING','STARTING','READINESS')`)
	return err
}
