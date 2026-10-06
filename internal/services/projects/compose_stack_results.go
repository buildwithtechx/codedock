package projects

import (
	"codedock.run/codedock/internal/models"
	"encoding/json"
	"fmt"
)

func validateStackResults(config string, results []models.ComposeServiceResult) error {
	var document struct {
		Services map[string]json.RawMessage `json:"services"`
	}
	if err := json.Unmarshal([]byte(config), &document); err != nil {
		return fmt.Errorf("read expected stack services: %w", err)
	}
	seen := map[string]bool{}
	for _, result := range results {
		seen[result.Name] = true
		if result.Health == "unhealthy" || result.Health == "starting" || (result.State != "running" && !(result.State == "exited" && result.ExitCode == 0)) {
			return fmt.Errorf("service %s is not ready: %s, health %s, exit %d", result.Name, result.State, result.Health, result.ExitCode)
		}
	}
	for name := range document.Services {
		if !seen[name] {
			return fmt.Errorf("service %s has no observed container", name)
		}
	}
	return nil
}
