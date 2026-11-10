package adapters

import (
	"context"
	"fmt"

	"codedock/internal/modules/audit/operations"
	modclusters "codedock/internal/modules/clusters"
)

type ClusterDataServiceAdapter struct {
	Service    *modclusters.ClusterDataService
	Operations *operations.OperationService
	Apps       modclusters.BindingApps
	Vars       modclusters.BindingVars
}

func (a *ClusterDataServiceAdapter) List(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.Service.List(ctx, projectID, clusterID)
}

func (a *ClusterDataServiceAdapter) Review(ctx context.Context, userID, projectID, clusterID string, req modclusters.ClusterDataRequest) (any, error) {
	return a.Service.Review(ctx, userID, projectID, clusterID, req)
}

func (a *ClusterDataServiceAdapter) Apply(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return nil, fmt.Errorf("database operation not found")
	}
	if err := a.Service.Apply(ctx, op.UserID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) GetOperation(ctx context.Context, userID, operationID string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return nil, fmt.Errorf("database operation not found")
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) CancelOperation(ctx context.Context, userID, operationID string) error {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return fmt.Errorf("database operation not found")
	}
	return a.Operations.Cancel(ctx, op.ID, op.UserID)
}

func (a *ClusterDataServiceAdapter) Credentials(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.Service.Credentials(ctx, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) Backups(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.Service.Backups(ctx, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) VerifyConnection(ctx context.Context, projectID, clusterID, databaseID string) error {
	return a.Service.VerifyConnection(ctx, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) ReviewBinding(ctx context.Context, userID, projectID, clusterID, appID, databaseID string) (any, error) {
	if a.Apps == nil || a.Vars == nil {
		return nil, fmt.Errorf("binding workflow unavailable")
	}
	return a.Service.ReviewBinding(ctx, userID, projectID, clusterID, appID, databaseID, a.Apps, a.Vars)
}

func (a *ClusterDataServiceAdapter) ApplyBinding(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-binding" || op.UserID != userID {
		return nil, fmt.Errorf("binding operation not found")
	}
	if a.Apps == nil || a.Vars == nil {
		return nil, fmt.Errorf("binding workflow unavailable")
	}
	if err := a.Service.ApplyBinding(ctx, userID, op.ID, confirmation, a.Apps, a.Vars); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) ReviewLifecycle(ctx context.Context, userID, projectID, clusterID, databaseID, action string) (any, error) {
	return a.Service.ReviewLifecycle(ctx, userID, projectID, clusterID, modclusters.LifecycleRequest{
		DatabaseID: databaseID,
		Action:     action,
	})
}

func (a *ClusterDataServiceAdapter) ApplyLifecycle(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data-lifecycle" || op.UserID != userID {
		return nil, fmt.Errorf("lifecycle operation not found")
	}
	if err := a.Service.ApplyLifecycle(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) ReviewRedisSnapshot(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error) {
	return a.Service.ReviewRedisSnapshot(ctx, userID, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) ApplyRedisSnapshot(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("snapshot operation not found")
	}
	if err := a.Service.ApplyRedisSnapshot(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) ListRedisSnapshots(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.Service.ListRedisSnapshots(ctx, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) ReviewRedisRestore(ctx context.Context, userID, projectID, clusterID, databaseID, snapshotID, targetID string) (any, error) {
	return a.Service.ReviewRedisRestore(ctx, userID, projectID, clusterID, databaseID, snapshotID, targetID)
}

func (a *ClusterDataServiceAdapter) ApplyRedisRestore(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("restore operation not found")
	}
	if err := a.Service.ApplyRedisRestore(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterDataServiceAdapter) ReviewRedisFailover(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error) {
	return a.Service.ReviewRedisFailover(ctx, userID, projectID, clusterID, databaseID)
}

func (a *ClusterDataServiceAdapter) ApplyRedisFailover(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("failover operation not found")
	}
	if err := a.Service.ApplyRedisFailover(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}
