package tasks

type ScheduledTaskStatus string

const (
	ScheduledTaskStatusActive    ScheduledTaskStatus = "active"
	ScheduledTaskStatusInactive  ScheduledTaskStatus = "inactive"
	ScheduledTaskStatusRunning   ScheduledTaskStatus = "running"
	ScheduledTaskStatusCompleted ScheduledTaskStatus = "completed"
	ScheduledTaskStatusFailed    ScheduledTaskStatus = "failed"
	ScheduledTaskStatusError     ScheduledTaskStatus = "error"
)
