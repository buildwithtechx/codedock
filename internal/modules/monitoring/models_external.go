package monitoring

import (
	"context"
	"time"
)

type DeploymentListFilter struct {
	OrganizationID string
	ProjectID      string
	ServiceID      string
	Limit          int
	Offset         int
}

type DeploymentRecord struct {
	ID          string
	ProjectID   string
	ServiceID   string
	ServiceName string
	Status      string
	CreatedAt   time.Time
	FinishedAt  *time.Time
}

type Project struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
}

type AppService struct {
	ID                   string            `json:"id"`
	ProjectID            string            `json:"projectId"`
	Name                 string            `json:"name"`
	Type                 string            `json:"type"`
	Status               string            `json:"status"`
	Replicas             int               `json:"replicas"`
	ContainerID          string            `json:"containerId,omitempty"`
	EnvironmentVariables map[string]string `json:"environmentVariables,omitempty"`
}

const MigrationStatusFailed = "failed"

type MigrationRun struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
	Phase     string `json:"phase"`
	Error     string `json:"error"`
}

type ResumeMigrationRequest struct{}

type MigrationResumeResult struct {
	RunID string `json:"runId"`
}

type User struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	Name                 string  `json:"name"`
	Role                 string  `json:"role"`
	PlanType             string  `json:"planType"`
	StripeCustomerID     *string `json:"stripeCustomerId"`
	StripeSubscriptionID *string `json:"stripeSubscriptionId"`
	StripePriceID        *string `json:"stripePriceId"`
}

type MemberPermission string

const (
	MemberPermissionOwner MemberPermission = "owner"
	MemberPermissionAdmin MemberPermission = "admin"
)

type OrgMember struct {
	OrganizationID string           `json:"organizationId"`
	UserID         string           `json:"userId"`
	Role           string           `json:"role"`
	Permission     MemberPermission `json:"permission"`
}

type BackupRecord struct {
	ID             string
	BackupConfigID string
	Status         string
	Logs           string
	CreatedAt      time.Time
}

type Server struct {
	ID         string `json:"id"`
	ServerType string `json:"serverType"`
}

type ManagedQuota struct {
	MaxServers  int `json:"maxServers"`
	MaxMemoryGB int `json:"maxMemoryGb"`
}

type ServerSettings struct {
	AutoUpdateEnabled bool   `json:"autoUpdateEnabled"`
	ReleaseChannel    string `json:"releaseChannel"`
	CurrentVersion    string `json:"currentVersion"`
	LatestVersion     string `json:"latestVersion"`
	LastUpdateCheck   string `json:"lastUpdateCheck,omitempty"`
	UpdateCheckCron   string `json:"updateCheckCron,omitempty"`
}

type DeviceToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExampleApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Icon        string `json:"icon,omitempty"`
	Logo        string `json:"logo,omitempty"`
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

type DeviceTokenRepository interface {
	RegisterToken(ctx context.Context, userID, token, platform string) (*DeviceToken, error)
	DeleteToken(ctx context.Context, userID, token string) error
	ListTokensByUser(ctx context.Context, userID string) ([]*DeviceToken, error)
	ListAllTokens(ctx context.Context) ([]*DeviceToken, error)
}
