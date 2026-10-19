CREATE TABLE IF NOT EXISTS notification_device_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL DEFAULT '',
    token TEXT NOT NULL UNIQUE,
    platform TEXT NOT NULL DEFAULT 'expo',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_notification_device_tokens_user ON notification_device_tokens(user_id);
