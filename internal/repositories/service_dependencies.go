package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
)

func (r *AppServiceRepo) DeploymentDependencies(ctx context.Context, app *models.AppService) ([]string, error) {
	sources := []string{}
	err := r.db.SelectContext(ctx, &sources, `SELECT source FROM topology_dependencies WHERE environment_id=$1 AND target=$2`, app.EnvironmentID, "app-"+app.ID)
	return sources, err
}
