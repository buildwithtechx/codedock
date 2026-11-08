package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
)

func (cm *CronManager) ScheduleDockerCleanup(schedule string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cleanSchedule := strings.TrimSpace(schedule)
	if len(strings.Fields(cleanSchedule)) == 5 && !strings.HasPrefix(cleanSchedule, "@") {
		cleanSchedule = "0 " + cleanSchedule
	}

	if entryID, exists := cm.entries[dockerCleanupEntryID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, dockerCleanupEntryID)
	}

	entryID, err := cm.cronEngine.AddFunc(cleanSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		cm.DockerCleanup(ctx)
	})
	if err != nil {
		return fmt.Errorf("invalid docker cleanup schedule '%s': %w", schedule, err)
	}
	cm.entries[dockerCleanupEntryID] = entryID
	slog.Info("docker cleanup scheduled", "schedule", schedule)
	return nil
}

func (cm *CronManager) DockerCleanup(ctx context.Context) {
	if cm.dockerClient == nil {
		slog.Warn("docker cleanup skipped", "reason", "no Docker client")
		return
	}
	slog.Info("running docker cleanup")

	reclaimed := uint64(0)

	report, err := cm.dockerClient.ContainersPrune(ctx, filters.NewArgs(filters.Arg("until", "24h")))
	if err == nil {
		reclaimed += report.SpaceReclaimed
	}

	imgReport, err := cm.dockerClient.ImagesPrune(ctx, filters.NewArgs(filters.Arg("dangling", "true")))
	if err == nil {
		reclaimed += imgReport.SpaceReclaimed
	}

	volReport, err := cm.dockerClient.VolumesPrune(ctx, filters.NewArgs(filters.Arg("until", "24h")))
	if err == nil {
		reclaimed += volReport.SpaceReclaimed
	}

	if reclaimed > 0 {
		slog.Info("docker cleanup reclaimed space", "bytes", reclaimed)
	} else {
		slog.Info("docker cleanup completed", "result", "nothing to clean")
	}
}

func (cm *CronManager) ScheduleDiskUsageCheck(schedule string, threshold int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cleanSchedule := strings.TrimSpace(schedule)
	if len(strings.Fields(cleanSchedule)) == 5 && !strings.HasPrefix(cleanSchedule, "@") {
		cleanSchedule = "0 " + cleanSchedule
	}

	if entryID, exists := cm.entries[diskUsageEntryID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, diskUsageEntryID)
	}

	entryID, err := cm.cronEngine.AddFunc(cleanSchedule, func() {
		usageCmd := exec.Command("df", "--output=pcent", "/")
		out, err := usageCmd.Output()
		if err != nil {
			return
		}
		pctStr := strings.TrimSpace(string(out))
		pctStr = strings.TrimSuffix(strings.TrimSpace(strings.Split(pctStr, "\n")[1]), "%")
		usage, err := strconv.Atoi(pctStr)
		if err != nil {
			return
		}
		if usage > threshold {
			slog.Warn("disk usage above threshold", "usage_pct", usage, "threshold_pct", threshold)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid disk usage check schedule '%s': %w", schedule, err)
	}
	cm.entries[diskUsageEntryID] = entryID
	return nil
}
