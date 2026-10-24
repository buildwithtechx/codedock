ALTER TABLE backup_configs ADD COLUMN sftp_destination_id TEXT;
ALTER TABLE backup_configs ADD COLUMN parent_batch_id TEXT;
ALTER TABLE backup_configs ADD COLUMN sftp_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE backup_configs ADD COLUMN incremental BOOLEAN DEFAULT FALSE;
ALTER TABLE backup_configs ADD COLUMN quiesce_command TEXT DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN unquiesce_command TEXT DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN custom_backup_command TEXT DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN custom_restore_command TEXT DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN file_source_path TEXT DEFAULT '';
ALTER TABLE backup_records ADD COLUMN sftp_destination_id TEXT;
ALTER TABLE backup_records ADD COLUMN sftp_url TEXT DEFAULT '';
ALTER TABLE backup_records ADD COLUMN parent_record_id TEXT;
CREATE TABLE IF NOT EXISTS sftp_destinations (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    host TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 22,
    username TEXT NOT NULL,
    password TEXT NOT NULL DEFAULT '',
    private_key TEXT NOT NULL DEFAULT '',
    path_prefix TEXT NOT NULL DEFAULT '',
    last_verified_at TIMESTAMPTZ,
    last_verify_error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS backup_policy_batches (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    schedule TEXT NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'UTC',
    timeout INTEGER NOT NULL DEFAULT 3600,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
