package http

import (
	"context"
	"strings"

	appsmod "codedock/internal/modules/apps"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/databases"
)

type appDatabaseAdapter struct {
	repo databases.Repository
}

func (a *appDatabaseAdapter) ListByProject(ctx context.Context, projectID string, search string) ([]*appsmod.Database, error) {
	if a.repo == nil {
		return nil, nil
	}
	dbs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]*appsmod.Database, 0, len(dbs))
	for _, d := range dbs {
		if search != "" && !strings.Contains(d.Name, search) {
			continue
		}
		res = append(res, &appsmod.Database{
			ID:           d.ID,
			ProjectID:    d.ProjectID,
			Name:         d.Name,
			Engine:       string(d.Engine),
			Username:     d.Username,
			Password:     d.Password,
			InternalDNS:  d.InternalDNS,
			DatabaseName: d.DatabaseName,
		})
	}
	return res, nil
}

type appClusterPlansAdapter struct {
	repo *modclusters.ClusterDataRepo
}

func (a *appClusterPlansAdapter) Get(ctx context.Context, id string) (*appsmod.ClusterDataPlan, error) {
	if a.repo == nil {
		return nil, nil
	}
	plan, err := a.repo.Get(ctx, id)
	if err != nil || plan == nil {
		return nil, err
	}
	return &appsmod.ClusterDataPlan{
		Record: appsmod.ClusterData{
			ID:        plan.Record.ID,
			ProjectID: plan.Record.ProjectID,
			Status:    plan.Record.Status,
			Spec: appsmod.ClusterDataSpec{
				ID:            plan.Record.Spec.ID,
				Name:          plan.Record.Spec.Name,
				EnvironmentID: plan.Record.Spec.EnvironmentID,
				Engine:        plan.Record.Spec.Engine,
				Instances:     plan.Record.Spec.Instances,
				Shards:        plan.Record.Spec.Shards,
			},
		},
		Password: plan.Password,
	}, nil
}

func (a *appClusterPlansAdapter) ListByProject(ctx context.Context, project string) ([]appsmod.ClusterData, error) {
	if a.repo == nil {
		return nil, nil
	}
	items, err := a.repo.ListByProject(ctx, project)
	if err != nil {
		return nil, err
	}
	res := make([]appsmod.ClusterData, len(items))
	for i, it := range items {
		res[i] = appsmod.ClusterData{
			ID:        it.ID,
			ProjectID: it.ProjectID,
			Status:    it.Status,
			Spec: appsmod.ClusterDataSpec{
				ID:            it.Spec.ID,
				Name:          it.Spec.Name,
				EnvironmentID: it.Spec.EnvironmentID,
				Engine:        it.Spec.Engine,
				Instances:     it.Spec.Instances,
				Shards:        it.Spec.Shards,
			},
		}
	}
	return res, nil
}
