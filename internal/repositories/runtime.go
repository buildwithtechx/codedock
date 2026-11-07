package repositories

import (
	"codedock/internal/models"
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
	return &RuntimeRepo{sqlx.NewDb(db, "pgx"), vault}
}
func (r *RuntimeRepo) Get(ctx context.Context, id string) (*models.ServiceRuntime, error) {
	var runtime models.ServiceRuntime
	err := r.db.GetContext(ctx, &runtime, `SELECT * FROM service_runtimes WHERE service_id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.ServiceRuntime{ServiceID: id, RuntimeKind: "docker", Target: models.RuntimeTarget{Kind: "docker"}, Status: "CONFIGURED"}, nil
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
	if runtime.Target.Kind == "" {
		runtime.Target.Kind = "docker"
	}
	if runtime.Target.Kind != "docker" && runtime.Target.Kind != "kubernetes" && runtime.Target.Kind != "bare" {
		return fmt.Errorf("unsupported runtime destination")
	}
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
		_, err := r.db.ExecContext(ctx, `INSERT INTO service_runtimes(service_id,project_id,encrypted_config,updated_at,runtime_kind) VALUES($1,$2,$3,$4,$5)`, runtime.ServiceID, runtime.ProjectID, encrypted, now, runtime.Target.Kind)
		if err != nil {
			return fmt.Errorf("runtime target already configured or invalid: %w", err)
		}
		return nil
	}
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET encrypted_config=$1,runtime_kind=$2,revision=revision+1,status='CONFIGURED',error='',updated_at=$3 WHERE service_id=$4 AND project_id=$5 AND revision=$6 AND status NOT IN ('DEPLOYING','RECOVERING') AND encrypted_journal=''`, encrypted, runtime.Target.Kind, now, runtime.ServiceID, runtime.ProjectID, revision)
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
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET status='DEPLOYING',encrypted_journal=$1,error='',updated_at=$2 WHERE service_id=$3 AND revision=$4 AND encrypted_journal='' AND status NOT IN ('DEPLOYING','RECOVERING')`, encrypted, time.Now().UTC().Format(time.RFC3339Nano), id, revision)
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
	result, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET status=$1,error=$2,encrypted_journal=CASE WHEN $3 THEN '' ELSE encrypted_journal END,updated_at=$4 WHERE service_id=$5`, status, message, clearJournal, time.Now().UTC().Format(time.RFC3339Nano), id)
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
