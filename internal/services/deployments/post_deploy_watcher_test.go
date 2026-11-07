package deployments

import (
	"testing"

	"codedock/internal/models"
)

func TestCheckDeploymentRestartLoop(t *testing.T) {
	if got := CheckDeploymentRestartLoop(nil); got != "" {
		t.Fatalf("expected empty for nil deployment, got %q", got)
	}

	dep := &models.Deployment{
		ID:          "dep-1",
		RestartLoop: "Container in restart loop (exit code 1)",
	}
	if got := CheckDeploymentRestartLoop(dep); got != "Container in restart loop (exit code 1)" {
		t.Fatalf("expected restart loop detail, got %q", got)
	}

	cleanDep := &models.Deployment{
		ID: "dep-2",
	}
	if got := CheckDeploymentRestartLoop(cleanDep); got != "" {
		t.Fatalf("expected empty for clean deployment, got %q", got)
	}
}
