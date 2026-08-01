package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type EnvironmentRepository interface {
	Get(ctx context.Context, id string) (*models.EnvironmentConfig, error)
	ListByProject(ctx context.Context, projectID string) ([]models.EnvironmentConfig, error)
	Create(ctx context.Context, env *models.EnvironmentConfig) error
	Delete(ctx context.Context, id string) error
}

type DomainRepository interface {
	ListByService(ctx context.Context, serviceID string) ([]models.DomainConfig, error)
	ListAll(ctx context.Context) ([]models.DomainConfig, error)
	GetByID(ctx context.Context, id string) (*models.DomainConfig, error)
	Create(ctx context.Context, d *models.DomainConfig) error
	UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string) error
	Delete(ctx context.Context, id string) error
}

type EnvironmentRepo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewEnvironmentRepo(db *sql.DB) *EnvironmentRepo {
	return &EnvironmentRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *EnvironmentRepo) Get(_ context.Context, id string) (*models.EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var env models.EnvironmentConfig
	err := r.db.Get(&env, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Environment", id)
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvironmentRepo) ListByProject(_ context.Context, projectID string) ([]models.EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var envs []models.EnvironmentConfig
	err := r.db.Select(&envs, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY is_default DESC, created_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}
	if envs == nil {
		envs = make([]models.EnvironmentConfig, 0)
	}
	return envs, nil
}

func (r *EnvironmentRepo) Create(_ context.Context, env *models.EnvironmentConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	env.CreatedAt = now
	env.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO environments (id, project_id, name, is_default, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		env.ID, env.ProjectID, env.Name, env.IsDefault, env.CreatedAt, env.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}
	return nil
}

func (r *EnvironmentRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}

type DomainRepo struct {
	db *sqlx.DB
}

func NewDomainRepo(db *sql.DB) *DomainRepo {
	return &DomainRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *DomainRepo) ListByService(ctx context.Context, serviceID string) ([]models.DomainConfig, error) {
	var domains []models.DomainConfig
	err := r.db.Select(&domains, `SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix, COALESCE(dns_provision_status, 'pending') AS dns_provision_status, COALESCE(dns_provider, '') AS dns_provider, created_at, updated_at FROM domains WHERE service_id = ? ORDER BY domain_name ASC`, serviceID)
	if err != nil {
		return nil, err
	}
	if domains == nil {
		domains = make([]models.DomainConfig, 0)
	}
	return domains, nil
}

func (r *DomainRepo) ListAll(ctx context.Context) ([]models.DomainConfig, error) {
	var domains []models.DomainConfig
	err := r.db.Select(&domains, `SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix, COALESCE(dns_provision_status, 'pending') AS dns_provision_status, COALESCE(dns_provider, '') AS dns_provider, created_at, updated_at FROM domains ORDER BY domain_name ASC`)
	if err != nil {
		return nil, err
	}
	if domains == nil {
		domains = make([]models.DomainConfig, 0)
	}
	return domains, nil
}

func (r *DomainRepo) Create(_ context.Context, d *models.DomainConfig) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.DNSProvisionStatus == "" {
		d.DNSProvisionStatus = models.DNSProvisionStatusPending
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO domains (id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix, dns_provision_status, dns_provider, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ServiceID, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix, d.DNSProvisionStatus, d.DNSProvider, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *DomainRepo) UpdateDNSProvisionStatus(_ context.Context, id string, status string, provider string) error {
	if provider != "" {
		_, err := r.db.Exec(`UPDATE domains SET dns_provision_status = ?, dns_provider = ?, updated_at = ? WHERE id = ?`, status, provider, time.Now(), id)
		if err != nil {
			return fmt.Errorf("failed to update domain dns status: %w", err)
		}
		return nil
	}
	_, err := r.db.Exec(`UPDATE domains SET dns_provision_status = ?, updated_at = ? WHERE id = ?`, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update domain dns status: %w", err)
	}
	return nil
}

func (r *DomainRepo) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}

func (r *DomainRepo) GetByID(_ context.Context, id string) (*models.DomainConfig, error) {
	var domain models.DomainConfig
	err := r.db.Get(&domain, `SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix, COALESCE(dns_provision_status, 'pending') AS dns_provision_status, COALESCE(dns_provider, '') AS dns_provider, created_at, updated_at FROM domains WHERE id = ?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}
