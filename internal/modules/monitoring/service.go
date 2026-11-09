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

type MetricsOperations interface {
	GetServiceMetrics(ctx context.Context, serviceID, rangeParam string) (any, error)
}

type LogOperations interface {
	GetHistoricalLogs(ctx context.Context, serviceID, rangeParam string, limit int) (any, error)
}

type SystemOperations interface {
	GetStats(ctx context.Context) (any, error)
	Restart(ctx context.Context) error
	Cleanup(ctx context.Context) error
}

type UpdaterOperations interface {
	GetStatus() any
	CheckUpdate(ctx context.Context) (any, error)
	ApplyUpdate(ctx context.Context) (any, error)
}

type OperationsService interface {
	Get(ctx context.Context, id string) (*Operation, error)
	List(ctx context.Context, project, user string) ([]Operation, error)
	Cancel(ctx context.Context, id, user string) error
}

type ExampleOperations interface {
	ListExamples() (any, error)
}

type ServiceLogStreamer interface {
	StreamServiceLogs(ctx context.Context, serviceID string, w io.Writer) error
}
