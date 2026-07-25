package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
)

type ProjectRepository interface {
	ListByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]models.ProjectConfig, int, error)
	ListAll(ctx context.Context, limit, offset int) ([]models.ProjectConfig, int, error)
	Get(ctx context.Context, id string) (*models.ProjectConfig, error)
	GetByOrganization(ctx context.Context, id, organizationID string) (*models.ProjectConfig, error)
	CountByUser(ctx context.Context, userID string) (int, error)
	Create(ctx context.Context, p *models.ProjectConfig) error
	Delete(ctx context.Context, id string) error
}

type EnvRepository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}

type ProjectRepo struct {
	db           *sqlx.DB
	environments EnvironmentRepository
}

func NewProjectRepo(db *sql.DB, envRepo EnvironmentRepository) *ProjectRepo {
	return &ProjectRepo{db: sqlx.NewDb(db, "sqlite"), environments: envRepo}
}

func (r *ProjectRepo) ListByOrganization(_ context.Context, organizationID string, limit, offset int) ([]models.ProjectConfig, int, error) {
	var total int
	var err error
	var projects []models.ProjectConfig

	if organizationID == "none" {
		if err = r.db.Get(&total, `SELECT COUNT(*) FROM projects WHERE organization_id IS NULL`); err != nil {
			return nil, 0, err
		}
		err = r.db.Select(&projects, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE organization_id IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	} else {
		if err = r.db.Get(&total, `SELECT COUNT(*) FROM projects WHERE organization_id = ?`, organizationID); err != nil {
			return nil, 0, err
		}
		err = r.db.Select(&projects, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE organization_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, organizationID, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	if projects == nil {
		projects = make([]models.ProjectConfig, 0)
	}
	return projects, total, nil
}

func (r *ProjectRepo) ListAll(_ context.Context, limit, offset int) ([]models.ProjectConfig, int, error) {
	var total int
	var err error
	var projects []models.ProjectConfig

	if err = r.db.Get(&total, `SELECT COUNT(*) FROM projects`); err != nil {
		return nil, 0, err
	}
	err = r.db.Select(&projects, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	if projects == nil {
		projects = make([]models.ProjectConfig, 0)
	}
	return projects, total, nil
}

func (r *ProjectRepo) CountByUser(ctx context.Context, userID string) (int, error) {
	var total int
	query := `
		SELECT COUNT(p.id) 
		FROM projects p 
		INNER JOIN organization_members om ON p.organization_id = om.organization_id 
		WHERE om.user_id = ?`
	err := r.db.GetContext(ctx, &total, query, userID)
	return total, err
}

func (r *ProjectRepo) Get(_ context.Context, id string) (*models.ProjectConfig, error) {
	var p models.ProjectConfig
	err := r.db.Get(&p, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) GetByOrganization(_ context.Context, id, organizationID string) (*models.ProjectConfig, error) {
	var p models.ProjectConfig
	var err error
	if organizationID == "none" {
		err = r.db.Get(&p, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE id = ? AND organization_id IS NULL`, id)
	} else {
		err = r.db.Get(&p, `SELECT id, COALESCE(organization_id, '') AS organization_id, name, COALESCE(server_id, '') AS server_id, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE id = ? AND organization_id = ?`, id, organizationID)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) Create(ctx context.Context, p *models.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	var serverID any
	if p.ServerID != "" {
		serverID = p.ServerID
	}
	var orgID any
	if p.OrganizationID != "" {
		orgID = p.OrganizationID
	}
	_, err := r.db.Exec(
		`INSERT INTO projects (id, organization_id, server_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, orgID, serverID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return err
	}
	defaultEnv := &models.EnvironmentConfig{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.environments.Create(ctx, defaultEnv)
}

func (r *ProjectRepo) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

type EnvRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewEnvRepo(db *sql.DB, vault Vault) *EnvRepo {
	return &EnvRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func (r *EnvRepo) GetVars(_ context.Context, projectID string) (map[string]string, error) {
	rows, err := r.db.Query(`SELECT key, encrypted_value FROM env_vars WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envs := make(map[string]string)
	for rows.Next() {
		var key, encrypted string
		if err := rows.Scan(&key, &encrypted); err != nil {
			return nil, err
		}
		plaintext, err := r.vault.Decrypt(encrypted)
		if err != nil {
			continue
		}
		envs[key] = plaintext
	}
	return envs, rows.Err()
}

func (r *EnvRepo) SetVar(_ context.Context, projectID, key, plaintextValue string) error {
	encrypted, err := r.vault.Encrypt(plaintextValue)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO env_vars (id, project_id, key, encrypted_value, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, key) DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at`,
		uuid.NewString(), projectID, key, encrypted, now, now,
	)
	return err
}
