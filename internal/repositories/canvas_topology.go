package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
)

type topologyReader interface {
	SelectContext(context.Context, any, string, ...any) error
}

func topologyRevision(ctx context.Context, reader topologyReader, environment string) (string, error) {
	values := []string{}
	err := reader.SelectContext(ctx, &values, `SELECT value FROM (
 SELECT 'app:'||id||':'||updated_at||':'||COALESCE(domain,'') AS value FROM app_services WHERE environment_id=$1
 UNION ALL SELECT 'db:'||id||':'||updated_at FROM databases WHERE environment_id=$2
 UNION ALL SELECT 'clusterdb:'||id||':'||updated_at||':'||status FROM cluster_databases WHERE id IN (SELECT id FROM cluster_databases WHERE project_id=(SELECT project_id FROM environments WHERE id=$3))
 UNION ALL SELECT 'var:'||v.id||':'||v.key||':'||v.value FROM service_vars v JOIN app_services a ON a.id=v.service_id WHERE a.environment_id=$4
 UNION ALL SELECT 'domain:'||d.id||':'||d.hostname||':'||COALESCE(d.service_id,'') FROM domains d JOIN app_services a ON a.id=d.service_id WHERE a.environment_id=$5
 UNION ALL SELECT 'dep:'||source||':'||target FROM topology_dependencies WHERE environment_id=$6
 ) ORDER BY value`, environment, environment, environment, environment, environment, environment)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func (r *CanvasRepo) projectTopology(ctx context.Context, canvas *models.EnvironmentCanvas) error {
	canvas.Nodes = []models.CanvasNode{}
	canvas.Edges = []models.CanvasEdge{}
	names := map[string]string{}
	for i, app := range canvas.Apps {
		id := "app-" + app.ID
		names[app.Name] = id
		names[app.ID] = id
		names[strings.ToLower(strings.ReplaceAll(app.Name, "_", "-"))] = id
		canvas.Nodes = append(canvas.Nodes, models.CanvasNode{ID: id, Type: "appService", Data: map[string]any{"name": app.Name, "status": app.Status, "serviceId": app.ID}, Pos: models.CanvasPosition{X: 100 + float64(i)*240, Y: 50}})
	}
	for i, database := range canvas.Databases {
		id := "db-" + database.ID
		names[database.Name] = id
		names[database.ID] = id
		names[strings.ToLower(strings.ReplaceAll(database.Name, "_", "-"))] = id
		canvas.Nodes = append(canvas.Nodes, models.CanvasNode{ID: id, Type: "database", Data: map[string]any{"name": database.Name, "status": database.Status, "engine": database.Engine, "databaseId": database.ID}, Pos: models.CanvasPosition{X: 100 + float64(i)*240, Y: 240}})
	}
	for i, record := range canvas.ClusterDatabases {
		id := "clusterdb-" + record.ID
		names[record.Spec.Name] = id
		names[record.ID] = id
		names[strings.ToLower(strings.ReplaceAll(record.Spec.Name, "_", "-"))] = id
		canvas.Nodes = append(canvas.Nodes, models.CanvasNode{ID: id, Type: "clusterDatabase", Data: map[string]any{"name": record.Spec.Name, "status": record.Status, "engine": record.Spec.Engine, "clusterDatabaseId": record.ID, "clusterId": record.ClusterID}, Pos: models.CanvasPosition{X: 100 + float64(i)*240, Y: 380}})
	}
	var dependencies []struct {
		Source string `db:"source"`
		Target string `db:"target"`
	}
	if err := r.db.SelectContext(ctx, &dependencies, `SELECT source,target FROM topology_dependencies WHERE environment_id=$1`, canvas.Environment.ID); err != nil {
		return err
	}
	for _, dependency := range dependencies {
		canvas.Edges = append(canvas.Edges, models.CanvasEdge{ID: "dependency:" + dependency.Source + ":" + dependency.Target, Source: dependency.Source, Target: dependency.Target, Kind: "dependency", Label: "depends on"})
	}
	var variables []struct {
		ID        string `db:"id"`
		ServiceID string `db:"service_id"`
		Key       string `db:"key"`
		Value     string `db:"value"`
	}
	if err := r.db.SelectContext(ctx, &variables, `SELECT v.id,v.service_id,v.key,v.value FROM service_vars v JOIN app_services a ON a.id=v.service_id WHERE a.environment_id=$1`, canvas.Environment.ID); err != nil {
		return err
	}
	pattern := regexp.MustCompile(`\$\{([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_]+)\}`)
	for _, variable := range variables {
		for index, match := range pattern.FindAllStringSubmatch(variable.Value, -1) {
			if source := names[match[1]]; source != "" {
				canvas.Edges = append(canvas.Edges, models.CanvasEdge{ID: fmt.Sprintf("binding:%s:%d", variable.ID, index), Source: source, Target: "app-" + variable.ServiceID, Kind: "binding", Label: variable.Key + " ← " + match[2]})
			}
		}
	}
	var domains []struct {
		ID        string `db:"id"`
		ServiceID string `db:"service_id"`
		Hostname  string `db:"hostname"`
	}
	if err := r.db.SelectContext(ctx, &domains, `SELECT d.id,d.service_id,d.hostname FROM domains d JOIN app_services a ON a.id=d.service_id WHERE a.environment_id=$1`, canvas.Environment.ID); err != nil {
		return err
	}
	for i, domain := range domains {
		id := "domain-" + domain.ID
		canvas.Nodes = append(canvas.Nodes, models.CanvasNode{ID: id, Type: "default", Data: map[string]any{"label": domain.Hostname, "name": domain.Hostname, "serviceId": domain.ServiceID}, Pos: models.CanvasPosition{X: 100 + float64(i)*240, Y: -140}})
		canvas.Edges = append(canvas.Edges, models.CanvasEdge{ID: "routing:" + domain.ID, Source: id, Target: "app-" + domain.ServiceID, Kind: "routing", Label: "routes to"})
	}
	if err := r.projectStacks(ctx, canvas); err != nil {
		return err
	}
	revision, err := topologyRevision(ctx, r.db, canvas.Environment.ID)
	canvas.Revision = revision
	return err
}

func (r *CanvasRepo) ApplyTopology(ctx context.Context, environment string, request models.TopologyApplyRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	revision, err := topologyRevision(ctx, tx, environment)
	if err != nil {
		return err
	}
	if request.Revision != revision {
		return fmt.Errorf("topology changed since review; reload before applying")
	}
	nodes := []string{}
	if err := tx.SelectContext(ctx, &nodes, `SELECT 'app-'||id FROM app_services WHERE environment_id=$1 UNION ALL SELECT 'db-'||id FROM databases WHERE environment_id=$2 UNION ALL SELECT 'clusterdb-'||cd.id FROM cluster_databases cd JOIN environments e ON e.project_id=cd.project_id WHERE e.id=$3`, environment, environment, environment); err != nil {
		return err
	}
	valid := map[string]bool{}
	for _, node := range nodes {
		valid[node] = true
	}
	if err := validateDependencies(request.Dependencies, valid); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM topology_dependencies WHERE environment_id=$1`, environment); err != nil {
		return err
	}
	for _, edge := range request.Dependencies {
		if _, err := tx.ExecContext(ctx, `INSERT INTO topology_dependencies(environment_id,source,target) VALUES($1,$2,$3)`, environment, edge.Source, edge.Target); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func validateDependencies(edges []models.CanvasEdge, nodes map[string]bool) error {
	if len(edges) > 500 {
		return fmt.Errorf("too many dependencies")
	}
	graph := map[string][]string{}
	seen := map[string]bool{}
	for _, edge := range edges {
		if edge.Kind != "dependency" || !nodes[edge.Source] || !nodes[edge.Target] || !strings.HasPrefix(edge.Target, "app-") {
			return fmt.Errorf("dependencies must connect resources in this environment to an application")
		}
		key := edge.Source + ":" + edge.Target
		if seen[key] {
			return fmt.Errorf("duplicate dependency")
		}
		seen[key] = true
		graph[edge.Source] = append(graph[edge.Source], edge.Target)
	}
	active, visited := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(node string) error {
		if active[node] {
			return fmt.Errorf("dependency cycle detected")
		}
		if visited[node] {
			return nil
		}
		active[node] = true
		for _, target := range graph[node] {
			if err := visit(target); err != nil {
				return err
			}
		}
		active[node], visited[node] = false, true
		return nil
	}
	for node := range graph {
		if err := visit(node); err != nil {
			return err
		}
	}
	return nil
}

var _ topologyReader = (*sqlx.Tx)(nil)
