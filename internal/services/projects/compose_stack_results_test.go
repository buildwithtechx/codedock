package projects

import (
	"codedock/internal/models"
	"testing"
)

func TestStackReadinessRequiresEveryConfiguredService(t *testing.T) {
	config := `{"services":{"api":{},"worker":{}}}`
	results := []models.ComposeServiceResult{{Name: "api", State: "running", Health: "healthy"}}
	if validateStackResults(config, results) == nil {
		t.Fatal("missing worker reported ready")
	}
	results = append(results, models.ComposeServiceResult{Name: "worker", State: "exited", ExitCode: 1})
	if validateStackResults(config, results) == nil {
		t.Fatal("failed worker reported ready")
	}
	results[1].State = "running"
	if err := validateStackResults(config, results); err != nil {
		t.Fatal(err)
	}
}
