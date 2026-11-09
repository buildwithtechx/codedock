package adapters

import (
	"context"
	"time"

	"codedock/internal/modules/apps"
	"codedock/internal/modules/auth"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/users"
	"codedock/internal/platform/ssh"
)

type SSHServerRepoAdapter struct {
	Repo servers.Repository
}

func (a *SSHServerRepoAdapter) UpdateStatus(ctx context.Context, id string, status ssh.ServerStatus) error {
	return a.Repo.UpdateStatus(ctx, id, servers.ServerStatus(status))
}

func (a *SSHServerRepoAdapter) UpdateMetrics(ctx context.Context, id string, metrics []byte) error {
	return a.Repo.UpdateMetrics(ctx, id, metrics)
}

func (a *SSHServerRepoAdapter) ListByUser(ctx context.Context, userID, orgID string) ([]*ssh.Server, error) {
	list, err := a.Repo.ListByUser(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	res := make([]*ssh.Server, len(list))
	for i, s := range list {
		res[i] = &ssh.Server{
			ID:                    s.ID,
			OrganizationID:        s.OrganizationID,
			UserID:                s.UserID,
			Name:                  s.Name,
			IPAddress:             s.IPAddress,
			IsLocal:               s.IsLocal,
			SSHHost:               s.SSHHost,
			SSHPort:               s.SSHPort,
			SSHUser:               s.SSHUser,
			SSHAuthMethod:         s.SSHAuthMethod,
			SSHKey:                s.SSHKey,
			SSHPrivateKey:         s.SSHPrivateKey,
			SSHPassword:           s.SSHPassword,
			SSHTransport:          s.SSHTransport,
			SSHJumpHost:           s.SSHJumpHost,
			IsControlPlane:        s.IsControlPlane,
			Status:                ssh.ServerStatus(s.Status),
			Provider:              s.Provider,
			ExternalID:            s.ExternalID,
			Region:                s.Region,
			ServerType:            s.ServerType,
			WorkerToken:           s.WorkerToken,
			LastSeenAt:            s.LastSeenAt,
			BuildCacheBudgetBytes: s.BuildCacheBudgetBytes,
			Metrics:               s.Metrics,
			CreatedAt:             s.CreatedAt,
			UpdatedAt:             s.UpdatedAt,
		}
	}
	return res, nil
}

type ServerUserLookupAdapter struct {
	Users users.UserRepository
}

func (a *ServerUserLookupAdapter) GetUserByID(ctx context.Context, id string) (*servers.UserInfo, error) {
	if a.Users == nil {
		return nil, nil
	}
	u, err := a.Users.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return &servers.UserInfo{
		ID:       u.ID,
		Role:     string(u.Role),
		PlanType: u.PlanType,
	}, nil
}

type ProjectCreatorAdapter struct {
	Projects projects.ProjectRepository
	Apps     apps.AppServiceRepository
}

func (a *ProjectCreatorAdapter) CreateProject(ctx context.Context, id, name, description string) error {
	now := time.Now().UTC()
	return a.Projects.Create(ctx, &projects.ProjectConfig{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (a *ProjectCreatorAdapter) CreateApp(ctx context.Context, id, projectID, name, imageRef, runtimeMode string, port int, status string) error {
	now := time.Now().UTC()
	return a.Apps.Create(ctx, &apps.AppService{
		ID:           id,
		ProjectID:    projectID,
		Name:         name,
		ImageRef:     imageRef,
		RuntimeMode:  apps.RuntimeMode(runtimeMode),
		InternalPort: port,
		Status:       apps.AppServiceStatus(status),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

type ServerAuditLoggerAdapter struct {
	Audit *AuditDomainAdapter
}

func (a *ServerAuditLoggerAdapter) Log(ctx context.Context, userID, action, resource, ip string) {
	if a.Audit != nil {
		a.Audit.Log(ctx, userID, action, resource, ip, nil)
	}
}

type ServerTokenValidatorAdapter struct {
	Tokens *auth.TokenService
}

func (a *ServerTokenValidatorAdapter) ValidateToken(tokenStr string) (map[string]any, error) {
	if a.Tokens == nil {
		return nil, nil
	}
	claims, err := a.Tokens.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	res := make(map[string]any, len(claims))
	for k, v := range claims {
		res[k] = v
	}
	return res, nil
}
