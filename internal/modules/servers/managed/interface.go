package managed

import "context"

// Vault encrypts and decrypts provider access tokens.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// ManagedRepository defines the persistence contract for cloud provider credentials, quotas, and server links.
type ManagedRepository interface {
	CreateCredential(ctx context.Context, credential *ManagedCredential) error
	GetCredential(ctx context.Context, id string) (*ManagedCredential, error)
	ListCredentialsByOrg(ctx context.Context, orgID string) ([]*ManagedCredential, error)
	DeleteCredential(ctx context.Context, id string) error
	GetQuota(ctx context.Context, orgID string) (*ManagedQuota, error)
	SetQuota(ctx context.Context, quota *ManagedQuota) error
	SaveLink(ctx context.Context, link *ManagedServerLink) error
	GetLink(ctx context.Context, serverID string) (*ManagedServerLink, error)
	DeleteLink(ctx context.Context, serverID string) error
	ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error)
	ListProvisioningServerIDs(ctx context.Context) ([]string, error)
}
