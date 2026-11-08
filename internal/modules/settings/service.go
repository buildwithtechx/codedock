package settings

import (
	"context"
	"fmt"
	"net"

	"github.com/robfig/cron/v3"

	"codedock/internal/config"
	"codedock/internal/telemetry"
)

type SettingsService struct {
	repo SettingsRepository
}

func NewSettingsService(r SettingsRepository) *SettingsService {
	return &SettingsService{repo: r}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*ServerSettings, error) {
	cfg, err := s.repo.GetServerSettings(ctx)
	if err != nil {
		return nil, err
	}
	cfg.TelemetryEnabled = config.Get().Telemetry.Enabled
	return cfg, nil
}

func (s *SettingsService) IsNotificationEmailEnabled(ctx context.Context) (bool, error) {
	return s.repo.IsNotificationEmailEnabled(ctx)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, cfg *ServerSettings) error {
	if cfg.PublicIPv4 != "" {
		ip := net.ParseIP(cfg.PublicIPv4)
		if ip == nil || ip.To4() == nil {
			return fmt.Errorf("invalid public IPv4 address")
		}
	}
	if cfg.PublicIPv6 != "" {
		ip := net.ParseIP(cfg.PublicIPv6)
		if ip == nil || ip.To4() != nil {
			return fmt.Errorf("invalid public IPv6 address")
		}
	}
	if cfg.DockerCleanupCron != "" {
		if _, err := cron.ParseStandard(cfg.DockerCleanupCron); err != nil {
			return fmt.Errorf("invalid docker cleanup cron expression: %w", err)
		}
	}
	if cfg.DiskUsageCron != "" {
		if _, err := cron.ParseStandard(cfg.DiskUsageCron); err != nil {
			return fmt.Errorf("invalid disk usage cron expression: %w", err)
		}
	}
	if cfg.UpdateCheckCron != "" {
		if _, err := cron.ParseStandard(cfg.UpdateCheckCron); err != nil {
			return fmt.Errorf("invalid update check cron expression: %w", err)
		}
	}

	config.Get().Telemetry.Enabled = cfg.TelemetryEnabled
	telemetry.Init()

	return s.repo.UpdateServerSettings(ctx, cfg)
}
