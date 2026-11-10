package domains

import "context"

// DomainRepository persists domain configurations and bindings to services.
type DomainRepository interface {
	ListByService(ctx context.Context, serviceID string) ([]DomainConfig, error)
	ListAll(ctx context.Context) ([]DomainConfig, error)
	GetByID(ctx context.Context, id string) (*DomainConfig, error)
	Create(ctx context.Context, d *DomainConfig) error
	UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string, provisionedIP string) error
	Delete(ctx context.Context, id string) error
	BindUnassignedVerifiedDomainsToSingleService(ctx context.Context, projectID, singleServiceID string) error
}

// DomainDNSService integrates with DNS providers to provision and deprovision A records.
type DomainDNSService interface {
	ProvisionARecord(ctx context.Context, domain string) (string, string, error)
	DeprovisionARecord(ctx context.Context, domain, provider, value string) error
}
