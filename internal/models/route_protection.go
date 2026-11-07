package models

import "time"

type RouteProtectionPolicy struct {
	Version        string `json:"version"`
	Salt           string `json:"salt"`
	Verifier       string `json:"verifier"`
	SigningKey     string `json:"signingKey"`
	Iterations     int    `json:"iterations"`
	SessionSeconds int    `json:"sessionSeconds"`
}

type RouteProtection struct {
	ID              string    `json:"id" db:"id"`
	ProjectID       string    `json:"projectId" db:"project_id"`
	Hostname        string    `json:"hostname" db:"hostname"`
	Enabled         bool      `json:"enabled" db:"enabled"`
	Secret          *string   `json:"-" db:"secret"`
	Revision        string    `json:"revision" db:"revision"`
	AppliedRevision *string   `json:"appliedRevision,omitempty" db:"applied_revision"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type RouteProtectionStatus struct {
	Hostname string `json:"hostname"`
	Enabled  bool   `json:"enabled"`
	Status   string `json:"status"` // "applied", "pending", "failed"
}

type RouteProtectionInput struct {
	Hostname string `json:"hostname"`
	Action   string `json:"action"` // "protect", "revoke", "retry", "unlock"
	Password string `json:"password,omitempty"`
}

type RouteProtectionDraft struct {
	Hostname string `json:"hostname"`
	Action   string `json:"action"`
	Password string `json:"password"`
}
