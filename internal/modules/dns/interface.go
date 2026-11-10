package dns

import "context"

// Repository defines persistence for DNS records.
type Repository interface {
	Create(ctx context.Context, record *DNSRecord) error
	GetByID(ctx context.Context, id string) (*DNSRecord, error)
	ListByDomain(ctx context.Context, domainName string) ([]*DNSRecord, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, record *DNSRecord) error
}

// DNSService provides high-level domain record management and external DNS provider propagation.
type DNSService interface {
	CreateRecord(ctx context.Context, record *DNSRecord) error
	GetRecord(ctx context.Context, id string) (*DNSRecord, error)
	ListRecords(ctx context.Context, domain string) ([]*DNSRecord, error)
	UpdateRecord(ctx context.Context, record *DNSRecord) error
	DeleteRecord(ctx context.Context, id string) error
	VerifyRecord(ctx context.Context, id string) (bool, error)
}
