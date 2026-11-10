package tasks

// ScheduledTaskResponse wraps a single scheduled task entity.
type ScheduledTaskResponse struct {
	Task *ScheduledTask `json:"task"`
}

// ScheduledTaskListResponse wraps a list of scheduled tasks.
type ScheduledTaskListResponse struct {
	Tasks []ScheduledTask `json:"tasks"`
	Total int             `json:"total"`
}

// ExecuteTaskResponse wraps the stdout/stderr output from a task run.
type ExecuteTaskResponse struct {
	Output string `json:"output"`
}
