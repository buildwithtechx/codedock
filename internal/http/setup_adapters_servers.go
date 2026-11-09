package http

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

type sshServerRepoAdapter struct {
	repo servers.Repository
}

func (a *sshServerRepoAdapter) UpdateStatus(ctx context.Context, id string, status ssh.ServerStatus) error {
	return a.repo.UpdateStatus(ctx, id, servers.ServerStatus(status))
}

func (a *sshServerRepoAdapter) UpdateMetrics(ctx context.Context, id string, metrics []byte) error {
	return a.repo.UpdateMetrics(ctx, id, metrics)
}

func (a *sshServerRepoAdapter) ListByUser(ctx context.Context, userID, orgID string) ([]*ssh.Server, error) {
	list, err := a.repo.ListByUser(ctx, userID, orgID)
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

type serverUserLookupAdapter struct {
	users users.UserRepository
}

func (a *serverUserLookupAdapter) GetUserByID(ctx context.Context, id string) (*servers.UserInfo, error) {
	if a.users == nil {
		return nil, nil
	}
	u, err := a.users.GetUserByID(ctx, id)
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

type projectCreatorAdapter struct {
	projects projects.ProjectRepository
	apps     apps.AppServiceRepository
}

func (a *projectCreatorAdapter) CreateProject(ctx context.Context, id, name, description string) error {
	now := time.Now().UTC()
	return a.projects.Create(ctx, &projects.ProjectConfig{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (a *projectCreatorAdapter) CreateApp(ctx context.Context, id, projectID, name, imageRef, runtimeMode string, port int, status string) error {
	now := time.Now().UTC()
	return a.apps.Create(ctx, &apps.AppService{
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

type serverAuditLoggerAdapter struct {
	audit *auditDomainAdapter
}

func (a *serverAuditLoggerAdapter) Log(ctx context.Context, userID, action, resource, ip string) {
	if a.audit != nil {
		a.audit.Log(ctx, userID, action, resource, ip, nil)
	}
}

type serverTokenValidatorAdapter struct {
	tokens *auth.TokenService
}

func (a *serverTokenValidatorAdapter) ValidateToken(tokenStr string) (map[string]any, error) {
	if a.tokens == nil {
		return nil, nil
	}
	claims, err := a.tokens.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	res := make(map[string]any, len(claims))
	for k, v := range claims {
		res[k] = v
	}
	return res, nil
}

func createServerHandlers(d *serverSetupDeps) (*servers.Handler, *servers.ServerMetricsWSHandler, *servers.TakeoverHandler, *servers.ServerNetworksHandler) {
	serversRepo := servers.NewRepository(d.db, d.v)
	serversUserLookup := &serverUserLookupAdapter{users: d.userRepo}
	serversService := servers.NewService(serversRepo, serversUserLookup, d.sshManager)
	serverAudit := &serverAuditLoggerAdapter{audit: &auditDomainAdapter{audit: d.auditService}}
	serversHandler := servers.NewHandler(serversService, serverAudit)
	takeoverRepo := servers.NewTakeoverRepository(d.db, d.v)
	takeoverCreator := &projectCreatorAdapter{projects: d.projectRepo, apps: d.appRepo}
	takeoverService := servers.NewTakeoverService(takeoverRepo, takeoverCreator)
	serversTakeoverHandler := servers.NewTakeoverHandler(takeoverService, takeoverRepo)
	serversNetworksHandler := servers.NewServerNetworksHandler()
	serversWSHandler := servers.NewServerMetricsWSHandler(&serverTokenValidatorAdapter{tokens: d.tokenService}, serversService, serversUserLookup)
	return serversHandler, serversWSHandler, serversTakeoverHandler, serversNetworksHandler
}
