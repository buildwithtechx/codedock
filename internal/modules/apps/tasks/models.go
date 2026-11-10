package tasks

import "time"

type ScheduledTask struct {
	ID         string    `json:"id" db:"id"`
	ServiceID  string    `json:"serviceId" db:"service_id"`
	Name       string    `json:"name" db:"name"`
	Schedule   string    `json:"schedule" db:"schedule"`
	Command    string    `json:"command" db:"command"`
	Status     string    `json:"status" db:"status"`
	LastRunAt  time.Time `json:"lastRunAt" db:"last_run_at"`
	LastOutput string    `json:"lastOutput" db:"last_output"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}
