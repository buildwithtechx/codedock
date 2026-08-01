DROP INDEX IF EXISTS idx_refresh_token_revocations_token_hash;

ALTER TABLE domains ADD COLUMN dns_provision_status TEXT DEFAULT 'pending';
ALTER TABLE domains ADD COLUMN dns_provider TEXT DEFAULT '';
ALTER TABLE domains ADD COLUMN dns_provisioned_ip TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN spaceship_api_secret TEXT DEFAULT '';
