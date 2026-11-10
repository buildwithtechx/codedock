package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"codedock/internal/adapters"
	"codedock/internal/config"
	"codedock/internal/deploy"
	"codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	"codedock/internal/platform/ssh"
)

func configureListeners(s *appServices, deployer *deploy.Deployer, accessLogPath string) {
	serverSettings, _ := s.settingsRepo.GetServerSettings(context.Background())
	if serverSettings != nil && serverSettings.DockerCleanupCron != "" {
		_ = s.cronManager.ScheduleDockerCleanup(serverSettings.DockerCleanupCron)
	}
	if serverSettings != nil && serverSettings.DiskUsageCron != "" {
		_ = s.cronManager.ScheduleDiskUsageCheck(serverSettings.DiskUsageCron, serverSettings.DiskUsageThreshold)
	}

	s.deploymentListeners.Register()

	s.sshManager.OnServerUnreachable = func(server *ssh.Server, err error) {
		baseURL := strings.TrimRight(config.Get().Server.DashboardURL, "/")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		s.dispatcherService.Dispatch(&notifications.NotificationEvent{
			EventType: notifications.EventServerUnreachable,
			Title:     "Server unreachable: " + server.Name,
			Message:   fmt.Sprintf("Cannot reach Docker daemon on server %s (%s): %v", server.Name, server.IPAddress, err),
			Level:     "error",
			URL:       fmt.Sprintf("%s/servers/%s", baseURL, server.ID),
		})
	}

	s.backupManager.OnBackupFailed = func(recordID, errStr string) {
		baseURL := strings.TrimRight(config.Get().Server.DashboardURL, "/")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		s.dispatcherService.Dispatch(&notifications.NotificationEvent{
			EventType: notifications.EventBackupFailed,
			Title:     "Backup failed: " + recordID,
			Message:   fmt.Sprintf("Backup execution failed: %s", errStr),
			Level:     "error",
			URL:       fmt.Sprintf("%s/backups", baseURL),
		})
	}

	s.deploymentService.BeforeProjectDeployment = func(ctx context.Context, projectID string) error {
		return s.backupService.BeforeProjectDeployment(ctx, projectID, "")
	}
	s.deploymentService.BeforeDeployment = s.backupService.BeforeDeployment

	s.autoscaler.Start()

	if accessLogPath != "" && s.trafficRepo != nil && s.analyticsService != nil {
		geoResolver, _ := observability.NewGeoIPResolver(os.Getenv("GEOIP_DB_PATH"))
		collector := observability.NewCollector(accessLogPath, &adapters.TrafficCollectorAdapter{Repo: s.trafficRepo}, s.analyticsService.AttributeDomain, geoResolver)
		go collector.Run(context.Background())
	}
}
