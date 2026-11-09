package http

import (
	"context"
	"fmt"

	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/monitoring"
)

type clusterServiceAdapter struct {
	service    *modclusters.ClusterService
	operations *monitoring.OperationService
}

func (a *clusterServiceAdapter) List(ctx context.Context, projectID string) (any, error) {
	return a.service.List(ctx, projectID)
}

func (a *clusterServiceAdapter) Review(ctx context.Context, userID, projectID string, req modclusters.ClusterReviewRequest) (any, error) {
	return a.service.Review(ctx, userID, projectID, req)
}

func (a *clusterServiceAdapter) Apply(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return nil, fmt.Errorf("cluster operation not found")
	}
	if err := a.service.Apply(ctx, userID, operationID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *clusterServiceAdapter) GetOperation(ctx context.Context, userID, operationID string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return nil, fmt.Errorf("cluster operation not found")
	}
	return op, nil
}

func (a *clusterServiceAdapter) CancelOperation(ctx context.Context, userID, operationID string) error {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return fmt.Errorf("cluster operation not found")
	}
	return a.operations.Cancel(ctx, operationID, userID)
}

func (a *clusterServiceAdapter) StorageInventory(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.service.StorageInventory(ctx, projectID, clusterID)
}

func (a *clusterServiceAdapter) StorageHealth(ctx context.Context, projectID, clusterID string) error {
	return a.service.StorageHealth(ctx, projectID, clusterID)
}

func (a *clusterServiceAdapter) NetworkPreflight(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.service.NetworkPreflight(ctx, projectID, clusterID)
}

func (a *clusterServiceAdapter) ReviewStorageSetup(ctx context.Context, userID, projectID, clusterID string, req modclusters.StorageSetupRequest) (any, error) {
	return a.service.ReviewStorageSetup(ctx, userID, projectID, clusterID, req)
}

func (a *clusterServiceAdapter) ApplyStorageSetup(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-storage" || op.UserID != userID {
		return nil, fmt.Errorf("storage operation not found")
	}
	if err := a.service.ApplyStorageSetup(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}
