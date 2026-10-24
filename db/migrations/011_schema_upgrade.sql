CREATE TABLE IF NOT EXISTS redis_snapshots (
 id TEXT PRIMARY KEY,
 database_id TEXT NOT NULL REFERENCES cluster_databases(id) ON DELETE CASCADE,
 project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
 cluster_id TEXT NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
 s3_destination_id TEXT REFERENCES s3_destinations(id) ON DELETE SET NULL,
 s3_key TEXT NOT NULL,
 size_bytes INTEGER NOT NULL DEFAULT 0,
 status TEXT NOT NULL DEFAULT 'COMPLETED',
 error TEXT NOT NULL DEFAULT '',
 created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_redis_snapshots_database ON redis_snapshots(database_id);
