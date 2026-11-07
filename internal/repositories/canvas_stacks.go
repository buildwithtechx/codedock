package repositories

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *CanvasRepo) SetVault(vault Vault) { r.vault = vault }

func (r *CanvasRepo) projectStacks(ctx context.Context, canvas *models.EnvironmentCanvas) error {
	var stacks []models.ComposeStack
	if err := r.db.SelectContext(ctx, &stacks, `SELECT * FROM compose_stacks WHERE environment_id=$1`, canvas.Environment.ID); err != nil {
		return err
	}
	for _, stack := range stacks {
		if r.vault == nil {
			return fmt.Errorf("stack configuration unavailable")
		}
		config, err := r.vault.Decrypt(stack.Config)
		if err != nil {
			return err
		}
		var document struct {
			Services map[string]struct {
				DependsOn map[string]any `json:"depends_on"`
			} `json:"services"`
		}
		if err := json.Unmarshal([]byte(config), &document); err != nil {
			return err
		}
		var results []models.ComposeServiceResult
		if err := json.Unmarshal([]byte(stack.Results), &results); err != nil {
			return err
		}
		index := 0
		for name, service := range document.Services {
			id := "stack-" + stack.ID + "-" + name
			state, containerID := "unobserved", ""
			for _, result := range results {
				if result.Name == name {
					state, containerID = result.State, result.ContainerID
				}
			}
			canvas.Nodes = append(canvas.Nodes, models.CanvasNode{ID: id, Type: "default", Data: map[string]any{"label": stack.Name + " / " + name, "name": name, "status": state, "stackId": stack.ID, "projectId": stack.ProjectID, "containerId": containerID}, Pos: models.CanvasPosition{X: 100 + float64(index)*240, Y: 430}})
			index++
			for dependency := range service.DependsOn {
				canvas.Edges = append(canvas.Edges, models.CanvasEdge{ID: "stack-dependency:" + stack.ID + ":" + dependency + ":" + name, Source: "stack-" + stack.ID + "-" + dependency, Target: id, Kind: "dependency", Label: "Compose prerequisite"})
			}
		}
	}
	return nil
}
