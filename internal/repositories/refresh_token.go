package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepository interface {
	StoreToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, tokenHash string) (bool, error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	PruneExpired(ctx context.Context) error
}

type RefreshTokenRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: sqlx.NewDb(db, "sqlite")}
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func (r *RefreshTokenRepo) StoreToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO refresh_token_revocations (id, token_hash, user_id, expires_at)
		 VALUES (?, ?, ?, ?)`,
		uuid.New().String(), tokenHash, userID, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepo) IsRevoked(ctx context.Context, tokenHash string) (bool, error) {
	var revokedAt *time.Time
	err := r.db.QueryRowContext(ctx,
		`SELECT revoked_at FROM refresh_token_revocations WHERE token_hash = ?`, tokenHash,
	).Scan(&revokedAt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check revocation: %w", err)
	}
	if revokedAt != nil {
		return true, nil
	}
	return false, nil
}

func (r *RefreshTokenRepo) RevokeToken(ctx context.Context, tokenHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_token_revocations SET revoked_at = CURRENT_TIMESTAMP WHERE token_hash = ?`, tokenHash,
	)
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_token_revocations SET revoked_at = CURRENT_TIMESTAMP WHERE user_id = ? AND revoked_at IS NULL`, userID,
	)
	if err != nil {
		return fmt.Errorf("revoke all tokens for user: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepo) PruneExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_token_revocations WHERE expires_at < ?`, time.Now())
	if err != nil {
		return fmt.Errorf("prune expired tokens: %w", err)
	}
	return nil
}
