package security

import "time"

// RefreshAttempt records the status and metadata of an automated remediation attempt
type RefreshAttempt struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	DeploymentID string    `json:"deploymentId"`
	AdvisoryID   string    `json:"advisoryId"`
	Status       string    `json:"status"` // pending, queued, verified, failed
	Error        string    `json:"error,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ProjectSecurityResult represents the complete security audit findings for a project
type ProjectSecurityResult struct {
	ProjectID            string            `json:"projectId"`
	DeploymentID         string            `json:"deploymentId"`
	CheckedAt            time.Time         `json:"checkedAt"`
	AutoRefresh          bool              `json:"autoRefresh"`
	Findings             []AdvisoryFinding `json:"findings"`
	UnknownServices      []string          `json:"unknownServices,omitempty"`
	ActiveRefreshAttempt *RefreshAttempt   `json:"activeRefreshAttempt,omitempty"`
}

// AdvisoriesResponse defines the API response payload for all platform security advisories
type AdvisoriesResponse struct {
	Revision         int64                   `json:"revision"`
	CheckedAt        time.Time               `json:"checkedAt"`
	ExpiresAt        time.Time               `json:"expiresAt"`
	Stale            bool                    `json:"stale"`
	Advisories       []SecurityAdvisory      `json:"advisories"`
	AffectedProjects []ProjectSecurityResult `json:"affectedProjects"`
}
