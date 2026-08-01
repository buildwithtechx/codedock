ALTER TABLE domains ADD COLUMN dns_provision_status TEXT DEFAULT 'pending';
ALTER TABLE domains ADD COLUMN dns_provider TEXT DEFAULT '';
ALTER TABLE domains ADD COLUMN dns_provisioned_ip TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN spaceship_api_secret TEXT DEFAULT '';

UPDATE domains SET dns_provision_status = 'provisioned' WHERE dns_provision_status = 'success';
