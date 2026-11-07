package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"codedock/internal/models"
)

const (
	managedStreamPoll      = 2 * time.Second
	managedStreamHeartbeat = 25 * time.Second
)

func (s *managedService) authorizeManagedOperation(ctx context.Context, userID, operationID string) (*models.Operation, error) {
	op, err := s.operations.Get(ctx, operationID)
	if err != nil || op == nil {
		return nil, fmt.Errorf("managed operation not found")
	}
	if !strings.HasPrefix(op.Kind, "managed-") {
		return nil, fmt.Errorf("managed operation not found")
	}
	var plan models.ManagedServerPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return nil, fmt.Errorf("managed operation payload is invalid: %w", err)
	}
	if err := s.requireOrgRole(ctx, userID, plan.OrganizationID, false); err != nil {
		return nil, err
	}
	return op, nil
}

func (s *managedService) StreamOperation(ctx context.Context, userID, operationID string, emit func(*models.Operation) error) error {
	if emit == nil {
		return fmt.Errorf("operation stream requires an emitter")
	}
	current, err := s.authorizeManagedOperation(ctx, userID, operationID)
	if err != nil {
		return err
	}
	last := ""
	quietSince := time.Now()
	observe := func(op *models.Operation) error {
		fingerprint := snapshotValue(op.Status, op.Phase, op.Error, op.Logs)
		if fingerprint == last && time.Since(quietSince) < managedStreamHeartbeat && !models.ManagedOperationComplete(op.Status) {
			return nil
		}
		last = fingerprint
		quietSince = time.Now()
		return emit(op)
	}
	if err := observe(current); err != nil {
		return err
	}
	if models.ManagedOperationComplete(current.Status) {
		return nil
	}
	ticker := time.NewTicker(managedStreamPoll)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			op, err := s.operations.Get(ctx, operationID)
			if err != nil || op == nil {
				return fmt.Errorf("managed operation not found")
			}
			if err := observe(op); err != nil {
				return err
			}
			if models.ManagedOperationComplete(op.Status) {
				return nil
			}
		}
	}
}
