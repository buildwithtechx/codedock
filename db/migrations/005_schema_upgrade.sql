CREATE TABLE clusters (
 id TEXT PRIMARY KEY,
 project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
 organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 name TEXT NOT NULL,
 version TEXT NOT NULL,
 nodes_json TEXT NOT NULL,
 encrypted_token TEXT NOT NULL,
 revision INTEGER NOT NULL DEFAULT 1,
 status TEXT NOT NULL DEFAULT 'saved',
 error TEXT NOT NULL DEFAULT '',
 updated_at TEXT NOT NULL
);
