CREATE TABLE IF NOT EXISTS mail_servers (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    server_id TEXT,
    domain TEXT NOT NULL,
    maild_url TEXT NOT NULL,
    webmail_url TEXT NOT NULL DEFAULT '',
    api_token TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_mail_server_org_domain UNIQUE (organization_id, domain)
);

CREATE INDEX IF NOT EXISTS idx_mail_server_org ON mail_servers(organization_id);
