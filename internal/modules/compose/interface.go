package compose

import "context"

// Vault provides cryptographic encrypt and decrypt functions for sensitive configurations.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// ComposeStackStore defines the persistence contract for compose stacks.
type ComposeStackStore interface {
	List(context.Context, string) ([]ComposeStack, error)
	Get(context.Context, string, string) (*ComposeStack, error)
	Save(context.Context, *ComposeStack, int) error
	Claim(context.Context, string, string, int) error
	Observe(context.Context, string, string, string, string) error
	Recover(context.Context) error
}

// ComposeStackOperations defines business actions on compose stacks.
type ComposeStackOperations interface {
	List(ctx context.Context, projectID string) (any, error)
	Config(ctx context.Context, projectID, stackID string) (any, error)
	Review(ctx context.Context, projectID string, req any) (any, error)
	Save(ctx context.Context, projectID string, req any) (any, error)
	Deploy(ctx context.Context, projectID, stackID string) (any, error)
	Cancel(ctx context.Context, projectID, stackID string) error
	Restart(ctx context.Context, projectID, stackID string) error
	Stop(ctx context.Context, projectID, stackID string) error
}
