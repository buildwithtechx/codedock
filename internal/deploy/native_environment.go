package deploy

import (
	"context"
	"io"

	appsmod "codedock/internal/modules/apps"
)

func (d *Deployer) NativeEnvironment(ctx context.Context, app *appsmod.AppService, logs io.Writer) (map[string]string, error) {
	return d.getEnvironmentVariables(ctx, app, logs)
}
