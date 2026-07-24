package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type OrganizationService struct {
	orgRepo repositories.OrganizationRepository
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo}
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, userID, name string) (*models.Organization, error) {
	if name == "" {
		return nil, errors.New("organization name is required")
	}
	org := &models.Organization{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	owner := &models.OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		UserID:         userID,
		Permission:     models.MemberPermissionOwner,
		Status:         models.MemberStatusAccepted,
		InvitedAt:      time.Now(),
		AcceptedAt:     time.Now(),
	}
	if err := s.orgRepo.AddMember(ctx, owner); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) ListOrganizationsByUser(ctx context.Context, userID string) ([]*models.Organization, error) {
	return s.orgRepo.ListByUser(ctx, userID)
}

func (s *OrganizationService) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	return s.orgRepo.Delete(ctx, id)
}
