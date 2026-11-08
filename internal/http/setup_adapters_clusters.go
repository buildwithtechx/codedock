package http

import (
	"context"
	"fmt"

	"codedock/internal/models"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/services/clusters"
	"codedock/internal/services/operations"
)

type clusterServiceAdapter struct {
	service    *clusters.Service
	operations *operations.Service
}

func (a *clusterServiceAdapter) List(ctx context.Context, projectID string) (any, error) {
	return a.service.List(ctx, projectID)
}

func (a *clusterServiceAdapter) Review(ctx context.Context, userID, projectID string, req modclusters.ClusterReviewRequest) (any, error) {
	legacyNodes := make([]models.ClusterNode, len(req.Cluster.Nodes))
	for i, n := range req.Cluster.Nodes {
		legacyNodes[i] = models.ClusterNode{
			ServerID:    n.ServerID,
			PrivateIP:   n.PrivateIP,
			Interface:   n.Interface,
			Fingerprint: n.Fingerprint,
		}
	}
	legacyReq := models.ClusterReviewRequest{
		Cluster: models.Cluster{
			ID:             req.Cluster.ID,
			ProjectID:      req.Cluster.ProjectID,
			OrganizationID: req.Cluster.OrganizationID,
			Name:           req.Cluster.Name,
			Version:        req.Cluster.Version,
			Controls:       req.Cluster.Controls,
			Nodes:          legacyNodes,
			Revision:       req.Cluster.Revision,
			Status:         req.Cluster.Status,
		},
		Action: req.Action,
	}
	return a.service.Review(ctx, userID, projectID, legacyReq)
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
	return a.service.ReviewStorageSetup(ctx, userID, projectID, clusterID, models.StorageSetupRequest{Version: req.Version})
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
