package clusters

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (s *Service) reviewUpgradeRecovery(ctx context.Context, user, project, id string) (*models.OperationReview, error) {
	store, ok := s.store.(UpgradeStore)
	if !ok {
		return nil, fmt.Errorf("upgrade journal storage unavailable")
	}
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project {
		return nil, fmt.Errorf("cluster project mismatch")
	}
	if err := s.validateNodes(ctx, user, cluster); err != nil {
		return nil, err
	}
	plans, err := store.PendingUpgrades(ctx)
	if err != nil {
		return nil, err
	}
	for _, plan := range plans {
		if plan.Cluster.ID != id {
			continue
		}
		plan.Action = "recover"
		payload, err := json.Marshal(plan)
		if err != nil {
			return nil, err
		}
		return s.operations.Review(ctx, user, project, "cluster", "cluster:"+id, string(payload), clusterSnapshot(cluster), "Resume the interrupted upgrade to its original reviewed release. Existing database and workload data remain in place. Each node must report the reviewed version and readiness before proceeding.")
	}
	return nil, fmt.Errorf("cluster has no interrupted upgrade")
}
