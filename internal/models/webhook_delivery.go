package models

import (
	"encoding/json"
	"time"
)

type WebhookDelivery struct {
	ID             string          `json:"id" db:"id"`
	OrganizationID *string         `json:"organizationId,omitempty" db:"organization_id"`
	ProjectID      *string         `json:"projectId,omitempty" db:"project_id"`
	Source         string          `json:"source" db:"source"`
	HookID         *string         `json:"hookId,omitempty" db:"hook_id"`
	DeliveryID     *string         `json:"deliveryId,omitempty" db:"delivery_id"`
	Event          string          `json:"event" db:"event"`
	AuthResult     *string         `json:"authResult,omitempty" db:"auth_result"`
	Outcome        string          `json:"outcome" db:"outcome"`
	ActionRef      *string         `json:"actionRef,omitempty" db:"action_ref"`
	StatusCode     *int            `json:"statusCode,omitempty" db:"status_code"`
	Error          *string         `json:"error,omitempty" db:"error"`
	Summary        json.RawMessage `json:"summary,omitempty" db:"summary"`
	ClientIP       *string         `json:"clientIp,omitempty" db:"client_ip"`
	UserAgent      *string         `json:"userAgent,omitempty" db:"user_agent"`
	DurationMs     *int            `json:"durationMs,omitempty" db:"duration_ms"`
	ReceivedAt     time.Time       `json:"receivedAt" db:"received_at"`
	ProcessedAt    *time.Time      `json:"processedAt,omitempty" db:"processed_at"`
}

type WebhookDeliveryInput struct {
	OrganizationID *string
	ProjectID      *string
	Source         string
	HookID         *string
	DeliveryID     *string
	Event          string
	AuthResult     *string
	Outcome        string
	ActionRef      *string
	StatusCode     *int
	Error          *string
	Summary        any
	ClientIP       *string
	UserAgent      *string
	DurationMs     *int
}

type WebhookDeliveryResult struct {
	Outcome    string
	ActionRef  *string
	StatusCode *int
	Error      *string
	Summary    any
	DurationMs *int
}

type ClaimResult struct {
	Claimed           bool
	ID                string
	HandledProjectIDs []string
}
