package projects

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"codedock/internal/models"
)

func (s *ProjectService) ListOrganizationsByUser(ctx context.Context, userID string) ([]*models.Organization, error) {
	if s.orgRepo == nil {
		return nil, nil
	}
	return s.orgRepo.ListByUser(ctx, userID)
}

func (s *ProjectService) CreateOrganizationWithOwner(ctx context.Context, org *models.Organization, owner *models.OrganizationMember) error {
	if s.orgRepo == nil {
		return errors.New("organization repository is not initialized")
	}
	return s.orgRepo.CreateWithOwner(ctx, org, owner)
}

func (s *ProjectService) GetOrCreateDefaultOrganization(ctx context.Context, userID string) (*models.Organization, error) {
	s.defaultOrganizationMu.Lock()
	defer s.defaultOrganizationMu.Unlock()
	organizations, err := s.ListOrganizationsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user organizations: %w", err)
	}
	if len(organizations) > 0 {
		return organizations[0], nil
	}
	now := time.Now().UTC()
	org := &models.Organization{ID: uuid.NewString(), Name: "Default Workspace", CreatedAt: now, UpdatedAt: now}
	member := &models.OrganizationMember{ID: uuid.NewString(), OrganizationID: org.ID, UserID: userID, Permission: models.MemberPermissionOwner, Status: models.MemberStatusAccepted, InvitedAt: now, AcceptedAt: now}
	if err := s.CreateOrganizationWithOwner(ctx, org, member); err != nil {
		return nil, fmt.Errorf("provision default organization: %w", err)
	}
	return org, nil
}
