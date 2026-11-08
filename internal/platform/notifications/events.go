package notifications

import "time"

type Event interface {
	Name() string
}

type DeploymentCompleted struct {
	ProjectID  string
	ServiceID  string
	Status     string
	CommitHash string
	Branch     string
	DeployTime time.Duration
}

func (e DeploymentCompleted) Name() string { return "deployment.completed" }
