CREATE TABLE IF NOT EXISTS managed_providers (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    provider TEXT NOT NULL DEFAULT 'hetzner',
    label TEXT NOT NULL DEFAULT '',
    encrypted_token TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS managed_quotas (
    organization_id TEXT PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    max_servers INTEGER NOT NULL DEFAULT 5,
    max_memory_gb INTEGER NOT NULL DEFAULT 32,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE servers ADD COLUMN provider TEXT NOT NULL DEFAULT '';
ALTER TABLE servers ADD COLUMN external_id TEXT NOT NULL DEFAULT '';
ALTER TABLE servers ADD COLUMN region TEXT NOT NULL DEFAULT '';
ALTER TABLE servers ADD COLUMN server_type TEXT NOT NULL DEFAULT '';
CREATE TABLE IF NOT EXISTS managed_servers (
    server_id TEXT PRIMARY KEY REFERENCES servers(id) ON DELETE CASCADE,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    credential_id TEXT NOT NULL REFERENCES managed_providers(id) ON DELETE CASCADE,
    ssh_key_name TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_managed_servers_org ON managed_servers(organization_id);
