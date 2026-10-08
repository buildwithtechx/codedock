package apischema

func analyticsSummarySchema() Schema {
	return Obj("Traffic summary",
		RF("projectId", Str("Project ID")),
		RF("from", Str("Window start")),
		RF("to", Str("Window end")),
		RF("requests", Int("Requests")),
		RF("bytes", Int64("Bytes")),
		F("errorRate", Num("Error rate")),
		F("avgDurationMs", Num("Average duration ms")),
	)
}

func analyticsOverviewSchema() Schema {
	return Obj("Traffic overview",
		RF("projectId", Str("Project ID")),
		RF("from", Str("Window start")),
		RF("to", Str("Window end")),
		F("series", Arr("Series", Obj("Point", RF("time", Str("Time")), RF("requests", Int("Requests")), RF("bytes", Int64("Bytes")), RF("errors", Int("Errors"))))),
		F("statuses", Arr("Statuses", Obj("Status", RF("status", Int("Status code")), RF("requests", Int("Requests"))))),
		F("topPaths", Arr("Top paths", Obj("Path", RF("path", Str("Path")), RF("requests", Int("Requests")), RF("bytes", Int64("Bytes"))))),
	)
}

func analyticsGeoSchema() Schema {
	return Obj("Traffic geography",
		RF("projectId", Str("Project ID")),
		RF("from", Str("Window start")),
		RF("to", Str("Window end")),
		F("countries", Arr("Countries", Obj("Country", RF("country", Str("Country")), RF("requests", Int("Requests")), RF("visitors", Int("Visitors")), RF("bytes", Int64("Bytes"))))),
	)
}

func analyticsDashboardSchema() Schema {
	return Obj("Analytics dashboard",
		RF("organizationId", Str("Organization ID")),
		RF("requests24h", Int("Requests in 24h")),
		RF("errorRate24h", Num("Error rate in 24h")),
		RF("deployments7d", Int("Deployments in 7d")),
		RF("deploySuccess", Num("Deploy success rate")),
		RF("openIssues", Int("Open attention issues")),
	)
}

func deploymentStatsSchema() Schema {
	return Obj("Deployment stats",
		RF("projectId", Str("Project ID")),
		RF("total", Int("Total deployments")),
		RF("succeeded", Int("Succeeded")),
		RF("failed", Int("Failed")),
		F("successRate", Num("Success rate")),
		F("avgDurationSeconds", Num("Average duration seconds")),
	)
}

func attentionIssueSchema() Schema {
	return Obj("Attention issue",
		RF("id", ID("Issue ID")),
		RF("organizationId", Str("Organization ID")),
		RF("kind", Str("Issue kind")),
		RF("subject", Str("Subject")),
		F("projectId", Str("Project ID")),
		F("serviceId", Str("Service ID")),
		RF("severity", Str("Severity")),
		RF("status", Str("Status")),
		RF("title", Str("Title")),
		F("detail", Str("Detail")),
		F("remediation", Str("Remediation")),
		F("action", Str("Action")),
		F("occurrences", Int("Occurrences")),
		F("firstSeen", Str("First seen")),
		F("lastSeen", Str("Last seen")),
		F("updatedAt", Str("Last update")),
	)
}

func analyticsOperations() []Operation {
	analytics := "analytics"
	attention := "attention"
	read := AuthOrg
	write := AuthOrg
	return []Operation{
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/summary", Summary: "Get the traffic summary", Tags: []string{analytics}, Auth: read, Response: analyticsSummarySchema(), Query: []Param{QueryParam("from", "Window start"), QueryParam("to", "Window end")}},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/overview", Summary: "Get the traffic overview", Tags: []string{analytics}, Auth: read, Response: analyticsOverviewSchema(), Query: []Param{QueryParam("from", "Window start"), QueryParam("to", "Window end")}},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/geo", Summary: "Get traffic geography", Tags: []string{analytics}, Auth: read, Response: analyticsGeoSchema(), Query: []Param{QueryParam("from", "Window start"), QueryParam("to", "Window end")}},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/paths", Summary: "Get path collection state", Tags: []string{analytics}, Auth: read, Response: Obj("Paths", RF("enabled", Bool("Path collection enabled")))},
		{Method: "PUT", Path: "/api/organizations/:id/projects/:projectId/analytics/paths", Summary: "Set path collection", Tags: []string{analytics}, Auth: write, Request: JSON(Obj("Set paths", RF("enabled", Bool("Collect paths")))), Response: Obj("Paths", RF("enabled", Bool("Path collection enabled")))},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/deployments", Summary: "Get deployment stats", Tags: []string{analytics}, Auth: read, Response: deploymentStatsSchema(), Query: []Param{QueryParamInt("days", "Lookback days")}},
		{Method: "GET", Path: "/api/organizations/:id/analytics/dashboard", Summary: "Get the analytics dashboard", Tags: []string{analytics}, Auth: read, Response: analyticsDashboardSchema()},
		{Method: "GET", Path: "/api/organizations/:id/analytics/periods", Summary: "List analytics periods", Tags: []string{analytics}, Auth: read, Response: Arr("Periods", Str("Period"))},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/usage", Summary: "Get service usage", Tags: []string{analytics}, Auth: read, Response: Arr("Usage", Any("Service usage"))},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/usage/history", Summary: "Get usage history", Tags: []string{analytics}, Auth: read, Response: Any("Usage history"), Query: []Param{QueryParam("from", "Window start"), QueryParam("to", "Window end")}},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/usage/stream", Summary: "Stream usage", Tags: []string{analytics}, Auth: read, Raw: true, Content: ContentSSE, Response: Str("Usage events")},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/container", Summary: "Get container analytics", Tags: []string{analytics}, Auth: read, Response: Any("Container analytics"), Query: []Param{QueryParam("containerId", "Container ID")}},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/live", Summary: "Get live traffic", Tags: []string{analytics}, Auth: read, Response: analyticsOverviewSchema()},
		{Method: "GET", Path: "/api/organizations/:id/projects/:projectId/analytics/resources", Summary: "Get resource analytics", Tags: []string{analytics}, Auth: read, Response: Any("Resource analytics")},
		{Method: "GET", Path: "/api/organizations/:id/attention", Summary: "List attention issues", Tags: []string{attention}, Auth: read, Response: Arr("Issues", attentionIssueSchema()), Query: []Param{QueryParam("status", "Filter by status")}},
		{Method: "POST", Path: "/api/organizations/:id/attention/evaluate", Summary: "Evaluate attention issues", Tags: []string{attention}, Auth: write, Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/attention/:issueId/ack", Summary: "Acknowledge an issue", Tags: []string{attention}, Auth: write, Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/attention/:issueId/resolve", Summary: "Resolve an issue", Tags: []string{attention}, Auth: write, Response: EmptySchema()},
		{Method: "POST", Path: "/api/organizations/:id/attention/:issueId/act", Summary: "Act on an issue", Tags: []string{attention}, Auth: write, Request: JSON(Obj("Act", RF("action", Str("Action")), F("params", Map("Params", Str("Value"))))), Response: Obj("Outcome", RF("outcome", Str("Outcome")))},
	}
}
