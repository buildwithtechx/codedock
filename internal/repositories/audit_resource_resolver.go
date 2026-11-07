package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"codedock/internal/models"
)

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// IsUUID checks whether the input string matches standard UUID v4 format.
func IsUUID(s string) bool {
	return uuidPattern.MatchString(strings.TrimSpace(s))
}

type resourceMeta struct {
	name         string
	resourceType string
}

func extractNameFromDetails(details string) (string, string) {
	d := strings.TrimSpace(details)
	if d == "" || !strings.HasPrefix(d, "{") {
		return "", ""
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(d), &m); err != nil {
		return "", ""
	}
	for _, key := range []string{"name", "appName", "projectName", "title", "domain", "hostname", "serverName", "slug"} {
		if val, ok := m[key]; ok {
			if s, ok := val.(string); ok && strings.TrimSpace(s) != "" && !IsUUID(s) {
				resType := ""
				if t, ok := m["type"].(string); ok {
					resType = t
				} else if t, ok := m["resourceType"].(string); ok {
					resType = t
				}
				return strings.TrimSpace(s), resType
			}
		}
	}
	for _, subKey := range []string{"after", "before"} {
		if sub, ok := m[subKey].(map[string]any); ok {
			for _, key := range []string{"name", "appName", "projectName", "title", "domain", "hostname", "serverName", "slug"} {
				if val, ok := sub[key]; ok {
					if s, ok := val.(string); ok && strings.TrimSpace(s) != "" && !IsUUID(s) {
						return strings.TrimSpace(s), ""
					}
				}
			}
		}
	}
	return "", ""
}

func resolveResourceNames(ctx context.Context, db *sql.DB, logs []models.AuditLog) {
	if len(logs) == 0 {
		return
	}
	unresolvedUUIDs := make(map[string]struct{})
	for i := range logs {
		if logs[i].ResourceName == "" && logs[i].Details != "" {
			name, resType := extractNameFromDetails(logs[i].Details)
			if name != "" {
				logs[i].ResourceName = name
			}
			if resType != "" && logs[i].ResourceType == "" {
				logs[i].ResourceType = resType
			}
		}
		if logs[i].ResourceName == "" && logs[i].Resource != "" {
			if !IsUUID(logs[i].Resource) {
				logs[i].ResourceName = logs[i].Resource
			} else {
				unresolvedUUIDs[logs[i].Resource] = struct{}{}
			}
		}
	}

	if len(unresolvedUUIDs) == 0 || db == nil {
		return
	}

	ids := make([]any, 0, len(unresolvedUUIDs))
	marks := make([]string, 0, len(unresolvedUUIDs))
	idx := 1
	for id := range unresolvedUUIDs {
		ids = append(ids, id)
		marks = append(marks, fmt.Sprintf("$%d", idx))
		idx++
	}

	inClause := strings.Join(marks, ", ")
	query := fmt.Sprintf(`
		SELECT id, name, type FROM (
			SELECT id, name, 'project' AS type FROM projects
			UNION ALL
			SELECT id, name, 'app' AS type FROM project_apps
			UNION ALL
			SELECT id, name, 'service' AS type FROM app_services
			UNION ALL
			SELECT id, name, 'server' AS type FROM servers
			UNION ALL
			SELECT id, name, 'database' AS type FROM databases
			UNION ALL
			SELECT id, name, 'organization' AS type FROM organizations
			UNION ALL
			SELECT id, name, 'user' AS type FROM users
		) resources
		WHERE id IN (%s)
	`, inClause)

	rows, err := db.QueryContext(ctx, query, ids...)
	if err != nil {
		return
	}
	defer rows.Close()

	resolvedMap := make(map[string]resourceMeta)
	for rows.Next() {
		var id, name, resType string
		if err := rows.Scan(&id, &name, &resType); err == nil && name != "" {
			resolvedMap[id] = resourceMeta{name: name, resourceType: resType}
		}
	}

	for i := range logs {
		if logs[i].ResourceName == "" {
			if info, ok := resolvedMap[logs[i].Resource]; ok {
				logs[i].ResourceName = info.name
				if logs[i].ResourceType == "" {
					logs[i].ResourceType = info.resourceType
				}
			}
		}
	}
}
