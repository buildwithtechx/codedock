package projects

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock/internal/modules/organizations"
	"codedock/internal/modules/users"
	"codedock/internal/utils"
)

type OrganizationRepository interface {
	GetMember(ctx context.Context, organizationID, userID string) (*organizations.OrganizationMember, error)
}

type UserOrgLookup interface {
	ListOrganizationsByUser(ctx context.Context, userID string) ([]*organizations.Organization, error)
}

type Service struct {
	projectRepo ProjectRepository
	orgRepo     OrganizationRepository
}

func NewService(pr ProjectRepository, orgRepo OrganizationRepository) *Service {
	return &Service{
		projectRepo: pr,
		orgRepo:     orgRepo,
	}
}

func (s *Service) CreateProject(ctx context.Context, name, description string) (*ProjectConfig, error) {
	if name == "" {
		name = utils.GenerateRandomName()
	}
	id := uuid.NewString()
	now := time.Now().UTC()
	p := &ProjectConfig{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) CreateProjectFromRequest(ctx context.Context, req *CreateProjectRequest) (*ProjectConfig, error) {
	if req.Name == "" {
		req.Name = utils.GenerateRandomName()
	}
	id := req.ID
	if id == "" {
		id = uuid.NewString()
	}
	orgID := req.OrganizationID
	if orgID == "none" {
		orgID = ""
	}
	serverID := req.ServerID
	if serverID == "local" {
		serverID = ""
	}
	p := &ProjectConfig{
		ID:              id,
		AppID:           req.AppID,
		OrganizationID:  orgID,
		ServerID:        serverID,
		Name:            req.Name,
		Description:     req.Description,
		EnvironmentName: req.EnvironmentName,
		EnvironmentType: req.EnvironmentType,
		GitProvider:     req.GitProvider,
		GitOwner:        req.GitOwner,
		GitRepo:         req.GitRepo,
		GitBranch:       req.GitBranch,
		GitURL:          req.GitURL,
		IsApp:           req.IsApp,
		AppTemplateID:   req.AppTemplateID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return p, nil
}

func (s *Service) GetProject(ctx context.Context, id string) (*ProjectConfig, error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}
	return s.projectRepo.Get(ctx, id)
}

func (s *Service) HasPermission(ctx context.Context, projectID, userID string, userRole users.UserRole, minPermission organizations.MemberPermission) bool {
	if userRole == users.UserRoleOwner || userRole == users.UserRoleAdmin {
		return true
	}
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil || project == nil || project.OrganizationID == "" {
		return false
	}
	if s.orgRepo == nil {
		return false
	}
	member, err := s.orgRepo.GetMember(ctx, project.OrganizationID, userID)
	if err != nil || member == nil || member.Status != organizations.MemberStatusAccepted {
		return false
	}

	switch minPermission {
	case organizations.MemberPermissionOwner:
		return member.Permission == organizations.MemberPermissionOwner
	case organizations.MemberPermissionAdmin:
		return member.Permission == organizations.MemberPermissionOwner || member.Permission == organizations.MemberPermissionAdmin
	case organizations.MemberPermissionMember, "":
		return true
	default:
		return false
	}
}

func (s *Service) HasOrgPermission(ctx context.Context, orgID, userID string, userRole users.UserRole, minPermission organizations.MemberPermission) bool {
	if userRole == users.UserRoleOwner || userRole == users.UserRoleAdmin {
		return true
	}
	if orgID == "" || s.orgRepo == nil {
		return false
	}
	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil || member == nil || member.Status != organizations.MemberStatusAccepted {
		return false
	}
	if minPermission != "" {
		if member.Permission != minPermission && member.Permission != organizations.MemberPermissionAdmin && member.Permission != organizations.MemberPermissionOwner {
			return false
		}
	}
	return true
}

func (s *Service) IsMemberOrOwner(ctx context.Context, projectID string, userID string, userRole users.UserRole) bool {
	return s.HasPermission(ctx, projectID, userID, userRole, "")
}

func (s *Service) ListProjectsByOrganization(ctx context.Context, organizationID string, search string, limit, offset int) ([]ProjectConfig, int, error) {
	return s.projectRepo.ListByOrganization(ctx, organizationID, search, limit, offset)
}

func (s *Service) ListProjects(ctx context.Context, search string, limit, offset int) ([]ProjectConfig, int, error) {
	return s.projectRepo.ListAll(ctx, search, limit, offset)
}

func (s *Service) CountProjectsByUser(ctx context.Context, userID string) (int, error) {
	return s.projectRepo.CountByUser(ctx, userID)
}

func (s *Service) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return ErrProjectIDRequired
	}
	return s.projectRepo.Delete(ctx, id)
}
