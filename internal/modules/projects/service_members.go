package projects

import (
	"context"
	"time"

	"github.com/google/uuid"

	"codedock/internal/modules/organizations"
)

func (s *Service) ListProjectMembers(ctx context.Context, projectID string) ([]*organizations.OrganizationMember, error) {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, ErrProjectNotFound
	}
	return s.orgRepo.ListMembers(ctx, proj.OrganizationID)
}

func (s *Service) AddProjectMember(ctx context.Context, projectID, email string, role organizations.MemberPermission) (*organizations.OrganizationMember, error) {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, ErrProjectNotFound
	}

	now := time.Now().UTC()
	member := &organizations.OrganizationMember{
		ID:             uuid.NewString(),
		OrganizationID: proj.OrganizationID,
		UserID:         uuid.NewString(),
		Email:          email,
		Permission:     role,
		Status:         organizations.MemberStatusAccepted,
		InvitedAt:      now,
		AcceptedAt:     now,
	}
	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *Service) RemoveProjectMember(ctx context.Context, projectID, memberID string) error {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	if proj == nil {
		return ErrProjectNotFound
	}
	return s.orgRepo.RemoveMember(ctx, memberID)
}
