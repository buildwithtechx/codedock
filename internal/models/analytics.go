package models

type TrafficBucket struct {
	Minute    string `json:"minute" db:"bucket_minute"`
	ProjectID string `json:"projectId" db:"project_id"`
	Domain    string `json:"domain" db:"domain"`
	Path      string `json:"path" db:"path"`
	Status    int    `json:"status" db:"status"`
	Requests  int    `json:"requests" db:"requests"`
	Bytes     int64  `json:"bytes" db:"bytes"`
	DurationMs int64 `json:"durationMs" db:"duration_ms"`
}

type TrafficVisitor struct {
	Day       string `json:"day" db:"day"`
	ProjectID string `json:"projectId" db:"project_id"`
	IP        string `json:"ip" db:"ip"`
	Country   string `json:"country" db:"country"`
	Requests  int    `json:"requests" db:"requests"`
	Bytes     int64  `json:"bytes" db:"bytes"`
}

type TrafficPathsToggle struct {
	ProjectID string `json:"projectId" db:"project_id"`
	Enabled   bool   `json:"enabled" db:"enabled"`
}

type TrafficSample struct {
	Time     string `json:"time"`
	ProjectID string `json:"projectId"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Bytes    int64  `json:"bytes"`
	DurationMs int64 `json:"durationMs"`
	ClientIP string `json:"clientIp"`
}

type AnalyticsSummary struct {
	ProjectID   string  `json:"projectId"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Requests    int     `json:"requests"`
	Bytes       int64   `json:"bytes"`
	ErrorRate   float64 `json:"errorRate"`
	AvgDurationMs float64 `json:"avgDurationMs"`
}

type AnalyticsSeriesPoint struct {
	Time     string `json:"time"`
	Requests int    `json:"requests"`
	Bytes    int64  `json:"bytes"`
	Errors   int    `json:"errors"`
}

type AnalyticsStatusBreakdown struct {
	Status   int `json:"status"`
	Requests int `json:"requests"`
}

type AnalyticsTopPath struct {
	Path     string `json:"path"`
	Requests int    `json:"requests"`
	Bytes    int64  `json:"bytes"`
}

type AnalyticsOverview struct {
	ProjectID string                   `json:"projectId"`
	From      string                   `json:"from"`
	To        string                   `json:"to"`
	Series    []AnalyticsSeriesPoint   `json:"series"`
	Statuses  []AnalyticsStatusBreakdown `json:"statuses"`
	TopPaths  []AnalyticsTopPath       `json:"topPaths"`
}

type AnalyticsGeoEntry struct {
	Country  string `json:"country"`
	Requests int    `json:"requests"`
	Visitors int    `json:"visitors"`
	Bytes    int64  `json:"bytes"`
}

type AnalyticsGeo struct {
	ProjectID string             `json:"projectId"`
	From      string             `json:"from"`
	To        string             `json:"to"`
	Countries []AnalyticsGeoEntry `json:"countries"`
}

type AnalyticsDeploymentStats struct {
	ProjectID   string  `json:"projectId"`
	Total       int     `json:"total"`
	Succeeded   int     `json:"succeeded"`
	Failed      int     `json:"failed"`
	SuccessRate float64 `json:"successRate"`
	AvgDurationSeconds float64 `json:"avgDurationSeconds"`
}

type AnalyticsDashboard struct {
	OrganizationID string  `json:"organizationId"`
	Requests24h    int     `json:"requests24h"`
	ErrorRate24h   float64 `json:"errorRate24h"`
	Deployments7d  int     `json:"deployments7d"`
	DeploySuccess  float64 `json:"deploySuccess"`
	OpenIssues     int     `json:"openIssues"`
}

type AttentionSeverity string

const (
	AttentionSeverityInfo     AttentionSeverity = "info"
	AttentionSeverityWarning  AttentionSeverity = "warning"
	AttentionSeverityCritical AttentionSeverity = "critical"
)

type AttentionStatus string

const (
	AttentionStatusOpen     AttentionStatus = "open"
	AttentionStatusAcked    AttentionStatus = "acked"
	AttentionStatusResolved AttentionStatus = "resolved"
)

type AttentionIssue struct {
	ID             string            `json:"id" db:"id"`
	OrganizationID string            `json:"organizationId" db:"organization_id"`
	Kind           string            `json:"kind" db:"kind"`
	Subject        string            `json:"subject" db:"subject"`
	ProjectID      string            `json:"projectId,omitempty" db:"project_id"`
	ServiceID      string            `json:"serviceId,omitempty" db:"service_id"`
	Severity       AttentionSeverity `json:"severity" db:"severity"`
	Status         AttentionStatus   `json:"status" db:"status"`
	Title          string            `json:"title" db:"title"`
	Detail         string            `json:"detail" db:"detail"`
	Remediation    string            `json:"remediation" db:"remediation"`
	Action         string            `json:"action,omitempty" db:"action"`
	ActionParams   string            `json:"-" db:"action_params"`
	Occurrences    int               `json:"occurrences" db:"occurrences"`
	FirstSeen      string            `json:"firstSeen" db:"first_seen"`
	LastSeen       string            `json:"lastSeen" db:"last_seen"`
	UpdatedAt      string            `json:"updatedAt" db:"updated_at"`
}

type AttentionActionRequest struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}
