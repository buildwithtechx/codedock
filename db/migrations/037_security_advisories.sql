CREATE TABLE IF NOT EXISTS advisory_feed (
    id TEXT PRIMARY KEY NOT NULL,
    body JSONB NOT NULL DEFAULT '{}'::jsonb,
    etag TEXT NOT NULL DEFAULT '',
    revision BIGINT NOT NULL DEFAULT 0,
    checked_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ,
    error TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS project_security (
    project_id TEXT PRIMARY KEY NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    deployment_id TEXT NOT NULL DEFAULT '',
    inventory JSONB NOT NULL DEFAULT '{}'::jsonb,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    auto_refresh BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS security_refresh_attempts (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    deployment_id TEXT NOT NULL DEFAULT '',
    advisory_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_security_refresh_project ON security_refresh_attempts(project_id);

