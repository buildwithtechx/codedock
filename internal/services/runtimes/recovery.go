package runtimes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"errors"
	"fmt"
	"time"
)

func (s *Service) recoverFailure(app *models.AppService, node models.ClusterNode, journal string, cause error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	rollback := s.engine.Rollback(ctx, node, app, journal)
	status := "FAILED"
	if rollback != nil {
		status = "RECOVERY_REQUIRED"
	}
	persist := s.store.Observe(ctx, app.ID, status, errors.Join(cause, rollback).Error(), rollback == nil)
	return errors.Join(cause, rollback, persist)
}
func (s *Service) Recover(ctx context.Context) error {
	pending, err := s.store.Pending(ctx)
	if err != nil {
		return err
	}
	for _, id := range pending {
		runtime, err := s.store.Get(ctx, id)
		if err != nil {
			return err
		}
		app, err := s.apps.GetByID(ctx, id)
		if err != nil {
			return err
		}
		cluster, err := s.clusters.Get(ctx, runtime.Target.ClusterID)
		if err != nil || len(cluster.Nodes) == 0 {
			if err := s.store.Observe(ctx, id, "RECOVERY_REQUIRED", "Cluster unavailable during recovery", false); err != nil {
				return err
			}
			continue
		}
		release, err := s.lockCluster(cluster)
		if err != nil {
			return err
		}
		recovery, cancel := context.WithTimeout(ctx, 6*time.Minute)
		err = s.engine.VerifyCluster(recovery, cluster, false)
		if err == nil {
			err = s.engine.Rollback(recovery, cluster.Nodes[0], app, runtime.Journal)
		}
		cancel()
		release()
		message := "Deployment interrupted; previous resources restored"
		if err != nil {
			message = fmt.Sprintf("Deployment recovery failed: %v", err)
		}
		status := "INTERRUPTED"
		if err != nil {
			status = "RECOVERY_REQUIRED"
		}
		if err := s.store.Observe(ctx, id, status, message, err == nil); err != nil {
			return err
		}
	}
	return nil
}
