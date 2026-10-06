package deploy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/google/uuid"
)

type deploymentProgressKey struct{}

func WithDeploymentProgress(ctx context.Context, progress func(string) error) context.Context {
	return context.WithValue(ctx, deploymentProgressKey{}, progress)
}

func reportDeploymentProgress(ctx context.Context, phase string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if progress, ok := ctx.Value(deploymentProgressKey{}).(func(string) error); ok {
		return progress(phase)
	}
	return nil
}

func (d *Deployer) replaceReplicas(ctx context.Context, app *models.AppService, opts ContainerRunOptions, replicas int, logs io.Writer) (result string, resultErr error) {
	if d.rolloutDirectory != "" {
		if _, err := os.Stat(d.rolloutPath(app.ID)); err == nil {
			return "", fmt.Errorf("previous rollout cleanup is pending; recover it before retrying")
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("check rollout recovery: %w", err)
		}
	}
	manager := d.containerManager
	if err := reportDeploymentProgress(ctx, "STARTING"); err != nil {
		return "", err
	}
	prefix := utils.NormalizeContainerName(app.ID)
	existing, err := manager.dockerClient.ContainerList(ctx, container.ListOptions{All: true, Filters: filters.NewArgs(filters.Arg("name", prefix))})
	if err != nil {
		return "", fmt.Errorf("list previous containers: %w", err)
	}
	journal := &models.RolloutJournal{ID: uuid.NewString(), PreviousApp: *app, Previous: []models.RolloutContainer{}}
	if err := d.saveRollout(journal); err != nil {
		return "", fmt.Errorf("save rollout recovery: %w", err)
	}
	if opts.ExtraLabels == nil {
		opts.ExtraLabels = map[string]string{}
	}
	opts.ExtraLabels["codedock.rollout"] = journal.ID
	created := []string{}
	committed := false
	defer func() {
		if committed {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		var rollbackErr error
		for _, id := range created {
			rollbackErr = errors.Join(rollbackErr, manager.StopAndRemove(cleanupCtx, id))
		}
		for _, old := range journal.Previous {
			inspected, err := manager.Inspect(cleanupCtx, old.ID)
			if err != nil {
				rollbackErr = errors.Join(rollbackErr, err)
				continue
			}
			if inspected.Name != "/"+old.Name {
				rollbackErr = errors.Join(rollbackErr, manager.dockerClient.ContainerRename(cleanupCtx, old.ID, old.Name))
			}
			if old.Running {
				rollbackErr = errors.Join(rollbackErr, manager.dockerClient.ContainerStart(cleanupCtx, old.ID, container.StartOptions{}))
			}
		}
		rollbackErr = errors.Join(rollbackErr, d.store.UpdateAppService(&journal.PreviousApp))
		if rollbackErr == nil {
			rollbackErr = d.clearRollout(app.ID)
		}
		resultErr = errors.Join(resultErr, rollbackErr)
	}()
	for _, old := range existing {
		if len(old.Names) == 0 {
			continue
		}
		name := old.Names[0][1:]
		if name != prefix && !isReplicaName(name, prefix) {
			continue
		}
		saved := models.RolloutContainer{ID: old.ID, Name: name, Running: old.State == "running"}
		journal.Previous = append(journal.Previous, saved)
		if err := d.saveRollout(journal); err != nil {
			return "", err
		}
		if err := manager.dockerClient.ContainerRename(ctx, old.ID, name+"-previous-"+uuid.NewString()[:8]); err != nil {
			return "", fmt.Errorf("retain previous container: %w", err)
		}
		if len(app.Volumes) > 0 && saved.Running {
			if err := manager.dockerClient.ContainerStop(ctx, old.ID, container.StopOptions{}); err != nil {
				return "", fmt.Errorf("stop volume consumer: %w", err)
			}
		}
	}
	for i := 0; i < replicas; i++ {
		opts.Name = prefix
		if replicas > 1 {
			opts.Name = fmt.Sprintf("%s-%d", prefix, i+1)
		}
		if opts.RuntimeMode != models.RuntimeModeWorker && opts.HealthCheckPath == "" {
			opts.HealthCheckPath = "/"
		}
		id, err := manager.createAndStart(ctx, opts)
		if err != nil {
			return "", err
		}
		created = append(created, id)
		if err := reportDeploymentProgress(ctx, "READINESS"); err != nil {
			return "", err
		}
		if app.RuntimeMode == models.RuntimeModeWorker {
			inspected, err := manager.Inspect(ctx, id)
			if err != nil {
				return "", err
			}
			if inspected.State == nil || !inspected.State.Running {
				return "", fmt.Errorf("worker exited before activation")
			}
		} else if err := d.verifyHealthCheck(ctx, app, id, logs); err != nil {
			return "", err
		}
	}
	if err := reportDeploymentProgress(ctx, "ROUTING"); err != nil {
		return "", err
	}
	app.ContainerID, app.Domain, app.Status = created[0], opts.Domain, models.AppServiceStatusRunning
	if err := d.store.UpdateAppService(app); err != nil {
		return "", fmt.Errorf("save activated service: %w", err)
	}
	journal.Committed = true
	if err := d.saveRollout(journal); err != nil {
		return "", fmt.Errorf("save rollout activation: %w", err)
	}
	committed = true
	cleanupCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for _, old := range journal.Previous {
		if err := manager.StopAndRemove(cleanupCtx, old.ID); err != nil && !errdefs.IsNotFound(err) {
			return created[0], fmt.Errorf("deployment activated; remove retained container: %w", err)
		}
	}
	if err := d.clearRollout(app.ID); err != nil {
		return created[0], err
	}
	return created[0], nil
}

func isReplicaName(name, prefix string) bool {
	if len(name) <= len(prefix)+1 || name[:len(prefix)+1] != prefix+"-" {
		return false
	}
	for _, char := range name[len(prefix)+1:] {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
