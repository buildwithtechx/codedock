-- 031_webhook_deliveries.sql
-- Webhook delivery ledger and idempotency claim table

CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    source TEXT NOT NULL,
    hook_id TEXT,
    delivery_id TEXT,
    event TEXT NOT NULL,
    auth_result TEXT,
    outcome TEXT NOT NULL,
    action_ref TEXT,
    status_code INTEGER,
    error TEXT,
    summary JSONB,
    client_ip TEXT,
    user_agent TEXT,
    duration_ms INTEGER,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_webhook_deliveries_github_delivery
ON webhook_deliveries (source, delivery_id)
WHERE source = 'github' AND delivery_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_project ON webhook_deliveries (project_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_org ON webhook_deliveries (organization_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_hook ON webhook_deliveries (hook_id, received_at DESC);
