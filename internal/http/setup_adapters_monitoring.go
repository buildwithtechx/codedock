package http

import (
	"context"
	"errors"
	"time"

	authservices "codedock/internal/modules/auth"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/users"
	"codedock/internal/repositories"
	"codedock/internal/services/operations"
	systemservices "codedock/internal/services/system"
)

type operationsServiceAdapter struct {
	svc *operations.Service
}

func (o *operationsServiceAdapter) Get(ctx context.Context, id string) (*monitoring.Operation, error) {
	if o.svc == nil {
		return nil, errors.New("operation service unavailable")
	}
	op, err := o.svc.Get(ctx, id)
	if err != nil || op == nil {
		return nil, err
	}
	return &monitoring.Operation{
		ID:        op.ID,
		UserID:    op.UserID,
		Kind:      op.Kind,
		Status:    op.Status,
		CreatedAt: op.UpdatedAt,
		UpdatedAt: op.UpdatedAt,
	}, nil
}

func (o *operationsServiceAdapter) List(ctx context.Context, project, user string) ([]monitoring.Operation, error) {
	if o.svc == nil {
		return []monitoring.Operation{}, nil
	}
	list, err := o.svc.List(ctx, project, user)
	if err != nil {
		return nil, err
	}
	result := make([]monitoring.Operation, len(list))
	for i, op := range list {
		result[i] = monitoring.Operation{
			ID:        op.ID,
			UserID:    op.UserID,
			Kind:      op.Kind,
			Status:    op.Status,
			CreatedAt: op.UpdatedAt,
			UpdatedAt: op.UpdatedAt,
		}
	}
	return result, nil
}

func (o *operationsServiceAdapter) Cancel(ctx context.Context, id, user string) error {
	if o.svc == nil {
		return errors.New("operation service unavailable")
	}
	return o.svc.Cancel(ctx, id, user)
}

type metricsServiceAdapter struct {
	svc *systemservices.MetricsService
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
	return m.svc.GetServiceMetrics(ctx, systemservices.ServiceMetricsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Step:      step,
	})
}

type logServiceAdapter struct {
	svc *systemservices.LogService
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
	return l.svc.GetHistoricalLogs(ctx, systemservices.HistoricalLogsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Limit:     limit,
	})
}

type updaterServiceAdapter struct {
	svc *systemservices.UpdaterService
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
	svc *systemservices.SystemService
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
	userRepo repositories.UserRepository
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
	svc *systemservices.ExampleService
}

func (e *exampleServiceAdapter) ListExamples() (any, error) {
	if e.svc == nil {
		return []any{}, nil
	}
	return e.svc.ListExamples()
}

type notificationServiceAdapter struct {
	svc *systemservices.NotificationService
}

func (n *notificationServiceAdapter) TestNotification(ctx context.Context, provider string) error {
	if n.svc == nil {
		return errors.New("notification service unavailable")
	}
	return n.svc.TestGlobalNotification(ctx, provider)
}
