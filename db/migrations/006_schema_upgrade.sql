CREATE TABLE service_runtimes (
 service_id TEXT PRIMARY KEY REFERENCES app_services(id) ON DELETE CASCADE,
 project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
 encrypted_config TEXT NOT NULL,
 encrypted_journal TEXT NOT NULL DEFAULT '',
 revision INTEGER NOT NULL DEFAULT 1,
 status TEXT NOT NULL DEFAULT 'CONFIGURED',
 error TEXT NOT NULL DEFAULT '',
 updated_at TEXT NOT NULL
);
