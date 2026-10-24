package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

type StorageInventory struct {
	StorageClasses []map[string]any `json:"storageClasses"`
	Claims         []map[string]any `json:"claims"`
	Snapshots      []map[string]any `json:"snapshots"`
	Provisioner    string           `json:"provisioner"`
	DefaultClass   string           `json:"defaultClass"`
}

func kubectlList(ctx context.Context, commands Commands, node models.ClusterNode, args []string) ([]map[string]any, error) {
	response, err := commands.Kubectl(ctx, node, append(args, "-o", "json"), "")
	if err != nil {
		return nil, err
	}
	var list map[string]any
	if err := json.Unmarshal([]byte(response), &list); err != nil {
		return nil, err
	}
	items, _ := list["items"].([]any)
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		object, _ := item.(map[string]any)
		result = append(result, object)
	}
	return result, nil
}

func (r *WorkloadRuntime) StorageInventory(ctx context.Context, node models.ClusterNode) (*StorageInventory, error) {
	classes, err := kubectlList(ctx, r.commands, node, []string{"get", "storageclass"})
	if err != nil {
		return nil, err
	}
	claims, err := kubectlList(ctx, r.commands, node, []string{"get", "pvc", "-A"})
	if err != nil {
		return nil, err
	}
	snapshots, err := kubectlList(ctx, r.commands, node, []string{"get", "volumesnapshot", "-A"})
	if err != nil {
		snapshots = []map[string]any{}
	}
	inventory := &StorageInventory{Claims: claims, Snapshots: snapshots}
	for _, class := range classes {
		inventory.StorageClasses = append(inventory.StorageClasses, class)
		metadata, _ := class["metadata"].(map[string]any)
		annotations, _ := metadata["annotations"].(map[string]any)
		if annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			if name, _ := metadata["name"].(string); name != "" {
				inventory.DefaultClass = name
			}
		}
		if provisioner, _ := class["provisioner"].(string); provisioner == "driver.longhorn.io" {
			inventory.Provisioner = provisioner
		}
	}
	return inventory, nil
}

func (r *WorkloadRuntime) StorageHealth(ctx context.Context, node models.ClusterNode) error {
	inventory, err := r.StorageInventory(ctx, node)
	if err != nil {
		return err
	}
	if inventory.Provisioner == "" {
		return fmt.Errorf("shared storage provisioner is not installed; declare storage classes explicitly per workload")
	}
	unbound := 0
	for _, claim := range inventory.Claims {
		status, _ := claim["status"].(map[string]any)
		if status["phase"] != "Bound" {
			unbound++
		}
	}
	if unbound > 0 {
		return fmt.Errorf("%d persistent claims are not bound", unbound)
	}
	return nil
}
