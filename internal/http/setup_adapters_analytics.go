package http

import (
	"context"
	"strconv"

	"codedock/internal/modules/monitoring"
)

type analyticsServiceAdapter struct {
	svc *monitoring.Service
}

func (a *analyticsServiceAdapter) Periods(ctx context.Context, orgID string) (any, error) {
	if a.svc == nil {
		return []string{}, nil
	}
	return a.svc.Periods(), nil
}

func (a *analyticsServiceAdapter) Dashboard(ctx context.Context, orgID string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.Dashboard(ctx, "", orgID)
}

func (a *analyticsServiceAdapter) Summary(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.Summary(ctx, userID, orgID, projectID, "", from, to)
}

func (a *analyticsServiceAdapter) Overview(ctx context.Context, userID, orgID, projectID, domain, from, to, period string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.Overview(ctx, userID, orgID, projectID, domain, from, to, period)
}

func (a *analyticsServiceAdapter) Geo(ctx context.Context, userID, orgID, projectID, from, to, period string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.Geo(ctx, userID, orgID, projectID, from, to)
}

func (a *analyticsServiceAdapter) GetPaths(ctx context.Context, userID, orgID, projectID string) (bool, error) {
	if a.svc == nil {
		return false, nil
	}
	return a.svc.PathsEnabled(ctx, userID, orgID, projectID)
}

func (a *analyticsServiceAdapter) SetPaths(ctx context.Context, userID, orgID, projectID string, enabled bool) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.SetPathsEnabled(ctx, userID, orgID, projectID, enabled)
}

func (a *analyticsServiceAdapter) DeploymentStats(ctx context.Context, userID, orgID, projectID string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.DeploymentStats(ctx, userID, orgID, projectID, 30)
}

func (a *analyticsServiceAdapter) Usage(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.svc == nil || a.svc.Usage() == nil {
		return nil, nil
	}
	return a.svc.Usage().Current(ctx, projectID)
}

func (a *analyticsServiceAdapter) Resources(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.svc == nil || a.svc.Usage() == nil {
		return nil, nil
	}
	return a.svc.Usage().Resources(ctx, projectID)
}

func (a *analyticsServiceAdapter) UsageHistory(ctx context.Context, userID, orgID, projectID, period string) (any, error) {
	if a.svc == nil || a.svc.Usage() == nil {
		return nil, nil
	}
	hours, _ := strconv.Atoi(period)
	if hours <= 0 {
		hours = 24
	}
	return a.svc.Usage().History(ctx, projectID, hours)
}

func (a *analyticsServiceAdapter) Container(ctx context.Context, userID, orgID, projectID, serviceID, period string) (any, error) {
	if a.svc == nil || a.svc.Usage() == nil {
		return nil, nil
	}
	return a.svc.Usage().Container(ctx, projectID, serviceID)
}

func (a *analyticsServiceAdapter) CurrentUsage(ctx context.Context, projectID string) (any, error) {
	if a.svc == nil || a.svc.Usage() == nil {
		return nil, nil
	}
	return a.svc.Usage().Current(ctx, projectID)
}

type attentionServiceAdapter struct {
	svc *monitoring.AttentionManager
}

func (a *attentionServiceAdapter) List(ctx context.Context, userID, orgID string) (any, error) {
	if a.svc == nil {
		return []any{}, nil
	}
	issues, err := a.svc.ListIssues(ctx, userID, orgID, "")
	if err != nil {
		return nil, err
	}
	canClear, _ := a.svc.CanClearHistory(ctx, userID, orgID)
	return map[string]any{"issues": issues, "canClearHistory": canClear}, nil
}

func (a *attentionServiceAdapter) Evaluate(ctx context.Context, userID, orgID string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return nil, a.svc.Evaluate(ctx, userID, orgID)
}

func (a *attentionServiceAdapter) ClearHistory(ctx context.Context, userID, orgID string) error {
	if a.svc == nil {
		return nil
	}
	_, err := a.svc.ClearHistory(ctx, userID, orgID)
	return err
}

func (a *attentionServiceAdapter) Acknowledge(ctx context.Context, userID, orgID, issueID string) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.Acknowledge(ctx, userID, orgID, issueID)
}

func (a *attentionServiceAdapter) Resolve(ctx context.Context, userID, orgID, issueID string) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.Resolve(ctx, userID, orgID, issueID)
}

func (a *attentionServiceAdapter) Act(ctx context.Context, userID, orgID, issueID string, req monitoring.AttentionActionRequest) error {
	if a.svc == nil {
		return nil
	}
	_, err := a.svc.Act(ctx, userID, orgID, issueID, req)
	return err
}
