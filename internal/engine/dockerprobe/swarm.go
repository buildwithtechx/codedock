package dockerprobe

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"codedock/internal/models"
)

type dockerInfoSwarm struct {
	NodeID           string `json:"NodeID"`
	LocalNodeState   string `json:"LocalNodeState"`
	ControlAvailable bool   `json:"ControlAvailable"`
	Error            string `json:"Error"`
	RemoteManagers   []struct {
		NodeID string `json:"NodeID"`
		Addr   string `json:"Addr"`
	} `json:"RemoteManagers"`
	Nodes    int `json:"Nodes"`
	Managers int `json:"Managers"`
	Cluster  struct {
		ID string `json:"ID"`
	} `json:"Cluster"`
}

// InspectSwarm probes Docker Swarm status via the provided runner.
func InspectSwarm(ctx context.Context, runner Runner) (*models.SwarmInspection, error) {
	out, err := runner.Run(ctx, "docker info --format '{{json .Swarm}}'")
	if err != nil {
		return &models.SwarmInspection{
			State:      models.SwarmStateUnknown,
			ObservedAt: time.Now().UTC(),
		}, nil
	}

	var raw dockerInfoSwarm
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &raw); err != nil {
		return &models.SwarmInspection{
			State:      models.SwarmStateUnknown,
			ObservedAt: time.Now().UTC(),
		}, nil
	}

	state := models.SwarmStateInactive
	switch strings.ToLower(raw.LocalNodeState) {
	case "active":
		state = models.SwarmStateActive
	case "pending":
		state = models.SwarmStatePending
	case "locked":
		state = models.SwarmStateLocked
	case "error":
		state = models.SwarmStateError
	}

	var role *string
	if state == models.SwarmStateActive {
		r := "worker"
		if raw.ControlAvailable {
			r = "manager"
		}
		role = &r
	}

	var nodeID *string
	if raw.NodeID != "" {
		nodeID = &raw.NodeID
	}

	var clusterID *string
	if raw.Cluster.ID != "" {
		clusterID = &raw.Cluster.ID
	}

	var nodeCount, managerCount *int
	if raw.Nodes > 0 {
		n := raw.Nodes
		nodeCount = &n
	}
	if raw.Managers > 0 {
		m := raw.Managers
		managerCount = &m
	}

	return &models.SwarmInspection{
		State:        state,
		Role:         role,
		NodeID:       nodeID,
		ClusterID:    clusterID,
		NodeCount:    nodeCount,
		ManagerCount: managerCount,
		ObservedAt:   time.Now().UTC(),
	}, nil
}
