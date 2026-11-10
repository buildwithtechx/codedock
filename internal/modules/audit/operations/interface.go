package operations

import "context"

// Vault provides encryption for operation payload snapshots.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// OperationStore defines database persistence for two-man rule reviewable operations.
type OperationStore interface {
	Create(context.Context, *Operation) error
	Get(context.Context, string) (*Operation, error)
	List(context.Context, string, string) ([]Operation, error)
	Claim(context.Context, string, string) error
	Observe(context.Context, string, string, string, string, string) error
	Recover(context.Context) error
}

// OperationsService defines high-level management of operations.
type OperationsService interface {
	Get(ctx context.Context, id string) (*Operation, error)
	List(ctx context.Context, project, user string) ([]Operation, error)
	Cancel(ctx context.Context, id, user string) error
}
