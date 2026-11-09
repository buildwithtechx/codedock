package adapters

import (
	"context"
	"fmt"

	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/operations"
)

type ClusterServiceAdapter struct {
	Service    *modclusters.ClusterService
	Operations *operations.OperationService
}

func (a *ClusterServiceAdapter) List(ctx context.Context, projectID string) (any, error) {
	return a.Service.List(ctx, projectID)
}

func (a *ClusterServiceAdapter) Review(ctx context.Context, userID, projectID string, req modclusters.ClusterReviewRequest) (any, error) {
	return a.Service.Review(ctx, userID, projectID, req)
}

func (a *ClusterServiceAdapter) Apply(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return nil, fmt.Errorf("cluster operation not found")
	}
	if err := a.Service.Apply(ctx, userID, operationID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}

func (a *ClusterServiceAdapter) GetOperation(ctx context.Context, userID, operationID string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return nil, fmt.Errorf("cluster operation not found")
	}
	return op, nil
}

func (a *ClusterServiceAdapter) CancelOperation(ctx context.Context, userID, operationID string) error {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster" || op.UserID != userID {
		return fmt.Errorf("cluster operation not found")
	}
	return a.Operations.Cancel(ctx, operationID, userID)
}

func (a *ClusterServiceAdapter) StorageInventory(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.Service.StorageInventory(ctx, projectID, clusterID)
}

func (a *ClusterServiceAdapter) StorageHealth(ctx context.Context, projectID, clusterID string) error {
	return a.Service.StorageHealth(ctx, projectID, clusterID)
}

func (a *ClusterServiceAdapter) NetworkPreflight(ctx context.Context, projectID, clusterID string) (any, error) {
	return a.Service.NetworkPreflight(ctx, projectID, clusterID)
}

func (a *ClusterServiceAdapter) ReviewStorageSetup(ctx context.Context, userID, projectID, clusterID string, req modclusters.StorageSetupRequest) (any, error) {
	return a.Service.ReviewStorageSetup(ctx, userID, projectID, clusterID, req)
}

func (a *ClusterServiceAdapter) ApplyStorageSetup(ctx context.Context, userID, operationID, confirmation string) (any, error) {
	op, err := a.Operations.Get(ctx, operationID)
	if err != nil || op.Kind != "cluster-storage" || op.UserID != userID {
		return nil, fmt.Errorf("storage operation not found")
	}
	if err := a.Service.ApplyStorageSetup(ctx, userID, op.ID, confirmation); err != nil {
		return nil, err
	}
	return op, nil
}
