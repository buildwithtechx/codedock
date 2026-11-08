package domains

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DomainRepository interface {
	ListByService(ctx context.Context, serviceID string) ([]DomainConfig, error)
	ListAll(ctx context.Context) ([]DomainConfig, error)
	GetByID(ctx context.Context, id string) (*DomainConfig, error)
	Create(ctx context.Context, d *DomainConfig) error
	UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string, provisionedIP string) error
	Delete(ctx context.Context, id string) error
	BindUnassignedVerifiedDomainsToSingleService(ctx context.Context, projectID, singleServiceID string) error
}

type DomainRepo struct {
	db *sqlx.DB
}

func NewDomainRepo(db *sql.DB) *DomainRepo {
	return &DomainRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *DomainRepo) ListByService(ctx context.Context, serviceID string) ([]DomainConfig, error) {
	var domains []DomainConfig
	err := r.db.SelectContext(ctx, &domains, `
		SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix,
		COALESCE(dns_provision_status, 'pending') AS dns_provision_status,
		COALESCE(dns_provider, '') AS dns_provider,
		COALESCE(dns_provisioned_ip, '') AS dns_provisioned_ip,
		created_at, updated_at
		FROM domains WHERE service_id = $1 ORDER BY domain_name ASC`, serviceID)
	if err != nil {
		return nil, err
	}
	if domains == nil {
		domains = make([]DomainConfig, 0)
	}
	return domains, nil
}

func (r *DomainRepo) ListAll(ctx context.Context) ([]DomainConfig, error) {
	var domains []DomainConfig
	err := r.db.SelectContext(ctx, &domains, `
		SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix,
		COALESCE(dns_provision_status, 'pending') AS dns_provision_status,
		COALESCE(dns_provider, '') AS dns_provider,
		COALESCE(dns_provisioned_ip, '') AS dns_provisioned_ip,
		created_at, updated_at
		FROM domains ORDER BY domain_name ASC`)
	if err != nil {
		return nil, err
	}
	if domains == nil {
		domains = make([]DomainConfig, 0)
	}
	return domains, nil
}

func (r *DomainRepo) Create(ctx context.Context, d *DomainConfig) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	d.DNSProvisionStatus = DNSProvisionStatusPending
	d.DNSProvider = ""
	d.DNSProvisionedIP = ""
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO domains (id, service_id, hostname, domain_name, redirect_to, ssl_cert_status, path_prefix, dns_provision_status, dns_provider, dns_provisioned_ip, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		d.ID, d.ServiceID, d.DomainName, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix,
		d.DNSProvisionStatus, d.DNSProvider, d.DNSProvisionedIP, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *DomainRepo) UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string, provisionedIP string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE domains SET dns_provision_status = $1, dns_provider = $2, dns_provisioned_ip = $3, updated_at = $4
		WHERE id = $5`, status, provider, provisionedIP, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update domain dns status: %w", err)
	}
	return nil
}

func (r *DomainRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM domains WHERE id = $1`, id)
	return err
}

func (r *DomainRepo) GetByID(ctx context.Context, id string) (*DomainConfig, error) {
	var domain DomainConfig
	err := r.db.GetContext(ctx, &domain, `
		SELECT id, service_id, domain_name, redirect_to, ssl_cert_status, path_prefix,
		COALESCE(dns_provision_status, 'pending') AS dns_provision_status,
		COALESCE(dns_provider, '') AS dns_provider,
		COALESCE(dns_provisioned_ip, '') AS dns_provisioned_ip,
		created_at, updated_at
		FROM domains WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}

func (r *DomainRepo) BindUnassignedVerifiedDomainsToSingleService(ctx context.Context, projectID, singleServiceID string) error {
	if projectID == "" || singleServiceID == "" {
		return nil
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE domains
		SET service_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE project_id = $2
		  AND (service_id = '' OR service_id IS NULL)
		  AND (dns_verification_status = 'verified' OR dns_provision_status = 'verified')
	`, singleServiceID, projectID)
	return err
}
