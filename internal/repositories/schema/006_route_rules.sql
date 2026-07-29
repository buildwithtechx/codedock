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
