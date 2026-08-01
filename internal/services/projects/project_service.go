package projects

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type ProjectService struct {
	projectRepo    repositories.ProjectRepository
	envRepo        repositories.EnvironmentRepository
	appRepo        repositories.AppServiceRepository
	serviceVarRepo repositories.ServiceVarRepository
	settingsRepo   repositories.SettingsRepository
	orgRepo        repositories.OrganizationRepository
}

func NewProjectService(pr repositories.ProjectRepository, er repositories.EnvironmentRepository, ar repositories.AppServiceRepository, svr repositories.ServiceVarRepository, sr repositories.SettingsRepository, orgRepo repositories.OrganizationRepository) *ProjectService {
	return &ProjectService{
		projectRepo:    pr,
		envRepo:        er,
		appRepo:        ar,
		serviceVarRepo: svr,
		settingsRepo:   sr,
		orgRepo:        orgRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, description string) (*models.ProjectConfig, error) {
	if name == "" {
		name = utils.GenerateRandomName()
	}
	id := uuid.NewString()
	now := time.Now()
	p := &models.ProjectConfig{
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

func (s *ProjectService) CreateProjectFromRequest(ctx context.Context, req *models.CreateProjectRequest) (*models.ProjectConfig, error) {
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
	p := &models.ProjectConfig{
		ID:             id,
		OrganizationID: orgID,
		ServerID:       serverID,
		Name:           req.Name,
		Description:    req.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) CreateProjectWithMemberFromRequest(ctx context.Context, req *models.CreateProjectRequest, userID, role string) (*models.ProjectConfig, error) {
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
	p := &models.ProjectConfig{
		ID:             id,
		OrganizationID: orgID,
		ServerID:       serverID,
		Name:           req.Name,
		Description:    req.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.ProjectConfig, error) {
	if id == "" {
		return nil, errors.New("project id is required")
	}
	return s.projectRepo.Get(ctx, id)
}

func (s *ProjectService) HasPermission(ctx context.Context, projectID, userID string, userRole models.UserRole, minPermission models.MemberPermission) bool {
	if userRole == models.UserRoleOwner || userRole == models.UserRoleAdmin {
		return true
	}
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil || project == nil {
		return false
	}
	if project.OrganizationID == "" {
		return false
	}

	member, err := s.orgRepo.GetMember(ctx, project.OrganizationID, userID)
	if err != nil || member == nil || member.Status != models.MemberStatusAccepted {
		return false
	}

	switch minPermission {
	case models.MemberPermissionOwner:
		return member.Permission == models.MemberPermissionOwner
	case models.MemberPermissionAdmin:
		return member.Permission == models.MemberPermissionOwner || member.Permission == models.MemberPermissionAdmin
	case models.MemberPermissionMember, "":
		return true
	default:
		return false
	}
}

func (s *ProjectService) AccessibleServiceIDs(ctx context.Context, userID string, userRole models.UserRole) (map[string]struct{}, error) {
	services, err := s.appRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	accessibleProjects := make(map[string]bool)
	accessibleServices := make(map[string]struct{})
	for _, service := range services {
		if service == nil || service.ProjectID == "" {
			continue
		}
		allowed, ok := accessibleProjects[service.ProjectID]
		if !ok {
			allowed = s.HasPermission(ctx, service.ProjectID, userID, userRole, "")
			accessibleProjects[service.ProjectID] = allowed
		}
		if allowed {
			accessibleServices[service.ID] = struct{}{}
		}
	}
	return accessibleServices, nil
}

func (s *ProjectService) HasOrgPermission(ctx context.Context, orgID, userID string, userRole models.UserRole, minPermission models.MemberPermission) bool {
	if userRole == models.UserRoleOwner || userRole == models.UserRoleAdmin {
		return true
	}
	if orgID == "" {
		return false
	}
	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil || member == nil || member.Status != models.MemberStatusAccepted {
		return false
	}
	if minPermission != "" {
		if member.Permission != minPermission && member.Permission != models.MemberPermissionAdmin && member.Permission != models.MemberPermissionOwner {
			return false
		}
	}
	return true
}

func (s *ProjectService) IsMemberOrOwner(ctx context.Context, projectID string, userID string, userRole models.UserRole) bool {
	hasPerm := s.HasPermission(ctx, projectID, userID, userRole, "")
	return hasPerm
}
func (s *ProjectService) ListProjectsByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]models.ProjectConfig, int, error) {
	return s.projectRepo.ListByOrganization(ctx, organizationID, limit, offset)
}

func (s *ProjectService) ListProjects(ctx context.Context, limit, offset int) ([]models.ProjectConfig, int, error) {
	return s.projectRepo.ListAll(ctx, limit, offset)
}

func (s *ProjectService) CountProjectsByUser(ctx context.Context, userID string) (int, error) {
	return s.projectRepo.CountByUser(ctx, userID)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("project id is required")
	}
	return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, projectID, name string) (*models.EnvironmentConfig, error) {
	if projectID == "" || name == "" {
		return nil, errors.New("project id and environment name required")
	}
	env := &models.EnvironmentConfig{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      name,
	}
	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID string) ([]models.EnvironmentConfig, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.envRepo.ListByProject(ctx, projectID)
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("environment id is required")
	}
	return s.envRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil || svc.Name == "" || svc.ProjectID == "" {
		return errors.New("valid app service with name and projectId required")
	}
	if svc.ID == "" {
		svc.ID = uuid.New().String()
	}
	if svc.Status == "" {
		svc.Status = "stopped"
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}
	svc.CreatedAt = time.Now()
	svc.UpdatedAt = time.Now()
	return s.appRepo.Create(ctx, svc)
}

func (s *ProjectService) GetAppService(ctx context.Context, id string) (*models.AppService, error) {
	if id == "" {
		return nil, errors.New("service id is required")
	}
	return s.appRepo.GetByID(ctx, id)
}

func (s *ProjectService) ListAppServicesByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.appRepo.ListByProject(ctx, projectID)
}

func (s *ProjectService) UpdateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil || svc.ID == "" {
		return errors.New("valid app service required for update")
	}
	svc.UpdatedAt = time.Now()
	return s.appRepo.Update(ctx, svc)
}

func (s *ProjectService) DeleteAppService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("service id is required")
	}
	return s.appRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateServiceVariable(ctx context.Context, v *models.Variable) error {
	if v == nil || v.ServiceID == "" || v.Key == "" {
		return errors.New("valid variable required")
	}
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return s.serviceVarRepo.Create(ctx, v)
}

func (s *ProjectService) ListServiceVariables(ctx context.Context, serviceID string) ([]*models.Variable, error) {
	if serviceID == "" {
		return nil, errors.New("service id is required")
	}
	return s.serviceVarRepo.ListByService(ctx, serviceID)
}

func (s *ProjectService) DeleteServiceVariable(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("variable id is required")
	}
	return s.serviceVarRepo.Delete(ctx, id)
}
