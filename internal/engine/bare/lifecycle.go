package bare

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (r *Runtime) Observe(ctx context.Context, app *models.AppService, target models.RuntimeTarget) (*models.WorkloadObservation, error) {
	script, err := ownedScript(app)
	if err != nil {
		return nil, err
	}
	raw, err := r.runner.Host(ctx, target.BareNode, script+"systemctl show \"$unit\" --property=ActiveState,SubState,NRestarts,MemoryCurrent\n")
	if err != nil {
		return nil, err
	}
	values := map[string]string{}
	for _, line := range strings.Split(raw, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			values[key] = value
		}
	}
	result := &models.WorkloadObservation{Kind: "bare", Status: "STOPPED", Desired: 1, MetricsError: "Native CPU utilization is not sampled", Pods: []models.RuntimePod{}}
	if app.Status == models.AppServiceStatusStopped {
		result.Desired = 0
	}
	if values["ActiveState"] == "active" {
		result.Status = "READY"
		result.Available = 1
	} else if values["ActiveState"] == "failed" {
		result.Status = "FAILED"
	}
	restarts, _ := strconv.Atoi(values["NRestarts"])
	memory, _ := strconv.ParseInt(values["MemoryCurrent"], 10, 64)
	result.Pods = append(result.Pods, models.RuntimePod{Name: "codedock-app-" + app.ID, Node: target.BareNode.ServerID, Phase: values["SubState"], Ready: result.Available == 1, Restarts: restarts, MemoryBytes: memory})
	return result, nil
}
func (r *Runtime) Lifecycle(ctx context.Context, app *models.AppService, target models.RuntimeTarget, action string) error {
	script, err := ownedScript(app)
	if err != nil {
		return err
	}
	script += "test ! -e \"$base/pending\"\n"
	switch action {
	case "stop":
		script += "systemctl stop \"$unit\"\n"
	case "restart":
		script += "systemctl restart \"$unit\"\n" + readinessScript(app)
	case "remove":
		script += "systemctl disable --now \"$unit\"\nrm -f \"/etc/systemd/system/$unit.service\"\nsystemctl daemon-reload\n"
	default:
		return fmt.Errorf("native runtime supports stop, restart and remove; it has one supervised process")
	}
	return r.runner.Script(ctx, target.BareNode, app.ID, script)
}
func (r *Runtime) Logs(ctx context.Context, app *models.AppService, target models.RuntimeTarget) (string, error) {
	script, err := ownedScript(app)
	if err != nil {
		return "", err
	}
	return r.runner.Host(ctx, target.BareNode, script+"journalctl --no-pager -n 200 -u \"$unit\"\n")
}
func (r *Runtime) Exec(ctx context.Context, app *models.AppService, target models.RuntimeTarget, request models.RuntimeExecRequest) (string, error) {
	if len(request.Command) == 0 || len(request.Command) > 32 {
		return "", fmt.Errorf("provide one to thirty-two arguments")
	}
	script, err := ownedScript(app)
	if err != nil {
		return "", err
	}
	args := []string{}
	for _, arg := range request.Command {
		if len(arg) > 4096 || strings.ContainsRune(arg, 0) {
			return "", fmt.Errorf("invalid native command argument")
		}
		args = append(args, ssh.ShellQuote(arg))
	}
	return r.runner.Host(ctx, target.BareNode, script+"cd \"$base/data\"\nrunuser -u \"$user\" -- "+strings.Join(args, " ")+"\n")
}
