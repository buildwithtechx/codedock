CREATE TABLE IF NOT EXISTS self_hosted_config (
    id TEXT PRIMARY KEY,
    jwt_secret TEXT NOT NULL DEFAULT '',
    refresh_secret TEXT NOT NULL DEFAULT '',
    telemetry_salt TEXT NOT NULL DEFAULT '',
    tls_email TEXT NOT NULL DEFAULT '',
    wildcard_domain TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
