package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"regexp"
)

func ValidateDataSpec(spec models.ClusterDataSpec, nodes int) error {
	for _, id := range []string{spec.ID, spec.EnvironmentID} {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("database and environment identities must be UUIDs")
		}
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9-]{0,39}$`).MatchString(spec.Name) {
		return fmt.Errorf("database name must use lowercase letters, numbers and hyphens")
	}
	if spec.Instances != 1 && spec.Instances != 3 {
		return fmt.Errorf("select standalone or three database instances")
	}
	if nodes < spec.Instances {
		return fmt.Errorf("replicated databases require three distinct cluster nodes")
	}
	if spec.StorageGiB < 1 || spec.StorageGiB > 16384 {
		return fmt.Errorf("select persistent storage from 1 to 16384 GiB")
	}
	if spec.StorageClass != "" && !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,62}$`).MatchString(spec.StorageClass) {
		return fmt.Errorf("invalid storage class")
	}
	switch spec.Engine {
	case "postgres":
		if !regexp.MustCompile(`^ghcr\.io/cloudnative-pg/postgresql:[0-9]+\.[0-9]+[-a-zA-Z0-9.]*$`).MatchString(spec.Image) {
			return fmt.Errorf("select an exact CloudNativePG PostgreSQL image tag")
		}
		if spec.Instances == 3 && !spec.Synchronous {
			return fmt.Errorf("three PostgreSQL instances require synchronous durability")
		}
	case "redis":
		if !regexp.MustCompile(`^redis:[0-9]+\.[0-9]+\.[0-9]+[-a-zA-Z0-9.]*$`).MatchString(spec.Image) {
			return fmt.Errorf("select an exact Redis image tag")
		}
		if spec.S3DestinationID != "" {
			return fmt.Errorf("Redis persistence uses AOF/PVCs; PostgreSQL object-store backups are a separate capability")
		}
	default:
		return fmt.Errorf("cluster databases support PostgreSQL and Redis")
	}
	return nil
}
func DataIdentity(spec models.ClusterDataSpec) (string, string) {
	return "codedock-" + spec.EnvironmentID, "db-" + spec.ID
}
func DataManifest(plan *models.ClusterDataPlan, destination, source *models.S3Destination, sourceName string) (string, error) {
	spec := plan.Record.Spec
	namespace, name := DataIdentity(spec)
	labels := map[string]string{"codedock.run/project": plan.Record.ProjectID, "codedock.run/database": spec.ID, "codedock.run/environment": spec.EnvironmentID}
	metadata := func(resource string) map[string]any {
		return map[string]any{"name": resource, "namespace": namespace, "labels": labels}
	}
	objects := []map[string]any{{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": namespace, "labels": map[string]string{"codedock.run/project": plan.Record.ProjectID, "codedock.run/environment": spec.EnvironmentID}}}, {"apiVersion": "v1", "kind": "Secret", "metadata": metadata(name + "-auth"), "type": "kubernetes.io/basic-auth", "stringData": map[string]string{"username": "app", "password": plan.Password}}}
	if spec.Engine == "redis" {
		objects = append(objects, redisObjects(spec, metadata, labels)...)
	} else {
		postgres, extra, err := postgresObjects(plan, metadata, destination, source, sourceName)
		if err != nil {
			return "", err
		}
		objects = append(objects, extra...)
		objects = append(objects, postgres)
	}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": objects})
	return string(data), err
}
