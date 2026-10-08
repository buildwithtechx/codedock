CREATE TABLE IF NOT EXISTS migration_sources (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT '',
    ssh_host TEXT NOT NULL DEFAULT '',
    ssh_port INTEGER NOT NULL DEFAULT 22,
    ssh_user TEXT NOT NULL DEFAULT 'root',
    ssh_auth_method TEXT NOT NULL DEFAULT 'key',
    ssh_key TEXT NOT NULL DEFAULT '',
    ssh_password TEXT NOT NULL DEFAULT '',
    fingerprint TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS migration_runs (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL DEFAULT '',
    source_id TEXT NOT NULL DEFAULT '',
    source_kind TEXT NOT NULL DEFAULT 'external',
    project_id TEXT NOT NULL DEFAULT '',
    target_server_id TEXT NOT NULL DEFAULT '',
    mode TEXT NOT NULL DEFAULT 'move',
    status TEXT NOT NULL DEFAULT 'pending',
    phase TEXT NOT NULL DEFAULT '',
    selection TEXT NOT NULL DEFAULT '',
    progress TEXT NOT NULL DEFAULT '',
    logs TEXT NOT NULL DEFAULT '',
    prompt TEXT NOT NULL DEFAULT '',
    token_hash TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    cancel_requested INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_migration_runs_org ON migration_runs(organization_id);
CREATE INDEX IF NOT EXISTS idx_migration_runs_source ON migration_runs(source_id);
