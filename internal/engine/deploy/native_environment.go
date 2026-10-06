package deploy

import (
	"codedock.run/codedock/internal/models"
	"context"
	"io"
)

func (d *Deployer) NativeEnvironment(ctx context.Context, app *models.AppService, logs io.Writer) (map[string]string, error) {
	return d.getEnvironmentVariables(ctx, app, logs)
}
