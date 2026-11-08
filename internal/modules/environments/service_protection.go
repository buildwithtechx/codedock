package environments

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ProjectLookup interface {
	Get(ctx context.Context, id string) (any, error)
}

type RouteProtectionService struct {
	repo        RouteProtectionRepository
	projectRepo ProjectLookup
}

func NewRouteProtectionService(repo RouteProtectionRepository, projRepo ProjectLookup) *RouteProtectionService {
	return &RouteProtectionService{
		repo:        repo,
		projectRepo: projRepo,
	}
}

func (s *RouteProtectionService) List(ctx context.Context, projectID string) ([]RouteProtectionStatus, error) {
	rows, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]RouteProtectionStatus, 0, len(rows))
	for _, row := range rows {
		st := RouteProtectionStatusPending
		if (!row.Enabled || row.Secret != nil) && row.AppliedRevision != nil && *row.AppliedRevision == row.Revision {
			st = RouteProtectionStatusApplied
		}
		result = append(result, RouteProtectionStatus{
			Hostname: row.Hostname,
			Enabled:  row.Enabled,
			Status:   st,
		})
	}
	return result, nil
}

func (s *RouteProtectionService) GetPolicyByHostname(ctx context.Context, hostname string) (*RouteProtection, *RouteProtectionPolicy, error) {
	clean, err := NormalizeHostname(hostname)
	if err != nil {
		return nil, nil, err
	}
	rp, err := s.repo.GetByHostname(ctx, clean)
	if err != nil || rp == nil || !rp.Enabled || rp.Secret == nil {
		return nil, nil, nil
	}
	var policy RouteProtectionPolicy
	if err := json.Unmarshal([]byte(*rp.Secret), &policy); err != nil {
		return nil, nil, err
	}
	return rp, &policy, nil
}

func (s *RouteProtectionService) Change(ctx context.Context, projectID string, input RouteProtectionInput) (*RouteProtectionStatus, error) {
	clean, err := NormalizeHostname(input.Hostname)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.Get(ctx, projectID, clean)
	if err != nil {
		return nil, err
	}

	switch input.Action {
	case ActionProtect:
		policy, err := CreatePolicy(input.Password, DefaultSessionSec)
		if err != nil {
			return nil, err
		}
		raw, _ := json.Marshal(policy)
		secretStr := string(raw)
		rev := fmt.Sprintf("%x", sha256.Sum256([]byte(secretStr+clean)))
		rp := &RouteProtection{
			ProjectID:       projectID,
			Hostname:        clean,
			Enabled:         true,
			Secret:          &secretStr,
			Revision:        rev,
			AppliedRevision: &rev,
		}
		if err := s.repo.Save(ctx, rp); err != nil {
			return nil, err
		}
		return &RouteProtectionStatus{Hostname: clean, Enabled: true, Status: RouteProtectionStatusApplied}, nil

	case ActionRevoke:
		if existing == nil || existing.Secret == nil {
			return nil, errors.New("no active route protection to revoke")
		}
		var policy RouteProtectionPolicy
		if err := json.Unmarshal([]byte(*existing.Secret), &policy); err != nil {
			return nil, err
		}
		rotated := RevokeRouteSessions(&policy)
		raw, _ := json.Marshal(rotated)
		secretStr := string(raw)
		rev := fmt.Sprintf("%x", sha256.Sum256([]byte(secretStr+clean)))
		existing.Secret = &secretStr
		existing.Revision = rev
		existing.AppliedRevision = &rev
		if err := s.repo.Save(ctx, existing); err != nil {
			return nil, err
		}
		return &RouteProtectionStatus{Hostname: clean, Enabled: true, Status: RouteProtectionStatusApplied}, nil

	case ActionRetry:
		if existing == nil {
			return nil, ErrProtectionNotFound
		}
		_ = s.repo.MarkApplied(ctx, projectID, clean, existing.Revision)
		return &RouteProtectionStatus{Hostname: clean, Enabled: existing.Enabled, Status: RouteProtectionStatusApplied}, nil

	case ActionUnlock:
		if existing != nil {
			existing.Enabled = false
			existing.Revision = fmt.Sprintf("unlocked-%d", time.Now().Unix())
			existing.AppliedRevision = &existing.Revision
			if err := s.repo.Save(ctx, existing); err != nil {
				return nil, err
			}
		}
		return &RouteProtectionStatus{Hostname: clean, Enabled: false, Status: RouteProtectionStatusApplied}, nil

	default:
		return nil, fmt.Errorf("unsupported route protection action: %s", input.Action)
	}
}
