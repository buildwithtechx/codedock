package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type PasskeyRepository interface {
	CreateCredential(ctx context.Context, cred *models.PasskeyCredential) error
	GetCredentialByID(ctx context.Context, id string) (*models.PasskeyCredential, error)
	GetCredentialByCredentialID(ctx context.Context, credentialID string) (*models.PasskeyCredential, error)
	ListCredentialsByUserID(ctx context.Context, userID string) ([]*models.PasskeyCredential, error)
	UpdateCredentialCounter(ctx context.Context, id string, counter int64, lastUsedAt time.Time) error
	RenameCredential(ctx context.Context, id, userID, name string) error
	DeleteCredential(ctx context.Context, id, userID string) error
	SaveChallenge(ctx context.Context, challenge *models.PasskeyChallenge) error
	GetAndConsumeChallenge(ctx context.Context, challenge, ceremonyType string) (*models.PasskeyChallenge, error)
}

type PasskeyRepo struct {
	db *sqlx.DB
}

func NewPasskeyRepo(db *sql.DB) *PasskeyRepo {
	return &PasskeyRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *PasskeyRepo) CreateCredential(ctx context.Context, cred *models.PasskeyCredential) error {
	if cred.ID == "" {
		cred.ID = uuid.NewString()
	}
	if cred.CreatedAt.IsZero() {
		cred.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO passkey_credentials 
	          (id, user_id, name, credential_id, public_key, counter, device_type, backed_up, transports, aaguid, created_at, last_used_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		cred.ID, cred.UserID, cred.Name, cred.CredentialID, cred.PublicKey,
		cred.Counter, cred.DeviceType, cred.BackedUp, cred.Transports, cred.AAGUID,
		cred.CreatedAt, cred.LastUsedAt,
	)
	if err != nil {
		return fmt.Errorf("create passkey credential: %w", err)
	}
	return nil
}

func (r *PasskeyRepo) GetCredentialByID(ctx context.Context, id string) (*models.PasskeyCredential, error) {
	query := `SELECT id, user_id, name, credential_id, public_key, counter, device_type, backed_up, transports, aaguid, created_at, last_used_at
	          FROM passkey_credentials WHERE id = $1`
	var cred models.PasskeyCredential
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cred.ID, &cred.UserID, &cred.Name, &cred.CredentialID, &cred.PublicKey,
		&cred.Counter, &cred.DeviceType, &cred.BackedUp, &cred.Transports, &cred.AAGUID,
		&cred.CreatedAt, &cred.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("passkey credential not found")
		}
		return nil, fmt.Errorf("get passkey credential by id: %w", err)
	}
	return &cred, nil
}

func (r *PasskeyRepo) GetCredentialByCredentialID(ctx context.Context, credentialID string) (*models.PasskeyCredential, error) {
	query := `SELECT id, user_id, name, credential_id, public_key, counter, device_type, backed_up, transports, aaguid, created_at, last_used_at
	          FROM passkey_credentials WHERE credential_id = $1`
	var cred models.PasskeyCredential
	err := r.db.QueryRowContext(ctx, query, credentialID).Scan(
		&cred.ID, &cred.UserID, &cred.Name, &cred.CredentialID, &cred.PublicKey,
		&cred.Counter, &cred.DeviceType, &cred.BackedUp, &cred.Transports, &cred.AAGUID,
		&cred.CreatedAt, &cred.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("passkey credential not found")
		}
		return nil, fmt.Errorf("get passkey credential by credential id: %w", err)
	}
	return &cred, nil
}

func (r *PasskeyRepo) ListCredentialsByUserID(ctx context.Context, userID string) ([]*models.PasskeyCredential, error) {
	query := `SELECT id, user_id, name, credential_id, public_key, counter, device_type, backed_up, transports, aaguid, created_at, last_used_at
	          FROM passkey_credentials WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list passkey credentials: %w", err)
	}
	defer rows.Close()

	var result []*models.PasskeyCredential
	for rows.Next() {
		var cred models.PasskeyCredential
		if err := rows.Scan(
			&cred.ID, &cred.UserID, &cred.Name, &cred.CredentialID, &cred.PublicKey,
			&cred.Counter, &cred.DeviceType, &cred.BackedUp, &cred.Transports, &cred.AAGUID,
			&cred.CreatedAt, &cred.LastUsedAt,
		); err != nil {
			return nil, fmt.Errorf("scan passkey credential: %w", err)
		}
		result = append(result, &cred)
	}
	return result, rows.Err()
}

func (r *PasskeyRepo) UpdateCredentialCounter(ctx context.Context, id string, counter int64, lastUsedAt time.Time) error {
	query := `UPDATE passkey_credentials SET counter = $1, last_used_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, counter, lastUsedAt, id)
	if err != nil {
		return fmt.Errorf("update passkey counter: %w", err)
	}
	return nil
}

func (r *PasskeyRepo) RenameCredential(ctx context.Context, id, userID, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("passkey name cannot be empty")
	}
	query := `UPDATE passkey_credentials SET name = $1 WHERE id = $2 AND user_id = $3`
	res, err := r.db.ExecContext(ctx, query, name, id, userID)
	if err != nil {
		return fmt.Errorf("rename passkey credential: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("passkey credential not found or unauthorized")
	}
	return nil
}

func (r *PasskeyRepo) DeleteCredential(ctx context.Context, id, userID string) error {
	query := `DELETE FROM passkey_credentials WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete passkey credential: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("passkey credential not found or unauthorized")
	}
	return nil
}

func (r *PasskeyRepo) SaveChallenge(ctx context.Context, challenge *models.PasskeyChallenge) error {
	if challenge.CreatedAt.IsZero() {
		challenge.CreatedAt = time.Now().UTC()
	}

	_, _ = r.db.ExecContext(ctx, `DELETE FROM passkey_challenges WHERE expires_at < $1`, time.Now().UTC())

	query := `INSERT INTO passkey_challenges (challenge, user_id, ceremony_type, session_data, created_at, expires_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		challenge.Challenge, challenge.UserID, challenge.CeremonyType,
		challenge.SessionData, challenge.CreatedAt, challenge.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("save passkey challenge: %w", err)
	}
	return nil
}

func (r *PasskeyRepo) GetAndConsumeChallenge(ctx context.Context, challenge, ceremonyType string) (*models.PasskeyChallenge, error) {
	now := time.Now().UTC()
	query := `DELETE FROM passkey_challenges
	          WHERE challenge = $1 AND ceremony_type = $2 AND expires_at > $3
	          RETURNING challenge, user_id, ceremony_type, session_data, created_at, expires_at`
	var pc models.PasskeyChallenge
	err := r.db.QueryRowContext(ctx, query, challenge, ceremonyType, now).
		Scan(&pc.Challenge, &pc.UserID, &pc.CeremonyType, &pc.SessionData, &pc.CreatedAt, &pc.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("challenge expired or not found")
		}
		return nil, fmt.Errorf("consume passkey challenge: %w", err)
	}
	return &pc, nil
}
