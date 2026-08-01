ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

ALTER TABLE users ADD COLUMN plan_type text NOT NULL DEFAULT 'free';
ALTER TABLE users ADD COLUMN stripe_customer_id text;
ALTER TABLE users ADD COLUMN stripe_subscription_id text;
ALTER TABLE users ADD COLUMN stripe_price_id text;

ALTER TABLE users ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE IF NOT EXISTS refresh_token_revocations (
	 id TEXT PRIMARY KEY,
	 token_hash TEXT UNIQUE NOT NULL,
	 user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	 expires_at DATETIME NOT NULL,
	 revoked_at DATETIME DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_refresh_token_revocations_user_id ON refresh_token_revocations (user_id);

CREATE TABLE IF NOT EXISTS takeover_runs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_host TEXT NOT NULL,
    source_platform TEXT NOT NULL DEFAULT 'docker',
    status TEXT NOT NULL DEFAULT 'scanning',
    discovered_json TEXT NOT NULL DEFAULT '',
    adopted_project_ids TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_takeover_runs_user ON takeover_runs(user_id);

CREATE TABLE IF NOT EXISTS route_rules (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    rule_type TEXT NOT NULL,
    spec_json TEXT NOT NULL DEFAULT '{}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_route_rules_service ON route_rules(service_id);

ALTER TABLE servers ADD COLUMN ssh_host TEXT DEFAULT '';
ALTER TABLE servers ADD COLUMN ssh_port INTEGER DEFAULT 22;
ALTER TABLE servers ADD COLUMN ssh_user TEXT DEFAULT 'root';
ALTER TABLE servers ADD COLUMN ssh_key TEXT DEFAULT '';
ALTER TABLE servers ADD COLUMN ssh_password TEXT DEFAULT '';
