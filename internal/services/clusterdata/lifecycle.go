package clusterdata

import (
	"codedock/internal/engine/kubernetes"
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

type LifecycleRequest struct {
	Action      string                 `json:"action"`
	DatabaseID  string                 `json:"databaseId"`
	ConfirmName string                 `json:"confirmName"`
	Spec        models.ClusterDataSpec `json:"spec"`
}

func (s *Service) ReviewLifecycle(ctx context.Context, user, project, clusterID string, request LifecycleRequest) (*models.OperationReview, error) {
	if request.Action != "settings" && request.Action != "stop" && request.Action != "delete" {
		return nil, fmt.Errorf("unsupported database lifecycle action")
	}
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	plan, err := s.owned(ctx, project, clusterID, request.DatabaseID)
	if err != nil {
		return nil, err
	}
	if plan.Record.Status == "APPLYING" {
		return nil, fmt.Errorf("database operation already in progress")
	}
	target := "cluster-data:" + plan.Record.ID
	var effects string
	payloadPlan := *plan
	payloadPlan.Action = request.Action
	switch request.Action {
	case "stop":
		if plan.Record.Status == "STOPPED" {
			return nil, fmt.Errorf("database is already stopped")
		}
		effects = fmt.Sprintf("Stop cluster database %s by removing its workloads while retaining PVCs, secrets and credentials. Restart with a reviewed create or recover.", plan.Record.Spec.Name)
	case "delete":
		if request.ConfirmName != plan.Record.Spec.Name {
			return nil, fmt.Errorf("type the exact database name to confirm deletion")
		}
		effects = fmt.Sprintf("Delete cluster database %s and its owned workloads, auth secret and PVCs. Backups in object storage are retained. This cannot be undone.", plan.Record.Spec.Name)
	case "settings":
		if request.Spec.Instances != 1 && request.Spec.Instances != 3 {
			return nil, fmt.Errorf("select standalone or three database instances")
		}
		if request.Spec.StorageGiB < plan.Record.Spec.StorageGiB {
			return nil, fmt.Errorf("storage can only stay the same or increase")
		}
		if len(cluster.Nodes) < request.Spec.Instances {
			return nil, fmt.Errorf("replicated databases require three distinct cluster nodes")
		}
		next := plan.Record.Spec
		next.Instances = request.Spec.Instances
		next.StorageGiB = request.Spec.StorageGiB
		if request.Spec.Image != "" {
			next.Image = request.Spec.Image
		}
		if request.Spec.StorageClass != "" {
			next.StorageClass = request.Spec.StorageClass
		}
		if err := kubernetes.ValidateDataSpec(next, len(cluster.Nodes)); err != nil {
			return nil, err
		}
		payloadPlan.Record.Spec = next
		destination, err := s.destination(ctx, next.S3DestinationID)
		if err != nil {
			return nil, err
		}
		payloadPlan.Manifest, err = kubernetes.DataManifest(&payloadPlan, destination, nil, "")
		if err != nil {
			return nil, err
		}
		effects = fmt.Sprintf("Update cluster database %s to %d instance(s) and %d GiB per instance with a reviewed manifest. Storage never shrinks; workloads roll forward.", plan.Record.Spec.Name, next.Instances, next.StorageGiB)
	}
	payload, err := json.Marshal(payloadPlan)
	if err != nil {
		return nil, err
	}
	snapshot, err := s.snapshot(ctx, cluster, plan.Record.ID)
	if err != nil {
		return nil, err
	}
	return s.operations.Review(ctx, user, project, "cluster-data-lifecycle", target, string(payload), snapshot, effects)
}

func (s *Service) ApplyLifecycle(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "cluster-data-lifecycle" || op.UserID != user {
		return fmt.Errorf("database operation not found")
	}
	var plan models.ClusterDataPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return err
	}
	cluster, err := s.cluster(ctx, op.ProjectID, plan.Record.ClusterID)
	if err != nil {
		return err
	}
	snapshot, err := s.snapshot(ctx, cluster, plan.Record.ID)
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, snapshot, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		release, err := s.gate.AcquireVolume("cluster:" + cluster.ID)
		if err != nil {
			return err
		}
		defer release()
		latest, err := s.cluster(ctx, op.ProjectID, cluster.ID)
		if err != nil {
			return err
		}
		current, err := s.snapshot(ctx, latest, plan.Record.ID)
		if err != nil {
			return err
		}
		if current != op.Snapshot {
			return fmt.Errorf("database or cluster changed since review")
		}
		if err := s.engine.VerifyCluster(ctx, latest, true); err != nil {
			return err
		}
		node := latest.Nodes[0]
		stored, err := s.store.Get(ctx, plan.Record.ID)
		if err != nil {
			return err
		}
		switch plan.Action {
		case "stop":
			if err := s.store.Observe(ctx, plan.Record.ID, "APPLYING", ""); err != nil {
				return err
			}
			if err := s.engine.StopOwnedData(ctx, node, &stored.Record); err != nil {
				_ = s.store.Observe(ctx, plan.Record.ID, "FAILED", err.Error())
				return err
			}
			if err := progress("VERIFIED", "Database workloads stopped; PVCs and credentials retained"); err != nil {
				return err
			}
			return s.store.Observe(context.Background(), plan.Record.ID, "STOPPED", "")
		case "delete":
			if err := s.store.Observe(ctx, plan.Record.ID, "APPLYING", ""); err != nil {
				return err
			}
			if err := s.engine.RemoveOwnedData(ctx, node, &stored.Record); err != nil {
				_ = s.store.Observe(ctx, plan.Record.ID, "FAILED", err.Error())
				return err
			}
			if deleter, ok := s.store.(interface {
				Delete(context.Context, string) error
			}); ok {
				if err := deleter.Delete(ctx, plan.Record.ID); err != nil {
					return err
				}
			}
			return progress("VERIFIED", "Database workloads, secrets and PVCs removed; object-store backups retained")
		case "settings":
			if err := s.store.Observe(ctx, plan.Record.ID, "APPLYING", ""); err != nil {
				return err
			}
			if err := s.engine.ApplyOwnedData(ctx, node, plan.Manifest, progress); err != nil {
				_ = s.store.Observe(ctx, plan.Record.ID, "FAILED", err.Error())
				return err
			}
			if err := s.engine.DataReady(ctx, node, &plan.Record); err != nil {
				_ = s.store.Observe(ctx, plan.Record.ID, "FAILED", err.Error())
				return err
			}
			if err := progress("VERIFIED", "Database settings applied and instances ready"); err != nil {
				return err
			}
			return s.store.Observe(context.Background(), plan.Record.ID, "READY", "")
		default:
			return fmt.Errorf("unsupported lifecycle action")
		}
	})
}
