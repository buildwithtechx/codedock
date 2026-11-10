package environments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EnvironmentRepo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewEnvironmentRepo(db *sql.DB) *EnvironmentRepo {
	return &EnvironmentRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *EnvironmentRepo) Get(ctx context.Context, id string) (*EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var env EnvironmentConfig
	err := r.db.GetContext(ctx, &env, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEnvironmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvironmentRepo) ListByProject(ctx context.Context, projectID string) ([]EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var envs []EnvironmentConfig
	err := r.db.SelectContext(ctx, &envs, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = $1 ORDER BY is_default DESC, created_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}
	if envs == nil {
		envs = make([]EnvironmentConfig, 0)
	}
	return envs, nil
}

func (r *EnvironmentRepo) Create(ctx context.Context, env *EnvironmentConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	env.CreatedAt = now
	env.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO environments (id, project_id, name, is_default, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		env.ID, env.ProjectID, env.Name, env.IsDefault, env.CreatedAt, env.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}
	return nil
}

func (r *EnvironmentRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM environments WHERE id = $1`, id)
	return err
}

type EnvVarRepo struct {
	mu   sync.RWMutex
	vars map[string]map[string]string
}

func NewEnvVarRepo() *EnvVarRepo {
	return &EnvVarRepo{vars: make(map[string]map[string]string)}
}

func (r *EnvVarRepo) GetVars(_ context.Context, projectID string) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	vars, ok := r.vars[projectID]
	if !ok {
		return make(map[string]string), nil
	}
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = v
	}
	return result, nil
}

func (r *EnvVarRepo) SetVar(_ context.Context, projectID, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.vars[projectID]; !ok {
		r.vars[projectID] = make(map[string]string)
	}
	r.vars[projectID][key] = value
	return nil
}
