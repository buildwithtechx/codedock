package clusters

import (
	"context"
)

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type OperationsService interface {
	Review(ctx context.Context, user, project, kind, target, payload, snapshot, effects string) (*OperationReview, error)
	Get(ctx context.Context, id string) (*Operation, error)
	Apply(ctx context.Context, id, user, confirmation, snapshot string, fn func(ctx context.Context, op *Operation, progress func(string, string) error) error) error
}

type Service interface {
	List(ctx context.Context, projectID string) (any, error)
	Review(ctx context.Context, userID, projectID string, req ClusterReviewRequest) (any, error)
	Apply(ctx context.Context, userID, operationID, confirmation string) (any, error)
	GetOperation(ctx context.Context, userID, operationID string) (any, error)
	CancelOperation(ctx context.Context, userID, operationID string) error
	StorageInventory(ctx context.Context, projectID, clusterID string) (any, error)
	StorageHealth(ctx context.Context, projectID, clusterID string) error
	NetworkPreflight(ctx context.Context, projectID, clusterID string) (any, error)
	ReviewStorageSetup(ctx context.Context, userID, projectID, clusterID string, req StorageSetupRequest) (any, error)
	ApplyStorageSetup(ctx context.Context, userID, operationID, confirmation string) (any, error)
}

type DataService interface {
	List(ctx context.Context, projectID, clusterID string) (any, error)
	Review(ctx context.Context, userID, projectID, clusterID string, req ClusterDataRequest) (any, error)
	Apply(ctx context.Context, userID, operationID, confirmation string) (any, error)
	GetOperation(ctx context.Context, userID, operationID string) (any, error)
	CancelOperation(ctx context.Context, userID, operationID string) error
	Credentials(ctx context.Context, projectID, clusterID, databaseID string) (any, error)
	Backups(ctx context.Context, projectID, clusterID, databaseID string) (any, error)
	VerifyConnection(ctx context.Context, projectID, clusterID, databaseID string) error
	ReviewBinding(ctx context.Context, userID, projectID, clusterID, appID, databaseID string) (any, error)
	ApplyBinding(ctx context.Context, userID, operationID, confirmation string) (any, error)
	ReviewLifecycle(ctx context.Context, userID, projectID, clusterID, databaseID, action string) (any, error)
	ApplyLifecycle(ctx context.Context, userID, operationID, confirmation string) (any, error)
	ReviewRedisSnapshot(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error)
	ApplyRedisSnapshot(ctx context.Context, userID, operationID, confirmation string) (any, error)
	ListRedisSnapshots(ctx context.Context, projectID, clusterID, databaseID string) (any, error)
	ReviewRedisRestore(ctx context.Context, userID, projectID, clusterID, databaseID, snapshotID, targetID string) (any, error)
	ApplyRedisRestore(ctx context.Context, userID, operationID, confirmation string) (any, error)
	ReviewRedisFailover(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error)
	ApplyRedisFailover(ctx context.Context, userID, operationID, confirmation string) (any, error)
}
