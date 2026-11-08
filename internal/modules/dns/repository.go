package dns

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, record *DNSRecord) error
	GetByID(ctx context.Context, id string) (*DNSRecord, error)
	ListByDomain(ctx context.Context, domainName string) ([]*DNSRecord, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, record *DNSRecord) error
}

type repo struct {
	db *sqlx.DB
	mu sync.Mutex
}

const dnsRecordColumns = `id, domain_name, record_type, record_name, record_value, ttl, created_at, updated_at`

func NewRepository(db *sql.DB) Repository {
	return &repo{db: sqlx.NewDb(db, "pgx")}
}

func (r *repo) Create(ctx context.Context, record *DNSRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if record.ID == "" {
		record.ID = uuid.NewString()
	}
	if record.TTL == 0 {
		record.TTL = 3600
	}
	now := time.Now().UTC()
	record.CreatedAt = now
	record.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, `INSERT INTO dns_records (
		id, domain_name, record_type, record_name, record_value, ttl, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		record.ID, record.DomainName, record.RecordType, record.RecordName, record.RecordValue, record.TTL, record.CreatedAt, record.UpdatedAt)
	return err
}

func (r *repo) GetByID(ctx context.Context, id string) (*DNSRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var record DNSRecord
	err := r.db.GetContext(ctx, &record, `SELECT `+dnsRecordColumns+` FROM dns_records WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRecordNotFound
	}
	return &record, err
}

func (r *repo) ListByDomain(ctx context.Context, domainName string) ([]*DNSRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []*DNSRecord
	var err error
	if domainName == "" {
		err = r.db.SelectContext(ctx, &list, `SELECT `+dnsRecordColumns+` FROM dns_records ORDER BY created_at ASC`)
	} else {
		err = r.db.SelectContext(ctx, &list, `SELECT `+dnsRecordColumns+` FROM dns_records WHERE domain_name = $1 ORDER BY created_at ASC`, domainName)
	}
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*DNSRecord, 0)
	}
	return list, nil
}

func (r *repo) Update(ctx context.Context, record *DNSRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `UPDATE dns_records SET record_type = $1, record_name = $2, record_value = $3, ttl = $4, updated_at = $5 WHERE id = $6`,
		record.RecordType, record.RecordName, record.RecordValue, record.TTL, record.UpdatedAt, record.ID)
	return err
}

func (r *repo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `DELETE FROM dns_records WHERE id = $1`, id)
	return err
}
