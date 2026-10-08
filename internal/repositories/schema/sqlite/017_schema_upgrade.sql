CREATE TABLE IF NOT EXISTS traffic_buckets (
    bucket_minute TEXT NOT NULL,
    project_id TEXT NOT NULL DEFAULT '',
    domain TEXT NOT NULL DEFAULT '',
    path TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    requests INTEGER NOT NULL DEFAULT 0,
    bytes INTEGER NOT NULL DEFAULT 0,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (bucket_minute, project_id, domain, path, status)
);
CREATE INDEX IF NOT EXISTS idx_traffic_buckets_project ON traffic_buckets(project_id, bucket_minute);
CREATE TABLE IF NOT EXISTS traffic_visitors (
    day TEXT NOT NULL,
    project_id TEXT NOT NULL DEFAULT '',
    ip TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    requests INTEGER NOT NULL DEFAULT 0,
    bytes INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (day, project_id, ip)
);
CREATE TABLE IF NOT EXISTS traffic_paths (
    project_id TEXT PRIMARY KEY,
    enabled INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS attention_issues (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    kind TEXT NOT NULL DEFAULT '',
    subject TEXT NOT NULL DEFAULT '',
    project_id TEXT NOT NULL DEFAULT '',
    service_id TEXT NOT NULL DEFAULT '',
    severity TEXT NOT NULL DEFAULT 'warning',
    status TEXT NOT NULL DEFAULT 'open',
    title TEXT NOT NULL DEFAULT '',
    detail TEXT NOT NULL DEFAULT '',
    remediation TEXT NOT NULL DEFAULT '',
    action TEXT NOT NULL DEFAULT '',
    action_params TEXT NOT NULL DEFAULT '',
    occurrences INTEGER NOT NULL DEFAULT 1,
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, kind, subject)
);
CREATE INDEX IF NOT EXISTS idx_attention_org_status ON attention_issues(organization_id, status);
