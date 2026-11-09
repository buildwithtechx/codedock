package http

import (
	"context"
	"fmt"

	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/operations"
)

type clusterDataServiceAdapter struct {
	service    *modclusters.ClusterDataService
	operations *operations.OperationService
	apps       modclusters.BindingApps
	vars       modclusters.BindingVars
}

func (a *clusterDataServiceAdapter) List(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.service.List(ctx, projectID, clusterID)
}

func (a *clusterDataServiceAdapter) Review(ctx context.Context, userID, projectID, clusterID string, req modclusters.ClusterDataRequest) (any, error) {
	return a.service.Review(ctx, userID, projectID, clusterID, req)
}

func (a *clusterDataServiceAdapter) Apply(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return nil, fmt.Errorf("database operation not found")
	}
	if err := a.service.Apply(ctx, op.UserID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) GetOperation(ctx context.Context, userID, operationID string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return nil, fmt.Errorf("database operation not found")
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) CancelOperation(ctx context.Context, userID, operationID string) error {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data" || op.UserID != userID {
		return fmt.Errorf("database operation not found")
	}
	return a.operations.Cancel(ctx, op.ID, op.UserID)
}

func (a *clusterDataServiceAdapter) Credentials(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.service.Credentials(ctx, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) Backups(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.service.Backups(ctx, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) VerifyConnection(ctx context.Context, projectID, clusterID, databaseID string) error {
	return a.service.VerifyConnection(ctx, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) ReviewBinding(ctx context.Context, userID, projectID, clusterID, appID, databaseID string) (any, error) {
	if a.apps == nil || a.vars == nil {
		return nil, fmt.Errorf("binding workflow unavailable")
	}
	return a.service.ReviewBinding(ctx, userID, projectID, clusterID, appID, databaseID, a.apps, a.vars)
}

func (a *clusterDataServiceAdapter) ApplyBinding(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-binding" || op.UserID != userID {
		return nil, fmt.Errorf("binding operation not found")
	}
	if a.apps == nil || a.vars == nil {
		return nil, fmt.Errorf("binding workflow unavailable")
	}
	if err := a.service.ApplyBinding(ctx, userID, op.ID, confirmation, a.apps, a.vars); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) ReviewLifecycle(ctx context.Context, userID, projectID, clusterID, databaseID, action string) (any, error) {
	return a.service.ReviewLifecycle(ctx, userID, projectID, clusterID, modclusters.LifecycleRequest{
		DatabaseID: databaseID,
		Action:     action,
	})
}

func (a *clusterDataServiceAdapter) ApplyLifecycle(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-data-lifecycle" || op.UserID != userID {
		return nil, fmt.Errorf("lifecycle operation not found")
	}
	if err := a.service.ApplyLifecycle(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) ReviewRedisSnapshot(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error) {
	return a.service.ReviewRedisSnapshot(ctx, userID, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) ApplyRedisSnapshot(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("snapshot operation not found")
	}
	if err := a.service.ApplyRedisSnapshot(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) ListRedisSnapshots(ctx context.Context, projectID, clusterID, databaseID string) (any, error) {
	return a.service.ListRedisSnapshots(ctx, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) ReviewRedisRestore(ctx context.Context, userID, projectID, clusterID, databaseID, snapshotID, targetID string) (any, error) {
	return a.service.ReviewRedisRestore(ctx, userID, projectID, clusterID, databaseID, snapshotID, targetID)
}

func (a *clusterDataServiceAdapter) ApplyRedisRestore(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("restore operation not found")
	}
	if err := a.service.ApplyRedisRestore(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterDataServiceAdapter) ReviewRedisFailover(ctx context.Context, userID, projectID, clusterID, databaseID string) (any, error) {
	return a.service.ReviewRedisFailover(ctx, userID, projectID, clusterID, databaseID)
}

func (a *clusterDataServiceAdapter) ApplyRedisFailover(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.UserID != userID {
		return nil, fmt.Errorf("failover operation not found")
	}
	if err := a.service.ApplyRedisFailover(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}
