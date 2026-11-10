ALTER TABLE app_services ADD COLUMN stack TEXT NOT NULL DEFAULT '';
ALTER TABLE app_services ADD COLUMN stop_grace_period_seconds INTEGER NOT NULL DEFAULT 0;
