package http

import (
	"context"
	"errors"
	"time"

	authservices "codedock/internal/modules/auth"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/operations"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/system"
	"codedock/internal/modules/users"
)

type operationsServiceAdapter struct {
	svc *operations.OperationService
}

func (o *operationsServiceAdapter) Get(ctx context.Context, id string) (*operations.Operation, error) {
	if o.svc == nil {
		return nil, errors.New("operation service unavailable")
	}
	return o.svc.Get(ctx, id)
}

func (o *operationsServiceAdapter) List(ctx context.Context, project, user string) ([]operations.Operation, error) {
	if o.svc == nil {
		return []operations.Operation{}, nil
	}
	return o.svc.List(ctx, project, user)
}

func (o *operationsServiceAdapter) Cancel(ctx context.Context, id, user string) error {
	if o.svc == nil {
		return errors.New("operation service unavailable")
	}
	return o.svc.Cancel(ctx, id, user)
}

type metricsServiceAdapter struct {
	svc *monitoring.MetricsService
}

func (m *metricsServiceAdapter) GetServiceMetrics(ctx context.Context, serviceID, rangeParam string) (any, error) {
	if m.svc == nil {
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
	return m.svc.GetServiceMetrics(ctx, monitoring.ServiceMetricsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Step:      step,
	})
}

type logServiceAdapter struct {
	svc *monitoring.LogService
}

func (l *logServiceAdapter) GetHistoricalLogs(ctx context.Context, serviceID, rangeParam string, limit int) (any, error) {
	if l.svc == nil {
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
	return l.svc.GetHistoricalLogs(ctx, monitoring.HistoricalLogsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Limit:     limit,
	})
}

type updaterServiceAdapter struct {
	svc *system.UpdaterService
}

func (u *updaterServiceAdapter) GetStatus() any {
	if u.svc == nil {
		return nil
	}
	return u.svc.GetStatus()
}

func (u *updaterServiceAdapter) CheckUpdate(ctx context.Context) (any, error) {
	if u.svc == nil {
		return nil, errors.New("updater unavailable")
	}
	return u.svc.CheckForUpdates(ctx)
}

func (u *updaterServiceAdapter) ApplyUpdate(ctx context.Context) (any, error) {
	if u.svc == nil {
		return nil, errors.New("updater unavailable")
	}
	return nil, u.svc.DeployUpdate(ctx)
}

type systemServiceAdapter struct {
	svc *system.SystemService
}

func (s *systemServiceAdapter) GetStats(ctx context.Context) (any, error) {
	if s.svc == nil {
		return nil, errors.New("system service unavailable")
	}
	return s.svc.GetStats()
}

func (s *systemServiceAdapter) Restart(ctx context.Context) error { return nil }
func (s *systemServiceAdapter) Cleanup(ctx context.Context) error {
	if s.svc == nil {
		return nil
	}
	return s.svc.Cleanup()
}

type systemAdminFinderAdapter struct {
	userRepo users.UserRepository
}

func (a *systemAdminFinderAdapter) FindFoundingAdmin(ctx context.Context) (id, email, role string, verified, totp bool, err error) {
	if a.userRepo == nil {
		return "", "", "", false, false, errors.New("user repository unavailable")
	}
	adminUser, err := a.userRepo.FindFoundingAdmin(ctx)
	if err != nil || adminUser == nil {
		return "", "", "", false, false, err
	}
	return adminUser.ID, adminUser.Email, string(adminUser.Role), adminUser.EmailVerified, adminUser.TOTPEnabled, nil
}

type systemTokenGenAdapter struct {
	tokenService *authservices.TokenService
}

func (a *systemTokenGenAdapter) GenerateCLIToken(id, email, role string, verified, totp bool) (string, error) {
	if a.tokenService == nil {
		return "", errors.New("token service unavailable")
	}
	u := &users.User{
		ID:            id,
		Email:         email,
		Role:          users.UserRole(role),
		EmailVerified: verified,
		TOTPEnabled:   totp,
	}
	return a.tokenService.GenerateToken(u)
}

type exampleServiceAdapter struct {
	svc *monitoring.ExampleService
}

func (e *exampleServiceAdapter) ListExamples() (any, error) {
	if e.svc == nil {
		return []any{}, nil
	}
	return e.svc.ListExamples()
}

type notificationServiceAdapter struct{}

func (n *notificationServiceAdapter) TestNotification(ctx context.Context, provider string) error {
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

type monitoringSettingsAdapter struct {
	repo settingsmodule.Repository
}

func (a *monitoringSettingsAdapter) GetServerSettings(ctx context.Context) (*monitoring.ServerSettings, error) {
	if a.repo == nil {
		return nil, nil
	}
	s, err := a.repo.GetServerSettings(ctx)
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

func (a *monitoringSettingsAdapter) UpdateServerSettings(ctx context.Context, settings *monitoring.ServerSettings) error {
	if a.repo == nil || settings == nil {
		return nil
	}
	s, err := a.repo.GetServerSettings(ctx)
	if err != nil || s == nil {
		return err
	}
	s.AutoUpdateEnabled = settings.AutoUpdateEnabled
	s.CurrentVersion = settings.CurrentVersion
	s.LatestVersion = settings.LatestVersion
	s.UpdateCheckCron = settings.UpdateCheckCron
	return a.repo.UpdateServerSettings(ctx, s)
}
