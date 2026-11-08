package migrations

type FingerprintResponse struct {
	Fingerprint string `json:"fingerprint"`
}

type MigrationSourceResponse struct {
	Source *MigrationSource `json:"source"`
}

type MigrationRunResponse struct {
	Run *MigrationRun `json:"run"`
}

type MigrationPreviewResponse struct {
	Preview *MigrationPreview `json:"preview"`
}

type ActiveMigrationsResponse struct {
	Runs []MigrationRun `json:"runs"`
}
