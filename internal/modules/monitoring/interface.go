package monitoring

import "context"

// AttentionRepository manages detected infrastructure health warnings and automated remediations.
type AttentionRepository interface {
	Upsert(ctx context.Context, issue *AttentionIssue) error
	Get(ctx context.Context, id string) (*AttentionIssue, error)
	List(ctx context.Context, orgID, status string, limit int) ([]*AttentionIssue, error)
	CountOpen(ctx context.Context, orgID string) (int, error)
	SetStatus(ctx context.Context, id, orgID string, status AttentionStatus) error
	ResolveMissing(ctx context.Context, orgID string, live []string) error
	AppendDetail(ctx context.Context, id, line string) error
	ClearResolved(ctx context.Context, orgID string) (int64, error)
}

// TrafficRepository stores and aggregates HTTP request telemetry and geo-IP metrics.
type TrafficRepository interface {
	RecordBatch(ctx context.Context, samples []TrafficSample) error
	Summary(ctx context.Context, projectID, domain, from, to string) (*AnalyticsSummary, error)
	Overview(ctx context.Context, projectID, domain, from, to, stepMinutes string) (*AnalyticsOverview, error)
	Geo(ctx context.Context, projectID, from, to string) (*AnalyticsGeo, error)
	PathsEnabled(ctx context.Context, projectID string) (bool, error)
	SetPathsEnabled(ctx context.Context, projectID string, enabled bool) error
	RetentionCutoff(days int) string
	DeleteBefore(ctx context.Context, cutoff string) error
}

type BackupRepository interface {
	ListRecords(ctx context.Context, orgID string, limit int) ([]BackupRecord, error)
	GetRecord(ctx context.Context, id string) (*BackupRecord, error)
}

type DeploymentRepository interface {
	ListByOrganization(ctx context.Context, filter DeploymentListFilter) ([]DeploymentRecord, int, error)
}

type ProjectRepository interface {
	GetByOrganization(ctx context.Context, projectID, orgID string) (*Project, error)
	ListByOrganization(ctx context.Context, orgID, search string, limit, offset int) ([]Project, int, error)
}

type AppServiceRepository interface {
	GetByID(ctx context.Context, id string) (*AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]AppService, error)
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByStripeCustomerID(ctx context.Context, customerID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

type OrganizationRepository interface {
	GetMember(ctx context.Context, orgID, userID string) (*OrgMember, error)
}

type ServerRepository interface {
	GetByID(ctx context.Context, id string) (*Server, error)
}

type ManagedRepository interface {
	GetQuota(ctx context.Context, orgID string) (*ManagedQuota, error)
	ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error)
}

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*ServerSettings, error)
	UpdateServerSettings(ctx context.Context, settings *ServerSettings) error
}
