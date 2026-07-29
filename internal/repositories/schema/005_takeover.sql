CREATE TABLE IF NOT EXISTS takeover_runs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_host TEXT NOT NULL,
    source_platform TEXT NOT NULL DEFAULT 'docker',
    status TEXT NOT NULL DEFAULT 'scanning',
    discovered_json TEXT NOT NULL DEFAULT '',
    adopted_project_ids TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_takeover_runs_user ON takeover_runs(user_id);
