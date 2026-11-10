package adapters

import (
	"context"
	"errors"
	"time"

	"codedock/internal/modules/audit/operations"
	authservices "codedock/internal/modules/auth"
	"codedock/internal/modules/monitoring"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/system"
	"codedock/internal/modules/users"
)

type OperationsServiceAdapter struct {
	Svc *operations.OperationService
}

func (o *OperationsServiceAdapter) Get(ctx context.Context, id string) (*operations.Operation, error) {
	if o.Svc == nil {
		return nil, errors.New("operation service unavailable")
	}
	return o.Svc.Get(ctx, id)
}

func (o *OperationsServiceAdapter) List(ctx context.Context, project, user string) ([]operations.Operation, error) {
	if o.Svc == nil {
		return []operations.Operation{}, nil
	}
	return o.Svc.List(ctx, project, user)
}

func (o *OperationsServiceAdapter) Cancel(ctx context.Context, id, user string) error {
	if o.Svc == nil {
		return errors.New("operation service unavailable")
	}
	return o.Svc.Cancel(ctx, id, user)
}

type MetricsServiceAdapter struct {
	Svc *monitoring.MetricsService
}

func (m *MetricsServiceAdapter) GetServiceMetrics(ctx context.Context, serviceID, rangeParam string) (any, error) {
	if m.Svc == nil {
		return nil, nil
	}
	end := time.Now()
	var start time.Time
	var step string
	switch rangeParam {
	case "7d":
		start = end.Add(-7 * 24 * time.Hour)
		step = "1h"
	case "24h":
		start = end.Add(-24 * time.Hour)
		step = "5m"
	case "1h":
		start = end.Add(-1 * time.Hour)
		step = "1m"
	default:
		start = end.Add(-24 * time.Hour)
		step = "5m"
	}
	return m.Svc.GetServiceMetrics(ctx, monitoring.ServiceMetricsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Step:      step,
	})
}

type LogServiceAdapter struct {
	Svc *monitoring.LogService
}

func (l *LogServiceAdapter) GetHistoricalLogs(ctx context.Context, serviceID, rangeParam string, limit int) (any, error) {
	if l.Svc == nil {
		return nil, nil
	}
	end := time.Now()
	var start time.Time
	switch rangeParam {
	case "7d":
		start = end.Add(-7 * 24 * time.Hour)
	case "24h":
		start = end.Add(-24 * time.Hour)
	case "1h":
		start = end.Add(-1 * time.Hour)
	default:
		start = end.Add(-24 * time.Hour)
	}
	return l.Svc.GetHistoricalLogs(ctx, monitoring.HistoricalLogsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Limit:     limit,
	})
}

type UpdaterServiceAdapter struct {
	Svc *system.UpdaterService
}

func (u *UpdaterServiceAdapter) GetStatus() any {
	if u.Svc == nil {
		return nil
	}
	return u.Svc.GetStatus()
}

func (u *UpdaterServiceAdapter) CheckUpdate(ctx context.Context) (any, error) {
	if u.Svc == nil {
		return nil, errors.New("updater unavailable")
	}
	return u.Svc.CheckForUpdates(ctx)
}

func (u *UpdaterServiceAdapter) ApplyUpdate(ctx context.Context) (any, error) {
	if u.Svc == nil {
		return nil, errors.New("updater unavailable")
	}
	return nil, u.Svc.DeployUpdate(ctx)
}

type SystemServiceAdapter struct {
	Svc *system.SystemService
}

func (s *SystemServiceAdapter) GetStats(ctx context.Context) (any, error) {
	if s.Svc == nil {
		return nil, errors.New("system service unavailable")
	}
	return s.Svc.GetStats()
}

func (s *SystemServiceAdapter) Restart(ctx context.Context) error { return nil }
func (s *SystemServiceAdapter) Cleanup(ctx context.Context) error {
	if s.Svc == nil {
		return nil
	}
	return s.Svc.Cleanup()
}

type SystemAdminFinderAdapter struct {
	UserRepo users.UserRepository
}

func (a *SystemAdminFinderAdapter) FindFoundingAdmin(ctx context.Context) (id, email, role string, verified, totp bool, err error) {
	if a.UserRepo == nil {
		return "", "", "", false, false, errors.New("user repository unavailable")
	}
	adminUser, err := a.UserRepo.FindFoundingAdmin(ctx)
	if err != nil || adminUser == nil {
		return "", "", "", false, false, err
	}
	return adminUser.ID, adminUser.Email, string(adminUser.Role), adminUser.EmailVerified, adminUser.TOTPEnabled, nil
}

type SystemTokenGenAdapter struct {
	TokenService *authservices.TokenService
}

func (a *SystemTokenGenAdapter) GenerateCLIToken(id, email, role string, verified, totp bool) (string, error) {
	if a.TokenService == nil {
		return "", errors.New("token service unavailable")
	}
	u := &users.User{
		ID:            id,
		Email:         email,
		Role:          users.UserRole(role),
		EmailVerified: verified,
		TOTPEnabled:   totp,
	}
	return a.TokenService.GenerateToken(u)
}

type ExampleServiceAdapter struct {
	Svc *monitoring.ExampleService
}

func (e *ExampleServiceAdapter) ListExamples() (any, error) {
	if e.Svc == nil {
		return []any{}, nil
	}
	return e.Svc.ListExamples()
}

type NotificationServiceAdapter struct{}

func (n *NotificationServiceAdapter) TestNotification(ctx context.Context, provider string) error {
	return nil
}

func toMonitoringUser(u *users.User) *monitoring.User {
	if u == nil {
		return nil
	}
	return &monitoring.User{
		ID:                   u.ID,
		Email:                u.Email,
		Name:                 u.Name,
		Role:                 string(u.Role),
		PlanType:             u.PlanType,
		StripeCustomerID:     u.StripeCustomerID,
		StripeSubscriptionID: u.StripeSubscriptionID,
		StripePriceID:        u.StripePriceID,
	}
}

type MonitoringSettingsAdapter struct {
	Repo settingsmodule.Repository
}

func (a *MonitoringSettingsAdapter) GetServerSettings(ctx context.Context) (*monitoring.ServerSettings, error) {
	if a.Repo == nil {
		return nil, nil
	}
	s, err := a.Repo.GetServerSettings(ctx)
	if err != nil || s == nil {
		return nil, err
	}
	return &monitoring.ServerSettings{
		AutoUpdateEnabled: s.AutoUpdateEnabled,
		CurrentVersion:    s.CurrentVersion,
		LatestVersion:     s.LatestVersion,
		UpdateCheckCron:   s.UpdateCheckCron,
	}, nil
}

func (a *MonitoringSettingsAdapter) UpdateServerSettings(ctx context.Context, settings *monitoring.ServerSettings) error {
	if a.Repo == nil || settings == nil {
		return nil
	}
	s, err := a.Repo.GetServerSettings(ctx)
	if err != nil || s == nil {
		return err
	}
	s.AutoUpdateEnabled = settings.AutoUpdateEnabled
	s.CurrentVersion = settings.CurrentVersion
	s.LatestVersion = settings.LatestVersion
	s.UpdateCheckCron = settings.UpdateCheckCron
	return a.Repo.UpdateServerSettings(ctx, s)
}
