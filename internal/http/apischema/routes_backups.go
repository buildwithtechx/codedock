package apischema

func backupConfigSchema() Schema {
	return Obj("Backup configuration",
		RF("id", ID("Backup ID")),
		RF("projectId", Str("Project ID")),
		F("databaseId", Str("Database ID")),
		F("serviceId", Str("Service ID")),
		F("volumeName", Str("Volume name")),
		F("name", Str("Backup name")),
		F("description", Str("Description")),
		F("backupEnabled", Bool("Backup enabled")),
		F("s3Enabled", Bool("S3 enabled")),
		F("sftpEnabled", Bool("SFTP enabled")),
		F("incremental", Bool("Incremental")),
		F("schedule", Str("Cron schedule")),
		F("timezone", Str("Timezone")),
		F("timeout", Int("Timeout seconds")),
		F("retentionDays", Int("Retention days")),
		F("maxBackups", Int("Max backups")),
		F("status", Str("Status")),
		F("createdAt", Str("Creation time")),
		F("updatedAt", Str("Last update")),
	)
}

func backupRecordSchema() Schema {
	return Obj("Backup record",
		RF("id", ID("Record ID")),
		RF("backupConfigId", Str("Backup ID")),
		F("databaseId", Str("Database ID")),
		RF("status", Str("Status")),
		F("fileSizeBytes", Int64("File size bytes")),
		F("s3Url", Str("S3 URL")),
		F("sftpUrl", Str("SFTP URL")),
		F("sha256", Str("SHA-256")),
		F("protectedUntil", Int64("Protection expiry")),
		F("logs", Str("Logs")),
		F("startedAt", Str("Start time")),
		F("completedAt", Str("Completion time")),
	)
}

func policyBatchSchema() Schema {
	return Obj("Policy batch",
		RF("id", ID("Batch ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Batch name")),
		F("description", Str("Description")),
		F("schedule", Str("Cron schedule")),
		F("timezone", Str("Timezone")),
		F("timeout", Int("Timeout seconds")),
		F("status", Str("Status")),
	)
}

func s3DestinationSchema() Schema {
	return Obj("S3 destination",
		RF("id", ID("Destination ID")),
		RF("name", Str("Name")),
		F("description", Str("Description")),
		RF("provider", Str("Provider")),
		F("endpoint", Str("Endpoint")),
		RF("bucket", Str("Bucket")),
		RF("region", Str("Region")),
		F("pathPrefix", Str("Path prefix")),
		F("isDefault", Bool("Default destination")),
		F("accessKeyId", Str("Access key ID")),
	)
}

func sftpDestinationSchema() Schema {
	return Obj("SFTP destination",
		RF("id", ID("Destination ID")),
		RF("name", Str("Name")),
		F("description", Str("Description")),
		RF("host", Str("Host")),
		RF("port", Int("Port")),
		RF("username", Str("Username")),
		F("pathPrefix", Str("Path prefix")),
	)
}

func backupOperations() []Operation {
	backups := "backups"
	destinations := "destinations"
	write := Scope("backup:write")
	return []Operation{
		{Method: "GET", Path: "/api/backups", Summary: "List backups", Tags: []string{backups}, Auth: AuthUser, Response: Arr("Backups", backupConfigSchema())},
		{Method: "POST", Path: "/api/backups", Summary: "Create a backup", Tags: []string{backups}, Auth: AuthUser, Code: 201, Request: JSON(Obj("Create backup", RF("name", Str("Backup name")), F("description", Str("Description")), F("databaseId", Str("Database ID")), F("serviceId", Str("Service ID")), F("volumeName", Str("Volume name")), F("schedule", Str("Cron schedule")), F("timezone", Str("Timezone")), F("retentionDays", Int("Retention days")), F("maxBackups", Int("Max backups")), F("s3DestinationId", Str("S3 destination")), F("sftpDestinationId", Str("SFTP destination")), F("backupEnabled", Bool("Backup enabled")), F("s3Enabled", Bool("S3 enabled")), F("disableLocal", Bool("Disable local")))), Response: backupConfigSchema()},
		{Method: "GET", Path: "/api/backups/:id", Summary: "Get a backup", Tags: []string{backups}, Auth: AuthUser, Response: backupConfigSchema()},
		{Method: "PUT", Path: "/api/backups/:id", Summary: "Update a backup", Tags: []string{backups}, Auth: write, Request: JSON(Obj("Update backup", F("name", Str("Name")), F("description", Str("Description")), F("schedule", Str("Cron schedule")), F("retentionDays", Int("Retention days")))), Response: backupConfigSchema()},
		{Method: "DELETE", Path: "/api/backups/:id", Summary: "Delete a backup", Tags: []string{backups}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/backups/:id/trigger", Summary: "Trigger a backup run", Tags: []string{backups}, Auth: AuthUser, Response: backupRecordSchema()},
		{Method: "POST", Path: "/api/backups/:id/runs", Summary: "Start a backup run", Tags: []string{backups}, Auth: write, Response: backupRecordSchema()},
		{Method: "POST", Path: "/api/backups/:id/restore", Summary: "Restore a backup", Tags: []string{backups}, Auth: AuthUser, Request: JSON(Obj("Restore", RF("operationId", Str("Operation ID")), RF("confirmation", Str("Confirmation token")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/backups/:id/records", Summary: "List backup records", Tags: []string{backups}, Auth: AuthUser, Response: Arr("Records", backupRecordSchema())},
		{Method: "GET", Path: "/api/backups/:id/records/:recordId/download", Summary: "Download a backup record", Tags: []string{backups}, Auth: AuthUser, Raw: true, Content: ContentBinary, Response: Str("Archive bytes")},
		{Method: "DELETE", Path: "/api/backups/:id/records/:recordId", Summary: "Delete a backup record", Tags: []string{backups}, Auth: AuthUser, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "GET", Path: "/api/backup-records", Summary: "List all backup records", Tags: []string{backups}, Auth: AuthUser, Response: Arr("Records", backupRecordSchema())},
		{Method: "POST", Path: "/api/backup-records/:recordId/restore/review", Summary: "Review a restore", Tags: []string{backups}, Auth: write, Request: JSON(Obj("Review restore", F("targetDatabaseId", Str("Target database ID")))), Response: Any("Restore review")},
		{Method: "POST", Path: "/api/backup-operations/:operationId/apply", Summary: "Apply a reviewed restore", Tags: []string{backups}, Auth: write, Raw: true, Code: 202, Request: JSON(operationApplySchema()), Response: Obj("Accepted", RF("status", Str("Status")), RF("data", operationSchema()))},
		{Method: "PUT", Path: "/api/backup-records/:recordId/protection", Summary: "Set record protection", Tags: []string{backups}, Auth: write, Request: JSON(Obj("Protection", RF("protectedUntil", Int64("Protection expiry unix time")))), Response: EmptySchema()},
		{Method: "GET", Path: "/api/backup-records/:recordId/volume-restore", Summary: "Get a volume restore target", Tags: []string{backups}, Auth: AuthUser, Response: Obj("Restore target", RF("recordId", Str("Record ID")), RF("volumeName", Str("Volume name")), F("timeoutSeconds", Int("Timeout seconds")))},
		{Method: "POST", Path: "/api/backup-records/:recordId/volume-restore", Summary: "Restore a volume", Tags: []string{backups}, Auth: AuthUser, Request: JSON(Obj("Restore volume", RF("volumeName", Str("Volume name")), RF("confirmOverwrite", Bool("Confirm overwrite")))), Response: EmptySchema()},
		{Method: "DELETE", Path: "/api/backup-records/:recordId/volume-restore", Summary: "Cancel a volume restore", Tags: []string{backups}, Auth: AuthUser, Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:id/policy-batches", Summary: "List policy batches", Tags: []string{backups}, Auth: AuthAdmin, Response: Arr("Batches", policyBatchSchema())},
		{Method: "POST", Path: "/api/projects/:id/policy-batches", Summary: "Create a policy batch", Tags: []string{backups}, Auth: AuthAdmin, Code: 201, Request: JSON(Obj("Create batch", RF("name", Str("Batch name")), F("description", Str("Description")), F("schedule", Str("Cron schedule")), F("timezone", Str("Timezone")), F("timeout", Int("Timeout seconds")))), Response: policyBatchSchema()},
		{Method: "POST", Path: "/api/policy-batches/:batchId/trigger", Summary: "Trigger a policy batch", Tags: []string{backups}, Auth: AuthAdmin, Code: 202, Response: Arr("Records", backupRecordSchema())},
		{Method: "GET", Path: "/api/s3-destinations", Summary: "List S3 destinations", Tags: []string{destinations}, Auth: AuthAdmin, Response: Arr("Destinations", s3DestinationSchema())},
		{Method: "POST", Path: "/api/s3-destinations", Summary: "Create an S3 destination", Tags: []string{destinations}, Auth: AuthAdmin, Code: 201, Request: JSON(Obj("Create S3", RF("name", Str("Name")), RF("bucket", Str("Bucket")), RF("region", Str("Region")), F("provider", Str("Provider")), F("endpoint", Str("Endpoint")), F("pathPrefix", Str("Path prefix")), F("accessKeyId", Str("Access key ID")), F("secretAccessKey", Str("Secret key")), F("isDefault", Bool("Default")))), Response: s3DestinationSchema()},
		{Method: "PUT", Path: "/api/s3-destinations/:id", Summary: "Update an S3 destination", Tags: []string{destinations}, Auth: AuthAdmin, Request: JSON(Obj("Update S3", F("name", Str("Name")), F("bucket", Str("Bucket")), F("region", Str("Region")))), Response: s3DestinationSchema()},
		{Method: "DELETE", Path: "/api/s3-destinations/:id", Summary: "Delete an S3 destination", Tags: []string{destinations}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/s3-destinations/:id/default", Summary: "Set the default S3 destination", Tags: []string{destinations}, Auth: AuthAdmin, Response: OkSchema()},
		{Method: "POST", Path: "/api/s3-destinations/:id/verify", Summary: "Verify an S3 destination", Tags: []string{destinations}, Auth: AuthAdmin, Response: OkSchema()},
		{Method: "POST", Path: "/api/s3-destinations/verify", Summary: "Verify an S3 draft", Tags: []string{destinations}, Auth: AuthAdmin, Request: JSON(Obj("Verify S3 draft", RF("bucket", Str("Bucket")), RF("region", Str("Region")), F("endpoint", Str("Endpoint")), F("accessKeyId", Str("Access key ID")), F("secretAccessKey", Str("Secret key")))), Response: OkSchema()},
		{Method: "GET", Path: "/api/sftp-destinations", Summary: "List SFTP destinations", Tags: []string{destinations}, Auth: AuthAdmin, Response: Arr("Destinations", sftpDestinationSchema())},
		{Method: "POST", Path: "/api/sftp-destinations", Summary: "Create an SFTP destination", Tags: []string{destinations}, Auth: AuthAdmin, Code: 201, Request: JSON(Obj("Create SFTP", RF("name", Str("Name")), RF("host", Str("Host")), RF("port", Int("Port")), RF("username", Str("Username")), F("password", Str("Password")), F("privateKey", Str("Private key")), F("pathPrefix", Str("Path prefix")), F("description", Str("Description")))), Response: sftpDestinationSchema()},
		{Method: "DELETE", Path: "/api/sftp-destinations/:id", Summary: "Delete an SFTP destination", Tags: []string{destinations}, Auth: AuthAdmin, Raw: true, Code: 204, Response: EmptySchema()},
		{Method: "POST", Path: "/api/sftp-destinations/:id/verify", Summary: "Verify an SFTP destination", Tags: []string{destinations}, Auth: AuthAdmin, Response: OkSchema()},
	}
}
