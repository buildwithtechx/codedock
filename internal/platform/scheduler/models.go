package scheduler

import "time"

type ScheduledTaskStatus string

const (
	ScheduledTaskStatusActive ScheduledTaskStatus = "active"
	ScheduledTaskStatusError  ScheduledTaskStatus = "error"
)

type ScheduledTask struct {
	ID         string              `json:"id"`
	ServiceID  string              `json:"serviceId"`
	Name       string              `json:"name"`
	Schedule   string              `json:"schedule"`
	Command    string              `json:"command"`
	Status     ScheduledTaskStatus `json:"status"`
	LastRunAt  *time.Time          `json:"lastRunAt,omitempty"`
	LastOutput string              `json:"lastOutput,omitempty"`
}

type ProjectConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AppService struct {
	ID          string `json:"id"`
	ContainerID string `json:"containerId"`
	Name        string `json:"name"`
}
