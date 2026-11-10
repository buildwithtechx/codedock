package webhooks

import "time"

type DrainType string

const (
	DrainTypeHTTP   DrainType = "http"
	DrainTypeSyslog DrainType = "syslog"
)

type Webhook struct {
	ID        string    `json:"id" db:"id"`
	ServiceID string    `json:"serviceId" db:"service_id"`
	URL       string    `json:"url" db:"url"`
	Secret    string    `json:"secret,omitempty" db:"secret"`
	Events    string    `json:"events" db:"events"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type LogDrain struct {
	ID          string    `json:"id" db:"id"`
	ProjectID   string    `json:"projectId" db:"project_id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	DrainType   DrainType `json:"drainType" db:"drain_type"`
	EndpointURL string    `json:"endpointUrl" db:"endpoint_url"`
	APIKey      string    `json:"apiKey,omitempty" db:"api_key"`
	AuthToken   string    `json:"authToken,omitempty" db:"auth_token"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
