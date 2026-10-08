CREATE TABLE IF NOT EXISTS autoscaling_policies (
    service_id TEXT PRIMARY KEY REFERENCES app_services(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT 0,
    min_replicas INTEGER NOT NULL DEFAULT 1 CHECK(min_replicas BETWEEN 1 AND 10),
    max_replicas INTEGER NOT NULL DEFAULT 5 CHECK(max_replicas BETWEEN min_replicas AND 10),
    scale_up_cpu REAL NOT NULL DEFAULT 80 CHECK(scale_up_cpu > 0 AND scale_up_cpu <= 100),
    scale_down_cpu REAL NOT NULL DEFAULT 20 CHECK(scale_down_cpu >= 0 AND scale_down_cpu < scale_up_cpu),
    cooldown_seconds INTEGER NOT NULL DEFAULT 300 CHECK(cooldown_seconds BETWEEN 120 AND 86400),
    last_cpu REAL NOT NULL DEFAULT 0,
    last_decision TEXT NOT NULL DEFAULT '',
    last_evaluated_at TEXT NOT NULL DEFAULT '',
    last_scaled_at TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL
);
