package models

import "time"

type SwarmState string

const (
	SwarmStateInactive SwarmState = "inactive"
	SwarmStateActive   SwarmState = "active"
	SwarmStatePending  SwarmState = "pending"
	SwarmStateLocked   SwarmState = "locked"
	SwarmStateError    SwarmState = "error"
	SwarmStateUnknown  SwarmState = "unknown"
)

type SwarmNodeSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Availability string `json:"availability"`
	State        string `json:"state"`
	Leader       bool   `json:"leader"`
}

type SwarmPortMapping struct {
	Published int    `json:"published"`
	Target    int    `json:"target"`
	Protocol  string `json:"protocol"`
	Mode      string `json:"mode"`
}

type SwarmVolumeMount struct {
	Type     string `json:"type"`
	Target   string `json:"target"`
	ReadOnly bool   `json:"readOnly"`
}

type SwarmServiceSummary struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	Mode         string             `json:"mode"`
	DesiredTasks *int               `json:"desiredTasks"`
	RunningTasks *int               `json:"runningTasks"`
	UpdateState  *string            `json:"updateState"`
	Ports        []SwarmPortMapping `json:"ports"`
	Mounts       []SwarmVolumeMount `json:"mounts"`
}

type SwarmStackSummary struct {
	Name     *string               `json:"name"`
	Services []SwarmServiceSummary `json:"services"`
}

type SwarmInventory struct {
	Nodes  []SwarmNodeSummary  `json:"nodes"`
	Stacks []SwarmStackSummary `json:"stacks"`
}

type SwarmInspection struct {
	State        SwarmState      `json:"state"`
	Role         *string         `json:"role"`
	NodeID       *string         `json:"nodeId"`
	ClusterID    *string         `json:"clusterId"`
	NodeCount    *int            `json:"nodeCount"`
	ManagerCount *int            `json:"managerCount"`
	ObservedAt   time.Time       `json:"observedAt"`
	Inventory    *SwarmInventory `json:"inventory"`
}
