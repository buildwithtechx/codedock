package bare

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"regexp"
	"strings"
)

type Runner interface {
	Script(context.Context, models.ClusterNode, string, string) error
	Host(context.Context, models.ClusterNode, string) (string, error)
}
type Runtime struct{ runner Runner }

func NewRuntime(runner Runner) *Runtime { return &Runtime{runner} }
func ValidateIdentity(app *models.AppService, target models.RuntimeTarget) error {
	if _, err := uuid.Parse(app.ID); err != nil {
		return fmt.Errorf("bare application requires a UUID")
	}
	if _, err := uuid.Parse(app.ProjectID); err != nil {
		return fmt.Errorf("bare project requires a UUID")
	}
	if err := kubernetes.ValidateNode(target.BareNode); err != nil {
		return err
	}
	return nil
}
func Validate(app *models.AppService, target models.RuntimeTarget) error {
	if err := ValidateIdentity(app, target); err != nil {
		return err
	}
	parsed, err := url.Parse(target.BareReleaseURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
		return fmt.Errorf("bare release requires an HTTPS artifact URL without credentials or fragment")
	}
	if !regexp.MustCompile(`^[a-fA-F0-9]{64}$`).MatchString(target.BareSHA256) {
		return fmt.Errorf("bare release requires its SHA256 checksum")
	}
	if len(target.BareCommand) == 0 || len(target.BareCommand) > 32 || target.BareCommand[0] != "./app" {
		return fmt.Errorf("bare releases execute ./app followed by up to thirty-one arguments")
	}
	for _, arg := range target.BareCommand {
		if strings.ContainsRune(arg, 0) || len(arg) > 4096 {
			return fmt.Errorf("invalid native argument")
		}
	}
	if app.Replicas > 1 || len(target.Volumes) > 0 || len(app.Volumes) > 0 || target.ClusterID != "" || len(target.NodeIDs) > 0 {
		return fmt.Errorf("bare runtime supports one supervised process and its private data directory; Docker volumes and cluster placement are unavailable")
	}
	if app.Domain != "" {
		return fmt.Errorf("bare runtime uses its native port; clear the managed HTTP domain before selecting it")
	}
	if app.RuntimeMode != models.RuntimeModeWorker && (app.InternalPort < 1024 || app.InternalPort > 65535) {
		return fmt.Errorf("bare HTTP services need an unprivileged port from 1024 to 65535")
	}
	if app.HealthCheckPath != "" && (!strings.HasPrefix(app.HealthCheckPath, "/") || strings.ContainsAny(app.HealthCheckPath, "\r\n")) {
		return fmt.Errorf("invalid native readiness path")
	}
	return nil
}
func (r *Runtime) Deploy(ctx context.Context, app *models.AppService, target models.RuntimeTarget, variables map[string]string, release string) error {
	if err := Validate(app, target); err != nil {
		return err
	}
	script, err := DeploymentScript(app, target, variables, release)
	if err != nil {
		return err
	}
	return r.runner.Script(ctx, target.BareNode, app.ID, script)
}
func (r *Runtime) Recover(ctx context.Context, app *models.AppService, target models.RuntimeTarget) error {
	script, err := recoveryScript(app)
	if err != nil {
		return err
	}
	return r.runner.Script(ctx, target.BareNode, app.ID, script)
}
func (r *Runtime) Finalize(ctx context.Context, app *models.AppService, target models.RuntimeTarget) error {
	script, err := ownedScript(app)
	if err != nil {
		return err
	}
	return r.runner.Script(ctx, target.BareNode, app.ID, script+"rm -f \"$base/pending\"\n")
}
