CREATE TABLE IF NOT EXISTS route_protection (
    id TEXT PRIMARY KEY NOT NULL,
    enabled BOOLEAN DEFAULT true NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    hostname TEXT NOT NULL,
    secret TEXT,
    revision TEXT NOT NULL,
    applied_revision TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT route_protection_project_host UNIQUE (project_id, hostname)
);
