package apischema

func queryResultSchema() Schema {
	return Obj("Query result",
		RF("columns", Arr("Columns", Str("Column"))),
		RF("rows", Arr("Rows", Map("Row", Any("Value")))),
		RF("rowCount", Int("Row count")),
		RF("executionTimeMs", Int64("Execution ms")),
	)
}

func tableSchemaSchema() Schema {
	return Obj("Table schema",
		RF("name", Str("Table name")),
		RF("columns", Arr("Columns", Obj("Column", RF("name", Str("Name")), F("type", Str("Type"))))),
	)
}

func clusterDataSpecSchema() Schema {
	return Obj("Cluster database spec",
		RF("name", Str("Database name")),
		RF("environmentId", Str("Environment ID")),
		RF("engine", StrEnum("Engine", "postgres", "redis")),
		F("image", Str("Image")),
		F("instances", Int("Instances")),
		F("shards", Int("Redis shards")),
		F("storageGiB", Int("Storage GiB")),
		F("storageClass", Str("Storage class")),
		F("s3DestinationId", Str("S3 destination")),
		F("synchronous", Bool("Synchronous replication")),
	)
}

func clusterDataSchema() Schema {
	return Obj("Cluster database",
		RF("id", ID("Database ID")),
		RF("clusterId", Str("Cluster ID")),
		RF("projectId", Str("Project ID")),
		RF("spec", clusterDataSpecSchema()),
		RF("status", Str("Status")),
		F("error", Str("Error message")),
		F("observedRoles", Arr("Roles", Str("Role"))),
		F("observedVolumes", Arr("Volumes", Str("Volume"))),
		F("updatedAt", Str("Last update")),
	)
}

func redisSnapshotSchema() Schema {
	return Obj("Redis snapshot",
		RF("id", ID("Snapshot ID")),
		RF("databaseId", Str("Database ID")),
		RF("projectId", Str("Project ID")),
		RF("clusterId", Str("Cluster ID")),
		F("s3DestinationId", Str("S3 destination")),
		F("s3Key", Str("S3 key")),
		F("sizeBytes", Int64("Size bytes")),
		RF("status", Str("Status")),
		F("error", Str("Error")),
		F("createdAt", Str("Creation time")),
	)
}

func operationApplySchema() Schema {
	return Obj("Apply operation", RF("operationId", Str("Operation ID")), RF("confirmation", Str("Confirmation token")))
}

func databaseOperations() []Operation {
	databases := "databases"
	clusterData := "cluster-databases"
	db := Scope("database:manage")
	return []Operation{
		{Method: "GET", Path: "/api/databases", Summary: "List databases", Tags: []string{databases}, Auth: db, Response: Arr("Databases", DatabaseSchema())},
		{Method: "POST", Path: "/api/databases", Summary: "Create a database", Tags: []string{databases}, Auth: db, Code: 201, Request: JSON(Obj("Create database", RF("projectId", Str("Project ID")), RF("name", Str("Database name")), RF("engine", Str("Engine")), F("environmentId", Str("Environment ID")), F("version", Str("Version")), F("databaseName", Str("Database name on server")), F("username", Str("Username")), F("password", Str("Password")))), Response: DatabaseSchema()},
		{Method: "GET", Path: "/api/databases/:id", Summary: "Get a database", Tags: []string{databases}, Auth: db, Response: DatabaseSchema()},
		{Method: "PUT", Path: "/api/databases/:id", Summary: "Update a database", Tags: []string{databases}, Auth: db, Request: JSON(Obj("Update database", F("name", Str("Name")), F("version", Str("Version")), F("externalDns", Str("External DNS")), F("cpuLimit", Num("CPU limit")), F("memoryLimit", Int("Memory MB")), F("logicalReplication", Bool("Logical replication")), F("customArgs", Str("Custom args")))), Response: DatabaseSchema()},
		{Method: "DELETE", Path: "/api/databases/:id", Summary: "Delete a database", Tags: []string{databases}, Auth: db, Response: StatusSchema("deleted")},
		{Method: "POST", Path: "/api/databases/:id/start", Summary: "Start a database", Tags: []string{databases}, Auth: db, Response: DatabaseSchema()},
		{Method: "POST", Path: "/api/databases/:id/stop", Summary: "Stop a database", Tags: []string{databases}, Auth: db, Response: StatusSchema("stopped")},
		{Method: "POST", Path: "/api/databases/:id/restart", Summary: "Restart a database", Tags: []string{databases}, Auth: db, Response: DatabaseSchema()},
		{Method: "POST", Path: "/api/databases/:id/query", Summary: "Query a database", Tags: []string{databases}, Auth: db, Request: JSON(Obj("Query", RF("query", Str("SQL query")))), Response: queryResultSchema()},
		{Method: "POST", Path: "/api/databases/:id/import", Summary: "Import SQL into a database", Tags: []string{databases}, Auth: db, Request: JSON(Obj("Import", RF("sql", Str("SQL statements")))), Response: EmptySchema()},
		{Method: "POST", Path: "/api/databases/:id/credentials/reveal", Summary: "Reveal database credentials", Tags: []string{databases}, Auth: db, Response: Obj("Credentials", RF("password", Str("Password")))},
		{Method: "GET", Path: "/api/databases/:id/schemas", Summary: "List database tables", Tags: []string{databases}, Auth: db, Response: Arr("Tables", tableSchemaSchema())},
		{Method: "GET", Path: "/api/databases/:id/data/:table", Summary: "Get table rows", Tags: []string{databases}, Auth: db, Response: Any("Table rows"), Query: []Param{QueryParamInt("page", "Page"), QueryParamInt("limit", "Page size")}},
		{Method: "POST", Path: "/api/databases/:id/data/:table", Summary: "Insert a table row", Tags: []string{databases}, Auth: db, Request: JSON(Map("Row", Any("Value"))), Response: Any("Insert result")},
		{Method: "PUT", Path: "/api/databases/:id/data/:table", Summary: "Update table rows", Tags: []string{databases}, Auth: db, Request: JSON(Map("Row", Any("Value"))), Response: Any("Update result")},
		{Method: "DELETE", Path: "/api/databases/:id/data/:table", Summary: "Delete table rows", Tags: []string{databases}, Auth: db, Request: JSON(Map("Row", Any("Value"))), Response: Any("Delete result")},
		{Method: "GET", Path: "/api/databases/:id/backups", Summary: "List database backups", Tags: []string{databases}, Auth: db, Response: Arr("Records", Any("Backup record"))},
		{Method: "POST", Path: "/api/databases/:id/backups", Summary: "Trigger a database backup", Tags: []string{databases}, Auth: db, Response: Any("Backup record")},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/databases", Summary: "List cluster databases", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: Arr("Databases", clusterDataSchema())},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/review", Summary: "Review a cluster database", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(Obj("Review database", RF("action", Str("Action")), RF("spec", clusterDataSpecSchema()), F("sourceId", Str("Source ID")), F("restoreTime", Str("Restore time")), F("operatorVersion", Str("Operator version")))), Response: operationSchema()},
		{Method: "GET", Path: "/api/cluster-data-operations/:operationId", Summary: "Get a cluster database operation", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: operationSchema()},
		{Method: "POST", Path: "/api/cluster-data-operations/:operationId/apply", Summary: "Apply a cluster database operation", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Raw: true, Code: 202, Request: JSON(operationApplySchema()), Response: Obj("Accepted", RF("status", Str("Status")), RF("data", operationSchema()))},
		{Method: "POST", Path: "/api/cluster-data-operations/:operationId/cancel", Summary: "Cancel a cluster database operation", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/credentials", Summary: "Get database credentials", Tags: []string{clusterData}, Auth: db, Response: Obj("Credentials", RF("username", Str("Username")), RF("password", Str("Password")), RF("host", Str("Host")), RF("port", Int("Port")))},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/backups", Summary: "List cluster database backups", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: Any("Backups")},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/verify", Summary: "Verify a cluster database", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: EmptySchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/bindings/review", Summary: "Review a binding change", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(Obj("Review binding", RF("appId", Str("App ID")), RF("databaseId", Str("Database ID")))), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/bindings/apply", Summary: "Apply a binding change", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(operationApplySchema()), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/lifecycle/review", Summary: "Review a lifecycle change", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(Obj("Review lifecycle", RF("action", Str("Action")), RF("databaseId", Str("Database ID")), F("confirmName", Str("Name confirmation")), F("spec", clusterDataSpecSchema()))), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/lifecycle/apply", Summary: "Apply a lifecycle change", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(operationApplySchema()), Response: operationSchema()},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots", Summary: "List Redis snapshots", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: Arr("Snapshots", redisSnapshotSchema())},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots/review", Summary: "Review a Redis snapshot", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots/apply", Summary: "Apply a Redis snapshot", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(operationApplySchema()), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/restores/review", Summary: "Review a Redis restore", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(Obj("Review restore", RF("snapshotId", Str("Snapshot ID")), F("targetId", Str("Target ID")))), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/restores/apply", Summary: "Apply a Redis restore", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(operationApplySchema()), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/failover/review", Summary: "Review a Redis failover", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Response: operationSchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/databases/:databaseId/failover/apply", Summary: "Apply a Redis failover", Tags: []string{clusterData}, Auth: ProjectRole("admin"), Request: JSON(operationApplySchema()), Response: operationSchema()},
	}
}
