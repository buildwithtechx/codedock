-- 032_rollback_retention.sql: Rollback artifact protection & disk-pressure reclamation (#1351)

ALTER TABLE projects ADD COLUMN IF NOT EXISTS default_rollback_strategy TEXT NOT NULL DEFAULT 'snapshot';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS rollback_window INTEGER;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS snapshot_size_bytes BIGINT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS capacity_measured_at TIMESTAMPTZ;

ALTER TABLE deployments ADD COLUMN IF NOT EXISTS artifact_retained_at TIMESTAMPTZ;
ALTER TABLE deployments ADD COLUMN IF NOT EXISTS pinned BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_deployments_retention ON deployments (project_id, created_at DESC)
    WHERE status IN ('READY', 'ACTIVE', 'SUCCESS', 'PARTIAL_FAILURE') OR artifact_retained_at IS NOT NULL OR pinned = TRUE;
