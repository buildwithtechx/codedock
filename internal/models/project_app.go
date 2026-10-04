package models

import "time"

type ProjectApp struct {
	ID               string     `json:"id" db:"id"`
	OrganizationID   string     `json:"organizationId" db:"organization_id"`
	Name             string     `json:"name" db:"name"`
	Slug             string     `json:"slug" db:"slug"`
	GitProvider      string     `json:"gitProvider" db:"git_provider"`
	GitOwner         string     `json:"gitOwner" db:"git_owner"`
	GitRepo          string     `json:"gitRepo" db:"git_repo"`
	GitURL           string     `json:"gitUrl" db:"git_url"`
	InstallationID   int        `json:"installationId" db:"installation_id"`
	Favicon          string     `json:"favicon" db:"favicon"`
	FaviconCheckedAt *time.Time `json:"faviconCheckedAt,omitempty" db:"favicon_checked_at"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`
}

type CreateProjectAppRequest struct {
	Name           string `json:"name"`
	Slug           string `json:"slug,omitempty"`
	GitProvider    string `json:"gitProvider,omitempty"`
	GitOwner       string `json:"gitOwner,omitempty"`
	GitRepo        string `json:"gitRepo,omitempty"`
	GitURL         string `json:"gitUrl,omitempty"`
	InstallationID int    `json:"installationId,omitempty"`
}

type UpdateProjectAppRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	GitProvider string `json:"gitProvider,omitempty"`
	GitOwner    string `json:"gitOwner,omitempty"`
	GitRepo     string `json:"gitRepo,omitempty"`
	GitURL      string `json:"gitUrl,omitempty"`
}
