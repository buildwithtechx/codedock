package adapters

import (
	"context"
	"strconv"

	"codedock/internal/modules/monitoring"
)

type AnalyticsServiceAdapter struct {
	Svc *monitoring.Service
}

func (a *AnalyticsServiceAdapter) Periods(ctx context.Context, orgID string) (any, error) {
	if a.Svc == nil {
		return []string{}, nil
	}
	return a.Svc.Periods(), nil
}

func (a *AnalyticsServiceAdapter) Dashboard(ctx context.Context, orgID string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return a.Svc.Dashboard(ctx, "", orgID)
}

func (a *AnalyticsServiceAdapter) Summary(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return a.Svc.Summary(ctx, userID, orgID, projectID, "", from, to)
}

func (a *AnalyticsServiceAdapter) Overview(ctx context.Context, userID, orgID, projectID, domain, from, to, period string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return a.Svc.Overview(ctx, userID, orgID, projectID, domain, from, to, period)
}

func (a *AnalyticsServiceAdapter) Geo(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return a.Svc.Geo(ctx, userID, orgID, projectID, from, to)
}

func (a *AnalyticsServiceAdapter) GetPaths(ctx context.Context, userID, orgID, projectID string) (bool, error) {
	if a.Svc == nil {
		return false, nil
	}
	return a.Svc.PathsEnabled(ctx, userID, orgID, projectID)
}

func (a *AnalyticsServiceAdapter) SetPaths(ctx context.Context, userID, orgID, projectID string, enabled bool) error {
	if a.Svc == nil {
		return nil
	}
	return a.Svc.SetPathsEnabled(ctx, userID, orgID, projectID, enabled)
}

func (a *AnalyticsServiceAdapter) DeploymentStats(ctx context.Context, userID, orgID, projectID string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return a.Svc.DeploymentStats(ctx, userID, orgID, projectID, 30)
}

func (a *AnalyticsServiceAdapter) Usage(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.Svc == nil || a.Svc.Usage() == nil {
		return nil, nil
	}
	return a.Svc.Usage().Current(ctx, projectID)
}

func (a *AnalyticsServiceAdapter) Resources(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.Svc == nil || a.Svc.Usage() == nil {
		return nil, nil
	}
	return a.Svc.Usage().Resources(ctx, projectID)
}

func (a *AnalyticsServiceAdapter) UsageHistory(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.Svc == nil || a.Svc.Usage() == nil {
		return nil, nil
	}
	hours, _ := strconv.Atoi(period)
	if hours <= 0 {
		hours = 24
	}
	return a.Svc.Usage().History(ctx, projectID, hours)
}

func (a *AnalyticsServiceAdapter) Container(ctx context.Context, userID, orgID, projectID, serviceID, period string) (any, error) {
	if a.Svc == nil || a.Svc.Usage() == nil {
		return nil, nil
	}
	return a.Svc.Usage().Container(ctx, projectID, serviceID)
}

func (a *AnalyticsServiceAdapter) CurrentUsage(ctx context.Context, projectID string) (any, error) {
	if a.Svc == nil || a.Svc.Usage() == nil {
		return nil, nil
	}
	return a.Svc.Usage().Current(ctx, projectID)
}

type AttentionServiceAdapter struct {
	Svc *monitoring.AttentionManager
}

func (a *AttentionServiceAdapter) List(ctx context.Context, userID, orgID string) (any, error) {
	if a.Svc == nil {
		return []any{}, nil
	}
	issues, err := a.Svc.ListIssues(ctx, userID, orgID, "")
	if err != nil {
		return nil, err
	}
	canClear, _ := a.Svc.CanClearHistory(ctx, userID, orgID)
	return map[string]any{"issues": issues, "canClearHistory": canClear}, nil
}

func (a *AttentionServiceAdapter) Evaluate(ctx context.Context, userID, orgID string) (any, error) {
	if a.Svc == nil {
		return nil, nil
	}
	return nil, a.Svc.Evaluate(ctx, userID, orgID)
}

func (a *AttentionServiceAdapter) ClearHistory(ctx context.Context, userID, orgID string) error {
	if a.Svc == nil {
		return nil
	}
	_, err := a.Svc.ClearHistory(ctx, userID, orgID)
	return err
}

func (a *AttentionServiceAdapter) Acknowledge(ctx context.Context, userID, orgID, issueID string) error {
	if a.Svc == nil {
		return nil
	}
	return a.Svc.Acknowledge(ctx, userID, orgID, issueID)
}

func (a *AttentionServiceAdapter) Resolve(ctx context.Context, userID, orgID, issueID string) error {
	if a.Svc == nil {
		return nil
	}
	return a.Svc.Resolve(ctx, userID, orgID, issueID)
}

func (a *AttentionServiceAdapter) Act(ctx context.Context, userID, orgID, issueID string, req monitoring.AttentionActionRequest) error {
	if a.Svc == nil {
		return nil
	}
	_, err := a.Svc.Act(ctx, userID, orgID, issueID, req)
	return err
}
