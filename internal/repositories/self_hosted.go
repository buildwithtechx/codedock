package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
)

type SelfHostedRepo struct {
	db *sqlx.DB
}

func NewSelfHostedRepo(db *sql.DB) *SelfHostedRepo {
	return &SelfHostedRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *SelfHostedRepo) Load(ctx context.Context) (*models.SelfHostedConfig, error) {
	var stored models.SelfHostedConfig
	err := r.db.QueryRowContext(ctx, `SELECT jwt_secret, refresh_secret, telemetry_salt, tls_email, wildcard_domain FROM self_hosted_config WHERE id = 'global'`).Scan(&stored.JWTSecret, &stored.RefreshSecret, &stored.TelemetrySalt, &stored.TLSEmail, &stored.WildcardDomain)
	if err == sql.ErrNoRows {
		return &models.SelfHostedConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load self-hosted configuration: %w", err)
	}
	return &stored, nil
}

func (r *SelfHostedRepo) Save(ctx context.Context, stored *models.SelfHostedConfig) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO self_hosted_config (id, jwt_secret, refresh_secret, telemetry_salt, tls_email, wildcard_domain, updated_at)
		VALUES ('global', $1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET jwt_secret = excluded.jwt_secret, refresh_secret = excluded.refresh_secret, telemetry_salt = excluded.telemetry_salt, tls_email = excluded.tls_email, wildcard_domain = excluded.wildcard_domain, updated_at = CURRENT_TIMESTAMP`,
		stored.JWTSecret, stored.RefreshSecret, stored.TelemetrySalt, stored.TLSEmail, stored.WildcardDomain)
	if err != nil {
		return fmt.Errorf("save self-hosted configuration: %w", err)
	}
	return nil
}
