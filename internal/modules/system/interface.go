package system

import "context"

// FoundingAdminFinder retrieves founding administrator credentials for local sessions.
type FoundingAdminFinder interface {
	FindFoundingAdmin(ctx context.Context) (id, email, role string, verified, totp bool, err error)
}

// TokenGenerator signs tokens for authenticated CLI sessions.
type TokenGenerator interface {
	GenerateCLIToken(id, email, role string, verified, totp bool) (string, error)
}

// SystemOperations handles system metrics, cleanup, and restart cycles.
type SystemOperations interface {
	GetStats(ctx context.Context) (any, error)
	Restart(ctx context.Context) error
	Cleanup(ctx context.Context) error
}

// UpdaterOperations manages platform version checks and self-updates.
type UpdaterOperations interface {
	GetStatus() any
	CheckUpdate(ctx context.Context) (any, error)
	ApplyUpdate(ctx context.Context) (any, error)
}
