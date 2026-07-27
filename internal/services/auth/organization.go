package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type OrganizationService struct {
	orgRepo  repositories.OrganizationRepository
	userRepo repositories.UserRepository
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository, userRepo repositories.UserRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo, userRepo: userRepo}
}

func (s *OrganizationService) isRequesterOwnerOrAdmin(ctx context.Context, orgID, requesterUserID string) bool {
	if requesterUserID == "" {
		return false
	}
	requester, err := s.orgRepo.GetMember(ctx, orgID, requesterUserID)
	if err == nil && requester != nil {
		return requester.Permission == models.MemberPermissionOwner
	}
	if s.userRepo != nil {
		u, err := s.userRepo.GetUserByID(ctx, requesterUserID)
		if err == nil && u != nil && u.Role == "admin" {
			return true
		}
	}
	return false
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

func (s *OrganizationService) InviteMember(ctx context.Context, requesterUserID, orgID, email string, permission models.MemberPermission) (*models.OrganizationMember, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if permission == models.MemberPermissionOwner {
		if !s.isRequesterOwnerOrAdmin(ctx, orgID, requesterUserID) {
			return nil, errors.New("only organization owners can invite new owners")
		}
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
	if memberID == "" {
		return errors.New("member id is required")
	}
	orgMembers, err := s.orgRepo.ListMembers(ctx, "")
	var target *models.OrganizationMember
	if err == nil {
		for _, m := range orgMembers {
			if m.ID == memberID {
				target = m
				break
			}
		}
	}
	if target != nil && target.Permission == models.MemberPermissionOwner {
		membersInOrg, err := s.orgRepo.ListMembers(ctx, target.OrganizationID)
		if err == nil {
			ownerCount := 0
			for _, m := range membersInOrg {
				if m.Permission == models.MemberPermissionOwner {
					ownerCount++
				}
			}
			if ownerCount <= 1 {
				return errors.New("cannot remove the last owner of an organization")
			}
		}
	}
	return s.orgRepo.RemoveMember(ctx, memberID)
}

func (s *OrganizationService) UpdateMemberPermission(ctx context.Context, requesterUserID, orgID, targetUserID string, permission models.MemberPermission) error {
	isOwnerOrAdmin := s.isRequesterOwnerOrAdmin(ctx, orgID, requesterUserID)

	targetMember, err := s.orgRepo.GetMember(ctx, orgID, targetUserID)
	if err != nil || targetMember == nil {
		return errors.New("target member not found")
	}

	if (permission == models.MemberPermissionOwner || targetMember.Permission == models.MemberPermissionOwner) && !isOwnerOrAdmin {
		return errors.New("only organization owners can manage owner permissions")
	}

	if targetMember.Permission == models.MemberPermissionOwner && permission != models.MemberPermissionOwner {
		membersInOrg, err := s.orgRepo.ListMembers(ctx, orgID)
		if err == nil {
			ownerCount := 0
			for _, m := range membersInOrg {
				if m.Permission == models.MemberPermissionOwner {
					ownerCount++
				}
			}
			if ownerCount <= 1 {
				return errors.New("cannot demote the last owner of an organization")
			}
		}
	}

	targetMember.Permission = permission
	return s.orgRepo.UpdateMember(ctx, targetMember)
}
