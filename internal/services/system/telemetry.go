package system

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"runtime"
	"time"

	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/telemetry"
)

func StartTelemetryReporter(db *sql.DB, version string) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		time.Sleep(5 * time.Minute)
		pingTelemetry(db, version)

		for range ticker.C {
			pingTelemetry(db, version)
		}
	}()
}

func pingTelemetry(db *sql.DB, version string) {
	settingsRepo := repositories.NewSettingsRepo(db)
	settings, err := settingsRepo.GetServerSettings(context.Background())
	if err != nil {
		return
	}

	if !settings.TelemetryEnabled {
		return
	}

	appRepo := repositories.NewAppServiceRepo(db)
	apps, err := appRepo.ListAll(context.Background())
	activeApps := 0
	if err == nil {
		activeApps = len(apps)
	}

	slog.Info("telemetry", "instance", settings.ID, "version", version, "os", runtime.GOOS, "arch", runtime.GOARCH, "apps", activeApps)

	telemetry.Track(settings.ID, "telemetry_heartbeat", map[string]any{
		"version":     version,
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"active_apps": activeApps,
		"cloud_mode":  os.Getenv("CODEDOCK_CLOUD_MODE") == "true",
	})
}
