CREATE TABLE IF NOT EXISTS compose_stacks (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id TEXT NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    encrypted_config TEXT NOT NULL,
    revision INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'saved',
    error TEXT NOT NULL DEFAULT '',
    results TEXT NOT NULL DEFAULT '[]',
    updated_at TEXT NOT NULL,
    UNIQUE(project_id, name)
);
CREATE TABLE IF NOT EXISTS topology_dependencies (
    environment_id TEXT NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    source TEXT NOT NULL,
    target TEXT NOT NULL,
    PRIMARY KEY(environment_id, source, target),
    CHECK(source <> target)
);
ALTER TABLE app_services ADD COLUMN git_user_id TEXT NOT NULL DEFAULT '';
