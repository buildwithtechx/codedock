package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type ServerlessRepository interface {
	SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error)
	GetCodeByServiceID(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error)
}

type serverlessRepo struct {
	db *sqlx.DB
}

func NewServerlessRepository(db *sql.DB) ServerlessRepository {
	return &serverlessRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *serverlessRepo) SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error) {
	existing, err := r.GetCodeByServiceID(ctx, serviceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now()
	if existing != nil {
		query := `UPDATE serverless_functions_code SET runtime = $1, code_content = $2, updated_at = $3 WHERE service_id = $4`
		_, err := r.db.ExecContext(ctx, query, runtime, codeContent, now, serviceID)
		if err != nil {
			return nil, err
		}
		existing.Runtime = runtime
		existing.CodeContent = codeContent
		existing.UpdatedAt = now
		return existing, nil
	}

	id := uuid.New().String()
	query := `INSERT INTO serverless_functions_code (id, service_id, runtime, code_content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.ExecContext(ctx, query, id, serviceID, runtime, codeContent, now, now)
	if err != nil {
		return nil, err
	}

	return &models.ServerlessFunctionCode{
		ID:          id,
		ServiceID:   serviceID,
		Runtime:     runtime,
		CodeContent: codeContent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *serverlessRepo) GetCodeByServiceID(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error) {
	query := `SELECT id, service_id, runtime, code_content, created_at, updated_at FROM serverless_functions_code WHERE service_id = $1`
	var code models.ServerlessFunctionCode
	err := r.db.GetContext(ctx, &code, query, serviceID)
	if err != nil {
		return nil, err
	}

	return &code, nil
}
