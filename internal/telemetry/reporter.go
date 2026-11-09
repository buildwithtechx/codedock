package telemetry

import (
	"context"
	"database/sql"
	"log/slog"
	"runtime"
	"time"

	"codedock/internal/config"
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
	var instanceID string
	var enabled bool
	row := db.QueryRowContext(context.Background(), `SELECT id, telemetry_enabled FROM server_settings LIMIT 1`)
	if err := row.Scan(&instanceID, &enabled); err != nil || !enabled {
		return
	}

	var activeApps int
	_ = db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM services`).Scan(&activeApps)

	slog.Info("telemetry", "instance", instanceID, "version", version, "os", runtime.GOOS, "arch", runtime.GOARCH, "apps", activeApps)

	Track(instanceID, "telemetry_heartbeat", map[string]any{
		"version":     version,
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"active_apps": activeApps,
		"cloud_mode":  config.Get().Cloud.Enabled,
	})
}
