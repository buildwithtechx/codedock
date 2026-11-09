package http

import (
	"context"
	"strings"
	"time"

	"codedock/internal/modules/apps"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/notifications"
)

type notifSettingsAdapter struct {
	repo settingsmodule.Repository
}

func (a *notifSettingsAdapter) GetServerSettings(ctx context.Context) (*notifications.ServerSettings, error) {
	if a.repo == nil {
		return nil, nil
	}
	s, err := a.repo.GetServerSettings(ctx)
	if err != nil || s == nil {
		return nil, err
	}
	return &notifications.ServerSettings{
		InstanceName: s.SiteName,
	}, nil
}

type notifUserLookupAdapter struct {
	repo users.UserRepository
}

func (a *notifUserLookupAdapter) GetUserByID(ctx context.Context, id string) (*notifications.User, error) {
	if a.repo == nil {
		return nil, nil
	}
	u, err := a.repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return &notifications.User{
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

type notifDeviceTokenAdapter struct {
	repo monitoring.DeviceTokenRepository
}

func (a *notifDeviceTokenAdapter) ListTokensByUser(ctx context.Context, userID string) ([]*notifications.DeviceToken, error) {
	if a.repo == nil {
		return []*notifications.DeviceToken{}, nil
	}
	tokens, err := a.repo.ListTokensByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*notifications.DeviceToken, len(tokens))
	for i, t := range tokens {
		res[i] = &notifications.DeviceToken{
			ID:        t.ID,
			UserID:    t.UserID,
			Token:     t.Token,
			Platform:  t.Platform,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		}
	}
	return res, nil
}

func (a *notifDeviceTokenAdapter) ListAllTokens(ctx context.Context) ([]*notifications.DeviceToken, error) {
	if a.repo == nil {
		return []*notifications.DeviceToken{}, nil
	}
	tokens, err := a.repo.ListAllTokens(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*notifications.DeviceToken, len(tokens))
	for i, t := range tokens {
		res[i] = &notifications.DeviceToken{
			ID:        t.ID,
			UserID:    t.UserID,
			Token:     t.Token,
			Platform:  t.Platform,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		}
	}
	return res, nil
}

type notifAppRepoAdapter struct {
	repo apps.AppServiceRepository
}

func (a *notifAppRepoAdapter) GetByID(ctx context.Context, id string) (*notifications.AppService, error) {
	if a.repo == nil {
		return nil, nil
	}
	app, err := a.repo.GetByID(ctx, id)
	if err != nil || app == nil {
		return nil, err
	}
	return &notifications.AppService{
		ID:        app.ID,
		ProjectID: app.ProjectID,
		Name:      app.Name,
		Branch:    app.Branch,
	}, nil
}

func (a *notifAppRepoAdapter) ListWebhooksByService(ctx context.Context, serviceID string) ([]*notifications.Webhook, error) {
	if a.repo == nil {
		return []*notifications.Webhook{}, nil
	}
	hooks, err := a.repo.ListWebhooksByService(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	res := make([]*notifications.Webhook, len(hooks))
	for i, h := range hooks {
		var events []string
		if h.Events != "" {
			events = strings.Split(h.Events, ",")
		}
		res[i] = &notifications.Webhook{
			ID:        h.ID,
			ServiceID: h.ServiceID,
			URL:       h.URL,
			Secret:    h.Secret,
			Events:    events,
			Enabled:   true,
		}
	}
	return res, nil
}

type notifSettingsServiceAdapter struct {
	svc     *monitoring.NotificationSettingsService
	orgRepo organizations.OrganizationRepository
}

func (n *notifSettingsServiceAdapter) GetNotificationSettings(ctx context.Context) (*monitoring.NotificationSettings, error) {
	if n.svc == nil {
		return &monitoring.NotificationSettings{}, nil
	}
	return n.svc.GetNotificationSettings(ctx)
}

func (n *notifSettingsServiceAdapter) UpdateNotificationSettings(ctx context.Context, settings *monitoring.NotificationSettings) error {
	if n.svc == nil {
		return nil
	}
	return n.svc.UpdateNotificationSettings(ctx, settings)
}

func (n *notifSettingsServiceAdapter) ListSubscriptions(ctx context.Context, userID, orgID string) (any, error) {
	if n.svc == nil {
		return []monitoring.NotificationSubscription{}, nil
	}
	return n.svc.ListSubscriptions(ctx, userID, orgID)
}

func (n *notifSettingsServiceAdapter) UpsertSubscription(ctx context.Context, userID, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	sub := &monitoring.NotificationSubscription{
		UserID:         userID,
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.svc.UpsertSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (n *notifSettingsServiceAdapter) ListDefaults(ctx context.Context, orgID string) (any, error) {
	if n.svc == nil {
		return []monitoring.NotificationDefault{}, nil
	}
	return n.svc.ListDefaults(ctx, orgID)
}

func (n *notifSettingsServiceAdapter) UpsertDefault(ctx context.Context, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	def := &monitoring.NotificationDefault{
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.svc.UpsertDefault(ctx, def); err != nil {
		return nil, err
	}
	return def, nil
}

func (n *notifSettingsServiceAdapter) ListDeliveries(ctx context.Context, limit int) (any, error) {
	if n.svc == nil {
		return []monitoring.NotificationDelivery{}, nil
	}
	return n.svc.ListDeliveries(ctx, limit)
}

func (n *notifSettingsServiceAdapter) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	return n.svc.RegisterDeviceToken(ctx, userID, token, platform)
}

func (n *notifSettingsServiceAdapter) DeleteDeviceToken(ctx context.Context, userID, token string) error {
	if n.svc == nil {
		return nil
	}
	return n.svc.DeleteDeviceToken(ctx, userID, token)
}

func (n *notifSettingsServiceAdapter) ListDeviceTokens(ctx context.Context, userID string) (any, error) {
	if n.svc == nil {
		return []*monitoring.DeviceToken{}, nil
	}
	return n.svc.ListDeviceTokens(ctx, userID)
}

func (n *notifSettingsServiceAdapter) ResolveOrgID(ctx context.Context, userID string) string {
	if n.orgRepo == nil {
		return ""
	}
	orgs, err := n.orgRepo.ListByUser(ctx, userID)
	if err != nil || len(orgs) == 0 {
		return ""
	}
	return orgs[0].ID
}

type wsAccessCheckerAdapter struct {
	appService     *apps.Service
	projectService *projects.Service
}

func (w *wsAccessCheckerAdapter) HasAccess(ctx context.Context, serviceID, userID, role string) (bool, error) {
	if role == "admin" {
		return true, nil
	}
	if w.appService == nil || w.projectService == nil {
		return true, nil
	}
	svc, err := w.appService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil {
		return false, err
	}
	return w.projectService.HasPermission(ctx, svc.ProjectID, userID, users.UserRole(role), ""), nil
}

type wsUserValidatorAdapter struct {
	userRepo users.UserRepository
}

func (w *wsUserValidatorAdapter) IsUserActive(ctx context.Context, userID string) (bool, error) {
	if w.userRepo == nil {
		return true, nil
	}
	u, err := w.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil || !u.IsActive {
		return false, err
	}
	return true, nil
}
