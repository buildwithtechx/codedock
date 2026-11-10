package ingress

import "context"

// IngressRepository persists edge ingress tunnels and hostname mappings.
type IngressRepository interface {
	ListConnections(ctx context.Context, serverID string) ([]*IngressConnection, error)
	GetConnection(ctx context.Context, id string) (*IngressConnection, error)
	SaveConnection(ctx context.Context, conn *IngressConnection) error
	DeleteConnection(ctx context.Context, id string) error

	GetHostname(ctx context.Context, domainID string) (*IngressHostname, error)
	SaveHostname(ctx context.Context, host *IngressHostname) error
	DeleteHostname(ctx context.Context, domainID string) error
}

// IngressProvider manages edge tunnel creation and discovery via external APIs.
type IngressProvider interface {
	ListAccounts(ctx context.Context, token string) ([]IngressDiscoveryAccount, error)
	ListTunnels(ctx context.Context, accountID, token string) ([]IngressDiscoveryTunnel, error)
	CreateTunnel(ctx context.Context, accountID, name, secret, token string) (*IngressDiscoveryTunnel, error)
	DeleteTunnel(ctx context.Context, accountID, tunnelID, token string) error
}
