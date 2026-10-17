package models

import (
	"fmt"
)

type RoleHandoverPhase string

const (
	RoleHandoverPhasePreparing  RoleHandoverPhase = "preparing"
	RoleHandoverPhaseCommitting RoleHandoverPhase = "committing"
	RoleHandoverPhaseKept       RoleHandoverPhase = "kept"
	RoleHandoverPhaseComplete   RoleHandoverPhase = "complete"
)

type RoleCandidate struct {
	Role    string `json:"role"`
	Replica int    `json:"replica"`
	Name    string `json:"name"`
}

type RoleContainerRef struct {
	ContainerID       string `json:"containerId"`
	Role              string `json:"role"`
	Replica           int    `json:"replica"`
	OwnerDeploymentID string `json:"ownerDeploymentId,omitempty"`
}

type PreserveWebRef struct {
	ContainerID       string `json:"containerId"`
	OwnerDeploymentID string `json:"ownerDeploymentId"`
}

type RoleHandoverJournal struct {
	Version              int                `json:"version"`
	Revision             int                `json:"revision"`
	Phase                RoleHandoverPhase  `json:"phase"`
	PreviousDeploymentID *string            `json:"previousDeploymentId"`
	Candidates           []RoleCandidate    `json:"candidates"`
	Stopped              []RoleContainerRef `json:"stopped"`
	WebName              string             `json:"webName"`
	WebContainerID       *string            `json:"webContainerId"`
	StopPreviousWeb      bool               `json:"stopPreviousWeb"`
	PreserveWeb          *PreserveWebRef    `json:"preserveWeb,omitempty"`
	Reused               []RoleContainerRef `json:"reused,omitempty"`
}

func (j *RoleHandoverJournal) IsPending() bool {
	if j == nil {
		return false
	}
	return j.Phase != RoleHandoverPhaseComplete
}

func (j *RoleHandoverJournal) Validate() error {
	if j == nil {
		return nil
	}
	if j.Version < 1 || j.Version > 2 {
		return fmt.Errorf("invalid handover version %d", j.Version)
	}
	if j.Revision < 1 {
		return fmt.Errorf("invalid handover revision %d", j.Revision)
	}
	switch j.Phase {
	case RoleHandoverPhasePreparing, RoleHandoverPhaseCommitting, RoleHandoverPhaseKept, RoleHandoverPhaseComplete:
	default:
		return fmt.Errorf("invalid handover phase %q", j.Phase)
	}
	return nil
}

type RoleSelectorQuery struct {
	Role         string `json:"role,omitempty"`
	Replica      int    `json:"replica,omitempty"`
	DeploymentID string `json:"deploymentId,omitempty"`
}
