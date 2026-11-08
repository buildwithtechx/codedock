package monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *AttentionManager) scanDeployments(ctx context.Context, _ string, orgID string) ([]*AttentionIssue, error) {
	records, _, err := s.deploys.ListByOrganization(ctx, DeploymentListFilter{OrganizationID: orgID, Limit: 200})
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var issues []*AttentionIssue
	for _, record := range records {
		if record.Status != DeploymentStatusFailed {
			continue
		}
		if record.CreatedAt.Before(time.Now().UTC().Add(-24 * time.Hour)) {
			continue
		}
		if seen[record.ServiceID] {
			continue
		}
		seen[record.ServiceID] = true
		params, _ := json.Marshal(map[string]string{"serviceId": record.ServiceID, "projectId": record.ProjectID})
		issues = append(issues, &AttentionIssue{
			Kind: "deployment-failed", Subject: record.ServiceID,
			ProjectID: record.ProjectID, ServiceID: record.ServiceID,
			Severity:    AttentionSeverityCritical,
			Title:       fmt.Sprintf("Deployments failing for %s", record.ServiceName),
			Detail:      fmt.Sprintf("Latest deployment %s failed at %s", record.ID, record.CreatedAt.Format(time.RFC3339)),
			Remediation: "Inspect the deployment logs, fix the build or config, then redeploy.",
			Action:      "redeploy", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *AttentionManager) scanWorkloads(ctx context.Context, _ string, orgID string) ([]*AttentionIssue, error) {
	if s.health == nil {
		return nil, nil
	}
	projects, _, err := s.projects.ListByOrganization(ctx, orgID, "", 200, 0)
	if err != nil {
		return nil, err
	}
	var issues []*AttentionIssue
	for _, project := range projects {
		services, err := s.apps.ListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		for _, service := range services {
			if service.ContainerID == "" {
				continue
			}
			health, err := s.health.GetHealth(ctx, service.ContainerID)
			if err != nil || health.Running {
				continue
			}
			params, _ := json.Marshal(map[string]string{"serviceId": service.ID})
			issues = append(issues, &AttentionIssue{
				Kind: "service-unhealthy", Subject: service.ID,
				ProjectID: project.ID, ServiceID: service.ID,
				Severity:    AttentionSeverityWarning,
				Title:       fmt.Sprintf("%s is not running", service.Name),
				Detail:      fmt.Sprintf("Container %s reports %s", service.ContainerID, health.Status),
				Remediation: "Check service logs for crash loops, then restart the workload.",
				Action:      "restart-service", ActionParams: string(params),
			})
		}
	}
	return issues, nil
}

func (s *AttentionManager) scanBackups(ctx context.Context, _ string, orgID string) ([]*AttentionIssue, error) {
	if s.backups == nil {
		return nil, nil
	}
	records, err := s.backups.ListRecordsByConfigs(ctx, nil, 100)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var issues []*AttentionIssue
	for _, record := range records {
		if string(record.Status) != "failed" {
			continue
		}
		if seen[record.BackupConfigID] {
			continue
		}
		seen[record.BackupConfigID] = true
		params, _ := json.Marshal(map[string]string{"configId": record.BackupConfigID})
		issues = append(issues, &AttentionIssue{
			Kind: "backup-failed", Subject: record.BackupConfigID,
			Severity:    AttentionSeverityWarning,
			Title:       "Backup run failed",
			Detail:      fmt.Sprintf("Backup %s failed: %s", record.BackupConfigID, record.Logs),
			Remediation: "Verify the database and destination, then retry the backup.",
			Action:      "retry-backup", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *AttentionManager) scanMigrations(ctx context.Context, userID, orgID string) ([]*AttentionIssue, error) {
	if s.migrations == nil {
		return nil, nil
	}
	runs, err := s.migrations.ListRuns(ctx, userID, orgID, "")
	if err != nil {
		return nil, err
	}
	var issues []*AttentionIssue
	for _, run := range runs {
		if run.Status != MigrationStatusFailed {
			continue
		}
		params, _ := json.Marshal(map[string]string{"runId": run.ID})
		issues = append(issues, &AttentionIssue{
			Kind: "migration-failed", Subject: run.ID,
			ProjectID:   run.ProjectID,
			Severity:    AttentionSeverityWarning,
			Title:       "Migration run failed",
			Detail:      fmt.Sprintf("Run %s failed in %s: %s", run.ID, run.Phase, run.Error),
			Remediation: "Review the run logs, then resume with overrides or skips.",
			Action:      "resume-migration", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *AttentionManager) scanQuota(ctx context.Context, _ string, orgID string) ([]*AttentionIssue, error) {
	if s.managed == nil {
		return nil, nil
	}
	quota, err := s.managed.GetQuota(ctx, orgID)
	if err != nil {
		return nil, err
	}
	ids, err := s.managed.ListOrgServerIDs(ctx, orgID)
	if err != nil {
		return nil, err
	}
	usedServers := 0
	usedMemory := 0
	for _, id := range ids {
		server, err := s.servers.GetByID(ctx, id)
		if err != nil || server == nil {
			continue
		}
		usedServers++
		if tier, err := hetzner.TierByName(server.ServerType); err == nil {
			usedMemory += tier.MemoryGB
		}
	}
	var issues []*AttentionIssue
	if quota.MaxServers > 0 && float64(usedServers)/float64(quota.MaxServers) >= 0.8 {
		issues = append(issues, &AttentionIssue{
			Kind: "quota-pressure", Subject: "managed-servers",
			Severity:    AttentionSeverityWarning,
			Title:       "Managed server quota nearly exhausted",
			Detail:      fmt.Sprintf("%d of %d managed servers in use", usedServers, quota.MaxServers),
			Remediation: "Raise the organization quota or delete unused managed servers.",
		})
	}
	if quota.MaxMemoryGB > 0 && float64(usedMemory)/float64(quota.MaxMemoryGB) >= 0.8 {
		issues = append(issues, &AttentionIssue{
			Kind: "quota-pressure", Subject: "managed-memory",
			Severity:    AttentionSeverityWarning,
			Title:       "Managed memory quota nearly exhausted",
			Detail:      fmt.Sprintf("%d of %d GB in use", usedMemory, quota.MaxMemoryGB),
			Remediation: "Raise the organization quota or downsize managed servers.",
		})
	}
	return issues, nil
}
