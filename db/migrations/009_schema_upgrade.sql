ALTER TABLE cluster_databases ADD COLUMN cluster_id TEXT REFERENCES clusters(id) ON DELETE CASCADE;
ALTER TABLE cluster_databases ADD COLUMN encrypted_config TEXT NOT NULL DEFAULT '';
ALTER TABLE cluster_databases ADD COLUMN error TEXT NOT NULL DEFAULT '';
CREATE INDEX idx_cluster_databases_cluster ON cluster_databases(cluster_id);
