package apischema

func clusterSchema() Schema {
	return Obj("Cluster",
		RF("id", ID("Cluster ID")),
		RF("projectId", Str("Project ID")),
		RF("name", Str("Cluster name")),
		F("version", Str("Kubernetes version")),
		F("controls", Int("Control-plane nodes")),
		F("nodes", Arr("Nodes", Obj("Node", RF("serverId", Str("Server ID")), F("privateIp", Str("Private IP")), F("interface", Str("Interface")), F("fingerprint", Str("Fingerprint"))))),
		F("revision", Int("Revision")),
		F("status", Str("Status")),
		F("error", Str("Error")),
		F("updatedAt", Str("Last update")),
	)
}

func clusterOperations() []Operation {
	clusters := "clusters"
	admin := ProjectRole("admin")
	return []Operation{
		{Method: "GET", Path: "/api/projects/:id/runtime-clusters", Summary: "List runtime clusters", Tags: []string{clusters}, Auth: admin, Response: Arr("Clusters", clusterSchema())},
		{Method: "GET", Path: "/api/projects/:id/clusters", Summary: "List clusters", Tags: []string{clusters}, Auth: admin, Response: Arr("Clusters", clusterSchema())},
		{Method: "POST", Path: "/api/projects/:id/clusters/review", Summary: "Review a cluster change", Tags: []string{clusters}, Auth: admin, Request: JSON(Obj("Review cluster", RF("cluster", clusterSchema()), RF("action", Str("Action")))), Response: operationSchema()},
		{Method: "GET", Path: "/api/cluster-operations/:operationId", Summary: "Get a cluster operation", Tags: []string{clusters}, Auth: admin, Response: operationSchema()},
		{Method: "POST", Path: "/api/cluster-operations/:operationId/apply", Summary: "Apply a cluster operation", Tags: []string{clusters}, Auth: admin, Raw: true, Code: 202, Request: JSON(operationApplySchema()), Response: Obj("Accepted", RF("status", Str("Status")), RF("data", operationSchema()))},
		{Method: "POST", Path: "/api/cluster-operations/:operationId/cancel", Summary: "Cancel a cluster operation", Tags: []string{clusters}, Auth: admin, Response: EmptySchema()},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/preflight", Summary: "Run cluster preflight checks", Tags: []string{clusters}, Auth: admin, Response: Any("Preflight result")},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/storage", Summary: "Get cluster storage", Tags: []string{clusters}, Auth: admin, Response: Any("Storage inventory")},
		{Method: "GET", Path: "/api/projects/:id/clusters/:clusterId/storage/health", Summary: "Get storage health", Tags: []string{clusters}, Auth: admin, Response: EmptySchema()},
		{Method: "POST", Path: "/api/projects/:id/clusters/:clusterId/storage/review", Summary: "Review storage setup", Tags: []string{clusters}, Auth: admin, Request: JSON(Obj("Review storage", RF("version", Str("Longhorn version")))), Response: operationSchema()},
		{Method: "POST", Path: "/api/cluster-storage-operations/:operationId/apply", Summary: "Apply a storage operation", Tags: []string{clusters}, Auth: admin, Request: JSON(operationApplySchema()), Response: operationSchema()},
	}
}
