ALTER TABLE backup_configs ADD COLUMN owner_id TEXT NOT NULL DEFAULT '';
ALTER TABLE backup_configs ADD COLUMN project_id TEXT NOT NULL DEFAULT '';
UPDATE backup_configs SET project_id=COALESCE((SELECT project_id FROM app_services WHERE id=backup_configs.service_id),(SELECT project_id FROM databases WHERE id=backup_configs.database_id),'');
UPDATE backup_configs SET owner_id=COALESCE((SELECT m.user_id FROM organization_members m JOIN projects p ON p.organization_id=m.organization_id JOIN users u ON u.id=m.user_id WHERE p.id=backup_configs.project_id AND m.permission='owner' AND m.status='accepted' AND u.is_active=1 ORDER BY m.accepted_at LIMIT 1),'') WHERE project_id<>'';
UPDATE backup_configs SET owner_id=COALESCE((SELECT id FROM users WHERE role IN ('owner','admin') AND is_active=1 ORDER BY CASE role WHEN 'owner' THEN 0 ELSE 1 END,created_at LIMIT 1),'') WHERE project_id='';
