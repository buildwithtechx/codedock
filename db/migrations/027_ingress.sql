CREATE TABLE IF NOT EXISTS ingress_connections (
    id VARCHAR(64) PRIMARY KEY,
    server_id VARCHAR(64) NOT NULL,
    credential_id VARCHAR(64) NOT NULL,
    account_id VARCHAR(64) NOT NULL,
    tunnel_id VARCHAR(64),
    name VARCHAR(255) NOT NULL,
    connector VARCHAR(64) NOT NULL DEFAULT 'managed',
    edge_port INT NOT NULL DEFAULT 80,
    managed_tunnel BOOLEAN NOT NULL DEFAULT TRUE,
    default_for_domains BOOLEAN NOT NULL DEFAULT FALSE,
    default_access_policy_id VARCHAR(64),
    status VARCHAR(64) NOT NULL DEFAULT 'ready',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ingress_hostnames (
    id VARCHAR(64) PRIMARY KEY,
    domain_id VARCHAR(64),
    connection_id VARCHAR(64) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    access_policy_id VARCHAR(64),
    status VARCHAR(64) NOT NULL DEFAULT 'ready',
    phase VARCHAR(64) NOT NULL DEFAULT 'dns',
    error TEXT,
    ownership JSONB DEFAULT '{}'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
