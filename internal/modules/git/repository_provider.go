package git

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type gitProviderRepo struct {
	db    *sqlx.DB
	vault Encryptor
}

func NewGitProviderRepository(db *sql.DB, vault Encryptor) GitProviderRepository {
	return &gitProviderRepo{db: sqlx.NewDb(db, "pgx"), vault: vault}
}

func (r *gitProviderRepo) SaveProvider(ctx context.Context, gp *GitProviderConfig) error {
	if gp.ID == "" {
		gp.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	gp.CreatedAt = now
	gp.UpdatedAt = now

	encryptedToken := gp.AccessToken
	if r.vault != nil && gp.AccessToken != "" {
		enc, err := r.vault.Encrypt(gp.AccessToken)
		if err != nil {
			return err
		}
		encryptedToken = enc
	}

	query := `INSERT INTO user_git_providers (id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(user_id, provider) DO UPDATE SET encrypted_access_token = excluded.encrypted_access_token, account_name = excluded.account_name, updated_at = excluded.updated_at`

	_, err := r.db.ExecContext(ctx, query, gp.ID, gp.UserID, gp.Provider, encryptedToken, gp.AccountName, gp.CreatedAt, gp.UpdatedAt)
	return err
}

func (r *gitProviderRepo) GetProvider(ctx context.Context, userID, provider string) (*GitProviderConfig, error) {
	if userID == "" {
		return r.GetAnyProviderByType(ctx, provider)
	}
	var gp GitProviderConfig
	err := r.db.GetContext(ctx, &gp, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = $1 AND provider = $2`, userID, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if r.vault != nil && gp.AccessToken != "" {
		decrypted, err := r.vault.Decrypt(gp.AccessToken)
		if err == nil {
			gp.AccessToken = decrypted
		}
	}
	return &gp, nil
}

func (r *gitProviderRepo) GetAnyProviderByType(ctx context.Context, provider string) (*GitProviderConfig, error) {
	var gp GitProviderConfig
	err := r.db.GetContext(ctx, &gp, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE provider = $1 LIMIT 1`, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	if r.vault != nil && gp.AccessToken != "" {
		decrypted, err := r.vault.Decrypt(gp.AccessToken)
		if err == nil {
			gp.AccessToken = decrypted
		}
	}
	return &gp, nil
}

func (r *gitProviderRepo) ListProvidersByUser(ctx context.Context, userID string) ([]*GitProviderConfig, error) {
	var list []*GitProviderConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	for _, gp := range list {
		if r.vault != nil && gp.AccessToken != "" {
			if dec, err := r.vault.Decrypt(gp.AccessToken); err == nil {
				gp.AccessToken = dec
			}
		}
	}
	return list, nil
}

func (r *gitProviderRepo) DeleteProvider(ctx context.Context, userID, provider string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_git_providers WHERE user_id = $1 AND provider = $2`, userID, provider)
	return err
}
