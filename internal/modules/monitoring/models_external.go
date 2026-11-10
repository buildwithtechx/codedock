package monitoring

import (
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

type ExampleApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Icon        string `json:"icon,omitempty"`
	Logo        string `json:"logo,omitempty"`
}
