package attention

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"codedock.run/codedock/internal/models"
)

func (s *Service) Act(ctx context.Context, userID, orgID, issueID string, req models.AttentionActionRequest) (string, error) {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return "", err
	}
	issue, err := s.issues.Get(ctx, issueID)
	if err != nil {
		return "", err
	}
	if issue == nil || issue.OrganizationID != orgID {
		return "", fmt.Errorf("attention issue not found")
	}
	if issue.Status == models.AttentionStatusResolved {
		return "", fmt.Errorf("issue is already resolved")
	}
	action := req.Action
	if action == "" {
		action = issue.Action
	}
	if action == "" || action != issue.Action {
		return "", fmt.Errorf("unknown remediation action")
	}
	params := map[string]string{}
	if issue.ActionParams != "" {
		_ = json.Unmarshal([]byte(issue.ActionParams), &params)
	}
	for key, value := range req.Params {
		params[key] = value
	}
	var outcome string
	switch action {
	case "restart-service":
		outcome, err = s.actRestart(ctx, params)
	case "redeploy":
		outcome, err = s.actRedeploy(ctx, params)
	case "resume-migration":
		outcome, err = s.actResumeMigration(ctx, userID, params)
	case "retry-backup":
		outcome, err = s.actRetryBackup(ctx, params)
	default:
		return "", fmt.Errorf("unknown remediation action")
	}
	stamp := time.Now().UTC().Format(time.RFC3339)
	if err != nil {
		_ = s.issues.AppendDetail(ctx, issueID, fmt.Sprintf("%s action %s failed: %v", stamp, action, err))
		return "", err
	}
	_ = s.issues.AppendDetail(ctx, issueID, fmt.Sprintf("%s action %s by %s: %s", stamp, action, userID, outcome))
	return outcome, nil
}

func (s *Service) actRestart(ctx context.Context, params map[string]string) (string, error) {
	if s.deployer == nil {
		return "", fmt.Errorf("deployer is not available")
	}
	app, err := s.apps.GetByID(ctx, params["serviceId"])
	if err != nil {
		return "", err
	}
	if app == nil {
		return "", fmt.Errorf("service not found")
	}
	if err := s.deployer.RestartAppService(ctx, app); err != nil {
		return "", err
	}
	return fmt.Sprintf("restarted %s", app.Name), nil
}

func (s *Service) actRedeploy(ctx context.Context, params map[string]string) (string, error) {
	if s.deployer == nil {
		return "", fmt.Errorf("deployer is not available")
	}
	app, err := s.apps.GetByID(ctx, params["serviceId"])
	if err != nil {
		return "", err
	}
	if app == nil {
		return "", fmt.Errorf("service not found")
	}
	var logs bytes.Buffer
	containerID, err := s.deployer.DeployAppService(ctx, app.ID, "", &logs)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("redeployed %s as %s", app.Name, containerID), nil
}

func (s *Service) actResumeMigration(ctx context.Context, userID string, params map[string]string) (string, error) {
	if s.migrations == nil {
		return "", fmt.Errorf("migrations are not available")
	}
	result, err := s.migrations.ResumeRun(ctx, userID, params["runId"], models.ResumeMigrationRequest{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("resumed migration %s", result.RunID), nil
}

func (s *Service) actRetryBackup(ctx context.Context, params map[string]string) (string, error) {
	if s.backups == nil {
		return "", fmt.Errorf("backups are not available")
	}
	record, err := s.backups.TriggerBackup(ctx, params["configId"])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("backup %s running", record.ID), nil
}
