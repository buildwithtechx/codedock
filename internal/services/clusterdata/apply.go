package clusterdata

import (
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

func (s *Service) Apply(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "cluster-data" || op.UserID != user {
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
	return s.operations.Apply(ctx, id, user, confirmation, snapshot, func(ctx context.Context, op *models.Operation, progress func(string, string) error) (runErr error) {
		defer func() { runErr = redactDataError(&plan, runErr) }()
		names := []string{"cluster:" + cluster.ID}
		for _, node := range cluster.Nodes {
			names = append(names, "server:"+node.ServerID)
		}
		sort.Strings(names)
		releases := []func(){}
		defer func() {
			for i := len(releases) - 1; i >= 0; i-- {
				releases[i]()
			}
		}()
		for _, name := range names {
			release, err := s.gate.AcquireVolume(name)
			if err != nil {
				return err
			}
			releases = append(releases, release)
		}
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
		if plan.Action == "operators" {
			for _, manifest := range plan.Operators {
				sum := sha256.Sum256([]byte(manifest.Manifest))
				if hex.EncodeToString(sum[:]) != manifest.SHA256 {
					return fmt.Errorf("operator checksum mismatch")
				}
				if err := s.engine.ApplyOwnedData(ctx, node, manifest.Manifest, progress); err != nil {
					return err
				}
				namespace := "cnpg-system"
				deployments := []string{"cnpg-controller-manager"}
				if manifest.Name == "cert-manager" {
					namespace = "cert-manager"
					deployments = []string{"cert-manager", "cert-manager-webhook", "cert-manager-cainjector"}
				}
				if manifest.Name == "Barman Cloud" {
					deployments = []string{"barman-cloud"}
				}
				for _, deployment := range deployments {
					if _, err := s.commands.Kubectl(ctx, node, []string{"-n", namespace, "rollout", "status", "deployment/" + deployment, "--timeout=300s"}, ""); err != nil {
						return err
					}
				}
			}
			return progress("VERIFIED", "Database operator dependencies are ready")
		}
		if plan.Action != "backup" {
			defer func() {
				final, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				status, message := "READY", ""
				if runErr != nil {
					status, message = "FAILED", redactDataError(&plan, runErr).Error()
					if ctx.Err() != nil {
						status = "INTERRUPTED"
					}
				}
				if err := s.store.Observe(final, plan.Record.ID, status, message); err != nil {
					slog.Error("save cluster database outcome", "database", plan.Record.ID, "error", err)
					runErr = err
				}
			}()
		}
		if plan.Action != "backup" {
			if err := s.store.Observe(ctx, plan.Record.ID, "APPLYING", ""); err != nil {
				return err
			}
		}
		if err := s.engine.ApplyOwnedData(ctx, node, plan.Manifest, progress); err != nil {
			return err
		}
		if plan.Action == "backup" {
			if err := s.engine.WaitDataBackup(ctx, node, &plan.Record, plan.Manifest); err != nil {
				return err
			}
			return progress("VERIFIED", "Object-store base backup completed")
		}
		if err := progress("READINESS", "Waiting for all database instances and replication readiness"); err != nil {
			return err
		}
		if err := s.engine.DataReady(ctx, node, &plan.Record); err != nil {
			return err
		}
		return progress("VERIFIED", "Database instances are ready; persistent storage retained")
	})
}
