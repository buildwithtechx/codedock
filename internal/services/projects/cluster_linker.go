package projects

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/models"
	"fmt"
)

func buildClusterEnvVars(plan *models.ClusterDataPlan) map[string]string {
	vars := make(map[string]string)
	if plan == nil {
		return vars
	}
	namespace, name := kubernetes.DataIdentity(plan.Record.Spec)
	port := 6379
	username := "default"
	database := "0"
	host := name + "." + namespace + ".svc.cluster.local"
	if plan.Record.Spec.Engine == "postgres" {
		port = 5432
		username = "app"
		database = "app"
		host = name + "-rw." + namespace + ".svc.cluster.local"
	}
	prefix := plan.Record.Spec.Name
	if prefix == "" {
		prefix = plan.Record.ID
	}
	var connStr string
	if plan.Record.Spec.Engine == "postgres" {
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, plan.Password, host, port, database)
		vars[prefix+"_DATABASE_URL"] = connStr
		vars["DATABASE_URL"] = connStr
	} else {
		connStr = fmt.Sprintf("redis://:%s@%s:%d", plan.Password, host, port)
		vars[prefix+"_REDIS_URL"] = connStr
		vars["REDIS_URL"] = connStr
	}
	vars[prefix+"_HOST"] = host
	vars[prefix+"_PORT"] = fmt.Sprint(port)
	vars[prefix+"_USER"] = username
	vars[prefix+"_PASSWORD"] = plan.Password
	return vars
}
