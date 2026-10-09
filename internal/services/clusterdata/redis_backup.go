package clusterdata

import (
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type RedisSnapshotStore interface {
	Create(context.Context, *models.RedisSnapshot) error
	List(context.Context, string) ([]models.RedisSnapshot, error)
	Get(context.Context, string) (*models.RedisSnapshot, error)
}

func (s *Service) redisStore() (RedisSnapshotStore, error) {
	if s.snapshots == nil {
		return nil, fmt.Errorf("snapshot storage unavailable")
	}
	return s.snapshots, nil
}

func (s *Service) ReviewRedisSnapshot(ctx context.Context, user, project, clusterID, databaseID string) (*models.OperationReview, error) {
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	plan, err := s.owned(ctx, project, clusterID, databaseID)
	if err != nil {
		return nil, err
	}
	if plan.Record.Spec.Engine != "redis" {
		return nil, fmt.Errorf("snapshots support Redis databases")
	}
	if plan.Record.Status != "READY" {
		return nil, fmt.Errorf("database is not ready")
	}
	if plan.Record.Spec.S3DestinationID == "" {
		_ = cluster
		return nil, fmt.Errorf("attach an object-store destination before snapshotting")
	}
	payload, err := json.Marshal(map[string]string{"databaseId": databaseID, "clusterId": clusterID, "projectId": project})
	if err != nil {
		return nil, err
	}
	snapshot, err := s.snapshot(ctx, cluster, plan.Record.ID)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("Capture a Redis RDB snapshot of %s with BGSAVE and upload it to object storage. Source keeps serving; replicas re-sync.", plan.Record.Spec.Name)
	return s.operations.Review(ctx, user, project, "redis-snapshot", "cluster-data:"+plan.Record.ID, string(payload), snapshot, effects)
}

func (s *Service) ApplyRedisSnapshot(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "redis-snapshot" || op.UserID != user {
		return fmt.Errorf("snapshot operation not found")
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(op.Payload), &payload); err != nil {
		return err
	}
	cluster, err := s.cluster(ctx, op.ProjectID, payload["clusterId"])
	if err != nil {
		return err
	}
	snapshot, err := s.snapshot(ctx, cluster, payload["databaseId"])
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, snapshot, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		plan, err := s.store.Get(ctx, payload["databaseId"])
		if err != nil {
			return err
		}
		dest, err := s.destination(ctx, plan.Record.Spec.S3DestinationID)
		if err != nil || dest == nil {
			return fmt.Errorf("snapshot destination unavailable")
		}
		key := fmt.Sprintf("codedock/%s/redis/%s/%s.rdb", plan.Record.ProjectID, plan.Record.ID, time.Now().UTC().Format("20060102T150405Z"))
		if err := progress("SNAPSHOT", "Capturing Redis RDB snapshot"); err != nil {
			return err
		}
		size, err := s.engine.RedisSnapshot(ctx, cluster.Nodes[0], &plan.Record, plan.Password, dest, key)
		if err != nil {
			return err
		}
		store, err := s.redisStore()
		if err != nil {
			return err
		}
		if err := store.Create(ctx, &models.RedisSnapshot{DatabaseID: plan.Record.ID, ProjectID: plan.Record.ProjectID, ClusterID: plan.Record.ClusterID, S3DestinationID: plan.Record.Spec.S3DestinationID, S3Key: key, SizeBytes: size, Status: "COMPLETED"}); err != nil {
			return err
		}
		return progress("VERIFIED", fmt.Sprintf("Redis snapshot uploaded (%d bytes)", size))
	})
}

func (s *Service) ListRedisSnapshots(ctx context.Context, project, clusterID, databaseID string) ([]models.RedisSnapshot, error) {
	if _, err := s.owned(ctx, project, clusterID, databaseID); err != nil {
		return nil, err
	}
	store, err := s.redisStore()
	if err != nil {
		return nil, err
	}
	return store.List(ctx, databaseID)
}

func (s *Service) ReviewRedisRestore(ctx context.Context, user, project, clusterID, databaseID, snapshotID, targetID string) (*models.OperationReview, error) {
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	_, err = s.owned(ctx, project, clusterID, databaseID)
	if err != nil {
		return nil, err
	}
	target, err := s.owned(ctx, project, clusterID, targetID)
	if err != nil {
		return nil, err
	}
	if target.Record.Spec.Engine != "redis" {
		return nil, fmt.Errorf("restore supports Redis targets")
	}
	if databaseID == targetID {
		return nil, fmt.Errorf("restore requires a separate target database")
	}
	store, err := s.redisStore()
	if err != nil {
		return nil, err
	}
	snapshot, err := store.Get(ctx, snapshotID)
	if err != nil {
		return nil, err
	}
	if snapshot.DatabaseID != databaseID {
		return nil, fmt.Errorf("snapshot belongs to another database")
	}
	payload, err := json.Marshal(map[string]string{"snapshotId": snapshotID, "sourceId": databaseID, "targetId": targetID, "clusterId": clusterID, "projectId": project})
	if err != nil {
		return nil, err
	}
	state, err := s.snapshot(ctx, cluster, targetID)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("Restore Redis snapshot into separate target %s by replacing its RDB and restarting its pods. Source %s is preserved; target data is replaced.", target.Record.Spec.Name, databaseID)
	sum := sha256.Sum256([]byte(snapshot.S3Key + targetID))
	_ = sum
	return s.operations.Review(ctx, user, project, "redis-restore", "cluster-data:"+targetID, string(payload), state, effects)
}

func (s *Service) ApplyRedisRestore(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "redis-restore" || op.UserID != user {
		return fmt.Errorf("restore operation not found")
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(op.Payload), &payload); err != nil {
		return err
	}
	cluster, err := s.cluster(ctx, op.ProjectID, payload["clusterId"])
	if err != nil {
		return err
	}
	state, err := s.snapshot(ctx, cluster, payload["targetId"])
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, state, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		store, err := s.redisStore()
		if err != nil {
			return err
		}
		snapshot, err := store.Get(ctx, payload["snapshotId"])
		if err != nil {
			return err
		}
		target, err := s.store.Get(ctx, payload["targetId"])
		if err != nil {
			return err
		}
		source, err := s.store.Get(ctx, payload["sourceId"])
		if err != nil {
			return err
		}
		dest, err := s.destination(ctx, source.Record.Spec.S3DestinationID)
		if err != nil || dest == nil {
			return fmt.Errorf("snapshot destination unavailable")
		}
		if err := progress("RESTORE", "Replacing target RDB and restarting pods"); err != nil {
			return err
		}
		if err := s.engine.RedisRestore(ctx, cluster.Nodes[0], &target.Record, target.Password, dest, snapshot.S3Key); err != nil {
			return err
		}
		return progress("VERIFIED", "Target restored from snapshot; source preserved")
	})
}

func (s *Service) ReviewRedisFailover(ctx context.Context, user, project, clusterID, databaseID string) (*models.OperationReview, error) {
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	plan, err := s.owned(ctx, project, clusterID, databaseID)
	if err != nil {
		return nil, err
	}
	if plan.Record.Spec.Engine != "redis" || plan.Record.Spec.Instances < 2 {
		return nil, fmt.Errorf("failover requires a replicated Redis database")
	}
	if plan.Record.Status != "READY" {
		return nil, fmt.Errorf("database is not ready")
	}
	payload, err := json.Marshal(map[string]string{"databaseId": databaseID, "clusterId": clusterID, "projectId": project})
	if err != nil {
		return nil, err
	}
	state, err := s.snapshot(ctx, cluster, databaseID)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("Fail over Redis %s by restarting its fixed primary; replicas re-sync from AOF. This is restart recovery, not Cluster sharding.", plan.Record.Spec.Name)
	return s.operations.Review(ctx, user, project, "redis-failover", "cluster-data:"+databaseID, string(payload), state, effects)
}

func (s *Service) ApplyRedisFailover(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "redis-failover" || op.UserID != user {
		return fmt.Errorf("failover operation not found")
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(op.Payload), &payload); err != nil {
		return err
	}
	cluster, err := s.cluster(ctx, op.ProjectID, payload["clusterId"])
	if err != nil {
		return err
	}
	state, err := s.snapshot(ctx, cluster, payload["databaseId"])
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, state, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		plan, err := s.store.Get(ctx, payload["databaseId"])
		if err != nil {
			return err
		}
		if err := progress("FAILOVER", "Restarting fixed primary"); err != nil {
			return err
		}
		if err := s.engine.RedisFailover(ctx, cluster.Nodes[0], &plan.Record); err != nil {
			return err
		}
		return progress("VERIFIED", "Primary restarted and replicas observed ready")
	})
}
