package monitoring

import (
	"context"
	"io"
)

type AnalyticsService interface {
	Periods(ctx context.Context, orgID string) (any, error)
	Dashboard(ctx context.Context, orgID string) (any, error)
	Summary(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error)
	Overview(ctx context.Context, userID, orgID, projectID, domain, from, to, period string) (any, error)
	Geo(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error)
	GetPaths(ctx context.Context, userID, orgID, projectID string) (bool, error)
	SetPaths(ctx context.Context, userID, orgID, projectID string, enabled bool) error
	DeploymentStats(ctx context.Context, userID, orgID, projectID string) (any, error)
	Usage(ctx context.Context, userID, orgID, projectID, period string) (any, error)
	Resources(ctx context.Context, userID, orgID, projectID, period string) (any, error)
	UsageHistory(ctx context.Context, userID, orgID, projectID, period string) (any, error)
	Container(ctx context.Context, userID, orgID, projectID, serviceID, period string) (any, error)
	CurrentUsage(ctx context.Context, projectID string) (any, error)
}

type AttentionService interface {
	List(ctx context.Context, userID, orgID string) (any, error)
	Evaluate(ctx context.Context, userID, orgID string) (any, error)
	ClearHistory(ctx context.Context, userID, orgID string) error
	Acknowledge(ctx context.Context, userID, orgID, issueID string) error
	Resolve(ctx context.Context, userID, orgID, issueID string) error
	Act(ctx context.Context, userID, orgID, issueID string, req AttentionActionRequest) error
}

type CreateCheckoutSessionOpts struct {
	PriceID        string
	SuccessURL     string
	CancelURL      string
	IdempotencyKey string
	PurchaseIntent string
	CancelPending  bool
}

type BillingService interface {
	GetConfig(ctx context.Context) (any, error)
	CreateCheckoutSession(ctx context.Context, userID string, opts CreateCheckoutSessionOpts) (string, error)
	GetPendingCheckout(ctx context.Context, userID string) (any, error)
	CancelCheckoutSession(ctx context.Context, userID, id, idempotencyKey string) (any, error)
	HandleWebhook(ctx context.Context, payload []byte, sig string) error
}

type MetricsService interface {
	GetServiceMetrics(ctx context.Context, serviceID, rangeParam string) (any, error)
}

type LogService interface {
	GetHistoricalLogs(ctx context.Context, serviceID, rangeParam string, limit int) (any, error)
}

type SystemService interface {
	GetStats(ctx context.Context) (any, error)
	Restart(ctx context.Context) error
	Cleanup(ctx context.Context) error
}

type UpdaterService interface {
	GetStatus() any
	CheckUpdate(ctx context.Context) (any, error)
	ApplyUpdate(ctx context.Context) (any, error)
}

type OperationsService interface {
	Get(ctx context.Context, id string) (*Operation, error)
	List(ctx context.Context, project, user string) ([]Operation, error)
	Cancel(ctx context.Context, id, user string) error
}

type NotificationService interface {
	TestNotification(ctx context.Context, provider string) error
}

type NotificationSettingsService interface {
	GetSettings(ctx context.Context) (*NotificationSettings, error)
	UpdateSettings(ctx context.Context, req UpdateNotificationSettingsRequest) (*NotificationSettings, error)
	ListDeliveries(ctx context.Context, limit int) (any, error)
	RegisterDeviceToken(ctx context.Context, userID, token, platform string) error
}

type IngressService interface {
	List(ctx context.Context, serverID string) (any, error)
	Discover(ctx context.Context, req IngressDiscoverRequest) (*IngressDiscoveryResult, error)
	Connect(ctx context.Context, req IngressConnectRequest) (*IngressConnection, error)
	Update(ctx context.Context, id string, req IngressUpdateRequest) (*IngressConnection, error)
	Retry(ctx context.Context, id string) (*IngressConnection, error)
	Disconnect(ctx context.Context, id string) error
	GetDomain(ctx context.Context, domainID string) (any, error)
	SetDomain(ctx context.Context, domainID string, req IngressDomainRequest) (any, error)
}

type ExampleService interface {
	ListExamples() (any, error)
}

type ServiceLogStreamer interface {
	StreamServiceLogs(ctx context.Context, serviceID string, w io.Writer) error
}

type ScheduledTaskService interface {
	ListTasks(ctx context.Context, projectID, serviceID string) (any, error)
	RunTask(ctx context.Context, taskID string) error
	CancelTask(ctx context.Context, taskID string) error
}
