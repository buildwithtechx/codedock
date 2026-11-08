package scheduler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	cronv3 "github.com/robfig/cron/v3"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type Store interface {
	ListScheduledTasks() ([]models.ScheduledTask, error)
	GetScheduledTask(id string) (*models.ScheduledTask, error)
	GetProject(id string) (*models.ProjectConfig, error)
	GetAppService(id string) (*models.AppService, error)
	UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

const (
	dockerCleanupEntryID = "docker-cleanup"
	diskUsageEntryID     = "disk-usage"
)

type CronManager struct {
	dockerClient *client.Client
	store        Store
	cronEngine   *cronv3.Cron
	entries      map[string]cronv3.EntryID
	mu           sync.Mutex
}

func NewCronManager(dockerClient *client.Client, s Store) *CronManager {
	return &CronManager{
		dockerClient: dockerClient,
		store:        s,
		cronEngine:   cronv3.New(cronv3.WithSeconds()),
		entries:      make(map[string]cronv3.EntryID),
	}
}

func (cm *CronManager) Start() error {
	if err := cm.Reconcile(); err != nil {
		return err
	}
	cm.cronEngine.Start()
	slog.Info("cron manager started")
	return nil
}

func (cm *CronManager) Reconcile() error {
	scheduledTasks, err := cm.store.ListScheduledTasks()
	if err != nil {
		return fmt.Errorf("failed to load scheduledTasks during cron reconcile: %w", err)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	active := make(map[string]*models.ScheduledTask, len(scheduledTasks))
	for i := range scheduledTasks {
		if scheduledTasks[i].Status == "active" {
			scheduledTaskCopy := scheduledTasks[i]
			active[scheduledTasks[i].ID] = &scheduledTaskCopy
		}
	}
	for _, scheduledTask := range active {
		if err := cm.registerScheduledTaskLocked(scheduledTask); err != nil {
			slog.Warn("failed to register cron scheduledTask", "name", scheduledTask.Name, "id", scheduledTask.ID, "err", err)
		}
	}
	for id, entryID := range cm.entries {
		if id == dockerCleanupEntryID || id == diskUsageEntryID {
			continue
		}
		if _, ok := active[id]; !ok {
			cm.cronEngine.Remove(entryID)
			delete(cm.entries, id)
		}
	}
	return nil
}

func (cm *CronManager) ServeElected(ctx context.Context) {
	if err := cm.Reconcile(); err != nil {
		slog.Warn("cron reconcile failed", "error", err)
	}
	cm.cronEngine.Start()
	defer cm.cronEngine.Stop()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := cm.Reconcile(); err != nil {
				slog.Warn("cron reconcile failed", "error", err)
			}
		}
	}
}

func (cm *CronManager) Stop() {
	if cm.cronEngine != nil {
		cm.cronEngine.Stop()
	}
}

func (cm *CronManager) RegisterScheduledTask(j *models.ScheduledTask) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.registerScheduledTaskLocked(j)
}

func (cm *CronManager) registerScheduledTaskLocked(j *models.ScheduledTask) error {
	if entryID, exists := cm.entries[j.ID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, j.ID)
	}
	if j.Status != "active" {
		return nil
	}
	schedule := strings.TrimSpace(j.Schedule)
	if len(strings.Fields(schedule)) == 5 && !strings.HasPrefix(schedule, "@") {
		schedule = "0 " + schedule
	}
	scheduledTaskID := j.ID
	entryID, err := cm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		_, _ = cm.ExecuteScheduledTask(ctx, scheduledTaskID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s': %w", j.Schedule, err)
	}
	cm.entries[j.ID] = entryID
	return nil
}

func (cm *CronManager) UnregisterScheduledTask(scheduledTaskID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if entryID, exists := cm.entries[scheduledTaskID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, scheduledTaskID)
	}
}

func (cm *CronManager) ExecuteScheduledTask(ctx context.Context, scheduledTaskID string) (string, error) {
	j, err := cm.store.GetScheduledTask(scheduledTaskID)
	if err != nil || j == nil {
		return "", fmt.Errorf("scheduledTask %s not found: %w", scheduledTaskID, err)
	}
	app, err := cm.store.GetAppService(j.ServiceID)
	if err != nil || app == nil {
		return "", fmt.Errorf("service %s for scheduledTask %s not found: %w", j.ServiceID, scheduledTaskID, err)
	}
	containerName := utils.NormalizeContainerName(app.ID)
	inspectResp, err := cm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		errMsg := fmt.Sprintf("cannot run scheduledTask: service container %s (%s) is stopped or not found", app.Name, containerName)
		now := time.Now()
		if updateErr := cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, errMsg); updateErr != nil {
			slog.Warn("failed to update scheduled task status", "error", updateErr)
		}
		return errMsg, errors.New(errMsg)
	}
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", j.Command},
	}
	execCreateResp, err := cm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		now := time.Now()
		if updateErr := cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, err.Error()); updateErr != nil {
			slog.Warn("failed to update scheduled task status", "error", updateErr)
		}
		return "", fmt.Errorf("failed to create container exec for scheduledTask %s: %w", j.Name, err)
	}
	attachResp, err := cm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		now := time.Now()
		if updateErr := cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, err.Error()); updateErr != nil {
			slog.Warn("failed to update scheduled task status", "error", updateErr)
		}
		return "", fmt.Errorf("failed to attach to container exec for scheduledTask %s: %w", j.Name, err)
	}
	defer attachResp.Close()
	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}
	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		if output != "" {
			output += "\nSTDERR:\n"
		}
		output += stderrBuf.String()
	}
	now := time.Now()
	if updateErr := cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusActive, &now, output); updateErr != nil {
		slog.Warn("failed to update scheduled task status", "error", updateErr)
	}
	return output, nil
}
