package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type RuntimeRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewRuntimeRepo(db *sql.DB, vault Vault) *RuntimeRepo {
	return &RuntimeRepo{sqlx.NewDb(db, "sqlite"), vault}
}
func (r *RuntimeRepo) Get(ctx context.Context, id string) (*models.ServiceRuntime, error) {
	var runtime models.ServiceRuntime
	err := r.db.GetContext(ctx, &runtime, `SELECT * FROM service_runtimes WHERE service_id=?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.ServiceRuntime{ServiceID: id, Target: models.RuntimeTarget{Kind: "docker"}, Status: "CONFIGURED"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load runtime target: %w", err)
	}
	config, err := r.vault.Decrypt(runtime.Config)
	if err != nil {
		return nil, fmt.Errorf("decrypt runtime target: %w", err)
	}
	if err := json.Unmarshal([]byte(config), &runtime.Target); err != nil {
		return nil, fmt.Errorf("decode runtime target: %w", err)
	}
	if runtime.Journal != "" {
		journal, err := r.vault.Decrypt(runtime.Journal)
		if err != nil {
			return nil, fmt.Errorf("decrypt runtime recovery journal: %w", err)
		}
		runtime.Journal = journal
	}
	return &runtime, nil
}
func (r *RuntimeRepo) Save(ctx context.Context, runtime *models.ServiceRuntime, revision int) error {
	config, err := json.Marshal(runtime.Target)
	if err != nil {
		return err
	}
	encrypted, err := r.vault.Encrypt(string(config))
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if revision == 0 {
		_, err := r.db.ExecContext(ctx, `INSERT INTO service_runtimes(service_id,project_id,encrypted_config,updated_at) VALUES(?,?,?,?)`, runtime.ServiceID, runtime.ProjectID, encrypted, now)
		if err != nil {
			return fmt.Errorf("runtime target already configured or invalid: %w", err)
		}
		return nil
	}
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET encrypted_config=?,revision=revision+1,status='CONFIGURED',error='',updated_at=? WHERE service_id=? AND project_id=? AND revision=? AND status NOT IN ('DEPLOYING','RECOVERING') AND encrypted_journal=''`, encrypted, now, runtime.ServiceID, runtime.ProjectID, revision)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("runtime changed or recovery is pending")
	}
	return nil
}
func (r *RuntimeRepo) Begin(ctx context.Context, id string, revision int, journal string) error {
	encrypted, err := r.vault.Encrypt(journal)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET status='DEPLOYING',encrypted_journal=?,error='',updated_at=? WHERE service_id=? AND revision=? AND encrypted_journal='' AND status NOT IN ('DEPLOYING','RECOVERING')`, encrypted, time.Now().UTC().Format(time.RFC3339Nano), id, revision)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("another runtime operation or recovery owns this service")
	}
	return nil
}
func (r *RuntimeRepo) Observe(ctx context.Context, id, status, message string, clearJournal bool) error {
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET status=?,error=?,encrypted_journal=CASE WHEN ? THEN '' ELSE encrypted_journal END,updated_at=? WHERE service_id=?`, status, message, clearJournal, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("runtime target disappeared")
	}
	return nil
}
func (r *RuntimeRepo) Pending(ctx context.Context) ([]string, error) {
	result := []string{}
	if err := r.db.SelectContext(ctx, &result, `SELECT service_id FROM service_runtimes WHERE encrypted_journal!='' ORDER BY updated_at`); err != nil {
		return nil, err
	}
	return result, nil
}
