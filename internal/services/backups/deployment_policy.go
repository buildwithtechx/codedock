package backups

import "codedock.run/codedock/internal/models"

func deploymentPolicyMatches(cfg *models.BackupConfig, projectID, serviceID string) bool {
	if cfg.ServiceID != "" {
		return serviceID != "" && cfg.ServiceID == serviceID
	}
	return projectID != "" && cfg.ProjectID == projectID
}
