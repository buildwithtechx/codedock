package clusters

import (
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type UpgradeStore interface {
	BeginUpgrade(context.Context, *models.ClusterPlan) error
	PendingUpgrades(context.Context) ([]models.ClusterPlan, error)
	FinishUpgrade(context.Context, string, string, string, string) error
}

func (s *Service) RecoverUpgrades(ctx context.Context) error {
	store, ok := s.store.(UpgradeStore)
	if !ok {
		return fmt.Errorf("cluster upgrade recovery storage unavailable")
	}
	plans, err := store.PendingUpgrades(ctx)
	if err != nil {
		return err
	}
	for _, plan := range plans {
		if plan.OperationID != "" {
			op, err := s.operations.Get(ctx, plan.OperationID)
			if err != nil {
				return err
			}
			if op.Status == "CANCELLED" || op.Status == "CANCELLING" {
				continue
			}
		}
		releases, err := s.acquireTargets(&plan.Cluster)
		if err != nil {
			return err
		}
		recovery, cancel := context.WithTimeout(ctx, 20*time.Minute)
		err = s.resumeUpgrade(recovery, &plan)
		cancel()
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
		if err != nil {
			if persist := s.store.Observe(ctx, plan.Cluster.ID, "RECOVERY_REQUIRED", err.Error()); persist != nil {
				return errors.Join(err, persist)
			}
			continue
		}
		if err := store.FinishUpgrade(ctx, plan.Cluster.ID, plan.Cluster.Version, "READY", "Interrupted upgrade resumed to the reviewed release"); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) resumeUpgrade(ctx context.Context, plan *models.ClusterPlan) error {
	digest := sha256.Sum256([]byte(plan.Installer))
	if hex.EncodeToString(digest[:]) != plan.InstallerSHA256 {
		return fmt.Errorf("upgrade installer checksum mismatch")
	}
	latest, err := s.store.Get(ctx, plan.Cluster.ID)
	if err != nil {
		return err
	}
	if latest.ProjectID != plan.Cluster.ProjectID || latest.OrganizationID != plan.Cluster.OrganizationID || !reflect.DeepEqual(latest.Nodes, plan.Cluster.Nodes) {
		return fmt.Errorf("upgrade target ownership changed")
	}
	if err := s.validateNodes(ctx, "", latest); err != nil {
		return err
	}
	cluster := plan.Cluster
	cluster.JoinToken = latest.JoinToken
	for i, node := range cluster.Nodes {
		if err := s.runner.Preflight(ctx, node, cluster.ID, i == 0); err != nil {
			return err
		}
		role := "agent"
		if i == 0 {
			role = "leader"
		} else if i < cluster.Controls {
			role = "control"
		}
		if err := s.runner.Script(ctx, node, "recover-"+cluster.ID, installationScript(&cluster, node, role, plan.Installer)); err != nil {
			return fmt.Errorf("resume node %s release: %w", node.ServerID, err)
		}
		if err := s.waitRelease(ctx, &cluster, node); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) finishUpgrade(plan *models.ClusterPlan, cause error) error {
	store, ok := s.store.(UpgradeStore)
	if !ok {
		return errors.Join(cause, fmt.Errorf("upgrade journal storage unavailable"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	if errors.Is(cause, context.Canceled) || errors.Is(cause, context.DeadlineExceeded) {
		return cause
	}
	if cause == nil {
		return store.FinishUpgrade(ctx, plan.Cluster.ID, plan.Cluster.Version, "READY", "")
	}
	if err := s.resumeUpgrade(ctx, plan); err != nil {
		return errors.Join(cause, fmt.Errorf("cluster upgrade recovery required: %w", err))
	}
	if err := store.FinishUpgrade(ctx, plan.Cluster.ID, plan.Cluster.Version, "READY", "Interrupted upgrade resumed to the reviewed release"); err != nil {
		return errors.Join(cause, err)
	}
	return nil
}
