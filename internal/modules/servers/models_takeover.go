package servers

import "time"

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
