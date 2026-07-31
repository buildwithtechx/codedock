package models

import "time"

type TakeoverPlatform string

const (
	TakeoverPlatformDokploy TakeoverPlatform = "dokploy"
	TakeoverPlatformCoolify TakeoverPlatform = "coolify"
	TakeoverPlatformDokku   TakeoverPlatform = "dokku"
	TakeoverPlatformDocker  TakeoverPlatform = "docker"
)

type TakeoverStatus string

const (
	TakeoverStatusScanning TakeoverStatus = "scanning"
	TakeoverStatusScanned  TakeoverStatus = "scanned"
	TakeoverStatusAdopting TakeoverStatus = "adopting"
	TakeoverStatusDone     TakeoverStatus = "done"
	TakeoverStatusFailed   TakeoverStatus = "failed"
)

type TakeoverRun struct {
	ID                string           `json:"id"`
	UserID            string           `json:"userId"`
	SourceHost        string           `json:"sourceHost"`
	SourcePlatform    TakeoverPlatform `json:"sourcePlatform"`
	Status            TakeoverStatus   `json:"status"`
	DiscoveredJSON    string           `json:"-"`
	AdoptedProjectIDs string           `json:"-"`
	Error             string           `json:"error,omitempty"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

type DiscoveredContainer struct {
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	Ports          []string          `json:"ports"`
	Env            map[string]string `json:"env"`
	Volumes        []string          `json:"volumes"`
	Labels         map[string]string `json:"labels"`
	Status         string            `json:"status"`
	ComposeProject string            `json:"composeProject,omitempty"`
}

type DiscoveredStack struct {
	Containers      []DiscoveredContainer `json:"containers"`
	ComposeProjects []string              `json:"composeProjects"`
	Platform        TakeoverPlatform      `json:"platform"`
	Host            string                `json:"host"`
}

type TakeoverScanRequest struct {
	Host           string           `json:"host"`
	SSHUser        string           `json:"sshUser"`
	SSHKey         string           `json:"sshKey"`
	SSHFingerprint string           `json:"sshFingerprint,omitempty"`
	Platform       TakeoverPlatform `json:"platform"`
}

type TakeoverAdoptRequest struct {
	RunID        string   `json:"runId"`
	ProjectName  string   `json:"projectName"`
	ServiceNames []string `json:"serviceNames"`
	ImportEnv    bool     `json:"importEnv"`
}
