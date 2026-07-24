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

func (s *OrganizationService) InviteMember(ctx context.Context, orgID, email string, permission models.MemberPermission) (*models.OrganizationMember, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	existing, _ := s.orgRepo.GetMemberByEmail(ctx, orgID, email)
	if existing != nil {
		return nil, errors.New("user already invited or is a member")
	}

	member := &models.OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          email,
		Permission:     permission,
		Status:         models.MemberStatusPending,
		InvitedAt:      time.Now(),
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *OrganizationService) ListMembers(ctx context.Context, orgID string) ([]*models.OrganizationMember, error) {
	return s.orgRepo.ListMembers(ctx, orgID)
}

func (s *OrganizationService) RemoveMember(ctx context.Context, memberID string) error {
	return s.orgRepo.RemoveMember(ctx, memberID)
}

func (s *OrganizationService) UpdateMemberPermission(ctx context.Context, orgID, userID string, permission models.MemberPermission) error {
	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil || member == nil {
		return errors.New("member not found")
	}
	if member.Permission == models.MemberPermissionOwner && permission != models.MemberPermissionOwner {
		return errors.New("cannot change permission of owner")
	}
	member.Permission = permission
	return s.orgRepo.UpdateMember(ctx, member)
}
