package kubernetes

import (
	"codedock/internal/models"
	"fmt"
	"strings"
)

func redisObjects(spec models.ClusterDataSpec, metadata func(string) map[string]any, labels map[string]string) []map[string]any {
	if spec.Shards > 1 {
		return redisClusterObjects(spec, metadata, labels)
	}
	_, name := DataIdentity(spec)
	headless := name + "-nodes"
	selector := map[string]string{"codedock.run/database": spec.ID}
	service := map[string]any{"apiVersion": "v1", "kind": "Service", "metadata": metadata(headless), "spec": map[string]any{"clusterIP": "None", "publishNotReadyAddresses": true, "selector": selector, "ports": []map[string]any{{"name": "redis", "port": 6379, "targetPort": 6379}}}}
	primary := map[string]any{"apiVersion": "v1", "kind": "Service", "metadata": metadata(name), "spec": map[string]any{"selector": map[string]string{"statefulset.kubernetes.io/pod-name": name + "-0"}, "ports": []map[string]any{{"name": "redis", "port": 6379, "targetPort": 6379}}}}
	script := `set -eu
set -- redis-server --appendonly yes --appendfsync everysec --requirepass "$REDIS_PASSWORD" --masterauth "$REDIS_PASSWORD" --protected-mode yes --bind 0.0.0.0 --port 6379
if test "$(hostname)" != ` + name + `-0; then set -- "$@" --replicaof ` + name + `-0.` + headless + ` 6379; fi
exec "$@"
`
	health := `set -eu
redis-cli --raw PING | grep -x PONG >/dev/null
if test "$(hostname)" != ` + name + `-0; then redis-cli --raw INFO replication | tr -d '\r' | grep -x 'master_link_status:up' >/dev/null; redis-cli --raw INFO replication | tr -d '\r' | grep -x 'master_sync_in_progress:0' >/dev/null; fi
`
	password := map[string]any{"valueFrom": map[string]any{"secretKeyRef": map[string]string{"name": name + "-auth", "key": "password"}}}
	env := []map[string]any{{"name": "REDIS_PASSWORD", "valueFrom": password["valueFrom"]}, {"name": "REDISCLI_AUTH", "valueFrom": password["valueFrom"]}}
	pod := map[string]any{"securityContext": map[string]any{"runAsUser": 999, "runAsNonRoot": true, "fsGroup": 999}, "containers": []map[string]any{{"name": "redis", "image": spec.Image, "command": []string{"sh", "-ec", script}, "env": env, "ports": []map[string]any{{"containerPort": 6379}}, "volumeMounts": []map[string]string{{"name": "data", "mountPath": "/data"}}, "readinessProbe": map[string]any{"exec": map[string]any{"command": []string{"sh", "-ec", health}}, "initialDelaySeconds": 5, "periodSeconds": 5, "timeoutSeconds": 3}, "livenessProbe": map[string]any{"exec": map[string]any{"command": []string{"redis-cli", "PING"}}, "initialDelaySeconds": 30, "periodSeconds": 10, "timeoutSeconds": 3}}}, "terminationGracePeriodSeconds": 60}
	if spec.Instances > 1 {
		pod["affinity"] = map[string]any{"podAntiAffinity": map[string]any{"requiredDuringSchedulingIgnoredDuringExecution": []map[string]any{{"labelSelector": map[string]any{"matchLabels": selector}, "topologyKey": "kubernetes.io/hostname"}}}}
	}
	claimSpec := map[string]any{"accessModes": []string{"ReadWriteOnce"}, "resources": map[string]any{"requests": map[string]string{"storage": fmt.Sprintf("%dGi", spec.StorageGiB)}}}
	if spec.StorageClass != "" {
		claimSpec["storageClassName"] = spec.StorageClass
	}
	stateful := map[string]any{"apiVersion": "apps/v1", "kind": "StatefulSet", "metadata": metadata(name), "spec": map[string]any{"serviceName": headless, "replicas": spec.Instances, "podManagementPolicy": "OrderedReady", "selector": map[string]any{"matchLabels": selector}, "template": map[string]any{"metadata": map[string]any{"labels": labels}, "spec": pod}, "volumeClaimTemplates": []map[string]any{{"metadata": map[string]any{"name": "data", "labels": labels}, "spec": claimSpec}}}}
	return []map[string]any{service, primary, stateful}
}

func redisClusterObjects(spec models.ClusterDataSpec, metadata func(string) map[string]any, labels map[string]string) []map[string]any {
	_, base := DataIdentity(spec)
	objects := []map[string]any{}
	password := map[string]any{"valueFrom": map[string]any{"secretKeyRef": map[string]string{"name": base + "-auth", "key": "password"}}}
	for shard := 0; shard < spec.Shards; shard++ {
		name := fmt.Sprintf("%s-s%d", base, shard)
		shardLabels := map[string]string{}
		for key, value := range labels {
			shardLabels[key] = value
		}
		shardLabels["codedock.run/shard"] = fmt.Sprint(shard)
		selector := map[string]string{"codedock.run/database": spec.ID, "codedock.run/shard": fmt.Sprint(shard)}
		headless := name + "-nodes"
		script := `set -eu
set -- redis-server --appendonly yes --appendfsync everysec --requirepass "$REDIS_PASSWORD" --masterauth "$REDIS_PASSWORD" --protected-mode yes --bind 0.0.0.0 --port 6379 --cluster-enabled yes --cluster-config-file /data/nodes.conf --cluster-node-timeout 5000
if test "$(hostname)" != ` + name + `-0; then set -- "$@" --replicaof ` + name + `-0.` + headless + ` 6379; fi
exec "$@"
`
		health := `set -eu
redis-cli --raw PING | grep -x PONG >/dev/null
redis-cli --raw CLUSTER INFO | tr -d '\r' | grep -E '^cluster_state:ok$' >/dev/null
`
		env := []map[string]any{{"name": "REDIS_PASSWORD", "valueFrom": password["valueFrom"]}, {"name": "REDISCLI_AUTH", "valueFrom": password["valueFrom"]}}
		pod := map[string]any{"securityContext": map[string]any{"runAsUser": 999, "runAsNonRoot": true, "fsGroup": 999}, "containers": []map[string]any{{"name": "redis", "image": spec.Image, "command": []string{"sh", "-ec", script}, "env": env, "ports": []map[string]any{{"containerPort": 6379}, {"containerPort": 16379}}, "volumeMounts": []map[string]string{{"name": "data", "mountPath": "/data"}}, "readinessProbe": map[string]any{"exec": map[string]any{"command": []string{"sh", "-ec", health}}, "initialDelaySeconds": 10, "periodSeconds": 10, "timeoutSeconds": 5}, "livenessProbe": map[string]any{"exec": map[string]any{"command": []string{"redis-cli", "PING"}}, "initialDelaySeconds": 30, "periodSeconds": 10, "timeoutSeconds": 3}}}, "terminationGracePeriodSeconds": 60}
		if spec.Instances > 1 {
			pod["affinity"] = map[string]any{"podAntiAffinity": map[string]any{"requiredDuringSchedulingIgnoredDuringExecution": []map[string]any{{"labelSelector": map[string]any{"matchLabels": selector}, "topologyKey": "kubernetes.io/hostname"}}}}
		}
		claimSpec := map[string]any{"accessModes": []string{"ReadWriteOnce"}, "resources": map[string]any{"requests": map[string]string{"storage": fmt.Sprintf("%dGi", spec.StorageGiB)}}}
		if spec.StorageClass != "" {
			claimSpec["storageClassName"] = spec.StorageClass
		}
		objects = append(objects,
			map[string]any{"apiVersion": "v1", "kind": "Service", "metadata": metadata(headless), "spec": map[string]any{"clusterIP": "None", "publishNotReadyAddresses": true, "selector": selector, "ports": []map[string]any{{"name": "redis", "port": 6379, "targetPort": 6379}, {"name": "bus", "port": 16379, "targetPort": 16379}}}},
			map[string]any{"apiVersion": "apps/v1", "kind": "StatefulSet", "metadata": metadata(name), "spec": map[string]any{"serviceName": headless, "replicas": spec.Instances, "podManagementPolicy": "OrderedReady", "selector": map[string]any{"matchLabels": selector}, "template": map[string]any{"metadata": map[string]any{"labels": shardLabels}, "spec": pod}, "volumeClaimTemplates": []map[string]any{{"metadata": map[string]any{"name": "data", "labels": shardLabels}, "spec": claimSpec}}}},
		)
	}
	primaries := []string{}
	for shard := 0; shard < spec.Shards; shard++ {
		primaries = append(primaries, fmt.Sprintf("%s-s%d-0.%s-s%d-nodes:6379", base, shard, base, shard))
	}
	objects = append(objects, map[string]any{"apiVersion": "batch/v1", "kind": "Job", "metadata": metadata(base + "-bootstrap"), "spec": map[string]any{"backoffLimit": 0, "template": map[string]any{"metadata": map[string]any{"labels": labels}, "spec": map[string]any{"restartPolicy": "Never", "containers": []map[string]any{{"name": "bootstrap", "image": spec.Image, "command": []string{"sh", "-ec", "redis-cli -a \"$REDISCLI_AUTH\" --cluster create " + strings.Join(primaries, " ") + " --cluster-replicas 0 --cluster-yes"}, "env": []map[string]any{{"name": "REDISCLI_AUTH", "valueFrom": password["valueFrom"]}}}}}}}})
	return objects
}
