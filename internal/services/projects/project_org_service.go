package projects

import (
	"context"
	"errors"

	"codedock.run/codedock/internal/models"
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
