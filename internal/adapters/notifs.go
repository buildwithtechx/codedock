package adapters

import (
	"context"
	"strings"
	"time"

	"codedock/internal/modules/apps"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/notifications"
)

type NotifSettingsAdapter struct {
	Repo settingsmodule.Repository
}

func (a *NotifSettingsAdapter) GetServerSettings(ctx context.Context) (*notifications.ServerSettings, error) {
	if a.Repo == nil {
		return nil, nil
	}
	s, err := a.Repo.GetServerSettings(ctx)
	if err != nil || s == nil {
		return nil, err
	}
	return &notifications.ServerSettings{
		InstanceName: s.SiteName,
	}, nil
}

type NotifUserLookupAdapter struct {
	Repo users.UserRepository
}

func (a *NotifUserLookupAdapter) GetUserByID(ctx context.Context, id string) (*notifications.User, error) {
	if a.Repo == nil {
		return nil, nil
	}
	u, err := a.Repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return &notifications.User{
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

type NotifDeviceTokenAdapter struct {
	Repo notifsmod.DeviceTokenRepository
}

func (a *NotifDeviceTokenAdapter) ListTokensByUser(ctx context.Context, userID string) ([]*notifications.DeviceToken, error) {
	if a.Repo == nil {
		return []*notifications.DeviceToken{}, nil
	}
	tokens, err := a.Repo.ListTokensByUser(ctx, userID)
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

func (a *NotifDeviceTokenAdapter) ListAllTokens(ctx context.Context) ([]*notifications.DeviceToken, error) {
	if a.Repo == nil {
		return []*notifications.DeviceToken{}, nil
	}
	tokens, err := a.Repo.ListAllTokens(ctx)
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

type NotifAppRepoAdapter struct {
	Repo apps.AppServiceRepository
}

func (a *NotifAppRepoAdapter) GetByID(ctx context.Context, id string) (*notifications.AppService, error) {
	if a.Repo == nil {
		return nil, nil
	}
	app, err := a.Repo.GetByID(ctx, id)
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

func (a *NotifAppRepoAdapter) ListWebhooksByService(ctx context.Context, serviceID string) ([]*notifications.Webhook, error) {
	if a.Repo == nil {
		return []*notifications.Webhook{}, nil
	}
	hooks, err := a.Repo.ListWebhooksByService(ctx, serviceID)
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

type NotifSettingsServiceAdapter struct {
	Svc     *notifsmod.NotificationSettingsService
	OrgRepo organizations.OrganizationRepository
}

func (n *NotifSettingsServiceAdapter) GetNotificationSettings(ctx context.Context) (*notifsmod.NotificationSettings, error) {
	if n.Svc == nil {
		return &notifsmod.NotificationSettings{}, nil
	}
	return n.Svc.GetNotificationSettings(ctx)
}

func (n *NotifSettingsServiceAdapter) UpdateNotificationSettings(ctx context.Context, settings *notifsmod.NotificationSettings) error {
	if n.Svc == nil {
		return nil
	}
	return n.Svc.UpdateNotificationSettings(ctx, settings)
}

func (n *NotifSettingsServiceAdapter) ListSubscriptions(ctx context.Context, userID, orgID string) (any, error) {
	if n.Svc == nil {
		return []notifsmod.NotificationSubscription{}, nil
	}
	return n.Svc.ListSubscriptions(ctx, userID, orgID)
}

func (n *NotifSettingsServiceAdapter) UpsertSubscription(ctx context.Context, userID, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.Svc == nil {
		return nil, nil
	}
	sub := &notifsmod.NotificationSubscription{
		UserID:         userID,
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.Svc.UpsertSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (n *NotifSettingsServiceAdapter) ListDefaults(ctx context.Context, orgID string) (any, error) {
	if n.Svc == nil {
		return []notifsmod.NotificationDefault{}, nil
	}
	return n.Svc.ListDefaults(ctx, orgID)
}

func (n *NotifSettingsServiceAdapter) UpsertDefault(ctx context.Context, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.Svc == nil {
		return nil, nil
	}
	def := &notifsmod.NotificationDefault{
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.Svc.UpsertDefault(ctx, def); err != nil {
		return nil, err
	}
	return def, nil
}

func (n *NotifSettingsServiceAdapter) ListDeliveries(ctx context.Context, limit int) (any, error) {
	if n.Svc == nil {
		return []notifsmod.NotificationDelivery{}, nil
	}
	return n.Svc.ListDeliveries(ctx, limit)
}

func (n *NotifSettingsServiceAdapter) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (any, error) {
	if n.Svc == nil {
		return nil, nil
	}
	return n.Svc.RegisterDeviceToken(ctx, userID, token, platform)
}

func (n *NotifSettingsServiceAdapter) DeleteDeviceToken(ctx context.Context, userID, token string) error {
	if n.Svc == nil {
		return nil
	}
	return n.Svc.DeleteDeviceToken(ctx, userID, token)
}

func (n *NotifSettingsServiceAdapter) ListDeviceTokens(ctx context.Context, userID string) (any, error) {
	if n.Svc == nil {
		return []*notifsmod.DeviceToken{}, nil
	}
	return n.Svc.ListDeviceTokens(ctx, userID)
}

func (n *NotifSettingsServiceAdapter) ResolveOrgID(ctx context.Context, userID string) string {
	if n.OrgRepo == nil {
		return ""
	}
	orgs, err := n.OrgRepo.ListByUser(ctx, userID)
	if err != nil || len(orgs) == 0 {
		return ""
	}
	return orgs[0].ID
}

type WSAccessCheckerAdapter struct {
	AppService     *apps.Service
	ProjectService *projects.Service
}

func (w *WSAccessCheckerAdapter) HasAccess(ctx context.Context, serviceID, userID, role string) (bool, error) {
	if role == "admin" {
		return true, nil
	}
	if w.AppService == nil || w.ProjectService == nil {
		return true, nil
	}
	svc, err := w.AppService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil {
		return false, err
	}
	return w.ProjectService.HasPermission(ctx, svc.ProjectID, userID, users.UserRole(role), ""), nil
}

type WSUserValidatorAdapter struct {
	UserRepo users.UserRepository
}

func (w *WSUserValidatorAdapter) IsUserActive(ctx context.Context, userID string) (bool, error) {
	if w.UserRepo == nil {
		return true, nil
	}
	u, err := w.UserRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil || !u.IsActive {
		return false, err
	}
	return true, nil
}
