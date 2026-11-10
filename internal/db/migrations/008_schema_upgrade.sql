CREATE TABLE runtime_desired (
 service_id TEXT PRIMARY KEY REFERENCES service_runtimes(service_id) ON DELETE CASCADE,
 encrypted_workload TEXT NOT NULL
);
CREATE TABLE cluster_upgrade_journals (
 cluster_id TEXT PRIMARY KEY REFERENCES clusters(id) ON DELETE CASCADE,
 encrypted_plan TEXT NOT NULL
);
