package projects

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
)

func (s *ProjectService) ListProjectMembers(ctx context.Context, projectID string) ([]*models.OrganizationMember, error) {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, errors.New("project not found")
	}
	return s.orgRepo.ListMembers(ctx, proj.OrganizationID)
}

func (s *ProjectService) AddProjectMember(ctx context.Context, projectID, email string, role models.MemberPermission) (*models.OrganizationMember, error) {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, errors.New("project not found")
	}

	now := time.Now().UTC()
	member := &models.OrganizationMember{
		ID:             uuid.NewString(),
		OrganizationID: proj.OrganizationID,
		UserID:         uuid.NewString(),
		Email:          email,
		Permission:     role,
		Status:         models.MemberStatusAccepted,
		InvitedAt:      now,
		AcceptedAt:     now,
	}
	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *ProjectService) RemoveProjectMember(ctx context.Context, projectID, memberID string) error {
	proj, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	if proj == nil {
		return errors.New("project not found")
	}
	return s.orgRepo.RemoveMember(ctx, memberID)
}
