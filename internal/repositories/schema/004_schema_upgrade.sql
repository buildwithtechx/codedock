CREATE TABLE operations (
 id TEXT PRIMARY KEY,
 user_id TEXT NOT NULL,
 project_id TEXT NOT NULL DEFAULT '',
 kind TEXT NOT NULL,
 target TEXT NOT NULL,
 status TEXT NOT NULL,
 phase TEXT NOT NULL DEFAULT '',
 effects TEXT NOT NULL,
 encrypted_payload TEXT NOT NULL,
 snapshot TEXT NOT NULL,
 token_hash TEXT NOT NULL,
 error TEXT NOT NULL DEFAULT '',
 logs TEXT NOT NULL DEFAULT '',
 expires_at INTEGER NOT NULL,
 updated_at TEXT NOT NULL
);
CREATE UNIQUE INDEX operations_active_target ON operations(target) WHERE status IN ('RUNNING','CANCELLING');
ALTER TABLE backup_records ADD COLUMN protected_until INTEGER NOT NULL DEFAULT 0;
ALTER TABLE backup_records ADD COLUMN sha256 TEXT NOT NULL DEFAULT '';
ALTER TABLE backup_records ADD COLUMN verified_at TEXT NOT NULL DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN pre_deployment INTEGER NOT NULL DEFAULT 0;
