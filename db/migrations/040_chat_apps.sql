CREATE TABLE IF NOT EXISTS chat_app_connections (
    id TEXT PRIMARY KEY NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    transport TEXT NOT NULL DEFAULT 'webhook',
    purpose TEXT NOT NULL DEFAULT 'both',
    name TEXT NOT NULL DEFAULT '',
    bot_username TEXT NOT NULL DEFAULT '',
    credentials_encrypted TEXT NOT NULL DEFAULT '',
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    status TEXT NOT NULL DEFAULT 'connected',
    last_error TEXT,
    last_event_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chat_app_connections_user ON chat_app_connections(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_app_connections_provider ON chat_app_connections(provider);

CREATE TABLE IF NOT EXISTS chat_app_identities (
    id TEXT PRIMARY KEY NOT NULL,
    connection_id TEXT NOT NULL REFERENCES chat_app_connections(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    external_user_id TEXT NOT NULL,
    external_username TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'operator',
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    paired_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_chat_identity_conn_ext UNIQUE(connection_id, external_user_id)
);

CREATE INDEX IF NOT EXISTS idx_chat_identities_user ON chat_app_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_identities_conn ON chat_app_identities(connection_id);

CREATE TABLE IF NOT EXISTS chat_app_challenges (
    id TEXT PRIMARY KEY NOT NULL,
    connection_id TEXT NOT NULL REFERENCES chat_app_connections(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    external_user_id TEXT NOT NULL,
    external_channel_id TEXT NOT NULL DEFAULT '',
    token TEXT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chat_challenges_token ON chat_app_challenges(token);
CREATE INDEX IF NOT EXISTS idx_chat_challenges_code ON chat_app_challenges(code);

CREATE TABLE IF NOT EXISTS chat_app_commands (
    id TEXT PRIMARY KEY NOT NULL,
    connection_id TEXT NOT NULL REFERENCES chat_app_connections(id) ON DELETE CASCADE,
    identity_id TEXT REFERENCES chat_app_identities(id) ON DELETE SET NULL,
    provider TEXT NOT NULL,
    command TEXT NOT NULL,
    args TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    response TEXT NOT NULL DEFAULT '',
    confirmation_token TEXT,
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chat_commands_conn ON chat_app_commands(connection_id);
CREATE INDEX IF NOT EXISTS idx_chat_commands_token ON chat_app_commands(confirmation_token);
