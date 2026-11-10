package takeover

import "context"

// Vault encrypts discovered container metadata.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// TakeoverRepository persists takeover runs and discovered workloads.
type TakeoverRepository interface {
	Create(ctx context.Context, run *TakeoverRun) error
	GetByID(ctx context.Context, id string) (*TakeoverRun, error)
	UpdateStatus(ctx context.Context, id string, status TakeoverStatus, errMsg string) error
	UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error
	UpdateAdopted(ctx context.Context, id string, projectIDs []string) error
	ListByUser(ctx context.Context, userID string) ([]*TakeoverRun, error)
}
