package servers

import "context"

// Vault provides encryption for SSH keys and remote server connection credentials.
type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// Repository defines database persistence for attached servers and nodes.
type Repository interface {
	Create(ctx context.Context, server *Server) error
	Update(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id string) (*Server, error)
	GetByToken(ctx context.Context, token string) (*Server, error)
	ListByUser(ctx context.Context, userID string, search string) ([]*Server, error)
	SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error)
	UpdateStatus(ctx context.Context, id string, status ServerStatus) error
	UpdateMetrics(ctx context.Context, id string, metricsJSON []byte) error
	Delete(ctx context.Context, id string) error
}
