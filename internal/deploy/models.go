package deploy

import (
	appsmod "codedock/internal/modules/apps"
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
	return j.Phase == RoleHandoverPhasePreparing || j.Phase == RoleHandoverPhaseCommitting || j.Phase == RoleHandoverPhaseKept
}

type RolloutContainer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Running bool   `json:"running"`
}

type RolloutJournal struct {
	ID          string             `json:"id"`
	PreviousApp appsmod.AppService `json:"previousApp"`
	Previous    []RolloutContainer `json:"previous"`
	Committed   bool               `json:"committed"`
}
