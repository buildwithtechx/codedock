CREATE TABLE IF NOT EXISTS organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT UNIQUE,
    logo TEXT DEFAULT '',
    is_team BOOLEAN DEFAULT 0,
    plan_tier TEXT DEFAULT 'free',
    stripe_customer_id TEXT,
    subscription_status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT DEFAULT '',
    password_hash TEXT NOT NULL,
    role TEXT DEFAULT 'developer',
    email_verified BOOLEAN DEFAULT 0,
    is_active BOOLEAN DEFAULT 1,
    totp_enabled BOOLEAN DEFAULT 0,
    totp_secret TEXT DEFAULT '',
    recovery_codes TEXT DEFAULT '',
    oauth_provider TEXT DEFAULT '',
    plan_type TEXT DEFAULT 'hobby',
    stripe_customer_id TEXT DEFAULT '',
    stripe_subscription_id TEXT DEFAULT '',
    stripe_price_id TEXT DEFAULT '',
    last_login DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_members (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    email TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'member',
    permission TEXT NOT NULL DEFAULT 'member',
    status TEXT NOT NULL DEFAULT 'active',
    invited_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME,
    UNIQUE(organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_org_members_org ON organization_members(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_members_user ON organization_members(user_id);

CREATE TABLE IF NOT EXISTS invites (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    role TEXT DEFAULT 'developer',
    token TEXT UNIQUE NOT NULL,
    invited_by TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    accepted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_token_revocations (
    id TEXT PRIMARY KEY,
    token_hash TEXT UNIQUE NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at DATETIME NOT NULL,
    revoked_at DATETIME DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_refresh_token_revocations_user_id ON refresh_token_revocations (user_id);

CREATE TABLE IF NOT EXISTS personal_access_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    prefix TEXT NOT NULL,
    access_level TEXT DEFAULT 'read_write',
    project_scope TEXT DEFAULT 'all',
    allowed_projects TEXT DEFAULT '[]',
    expires_at DATETIME,
    last_used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS servers (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL DEFAULT '',
    status TEXT DEFAULT 'online',
    is_local BOOLEAN NOT NULL DEFAULT 0,
    ssh_host TEXT DEFAULT '',
    ssh_port INTEGER DEFAULT 22,
    ssh_user TEXT DEFAULT 'root',
    ssh_auth_method TEXT DEFAULT 'key',
    ssh_key TEXT DEFAULT '',
    ssh_private_key TEXT DEFAULT '',
    ssh_password TEXT DEFAULT '',
    ssh_transport TEXT NOT NULL DEFAULT 'direct',
    ssh_jump_host TEXT DEFAULT '',
    worker_token TEXT NOT NULL DEFAULT '',
    last_seen_at DATETIME,
    metrics TEXT DEFAULT '{}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_servers_user ON servers(user_id);
CREATE INDEX IF NOT EXISTS idx_servers_org ON servers(organization_id);

CREATE TABLE IF NOT EXISTS project_apps (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    git_provider TEXT DEFAULT 'github',
    git_owner TEXT DEFAULT '',
    git_repo TEXT DEFAULT '',
    git_url TEXT DEFAULT '',
    installation_id INTEGER DEFAULT 0,
    favicon TEXT DEFAULT '',
    favicon_checked_at DATETIME,
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_project_apps_org ON project_apps(organization_id);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL REFERENCES project_apps(id) ON DELETE CASCADE,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    server_id TEXT REFERENCES servers(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT DEFAULT '',
    environment_name TEXT NOT NULL DEFAULT 'Production',
    environment_slug TEXT NOT NULL DEFAULT 'production',
    environment_type TEXT NOT NULL DEFAULT 'production',
    is_app BOOLEAN NOT NULL DEFAULT 0,
    app_template_id TEXT DEFAULT '',
    local_path TEXT DEFAULT '',
    git_provider TEXT DEFAULT 'github',
    git_owner TEXT DEFAULT '',
    git_repo TEXT DEFAULT '',
    git_branch TEXT DEFAULT 'main',
    git_url TEXT DEFAULT '',
    clone_token_encrypted TEXT DEFAULT '',
    clone_token_set_at DATETIME,
    routing_config_json TEXT DEFAULT '{}',
    readiness_json TEXT DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'ready',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(app_id, environment_slug)
);

CREATE INDEX IF NOT EXISTS idx_projects_app ON projects(app_id);
CREATE INDEX IF NOT EXISTS idx_projects_org ON projects(organization_id);

CREATE TABLE IF NOT EXISTS app_services (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id TEXT DEFAULT '',
    kind TEXT NOT NULL DEFAULT 'compose',
    name TEXT NOT NULL,
    icon TEXT DEFAULT 'git',
    image TEXT DEFAULT '',
    image_ref TEXT DEFAULT '',
    build TEXT DEFAULT '',
    dockerfile TEXT DEFAULT '',
    dockerfile_path TEXT DEFAULT '',
    repository_url TEXT DEFAULT '',
    branch TEXT DEFAULT 'main',
    root_directory TEXT DEFAULT '/',
    runtime_mode TEXT DEFAULT 'web',
    install_command TEXT DEFAULT '',
    build_command TEXT DEFAULT '',
    start_command TEXT DEFAULT '',
    build_engine TEXT DEFAULT 'railpack',
    command TEXT DEFAULT '',
    static_output TEXT DEFAULT '',
    health_check_path TEXT DEFAULT '/',
    build_args_json TEXT DEFAULT '{}',
    ports_json TEXT DEFAULT '[]',
    depends_on_json TEXT DEFAULT '[]',
    environment_json TEXT DEFAULT '{}',
    volumes_json TEXT DEFAULT '[]',
    namespace_volumes BOOLEAN NOT NULL DEFAULT 1,
    command_argv_json TEXT DEFAULT '[]',
    restart TEXT DEFAULT 'unless-stopped',
    advanced_json TEXT DEFAULT '{}',
    exposed BOOLEAN NOT NULL DEFAULT 0,
    exposed_port TEXT DEFAULT '',
    internal_port INTEGER DEFAULT 3000,
    domain TEXT DEFAULT '',
    custom_domain TEXT DEFAULT '',
    domain_type TEXT DEFAULT 'free',
    public_endpoints_json TEXT DEFAULT '[]',
    replicas INTEGER DEFAULT 1,
    cpu_request REAL DEFAULT 0.5,
    memory_limit_mb INTEGER DEFAULT 512,
    cpu_limit REAL DEFAULT 0,
    memory_limit INTEGER DEFAULT 0,
    deploy_token TEXT DEFAULT '',
    enable_pr_previews BOOLEAN DEFAULT 0,
    maintenance_mode BOOLEAN DEFAULT 0,
    registry_id TEXT REFERENCES registries(id) ON DELETE SET NULL,
    status TEXT DEFAULT 'stopped',
    container_id TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, name)
);

CREATE VIEW IF NOT EXISTS services AS SELECT * FROM app_services;

CREATE INDEX IF NOT EXISTS idx_services_project ON app_services(project_id);
CREATE INDEX IF NOT EXISTS idx_services_env ON app_services(environment_id);

CREATE TABLE IF NOT EXISTS env_vars (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id TEXT REFERENCES app_services(id) ON DELETE CASCADE,
    environment_id TEXT DEFAULT '',
    key TEXT NOT NULL,
    encrypted_value TEXT NOT NULL,
    is_secret BOOLEAN NOT NULL DEFAULT 1,
    target TEXT DEFAULT 'production',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, service_id, key)
);

CREATE INDEX IF NOT EXISTS idx_env_vars_project ON env_vars(project_id);
CREATE INDEX IF NOT EXISTS idx_env_vars_service ON env_vars(service_id);

CREATE TABLE IF NOT EXISTS deployments (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    service_id TEXT REFERENCES app_services(id) ON DELETE SET NULL,
    branch TEXT NOT NULL DEFAULT 'main',
    commit_hash TEXT DEFAULT '',
    commit_message TEXT DEFAULT '',
    commit_sha_before TEXT DEFAULT '',
    trigger TEXT NOT NULL DEFAULT 'manual',
    environment TEXT NOT NULL DEFAULT 'production',
    framework TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'queued',
    image_ref TEXT DEFAULT '',
    build_duration_ms INTEGER DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 1,
    build_logs TEXT DEFAULT '',
    container_id TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_deployments_project ON deployments(project_id);
CREATE INDEX IF NOT EXISTS idx_deployments_service ON deployments(service_id);

CREATE TABLE IF NOT EXISTS service_deployments (
    id TEXT PRIMARY KEY,
    deployment_id TEXT NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    image_ref TEXT DEFAULT '',
    container_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'queued',
    error TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(deployment_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_service_deployments_dep ON service_deployments(deployment_id);

CREATE TABLE IF NOT EXISTS domains (
    id TEXT PRIMARY KEY,
    owner_type TEXT NOT NULL DEFAULT 'project',
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    service_id TEXT REFERENCES app_services(id) ON DELETE CASCADE,
    hostname TEXT NOT NULL UNIQUE,
    domain_name TEXT NOT NULL DEFAULT '',
    dns_provision_status TEXT DEFAULT 'pending',
    target_port INTEGER DEFAULT 80,
    target_path TEXT DEFAULT '/',
    domain_type TEXT DEFAULT 'custom',
    is_primary BOOLEAN NOT NULL DEFAULT 0,
    redirect_to TEXT DEFAULT '',
    redirect_status INTEGER DEFAULT 301,
    external_ingress BOOLEAN NOT NULL DEFAULT 0,
    ssl_cert_status TEXT DEFAULT 'pending',
    dns_verification_status TEXT DEFAULT 'pending',
    dns_provider TEXT DEFAULT '',
    dns_provisioned_ip TEXT DEFAULT '',
    path_prefix TEXT DEFAULT '/',
    server_id TEXT REFERENCES servers(id) ON DELETE SET NULL,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_domains_hostname ON domains(hostname);
CREATE INDEX IF NOT EXISTS idx_domains_project ON domains(project_id);
CREATE INDEX IF NOT EXISTS idx_domains_service ON domains(service_id);

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

CREATE TABLE IF NOT EXISTS cluster_databases (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    server_id TEXT REFERENCES servers(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    engine TEXT NOT NULL,
    version TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'provisioning',
    intent TEXT NOT NULL DEFAULT 'apply',
    port INTEGER NOT NULL,
    username TEXT NOT NULL,
    database_name TEXT NOT NULL,
    secret_encrypted TEXT NOT NULL,
    env_key TEXT DEFAULT '',
    container_id TEXT DEFAULT '',
    internal_dns TEXT DEFAULT '',
    external_dns TEXT DEFAULT '',
    progress_json TEXT DEFAULT '{"steps":[],"logs":[]}',
    backup_request_id TEXT DEFAULT '',
    cpu_limit REAL DEFAULT 0,
    memory_limit INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cluster_db_project ON cluster_databases(project_id);

CREATE TABLE IF NOT EXISTS backup_configs (
    id TEXT PRIMARY KEY,
    database_id TEXT REFERENCES databases(id) ON DELETE SET NULL,
    service_id TEXT REFERENCES app_services(id) ON DELETE SET NULL,
    volume_name TEXT DEFAULT '',
    s3_destination_id TEXT REFERENCES s3_destinations(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    db_user TEXT DEFAULT '',
    db_password TEXT DEFAULT '',
    backup_enabled BOOLEAN DEFAULT 1,
    s3_enabled BOOLEAN DEFAULT 0,
    disable_local BOOLEAN DEFAULT 0,
    schedule TEXT NOT NULL,
    timezone TEXT DEFAULT 'UTC',
    timeout INTEGER DEFAULT 3600,
    retention_days INTEGER DEFAULT 7,
    max_backups INTEGER DEFAULT 0,
    max_storage_gb REAL DEFAULT 0,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS s3_destinations (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    provider TEXT DEFAULT 's3',
    endpoint TEXT NOT NULL,
    bucket TEXT NOT NULL,
    region TEXT NOT NULL,
    access_key_id TEXT NOT NULL,
    secret_access_key TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backup_records (
    id TEXT PRIMARY KEY,
    backup_config_id TEXT NOT NULL REFERENCES backup_configs(id) ON DELETE CASCADE,
    database_id TEXT,
    s3_destination_id TEXT,
    status TEXT DEFAULT 'running',
    file_path TEXT DEFAULT '',
    file_size_bytes INTEGER DEFAULT 0,
    s3_url TEXT DEFAULT '',
    logs TEXT DEFAULT '',
    started_at DATETIME,
    completed_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_backup_records_config ON backup_records(backup_config_id);

CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id TEXT PRIMARY KEY,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    service_id TEXT REFERENCES app_services(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    schedule TEXT NOT NULL,
    command TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    last_run_at DATETIME,
    last_output TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_project ON scheduled_tasks(project_id);

CREATE TABLE IF NOT EXISTS service_incidents (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id TEXT REFERENCES app_services(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    severity TEXT NOT NULL DEFAULT 'warning',
    summary TEXT NOT NULL,
    detail TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'open',
    resolved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_service_incidents_project ON service_incidents(project_id);

CREATE TABLE IF NOT EXISTS mail_servers (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    server_id TEXT REFERENCES servers(id) ON DELETE CASCADE,
    domain TEXT NOT NULL UNIQUE,
    hostname TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'installing',
    dkim_selector TEXT DEFAULT 'mail',
    dkim_private_key TEXT DEFAULT '',
    dkim_public_key TEXT DEFAULT '',
    spf_record TEXT DEFAULT '',
    dmarc_record TEXT DEFAULT '',
    smtp_port INTEGER DEFAULT 587,
    imap_port INTEGER DEFAULT 993,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mail_inbound_rules (
    id TEXT PRIMARY KEY,
    mail_server_id TEXT NOT NULL REFERENCES mail_servers(id) ON DELETE CASCADE,
    recipient_pattern TEXT NOT NULL,
    action TEXT NOT NULL DEFAULT 'webhook',
    forward_target TEXT DEFAULT '',
    webhook_url TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_git_providers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    encrypted_access_token TEXT NOT NULL,
    account_name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, provider)
);

CREATE TABLE IF NOT EXISTS github_apps (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    app_id TEXT NOT NULL,
    installation_id TEXT DEFAULT '',
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    webhook_secret TEXT NOT NULL,
    private_key TEXT NOT NULL,
    is_public BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS oauth_providers (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    provider_name TEXT UNIQUE NOT NULL,
    enabled BOOLEAN DEFAULT 0,
    client_id TEXT DEFAULT '',
    client_secret TEXT DEFAULT '',
    redirect_uri TEXT DEFAULT '',
    base_url TEXT DEFAULT '',
    tenant TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS server_settings (
    id TEXT PRIMARY KEY,
    traefik_wildcard_ip TEXT DEFAULT '127.0.0.1',
    registration_enabled BOOLEAN DEFAULT 1,
    registration_domain_allowlist TEXT DEFAULT '',
    custom_dns_resolvers TEXT DEFAULT '',
    dns_validation_enabled BOOLEAN DEFAULT 1,
    ip_allowlist TEXT DEFAULT '',
    mcp_server_enabled BOOLEAN DEFAULT 1,
    default_wildcard_domain TEXT DEFAULT '',
    update_check_cron TEXT DEFAULT '0 * * * *',
    auto_update_enabled BOOLEAN DEFAULT 0,
    current_version TEXT DEFAULT '0.1.0',
    latest_version TEXT DEFAULT '0.1.0',
    last_update_check TEXT DEFAULT '',
    updated_at TEXT,
    site_name TEXT DEFAULT '',
    public_ipv4 TEXT DEFAULT '',
    public_ipv6 TEXT DEFAULT '',
    show_sponsorship_popup BOOLEAN DEFAULT 1,
    disable_two_step_confirmation BOOLEAN DEFAULT 0,
    panel_domain TEXT NOT NULL DEFAULT '',
    concurrent_builds INTEGER NOT NULL DEFAULT 2,
    deployment_timeout INTEGER NOT NULL DEFAULT 3600,
    server_timezone TEXT NOT NULL DEFAULT 'UTC',
    docker_cleanup_cron TEXT NOT NULL DEFAULT '0 0 * * *',
    disk_usage_threshold INTEGER NOT NULL DEFAULT 80,
    disk_usage_cron TEXT NOT NULL DEFAULT '0 23 * * *',
    cloudflare_api_token TEXT DEFAULT '',
    namecheap_api_user TEXT DEFAULT '',
    namecheap_api_key TEXT DEFAULT '',
    namecheap_client_ip TEXT DEFAULT '',
    spaceship_api_key TEXT DEFAULT '',
    spaceship_api_secret TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS notification_settings (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    discord_webhook_url TEXT DEFAULT '',
    discord_ping_enabled BOOLEAN DEFAULT 0,
    discord_enabled BOOLEAN DEFAULT 0,
    slack_webhook_url TEXT DEFAULT '',
    slack_enabled BOOLEAN DEFAULT 0,
    telegram_bot_token TEXT DEFAULT '',
    telegram_chat_id TEXT DEFAULT '',
    telegram_enabled BOOLEAN DEFAULT 0,
    smtp_host TEXT DEFAULT '',
    smtp_port INTEGER DEFAULT 587,
    smtp_user TEXT DEFAULT '',
    smtp_password TEXT DEFAULT '',
    smtp_from_name TEXT DEFAULT '',
    smtp_from_address TEXT DEFAULT '',
    smtp_enabled BOOLEAN DEFAULT 0,
    resend_api_key TEXT DEFAULT '',
    resend_enabled BOOLEAN DEFAULT 0,
    pushover_user_key TEXT DEFAULT '',
    pushover_api_token TEXT DEFAULT '',
    pushover_enabled BOOLEAN DEFAULT 0,
    generic_webhook_url TEXT DEFAULT '',
    generic_webhook_enabled BOOLEAN DEFAULT 0,
    notification_alerts BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    details TEXT DEFAULT '',
    ip_address TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS ai_settings (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    default_provider TEXT DEFAULT 'none',
    openai_key TEXT DEFAULT '',
    openai_model TEXT DEFAULT '',
    anthropic_key TEXT DEFAULT '',
    anthropic_model TEXT DEFAULT '',
    google_key TEXT DEFAULT '',
    google_model TEXT DEFAULT '',
    mistral_key TEXT DEFAULT '',
    mistral_model TEXT DEFAULT '',
    groq_key TEXT DEFAULT '',
    groq_model TEXT DEFAULT '',
    deepseek_key TEXT DEFAULT '',
    deepseek_model TEXT DEFAULT '',
    xai_key TEXT DEFAULT '',
    xai_model TEXT DEFAULT '',
    moonshot_key TEXT DEFAULT '',
    moonshot_model TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE IF NOT EXISTS serverless_functions_code (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    runtime TEXT NOT NULL,
    code_content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(service_id)
);

CREATE TABLE IF NOT EXISTS pr_previews (
    id TEXT PRIMARY KEY,
    service_id TEXT REFERENCES app_services(id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    pr_number INTEGER,
    branch TEXT,
    commit_hash TEXT,
    status TEXT,
    preview_domain TEXT,
    container_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS log_drains (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL DEFAULT '',
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    drain_type TEXT NOT NULL,
    endpoint_url TEXT NOT NULL,
    auth_token TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS registries (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    registry_url TEXT NOT NULL,
    username TEXT DEFAULT '',
    password_token TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS environments (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(project_id, name)
);

CREATE INDEX IF NOT EXISTS idx_environments_project ON environments(project_id);

CREATE TABLE IF NOT EXISTS databases (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    engine TEXT NOT NULL,
    version TEXT NOT NULL,
    port INTEGER NOT NULL,
    username TEXT NOT NULL,
    encrypted_password TEXT NOT NULL,
    database_name TEXT NOT NULL,
    volume_path TEXT NOT NULL DEFAULT '',
    container_id TEXT DEFAULT '',
    status TEXT DEFAULT 'stopped',
    internal_dns TEXT DEFAULT '',
    external_dns TEXT DEFAULT '',
    custom_args TEXT DEFAULT '',
    environment_id TEXT DEFAULT '',
    project_id TEXT DEFAULT '',
    cpu_limit REAL DEFAULT 0,
    memory_limit INTEGER DEFAULT 0,
    logical_replication INTEGER DEFAULT 0,
    server_id TEXT REFERENCES servers(id) ON DELETE SET NULL,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_databases_project ON databases(project_id);

CREATE TABLE IF NOT EXISTS service_volumes (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    host_path TEXT NOT NULL,
    container_path TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_service_volumes_service ON service_volumes(service_id);

CREATE TABLE IF NOT EXISTS service_webhooks (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    event_types TEXT DEFAULT '',
    include_pr_environments BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_service_webhooks_service ON service_webhooks(service_id);

CREATE TABLE IF NOT EXISTS service_vars (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL,
    environment_id TEXT DEFAULT '',
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    is_secret INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(service_id, key)
);

CREATE INDEX IF NOT EXISTS idx_service_vars_service ON service_vars(service_id);

CREATE TABLE IF NOT EXISTS project_tokens (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    token_prefix TEXT NOT NULL DEFAULT '',
    token_hash TEXT NOT NULL,
    scopes TEXT DEFAULT '',
    ip_allowlist TEXT DEFAULT '',
    expires_at TEXT,
    last_used_at TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_project_tokens_project ON project_tokens(project_id);

CREATE TABLE IF NOT EXISTS dns_records (
    id TEXT PRIMARY KEY,
    domain_name TEXT NOT NULL,
    record_type TEXT NOT NULL,
    record_name TEXT NOT NULL,
    record_value TEXT NOT NULL,
    ttl INTEGER DEFAULT 3600,
    server_id TEXT REFERENCES servers(id) ON DELETE SET NULL,
    organization_id TEXT REFERENCES organizations(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dns_records_domain ON dns_records(domain_name);
