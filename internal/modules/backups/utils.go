package backups

import "context"

type UserContext struct {
	UserID string
	Role   string
}

type ResourceProjectResolver interface {
	GetDatabaseProjectID(ctx context.Context, databaseID string) (string, error)
	GetServiceProjectID(ctx context.Context, serviceID string) (string, error)
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
}
