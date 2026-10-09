CREATE TABLE IF NOT EXISTS notification_subscriptions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    organization_id TEXT NOT NULL,
    category TEXT NOT NULL,
    channels JSONB NOT NULL DEFAULT '[]',
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, organization_id, category)
);

CREATE TABLE IF NOT EXISTS notification_defaults (
    organization_id TEXT NOT NULL,
    category TEXT NOT NULL,
    channels JSONB NOT NULL DEFAULT '[]',
    enabled BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(organization_id, category)
);
