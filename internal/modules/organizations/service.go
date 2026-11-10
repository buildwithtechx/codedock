package organizations

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/utils"
)

type OrganizationService struct {
	orgRepo OrganizationRepository
}

func NewOrganizationService(orgRepo OrganizationRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo}
}

func (s *OrganizationService) isRequesterOwnerOrAdmin(ctx context.Context, orgID, requesterUserID string) bool {
	if requesterUserID == "" {
		return false
	}
	requester, err := s.orgRepo.GetMember(ctx, orgID, requesterUserID)
	if err == nil && requester != nil && requester.Permission == MemberPermissionOwner {
		return true
	}
	role, err := s.orgRepo.GetUserRole(ctx, requesterUserID)
	if err == nil && (role == "admin" || role == "owner") {
		return true
	}
	return false
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, userID, name string) (*Organization, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = utils.GenerateRandomName()
	}
	org := &Organization{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	owner := &OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		UserID:         userID,
		Permission:     MemberPermissionOwner,
		Status:         MemberStatusAccepted,
		InvitedAt:      time.Now(),
		AcceptedAt:     time.Now(),
	}
	if err := s.orgRepo.CreateWithOwner(ctx, org, owner); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) ListOrganizationsByUser(ctx context.Context, userID string) ([]*Organization, error) {
	return s.orgRepo.ListByUser(ctx, userID)
}

func (s *OrganizationService) GetOrCreateDefaultOrganization(ctx context.Context, userID string) (*Organization, error) {
	orgs, err := s.ListOrganizationsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(orgs) > 0 {
		return orgs[0], nil
	}
	return s.CreateOrganization(ctx, userID, "Default")
}

func (s *OrganizationService) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, requesterUserID, id, name string) (*Organization, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if !s.isRequesterOwnerOrAdmin(ctx, id, requesterUserID) {
		return nil, ErrInsufficientPermissions
	}
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	org.Name = name
	org.UpdatedAt = time.Now()
	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	return s.orgRepo.Delete(ctx, id)
}

func (s *OrganizationService) InviteMember(ctx context.Context, requesterUserID, orgID, email string, permission MemberPermission) (*OrganizationMember, error) {
	if email == "" {
		return nil, ErrEmailRequired
	}
	if permission == MemberPermissionOwner {
		if !s.isRequesterOwnerOrAdmin(ctx, orgID, requesterUserID) {
			return nil, ErrOwnerInviteOnly
		}
	}
	existing, _ := s.orgRepo.GetMemberByEmail(ctx, orgID, email)
	if existing != nil {
		return nil, ErrAlreadyInvitedOrMember
	}

	member := &OrganizationMember{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          email,
		Permission:     permission,
		Status:         MemberStatusPending,
		InvitedAt:      time.Now(),
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *OrganizationService) ListMembers(ctx context.Context, orgID string) ([]*OrganizationMember, error) {
	return s.orgRepo.ListMembers(ctx, orgID)
}

func (s *OrganizationService) RemoveMember(ctx context.Context, memberID string) error {
	if memberID == "" {
		return ErrMemberIDRequired
	}
	target, err := s.orgRepo.GetMemberByID(ctx, memberID)
	if err != nil {
		return err
	}
	if target != nil && target.Permission == MemberPermissionOwner {
		membersInOrg, err := s.orgRepo.ListMembers(ctx, target.OrganizationID)
		if err == nil {
			activeOwnerCount := 0
			for _, m := range membersInOrg {
				if m.Permission == MemberPermissionOwner && m.UserID != "" && m.Status == MemberStatusAccepted {
					activeOwnerCount++
				}
			}
			if activeOwnerCount <= 1 {
				return ErrCannotRemoveLastOwner
			}
		}
	}
	return s.orgRepo.RemoveMember(ctx, memberID)
}

func (s *OrganizationService) UpdateMemberPermission(ctx context.Context, requesterUserID, orgID, targetUserID string, permission MemberPermission) error {
	isOwnerOrAdmin := s.isRequesterOwnerOrAdmin(ctx, orgID, requesterUserID)

	targetMember, err := s.orgRepo.GetMember(ctx, orgID, targetUserID)
	if err != nil || targetMember == nil {
		return ErrTargetMemberNotFound
	}

	if (permission == MemberPermissionOwner || targetMember.Permission == MemberPermissionOwner) && !isOwnerOrAdmin {
		return ErrOwnerManageOnly
	}

	if targetMember.Permission == MemberPermissionOwner && permission != MemberPermissionOwner {
		membersInOrg, err := s.orgRepo.ListMembers(ctx, orgID)
		if err == nil {
			activeOwnerCount := 0
			for _, m := range membersInOrg {
				if m.Permission == MemberPermissionOwner && m.UserID != "" && m.Status == MemberStatusAccepted {
					activeOwnerCount++
				}
			}
			if activeOwnerCount <= 1 {
				return ErrCannotDemoteLastOwner
			}
		}
	}

	targetMember.Permission = permission
	return s.orgRepo.UpdateMember(ctx, targetMember)
}
