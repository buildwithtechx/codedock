package runtimes

import (
	"codedock.run/codedock/internal/models"
	"encoding/json"
	"testing"
)

func TestDesiredReplicasPreservesStoppedWorkloadsAndRejectsForeignOwnership(t *testing.T) {
	manifest := `{"items":[{"kind":"Deployment","metadata":{"labels":{"codedock.run/service":"service","codedock.run/project":"project"}},"spec":{"replicas":3}}]}`
	app := &models.AppService{ID: "service", ProjectID: "project", Replicas: 4, Status: models.AppServiceStatusStopped}
	raw, err := desiredReplicas(manifest, app)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Items []struct {
			Spec struct {
				Replicas int `json:"replicas"`
			} `json:"spec"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Items[0].Spec.Replicas != 0 {
		t.Fatal("reconciliation started a stopped service")
	}
	app.ProjectID = "foreign"
	if _, err := desiredReplicas(manifest, app); err == nil {
		t.Fatal("foreign desired workload accepted")
	}
}
