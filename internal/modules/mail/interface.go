package mail

import "context"

// Vault provides encryption for mail provider tokens.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// Repository defines database persistence for mail servers and accounts.
type Repository interface {
	Create(ctx context.Context, server *MailServer) error
	GetByID(ctx context.Context, id string) (*MailServer, error)
	ListByOrganization(ctx context.Context, orgID string) ([]*MailServer, error)
	Delete(ctx context.Context, id string) error
}

// MailService provides high-level mail server provisioning, DKIM, and mailbox management.
type MailService interface {
	CreateServer(ctx context.Context, req *CreateMailServerRequest) (*MailServer, error)
	GetServer(ctx context.Context, id string) (*MailServer, error)
	ListServers(ctx context.Context, orgID string) ([]*MailServer, error)
	DeleteServer(ctx context.Context, id string) error
}
