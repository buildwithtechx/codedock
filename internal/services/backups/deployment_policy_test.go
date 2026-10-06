package backups

import (
	"codedock.run/codedock/internal/models"
	"testing"
)

func TestDeploymentPolicyDoesNotCrossProjectsOrServices(t *testing.T) {
	policy := &models.BackupConfig{ProjectID: "project"}
	if !deploymentPolicyMatches(policy, "project", "") || deploymentPolicyMatches(policy, "foreign", "") {
		t.Fatal("project policy scope lost")
	}
	policy.ServiceID = "application"
	if !deploymentPolicyMatches(policy, "project", "application") || deploymentPolicyMatches(policy, "project", "") || deploymentPolicyMatches(policy, "project", "foreign") {
		t.Fatal("application policy applied to another workload")
	}
}
