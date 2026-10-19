CREATE TABLE IF NOT EXISTS custom_app_templates (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id TEXT NOT NULL,
    template JSONB NOT NULL,
    created_by_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_custom_app_template_org_app UNIQUE (organization_id, app_id)
);

CREATE INDEX IF NOT EXISTS idx_custom_app_template_org ON custom_app_templates(organization_id);
