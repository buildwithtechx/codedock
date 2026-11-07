package models

type BuildAdmissionStatus struct {
	CanAdmit             bool   `json:"canAdmit"`
	ActiveBuilds         int    `json:"activeBuilds"`
	MaxConcurrentBuilds  int    `json:"maxConcurrentBuilds"`
	QueuePosition        int    `json:"queuePosition"`
	EstimatedWaitSeconds int    `json:"estimatedWaitSeconds"`
	Reason               string `json:"reason"`
}

type ProjectChecksConfig struct {
	ProjectID            string   `json:"projectId"`
	Enabled              bool     `json:"enabled"`
	ReportCommitStatus   bool     `json:"reportCommitStatus"`
	RequirePassingChecks bool     `json:"requirePassingChecks"`
	TargetBranches       []string `json:"targetBranches"`
	ContextName          string   `json:"contextName"`
	UpdatedAt            string   `json:"updatedAt"`
}

type UpdateChecksConfigRequest struct {
	Enabled              bool     `json:"enabled"`
	ReportCommitStatus   bool     `json:"reportCommitStatus"`
	RequirePassingChecks bool     `json:"requirePassingChecks"`
	TargetBranches       []string `json:"targetBranches"`
	ContextName          string   `json:"contextName"`
}
